package model

import (
	"time"
)

type SyncCursor struct {
	ID            uint `gorm:"primaryKey"`
	Provider      string
	ResourceType  string
	ScopeKey      string
	LastSuccessAt *time.Time
	LastAttemptAt *time.Time
	LastErrorAt   *time.Time
	LastError     *string
	Metadata      []byte    `gorm:"type:jsonb;not null"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}

func (SyncCursor) TableName() string {
	return "sync_cursors"
}
