package workers

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/gianghp/statify/internal/core/enums"
	"github.com/gianghp/statify/internal/core/repository/transaction"
	"github.com/gianghp/statify/internal/database/models"
	"github.com/gianghp/statify/internal/modules/deployment/repository"
	"github.com/gianghp/statify/internal/modules/file/processor"
	jobQueueRepo "github.com/gianghp/statify/internal/modules/job-queue/repository"
	projectRepo "github.com/gianghp/statify/internal/modules/project/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestDeploymentWorker_ProcessBuild(t *testing.T) {
	tests := []struct {
		name       string
		setupMocks func(jobRepo *jobQueueRepo.JobQueueRepositoryMock, depRepo *repository.DeploymentRepositoryMock, projRepo *projectRepo.ProjectRepositoryMock, fileProc *processor.FileProcessorMock, txManager *transaction.TransactionManagerMock)
	}{
		{
			name: "Process deployment build successfully",
			setupMocks: func(jobRepo *jobQueueRepo.JobQueueRepositoryMock, depRepo *repository.DeploymentRepositoryMock, projRepo *projectRepo.ProjectRepositoryMock, fileProc *processor.FileProcessorMock, txManager *transaction.TransactionManagerMock) {
				fileProc.On("ProcessDeploymentFiles", mock.Anything, mock.Anything).Return(nil)

				// Mock Transaction
				txManager.On("Transaction", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
					fn := args.Get(1).(func(tx *gorm.DB) error)
					_ = fn(&gorm.DB{})
				})

				depRepo.On("WithTx", mock.Anything).Return(depRepo)
				projRepo.On("WithTx", mock.Anything).Return(projRepo)

				depRepo.On("MarkReady", mock.Anything, uint(1)).Return(nil)
				projRepo.On("MarkReady", mock.Anything, uint(1), uint(1)).Return(nil)

				projRepo.On("FindByID", mock.Anything, mock.Anything).Return(&models.Project{LatestDeploymentID: 1}, nil)
				projRepo.On("Update", mock.Anything, mock.Anything).Return(nil)
			},
		},
		{
			name: "Process deployment build failure",
			setupMocks: func(jobRepo *jobQueueRepo.JobQueueRepositoryMock, depRepo *repository.DeploymentRepositoryMock, projRepo *projectRepo.ProjectRepositoryMock, fileProc *processor.FileProcessorMock, txManager *transaction.TransactionManagerMock) {
				fileProc.On("ProcessDeploymentFiles", mock.Anything, mock.Anything).Return(errors.New("build error"))

				// Mock Transaction
				txManager.On("Transaction", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
					fn := args.Get(1).(func(tx *gorm.DB) error)
					_ = fn(&gorm.DB{})
				})

				depRepo.On("WithTx", mock.Anything).Return(depRepo)
				projRepo.On("WithTx", mock.Anything).Return(projRepo)

				depRepo.On("MarkFailed", mock.Anything, uint(1), "build error").Return(nil)
				projRepo.On("MarkFailed", mock.Anything, uint(1), uint(1)).Return(nil)

				projRepo.On("FindByID", mock.Anything, mock.Anything).Return(&models.Project{LatestDeploymentID: 1}, nil)
				projRepo.On("Update", mock.Anything, mock.Anything).Return(nil)
			},
		},
		{
			name: "Process deployment build panic recovery",
			setupMocks: func(jobRepo *jobQueueRepo.JobQueueRepositoryMock, depRepo *repository.DeploymentRepositoryMock, projRepo *projectRepo.ProjectRepositoryMock, fileProc *processor.FileProcessorMock, txManager *transaction.TransactionManagerMock) {
				fileProc.On("ProcessDeploymentFiles", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
					panic("unexpected panic")
				})
				depRepo.On("MarkFailed", mock.Anything, uint(1), "Error processing deployment").Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jobRepo := new(jobQueueRepo.JobQueueRepositoryMock)
			depRepo := new(repository.DeploymentRepositoryMock)
			projRepo := new(projectRepo.ProjectRepositoryMock)
			fileProc := new(processor.FileProcessorMock)
			txManager := new(transaction.TransactionManagerMock)
			tt.setupMocks(jobRepo, depRepo, projRepo, fileProc, txManager)

			worker := NewDeploymentWorker(jobRepo, depRepo, projRepo, fileProc, txManager)
			job := &models.JobQueue{
				Deployment: models.Deployment{Model: gorm.Model{ID: 1}, ProjectID: 1},
			}
			worker.processBuild(context.Background(), job, 1*time.Millisecond, 1)

			jobRepo.AssertExpectations(t)
			depRepo.AssertExpectations(t)
			projRepo.AssertExpectations(t)
			fileProc.AssertExpectations(t)
			txManager.AssertExpectations(t)
		})
	}
}

func TestDeploymentWorker_ProcessDelete(t *testing.T) {
	tests := []struct {
		name       string
		setupMocks func(jobRepo *jobQueueRepo.JobQueueRepositoryMock, depRepo *repository.DeploymentRepositoryMock, projRepo *projectRepo.ProjectRepositoryMock, fileProc *processor.FileProcessorMock, txManager *transaction.TransactionManagerMock)
		checkJob   func(t *testing.T, job *models.JobQueue)
	}{
		{
			name: "Delete deployment folder successfully",
			setupMocks: func(jobRepo *jobQueueRepo.JobQueueRepositoryMock, depRepo *repository.DeploymentRepositoryMock, projRepo *projectRepo.ProjectRepositoryMock, fileProc *processor.FileProcessorMock, txManager *transaction.TransactionManagerMock) {
				fileProc.On("DeleteMinioFolder", mock.Anything, "prefix/").Return(nil)
				jobRepo.On("Update", mock.Anything, mock.MatchedBy(func(j *models.JobQueue) bool {
					return j.Status == enums.JobQueueStatusSuccess
				})).Return(nil)
			},
			checkJob: func(t *testing.T, job *models.JobQueue) {
				assert.Equal(t, enums.JobQueueStatusSuccess, job.Status)
			},
		},
		{
			name: "Delete deployment folder failure after retries",
			setupMocks: func(jobRepo *jobQueueRepo.JobQueueRepositoryMock, depRepo *repository.DeploymentRepositoryMock, projRepo *projectRepo.ProjectRepositoryMock, fileProc *processor.FileProcessorMock, txManager *transaction.TransactionManagerMock) {
				fileProc.On("DeleteMinioFolder", mock.Anything, "prefix/").Return(errors.New("delete error"))
				jobRepo.On("Update", mock.Anything, mock.MatchedBy(func(j *models.JobQueue) bool {
					return j.Status == enums.JobQueueStatusFailed && j.Error == "delete error"
				})).Return(nil)
			},
			checkJob: func(t *testing.T, job *models.JobQueue) {
				assert.Equal(t, enums.JobQueueStatusFailed, job.Status)
				assert.Equal(t, "delete error", job.Error)
				assert.Equal(t, 2, job.RetryCount) // maxRetry is 2 in this test
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jobRepo := new(jobQueueRepo.JobQueueRepositoryMock)
			depRepo := new(repository.DeploymentRepositoryMock)
			projRepo := new(projectRepo.ProjectRepositoryMock)
			fileProc := new(processor.FileProcessorMock)
			tt.setupMocks(jobRepo, depRepo, projRepo, fileProc, nil)

			worker := NewDeploymentWorker(jobRepo, depRepo, projRepo, fileProc, nil)
			job := &models.JobQueue{
				Deployment: models.Deployment{OutputPrefix: "prefix/"},
			}
			worker.processDelete(context.Background(), job, 1*time.Millisecond, 2)

			if tt.checkJob != nil {
				tt.checkJob(t, job)
			}

			jobRepo.AssertExpectations(t)
			depRepo.AssertExpectations(t)
			fileProc.AssertExpectations(t)
		})
	}
}
