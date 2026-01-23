package policy

import (
	"context"

	"github.com/gianghp/statify/internal/database/models"
)

type IAccessPolicy interface {
	CheckProjectAccess(ctx context.Context, userID uint, projectID uint) (*models.Project, error)
	CheckDeploymentAccess(ctx context.Context, userID uint, deploymentID uint) (*models.Project, *models.Deployment, error)
}
