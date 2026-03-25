package model

type Tag struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:50;not null"`
	Slug      string `gorm:"size:50;not null;uniqueIndex"`
	SortOrder int    `gorm:"not null;default:0;index"`
	IsActive  bool   `gorm:"not null;default:true;index"`
}

// TableName specifies the table name for GORM
func (Tag) TableName() string {
	return "tags"
}