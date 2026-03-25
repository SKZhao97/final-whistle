package model

import (
	"time"
)

type Session struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null;index"`
	Token     string    `gorm:"size:255;not null;uniqueIndex"`
	ExpiredAt time.Time `gorm:"not null;index"`
	CreatedAt time.Time `gorm:"autoCreateTime"`

	// Relationships
	User User `gorm:"foreignKey:UserID"`
}

// TableName specifies the table name for GORM
func (Session) TableName() string {
	return "sessions"
}