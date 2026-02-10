package service

import (
	"context"
	"testing"

	"github.com/gianghp/statify/internal/core/enums"
	"github.com/gianghp/statify/internal/database/models"
	policy "github.com/gianghp/statify/internal/modules/auth/policy"
	"github.com/gianghp/statify/internal/modules/deployment/repository"
	jobQueueRepository "github.com/gianghp/statify/internal/modules/job-queue/repository"
	projectRepository "github.com/gianghp/statify/internal/modules/project/repository"
	"github.com/gianghp/statify/internal/storage/minio"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestDeploymentService_DeleteDeployment(t *testing.T) {
	tests := []struct {
		name         string
		setupMocks   func(repo *repository.DeploymentRepositoryMock, projectRepo *projectRepository.ProjectRepositoryMock, jobQueueRepo *jobQueueRepository.JobQueueRepositoryMock, minioClient *minio.Mock, policyMock *policy.AccessPolicyMock)
		expectedFunc func(t *testing.T, err error)
	}{
		{
			name: "Delete deployment successfully",
			setupMocks: func(repo *repository.DeploymentRepositoryMock, projectRepo *projectRepository.ProjectRepositoryMock, jobQueueRepo *jobQueueRepository.JobQueueRepositoryMock, minioClient *minio.Mock, policyMock *policy.AccessPolicyMock) {
				projectModel := &models.Project{Model: gorm.Model{ID: 1}, UserID: 1}
				deploymentModel := &models.Deployment{Model: gorm.Model{ID: 1}, ProjectID: 1, Status: enums.DeploymentStatusReady}
				policyMock.On("CheckDeploymentAccess", mock.Anything, uint(1), uint(1)).Return(projectModel, deploymentModel, nil)
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
			test.setupMocks(repo, projectRepo, jobQueueRepo, minioClient, policyMock)
			s := NewDeploymentService(repo, projectRepo, jobQueueRepo, nil, minioClient, policyMock)
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
