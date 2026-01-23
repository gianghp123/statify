package policy

import (
	"context"

	"github.com/gianghp/statify/internal/database/models"
	"github.com/stretchr/testify/mock"
)

type AccessPolicyMock struct {
	mock.Mock
}

func (m *AccessPolicyMock) CheckProjectAccess(ctx context.Context, userID uint, projectID uint) (*models.Project, error) {
	args := m.Called(ctx, userID, projectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Project), args.Error(1)
}

func (m *AccessPolicyMock) CheckDeploymentAccess(ctx context.Context, userID uint, deploymentID uint) (*models.Project, *models.Deployment, error) {
	args := m.Called(ctx, userID, deploymentID)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	if args.Get(1) == nil {
		return args.Get(0).(*models.Project), nil, args.Error(2)
	}
	return args.Get(0).(*models.Project), args.Get(1).(*models.Deployment), args.Error(2)
}
