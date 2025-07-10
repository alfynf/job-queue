package job

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type Status string
type Type string

const (
	StatusPending Status = "pending"
	StatusRunning Status = "running"
	StatusSuccess Status = "success"
	StatusFailed  Status = "failed"
)

const (
	TypeSendingEmail Type = "sending_mail"
	TypeGeneratePdf  Type = "generate_pdf"
)

type Job struct {
	UUID             string `gorm:"primaryKey;default:uuid_generate_v4()"`
	Type             Type
	Payload          JSONB `gorm:"type:jsonb"`
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

// JSONB Interface for JSONB Field of Job Table

type JSONB map[string]interface{}

func (j *JSONB) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, j)
}

func (j JSONB) Value() (driver.Value, error) {
	return json.Marshal(j)
}
