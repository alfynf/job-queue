package worker

import (
	"context"
	"log"
	"time"

	"github.com/alfynf/job-queue/internal/job"
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
	log.Println("masuk sini")
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
		if j.RetryCount+1 >= j.MaxRetry {
			j.Status = job.StatusFailed
			j.LastErrorMessage = &errMessage
			err := w.repo.Update(ctx, j.UUID.String(), j)
			if err != nil {
				log.Printf("error update status %s on worker: %v\n", string(job.StatusFailed), err)
			}
		} else {
			j.Status = job.StatusPending
			j.LastErrorMessage = &errMessage
			err := w.repo.Update(ctx, j.UUID.String(), j)
			if err != nil {
				log.Printf("error update status %s on worker: %v\n", string(job.StatusPending), err)
			}
		}
	} else {
		j.Status = job.StatusSuccess
		err := w.repo.Update(ctx, j.UUID.String(), j)
		if err != nil {
			log.Printf("error update status %s on worker: %v\n", string(job.StatusSuccess), err)
		}

	}
}
