package worker

import (
	"context"
	"log"
	"time"

	"github.com/alfynf/job-queue/internal/job"
	"github.com/alfynf/job-queue/pkg/logger"
	"go.uber.org/zap"
)

type JobRepository interface {
	GetJobsByStatus(ctx context.Context, status job.Status, limit int) ([]job.Job, error)
	Update(ctx context.Context, uuid string, job job.Job) error
}

type JobHandler func(ctx context.Context, payload map[string]interface{}) error

type JobWorker struct {
	repo     JobRepository
	handlers map[job.Type]JobHandler
	interval time.Duration
	batch    int
}

func New(repo JobRepository, interval time.Duration, batch int) *JobWorker {
	return &JobWorker{
		repo:     repo,
		handlers: make(map[job.Type]JobHandler),
		interval: interval,
		batch:    batch,
	}
}

func (w *JobWorker) Register(jobType job.Type, handler JobHandler) {
	log.Printf("register handler %s", string(jobType))
	w.handlers[jobType] = handler
}

func (w *JobWorker) Start(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			w.processBatch(ctx)
		case <-ctx.Done():
			return
		}
	}
}

func (w *JobWorker) processBatch(ctx context.Context) {
	// TODO: lock pending job so it wont be processed more than once
	jobs, err := w.repo.GetJobsByStatus(ctx, job.StatusPending, w.batch)
	if err != nil {
		log.Printf("worker: failded to get pending jobs: %v\n", err)
		return
	}

	for _, j := range jobs {
		go w.handleJob(ctx, j)
		// TODO: add worker pool, mutex, and mark status
	}
}

func (w *JobWorker) handleJob(ctx context.Context, j job.Job) {

	start := time.Now()
	logger.Info("Starting job",
		zap.String("job_uuid", j.UUID.String()),
		zap.String("type", string(j.Type)),
	)

	j.Status = job.StatusRunning
	err := w.repo.Update(ctx, j.UUID.String(), j)
	if err != nil {
		log.Printf("error update status %s on worker: %v\n", string(job.StatusRunning), err)
	}

	handler, ok := w.handlers[j.Type]
	if !ok {
		errMessage := "no handler for job type: " + string(j.Type)
		j.Status = job.StatusFailed
		j.LastErrorMessage = &errMessage
		err := w.repo.Update(ctx, j.UUID.String(), j)
		if err != nil {
			log.Printf("error update status %s on worker: %v\n", job.StatusFailed, err)
			return

		}
		log.Printf("error handler on worker: %v\n", err)
		return
	}

	err = handler(ctx, j.Payload)
	if err != nil {
		j.RetryCount += 1
		errMessage := err.Error()
		j.LastErrorMessage = &errMessage

		if j.RetryCount >= j.MaxRetry {
			j.Status = job.StatusFailed
			j.LastErrorMessage = &errMessage
			timeNow := time.Now()
			j.FinishedAt = &timeNow
			err := w.repo.Update(ctx, j.UUID.String(), j)
			if err != nil {
				log.Printf("error update status %s on worker: %v\n", string(job.StatusFailed), err)
			}
			logger.Error("Job permanently failed",
				zap.String("job_uuid", j.UUID.String()),
				zap.String("type", string(j.Type)),
				zap.Int("retry", j.RetryCount),
				zap.Error(err),
			)

		} else {
			j.Status = job.StatusPending
			j.LastErrorMessage = &errMessage
			logger.Warn("Job failed, will retry",
				zap.String("job_uuid", j.UUID.String()),
				zap.String("type", string(j.Type)),
				zap.Error(err),
				zap.Int("retry", j.RetryCount),
			)
			err := w.repo.Update(ctx, j.UUID.String(), j)
			if err != nil {
				log.Printf("error update status %s on worker: %v\n", string(job.StatusPending), err)
			}
		}
	} else {
		j.Status = job.StatusSuccess
		timeNow := time.Now()
		j.FinishedAt = &timeNow
		err := w.repo.Update(ctx, j.UUID.String(), j)
		if err != nil {
			log.Printf("error update status %s on worker: %v\n", string(job.StatusSuccess), err)
		}
		logger.Info("Job completed",
			zap.String("job_uuid", j.UUID.String()),
			zap.String("type", string(j.Type)),
			zap.Duration("duration", time.Since(start)),
		)

	}
}
