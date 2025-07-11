package handler

import (
	"net/http"

	"github.com/alfynf/job-queue/internal/job"
	"github.com/alfynf/job-queue/internal/service"
	"github.com/gin-gonic/gin"
)

type JobHandler struct {
	service service.JobService
}

func NewJobHandler(s service.JobService) *JobHandler {
	return &JobHandler{service: s}
}

type SubmitJobRequest struct {
	Type    string                 `json:"type"`
	Payload map[string]interface{} `json:"payload"`
}

func (h *JobHandler) SubmitJob(c *gin.Context) {
	var req SubmitJobRequest
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to sumbit job",
			"error":   err.Error(),
		})
		return
	}

	job := job.Job{
		Type:    job.Type(req.Type),
		Payload: req.Payload,
	}

	uuid, err := h.service.SubmitJob(c.Request.Context(), job)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to sumbit job",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Job Submitted",
		"job_uuid": uuid,
	})

}

func (h *JobHandler) GetJobStatus(c *gin.Context) {
	uuid := c.Param("uuid")
	if uuid == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Failed to get job status",
			"error":   "UUID not found",
		})
		return
	}

	job, err := h.service.GetJobStatus(c.Request.Context(), uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to get job status",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success get status job",
		"job":     job,
	})

}
