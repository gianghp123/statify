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
	"github.com/gianghp/statify/internal/modules/auth/policy"
	deploymentRepo "github.com/gianghp/statify/internal/modules/deployment/repository"
	jobQueueRepo "github.com/gianghp/statify/internal/modules/job-queue/repository"
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
		setupMocks   func(repo *repository.ProjectRepositoryMock, dRepo *deploymentRepo.DeploymentRepositoryMock, jRepo *jobQueueRepo.JobQueueRepositoryMock, minioMock *minio.Mock, policyMock *policy.AccessPolicyMock)
		expectedFunc func(t *testing.T, project *response.ProjectDto, err error)
	}{
		{
			name:   "Create project successfully",
			userID: 1,
			request: &request.CreateProjectRequest{
				Name:      "Test Project",
				Subdomain: "test-project",
			},
			setupMocks: func(repo *repository.ProjectRepositoryMock, dRepo *deploymentRepo.DeploymentRepositoryMock, jRepo *jobQueueRepo.JobQueueRepositoryMock, minioMock *minio.Mock, policyMock *policy.AccessPolicyMock) {
				repo.On("FindBySubdomain", mock.Anything, "test-project").Return((*models.Project)(nil), gorm.ErrRecordNotFound)
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
			request: &request.CreateProjectRequest{Name: "Fail", Subdomain: "test-project"},
			setupMocks: func(repo *repository.ProjectRepositoryMock, dRepo *deploymentRepo.DeploymentRepositoryMock, jRepo *jobQueueRepo.JobQueueRepositoryMock, minioMock *minio.Mock, policyMock *policy.AccessPolicyMock) {
				repo.On("FindBySubdomain", mock.Anything, "test-project").Return((*models.Project)(nil), gorm.ErrRecordNotFound)
				repo.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error"))
			},
			expectedFunc: func(t *testing.T, project *response.ProjectDto, err error) {
				assert.Error(t, err)
				apiErr, ok := err.(*core.ApiError)
				assert.True(t, ok)
				assert.Equal(t, http.StatusInternalServerError, apiErr.Code)
			},
		},
		{
			name:   "Subdomain already exists",
			userID: 1,
			request: &request.CreateProjectRequest{
				Name:      "Test Project",
				Subdomain: "test-project",
			},
			setupMocks: func(repo *repository.ProjectRepositoryMock, dRepo *deploymentRepo.DeploymentRepositoryMock, jRepo *jobQueueRepo.JobQueueRepositoryMock, minioMock *minio.Mock, policyMock *policy.AccessPolicyMock) {
				repo.On("FindBySubdomain", mock.Anything, "test-project").Return(&models.Project{}, nil)
			},
			expectedFunc: func(t *testing.T, project *response.ProjectDto, err error) {
				assert.Error(t, err)
				apiErr, ok := err.(*core.ApiError)
				assert.True(t, ok)
				assert.Equal(t, http.StatusBadRequest, apiErr.Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(repository.ProjectRepositoryMock)
			dRepo := new(deploymentRepo.DeploymentRepositoryMock)
			jRepo := new(jobQueueRepo.JobQueueRepositoryMock)
			minioMock := new(minio.Mock)
			policyMock := new(policy.AccessPolicyMock)
			tt.setupMocks(repo, dRepo, jRepo, minioMock, policyMock)
			s := NewProjectService(repo, dRepo, jRepo, minioMock, policyMock)
			project, err := s.CreateProject(context.TODO(), tt.userID, tt.request)
			tt.expectedFunc(t, project, err)
		})
	}
}

func TestProjectService_ListProjects(t *testing.T) {
	tests := []struct {
		name             string
		userID           uint
		deploymentStatus enums.DeploymentStatus
		setupMocks       func(repo *repository.ProjectRepositoryMock, dRepo *deploymentRepo.DeploymentRepositoryMock, jRepo *jobQueueRepo.JobQueueRepositoryMock, policyMock *policy.AccessPolicyMock)
		expectedFunc     func(t *testing.T, projects coreRepo.PaginatedEntities[*response.ProjectDto], err error)
	}{
		{
			name:             "List projects successfully",
			userID:           1,
			deploymentStatus: "",
			setupMocks: func(repo *repository.ProjectRepositoryMock, dRepo *deploymentRepo.DeploymentRepositoryMock, jRepo *jobQueueRepo.JobQueueRepositoryMock, policyMock *policy.AccessPolicyMock) {
				repo.On("FindAllByUserID", context.TODO(), uint(1), 1, 10, enums.DeploymentStatus("")).Return(coreRepo.PaginatedEntities[*models.Project]{
					Entities: []*models.Project{{Model: gorm.Model{ID: 1}, Name: "P1"}},
				}, nil)
			},
			expectedFunc: func(t *testing.T, projects coreRepo.PaginatedEntities[*response.ProjectDto], err error) {
				assert.NoError(t, err)
				assert.NotNil(t, projects)
			},
		},
		{
			name:             "List live project successfully",
			userID:           1,
			deploymentStatus: enums.DeploymentStatusLive,
			setupMocks: func(repo *repository.ProjectRepositoryMock, dRepo *deploymentRepo.DeploymentRepositoryMock, jRepo *jobQueueRepo.JobQueueRepositoryMock, policyMock *policy.AccessPolicyMock) {
				repo.On("FindAllByUserID", context.TODO(), uint(1), 1, 10, enums.DeploymentStatusLive).Return(coreRepo.PaginatedEntities[*models.Project]{
					Entities: []*models.Project{{Model: gorm.Model{ID: 1}, Name: "P1", Deployments: []models.Deployment{{Status: enums.DeploymentStatusLive}}}},
				}, nil)
			},
			expectedFunc: func(t *testing.T, projects coreRepo.PaginatedEntities[*response.ProjectDto], err error) {
				assert.NoError(t, err)
				assert.NotNil(t, projects)
				assert.Equal(t, enums.DeploymentStatusLive, projects.Entities[0].Status)
			},
		},
		{
			name:             "List projects successfully with no deployment",
			userID:           1,
			deploymentStatus: "",
			setupMocks: func(repo *repository.ProjectRepositoryMock, dRepo *deploymentRepo.DeploymentRepositoryMock, jRepo *jobQueueRepo.JobQueueRepositoryMock, policyMock *policy.AccessPolicyMock) {
				repo.On("FindAllByUserID", context.TODO(), uint(1), 1, 10, enums.DeploymentStatus("")).Return(coreRepo.PaginatedEntities[*models.Project]{
					Entities: []*models.Project{{Model: gorm.Model{ID: 1}, Name: "P1"}},
				}, nil)
			},
			expectedFunc: func(t *testing.T, projects coreRepo.PaginatedEntities[*response.ProjectDto], err error) {
				assert.NoError(t, err)
				assert.NotNil(t, projects)
			},
		},
		{
			name:             "List projects failure",
			userID:           1,
			deploymentStatus: "",
			setupMocks: func(repo *repository.ProjectRepositoryMock, dRepo *deploymentRepo.DeploymentRepositoryMock, jRepo *jobQueueRepo.JobQueueRepositoryMock, policyMock *policy.AccessPolicyMock) {
				repo.On("FindAllByUserID", mock.Anything, uint(1), 1, 10, enums.DeploymentStatus("")).Return(coreRepo.PaginatedEntities[*models.Project]{}, errors.New("db error"))
			},
			expectedFunc: func(t *testing.T, projects coreRepo.PaginatedEntities[*response.ProjectDto], err error) {
				assert.Error(t, err)
				assert.Equal(t, http.StatusInternalServerError, err.(*core.ApiError).Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(repository.ProjectRepositoryMock)
			dRepo := new(deploymentRepo.DeploymentRepositoryMock)
			jRepo := new(jobQueueRepo.JobQueueRepositoryMock)
			minioMock := new(minio.Mock)
			policyMock := new(policy.AccessPolicyMock)
			tt.setupMocks(repo, dRepo, jRepo, policyMock)
			s := NewProjectService(repo, dRepo, jRepo, minioMock, policyMock)
			projects, err := s.ListProjects(context.TODO(), tt.userID, 1, 10, tt.deploymentStatus)
			tt.expectedFunc(t, projects, err)
		})
	}
}

func TestProjectService_GetProjectByID(t *testing.T) {
	tests := []struct {
		name             string
		id               uint
		deploymentStatus enums.DeploymentStatus
		setupMocks       func(repo *repository.ProjectRepositoryMock, jRepo *jobQueueRepo.JobQueueRepositoryMock, policyMock *policy.AccessPolicyMock)
		expectedFunc     func(t *testing.T, project *response.ProjectDto, err error)
	}{
		{
			name:             "Get project successfully",
			id:               1,
			deploymentStatus: "",
			setupMocks: func(repo *repository.ProjectRepositoryMock, jRepo *jobQueueRepo.JobQueueRepositoryMock, policyMock *policy.AccessPolicyMock) {
				policyMock.On("CheckProjectAccess", mock.Anything, uint(1), uint(1)).Return(&models.Project{Model: gorm.Model{ID: 1}, Name: "P1", UserID: 1, Deployments: []models.Deployment{{Status: enums.DeploymentStatusReady}}, Subdomain: "p1"}, nil)
			},
			expectedFunc: func(t *testing.T, project *response.ProjectDto, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "P1", project.Name)
				assert.Equal(t, enums.DeploymentStatusReady, project.Status)
				assert.Equal(t, "https://p1.statify.app", project.URL)
			},
		},
		{
			name:             "Project not found",
			id:               1,
			deploymentStatus: "",
			setupMocks: func(repo *repository.ProjectRepositoryMock, jRepo *jobQueueRepo.JobQueueRepositoryMock, policyMock *policy.AccessPolicyMock) {
				policyMock.On("CheckProjectAccess", mock.Anything, uint(1), uint(1)).Return((*models.Project)(nil), core.NotFoundError())
			},
			expectedFunc: func(t *testing.T, project *response.ProjectDto, err error) {
				assert.Error(t, err)
				assert.Equal(t, http.StatusNotFound, err.(*core.ApiError).Code)
			},
		},
		{
			name:             "Project not belong to user",
			id:               1,
			deploymentStatus: "",
			setupMocks: func(repo *repository.ProjectRepositoryMock, jRepo *jobQueueRepo.JobQueueRepositoryMock, policyMock *policy.AccessPolicyMock) {
				policyMock.On("CheckProjectAccess", mock.Anything, uint(1), uint(1)).Return((*models.Project)(nil), core.ForbiddenError())
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
			jRepo := new(jobQueueRepo.JobQueueRepositoryMock)
			minioMock := new(minio.Mock)
			policyMock := new(policy.AccessPolicyMock)
			tt.setupMocks(repo, jRepo, policyMock)
			s := NewProjectService(repo, nil, jRepo, minioMock, policyMock)
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
		setupMocks   func(repo *repository.ProjectRepositoryMock, dRepo *deploymentRepo.DeploymentRepositoryMock, jRepo *jobQueueRepo.JobQueueRepositoryMock, minioMock *minio.Mock, policyMock *policy.AccessPolicyMock)
		expectedFunc func(t *testing.T, err error)
	}{
		{
			name:      "Delete project successfully",
			projectID: 1,
			userID:    1,
			setupMocks: func(repo *repository.ProjectRepositoryMock, dRepo *deploymentRepo.DeploymentRepositoryMock, jRepo *jobQueueRepo.JobQueueRepositoryMock, minioMock *minio.Mock, policyMock *policy.AccessPolicyMock) {
				project := &models.Project{Model: gorm.Model{ID: 1}, UserID: 1}
				policyMock.On("CheckProjectAccess", mock.Anything, uint(1), uint(1)).Return(project, nil)
				repo.On("Delete", mock.Anything, project).Return(nil)
				jRepo.On("Create", mock.Anything, mock.MatchedBy(func(j *models.JobQueue) bool {
					return j.Type == enums.JobQueueTypeProjectDelete && j.Payload == "1"
				})).Return(nil)
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
			jRepo := new(jobQueueRepo.JobQueueRepositoryMock)
			minioMock := new(minio.Mock)
			policyMock := new(policy.AccessPolicyMock)
			tt.setupMocks(repo, dRepo, jRepo, minioMock, policyMock)
			s := NewProjectService(repo, dRepo, jRepo, minioMock, policyMock)
			err := s.DeleteProject(context.TODO(), tt.projectID, tt.userID)
			tt.expectedFunc(t, err)

			repo.AssertExpectations(t)
			dRepo.AssertExpectations(t)
			minioMock.AssertExpectations(t)
			jRepo.AssertExpectations(t)
		})
	}
}
