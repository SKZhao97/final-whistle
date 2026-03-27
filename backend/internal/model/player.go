package model

import (
	"time"
)

type Player struct {
	ID                uint    `gorm:"primaryKey"`
	TeamID            uint    `gorm:"not null;index"`
	Name              string  `gorm:"size:100;not null"`
	Slug              string  `gorm:"size:100;not null;uniqueIndex"`
	Position          *string `gorm:"size:50"`
	AvatarURL         *string `gorm:"type:text"`
	ExternalSource    *string `gorm:"size:50"`
	ExternalID        *string `gorm:"size:100"`
	ExternalUpdatedAt *time.Time
	CreatedAt         time.Time `gorm:"autoCreateTime"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime"`

	// Relationships
	Team          Team           `gorm:"foreignKey:TeamID"`
	MatchPlayers  []MatchPlayer  `gorm:"foreignKey:PlayerID"`
	PlayerRatings []PlayerRating `gorm:"foreignKey:PlayerID"`
}

// TableName specifies the table name for GORM
func (Player) TableName() string {
	return "players"
}
