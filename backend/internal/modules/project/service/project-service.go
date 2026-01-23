package service

import (
	"context"
	"fmt"
	"log"

	"github.com/gianghp/statify/internal/core"
	"github.com/gianghp/statify/internal/core/enums"
	coreRepo "github.com/gianghp/statify/internal/core/repository"
	"github.com/gianghp/statify/internal/database/models"
	"github.com/gianghp/statify/internal/modules/auth/policy"
	deploymentRepo "github.com/gianghp/statify/internal/modules/deployment/repository"
	"github.com/gianghp/statify/internal/modules/project/dtos/request"
	"github.com/gianghp/statify/internal/modules/project/dtos/response"
	"github.com/gianghp/statify/internal/modules/project/repository"
	"github.com/gianghp/statify/internal/storage/minio"
	"github.com/gianghp/statify/internal/utils"
	"gorm.io/gorm"
)

type ProjectService struct {
	repo           repository.IProjectRepository
	deploymentRepo deploymentRepo.IDeploymentRepository
	minioClient    minio.Interface
	policy         policy.IAccessPolicy
}

func NewProjectService(repo repository.IProjectRepository, deploymentRepo deploymentRepo.IDeploymentRepository, minioClient minio.Interface, policy policy.IAccessPolicy) *ProjectService {
	return &ProjectService{repo: repo, deploymentRepo: deploymentRepo, minioClient: minioClient, policy: policy}
}

func (s *ProjectService) CreateProject(ctx context.Context, userID uint, req *request.CreateProjectRequest) (*response.ProjectDto, error) {
	if _, err := s.repo.FindBySubdomain(ctx, req.Subdomain); err == nil {
		return nil, core.BadRequestError("Subdomain already exists, please choose another one.")
	}

	project := &models.Project{
		Name:      req.Name,
		Subdomain: req.Subdomain,
		UserID:    userID,
	}

	if err := s.repo.Create(ctx, project); err != nil {
		return nil, core.ParseDatabaseError(err)
	}

	projectDto, err := utils.EntityToDto[*response.ProjectDto](project)
	if err != nil {
		return nil, core.InternalError()
	}

	return projectDto, nil
}

func (s *ProjectService) ListProjects(ctx context.Context, userID uint, page int, limit int, status enums.DeploymentStatus) (coreRepo.PaginatedEntities[*response.ProjectDto], error) {
	result := coreRepo.PaginatedEntities[*response.ProjectDto]{
		Entities:   []*response.ProjectDto{},
		Pagination: coreRepo.Pagination{Page: page, Limit: limit},
	}

	projects, err := s.repo.FindAllByUserID(ctx, userID, page, limit, status)
	if err != nil {
		return result, core.ParseDatabaseError(err)
	}

	for _, project := range projects.Entities {
		projectDto, err := s.transformProjectToDto(project)
		if err != nil {
			return result, core.InternalError()
		}
		result.Entities = append(result.Entities, projectDto)
	}

	result.Pagination = projects.Pagination

	return result, nil
}

func (s *ProjectService) GetProjectByID(ctx context.Context, id uint, userID uint) (*response.ProjectDto, error) {
	project, err := s.policy.CheckProjectAccess(ctx, userID, id)

	if err != nil {
		return nil, err
	}

	projectDto, err := s.transformProjectToDto(project)
	if err != nil {
		return nil, core.InternalError()
	}

	return projectDto, nil
}

func (s *ProjectService) transformProjectToDto(project *models.Project) (*response.ProjectDto, error) {
	projectDto, err := utils.EntityToDto[*response.ProjectDto](project)
	if err != nil {
		return nil, core.InternalError()
	}

	projectDto.URL = fmt.Sprintf("https://%s.%s", project.Subdomain, utils.GetEnv("DOMAIN", "statify.app"))

	if len(project.Deployments) == 0 {
		return projectDto, nil
	}

	latestDeployments := project.Deployments[0]

	projectDto.Status = latestDeployments.Status

	return projectDto, nil
}

func (s *ProjectService) UpdateProject(ctx context.Context, userID uint, id uint, req *request.UpdateProjectRequest) error {
	project, err := s.policy.CheckProjectAccess(ctx, userID, id)
	if err != nil {
		return err
	}

	project.Name = req.Name
	project.Subdomain = req.Subdomain

	if err := s.repo.Update(ctx, project); err != nil {
		return core.ParseDatabaseError(err)
	}

	return nil
}

func (s *ProjectService) DeleteProject(ctx context.Context, projectID uint, userID uint) error {
	return s.repo.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)
		txDeploymentRepo := s.deploymentRepo.WithTx(tx)

		// 1. Fetch the project and its deployments first to get the MinIO keys
		project, err := s.policy.CheckProjectAccess(ctx, userID, projectID)
		if err != nil {
			return err
		}

		deployments, err := txDeploymentRepo.FindAllByProjectID(ctx, projectID, 1, 1000)
		if err != nil {
			return core.ParseDatabaseError(err)
		}

		// 2. Delete from DB (The migration handles cascading to deployments)
		if err := txRepo.Delete(ctx, project); err != nil {
			return core.ParseDatabaseError(err)
		}

		// 3. Database transaction is ready to commit.
		// Now attempt MinIO cleanup.
		bucket := utils.GetEnv("MINIO_BUCKET", "static-sites")
		for _, d := range deployments.Entities {
			// Delete the versioned folder
			if d.OutputPrefix != "" {
				err := s.minioClient.RemoveObjectsByPrefix(ctx, bucket, d.OutputPrefix)
				if err != nil {
					log.Printf("Failed to delete MinIO prefix %s: %v", d.OutputPrefix, err)
				}
			} else {
				// Fallback for old style paths if any
				prefix := fmt.Sprintf("deployments/%d/%d", d.ProjectID, d.ID)
				err := s.minioClient.RemoveObjectsByPrefix(ctx, bucket, prefix)
				if err != nil {
					log.Printf("Failed to delete MinIO fallback prefix %s: %v", prefix, err)
				}
			}
		}

		return nil // If nil is returned, GORM commits the DB transaction
	})
}
