package model

type PlayerRating struct {
	ID        uint    `gorm:"primaryKey"`
	CheckInID uint    `gorm:"not null;index"`
	PlayerID  uint    `gorm:"not null;index"`
	Rating    int     `gorm:"not null;check:rating BETWEEN 1 AND 10"`
	Note      *string `gorm:"size:80"`

	// Relationships
	CheckIn CheckIn `gorm:"foreignKey:CheckInID"`
	Player  Player  `gorm:"foreignKey:PlayerID"`
}

// TableName specifies the table name for GORM
func (PlayerRating) TableName() string {
	return "player_ratings"
}