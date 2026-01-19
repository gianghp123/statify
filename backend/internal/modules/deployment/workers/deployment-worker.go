package workers

import (
	"context"
	"time"

	"github.com/gianghp/statify/internal/database/models"
	"github.com/gianghp/statify/internal/modules/deployment/repository"
	"github.com/gianghp/statify/internal/modules/deployment/service"
)

type DeploymentWorker struct {
	repository    repository.IDeploymentRepository
	fileProcessor *service.FileProcessor
}

func NewDeploymentWorker(repository repository.IDeploymentRepository, fileProcessor *service.FileProcessor) *DeploymentWorker {
	return &DeploymentWorker{
		repository:    repository,
		fileProcessor: fileProcessor,
	}
}

func (w *DeploymentWorker) Run(ctx context.Context, maxGoroutines int, sleepTime time.Duration) {
	for i := 0; i < maxGoroutines; i++ {
		go w.workerLoop(ctx, sleepTime)
	}
	<-ctx.Done()
}

func (w *DeploymentWorker) workerLoop(ctx context.Context, sleep time.Duration) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		job, err := w.repository.ClaimNextQueued(ctx)
		if err != nil {
			time.Sleep(sleep)
			continue
		}

		if job == nil {
			time.Sleep(sleep)
			continue
		}

		w.processOne(ctx, job)
	}
}

func (w *DeploymentWorker) processOne(ctx context.Context, job *models.Deployment) {
	defer func() {
		if r := recover(); r != nil {
			w.repository.MarkFailed(ctx, job.ID, "Error processing deployment")
		}
	}()

	err := w.fileProcessor.ProcessDeploymentFiles(ctx, job)
	if err != nil {
		w.repository.MarkFailed(ctx, job.ID, err.Error())
		return
	}

	w.repository.MarkReady(ctx, job.ID)
}
