package model

import (
	"time"
)

type SyncJobStatus string

const (
	SyncJobStatusPending   SyncJobStatus = "pending"
	SyncJobStatusRunning   SyncJobStatus = "running"
	SyncJobStatusSucceeded SyncJobStatus = "succeeded"
	SyncJobStatusFailed    SyncJobStatus = "failed"
	SyncJobStatusCanceled  SyncJobStatus = "canceled"
)

type SyncTriggerMode string

const (
	SyncTriggerModeAutomatic SyncTriggerMode = "automatic"
	SyncTriggerModeManual    SyncTriggerMode = "manual"
)

type SyncJob struct {
	ID          uint            `gorm:"primaryKey"`
	JobType     string          `gorm:"size:50;not null"`
	ScopeType   string          `gorm:"size:50;not null"`
	ScopeKey    string          `gorm:"size:200;not null"`
	DedupeKey   string          `gorm:"size:255;not null"`
	TriggerMode SyncTriggerMode `gorm:"size:20;not null"`
	Priority    int             `gorm:"not null;default:100"`
	Status      SyncJobStatus   `gorm:"size:20;not null"`
	ScheduledAt time.Time       `gorm:"not null"`
	StartedAt   *time.Time
	FinishedAt  *time.Time
	Attempt     int    `gorm:"not null;default:0"`
	MaxAttempts int    `gorm:"not null;default:3"`
	Payload     []byte `gorm:"type:jsonb;not null"`
	LastError   *string
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

func (SyncJob) TableName() string {
	return "sync_jobs"
}
