package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/alfynf/job-queue/internal/handler"
	"github.com/alfynf/job-queue/internal/job"
	"github.com/alfynf/job-queue/internal/repository"
	"github.com/alfynf/job-queue/internal/router"
	"github.com/alfynf/job-queue/internal/service"
	"github.com/alfynf/job-queue/internal/worker"
	"github.com/alfynf/job-queue/pkg/logger"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	logger.Init()
	defer logger.Sync()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db_dsn := os.Getenv("DB_DSN")
	db, err := gorm.Open(postgres.Open(db_dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	db.AutoMigrate(&job.Job{})

	repo := repository.NewJobRepositoryGorm(db)
	service := service.NewJobService(repo)
	handler := handler.NewJobHandler(service)

	worker := worker.New(repo, 3*time.Second, 10)

	worker.Register(job.TypeSendingEmail, func(ctx context.Context, payload map[string]interface{}) error {
		log.Printf("sending mail to %s", payload["to"])
		time.Sleep(2 * time.Second)
		return nil
	})

	worker.Register(job.TypeGeneratePdf, func(ctx context.Context, payload map[string]interface{}) error {
		log.Printf("Generating pdf %s", payload["title"])
		return fmt.Errorf("mock error to test retry")
	})

	go worker.Start(context.Background())

	r := router.SetupRouter(handler)
	r.Run(":8080")

}
