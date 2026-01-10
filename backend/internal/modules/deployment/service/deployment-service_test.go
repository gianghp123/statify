package service

import (
	"archive/zip"
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"

	"github.com/gianghp/statify/internal/core"
	"github.com/gianghp/statify/internal/core/enums"
	coreRepo "github.com/gianghp/statify/internal/core/repository"
	"github.com/gianghp/statify/internal/database/models"
	"github.com/gianghp/statify/internal/modules/deployment/dtos/request"
	"github.com/gianghp/statify/internal/modules/deployment/dtos/response"
	"github.com/gianghp/statify/internal/modules/deployment/repository"
	projectRepository "github.com/gianghp/statify/internal/modules/project/repository"
	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func createMockMultipartFileHeader(t *testing.T, filenames []string) (*multipart.FileHeader, []byte) {
	// 1. Create a ZIP in memory
	zipBuf := new(bytes.Buffer)
	zw := zip.NewWriter(zipBuf)
	for _, name := range filenames {
		f, _ := zw.Create(name)
		f.Write([]byte("dummy content for " + name))
	}
	zw.Close()
	zipBytes := zipBuf.Bytes()

	// 2. Create a Multipart Form in memory to get a *multipart.FileHeader
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test.zip")
	if err != nil {
		t.Fatal(err)
	}
	part.Write(zipBytes)
	writer.Close()

	// 3. Parse it back to extract the FileHeader
	reader := multipart.NewReader(body, writer.Boundary())
	form, err := reader.ReadForm(1 << 20) // 1MB limit
	if err != nil {
		t.Fatal(err)
	}

	return form.File["file"][0], zipBytes
}

func TestDeploymentService_CreateDeployment(t *testing.T) {
	fileHeader, _ := createMockMultipartFileHeader(t, []string{"index.html", "css/style.css"})

	tests := []struct {
		name         string
		userID       uint
		projectID    uint
		request      *request.CreateDeploymentRequest
		setupMocks   func(repo *repository.DeploymentRepositoryMock, projectRepo *projectRepository.ProjectRepositoryMock, minioClient *MinioMock)
		expectedFunc func(t *testing.T, deployment *response.DeploymentDto, err error)
	}{
		{
			name:      "Create deployment successfully",
			userID:    1,
			projectID: 1,
			request: &request.CreateDeploymentRequest{
				File: fileHeader,
			},
			setupMocks: func(repo *repository.DeploymentRepositoryMock, projectRepo *projectRepository.ProjectRepositoryMock, minioClient *MinioMock) {
				project := &models.Project{
					Model:     gorm.Model{ID: 1},
					UserID:    1,
					Subdomain: "test-site",
				}
				projectRepo.On("FindByID", mock.Anything, uint(1)).Return(project, nil)

				// 1. Expect initial DB creation (Status: Uploaded/Processing)
				repo.On("Create", mock.Anything, mock.MatchedBy(func(d *models.Deployment) bool {
					return d.ProjectID == 1 && d.Status == enums.DeploymentStatusUploaded
				})).Return(nil).Run(func(args mock.Arguments) {
					// Simulate GORM setting the ID after creation
					d := args.Get(1).(*models.Deployment)
					d.ID = 10
				})

				// 2. Expect PutObject for each file
				minioClient.On("PutObject", mock.Anything, "static-sites",
					mock.MatchedBy(func(s string) bool { return strings.Contains(s, "index.html") }),
					mock.Anything, mock.Anything, mock.Anything,
				).Return(minio.UploadInfo{}, nil)

				minioClient.On("PutObject", mock.Anything, "static-sites",
					mock.MatchedBy(func(s string) bool { return strings.Contains(s, "css/style.css") }),
					mock.Anything, mock.Anything, mock.Anything,
				).Return(minio.UploadInfo{}, nil)

				// 3. Expect Deployment status update to Ready
				repo.On("Update", mock.Anything, mock.MatchedBy(func(d *models.Deployment) bool {
					return d.ID == 10 && d.Status == enums.DeploymentStatusReady
				})).Return(nil)

				// 4. Expect Project update to set CurrentDeploymentID
				projectRepo.On("Update", mock.Anything, mock.MatchedBy(func(p *models.Project) bool {
					return p.ID == 1 && p.CurrentDeploymentID == 10
				})).Return(nil)
			},
			expectedFunc: func(t *testing.T, deployment *response.DeploymentDto, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, deployment)
				assert.Equal(t, string(enums.DeploymentStatusReady), deployment.Status)
			},
		},
		{
			name:      "Create deployment failure - forbidden",
			userID:    2, // Wrong user
			projectID: 1,
			request:   &request.CreateDeploymentRequest{File: fileHeader},
			setupMocks: func(repo *repository.DeploymentRepositoryMock, projectRepo *projectRepository.ProjectRepositoryMock, minioClient *MinioMock) {
				projectRepo.On("FindByID", mock.Anything, uint(1)).Return(&models.Project{
					Model:  gorm.Model{ID: 1},
					UserID: 1,
				}, nil)
			},
			expectedFunc: func(t *testing.T, deployment *response.DeploymentDto, err error) {
				assert.Error(t, err)
				apiErr, ok := err.(*core.ApiError)
				assert.True(t, ok)
				assert.Equal(t, http.StatusForbidden, apiErr.Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(repository.DeploymentRepositoryMock)
			projectRepo := new(projectRepository.ProjectRepositoryMock)
			minioClient := new(MinioMock)

			tt.setupMocks(repo, projectRepo, minioClient)

			s := NewDeploymentService(repo, projectRepo, minioClient)
			deployment, err := s.CreateDeployment(context.TODO(), tt.userID, tt.projectID, tt.request)

			tt.expectedFunc(t, deployment, err)

			repo.AssertExpectations(t)
			projectRepo.AssertExpectations(t)
			minioClient.AssertExpectations(t)
		})
	}
}

func TestDeploymentService_GetHistory(t *testing.T) {
	tests := []struct {
		name         string
		userID       uint
		projectID    uint
		setupMocks   func(repo *repository.DeploymentRepositoryMock, projectRepo *projectRepository.ProjectRepositoryMock)
		expectedFunc func(t *testing.T, deployments *coreRepo.PaginatedEntities[*response.DeploymentDto], err error)
	}{
		{
			name:      "Get history successfully",
			userID:    1,
			projectID: 1,
			setupMocks: func(repo *repository.DeploymentRepositoryMock, projectRepo *projectRepository.ProjectRepositoryMock) {
				projectRepo.On("FindByID", mock.Anything, uint(1)).Return(&models.Project{Model: gorm.Model{ID: 1}, UserID: 1}, nil)
				repo.On("FindAllByProjectID", mock.Anything, uint(1)).Return(&coreRepo.PaginatedEntities[models.Deployment]{
					Entities: []models.Deployment{{ProjectID: 1}},
				}, nil)
			},
			expectedFunc: func(t *testing.T, deployments *coreRepo.PaginatedEntities[*response.DeploymentDto], err error) {
				assert.NoError(t, err)
				assert.NotNil(t, deployments)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(repository.DeploymentRepositoryMock)
			projectRepo := new(projectRepository.ProjectRepositoryMock)
			tt.setupMocks(repo, projectRepo)
			s := NewDeploymentService(repo, projectRepo, nil)
			deployments, err := s.GetHistory(context.TODO(), tt.userID, tt.projectID)
			tt.expectedFunc(t, deployments, err)
		})
	}
}

func TestDeploymentService_GetDeploymentByID(t *testing.T) {
	tests := []struct {
		name         string
		userID       uint
		id           uint
		setupMocks   func(repo *repository.DeploymentRepositoryMock, projectRepo *projectRepository.ProjectRepositoryMock)
		expectedFunc func(t *testing.T, deployment *response.DeploymentDto, err error)
	}{
		{
			name:   "Get deployment successfully",
			userID: 1,
			id:     1,
			setupMocks: func(repo *repository.DeploymentRepositoryMock, projectRepo *projectRepository.ProjectRepositoryMock) {
				repo.On("FindByID", mock.Anything, uint(1)).Return(&models.Deployment{Model: gorm.Model{ID: 1}, ProjectID: 1}, nil)
				projectRepo.On("FindByID", mock.Anything, uint(1)).Return(&models.Project{Model: gorm.Model{ID: 1}, UserID: 1}, nil)
			},
			expectedFunc: func(t *testing.T, deployment *response.DeploymentDto, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, deployment)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(repository.DeploymentRepositoryMock)
			projectRepo := new(projectRepository.ProjectRepositoryMock)
			tt.setupMocks(repo, projectRepo)
			s := NewDeploymentService(repo, projectRepo, nil)
			deployment, err := s.GetDeploymentByID(context.TODO(), tt.userID, tt.id)
			tt.expectedFunc(t, deployment, err)
		})
	}
}
