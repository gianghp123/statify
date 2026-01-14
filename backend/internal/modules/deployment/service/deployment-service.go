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
	"github.com/gianghp/statify/internal/storage/minio"
	"github.com/gianghp/statify/internal/utils"
	minioGo "github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

type DeploymentService struct {
	repo        repository.IDeploymentRepository
	projectRepo projectRepo.IProjectRepository
	minioClient minio.Interface
}

func NewDeploymentService(repo repository.IDeploymentRepository, projectRepo projectRepo.IProjectRepository, minioClient minio.Interface) *DeploymentService {
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

	err = s.repo.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)

		if err := txRepo.Create(ctx, deployment); err != nil {
			return err
		}

		// Update output prefix with ID and Timestamp
		deployment.OutputPrefix = fmt.Sprintf("deployments/%d/%d-%s", projectID, deployment.ID, deployment.CreatedAt.Format("20060102150405"))
		if err := txRepo.Update(ctx, deployment); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
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

		if err := s.uploadFileToMinio(ctx, deployment.OutputPrefix, file); err != nil {

			// Mark as failed in DB
			deployment.Status = enums.DeploymentStatusFailed
			_ = s.repo.Update(ctx, deployment)
			return nil, core.InternalError(err.Error())
		}
	}

	// 5. Success! Update status to Ready
	deployment.Status = enums.DeploymentStatusReady
	if err := s.repo.Update(ctx, deployment); err != nil {
		return nil, core.ParseDatabaseError(err)
	}

	return utils.EntityToDto[response.DeploymentDto](deployment)
}

// Helper function to handle individual file uploads
func (s *DeploymentService) uploadFileToMinio(ctx context.Context, outputPrefix string, zf *zip.File) error {
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

	// Construct path: deployments/1/10-20260112151513/css/style.css
	objectName := fmt.Sprintf("%s/%s", outputPrefix, zf.Name)

	_, err = s.minioClient.PutObject(ctx, "static-sites", objectName, rc, zf.FileInfo().Size(), minioGo.PutObjectOptions{
		ContentType: contentType,
	})

	return err
}

func (s *DeploymentService) GetHistory(ctx context.Context, userID uint, projectID uint, page int, limit int) (*coreRepo.PaginatedEntities[*response.DeploymentDto], error) {
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

	deployments, err := s.repo.FindAllByProjectID(ctx, projectID, page, limit)
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

	if deployment.Status != enums.DeploymentStatusLive {
		return nil, core.NotFoundError()
	}

	if deployment.OutputPrefix == "" {
		deployment.OutputPrefix = fmt.Sprintf("deployments/%d/%d", deployment.ProjectID, deployment.ID)
	}

	objectPath := fmt.Sprintf("%s/%s", deployment.OutputPrefix, fileName)

	stat, err := s.minioClient.StatObject(context.Background(), utils.GetEnv("MINIO_BUCKET", "static-sites"), objectPath, minioGo.StatObjectOptions{})

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

	// 3. Return DTO with stream
	obj, err := s.minioClient.GetObject(context.Background(), utils.GetEnv("MINIO_BUCKET", "static-sites"), objectPath, minioGo.GetObjectOptions{})
	if err != nil {
		return nil, core.ParseMinioError(err)
	}

	return &response.FileDownloadDto{
		Stream:      obj,
		Size:        stat.Size,
		ContentType: stat.ContentType,
		Headers:     headers,
		NotModified: false,
	}, nil
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
