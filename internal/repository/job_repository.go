package repository

import (
	"context"

	"github.com/alfynf/job-queue/internal/job"
)

type JobRepository interface {
	Save(ctx context.Context, j job.Job) error
	GetByUUID(ctx context.Context, uuid string) (job.Job, error)
	UpdateStatus(ctx context.Context, uuid string, status job.Status, errMessage *string) error
	GetPendingJobs(ctx context.Context, limit int) ([]job.Job, error)
}
