package router

import (
	"github.com/alfynf/job-queue/internal/handler"
	"github.com/alfynf/job-queue/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter(jobHandler *handler.JobHandler) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Logging())

	r.POST("/jobs", jobHandler.SubmitJob)
	r.GET("/jobs/:uuid", jobHandler.GetJobById)

	return r
}
