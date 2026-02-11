package service

import (
	"archive/zip"
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"testing"

	"github.com/gianghp/statify/internal/core/enums"
	coreRepo "github.com/gianghp/statify/internal/core/repository"
	"github.com/gianghp/statify/internal/database/models"
	"github.com/gianghp/statify/internal/modules/auth/policy"
	"github.com/gianghp/statify/internal/modules/deployment/dtos/response"
	"github.com/gianghp/statify/internal/modules/deployment/repository"
	projectRepository "github.com/gianghp/statify/internal/modules/project/repository"
	"github.com/gianghp/statify/internal/storage/minio"
	minioGo "github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestDeploymentService_GetGlobalDeploymentHistory(t *testing.T) {
	tests := []struct {
		name         string
		userID       uint
		page         int
		limit        int
		SetupMocks   func(t *testing.T, repo *repository.DeploymentRepositoryMock, policyMock *policy.AccessPolicyMock)
		expectedFunc func(t *testing.T, deployments coreRepo.PaginatedEntities[*response.DeploymentDto], err error)
	}{
		{
			name:   "Get global deployment history successfully",
			userID: 1,
			page:   1,
			limit:  10,
			SetupMocks: func(t *testing.T, repo *repository.DeploymentRepositoryMock, policyMock *policy.AccessPolicyMock) {
				repo.On("FindAllByUserID", mock.Anything, uint(1), 1, 10).Return(coreRepo.PaginatedEntities[*models.Deployment]{
					Entities: []*models.Deployment{{
						Model: gorm.Model{
							ID: 1,
						},
						Status:       enums.DeploymentStatusReady,
						OutputPrefix: "deployments/1/1-",
						ProjectID:    1,
					}},
					Pagination: coreRepo.Pagination{
						TotalCount: 0,
						Page:       1,
						Limit:      10,
					},
				}, nil)
			},
			expectedFunc: func(t *testing.T, deployments coreRepo.PaginatedEntities[*response.DeploymentDto], err error) {
				assert.NoError(t, err)
				assert.NotNil(t, deployments)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := new(repository.DeploymentRepositoryMock)
			policyMock := new(policy.AccessPolicyMock)
			test.SetupMocks(t, repo, policyMock)
			deploymentService := NewDeploymentService(repo, nil, nil, nil, nil, policyMock)
			deployments, err := deploymentService.GetGlobalDeploymentHistory(context.Background(), test.userID, test.page, test.limit)
			test.expectedFunc(t, deployments, err)
			repo.AssertExpectations(t)

		})
	}
}

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

func TestDeploymentService_GetHistory(t *testing.T) {
	tests := []struct {
		name         string
		userID       uint
		projectID    uint
		setupMocks   func(repo *repository.DeploymentRepositoryMock, projectRepo *projectRepository.ProjectRepositoryMock, policyMock *policy.AccessPolicyMock)
		expectedFunc func(t *testing.T, deployments coreRepo.PaginatedEntities[*response.DeploymentDto], err error)
	}{
		{
			name:      "Get history successfully",
			userID:    1,
			projectID: 1,
			setupMocks: func(repo *repository.DeploymentRepositoryMock, projectRepo *projectRepository.ProjectRepositoryMock, policyMock *policy.AccessPolicyMock) {
				policyMock.On("CheckProjectAccess", mock.Anything, uint(1), uint(1)).Return(&models.Project{Model: gorm.Model{ID: 1}, UserID: 1}, nil)
				repo.On("FindAllByProjectID", mock.Anything, uint(1), 1, 10).Return(coreRepo.PaginatedEntities[*models.Deployment]{
					Entities:   []*models.Deployment{{ProjectID: 1}},
					Pagination: coreRepo.Pagination{TotalCount: 1, Page: 1, Limit: 10},
				}, nil)
			},
			expectedFunc: func(t *testing.T, deployments coreRepo.PaginatedEntities[*response.DeploymentDto], err error) {
				assert.NoError(t, err)
				assert.NotNil(t, deployments)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(repository.DeploymentRepositoryMock)
			projectRepo := new(projectRepository.ProjectRepositoryMock)
			policyMock := new(policy.AccessPolicyMock)
			tt.setupMocks(repo, projectRepo, policyMock)
			s := NewDeploymentService(repo, projectRepo, nil, nil, nil, policyMock)
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
		setupMocks   func(repo *repository.DeploymentRepositoryMock, projectRepo *projectRepository.ProjectRepositoryMock, policyMock *policy.AccessPolicyMock)
		expectedFunc func(t *testing.T, deployment *response.DeploymentDto, err error)
	}{
		{
			name:   "Get deployment successfully",
			userID: 1,
			id:     1,
			setupMocks: func(repo *repository.DeploymentRepositoryMock, projectRepo *projectRepository.ProjectRepositoryMock, policyMock *policy.AccessPolicyMock) {
				policyMock.On("CheckDeploymentAccess", mock.Anything, uint(1), uint(1)).Return(&models.Project{Model: gorm.Model{ID: 1}, UserID: 1}, &models.Deployment{Model: gorm.Model{ID: 1}, ProjectID: 1}, nil)
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
			policyMock := new(policy.AccessPolicyMock)
			tt.setupMocks(repo, projectRepo, policyMock)
			s := NewDeploymentService(repo, projectRepo, nil, nil, nil, policyMock)
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
		setupMocks   func(repo *repository.DeploymentRepositoryMock, projectRepo *projectRepository.ProjectRepositoryMock, minioClient *minio.Mock, policyMock *policy.AccessPolicyMock)
		expectedFunc func(t *testing.T, fileDTO *response.FileDownloadDto, err error)
	}{
		{
			name:       "Get files successfully - 200 OK",
			subdomain:  "testing",
			fileName:   "index.html",
			clientEtag: "",
			setupMocks: func(repo *repository.DeploymentRepositoryMock, projectRepo *projectRepository.ProjectRepositoryMock, minioClient *minio.Mock, policyMock *policy.AccessPolicyMock) {
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
			},
		},
		{
			name:       "File not modified - 304 Not Modified",
			subdomain:  "testing",
			fileName:   "style.css",
			clientEtag: "v1-cache",
			setupMocks: func(repo *repository.DeploymentRepositoryMock, projectRepo *projectRepository.ProjectRepositoryMock, minioClient *minio.Mock, policyMock *policy.AccessPolicyMock) {
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
			setupMocks: func(repo *repository.DeploymentRepositoryMock, projectRepo *projectRepository.ProjectRepositoryMock, minioClient *minio.Mock, policyMock *policy.AccessPolicyMock) {
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
			setupMocks: func(repo *repository.DeploymentRepositoryMock, projectRepo *projectRepository.ProjectRepositoryMock, minioClient *minio.Mock, policyMock *policy.AccessPolicyMock) {
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
		{
			name:      "Page not found with custom 404",
			subdomain: "testing",
			fileName:  "about/",
			setupMocks: func(repo *repository.DeploymentRepositoryMock, projectRepo *projectRepository.ProjectRepositoryMock, minioClient *minio.Mock, policyMock *policy.AccessPolicyMock) {
				project := &models.Project{Model: gorm.Model{ID: 10}, CurrentDeploymentID: 50}
				deployment := &models.Deployment{Model: gorm.Model{ID: 50}, ProjectID: 10, Status: enums.DeploymentStatusLive, OutputPrefix: "deployments/10/50", IsSPA: false}

				projectRepo.On("FindBySubdomain", mock.Anything, "testing").Return(project, nil)
				repo.On("FindByID", mock.Anything, uint(50)).Return(deployment, nil)

				// Initial call with trailing slash
				minioClient.On("StatObject", mock.Anything, mock.Anything, "deployments/10/50/about/", mock.Anything).
					Return(minioGo.ObjectInfo{}, minioGo.ErrorResponse{Code: "NoSuchKey"}).Once()

				minioClient.On("StatObject", mock.Anything, mock.Anything, "deployments/10/50/about/.html", mock.Anything).
					Return(minioGo.ObjectInfo{}, minioGo.ErrorResponse{Code: "NoSuchKey"}).Once()

				// Folder index fallback (TrimSuffix ensures about/ -> about/index.html)
				minioClient.On("StatObject", mock.Anything, mock.Anything, "deployments/10/50/about/index.html", mock.Anything).
					Return(minioGo.ObjectInfo{}, minioGo.ErrorResponse{Code: "NoSuchKey"}).Once()

				// 404.html fallback
				minioClient.On("StatObject", mock.Anything, mock.Anything, "deployments/10/50/404.html", mock.Anything).
					Return(minioGo.ObjectInfo{}, nil).Once()

				minioClient.On("GetObject", mock.Anything, mock.Anything, "deployments/10/50/404.html", mock.Anything).
					Return(&minioGo.Object{}, nil).Once()
			},
			expectedFunc: func(t *testing.T, fileDTO *response.FileDownloadDto, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, fileDTO)
			},
		},
		{
			name:      "SPA Fallback successfully",
			subdomain: "testing",
			fileName:  "/some-route",
			setupMocks: func(repo *repository.DeploymentRepositoryMock, projectRepo *projectRepository.ProjectRepositoryMock, minioClient *minio.Mock, policyMock *policy.AccessPolicyMock) {
				project := &models.Project{Model: gorm.Model{ID: 10}, CurrentDeploymentID: 50}
				deployment := &models.Deployment{Model: gorm.Model{ID: 50}, ProjectID: 10, Status: enums.DeploymentStatusLive, OutputPrefix: "deployments/10/50", IsSPA: true}

				projectRepo.On("FindBySubdomain", mock.Anything, "testing").Return(project, nil)
				repo.On("FindByID", mock.Anything, uint(50)).Return(deployment, nil)

				// Initial call
				minioClient.On("StatObject", mock.Anything, mock.Anything, "deployments/10/50/some-route", mock.Anything).
					Return(minioGo.ObjectInfo{}, minioGo.ErrorResponse{Code: "NoSuchKey"}).Once()

				// .html fallback
				minioClient.On("StatObject", mock.Anything, mock.Anything, "deployments/10/50/some-route.html", mock.Anything).
					Return(minioGo.ObjectInfo{}, minioGo.ErrorResponse{Code: "NoSuchKey"}).Once()

				// SPA fallback
				minioClient.On("StatObject", mock.Anything, mock.Anything, "deployments/10/50/index.html", mock.Anything).
					Return(minioGo.ObjectInfo{}, nil).Once()

				minioClient.On("GetObject", mock.Anything, mock.Anything, "deployments/10/50/index.html", mock.Anything).
					Return(&minioGo.Object{}, nil).Once()
			},
			expectedFunc: func(t *testing.T, fileDTO *response.FileDownloadDto, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, fileDTO)
				assert.Equal(t, http.StatusOK, fileDTO.StatusCode)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize Mocks
			repo := new(repository.DeploymentRepositoryMock)
			projectRepo := new(projectRepository.ProjectRepositoryMock)
			minioClient := new(minio.Mock)
			policyMock := new(policy.AccessPolicyMock)

			// Setup expectations
			tt.setupMocks(repo, projectRepo, minioClient, policyMock)

			// Initialize Service
			s := NewDeploymentService(repo, projectRepo, nil, nil, minioClient, policyMock)

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
