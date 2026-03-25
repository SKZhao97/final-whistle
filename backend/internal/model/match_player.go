package model

type MatchPlayer struct {
	ID       uint `gorm:"primaryKey"`
	MatchID  uint `gorm:"not null;index"`
	PlayerID uint `gorm:"not null;index"`
	TeamID   uint `gorm:"not null;index"`

	// Relationships
	Match  Match  `gorm:"foreignKey:MatchID"`
	Player Player `gorm:"foreignKey:PlayerID"`
	Team   Team   `gorm:"foreignKey:TeamID"`
}

// TableName specifies the table name for GORM
func (MatchPlayer) TableName() string {
	return "match_players"
}