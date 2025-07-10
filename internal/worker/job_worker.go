package worker

import (
	"context"
	"log"
	"time"

	"github.com/alfynf/job-queue/internal/job"
)

type JobRepository interface {
	GetPendingJobs(ctx context.Context, limit int) ([]job.Job, error)
	UpdateStatus(ctx context.Context, uuid string, status job.Status, errMessage *string) error
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
	jobs, err := w.repo.GetPendingJobs(ctx, w.batch)
	if err != nil {
		log.Printf("worker: failded to get pending jobs: %v\n", err)
		return
	}

	log.Println("ini jumlah pending job")
	log.Println(len(jobs))

	for _, j := range jobs {
		go w.handleJob(ctx, j)
		// TODO: add worker pool, mutex, and mark status
	}
}

func (w *JobWorker) handleJob(ctx context.Context, j job.Job) {

	err := w.repo.UpdateStatus(ctx, j.UUID, job.StatusRunning, nil)
	if err != nil {
		log.Printf("error update status %s on worker: %v\n", string(job.StatusRunning), err)
	}

	handler, ok := w.handlers[j.Type]
	if !ok {
		errMessage := "no handler for job type: " + string(j.Type)
		err := w.repo.UpdateStatus(ctx, j.UUID, job.StatusFailed, &errMessage)
		log.Printf("error update status on worker: %v\n", err)
	}

	err = handler(ctx, j.Payload)
	if err != nil {
		errMessage := err.Error()
		if j.RetryCount+1 >= j.MaxRetry {
			err := w.repo.UpdateStatus(ctx, j.UUID, job.StatusFailed, &errMessage)
			if err != nil {
				log.Printf("error update status %s on worker: %v\n", string(job.StatusFailed), err)
			}
		} else {
			err := w.repo.UpdateStatus(ctx, j.UUID, job.StatusPending, &errMessage)
			if err != nil {
				log.Printf("error update status %s on worker: %v\n", string(job.StatusPending), err)
			}
		}
	} else {
		err := w.repo.UpdateStatus(ctx, j.UUID, job.StatusSuccess, nil)
		if err != nil {
			log.Printf("error update status %s on worker: %v\n", string(job.StatusSuccess), err)
		}

	}
}
