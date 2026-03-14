package database

import (
	"time"

	"gorm.io/gorm"
)

type JobStatus string

const (
	StatusPending    JobStatus = "PENDING"
	StatusProcessing JobStatus = "PROCESSING"
	StatusCompleted  JobStatus = "COMPLETED"
	StatusFailed     JobStatus = "FAILED"
)

type Job struct {
	ID          string         `gorm:"primaryKey;type:uuid" json:"id"`
	Name        string         `gorm:"index;not null" json:"name"`
	Payload     string         `gorm:"type:jsonb" json:"payload"`
	Status      JobStatus      `gorm:"index;default:'PENDING'" json:"status"`
	Retries     int            `gorm:"default:0" json:"retries"`
	MaxRetries  int            `gorm:"default:3" json:"max_retries"`
	LastError   string         `gorm:"type:text" json:"last_error,omitempty"`
	ScheduledAt time.Time      `gorm:"index" json:"scheduled_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type JobLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	JobID     string    `gorm:"index;type:uuid;not null" json:"job_id"`
	Message   string    `gorm:"type:text" json:"message"`
	Level     string    `json:"level"` // info, error
	CreatedAt time.Time `json:"created_at"`
}
