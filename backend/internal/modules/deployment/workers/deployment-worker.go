package workers

import (
	"context"
	"fmt"
	"time"

	"log/slog"

	"github.com/gianghp/statify/internal/core/enums"
	"github.com/gianghp/statify/internal/database/models"
	"github.com/gianghp/statify/internal/modules/deployment/repository"
	"github.com/gianghp/statify/internal/modules/file/processor"
	jobQueueRepo "github.com/gianghp/statify/internal/modules/job-queue/repository"
	"github.com/gianghp/statify/internal/utils"
)

type DeploymentWorker struct {
	jobQueueRepository   jobQueueRepo.IJobQueueRepository
	deploymentRepository repository.IDeploymentRepository
	fileProcessor        processor.IFileProcessor
}

func NewDeploymentWorker(jobQueueRepository jobQueueRepo.IJobQueueRepository, deploymentRepository repository.IDeploymentRepository, fileProcessor processor.IFileProcessor) *DeploymentWorker {
	return &DeploymentWorker{
		jobQueueRepository:   jobQueueRepository,
		deploymentRepository: deploymentRepository,
		fileProcessor:        fileProcessor,
	}
}

func (w *DeploymentWorker) Run(ctx context.Context, maxProcessGoroutines int, maxDeleteGoroutines int, maxProjectDeleteGoroutines int, sleepTime time.Duration, maxRetry int) {
	for range maxDeleteGoroutines {
		go w.workerLoop(ctx, enums.JobQueueTypeDeploymentDelete, sleepTime, maxRetry)
	}

	for range maxProcessGoroutines {
		go w.workerLoop(ctx, enums.JobQueueTypeDeploymentProcess, sleepTime, maxRetry)
	}

	for range maxProjectDeleteGoroutines {
		go w.workerLoop(ctx, enums.JobQueueTypeProjectDelete, sleepTime, maxRetry)
	}

	<-ctx.Done()
}

func (w *DeploymentWorker) workerLoop(ctx context.Context, jobType enums.JobQueueType, sleep time.Duration, maxRetry int) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		job, err := w.jobQueueRepository.ClaimNextQueueByType(ctx, jobType)
		if err != nil {
			time.Sleep(sleep)
			continue
		}

		if job == nil {
			time.Sleep(sleep)
			continue
		}

		switch job.Type {
		case enums.JobQueueTypeDeploymentProcess:
			w.processBuild(ctx, job, sleep, maxRetry)

		case enums.JobQueueTypeDeploymentDelete:
			w.processDelete(ctx, job, sleep, maxRetry)

		case enums.JobQueueTypeProjectDelete:
			w.processProjectDelete(ctx, job, sleep, maxRetry)

		default:
			slog.Warn("unknown job type", "type", job.Type)
		}
	}
}

func (w *DeploymentWorker) processBuild(ctx context.Context, job *models.JobQueue, sleepTime time.Duration, maxRetry int) {
	deployment := job.Deployment
	defer func() {
		if r := recover(); r != nil {
			w.deploymentRepository.MarkFailed(ctx, deployment.ID, "Error processing deployment")
		}
	}()

	err := w.fileProcessor.ProcessDeploymentFiles(ctx, &deployment)

	var updateErr error
	if err != nil {
		updateErr = utils.Retry(maxRetry, sleepTime, func() error {
			return w.deploymentRepository.MarkFailed(ctx, deployment.ID, err.Error())
		})
	} else {
		updateErr = utils.Retry(maxRetry, sleepTime, func() error {
			return w.deploymentRepository.MarkReady(ctx, deployment.ID)
		})
	}

	if updateErr != nil {
		slog.Error("failed to update job", "error", updateErr)
	}
}

func (w *DeploymentWorker) processDelete(ctx context.Context, job *models.JobQueue, sleepTime time.Duration, maxRetry int) {
	defer func() {
		if r := recover(); r != nil {
			job.Status = enums.JobQueueStatusFailed
			job.Error = "panic during delete"
			_ = w.jobQueueRepository.Update(ctx, job)
		}
	}()

	deployment := job.Deployment

	var deleteErr error

	for range maxRetry {
		select {
		case <-ctx.Done():
			return
		default:
		}

		deleteErr = w.fileProcessor.DeleteMinioFolder(ctx, deployment.OutputPrefix)

		if deleteErr == nil {
			break
		}

		job.RetryCount++
		time.Sleep(sleepTime)
	}

	if deleteErr != nil {
		job.Status = enums.JobQueueStatusFailed
		job.Error = deleteErr.Error()
	} else {
		job.Status = enums.JobQueueStatusSuccess
	}

	updateErr := utils.Retry(maxRetry, sleepTime, func() error {
		return w.jobQueueRepository.Update(ctx, job)
	})

	if updateErr != nil {
		slog.Error("failed to update job", "error", updateErr)
	}
}
func (w *DeploymentWorker) processProjectDelete(ctx context.Context, job *models.JobQueue, sleepTime time.Duration, maxRetry int) {
	defer func() {
		if r := recover(); r != nil {
			job.Status = enums.JobQueueStatusFailed
			job.Error = "panic during project delete"
			_ = w.jobQueueRepository.Update(ctx, job)
		}
	}()

	// Payload is projectID
	prefix := fmt.Sprintf("deployments/%s", job.Payload)

	var deleteErr error

	for range maxRetry {
		select {
		case <-ctx.Done():
			return
		default:
		}

		deleteErr = w.fileProcessor.DeleteMinioFolder(ctx, prefix)

		if deleteErr == nil {
			break
		}

		job.RetryCount++
		time.Sleep(sleepTime)
	}

	if deleteErr != nil {
		job.Status = enums.JobQueueStatusFailed
		job.Error = deleteErr.Error()
	} else {
		job.Status = enums.JobQueueStatusSuccess
	}

	updateErr := utils.Retry(maxRetry, sleepTime, func() error {
		return w.jobQueueRepository.Update(ctx, job)
	})

	if updateErr != nil {
		slog.Error("failed to update job", "error", updateErr)
	}
}
