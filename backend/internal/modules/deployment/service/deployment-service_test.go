package service

import (
	"errors"
	"net/http"
	"testing"

	"github.com/gianghp/statify/internal/core"
	coreRepo "github.com/gianghp/statify/internal/core/repository"
	"github.com/gianghp/statify/internal/database/models"
	"github.com/gianghp/statify/internal/modules/deployment/dtos/request"
	"github.com/gianghp/statify/internal/modules/deployment/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestDeploymentService_CreateDeployment(t *testing.T) {
	tests := []struct {
		name         string
		projectID    uint
		request      *request.CreateDeploymentRequest
		setupMocks   func(repo *repository.DeploymentRepositoryMock)
		expectedFunc func(t *testing.T, deployment *models.Deployment, err error)
	}{
		{
			name:      "Create deployment successfully",
			projectID: 1,
			request:   &request.CreateDeploymentRequest{},
			setupMocks: func(repo *repository.DeploymentRepositoryMock) {
				repo.On("Create", mock.Anything).Return(nil)
			},
			expectedFunc: func(t *testing.T, deployment *models.Deployment, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, deployment)
			},
		},
		{
			name:      "Create deployment failure",
			projectID: 1,
			setupMocks: func(repo *repository.DeploymentRepositoryMock) {
				repo.On("Create", mock.Anything).Return(errors.New("db error"))
			},
			expectedFunc: func(t *testing.T, deployment *models.Deployment, err error) {
				assert.Error(t, err)
				assert.Equal(t, http.StatusInternalServerError, err.(*core.ApiError).Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(repository.DeploymentRepositoryMock)
			tt.setupMocks(repo)
			s := NewDeploymentService(repo)
			deployment, err := s.CreateDeployment(tt.projectID, tt.request)
			tt.expectedFunc(t, deployment, err)
		})
	}
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
					Entities: []models.Deployment{{ProjectID: 1}},
				}, nil)
			},
			expectedFunc: func(t *testing.T, deployments *coreRepo.PaginatedEntities[models.Deployment], err error) {
				assert.NoError(t, err)
				assert.NotNil(t, deployments)
			},
		},
		{
			name:      "Get history failure",
			projectID: 1,
			setupMocks: func(repo *repository.DeploymentRepositoryMock) {
				repo.On("FindAllByProjectID", uint(1)).Return((*coreRepo.PaginatedEntities[models.Deployment])(nil), errors.New("db error"))
			},
			expectedFunc: func(t *testing.T, deployments *coreRepo.PaginatedEntities[models.Deployment], err error) {
				assert.Error(t, err)
				assert.Equal(t, http.StatusInternalServerError, err.(*core.ApiError).Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(repository.DeploymentRepositoryMock)
			tt.setupMocks(repo)
			s := NewDeploymentService(repo)
			deployments, err := s.GetHistory(tt.projectID)
			tt.expectedFunc(t, deployments, err)
		})
	}
}

func TestDeploymentService_GetDeploymentByID(t *testing.T) {
	tests := []struct {
		name         string
		id           uint
		setupMocks   func(repo *repository.DeploymentRepositoryMock)
		expectedFunc func(t *testing.T, deployment *models.Deployment, err error)
	}{
		{
			name: "Get deployment successfully",
			id:   1,
			setupMocks: func(repo *repository.DeploymentRepositoryMock) {
				repo.On("FindByID", uint(1)).Return(&models.Deployment{ProjectID: 1}, nil)
			},
			expectedFunc: func(t *testing.T, deployment *models.Deployment, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, deployment)
			},
		},
		{
			name: "Deployment not found",
			id:   1,
			setupMocks: func(repo *repository.DeploymentRepositoryMock) {
				repo.On("FindByID", uint(1)).Return((*models.Deployment)(nil), gorm.ErrRecordNotFound)
			},
			expectedFunc: func(t *testing.T, deployment *models.Deployment, err error) {
				assert.Error(t, err)
				assert.Equal(t, http.StatusNotFound, err.(*core.ApiError).Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(repository.DeploymentRepositoryMock)
			tt.setupMocks(repo)
			s := NewDeploymentService(repo)
			deployment, err := s.GetDeploymentByID(tt.id)
			tt.expectedFunc(t, deployment, err)
		})
	}
}
