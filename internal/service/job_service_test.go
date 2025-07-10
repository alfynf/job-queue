package service

import (
	"context"
	"testing"
	"time"

	"github.com/alfynf/job-queue/internal/job"
	mockRepo "github.com/alfynf/job-queue/internal/repository/mock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSubmitJob(t *testing.T) {
	ctx := context.Background()
	mockJobRepo := new(mockRepo.JobRepository)
	jobService := NewJobService(mockJobRepo)

	inputJob := job.Job{
		Type:      job.TypeSendingEmail,
		Payload:   map[string]interface{}{"to": "user@mail.com"},
		MaxRetry:  3,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Status:    job.StatusPending,
	}

	// mockJobRepo.On("Save", ctx, mock.Anything).Return(nil) <- bypassing all input content
	mockJobRepo.On("Save", ctx, mock.MatchedBy(func(j job.Job) bool {
		return j.Status == job.StatusPending
	})).Return(nil)

	uuid, err := jobService.SubmitJob(ctx, inputJob)

	assert.NoError(t, err)
	assert.NotNil(t, uuid)
	mockJobRepo.AssertCalled(t, "Save", ctx, mock.Anything)

}

func TestGetJobStatus(t *testing.T) {
	ctx := context.Background()
	mockJobRepo := new(mockRepo.JobRepository)
	jobService := NewJobService(mockJobRepo)

	uuid := uuid.New()

	expectedJob := job.Job{
		UUID:   uuid,
		Status: job.StatusSuccess,
	}

	mockJobRepo.On("GetByUUID", ctx, uuid.String()).Return(expectedJob, nil)

	actualJob, err := jobService.GetJobStatus(ctx, uuid.String())

	assert.NoError(t, err)
	assert.Equal(t, expectedJob.UUID.String(), actualJob.UUID.String())
	assert.Equal(t, expectedJob.Status, actualJob.Status)
	mockJobRepo.AssertCalled(t, "GetByUUID", ctx, uuid.String())
}
