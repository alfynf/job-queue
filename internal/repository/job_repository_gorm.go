package repository

import (
	"context"
	"log"
	"time"

	"github.com/alfynf/job-queue/internal/job"
	"gorm.io/gorm"
)

type jobRepositoryGorm struct {
	db *gorm.DB
}

func NewJobRepositoryGorm(db *gorm.DB) JobRepository {
	return &jobRepositoryGorm{db: db}
}

func (r *jobRepositoryGorm) Save(ctx context.Context, j job.Job) error {
	return r.db.WithContext(ctx).Create(&j).Error
}

func (r *jobRepositoryGorm) GetByUUID(ctx context.Context, uuid string) (job.Job, error) {
	var j job.Job
	err := r.db.WithContext(ctx).First(&j, "uuid = ?", uuid).Error
	return j, err
}

func (r *jobRepositoryGorm) UpdateStatus(ctx context.Context, uuid string, status job.Status, errMessage *string) error {
	log.Println("masuk db")
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}

	if errMessage != nil {
		updates["last_error_message"] = *errMessage
	}

	log.Println(updates)

	return r.db.WithContext(ctx).Model(&job.Job{}).Where("uuid = ?", uuid).Updates(updates).Error
}

func (r *jobRepositoryGorm) GetPendingJobs(ctx context.Context, limit int) ([]job.Job, error) {
	var jobs []job.Job
	err := r.db.WithContext(ctx).
		Where("status = ?", job.StatusPending).
		Order("scheduled_at NULLS FIRST").
		Limit(limit).
		Find(&jobs).Error
	return jobs, err
}
