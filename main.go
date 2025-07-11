package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/alfynf/job-queue/internal/handler"
	"github.com/alfynf/job-queue/internal/job"
	"github.com/alfynf/job-queue/internal/repository"
	"github.com/alfynf/job-queue/internal/router"
	"github.com/alfynf/job-queue/internal/service"
	"github.com/alfynf/job-queue/internal/worker"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
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

	/*
		=====WHEN HTTP HANDLER WASN'T IMPLEMENTED YET

		j := job.Job{
			UUID:     uuid.New().String(),
			Type:     job.TypeSendingEmail,
			Payload:  map[string]interface{}{"to": "user@example.com"},
			MaxRetry: 3,
		}

		uuid, err := service.SubmitJob(context.Background(), j)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Job submitted id: %s\n", uuid)

		fetchedJob, err := service.GetJobStatus(context.Background(), uuid)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Fetched job: %v\n", fetchedJob)

		worker := worker.New(repo, 3*time.Second, 10)

		worker.Register("sending_mail", func(ctx context.Context, payload map[string]interface{}) error {
			log.Printf("sending mail to %s", payload["to"])
			return nil
		})

		go worker.Start(context.Background())
		time.Sleep(20 * time.Second)

	*/

	repo := repository.NewJobRepositoryGorm(db)
	service := service.NewJobService(repo)
	handler := handler.NewJobHandler(service)

	worker := worker.New(repo, 3*time.Second, 10)

	worker.Register(job.TypeSendingEmail, func(ctx context.Context, payload map[string]interface{}) error {
		log.Printf("sending mail to %s", payload["to"])
		return nil
	})

	go worker.Start(context.Background())

	r := router.SetupRouter(handler)
	r.Run(":8080")

}
