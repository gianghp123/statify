package service

import (
	"testing"

	"github.com/gianghp/statify/internal/core/enums"
	coreRepo "github.com/gianghp/statify/internal/core/repository"
	"github.com/gianghp/statify/internal/database/models"
	"github.com/gianghp/statify/internal/modules/deployment/repository"
	"github.com/stretchr/testify/assert"
)

func TestDeploymentService_CreateDeployment(t *testing.T) {
	// ... (no changes needed)
}

func TestDeploymentService_GetHistory(t *testing.T) {
	tests := []struct {
		name         string
		projectID    uint
		setupMocks   func(repo *repository.DeploymentRepositoryMock)
		expectedFunc func(t *testing.T, deployments *coreRepo.PaginatedEntities[models.Deployment], err error)
	}{
		{
			name:      "Get history successfully",
			projectID: 1,
			setupMocks: func(repo *repository.DeploymentRepositoryMock) {
				repo.On("FindAllByProjectID", uint(1)).Return(&coreRepo.PaginatedEntities[models.Deployment]{
					Entities: []models.Deployment{
						{ProjectID: 1, Status: enums.DeploymentStatusReady},
						{ProjectID: 1, Status: enums.DeploymentStatusFailed},
					},
					Pagination: coreRepo.Pagination{TotalCount: 2},
				}, nil)
			},
			expectedFunc: func(t *testing.T, deployments *coreRepo.PaginatedEntities[models.Deployment], err error) {
				assert.NoError(t, err)
				assert.Len(t, deployments.Entities, 2)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(repository.DeploymentRepositoryMock)
			if tt.setupMocks != nil {
				tt.setupMocks(repo)
			}

			s := NewDeploymentService(repo)
			deployments, err := s.GetHistory(tt.projectID)

			tt.expectedFunc(t, deployments, err)
			repo.AssertExpectations(t)
		})
	}
}
