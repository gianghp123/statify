package workers

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/gianghp/statify/internal/core/enums"
	"github.com/gianghp/statify/internal/database/models"
	"github.com/gianghp/statify/internal/modules/deployment/repository"
	"github.com/gianghp/statify/internal/modules/file/processor"
	jobQueueRepo "github.com/gianghp/statify/internal/modules/job-queue/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestDeploymentWorker_ProcessBuild(t *testing.T) {
	tests := []struct {
		name       string
		setupMocks func(jobRepo *jobQueueRepo.JobQueueRepositoryMock, depRepo *repository.DeploymentRepositoryMock, fileProc *processor.FileProcessorMock)
	}{
		{
			name: "Process deployment build successfully",
			setupMocks: func(jobRepo *jobQueueRepo.JobQueueRepositoryMock, depRepo *repository.DeploymentRepositoryMock, fileProc *processor.FileProcessorMock) {
				fileProc.On("ProcessDeploymentFiles", mock.Anything, mock.Anything).Return(nil)
				depRepo.On("MarkReady", mock.Anything, uint(1)).Return(nil)
			},
		},
		{
			name: "Process deployment build failure",
			setupMocks: func(jobRepo *jobQueueRepo.JobQueueRepositoryMock, depRepo *repository.DeploymentRepositoryMock, fileProc *processor.FileProcessorMock) {
				fileProc.On("ProcessDeploymentFiles", mock.Anything, mock.Anything).Return(errors.New("build error"))
				depRepo.On("MarkFailed", mock.Anything, uint(1), "build error").Return(nil)
			},
		},
		{
			name: "Process deployment build panic recovery",
			setupMocks: func(jobRepo *jobQueueRepo.JobQueueRepositoryMock, depRepo *repository.DeploymentRepositoryMock, fileProc *processor.FileProcessorMock) {
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
			fileProc := new(processor.FileProcessorMock)
			tt.setupMocks(jobRepo, depRepo, fileProc)

			worker := NewDeploymentWorker(jobRepo, depRepo, fileProc)
			job := &models.JobQueue{
				Deployment: models.Deployment{Model: gorm.Model{ID: 1}},
			}
			worker.processBuild(context.Background(), job, 1*time.Millisecond, 1)

			jobRepo.AssertExpectations(t)
			depRepo.AssertExpectations(t)
			fileProc.AssertExpectations(t)
		})
	}
}

func TestDeploymentWorker_ProcessDelete(t *testing.T) {
	tests := []struct {
		name       string
		setupMocks func(jobRepo *jobQueueRepo.JobQueueRepositoryMock, depRepo *repository.DeploymentRepositoryMock, fileProc *processor.FileProcessorMock)
		checkJob   func(t *testing.T, job *models.JobQueue)
	}{
		{
			name: "Delete deployment folder successfully",
			setupMocks: func(jobRepo *jobQueueRepo.JobQueueRepositoryMock, depRepo *repository.DeploymentRepositoryMock, fileProc *processor.FileProcessorMock) {
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
			setupMocks: func(jobRepo *jobQueueRepo.JobQueueRepositoryMock, depRepo *repository.DeploymentRepositoryMock, fileProc *processor.FileProcessorMock) {
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
			fileProc := new(processor.FileProcessorMock)
			tt.setupMocks(jobRepo, depRepo, fileProc)

			worker := NewDeploymentWorker(jobRepo, depRepo, fileProc)
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
