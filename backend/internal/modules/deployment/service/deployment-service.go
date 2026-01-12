package service

import (
	"archive/zip"
	"context"
	"fmt"
	"log"
	"mime"
	"path/filepath"

	"github.com/gianghp/statify/internal/core"
	"github.com/gianghp/statify/internal/core/enums"
	coreRepo "github.com/gianghp/statify/internal/core/repository"
	"github.com/gianghp/statify/internal/database/models"
	"github.com/gianghp/statify/internal/modules/deployment/dtos/request"
	"github.com/gianghp/statify/internal/modules/deployment/dtos/response"
	"github.com/gianghp/statify/internal/modules/deployment/repository"
	projectRepo "github.com/gianghp/statify/internal/modules/project/repository"
	"github.com/gianghp/statify/internal/utils"
	"github.com/minio/minio-go/v7"
)

type DeploymentService struct {
	repo        repository.IDeploymentRepository
	projectRepo projectRepo.IProjectRepository
	minioClient MinioInterface
}

func NewDeploymentService(repo repository.IDeploymentRepository, projectRepo projectRepo.IProjectRepository, minioClient MinioInterface) *DeploymentService {
	return &DeploymentService{repo: repo, projectRepo: projectRepo, minioClient: minioClient}
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

	// 2. Create the Database Record (Initially "Processing")
	deployment := &models.Deployment{
		ProjectID: projectID,
		Status:    enums.DeploymentStatusProcessing,
	}
	if err := s.repo.Create(ctx, deployment); err != nil {
		return nil, core.ParseDatabaseError(err)
	}

	// 3. Process the Zip File
	// Open the uploaded multipart file
	multipartFile, err := req.File.Open()
	if err != nil {
		return nil, core.BadRequestError("Could not open uploaded file")
	}
	defer multipartFile.Close()

	// zip.NewReader requires an io.ReaderAt.
	// Since multipart.File is an interface, we read it into a SectionReader or a buffer.
	// For MVP local dev, reading into memory is fine, or use req.File directly if it supports ReaderAt.
	zipReader, err := zip.NewReader(multipartFile, req.File.Size)
	if err != nil {
		return nil, core.BadRequestError("Invalid zip file")
	}

	// Should be go routine
	log.Println("Uploading files to MinIO")
	for _, file := range zipReader.File {
		if file.FileInfo().IsDir() {
			continue
		}

		if err := s.uploadFileToMinio(ctx, projectID, deployment.ID, file); err != nil {

			// Mark as failed in DB
			deployment.Status = enums.DeploymentStatusFailed
			_ = s.repo.Update(ctx, deployment)
			return nil, core.InternalError(err.Error())
		}
	}

	// 5. Success! Update status to Ready and set this as the project's active deployment
	deployment.Status = enums.DeploymentStatusReady
	if err := s.repo.Update(ctx, deployment); err != nil {
		return nil, core.ParseDatabaseError(err)
	}

	// Set current deployment to offline
	if project.CurrentDeploymentID != 0 {
		currentDeployment, err := s.repo.FindByID(ctx, project.CurrentDeploymentID)

		if err != nil {
			return nil, core.ParseDatabaseError(err)
		}

		if currentDeployment == nil {
			return nil, core.NotFoundError()
		}

		if currentDeployment.Status == enums.DeploymentStatusReady {
			currentDeployment.Status = enums.DeploymentStatusUploaded
			if err := s.repo.Update(ctx, currentDeployment); err != nil {
				return nil, core.ParseDatabaseError(err)
			}
		}
	}

	// Update the project to point to this new active deployment
	project.CurrentDeploymentID = deployment.ID
	if err := s.projectRepo.Update(ctx, project); err != nil {
		return nil, core.ParseDatabaseError(err)
	}

	return utils.EntityToDto[response.DeploymentDto](deployment)
}

// Helper function to handle individual file uploads
func (s *DeploymentService) uploadFileToMinio(ctx context.Context, projID, depID uint, zf *zip.File) error {
	rc, err := zf.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	// Detect MIME type
	contentType := mime.TypeByExtension(filepath.Ext(zf.Name))
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Construct path: deployments/1/10/css/style.css
	objectName := fmt.Sprintf("deployments/%d/%d/%s", projID, depID, zf.Name)

	_, err = s.minioClient.PutObject(ctx, "static-sites", objectName, rc, zf.FileInfo().Size(), minio.PutObjectOptions{
		ContentType: contentType,
	})

	return err
}

func (s *DeploymentService) GetHistory(ctx context.Context, userID uint, projectID uint) (*coreRepo.PaginatedEntities[*response.DeploymentDto], error) {
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

	deployments, err := s.repo.FindAllByProjectID(ctx, projectID)
	if err != nil {
		return nil, core.ParseDatabaseError(err)
	}

	deploymentDtos, err := utils.EntitiesToDto[response.DeploymentDto](deployments.Entities)
	if err != nil {
		return nil, core.InternalError()
	}

	return &coreRepo.PaginatedEntities[*response.DeploymentDto]{
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

	deploymentDto, err := utils.EntityToDto[response.DeploymentDto](deployment)
	if err != nil {
		return nil, core.InternalError()
	}

	return deploymentDto, nil
}

func (s *DeploymentService) GetCurrentDeploymentFilesByProjectSubdomain(ctx context.Context, subdomain string, fileName string, clientETag string) (*response.FileDownloadDto, error) {
	project, err := s.projectRepo.FindBySubdomain(ctx, subdomain)
	if err != nil {
		return nil, core.ParseDatabaseError(err)
	}

	if project == nil {
		return nil, core.NotFoundError()
	}

	if project.CurrentDeploymentID == 0 {
		return nil, core.NotFoundError()
	}

	deployment, err := s.repo.FindByID(ctx, project.CurrentDeploymentID)
	if err != nil {
		return nil, core.ParseDatabaseError(err)
	}

	opjectPath := fmt.Sprintf("deployments/%d/%d/%s", project.ID, deployment.ID, fileName)

	stat, err := s.minioClient.StatObject(context.Background(), utils.GetEnv("MINIO_BUCKET", "static-sites"), opjectPath, minio.StatObjectOptions{})

	if err != nil {
		// Return specific error so controller knows if it's 404 or 500
		return nil, core.ParseMinioError(err)
	}

	// 2. Prepare Headers
	headers := map[string]string{
		"Cache-Control":          "public, max-age=3600",
		"ETag":                   stat.ETag,
		"X-Content-Type-Options": "nosniff",
	}

	if clientETag == stat.ETag {
		return &response.FileDownloadDto{
			NotModified: true,
			Headers:     headers,
		}, nil
	}

	object, err := s.minioClient.GetObject(context.Background(), utils.GetEnv("MINIO_BUCKET", "static-sites"), opjectPath, minio.GetObjectOptions{})
	if err != nil {
		return nil, core.ParseMinioError(err)
	}

	return &response.FileDownloadDto{
		Stream:      object,
		Size:        stat.Size,
		ContentType: stat.ContentType,
		Headers:     headers,
		NotModified: false,
	}, nil
}
