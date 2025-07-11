package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alfynf/job-queue/internal/job"
	mockSvc "github.com/alfynf/job-queue/internal/service/mock"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSubmitJobSuccess(t *testing.T) {

	gin.SetMode(gin.TestMode)
	ctx := context.Background()
	mockService := new(mockSvc.JobService)

	uuid := uuid.New()

	expectedJob := job.Job{UUID: uuid}
	mockService.On("SubmitJob", ctx, mock.Anything).Return(expectedJob.UUID.String(), nil)

	jobHandler := NewJobHandler(mockService)
	router := gin.Default()
	router.POST("/jobs", jobHandler.SubmitJob)

	body := map[string]interface{}{
		"type":    "send_email",
		"payload": map[string]interface{}{"to": "user@example.com"},
	}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/jobs", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	mockService.AssertExpectations(t)
}

func TestSubmitJobInvalidReqBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx := context.Background()
	mockService := new(mockSvc.JobService)

	uuid := uuid.New()

	expectedJob := job.Job{UUID: uuid}
	mockService.On("SubmitJob", ctx, mock.Anything).Return(expectedJob.UUID, nil)

	jobHandler := NewJobHandler(mockService)
	router := gin.Default()
	router.POST("/jobs", jobHandler.SubmitJob)

	body := `{"type": 123}`
	req, _ := http.NewRequest("POST", "/jobs", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

}

func TestGetJobStatusSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx := context.Background()
	mockService := new(mockSvc.JobService)

	uuid := uuid.New()

	expectedJob := job.Job{
		UUID:   uuid,
		Status: job.StatusPending,
		Type:   job.TypeSendingEmail,
	}

	mockService.On("GetJobStatus", ctx, mock.Anything).Return(expectedJob, nil)

	jobHandler := NewJobHandler(mockService)
	router := gin.Default()
	path := fmt.Sprintf("/jobs/%s", uuid.String())
	method := "GET"
	router.GET("/jobs/:uuid", jobHandler.GetJobById)

	req, _ := http.NewRequest(method, path, nil)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	type resBody struct {
		UUID string  `json:"uuid"`
		Job  job.Job `json:"job"`
	}
	var resp resBody
	_ = json.NewDecoder(rec.Body).Decode(&resp)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, job.StatusPending, resp.Job.Status)
	mockService.AssertExpectations(t)

}
