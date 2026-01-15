package service

import (
	"archive/zip"
	"context"
	"fmt"
	"log"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

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

	return utils.EntityToDto[*response.DeploymentDto](deployment)
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
		// The request is a "Route" (No extension like .js, .png) ---
		if !strings.Contains(cleanFileName, ".") {

			// Clean URL (e.g. /about -> about.html)
			tryPath := objectPath + ".html"
			if statH, errH := s.minioClient.StatObject(ctx, bucket, tryPath, minioGo.StatObjectOptions{}); errH == nil {
				return s.serveS3Object(ctx, bucket, tryPath, statH, http.StatusOK, clientETag)
			}

			// Try Folder Index (e.g. /contact -> contact/index.html)
			tryPath = strings.TrimSuffix(objectPath, "/") + "/index.html"
			if statF, errF := s.minioClient.StatObject(ctx, bucket, tryPath, minioGo.StatObjectOptions{}); errF == nil {
				return s.serveS3Object(ctx, bucket, tryPath, statF, http.StatusOK, clientETag)
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
