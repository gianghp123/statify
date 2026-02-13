package processor

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"log/slog"
	"mime"
	"path/filepath"
	"strings"

	"github.com/gianghp/statify/internal/database/models"
	"github.com/gianghp/statify/internal/modules/deployment/repository"
	"github.com/gianghp/statify/internal/storage/minio"
	"github.com/gianghp/statify/internal/utils"
	minioGo "github.com/minio/minio-go/v7"
)

type FileProcessor struct {
	minioClient minio.Interface
	repo        repository.IDeploymentRepository
	bucketName  string
}

func NewFileProcessor(minioClient minio.Interface, repo repository.IDeploymentRepository, bucketName string) *FileProcessor {
	return &FileProcessor{
		minioClient: minioClient,
		repo:        repo,
		bucketName:  bucketName,
	}
}

func (s *FileProcessor) ProcessDeploymentFiles(ctx context.Context, deployment *models.Deployment) error {
	var uploadedObjects []string

	objectName := fmt.Sprintf("%s/temp", deployment.OutputPrefix)

	stat, err := s.minioClient.StatObject(ctx, s.bucketName, objectName, minioGo.StatObjectOptions{})
	if err != nil {
		return err
	}

	object, err := s.minioClient.GetObject(ctx, s.bucketName, objectName, minioGo.GetObjectOptions{})
	if err != nil {
		return err
	}

	defer object.Close()

	zipReader, err := zip.NewReader(object, stat.Size)
	if err != nil {
		return err
	}

	// Validate zip archive
	if err := utils.ValidateZipArchive(zipReader.File); err != nil {
		return err
	}

	// Validate file types
	if err := utils.ValidateStaticFileTypes(zipReader.File); err != nil {
		return err
	}

	// Validate entrypoint
	if err := utils.ValidateEntrypoint(zipReader.File); err != nil {
		return err
	}

	for _, file := range zipReader.File {
		if file.FileInfo().IsDir() {
			continue
		}

		err = s.uploadFileToMinio(ctx, deployment.OutputPrefix, file)
		if err != nil {
			slog.Error("Upload failed", "file", file.Name, "error", err)
			s.rollbackMinioUploads(ctx, uploadedObjects)
			return err
		}
		fullKey := fmt.Sprintf("%s/%s", deployment.OutputPrefix, file.Name)
		uploadedObjects = append(uploadedObjects, fullKey)
	}

	_ = s.minioClient.RemoveObject(ctx, s.bucketName, objectName, minioGo.RemoveObjectOptions{})

	return nil
}

func (s *FileProcessor) uploadFileToMinio(ctx context.Context, prefix string, zf *zip.File) error {
	rc, err := zf.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	// Detect Content-Type
	contentType := mime.TypeByExtension(filepath.Ext(zf.Name))
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	limitReader := io.LimitReader(rc, int64(zf.UncompressedSize64))

	objectName := fmt.Sprintf("%s/%s", prefix, zf.Name)

	_, err = s.minioClient.PutObject(ctx, s.bucketName, objectName, limitReader, int64(zf.UncompressedSize64), minioGo.PutObjectOptions{
		ContentType: contentType,
	})

	return err
}

// Helper: Rollback (Compensating Transaction)
func (s *FileProcessor) rollbackMinioUploads(ctx context.Context, objects []string) {
	// Create a channel for objects to delete (MinIO optimized batch deletion)
	objectsCh := make(chan minioGo.ObjectInfo)

	go func() {
		defer close(objectsCh)
		for _, objName := range objects {
			objectsCh <- minioGo.ObjectInfo{Key: objName}
		}
	}()

	errorCh := s.minioClient.RemoveObjects(ctx, s.bucketName, objectsCh, minioGo.RemoveObjectsOptions{})

	for err := range errorCh {
		slog.Error("CRITICAL: Failed to delete object during rollback", "object", err.ObjectName, "error", err.Err)
	}
}

func (s *FileProcessor) DeleteMinioFolder(ctx context.Context, prefix string) error {
	if prefix == "" || prefix == "/" {
		return fmt.Errorf("refusing to delete unsafe prefix: %s", prefix)
	}
	objectsCh := make(chan minioGo.ObjectInfo)

	go func() {
		defer close(objectsCh)
		for object := range s.minioClient.ListObjects(ctx, s.bucketName, minioGo.ListObjectsOptions{
			Prefix:    prefix,
			Recursive: true,
		}) {
			if object.Err != nil {
				slog.Error("Error listing object for deletion", "error", object.Err)
				continue
			}
			objectsCh <- object
		}
	}()

	errorCh := s.minioClient.RemoveObjects(ctx, s.bucketName, objectsCh, minioGo.RemoveObjectsOptions{})

	var deleteErrors []string
	for err := range errorCh {
		deleteErrors = append(deleteErrors, fmt.Sprintf("%s: %v", err.ObjectName, err.Err))
	}

	if len(deleteErrors) > 0 {
		return fmt.Errorf("failed to delete some objects: %s", strings.Join(deleteErrors, ", "))
	}

	return nil
}
