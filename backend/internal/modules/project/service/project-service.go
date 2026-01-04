package service

import (
	"context"
	"fmt"

	// ... (omitted similar imports if not needed, but I should probably keep them)
	"github.com/gianghp/statify/internal/core"
	coreRepo "github.com/gianghp/statify/internal/core/repository"
	"github.com/gianghp/statify/internal/database/models"
	deploymentRepo "github.com/gianghp/statify/internal/modules/deployment/repository"
	"github.com/gianghp/statify/internal/modules/project/dtos/request"
	"github.com/gianghp/statify/internal/modules/project/dtos/response"
	"github.com/gianghp/statify/internal/modules/project/repository"
	"github.com/gianghp/statify/internal/utils"
)

type ProjectService struct {
	repo           repository.IProjectRepository
	deploymentRepo deploymentRepo.IDeploymentRepository
}

func NewProjectService(repo repository.IProjectRepository, deploymentRepo deploymentRepo.IDeploymentRepository) *ProjectService {
	return &ProjectService{repo: repo, deploymentRepo: deploymentRepo}
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

func (s *ProjectService) ListProjects(ctx context.Context, userID uint) (*coreRepo.PaginatedEntities[*response.ProjectDto], error) {
	projects, err := s.repo.FindAllByUserID(ctx, userID)
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
	dto.URL = fmt.Sprintf("https://%s.statify.app", dto.Subdomain)
	latest, err := s.deploymentRepo.FindLatestByProjectID(ctx, dto.ID)
	if err == nil && latest != nil {
		dto.Status = latest.Status
	}
}
