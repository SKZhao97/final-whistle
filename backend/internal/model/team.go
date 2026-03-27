package model

import (
	"time"
)

type Team struct {
	ID                uint    `gorm:"primaryKey"`
	Name              string  `gorm:"size:100;not null"`
	ShortName         *string `gorm:"size:50"`
	Slug              string  `gorm:"size:100;not null;uniqueIndex"`
	LogoURL           *string `gorm:"type:text"`
	ExternalSource    *string `gorm:"size:50"`
	ExternalID        *string `gorm:"size:100"`
	ExternalUpdatedAt *time.Time
	CreatedAt         time.Time `gorm:"autoCreateTime"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime"`

	// Relationships
	HomeMatches []Match  `gorm:"foreignKey:HomeTeamID"`
	AwayMatches []Match  `gorm:"foreignKey:AwayTeamID"`
	Players     []Player `gorm:"foreignKey:TeamID"`
}

// TableName specifies the table name for GORM
func (Team) TableName() string {
	return "teams"
}
