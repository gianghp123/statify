package service

import (
	"archive/zip"
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gianghp/statify/internal/core"
	"github.com/gianghp/statify/internal/core/enums"
	coreRepo "github.com/gianghp/statify/internal/core/repository"
	"github.com/gianghp/statify/internal/database/models"
	"github.com/gianghp/statify/internal/modules/auth/policy"
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
	policy      policy.IAccessPolicy
}

func NewDeploymentService(repo repository.IDeploymentRepository, projectRepo projectRepo.IProjectRepository, minioClient minio.Interface, policy policy.IAccessPolicy) *DeploymentService {
	return &DeploymentService{repo: repo, projectRepo: projectRepo, minioClient: minioClient, policy: policy}
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

	_, err := s.policy.CheckProjectAccess(ctx, userID, projectID)
	if err != nil {
		return result, err
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
	_, deployment, err := s.policy.CheckDeploymentAccess(ctx, userID, id)
	if err != nil {
		return nil, err
	}

	deploymentDto, err := utils.EntityToDto[*response.DeploymentDto](deployment)
	if err != nil {
		return nil, core.InternalError()
	}

	return deploymentDto, nil
}

func (s *DeploymentService) CreateDeployment(ctx context.Context, userID uint, projectID uint, req *request.CreateDeploymentRequest) (*response.DeploymentDto, error) {
	// 1. Authorization & Validation
	_, err := s.policy.CheckProjectAccess(ctx, userID, projectID)
	if err != nil {
		return nil, err
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

	objectName := fmt.Sprintf("%s/temp", outputPrefix)

	_, err = s.minioClient.StatObject(ctx, "static-sites", objectName, minioGo.StatObjectOptions{})
	if err == nil {
		return nil, core.BadRequestError("Deployment already exists")
	}

	_, err = s.minioClient.PutObject(ctx, "static-sites", objectName, multipartFile, req.File.Size, minioGo.PutObjectOptions{
		ContentType: "application/zip",
	})

	if err != nil {
		return nil, core.InternalError()
	}

	deployment := &models.Deployment{
		ProjectID:    projectID,
		Status:       enums.DeploymentStatusQueued,
		OutputPrefix: outputPrefix,
	}

	err = s.repo.Create(ctx, deployment)
	if err != nil {
		return nil, core.InternalError()
	}

	log.Printf("Deployment %d created successfully", deployment.ID)
	return utils.EntityToDto[*response.DeploymentDto](deployment)
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
				return s.serveS3Object(ctx, deployment.ProjectID, deployment.ID, bucket, tryPath, statH, http.StatusOK, clientETag)
			}

			// Try Folder Index (e.g. contact/ -> contact/index.html)
			if strings.HasSuffix(cleanFileName, "/") {
				tryPath = strings.TrimSuffix(objectPath, "/") + "/index.html"
				if statF, errF := s.minioClient.StatObject(ctx, bucket, tryPath, minioGo.StatObjectOptions{}); errF == nil {
					return s.serveS3Object(ctx, deployment.ProjectID, deployment.ID, bucket, tryPath, statF, http.StatusOK, clientETag)
				}
			}

			// SPA FALLBACK
			if deployment.IsSPA {
				spaPath := fmt.Sprintf("%s/index.html", deployment.OutputPrefix)
				if statSPA, errSPA := s.minioClient.StatObject(ctx, bucket, spaPath, minioGo.StatObjectOptions{}); errSPA == nil {
					return s.serveS3Object(ctx, deployment.ProjectID, deployment.ID, bucket, spaPath, statSPA, http.StatusOK, clientETag)
				}
			}
		}

		// 404.html
		return s.handleErrorFallback(ctx, deployment)
	}
	return s.serveS3Object(ctx, deployment.ProjectID, deployment.ID, bucket, objectPath, stat, http.StatusOK, clientETag)
}

func (s *DeploymentService) serveS3Object(ctx context.Context, projectID, deploymentID uint, bucket, path string, stat minioGo.ObjectInfo, status int, clientETag string) (*response.FileDownloadDto, error) {

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
		ProjectID:    projectID,
		DeploymentID: deploymentID,
		Stream:       obj,
		Size:         stat.Size,
		ContentType:  stat.ContentType,
		NotModified:  false,
		StatusCode:   status,
	}, nil
}

func (s *DeploymentService) handleErrorFallback(ctx context.Context, deployment *models.Deployment) (*response.FileDownloadDto, error) {
	bucket := utils.GetEnv("MINIO_BUCKET", "static-sites")
	notFoundPath := fmt.Sprintf("%s/404.html", deployment.OutputPrefix)

	if stat, err := s.minioClient.StatObject(ctx, bucket, notFoundPath, minioGo.StatObjectOptions{}); err == nil {
		return s.serveS3Object(ctx, deployment.ProjectID, deployment.ID, bucket, notFoundPath, stat, http.StatusNotFound, "")
	}

	return nil, core.NotFoundError()
}

func (s *DeploymentService) TurnDeploymentLive(ctx context.Context, deploymentID uint, userID uint) error {
	project, deployment, err := s.policy.CheckDeploymentAccess(ctx, userID, deploymentID)
	if err != nil {
		return err
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
	project, _, err := s.policy.CheckDeploymentAccess(ctx, userID, deploymentID)
	if err != nil {
		return err
	}

	if project.CurrentDeploymentID != deploymentID {
		return core.BadRequestError("Deployment is not the current deployment")
	}

	deployment := &models.Deployment{Model: gorm.Model{ID: deploymentID}}
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
	_, deployment, err := s.policy.CheckDeploymentAccess(ctx, userID, deploymentID)
	if err != nil {
		return err
	}

	deployment.IsSPA = !deployment.IsSPA
	if err := s.repo.Update(ctx, deployment); err != nil {
		return core.ParseDatabaseError(err)
	}

	return nil
}

func (s *DeploymentService) DeleteDeployment(ctx context.Context, deploymentID uint, userID uint) error {
	_, deployment, err := s.policy.CheckDeploymentAccess(ctx, userID, deploymentID)
	if err != nil {
		return err
	}

	if deployment.Status == enums.DeploymentStatusLive {
		return core.BadRequestError("Deployment is currently live, please turn it offline first")
	}

	deployment.Status = enums.DeploymentStatusPendingDelete
	if err := s.repo.Update(ctx, deployment); err != nil {
		return core.ParseDatabaseError(err)
	}

	return nil
}
