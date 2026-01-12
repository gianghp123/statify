package service

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/gianghp/statify/internal/core"
	"github.com/gianghp/statify/internal/core/enums"
	coreRepo "github.com/gianghp/statify/internal/core/repository"
	"github.com/gianghp/statify/internal/database/models"
	deploymentRepo "github.com/gianghp/statify/internal/modules/deployment/repository"
	"github.com/gianghp/statify/internal/modules/project/dtos/request"
	"github.com/gianghp/statify/internal/modules/project/dtos/response"
	"github.com/gianghp/statify/internal/modules/project/repository"
	"github.com/gianghp/statify/internal/storage/minio"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestProjectService_CreateProject(t *testing.T) {
	tests := []struct {
		name         string
		userID       uint
		request      *request.CreateProjectRequest
		setupMocks   func(repo *repository.ProjectRepositoryMock, dRepo *deploymentRepo.DeploymentRepositoryMock, minioMock *minio.Mock)
		expectedFunc func(t *testing.T, project *response.ProjectDto, err error)
	}{
		{
			name:   "Create project successfully",
			userID: 1,
			request: &request.CreateProjectRequest{
				Name:      "Test Project",
				Subdomain: "test-project",
			},
			setupMocks: func(repo *repository.ProjectRepositoryMock, dRepo *deploymentRepo.DeploymentRepositoryMock, minioMock *minio.Mock) {
				repo.On("Create", mock.Anything, mock.MatchedBy(func(p *models.Project) bool {
					return p.Name == "Test Project" && p.Subdomain == "test-project" && p.UserID == 1
				})).Return(nil)
			},
			expectedFunc: func(t *testing.T, project *response.ProjectDto, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, project)
				assert.Equal(t, "Test Project", project.Name)
			},
		},
		{
			name:    "Create project failure",
			userID:  1,
			request: &request.CreateProjectRequest{Name: "Fail"},
			setupMocks: func(repo *repository.ProjectRepositoryMock, dRepo *deploymentRepo.DeploymentRepositoryMock, minioMock *minio.Mock) {
				repo.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error"))
			},
			expectedFunc: func(t *testing.T, project *response.ProjectDto, err error) {
				assert.Error(t, err)
				apiErr, ok := err.(*core.ApiError)
				assert.True(t, ok)
				assert.Equal(t, http.StatusInternalServerError, apiErr.Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(repository.ProjectRepositoryMock)
			dRepo := new(deploymentRepo.DeploymentRepositoryMock)
			minioMock := new(minio.Mock)
			tt.setupMocks(repo, dRepo, minioMock)
			s := NewProjectService(repo, dRepo, minioMock)
			project, err := s.CreateProject(context.TODO(), tt.userID, tt.request)
			tt.expectedFunc(t, project, err)
		})
	}
}

func TestProjectService_ListProjects(t *testing.T) {
	tests := []struct {
		name         string
		userID       uint
		setupMocks   func(repo *repository.ProjectRepositoryMock, dRepo *deploymentRepo.DeploymentRepositoryMock)
		expectedFunc func(t *testing.T, projects *coreRepo.PaginatedEntities[*response.ProjectDto], err error)
	}{
		{
			name:   "List projects successfully",
			userID: 1,
			setupMocks: func(repo *repository.ProjectRepositoryMock, dRepo *deploymentRepo.DeploymentRepositoryMock) {
				repo.On("FindAllByUserID", mock.Anything, uint(1)).Return(&coreRepo.PaginatedEntities[models.Project]{
					Entities: []models.Project{{Model: gorm.Model{ID: 1}, Name: "P1"}},
				}, nil)
				dRepo.On("FindLatestByProjectID", mock.Anything, uint(1)).Return(&models.Deployment{Status: enums.DeploymentStatusReady}, nil)
			},
			expectedFunc: func(t *testing.T, projects *coreRepo.PaginatedEntities[*response.ProjectDto], err error) {
				assert.NoError(t, err)
				assert.NotNil(t, projects)
				assert.Equal(t, enums.DeploymentStatusReady, projects.Entities[0].Status)
			},
		},
		{
			name:   "List projects successfully with no deployment",
			userID: 1,
			setupMocks: func(repo *repository.ProjectRepositoryMock, dRepo *deploymentRepo.DeploymentRepositoryMock) {
				repo.On("FindAllByUserID", mock.Anything, uint(1)).Return(&coreRepo.PaginatedEntities[models.Project]{
					Entities: []models.Project{{Model: gorm.Model{ID: 1}, Name: "P1"}},
				}, nil)
				dRepo.On("FindLatestByProjectID", mock.Anything, uint(1)).Return((*models.Deployment)(nil), gorm.ErrRecordNotFound)
			},
			expectedFunc: func(t *testing.T, projects *coreRepo.PaginatedEntities[*response.ProjectDto], err error) {
				assert.NoError(t, err)
				assert.NotNil(t, projects)
				assert.Equal(t, enums.DeploymentStatus(""), projects.Entities[0].Status)
			},
		},
		{
			name:   "List projects failure",
			userID: 1,
			setupMocks: func(repo *repository.ProjectRepositoryMock, dRepo *deploymentRepo.DeploymentRepositoryMock) {
				repo.On("FindAllByUserID", mock.Anything, uint(1)).Return((*coreRepo.PaginatedEntities[models.Project])(nil), errors.New("db error"))
			},
			expectedFunc: func(t *testing.T, projects *coreRepo.PaginatedEntities[*response.ProjectDto], err error) {
				assert.Error(t, err)
				assert.Equal(t, http.StatusInternalServerError, err.(*core.ApiError).Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(repository.ProjectRepositoryMock)
			dRepo := new(deploymentRepo.DeploymentRepositoryMock)
			minioMock := new(minio.Mock)
			tt.setupMocks(repo, dRepo)
			s := NewProjectService(repo, dRepo, minioMock)
			projects, err := s.ListProjects(context.TODO(), tt.userID)
			tt.expectedFunc(t, projects, err)
		})
	}
}

func TestProjectService_GetProjectByID(t *testing.T) {
	tests := []struct {
		name         string
		id           uint
		setupMocks   func(repo *repository.ProjectRepositoryMock, dRepo *deploymentRepo.DeploymentRepositoryMock)
		expectedFunc func(t *testing.T, project *response.ProjectDto, err error)
	}{
		{
			name: "Get project successfully",
			id:   1,
			setupMocks: func(repo *repository.ProjectRepositoryMock, dRepo *deploymentRepo.DeploymentRepositoryMock) {
				repo.On("FindByID", mock.Anything, uint(1)).Return(&models.Project{Model: gorm.Model{ID: 1}, Name: "P1", UserID: 1}, nil)
				dRepo.On("FindLatestByProjectID", mock.Anything, uint(1)).Return(&models.Deployment{Status: enums.DeploymentStatusReady}, nil)
			},
			expectedFunc: func(t *testing.T, project *response.ProjectDto, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "P1", project.Name)
				assert.Equal(t, enums.DeploymentStatusReady, project.Status)
			},
		},
		{
			name: "Project not found",
			id:   1,
			setupMocks: func(repo *repository.ProjectRepositoryMock, dRepo *deploymentRepo.DeploymentRepositoryMock) {
				repo.On("FindByID", mock.Anything, uint(1)).Return((*models.Project)(nil), gorm.ErrRecordNotFound)
			},
			expectedFunc: func(t *testing.T, project *response.ProjectDto, err error) {
				assert.Error(t, err)
				assert.Equal(t, http.StatusNotFound, err.(*core.ApiError).Code)
			},
		},
		{
			name: "Project not belong to user",
			id:   1,
			setupMocks: func(repo *repository.ProjectRepositoryMock, dRepo *deploymentRepo.DeploymentRepositoryMock) {
				repo.On("FindByID", mock.Anything, uint(1)).Return(&models.Project{Model: gorm.Model{ID: 1}, Name: "P1", UserID: 2}, nil)
			},
			expectedFunc: func(t *testing.T, project *response.ProjectDto, err error) {
				assert.Error(t, err)
				assert.Equal(t, http.StatusForbidden, err.(*core.ApiError).Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(repository.ProjectRepositoryMock)
			dRepo := new(deploymentRepo.DeploymentRepositoryMock)
			minioMock := new(minio.Mock)
			tt.setupMocks(repo, dRepo)
			s := NewProjectService(repo, dRepo, minioMock)
			project, err := s.GetProjectByID(context.TODO(), tt.id, 1)
			tt.expectedFunc(t, project, err)
		})
	}
}

func TestProjectService_DeleteProject(t *testing.T) {
	tests := []struct {
		name         string
		projectID    uint
		userID       uint
		setupMocks   func(repo *repository.ProjectRepositoryMock, dRepo *deploymentRepo.DeploymentRepositoryMock, minioMock *minio.Mock)
		expectedFunc func(t *testing.T, err error)
	}{
		{
			name:      "Delete project successfully",
			projectID: 1,
			userID:    1,
			setupMocks: func(repo *repository.ProjectRepositoryMock, dRepo *deploymentRepo.DeploymentRepositoryMock, minioMock *minio.Mock) {
				project := &models.Project{Model: gorm.Model{ID: 1}, UserID: 1}
				repo.On("FindByID", mock.Anything, uint(1)).Return(project, nil)
				dRepo.On("FindAllByProjectID", mock.Anything, uint(1)).Return(&coreRepo.PaginatedEntities[models.Deployment]{
					Entities: []models.Deployment{{Model: gorm.Model{ID: 10}, OutputPrefix: "deployments/1/10"}},
				}, nil)
				repo.On("Delete", mock.Anything, project).Return(nil)
				minioMock.On("RemoveObjectsByPrefix", mock.Anything, "static-sites", "deployments/1/10").Return(nil)

				// Mock Transaction and WithTx
				repo.On("Transaction", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
					fn := args.Get(1).(func(*gorm.DB) error)
					_ = fn(nil)
				})
				repo.On("WithTx", mock.Anything).Return(repo)
				dRepo.On("WithTx", mock.Anything).Return(dRepo)
			},
			expectedFunc: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(repository.ProjectRepositoryMock)
			dRepo := new(deploymentRepo.DeploymentRepositoryMock)
			minioMock := new(minio.Mock)
			tt.setupMocks(repo, dRepo, minioMock)
			s := NewProjectService(repo, dRepo, minioMock)
			err := s.DeleteProject(context.TODO(), tt.projectID, tt.userID)
			tt.expectedFunc(t, err)

			repo.AssertExpectations(t)
			dRepo.AssertExpectations(t)
			minioMock.AssertExpectations(t)
		})
	}
}
