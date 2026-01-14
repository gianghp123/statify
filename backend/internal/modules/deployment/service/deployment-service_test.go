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
	"github.com/gianghp/statify/internal/storage/minio"
	minioGo "github.com/minio/minio-go/v7"
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
		setupMocks   func(repo *repository.DeploymentRepositoryMock, projectRepo *projectRepository.ProjectRepositoryMock, minioClient *minio.Mock)
		expectedFunc func(t *testing.T, deployment *response.DeploymentDto, err error)
	}{
		{
			name:      "Create deployment successfully",
			userID:    1,
			projectID: 1,
			request: &request.CreateDeploymentRequest{
				File: fileHeader,
			},
			setupMocks: func(repo *repository.DeploymentRepositoryMock, projectRepo *projectRepository.ProjectRepositoryMock, minioClient *minio.Mock) {
				project := &models.Project{
					Model:               gorm.Model{ID: 1},
					UserID:              1,
					Subdomain:           "test-site",
					CurrentDeploymentID: 1,
				}
				projectRepo.On("FindByID", mock.Anything, uint(1)).Return(project, nil)

				// 1. Expect initial DB creation (Status: Uploaded/Processing)
				repo.On("Create", mock.Anything, mock.MatchedBy(func(d *models.Deployment) bool {
					return d.ProjectID == 1 && d.Status == enums.DeploymentStatusProcessing
				})).Return(nil).Run(func(args mock.Arguments) {
					// Simulate GORM setting the ID after creation
					d := args.Get(1).(*models.Deployment)
					d.ID = 10
				})

				// 2. Expect PutObject for each file
				minioClient.On("PutObject", mock.Anything, "static-sites",
					mock.MatchedBy(func(s string) bool { return strings.Contains(s, "index.html") }),
					mock.Anything, mock.Anything, mock.Anything,
				).Return(minioGo.UploadInfo{}, nil)

				minioClient.On("PutObject", mock.Anything, "static-sites",
					mock.MatchedBy(func(s string) bool { return strings.Contains(s, "css/style.css") }),
					mock.Anything, mock.Anything, mock.Anything,
				).Return(minioGo.UploadInfo{}, nil)

				// 3. Expect Deployment status update to Ready
				repo.On("Update", mock.Anything, mock.MatchedBy(func(d *models.Deployment) bool {
					return d.ID == 10 && d.Status == enums.DeploymentStatusReady
				})).Return(nil)

				// 5. Expect Update for OutputPrefix setting (inside transaction)
				repo.On("Update", mock.Anything, mock.MatchedBy(func(d *models.Deployment) bool {
					return d.ID == 10 && strings.HasPrefix(d.OutputPrefix, "deployments/1/10-")
				})).Return(nil)

				// Mock Transaction and WithTx
				repo.On("Transaction", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
					fn := args.Get(1).(func(*gorm.DB) error)
					_ = fn(nil)
				})
				repo.On("WithTx", mock.Anything).Return(repo)
			},
			expectedFunc: func(t *testing.T, deployment *response.DeploymentDto, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, deployment)
				assert.Equal(t, string(enums.DeploymentStatusReady), deployment.Status)
			},
		},
		{
			name:      "Create deployment successfully when project is empty",
			userID:    1,
			projectID: 1,
			request: &request.CreateDeploymentRequest{
				File: fileHeader,
			},
			setupMocks: func(repo *repository.DeploymentRepositoryMock, projectRepo *projectRepository.ProjectRepositoryMock, minioClient *minio.Mock) {
				project := &models.Project{
					Model:     gorm.Model{ID: 1},
					UserID:    1,
					Subdomain: "test-site",
				}
				projectRepo.On("FindByID", mock.Anything, uint(1)).Return(project, nil)

				// 1. Expect initial DB creation (Status: Uploaded/Processing)
				repo.On("Create", mock.Anything, mock.MatchedBy(func(d *models.Deployment) bool {
					return d.ProjectID == 1 && d.Status == enums.DeploymentStatusProcessing
				})).Return(nil).Run(func(args mock.Arguments) {
					// Simulate GORM setting the ID after creation
					d := args.Get(1).(*models.Deployment)
					d.ID = 10
				})

				// 2. Expect PutObject for each file
				minioClient.On("PutObject", mock.Anything, "static-sites",
					mock.MatchedBy(func(s string) bool { return strings.Contains(s, "index.html") }),
					mock.Anything, mock.Anything, mock.Anything,
				).Return(minioGo.UploadInfo{}, nil)

				minioClient.On("PutObject", mock.Anything, "static-sites",
					mock.MatchedBy(func(s string) bool { return strings.Contains(s, "css/style.css") }),
					mock.Anything, mock.Anything, mock.Anything,
				).Return(minioGo.UploadInfo{}, nil)

				// 3. Expect Deployment status update to Ready
				repo.On("Update", mock.Anything, mock.MatchedBy(func(d *models.Deployment) bool {
					return d.ID == 10 && d.Status == enums.DeploymentStatusReady
				})).Return(nil)

				// 5. Expect Update for OutputPrefix setting (inside transaction)
				repo.On("Update", mock.Anything, mock.MatchedBy(func(d *models.Deployment) bool {
					return d.ID == 10 && strings.HasPrefix(d.OutputPrefix, "deployments/1/10-")
				})).Return(nil)

				// Mock Transaction and WithTx
				repo.On("Transaction", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
					fn := args.Get(1).(func(*gorm.DB) error)
					_ = fn(nil)
				})
				repo.On("WithTx", mock.Anything).Return(repo)
			},
			expectedFunc: func(t *testing.T, deployment *response.DeploymentDto, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, deployment)
				assert.Equal(t, string(enums.DeploymentStatusReady), deployment.Status)
			},
		},
		{
			name:      "Create deployment successfully when old deployment status is not ready",
			userID:    1,
			projectID: 1,
			request: &request.CreateDeploymentRequest{
				File: fileHeader,
			},
			setupMocks: func(repo *repository.DeploymentRepositoryMock, projectRepo *projectRepository.ProjectRepositoryMock, minioClient *minio.Mock) {
				project := &models.Project{
					Model:               gorm.Model{ID: 1},
					UserID:              1,
					Subdomain:           "test-site",
					CurrentDeploymentID: 1,
				}
				projectRepo.On("FindByID", mock.Anything, uint(1)).Return(project, nil)

				// 1. Expect initial DB creation (Status: Uploaded/Processing)
				repo.On("Create", mock.Anything, mock.MatchedBy(func(d *models.Deployment) bool {
					return d.ProjectID == 1 && d.Status == enums.DeploymentStatusProcessing
				})).Return(nil).Run(func(args mock.Arguments) {
					// Simulate GORM setting the ID after creation
					d := args.Get(1).(*models.Deployment)
					d.ID = 10
				})

				// 2. Expect PutObject for each file
				minioClient.On("PutObject", mock.Anything, "static-sites",
					mock.MatchedBy(func(s string) bool { return strings.Contains(s, "index.html") }),
					mock.Anything, mock.Anything, mock.Anything,
				).Return(minioGo.UploadInfo{}, nil)

				minioClient.On("PutObject", mock.Anything, "static-sites",
					mock.MatchedBy(func(s string) bool { return strings.Contains(s, "css/style.css") }),
					mock.Anything, mock.Anything, mock.Anything,
				).Return(minioGo.UploadInfo{}, nil)

				// 3. Expect Deployment status update to Ready
				repo.On("Update", mock.Anything, mock.MatchedBy(func(d *models.Deployment) bool {
					return d.ID == 10 && d.Status == enums.DeploymentStatusReady
				})).Return(nil)

				// 5. Expect Update for OutputPrefix setting (inside transaction)
				repo.On("Update", mock.Anything, mock.MatchedBy(func(d *models.Deployment) bool {
					return d.ID == 10 && strings.HasPrefix(d.OutputPrefix, "deployments/1/10-")
				})).Return(nil)

				// Mock Transaction and WithTx
				repo.On("Transaction", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
					fn := args.Get(1).(func(*gorm.DB) error)
					_ = fn(nil)
				})
				repo.On("WithTx", mock.Anything).Return(repo)
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
			setupMocks: func(repo *repository.DeploymentRepositoryMock, projectRepo *projectRepository.ProjectRepositoryMock, minioClient *minio.Mock) {
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
			minioClient := new(minio.Mock)

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
				repo.On("FindAllByProjectID", mock.Anything, uint(1), 1, 10).Return(&coreRepo.PaginatedEntities[models.Deployment]{
					Entities:   []models.Deployment{{ProjectID: 1}},
					Pagination: coreRepo.Pagination{TotalCount: 1, Page: 1, Limit: 10},
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
			deployments, err := s.GetHistory(context.TODO(), tt.userID, tt.projectID, 1, 10)
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

func TestDeploymentService_GetCurrentDeploymentFilesByProjectSubdomain(t *testing.T) {
	tests := []struct {
		name         string
		subdomain    string
		fileName     string
		clientEtag   string
		setupMocks   func(repo *repository.DeploymentRepositoryMock, projectRepo *projectRepository.ProjectRepositoryMock, minioClient *minio.Mock)
		expectedFunc func(t *testing.T, fileDTO *response.FileDownloadDto, err error)
	}{
		{
			name:       "Get files successfully - 200 OK",
			subdomain:  "testing",
			fileName:   "index.html",
			clientEtag: "",
			setupMocks: func(repo *repository.DeploymentRepositoryMock, projectRepo *projectRepository.ProjectRepositoryMock, minioClient *minio.Mock) {
				project := &models.Project{Model: gorm.Model{ID: 10}, CurrentDeploymentID: 50}
				deployment := &models.Deployment{Model: gorm.Model{ID: 50}, ProjectID: 10, Status: enums.DeploymentStatusLive}

				expectedPath := "deployments/10/50/index.html"
				bucket := "static-sites"

				projectRepo.On("FindBySubdomain", mock.Anything, "testing").Return(project, nil)
				repo.On("FindByID", mock.Anything, uint(50)).Return(deployment, nil)

				minioClient.On("StatObject", mock.Anything, bucket, expectedPath, mock.Anything).
					Return(minioGo.ObjectInfo{Size: 100, ContentType: "text/html", ETag: "v1"}, nil)

				minioClient.On("GetObject", mock.Anything, bucket, expectedPath, mock.Anything).
					Return(&minioGo.Object{}, nil)
			},
			expectedFunc: func(t *testing.T, fileDTO *response.FileDownloadDto, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, fileDTO)
				assert.False(t, fileDTO.NotModified)
				assert.Equal(t, "v1", fileDTO.Headers["ETag"])
			},
		},
		{
			name:       "File not modified - 304 Not Modified",
			subdomain:  "testing",
			fileName:   "style.css",
			clientEtag: "v1-cache",
			setupMocks: func(repo *repository.DeploymentRepositoryMock, projectRepo *projectRepository.ProjectRepositoryMock, minioClient *minio.Mock) {
				project := &models.Project{Model: gorm.Model{ID: 10}, CurrentDeploymentID: 50}
				deployment := &models.Deployment{Model: gorm.Model{ID: 50}, ProjectID: 10, Status: enums.DeploymentStatusLive}

				expectedPath := "deployments/10/50/style.css"

				projectRepo.On("FindBySubdomain", mock.Anything, "testing").Return(project, nil)
				repo.On("FindByID", mock.Anything, uint(50)).Return(deployment, nil)

				minioClient.On("StatObject", mock.Anything, "static-sites", expectedPath, mock.Anything).
					Return(minioGo.ObjectInfo{ETag: "v1-cache"}, nil)
			},
			expectedFunc: func(t *testing.T, fileDTO *response.FileDownloadDto, err error) {
				assert.NoError(t, err)
				assert.True(t, fileDTO.NotModified)
				assert.Nil(t, fileDTO.Stream) // Should not call GetObject
			},
		},
		{
			name:      "Project not found - 404",
			subdomain: "unknown",
			fileName:  "index.html",
			setupMocks: func(repo *repository.DeploymentRepositoryMock, projectRepo *projectRepository.ProjectRepositoryMock, minioClient *minio.Mock) {
				projectRepo.On("FindBySubdomain", mock.Anything, "unknown").Return((*models.Project)(nil), nil)
			},
			expectedFunc: func(t *testing.T, fileDTO *response.FileDownloadDto, err error) {
				assert.Error(t, err)
				assert.Nil(t, fileDTO)
			},
		},
		{
			name:      "Status is not live",
			subdomain: "testing",
			fileName:  "index.html",
			setupMocks: func(repo *repository.DeploymentRepositoryMock, projectRepo *projectRepository.ProjectRepositoryMock, minioClient *minio.Mock) {
				project := &models.Project{Model: gorm.Model{ID: 10}, CurrentDeploymentID: 50}
				deployment := &models.Deployment{Model: gorm.Model{ID: 50}, ProjectID: 10, Status: enums.DeploymentStatusProcessing}

				projectRepo.On("FindBySubdomain", mock.Anything, "testing").Return(project, nil)
				repo.On("FindByID", mock.Anything, uint(50)).Return(deployment, nil)
			},
			expectedFunc: func(t *testing.T, fileDTO *response.FileDownloadDto, err error) {
				assert.Error(t, err)
				assert.Nil(t, fileDTO)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize Mocks
			repo := new(repository.DeploymentRepositoryMock)
			projectRepo := new(projectRepository.ProjectRepositoryMock)
			minioClient := new(minio.Mock)

			// Setup expectations
			tt.setupMocks(repo, projectRepo, minioClient)

			// Initialize Service
			s := NewDeploymentService(repo, projectRepo, minioClient)

			// Execute
			fileDTO, err := s.GetCurrentDeploymentFilesByProjectSubdomain(context.TODO(), tt.subdomain, tt.fileName, tt.clientEtag)

			// Assert
			tt.expectedFunc(t, fileDTO, err)

			repo.AssertExpectations(t)
			projectRepo.AssertExpectations(t)
			minioClient.AssertExpectations(t)
		})
	}
}

func TestDeploymentService_TurnDeploymentLive(t *testing.T) {
	tests := []struct {
		name         string
		userID       uint
		deploymentID uint
		setupMocks   func(repo *repository.DeploymentRepositoryMock, projectRepo *projectRepository.ProjectRepositoryMock)
		expectedFunc func(t *testing.T, err error)
	}{
		{
			name:         "Turn deployment live successfully",
			userID:       1,
			deploymentID: 1,
			setupMocks: func(repo *repository.DeploymentRepositoryMock, projectRepo *projectRepository.ProjectRepositoryMock) {
				repo.On("FindByID", mock.Anything, uint(1)).Return(&models.Deployment{Model: gorm.Model{ID: 1}, ProjectID: 1, Status: enums.DeploymentStatusReady}, nil)
				projectRepo.On("FindByID", mock.Anything, uint(1)).Return(&models.Project{Model: gorm.Model{ID: 1}, UserID: 1}, nil)
				repo.On("Update", mock.Anything, &models.Deployment{Model: gorm.Model{ID: 1}, ProjectID: 1, Status: enums.DeploymentStatusLive}).Return(nil)
				projectRepo.On("Update", mock.Anything, &models.Project{Model: gorm.Model{ID: 1}, UserID: 1, CurrentDeploymentID: 1}).Return(nil)
			},
			expectedFunc: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:         "Status is already live",
			userID:       1,
			deploymentID: 1,
			setupMocks: func(repo *repository.DeploymentRepositoryMock, projectRepo *projectRepository.ProjectRepositoryMock) {
				repo.On("FindByID", mock.Anything, uint(1)).Return(&models.Deployment{Model: gorm.Model{ID: 1}, ProjectID: 1, Status: enums.DeploymentStatusLive}, nil)
				projectRepo.On("FindByID", mock.Anything, uint(1)).Return(&models.Project{Model: gorm.Model{ID: 1}, UserID: 1}, nil)
			},
			expectedFunc: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name:         "Project already has a live deployment",
			userID:       1,
			deploymentID: 2,
			setupMocks: func(repo *repository.DeploymentRepositoryMock, projectRepo *projectRepository.ProjectRepositoryMock) {
				repo.On("FindByID", mock.Anything, uint(2)).Return(&models.Deployment{Model: gorm.Model{ID: 2}, ProjectID: 1, Status: enums.DeploymentStatusReady}, nil)
				projectRepo.On("FindByID", mock.Anything, uint(1)).Return(&models.Project{Model: gorm.Model{ID: 1}, UserID: 1, CurrentDeploymentID: 1}, nil)
			},
			expectedFunc: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name:         "Forbidden User",
			userID:       2,
			deploymentID: 1,
			setupMocks: func(repo *repository.DeploymentRepositoryMock, projectRepo *projectRepository.ProjectRepositoryMock) {
				repo.On("FindByID", mock.Anything, uint(1)).Return(&models.Deployment{Model: gorm.Model{ID: 1}, ProjectID: 1, Status: enums.DeploymentStatusReady}, nil)
				projectRepo.On("FindByID", mock.Anything, uint(1)).Return(&models.Project{Model: gorm.Model{ID: 1}, UserID: 1}, nil)
			},
			expectedFunc: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize Mocks
			repo := new(repository.DeploymentRepositoryMock)
			projectRepo := new(projectRepository.ProjectRepositoryMock)
			tt.setupMocks(repo, projectRepo)
			s := NewDeploymentService(repo, projectRepo, nil)
			err := s.TurnDeploymentLive(context.TODO(), tt.deploymentID, tt.userID)
			tt.expectedFunc(t, err)
		})
	}
}

func TestDeploymentService_TurnDeploymentOffline(t *testing.T) {
	tests := []struct {
		name         string
		userID       uint
		deploymentID uint
		setupMocks   func(repo *repository.DeploymentRepositoryMock, projectRepo *projectRepository.ProjectRepositoryMock)
		expectedFunc func(t *testing.T, err error)
	}{
		{
			name:         "Turn deployment offline successfully",
			userID:       1,
			deploymentID: 1,
			setupMocks: func(repo *repository.DeploymentRepositoryMock, projectRepo *projectRepository.ProjectRepositoryMock) {
				repo.On("FindByID", mock.Anything, uint(1)).Return(&models.Deployment{Model: gorm.Model{ID: 1}, ProjectID: 1, Status: enums.DeploymentStatusLive}, nil)
				projectRepo.On("FindByID", mock.Anything, uint(1)).Return(&models.Project{Model: gorm.Model{ID: 1}, UserID: 1, CurrentDeploymentID: 1}, nil)
				repo.On("Update", mock.Anything, &models.Deployment{Model: gorm.Model{ID: 1}, ProjectID: 1, Status: enums.DeploymentStatusReady}).Return(nil)
				projectRepo.On("Update", mock.Anything, &models.Project{Model: gorm.Model{ID: 1}, UserID: 1, CurrentDeploymentID: 0}).Return(nil)
			},
			expectedFunc: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:         "Deployment is not the current deployment",
			userID:       1,
			deploymentID: 1,
			setupMocks: func(repo *repository.DeploymentRepositoryMock, projectRepo *projectRepository.ProjectRepositoryMock) {
				repo.On("FindByID", mock.Anything, uint(1)).Return(&models.Deployment{Model: gorm.Model{ID: 1}, ProjectID: 1, Status: enums.DeploymentStatusLive}, nil)
				projectRepo.On("FindByID", mock.Anything, uint(1)).Return(&models.Project{Model: gorm.Model{ID: 1}, UserID: 1, CurrentDeploymentID: 2}, nil)
			},
			expectedFunc: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name:         "Forbidden User",
			userID:       2,
			deploymentID: 1,
			setupMocks: func(repo *repository.DeploymentRepositoryMock, projectRepo *projectRepository.ProjectRepositoryMock) {
				repo.On("FindByID", mock.Anything, uint(1)).Return(&models.Deployment{Model: gorm.Model{ID: 1}, ProjectID: 1, Status: enums.DeploymentStatusLive}, nil)
				projectRepo.On("FindByID", mock.Anything, uint(1)).Return(&models.Project{Model: gorm.Model{ID: 1}, UserID: 1}, nil)
			},
			expectedFunc: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize Mocks
			repo := new(repository.DeploymentRepositoryMock)
			projectRepo := new(projectRepository.ProjectRepositoryMock)
			tt.setupMocks(repo, projectRepo)
			s := NewDeploymentService(repo, projectRepo, nil)
			err := s.TurnDeploymentOffline(context.TODO(), tt.deploymentID, tt.userID)
			tt.expectedFunc(t, err)
		})
	}
}
