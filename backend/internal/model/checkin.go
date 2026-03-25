package model

import (
	"time"
)

type WatchedType string

const (
	WatchedTypeFull       WatchedType = "FULL"
	WatchedTypePartial    WatchedType = "PARTIAL"
	WatchedTypeHighlights WatchedType = "HIGHLIGHTS"
)

type SupporterSide string

const (
	SupporterSideHome    SupporterSide = "HOME"
	SupporterSideAway    SupporterSide = "AWAY"
	SupporterSideNeutral SupporterSide = "NEUTRAL"
)

type CheckIn struct {
	ID              uint          `gorm:"primaryKey"`
	UserID          uint          `gorm:"not null;index"`
	MatchID         uint          `gorm:"not null;index"`
	WatchedType     WatchedType   `gorm:"size:20;not null"`
	SupporterSide   SupporterSide `gorm:"size:20;not null"`
	MatchRating     int           `gorm:"not null;check:match_rating BETWEEN 1 AND 10"`
	HomeTeamRating  int           `gorm:"not null;check:home_team_rating BETWEEN 1 AND 10"`
	AwayTeamRating  int           `gorm:"not null;check:away_team_rating BETWEEN 1 AND 10"`
	ShortReview     *string       `gorm:"size:280"`
	WatchedAt       time.Time     `gorm:"not null"`
	CreatedAt       time.Time     `gorm:"autoCreateTime;index"`
	UpdatedAt       time.Time     `gorm:"autoUpdateTime"`

	// Relationships
	User          User            `gorm:"foreignKey:UserID"`
	Match         Match           `gorm:"foreignKey:MatchID"`
	PlayerRatings []PlayerRating  `gorm:"foreignKey:CheckInID"`
	Tags          []Tag           `gorm:"many2many:checkin_tags;foreignKey:ID;joinForeignKey:CheckInID;joinReferences:TagID"`
}

// TableName specifies the table name for GORM
func (CheckIn) TableName() string {
	return "check_ins"
}