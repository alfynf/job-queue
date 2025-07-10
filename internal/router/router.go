package router

import (
	"github.com/alfynf/job-queue/internal/handler"
	"github.com/gin-gonic/gin"
)

func SetupRouter(jobHandler *handler.JobHandler) *gin.Engine {
	r := gin.Default()

	r.POST("/jobs", jobHandler.SubmitJob)

	return r
}
