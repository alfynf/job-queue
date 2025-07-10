package service

import (
	"context"
	"time"

	"github.com/alfynf/job-queue/internal/job"
	"github.com/alfynf/job-queue/internal/repository"
)

type JobService interface {
	SubmitJob(ctx context.Context, j job.Job) (string, error)
	GetJobStatus(ctx context.Context, uuid string) (job.Job, error)
}

type jobService struct {
	repo repository.JobRepository
}

func NewJobService(repo repository.JobRepository) JobService {
	return &jobService{
		repo: repo,
	}
}

func (s *jobService) SubmitJob(ctx context.Context, j job.Job) (string, error) {
	j.Status = job.StatusPending
	j.CreatedAt = time.Now()
	j.UpdatedAt = time.Now()

	err := s.repo.Save(ctx, j)
	if err != nil {
		return "", err
	}
	return j.UUID, nil
}

func (s *jobService) GetJobStatus(ctx context.Context, uuid string) (job.Job, error) {
	return job.Job{}, nil
}
