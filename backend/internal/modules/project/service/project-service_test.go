package service

import (
	"errors"
	"net/http"
	"testing"

	"github.com/gianghp/statify/internal/core"
	coreRepo "github.com/gianghp/statify/internal/core/repository"
	"github.com/gianghp/statify/internal/database/models"
	"github.com/gianghp/statify/internal/modules/project/dtos/request"
	"github.com/gianghp/statify/internal/modules/project/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestProjectService_CreateProject(t *testing.T) {
	tests := []struct {
		name         string
		userID       uint
		request      *request.CreateProjectRequest
		setupMocks   func(repo *repository.ProjectRepositoryMock)
		expectedFunc func(t *testing.T, project *models.Project, err error)
	}{
		{
			name:   "Create project successfully",
			userID: 1,
			request: &request.CreateProjectRequest{
				Name:      "Test Project",
				Subdomain: "test-project",
			},
			setupMocks: func(repo *repository.ProjectRepositoryMock) {
				repo.On("Create", mock.MatchedBy(func(p *models.Project) bool {
					return p.Name == "Test Project" && p.Subdomain == "test-project" && p.UserID == 1
				})).Return(nil)
			},
			expectedFunc: func(t *testing.T, project *models.Project, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, project)
				assert.Equal(t, "Test Project", project.Name)
			},
		},
		{
			name:    "Create project failure",
			userID:  1,
			request: &request.CreateProjectRequest{Name: "Fail"},
			setupMocks: func(repo *repository.ProjectRepositoryMock) {
				repo.On("Create", mock.Anything).Return(errors.New("db error"))
			},
			expectedFunc: func(t *testing.T, project *models.Project, err error) {
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
			tt.setupMocks(repo)
			s := NewProjectService(repo)
			project, err := s.CreateProject(tt.userID, tt.request)
			tt.expectedFunc(t, project, err)
		})
	}
}

func TestProjectService_ListProjects(t *testing.T) {
	tests := []struct {
		name         string
		userID       uint
		setupMocks   func(repo *repository.ProjectRepositoryMock)
		expectedFunc func(t *testing.T, projects *coreRepo.PaginatedEntities[models.Project], err error)
	}{
		{
			name:   "List projects successfully",
			userID: 1,
			setupMocks: func(repo *repository.ProjectRepositoryMock) {
				repo.On("FindAllByUserID", uint(1)).Return(&coreRepo.PaginatedEntities[models.Project]{
					Entities: []models.Project{{Name: "P1"}},
				}, nil)
			},
			expectedFunc: func(t *testing.T, projects *coreRepo.PaginatedEntities[models.Project], err error) {
				assert.NoError(t, err)
				assert.NotNil(t, projects)
			},
		},
		{
			name:   "List projects failure",
			userID: 1,
			setupMocks: func(repo *repository.ProjectRepositoryMock) {
				repo.On("FindAllByUserID", uint(1)).Return((*coreRepo.PaginatedEntities[models.Project])(nil), errors.New("db error"))
			},
			expectedFunc: func(t *testing.T, projects *coreRepo.PaginatedEntities[models.Project], err error) {
				assert.Error(t, err)
				assert.Equal(t, http.StatusInternalServerError, err.(*core.ApiError).Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(repository.ProjectRepositoryMock)
			tt.setupMocks(repo)
			s := NewProjectService(repo)
			projects, err := s.ListProjects(tt.userID)
			tt.expectedFunc(t, projects, err)
		})
	}
}

func TestProjectService_GetProjectByID(t *testing.T) {
	tests := []struct {
		name         string
		id           uint
		setupMocks   func(repo *repository.ProjectRepositoryMock)
		expectedFunc func(t *testing.T, project *models.Project, err error)
	}{
		{
			name: "Get project successfully",
			id:   1,
			setupMocks: func(repo *repository.ProjectRepositoryMock) {
				repo.On("FindByID", uint(1)).Return(&models.Project{Name: "P1"}, nil)
			},
			expectedFunc: func(t *testing.T, project *models.Project, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "P1", project.Name)
			},
		},
		{
			name: "Project not found",
			id:   1,
			setupMocks: func(repo *repository.ProjectRepositoryMock) {
				repo.On("FindByID", uint(1)).Return((*models.Project)(nil), gorm.ErrRecordNotFound)
			},
			expectedFunc: func(t *testing.T, project *models.Project, err error) {
				assert.Error(t, err)
				assert.Equal(t, http.StatusNotFound, err.(*core.ApiError).Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(repository.ProjectRepositoryMock)
			tt.setupMocks(repo)
			s := NewProjectService(repo)
			project, err := s.GetProjectByID(tt.id)
			tt.expectedFunc(t, project, err)
		})
	}
}
