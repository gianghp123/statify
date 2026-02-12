package policy

import (
	"context"

	"github.com/gianghp/statify/internal/core"
	"github.com/gianghp/statify/internal/database/models"
	deploymentRepo "github.com/gianghp/statify/internal/modules/deployment/repository"
	projectRepo "github.com/gianghp/statify/internal/modules/project/repository"
)

type AccessPolicy struct {
	projectRepository    projectRepo.IProjectRepository
	deploymentRepository deploymentRepo.IDeploymentRepository
}

func NewAccessPolicy(
	projectRepository projectRepo.IProjectRepository,
	deploymentRepository deploymentRepo.IDeploymentRepository,
) *AccessPolicy {
	return &AccessPolicy{
		projectRepository:    projectRepository,
		deploymentRepository: deploymentRepository,
	}
}

func (p *AccessPolicy) CheckProjectAccess(ctx context.Context, userID uint, projectID uint) (*models.Project, error) {
	project, err := p.projectRepository.FindByID(ctx, projectID)

	if err != nil {
		return nil, core.NotFoundError()
	}

	if project.UserID != userID {
		return nil, core.ForbiddenError()
	}

	return project, nil
}

func (p *AccessPolicy) CheckDeploymentAccess(ctx context.Context, userID uint, deploymentID uint) (*models.Project, *models.Deployment, error) {
	deployment, err := p.deploymentRepository.FindByID(ctx, deploymentID)
	if err != nil {
		return nil, nil, core.NotFoundError()
	}

	if deployment.Project.UserID != userID {
		return nil, nil, core.ForbiddenError()
	}

	return &deployment.Project, deployment, nil
}
