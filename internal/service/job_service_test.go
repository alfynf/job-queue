package service

import (
	"context"
	"testing"
	"time"

	"github.com/alfynf/job-queue/internal/job"
	mockRepo "github.com/alfynf/job-queue/internal/repository/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSubmitJob(t *testing.T) {
	ctx := context.Background()
	mockJobRepo := new(mockRepo.JobRepository)
	jobService := NewJobService(mockJobRepo)

	inputJob := job.Job{
		UUID:      "f6c04811-d7f0-49d3-b345-16be7cab99f8",
		Type:      int(job.TypeSendingEmail),
		Payload:   map[string]interface{}{"to": "user@mail.com"},
		MaxRetry:  3,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Status:    job.StatusPending,
	}

	mockJobRepo.On("Save", ctx, mock.Anything).Return(nil)

	uuid, err := jobService.SubmitJob(ctx, inputJob)

	assert.NoError(t, err)
	assert.Equal(t, inputJob.UUID, uuid)
	// mockJobRepo.AssertCalled(t, "Save", ctx, inputJob)
}
