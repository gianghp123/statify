package service

import (
	coreRepo "github.com/gianghp/statify/internal/core/repository"
	"github.com/gianghp/statify/internal/database/models"
	"github.com/gianghp/statify/internal/modules/project/dtos/request"
	"github.com/gianghp/statify/internal/modules/project/repository"
)

type ProjectService struct {
	repo repository.IProjectRepository
}

func NewProjectService(repo repository.IProjectRepository) *ProjectService {
	return &ProjectService{repo: repo}
}

func (s *ProjectService) CreateProject(userID uint, req *request.CreateProjectRequest) (*models.Project, error) {
	project := &models.Project{
		Name:      req.Name,
		Subdomain: req.Subdomain,
		UserID:    userID,
	}

	if err := s.repo.Create(project); err != nil {
		return nil, err
	}

	return project, nil
}

func (s *ProjectService) ListProjects(userID uint) (*coreRepo.PaginatedEntities[models.Project], error) {
	return s.repo.FindAllByUserID(userID)
}

func (s *ProjectService) GetProjectByID(id uint) (*models.Project, error) {
	return s.repo.FindByID(id)
}
