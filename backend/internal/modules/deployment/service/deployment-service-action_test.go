package service

import (
	"context"
	"net/url"
	"testing"

	"github.com/gianghp/statify/internal/core"
	"github.com/gianghp/statify/internal/core/enums"
	"github.com/gianghp/statify/internal/core/repository/transaction"
	"github.com/gianghp/statify/internal/database/models"
	"github.com/gianghp/statify/internal/modules/auth/policy"
	"github.com/gianghp/statify/internal/modules/deployment/repository"
	jobQueueRepository "github.com/gianghp/statify/internal/modules/job-queue/repository"
	projectRepository "github.com/gianghp/statify/internal/modules/project/repository"
	uploadSessionRepository "github.com/gianghp/statify/internal/modules/upload-session/repository"
	"github.com/gianghp/statify/internal/storage/minio"
	minioGo "github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestDeploymentService_CreatePresignedUrl(t *testing.T) {
	tests := []struct {
		name         string
		userID       uint
		projectID    uint
		setupMocks   func(uploadSessionRepo *uploadSessionRepository.UploadSessionRepositoryMock, minioClient *minio.Mock, policyMock *policy.AccessPolicyMock)
		expectedFunc func(t *testing.T, err error)
	}{
		{
			name:      "Existing unexpired session found",
			userID:    1,
			projectID: 1,
			setupMocks: func(uploadSessionRepo *uploadSessionRepository.UploadSessionRepositoryMock, minioClient *minio.Mock, policyMock *policy.AccessPolicyMock) {
				policyMock.On("CheckProjectAccess", mock.Anything, uint(1), uint(1)).Return(&models.Project{Model: gorm.Model{ID: 1}}, nil)
				uploadSessionRepo.On("FindUnexpiredByProjectID", mock.Anything, uint(1)).Return(&models.DeploymentUploadSession{Model: gorm.Model{ID: 1}}, nil)
			},
			expectedFunc: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:      "Create new session successfully",
			userID:    1,
			projectID: 1,
			setupMocks: func(uploadSessionRepo *uploadSessionRepository.UploadSessionRepositoryMock, minioClient *minio.Mock, policyMock *policy.AccessPolicyMock) {
				policyMock.On("CheckProjectAccess", mock.Anything, uint(1), uint(1)).Return(&models.Project{Model: gorm.Model{ID: 1}}, nil)
				uploadSessionRepo.On("FindUnexpiredByProjectID", mock.Anything, uint(1)).Return((*models.DeploymentUploadSession)(nil), nil)

				parsedURL, _ := url.Parse("https://storage.example.com/upload")
				minioClient.On("PresignedPostPolicy", mock.Anything, mock.Anything).Return(parsedURL, map[string]string{}, nil)

				uploadSessionRepo.On("Create", mock.Anything, mock.MatchedBy(func(s *models.DeploymentUploadSession) bool {
					return s.PresignedUrl == "https://storage.example.com/upload"
				})).Return(nil)
			},
			expectedFunc: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uploadSessionRepo := new(uploadSessionRepository.UploadSessionRepositoryMock)
			minioClient := new(minio.Mock)
			policyMock := new(policy.AccessPolicyMock)
			tt.setupMocks(uploadSessionRepo, minioClient, policyMock)

			tmMock := new(transaction.TransactionManagerMock)
			s := NewDeploymentService(nil, nil, nil, uploadSessionRepo, minioClient, policyMock, tmMock)
			_, err := s.CreatePresignedUrl(context.TODO(), tt.userID, tt.projectID)
			tt.expectedFunc(t, err)
		})
	}
}

func TestDeploymentService_ConfirmCreateDeployment(t *testing.T) {
	tests := []struct {
		name            string
		uploadSessionID uint
		setupMocks      func(repo *repository.DeploymentRepositoryMock, jqRepo *jobQueueRepository.JobQueueRepositoryMock, uploadSessionRepo *uploadSessionRepository.UploadSessionRepositoryMock, minioClient *minio.Mock, tmMock *transaction.TransactionManagerMock)
		expectedFunc    func(t *testing.T, err error)
	}{
		{
			name:            "Confirm deployment successfully",
			uploadSessionID: 1,
			setupMocks: func(repo *repository.DeploymentRepositoryMock, jqRepo *jobQueueRepository.JobQueueRepositoryMock, uploadSessionRepo *uploadSessionRepository.UploadSessionRepositoryMock, minioClient *minio.Mock, tmMock *transaction.TransactionManagerMock) {
				session := &models.DeploymentUploadSession{
					Model:        gorm.Model{ID: 1},
					ProjectID:    1,
					UploadKey:    "test/source.zip",
					OutputPrefix: "deployments/1/test",
				}
				uploadSessionRepo.On("FindByID", mock.Anything, uint(1)).Return(session, nil)
				minioClient.On("StatObject", mock.Anything, mock.Anything, "test/source.zip", mock.Anything).Return(minioGo.ObjectInfo{}, nil)

				tmMock.On("Transaction", mock.Anything, mock.Anything).Return(func(ctx context.Context, fn func(tx *gorm.DB) error) error {
					return fn(nil)
				})
				repo.On("WithTx", mock.Anything).Return(repo)
				jqRepo.On("WithTx", mock.Anything).Return(jqRepo)

				repo.On("Create", mock.Anything, mock.MatchedBy(func(d *models.Deployment) bool {
					return d.ProjectID == 1 && d.Status == enums.DeploymentStatusQueued
				})).Return(nil)
				jqRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
			},
			expectedFunc: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:            "File not uploaded to MinIO",
			uploadSessionID: 1,
			setupMocks: func(repo *repository.DeploymentRepositoryMock, jqRepo *jobQueueRepository.JobQueueRepositoryMock, uploadSessionRepo *uploadSessionRepository.UploadSessionRepositoryMock, minioClient *minio.Mock, tmMock *transaction.TransactionManagerMock) {
				session := &models.DeploymentUploadSession{
					Model:     gorm.Model{ID: 1},
					UploadKey: "test/source.zip",
				}
				uploadSessionRepo.On("FindByID", mock.Anything, uint(1)).Return(session, nil)
				minioClient.On("StatObject", mock.Anything, mock.Anything, "test/source.zip", mock.Anything).Return(minioGo.ObjectInfo{}, minioGo.ErrorResponse{Code: "NoSuchKey"})
			},
			expectedFunc: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(repository.DeploymentRepositoryMock)
			jqRepo := new(jobQueueRepository.JobQueueRepositoryMock)
			uploadSessionRepo := new(uploadSessionRepository.UploadSessionRepositoryMock)
			minioClient := new(minio.Mock)
			tmMock := new(transaction.TransactionManagerMock)
			tt.setupMocks(repo, jqRepo, uploadSessionRepo, minioClient, tmMock)

			s := NewDeploymentService(repo, nil, jqRepo, uploadSessionRepo, minioClient, nil, tmMock)
			_, err := s.ConfirmCreateDeployment(context.TODO(), tt.uploadSessionID)
			tt.expectedFunc(t, err)
		})
	}
}

func TestDeploymentService_TurnDeploymentLive(t *testing.T) {
	tests := []struct {
		name         string
		userID       uint
		deploymentID uint
		setupMocks   func(repo *repository.DeploymentRepositoryMock, projectRepo *projectRepository.ProjectRepositoryMock, policyMock *policy.AccessPolicyMock, tmMock *transaction.TransactionManagerMock)
		expectedFunc func(t *testing.T, err error)
	}{
		{
			name:         "Turn deployment live successfully",
			userID:       1,
			deploymentID: 1,
			setupMocks: func(repo *repository.DeploymentRepositoryMock, projectRepo *projectRepository.ProjectRepositoryMock, policyMock *policy.AccessPolicyMock, tmMock *transaction.TransactionManagerMock) {
				policyMock.On("CheckDeploymentAccess", mock.Anything, uint(1), uint(1)).Return(&models.Project{Model: gorm.Model{ID: 1}, UserID: 1}, &models.Deployment{Model: gorm.Model{ID: 1}, ProjectID: 1, Status: enums.DeploymentStatusReady}, nil)
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
			setupMocks: func(repo *repository.DeploymentRepositoryMock, projectRepo *projectRepository.ProjectRepositoryMock, policyMock *policy.AccessPolicyMock, tmMock *transaction.TransactionManagerMock) {
				policyMock.On("CheckDeploymentAccess", mock.Anything, uint(1), uint(1)).Return(&models.Project{Model: gorm.Model{ID: 1}, UserID: 1}, &models.Deployment{Model: gorm.Model{ID: 1}, ProjectID: 1, Status: enums.DeploymentStatusLive}, nil)
			},
			expectedFunc: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name:         "Project already has a live deployment",
			userID:       1,
			deploymentID: 2,
			setupMocks: func(repo *repository.DeploymentRepositoryMock, projectRepo *projectRepository.ProjectRepositoryMock, policyMock *policy.AccessPolicyMock, tmMock *transaction.TransactionManagerMock) {
				policyMock.On("CheckDeploymentAccess", mock.Anything, uint(1), uint(2)).Return(&models.Project{Model: gorm.Model{ID: 1}, UserID: 1, CurrentDeploymentID: 1}, &models.Deployment{Model: gorm.Model{ID: 2}, ProjectID: 1, Status: enums.DeploymentStatusReady}, nil)
			},
			expectedFunc: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name:         "Forbidden User",
			userID:       2,
			deploymentID: 1,
			setupMocks: func(repo *repository.DeploymentRepositoryMock, projectRepo *projectRepository.ProjectRepositoryMock, policyMock *policy.AccessPolicyMock, tmMock *transaction.TransactionManagerMock) {
				policyMock.On("CheckDeploymentAccess", mock.Anything, uint(2), uint(1)).Return((*models.Project)(nil), (*models.Deployment)(nil), core.ForbiddenError())
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
			policyMock := new(policy.AccessPolicyMock)
			tmMock := new(transaction.TransactionManagerMock)
			tt.setupMocks(repo, projectRepo, policyMock, tmMock)
			s := NewDeploymentService(repo, projectRepo, nil, nil, nil, policyMock, tmMock)
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
		setupMocks   func(repo *repository.DeploymentRepositoryMock, projectRepo *projectRepository.ProjectRepositoryMock, policyMock *policy.AccessPolicyMock, tmMock *transaction.TransactionManagerMock)
		expectedFunc func(t *testing.T, err error)
	}{
		{
			name:         "Turn deployment offline successfully",
			userID:       1,
			deploymentID: 1,
			setupMocks: func(repo *repository.DeploymentRepositoryMock, projectRepo *projectRepository.ProjectRepositoryMock, policyMock *policy.AccessPolicyMock, tmMock *transaction.TransactionManagerMock) {
				policyMock.On("CheckDeploymentAccess", mock.Anything, uint(1), uint(1)).Return(&models.Project{Model: gorm.Model{ID: 1}, UserID: 1, CurrentDeploymentID: 1}, &models.Deployment{Model: gorm.Model{ID: 1}, ProjectID: 1, Status: enums.DeploymentStatusLive}, nil)
				repo.On("Update", mock.Anything, mock.Anything).Return(nil)
				projectRepo.On("Update", mock.Anything, mock.Anything).Return(nil)
			},
			expectedFunc: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:         "Deployment is not the current deployment",
			userID:       1,
			deploymentID: 1,
			setupMocks: func(repo *repository.DeploymentRepositoryMock, projectRepo *projectRepository.ProjectRepositoryMock, policyMock *policy.AccessPolicyMock, tmMock *transaction.TransactionManagerMock) {
				policyMock.On("CheckDeploymentAccess", mock.Anything, uint(1), uint(1)).Return(&models.Project{Model: gorm.Model{ID: 1}, UserID: 1, CurrentDeploymentID: 2}, &models.Deployment{Model: gorm.Model{ID: 1}, ProjectID: 1, Status: enums.DeploymentStatusLive}, nil)
			},
			expectedFunc: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name:         "Forbidden User",
			userID:       2,
			deploymentID: 1,
			setupMocks: func(repo *repository.DeploymentRepositoryMock, projectRepo *projectRepository.ProjectRepositoryMock, policyMock *policy.AccessPolicyMock, tmMock *transaction.TransactionManagerMock) {
				policyMock.On("CheckDeploymentAccess", mock.Anything, uint(2), uint(1)).Return((*models.Project)(nil), (*models.Deployment)(nil), core.ForbiddenError())
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
			policyMock := new(policy.AccessPolicyMock)
			tmMock := new(transaction.TransactionManagerMock)
			tt.setupMocks(repo, projectRepo, policyMock, tmMock)
			s := NewDeploymentService(repo, projectRepo, nil, nil, nil, policyMock, tmMock)
			err := s.TurnDeploymentOffline(context.TODO(), tt.deploymentID, tt.userID)
			tt.expectedFunc(t, err)
		})
	}
}

func TestDeploymentService_ToggleIsSPAMode(t *testing.T) {
	tests := []struct {
		name         string
		userID       uint
		deploymentID uint
		setupMocks   func(repo *repository.DeploymentRepositoryMock, projectRepo *projectRepository.ProjectRepositoryMock, policyMock *policy.AccessPolicyMock, tmMock *transaction.TransactionManagerMock)
		expectedFunc func(t *testing.T, err error)
	}{
		{
			name:         "Toggle SPA mode successfully",
			userID:       1,
			deploymentID: 10,
			setupMocks: func(repo *repository.DeploymentRepositoryMock, projectRepo *projectRepository.ProjectRepositoryMock, policyMock *policy.AccessPolicyMock, tmMock *transaction.TransactionManagerMock) {
				deployment := &models.Deployment{
					Model:     gorm.Model{ID: 10},
					ProjectID: 1,
					IsSPA:     false,
				}
				policyMock.On("CheckDeploymentAccess", mock.Anything, uint(1), uint(10)).Return(&models.Project{Model: gorm.Model{ID: 1}, UserID: 1}, deployment, nil)

				repo.On("Update", mock.Anything, mock.MatchedBy(func(d *models.Deployment) bool {
					return d.ID == 10 && d.IsSPA == true
				})).Return(nil)
			},
			expectedFunc: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:         "Toggle SPA mode failure - forbidden user",
			userID:       2,
			deploymentID: 10,
			setupMocks: func(repo *repository.DeploymentRepositoryMock, projectRepo *projectRepository.ProjectRepositoryMock, policyMock *policy.AccessPolicyMock, tmMock *transaction.TransactionManagerMock) {
				policyMock.On("CheckDeploymentAccess", mock.Anything, uint(2), uint(10)).Return((*models.Project)(nil), (*models.Deployment)(nil), core.ForbiddenError())
			},
			expectedFunc: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(repository.DeploymentRepositoryMock)
			projectRepo := new(projectRepository.ProjectRepositoryMock)
			policyMock := new(policy.AccessPolicyMock)
			tmMock := new(transaction.TransactionManagerMock)
			tt.setupMocks(repo, projectRepo, policyMock, tmMock)
			s := NewDeploymentService(repo, projectRepo, nil, nil, nil, policyMock, tmMock)
			err := s.ToggleIsSPAMode(context.TODO(), tt.deploymentID, tt.userID)
			tt.expectedFunc(t, err)
		})
	}
}

func TestDeploymentService_DeleteDeployment(t *testing.T) {
	tests := []struct {
		name         string
		setupMocks   func(repo *repository.DeploymentRepositoryMock, projectRepo *projectRepository.ProjectRepositoryMock, jobQueueRepo *jobQueueRepository.JobQueueRepositoryMock, minioClient *minio.Mock, policyMock *policy.AccessPolicyMock, tmMock *transaction.TransactionManagerMock)
		expectedFunc func(t *testing.T, err error)
	}{
		{
			name: "Delete deployment successfully",
			setupMocks: func(repo *repository.DeploymentRepositoryMock, projectRepo *projectRepository.ProjectRepositoryMock, jobQueueRepo *jobQueueRepository.JobQueueRepositoryMock, minioClient *minio.Mock, policyMock *policy.AccessPolicyMock, tmMock *transaction.TransactionManagerMock) {
				projectModel := &models.Project{Model: gorm.Model{ID: 1}, UserID: 1}
				deploymentModel := &models.Deployment{Model: gorm.Model{ID: 1}, ProjectID: 1, Status: enums.DeploymentStatusReady}
				policyMock.On("CheckDeploymentAccess", mock.Anything, uint(1), uint(1)).Return(projectModel, deploymentModel, nil)

				tmMock.On("Transaction", mock.Anything, mock.Anything).Return(func(ctx context.Context, fn func(tx *gorm.DB) error) error {
					return fn(nil)
				})
				repo.On("WithTx", mock.Anything).Return(repo)
				jobQueueRepo.On("WithTx", mock.Anything).Return(jobQueueRepo)

				repo.On("Delete", mock.Anything, deploymentModel).Return(nil)
				jobQueueRepo.On("Create", mock.Anything, mock.MatchedBy(func(job *models.JobQueue) bool {
					return job.Type == enums.JobQueueTypeDeploymentDelete &&
						job.DeploymentID == deploymentModel.ID &&
						job.Status == enums.JobQueueStatusPending
				})).Return(nil)
			},
			expectedFunc: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := new(repository.DeploymentRepositoryMock)
			projectRepo := new(projectRepository.ProjectRepositoryMock)
			jobQueueRepo := new(jobQueueRepository.JobQueueRepositoryMock)
			minioClient := new(minio.Mock)
			policyMock := new(policy.AccessPolicyMock)
			tmMock := new(transaction.TransactionManagerMock)
			test.setupMocks(repo, projectRepo, jobQueueRepo, minioClient, policyMock, tmMock)
			s := NewDeploymentService(repo, projectRepo, jobQueueRepo, nil, minioClient, policyMock, tmMock)
			err := s.DeleteDeployment(context.TODO(), uint(1), uint(1))
			test.expectedFunc(t, err)
			repo.AssertExpectations(t)
			projectRepo.AssertExpectations(t)
			jobQueueRepo.AssertExpectations(t)
			minioClient.AssertExpectations(t)
			policyMock.AssertExpectations(t)
		})
	}
}
