package model

type CheckInTag struct {
	ID        uint `gorm:"primaryKey"`
	CheckInID uint `gorm:"not null;index"`
	TagID     uint `gorm:"not null;index"`

	// Relationships
	CheckIn CheckIn `gorm:"foreignKey:CheckInID"`
	Tag     Tag     `gorm:"foreignKey:TagID"`
}

// TableName specifies the table name for GORM
func (CheckInTag) TableName() string {
	return "checkin_tags"
}