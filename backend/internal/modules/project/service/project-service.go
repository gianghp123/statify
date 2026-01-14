package service

import (
	"context"
	"fmt"
	"log"

	"github.com/gianghp/statify/internal/core"
	coreRepo "github.com/gianghp/statify/internal/core/repository"
	"github.com/gianghp/statify/internal/database/models"
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
}

func NewProjectService(repo repository.IProjectRepository, deploymentRepo deploymentRepo.IDeploymentRepository, minioClient minio.Interface) *ProjectService {
	return &ProjectService{repo: repo, deploymentRepo: deploymentRepo, minioClient: minioClient}
}

func (s *ProjectService) CreateProject(ctx context.Context, userID uint, req *request.CreateProjectRequest) (*response.ProjectDto, error) {
	project := &models.Project{
		Name:      req.Name,
		Subdomain: req.Subdomain,
		UserID:    userID,
	}

	if err := s.repo.Create(ctx, project); err != nil {
		return nil, core.ParseDatabaseError(err)
	}

	projectDto, err := utils.EntityToDto[response.ProjectDto](project)
	if err != nil {
		return nil, core.InternalError()
	}

	return projectDto, nil
}

func (s *ProjectService) ListProjects(ctx context.Context, userID uint, page int, limit int) (*coreRepo.PaginatedEntities[*response.ProjectDto], error) {
	projects, err := s.repo.FindAllByUserID(ctx, userID, page, limit)
	if err != nil {
		return nil, core.ParseDatabaseError(err)
	}

	projectDtos, err := utils.EntitiesToDto[response.ProjectDto](projects.Entities)
	if err != nil {
		return nil, core.InternalError()
	}

	for _, dto := range projectDtos {
		s.enrichProjectDto(ctx, dto)
	}

	return &coreRepo.PaginatedEntities[*response.ProjectDto]{
		Entities:   projectDtos,
		Pagination: projects.Pagination,
	}, nil
}

func (s *ProjectService) GetProjectByID(ctx context.Context, id uint, userID uint) (*response.ProjectDto, error) {
	project, err := s.repo.FindByID(ctx, id)

	if err != nil {
		return nil, core.ParseDatabaseError(err)
	}

	if project.UserID != userID {
		return nil, core.ForbiddenError()
	}

	if project == nil {
		return nil, core.NotFoundError()
	}

	projectDto, err := utils.EntityToDto[response.ProjectDto](project)
	if err != nil {
		return nil, core.InternalError()
	}

	s.enrichProjectDto(ctx, projectDto)

	return projectDto, nil
}

func (s *ProjectService) enrichProjectDto(ctx context.Context, dto *response.ProjectDto) {
	dto.URL = fmt.Sprintf("https://%s.%s", dto.Subdomain, utils.GetEnv("DOMAIN", "statify.online"))
	latest, err := s.deploymentRepo.FindLatestByProjectID(ctx, dto.ID)
	if err == nil && latest != nil {
		dto.Status = latest.Status
	}
}

func (s *ProjectService) UpdateProject(ctx context.Context, id uint, req *request.UpdateProjectRequest) error {
	project, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return core.ParseDatabaseError(err)
	}

	if project == nil {
		return core.NotFoundError()
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
		project, err := txRepo.FindByID(ctx, projectID)
		if err != nil {
			return core.ParseDatabaseError(err)
		}

		if project.UserID != userID {
			return core.ForbiddenError()
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
