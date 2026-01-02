package service

import (
	"github.com/gianghp/statify/internal/core/enums"
	coreRepo "github.com/gianghp/statify/internal/core/repository"
	"github.com/gianghp/statify/internal/database/models"
	"github.com/gianghp/statify/internal/modules/deployment/dtos/request"
	"github.com/gianghp/statify/internal/modules/deployment/repository"
)

type DeploymentService struct {
	repo repository.IDeploymentRepository
}

func NewDeploymentService(repo repository.IDeploymentRepository) *DeploymentService {
	return &DeploymentService{repo: repo}
}

func (s *DeploymentService) CreateDeployment(projectID uint, req *request.CreateDeploymentRequest) (*models.Deployment, error) {
	deployment := &models.Deployment{
		ProjectID: projectID,
		Status:    enums.DeploymentStatusQueued,
	}

	if err := s.repo.Create(deployment); err != nil {
		return nil, err
	}

	return deployment, nil
}

func (s *DeploymentService) GetHistory(projectID uint) (*coreRepo.PaginatedEntities[models.Deployment], error) {
	return s.repo.FindAllByProjectID(projectID)
}

func (s *DeploymentService) GetDeploymentByID(id uint) (*models.Deployment, error) {
	return s.repo.FindByID(id)
}
