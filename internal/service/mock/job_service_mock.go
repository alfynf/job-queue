package mock

import (
	"context"

	"github.com/alfynf/job-queue/internal/job"
	"github.com/stretchr/testify/mock"
)

type JobService struct {
	mock.Mock
}

func (m *JobService) SubmitJob(ctx context.Context, j job.Job) (string, error) {
	args := m.Called(ctx, j)
	return args.String(0), args.Error(1)
}

func (m *JobService) GetJobStatus(ctx context.Context, uuid string) (job.Job, error) {
	args := m.Called(ctx, uuid)
	return args.Get(0).(job.Job), args.Error(1)
}
