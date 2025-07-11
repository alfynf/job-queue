package repository

import (
	"context"

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

func (r *jobRepositoryGorm) Update(ctx context.Context, uuid string, job job.Job) error {
	return r.db.WithContext(ctx).Save(&job).Error
}

func (r *jobRepositoryGorm) GetJobsByStatus(ctx context.Context, status job.Status, limit int) ([]job.Job, error) {
	var jobs []job.Job
	err := r.db.WithContext(ctx).
		Where("status = ?", status).
		Order("scheduled_at NULLS FIRST").
		Limit(limit).
		Find(&jobs).Error
	return jobs, err
}
