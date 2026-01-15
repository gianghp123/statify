package service

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gianghp/statify/internal/core"
	"github.com/gianghp/statify/internal/core/enums"
	coreRepo "github.com/gianghp/statify/internal/core/repository"
	"github.com/gianghp/statify/internal/database/models"
	"github.com/gianghp/statify/internal/modules/deployment/dtos/request"
	"github.com/gianghp/statify/internal/modules/deployment/dtos/response"
	"github.com/gianghp/statify/internal/modules/deployment/repository"
	projectRepo "github.com/gianghp/statify/internal/modules/project/repository"
	"github.com/gianghp/statify/internal/storage/minio"
	"github.com/gianghp/statify/internal/utils"
	minioGo "github.com/minio/minio-go/v7"
)

type DeploymentService struct {
	repo        repository.IDeploymentRepository
	projectRepo projectRepo.IProjectRepository
	minioClient minio.Interface
}

func NewDeploymentService(repo repository.IDeploymentRepository, projectRepo projectRepo.IProjectRepository, minioClient minio.Interface) *DeploymentService {
	return &DeploymentService{repo: repo, projectRepo: projectRepo, minioClient: minioClient}
}

func (s *DeploymentService) GetGlobalDeploymentHistory(ctx context.Context, userId uint, page int, limit int) (coreRepo.PaginatedEntities[*response.DeploymentDto], error) {
	result := coreRepo.PaginatedEntities[*response.DeploymentDto]{
		Entities:   []*response.DeploymentDto{},
		Pagination: coreRepo.Pagination{Page: page, Limit: limit},
	}

	deployments, err := s.repo.FindAllByUserID(ctx, userId, page, limit)
	if err != nil {
		return result, core.ParseDatabaseError(err)
	}

	deploymentDtos, err := utils.EntitiesToDto[*response.DeploymentDto](deployments.Entities)
	if err != nil {
		return result, core.InternalError()
	}

	return coreRepo.PaginatedEntities[*response.DeploymentDto]{
		Entities:   deploymentDtos,
		Pagination: deployments.Pagination,
	}, nil
}

func (s *DeploymentService) GetHistory(ctx context.Context, userID uint, projectID uint, page int, limit int) (coreRepo.PaginatedEntities[*response.DeploymentDto], error) {
	result := coreRepo.PaginatedEntities[*response.DeploymentDto]{
		Entities:   []*response.DeploymentDto{},
		Pagination: coreRepo.Pagination{Page: page, Limit: limit},
	}

	project, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return result, core.ParseDatabaseError(err)
	}

	if project == nil {
		return result, core.NotFoundError()
	}

	if project.UserID != userID {
		return result, core.ForbiddenError()
	}

	deployments, err := s.repo.FindAllByProjectID(ctx, projectID, page, limit)
	if err != nil {
		return result, core.ParseDatabaseError(err)
	}

	deploymentDtos, err := utils.EntitiesToDto[*response.DeploymentDto](deployments.Entities)
	if err != nil {
		return result, core.InternalError()
	}

	return coreRepo.PaginatedEntities[*response.DeploymentDto]{
		Entities:   deploymentDtos,
		Pagination: deployments.Pagination,
	}, nil
}

func (s *DeploymentService) GetDeploymentByID(ctx context.Context, userID uint, id uint) (*response.DeploymentDto, error) {
	deployment, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, core.ParseDatabaseError(err)
	}

	if deployment == nil {
		return nil, core.NotFoundError()
	}

	project, err := s.projectRepo.FindByID(ctx, deployment.ProjectID)
	if err != nil {
		return nil, core.ParseDatabaseError(err)
	}

	if project == nil {
		return nil, core.NotFoundError()
	}

	if project.UserID != userID {
		return nil, core.ForbiddenError()
	}

	deploymentDto, err := utils.EntityToDto[*response.DeploymentDto](deployment)
	if err != nil {
		return nil, core.InternalError()
	}

	return deploymentDto, nil
}

func (s *DeploymentService) CreateDeployment(ctx context.Context, userID uint, projectID uint, req *request.CreateDeploymentRequest) (*response.DeploymentDto, error) {
	// 1. Authorization & Validation
	project, err := s.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return nil, core.ParseDatabaseError(err)
	}
	if project == nil {
		return nil, core.NotFoundError()
	}
	if project.UserID != userID {
		return nil, core.ForbiddenError()
	}

	// 2. Zip Safety Analysis (Prevent Disk Fill / Zip Bomb)
	multipartFile, err := req.File.Open()
	if err != nil {
		return nil, core.BadRequestError("Could not open uploaded file")
	}
	defer multipartFile.Close()

	zipReader, err := zip.NewReader(multipartFile, req.File.Size)
	if err != nil {
		return nil, core.BadRequestError("Invalid zip file")
	}

	// Security Check: Validate size, count, and paths before processing
	if err := utils.ValidateZipArchive(zipReader.File); err != nil {
		return nil, core.BadRequestError(err.Error())
	}

	// Index html must at root
	if err := utils.ValidateEntrypoint(zipReader.File); err != nil {
		return nil, core.BadRequestError(err.Error())
	}

	// 2. Content Checks (Static files only)
	if err := utils.ValidateStaticFileTypes(zipReader.File); err != nil {
		return nil, core.BadRequestError(err.Error())
	}

	// 3. Upload to MinIO (With Rollback Tracking)
	// Generate a unique path upfront (e.g., deployments/1/20260112-1530-nano)
	timestamp := time.Now().Format("20060102-150405")
	uniqueSuffix := time.Now().UnixNano()
	outputPrefix := fmt.Sprintf("deployments/%d/%s-%d", projectID, timestamp, uniqueSuffix)

	var uploadedObjects []string // Track files for rollback

	log.Printf("Starting upload to %s", outputPrefix)

	for _, file := range zipReader.File {
		if file.FileInfo().IsDir() {
			continue
		}

		// Upload file
		err := s.uploadFileToMinio(ctx, outputPrefix, file)
		if err != nil {
			// MINIO FAILURE: Trigger Rollback
			log.Printf("Upload failed for %s: %v. Rolling back...", file.Name, err)
			s.rollbackMinioUploads(ctx, uploadedObjects)
			return nil, core.InternalError("Failed to upload files: " + err.Error())
		}

		// Track successful upload (construct the full key)
		fullKey := fmt.Sprintf("%s/%s", outputPrefix, file.Name)
		uploadedObjects = append(uploadedObjects, fullKey)
	}

	// 4. Create Database Record (The "Commit" Point)
	deployment := &models.Deployment{
		ProjectID:    projectID,
		Status:       enums.DeploymentStatusReady,
		OutputPrefix: outputPrefix,
	}

	// Attempt to save to DB
	if err := s.repo.Create(ctx, deployment); err != nil {
		// DB FAILURE: Trigger Rollback
		// The files are in MinIO, but the DB failed. We must delete the files to avoid "ghost" storage.
		log.Printf("Database creation failed. Rolling back MinIO storage...")
		s.rollbackMinioUploads(ctx, uploadedObjects)
		return nil, core.ParseDatabaseError(err)
	}

	log.Printf("Deployment %d created successfully", deployment.ID)
	return utils.EntityToDto[*response.DeploymentDto](deployment)
}

func (s *DeploymentService) uploadFileToMinio(ctx context.Context, prefix string, zf *zip.File) error {
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

	// Use LimitReader to enforce the uncompressed size validated earlier
	// This prevents the file from reading more bytes than the header claimed
	limitReader := io.LimitReader(rc, int64(zf.UncompressedSize64))

	objectName := fmt.Sprintf("%s/%s", prefix, zf.Name)

	_, err = s.minioClient.PutObject(ctx, "static-sites", objectName, limitReader, int64(zf.UncompressedSize64), minioGo.PutObjectOptions{
		ContentType: contentType,
	})

	return err
}

// ---------------------------------------------------------
// Helper: Rollback (Compensating Transaction)
// ---------------------------------------------------------
func (s *DeploymentService) rollbackMinioUploads(ctx context.Context, objects []string) {
	// Execute deletion in a background goroutine or immediately.
	// For critical consistency, doing it immediately is safer, though it adds latency to the error response.

	// Create a channel for objects to delete (MinIO optimized batch deletion)
	objectsCh := make(chan minioGo.ObjectInfo)

	go func() {
		defer close(objectsCh)
		for _, objName := range objects {
			objectsCh <- minioGo.ObjectInfo{Key: objName}
		}
	}()

	// Call RemoveObjects API
	errorCh := s.minioClient.RemoveObjects(ctx, "static-sites", objectsCh, minioGo.RemoveObjectsOptions{})

	// Log any errors during rollback
	for err := range errorCh {
		log.Printf("CRITICAL: Failed to delete object during rollback: %s, Error: %v", err.ObjectName, err.Err)
	}
}

func (s *DeploymentService) GetCurrentDeploymentFilesByProjectSubdomain(ctx context.Context, subdomain string, fileName string, clientETag string) (*response.FileDownloadDto, error) {
	project, err := s.projectRepo.FindBySubdomain(ctx, subdomain)
	if err != nil || project == nil || project.CurrentDeploymentID == 0 {
		return nil, core.NotFoundError()
	}

	deployment, err := s.repo.FindByID(ctx, project.CurrentDeploymentID)
	if err != nil || deployment.Status != enums.DeploymentStatusLive {
		return nil, core.NotFoundError()
	}

	if deployment.OutputPrefix == "" {
		deployment.OutputPrefix = fmt.Sprintf("deployments/%d/%d", deployment.ProjectID, deployment.ID)
	}

	cleanFileName := strings.TrimPrefix(fileName, "/")
	if cleanFileName == "" {
		cleanFileName = "index.html"
	}

	bucket := utils.GetEnv("MINIO_BUCKET", "static-sites")
	objectPath := fmt.Sprintf("%s/%s", deployment.OutputPrefix, cleanFileName)

	stat, err := s.minioClient.StatObject(ctx, bucket, objectPath, minioGo.StatObjectOptions{})

	if err != nil {
		// The request is a "Route" (No extension like .js, png) ---
		if !strings.Contains(cleanFileName, ".") {

			// Clean URL (e.g. /about -> about.html)
			tryPath := objectPath + ".html"
			if statH, errH := s.minioClient.StatObject(ctx, bucket, tryPath, minioGo.StatObjectOptions{}); errH == nil {
				return s.serveS3Object(ctx, bucket, tryPath, statH, http.StatusOK, clientETag)
			}

			// Try Folder Index (e.g. contact/ -> contact/index.html)
			if strings.HasSuffix(cleanFileName, "/") {
				tryPath = strings.TrimSuffix(objectPath, "/") + "/index.html"
				if statF, errF := s.minioClient.StatObject(ctx, bucket, tryPath, minioGo.StatObjectOptions{}); errF == nil {
					return s.serveS3Object(ctx, bucket, tryPath, statF, http.StatusOK, clientETag)
				}
			}

			// SPA FALLBACK
			if deployment.IsSPA {
				spaPath := fmt.Sprintf("%s/index.html", deployment.OutputPrefix)
				if statSPA, errSPA := s.minioClient.StatObject(ctx, bucket, spaPath, minioGo.StatObjectOptions{}); errSPA == nil {
					return s.serveS3Object(ctx, bucket, spaPath, statSPA, http.StatusOK, clientETag)
				}
			}
		}

		// 404.html
		return s.handleErrorFallback(ctx, deployment)
	}
	return s.serveS3Object(ctx, bucket, objectPath, stat, http.StatusOK, clientETag)
}

func (s *DeploymentService) serveS3Object(ctx context.Context, bucket, path string, stat minioGo.ObjectInfo, status int, clientETag string) (*response.FileDownloadDto, error) {

	headers := map[string]string{
		"Cache-Control":          "public, max-age=3600",
		"ETag":                   stat.ETag,
		"X-Content-Type-Options": "nosniff",
	}

	if strings.Contains(path, "/assets/") {
		headers["Cache-Control"] = "public, max-age=31536000, immutable"
	}

	if clientETag != "" && clientETag == stat.ETag {
		return &response.FileDownloadDto{
			NotModified: true,
			Headers:     headers,
			StatusCode:  http.StatusNotModified,
		}, nil
	}

	obj, err := s.minioClient.GetObject(ctx, bucket, path, minioGo.GetObjectOptions{})
	if err != nil {
		return nil, core.ParseMinioError(err)
	}

	return &response.FileDownloadDto{
		Stream:      obj,
		Size:        stat.Size,
		ContentType: stat.ContentType,
		NotModified: false,
		StatusCode:  status,
	}, nil
}

func (s *DeploymentService) handleErrorFallback(ctx context.Context, deployment *models.Deployment) (*response.FileDownloadDto, error) {
	bucket := utils.GetEnv("MINIO_BUCKET", "static-sites")
	notFoundPath := fmt.Sprintf("%s/404.html", deployment.OutputPrefix)

	if stat, err := s.minioClient.StatObject(ctx, bucket, notFoundPath, minioGo.StatObjectOptions{}); err == nil {
		return s.serveS3Object(ctx, bucket, notFoundPath, stat, http.StatusNotFound, "")
	}

	return nil, core.NotFoundError()
}

func (s *DeploymentService) TurnDeploymentLive(ctx context.Context, deploymentID uint, userID uint) error {
	deployment, err := s.repo.FindByID(ctx, deploymentID)

	if err != nil {
		return core.ParseDatabaseError(err)
	}

	project, err := s.projectRepo.FindByID(ctx, deployment.ProjectID)
	if err != nil {
		return core.ParseDatabaseError(err)
	}

	if project == nil {
		return core.NotFoundError()
	}

	if project.UserID != userID {
		return core.ForbiddenError()
	}

	if project.CurrentDeploymentID != 0 {
		return core.BadRequestError("Project already has a live deployment")
	}

	if deployment == nil {
		return core.NotFoundError()
	}

	if deployment.Status != enums.DeploymentStatusReady {
		return core.BadRequestError("Deployment is not ready or is living")
	}

	deployment.Status = enums.DeploymentStatusLive
	if err := s.repo.Update(ctx, deployment); err != nil {
		return core.ParseDatabaseError(err)
	}

	project.CurrentDeploymentID = deploymentID
	if err := s.projectRepo.Update(ctx, project); err != nil {
		return core.ParseDatabaseError(err)
	}

	return nil
}

func (s *DeploymentService) TurnDeploymentOffline(ctx context.Context, deploymentID uint, userID uint) error {
	deployment, err := s.repo.FindByID(ctx, deploymentID)

	if err != nil {
		return core.ParseDatabaseError(err)
	}

	project, err := s.projectRepo.FindByID(ctx, deployment.ProjectID)
	if err != nil {
		return core.ParseDatabaseError(err)
	}

	if project == nil {
		return core.NotFoundError()
	}

	if project.UserID != userID {
		return core.ForbiddenError()
	}

	if project.CurrentDeploymentID != deploymentID {
		return core.BadRequestError("Deployment is not the current deployment")
	}

	deployment.Status = enums.DeploymentStatusReady
	if err := s.repo.Update(ctx, deployment); err != nil {
		return core.ParseDatabaseError(err)
	}

	project.CurrentDeploymentID = 0
	if err := s.projectRepo.Update(ctx, project); err != nil {
		return core.ParseDatabaseError(err)
	}

	return nil
}

func (s *DeploymentService) ToggleIsSPAMode(ctx context.Context, deploymentID uint, userID uint) error {
	deployment, err := s.repo.FindByID(ctx, deploymentID)

	if err != nil {
		return core.ParseDatabaseError(err)
	}

	project, err := s.projectRepo.FindByID(ctx, deployment.ProjectID)
	if err != nil {
		return core.ParseDatabaseError(err)
	}

	if project == nil {
		return core.NotFoundError()
	}

	if project.UserID != userID {
		return core.ForbiddenError()
	}

	deployment.IsSPA = !deployment.IsSPA
	if err := s.repo.Update(ctx, deployment); err != nil {
		return core.ParseDatabaseError(err)
	}

	return nil
}

func (s *DeploymentService) DeleteDeployment(ctx context.Context, deploymentID uint, userID uint) error {
	// 1. Fetch deployment and verify ownership
	// We need the OutputPrefix from the DB before we delete the record
	deployment, err := s.repo.FindByID(ctx, deploymentID)
	if err != nil {
		return core.ParseDatabaseError(err)
	}
	if deployment == nil {
		return core.NotFoundError()
	}

	// Security: Ensure the user owns the project this deployment belongs to
	project, err := s.projectRepo.FindByID(ctx, deployment.ProjectID)
	if err != nil || project.UserID != userID {
		return core.ForbiddenError()
	}

	// 2. Perform Physical Deletion (MinIO)
	// We do this first or in a way that ensures we don't lose the Prefix string
	log.Printf("Hard deleting deployment %d from storage at prefix: %s", deploymentID, deployment.OutputPrefix)

	err = s.deleteMinioFolder(ctx, deployment.OutputPrefix)
	if err != nil {
		// We log the error but might choose to continue to DB deletion
		// depending on how "stuck" you want the UI to get.
		log.Printf("Warning: MinIO cleanup failed for %s: %v", deployment.OutputPrefix, err)
	}

	// 3. Perform Logical Deletion (Database)
	// Hard delete the record from the DB
	if err := s.repo.Delete(ctx, deployment); err != nil {
		return core.ParseDatabaseError(err)
	}

	return nil
}

// Internal helper to handle the MinIO prefix deletion
func (s *DeploymentService) deleteMinioFolder(ctx context.Context, prefix string) error {
	bucketName := "static-sites"

	// Guard: Never allow empty prefix deletion
	if prefix == "" || prefix == "/" {
		return fmt.Errorf("refusing to delete unsafe prefix: %s", prefix)
	}

	// 1. List all objects with the prefix
	objectsCh := make(chan minioGo.ObjectInfo)

	go func() {
		defer close(objectsCh)
		// Recursive: true ensures we get files in subfolders like /css/*
		for object := range s.minioClient.ListObjects(ctx, bucketName, minioGo.ListObjectsOptions{
			Prefix:    prefix,
			Recursive: true,
		}) {
			if object.Err != nil {
				log.Printf("Error listing object for deletion: %v", object.Err)
				continue
			}
			objectsCh <- object
		}
	}()

	// 2. Batch Delete
	errorCh := s.minioClient.RemoveObjects(ctx, bucketName, objectsCh, minioGo.RemoveObjectsOptions{})

	// 3. Collect errors
	var deleteErrors []string
	for err := range errorCh {
		deleteErrors = append(deleteErrors, fmt.Sprintf("%s: %v", err.ObjectName, err.Err))
	}

	if len(deleteErrors) > 0 {
		return fmt.Errorf("failed to delete some objects: %s", strings.Join(deleteErrors, ", "))
	}

	return nil
}
