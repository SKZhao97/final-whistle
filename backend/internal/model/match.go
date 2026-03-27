package model

import (
	"time"
)

type MatchStatus string

const (
	MatchStatusScheduled MatchStatus = "SCHEDULED"
	MatchStatusFinished  MatchStatus = "FINISHED"
)

type Match struct {
	ID                uint        `gorm:"primaryKey"`
	Competition       string      `gorm:"size:100;not null"`
	Season            string      `gorm:"size:50;not null"`
	Round             *string     `gorm:"size:50"`
	Status            MatchStatus `gorm:"size:20;not null;index"`
	KickoffAt         time.Time   `gorm:"not null;index"`
	HomeTeamID        uint        `gorm:"not null;index"`
	AwayTeamID        uint        `gorm:"not null;index"`
	HomeScore         *int        `gorm:"type:integer"`
	AwayScore         *int        `gorm:"type:integer"`
	Venue             *string     `gorm:"size:200"`
	ExternalSource    *string     `gorm:"size:50"`
	ExternalID        *string     `gorm:"size:100"`
	ExternalUpdatedAt *time.Time
	CreatedAt         time.Time `gorm:"autoCreateTime"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime"`

	// Relationships
	HomeTeam     Team          `gorm:"foreignKey:HomeTeamID"`
	AwayTeam     Team          `gorm:"foreignKey:AwayTeamID"`
	MatchPlayers []MatchPlayer `gorm:"foreignKey:MatchID"`
	CheckIns     []CheckIn     `gorm:"foreignKey:MatchID"`
}

// TableName specifies the table name for GORM
func (Match) TableName() string {
	return "matches"
}
