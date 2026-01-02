package service

import (
	"testing"

	coreRepo "github.com/gianghp/statify/internal/core/repository"
	"github.com/gianghp/statify/internal/database/models"
	"github.com/gianghp/statify/internal/modules/project/dtos/request"
	"github.com/gianghp/statify/internal/modules/project/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
				assert.Equal(t, "test-project", project.Subdomain)
				assert.Equal(t, uint(1), project.UserID)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(repository.ProjectRepositoryMock)
			if tt.setupMocks != nil {
				tt.setupMocks(repo)
			}

			s := NewProjectService(repo)
			project, err := s.CreateProject(tt.userID, tt.request)

			tt.expectedFunc(t, project, err)
			repo.AssertExpectations(t)
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
					Entities: []models.Project{
						{Name: "Project 1", UserID: 1},
						{Name: "Project 2", UserID: 1},
					},
					Pagination: coreRepo.Pagination{TotalCount: 2},
				}, nil)
			},
			expectedFunc: func(t *testing.T, projects *coreRepo.PaginatedEntities[models.Project], err error) {
				assert.NoError(t, err)
				assert.Len(t, projects.Entities, 2)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(repository.ProjectRepositoryMock)
			if tt.setupMocks != nil {
				tt.setupMocks(repo)
			}

			s := NewProjectService(repo)
			projects, err := s.ListProjects(tt.userID)

			tt.expectedFunc(t, projects, err)
			repo.AssertExpectations(t)
		})
	}
}
