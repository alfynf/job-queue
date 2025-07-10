package job

import "time"

type Status string
type Type int

const (
	StatusPending Status = "pending"
	StatusRunning Status = "running"
	StatusSuccess Status = "success"
	StatusFailed  Status = "failed"
)

const (
	TypeSendingEmail Type = iota + 1
	TypeGeneratePdf
)

type Job struct {
	UUID             string
	Type             int
	Payload          map[string]interface{}
	Status           Status
	RetryCount       int
	MaxRetry         int
	LastErrorMessage *string
	ScheduledAt      *time.Time
	StartedAt        *time.Time
	FinishedAt       *time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
