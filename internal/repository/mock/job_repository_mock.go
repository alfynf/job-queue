package mock

import (
	"context"

	"github.com/alfynf/job-queue/internal/job"
	"github.com/stretchr/testify/mock"
)

type JobRepository struct {
	mock.Mock
}

func (m *JobRepository) Save(ctx context.Context, j job.Job) error {
	args := m.Called(ctx, j)
	return args.Error(0)
}

func (m *JobRepository) GetByUUID(ctx context.Context, uuid string) (job.Job, error) {
	args := m.Called(ctx, uuid)
	return args.Get(0).(job.Job), args.Error(1)
}

func (m *JobRepository) UpdateStatus(ctx context.Context, uuid string, status job.Status, errMessage *string) error {
	args := m.Called(ctx, uuid, status, errMessage)
	return args.Error(0)
}

func (m *JobRepository) GetPendingJobs(ctx context.Context, limit int) ([]job.Job, error) {
	args := m.Called(ctx, limit)
	return []job.Job{}, args.Error(0)

}
