package repository

import (
	"context"

	"github.com/alfynf/job-queue/internal/job"
)

type JobRepository interface {
	Save(ctx context.Context, j job.Job) error
	GetByUUID(ctx context.Context, uuid string) (job.Job, error)
	Update(ctx context.Context, uuid string, job job.Job) error
	GetJobsByStatus(ctx context.Context, status job.Status, limit int) ([]job.Job, error)
}
