package seed

import (
	"final-whistle/backend/internal/model"
	"gorm.io/gorm"
)

// SeedTags creates predefined tag dictionary
func SeedTags(db *gorm.DB) error {
	tags := []model.Tag{
		{
			Name:      "热血",
			Slug:      "hot-blooded",
			SortOrder: 1,
			IsActive:  true,
		},
		{
			Name:      "无聊",
			Slug:      "boring",
			SortOrder: 2,
			IsActive:  true,
		},
		{
			Name:      "窒息",
			Slug:      "suffocating",
			SortOrder: 3,
			IsActive:  true,
		},
		{
			Name:      "经典",
			Slug:      "classic",
			SortOrder: 4,
			IsActive:  true,
		},
		{
			Name:      "离谱",
			Slug:      "unbelievable",
			SortOrder: 5,
			IsActive:  true,
		},
		{
			Name:      "可惜",
			Slug:      "regrettable",
			SortOrder: 6,
			IsActive:  true,
		},
		{
			Name:      "统治力",
			Slug:      "dominance",
			SortOrder: 7,
			IsActive:  true,
		},
		{
			Name:      "折磨",
			Slug:      "torture",
			SortOrder: 8,
			IsActive:  true,
		},
		{
			Name:      "逆转",
			Slug:      "comeback",
			SortOrder: 9,
			IsActive:  true,
		},
		{
			Name:      "宿命感",
			Slug:      "destiny",
			SortOrder: 10,
			IsActive:  true,
		},
	}

	for _, tag := range tags {
		if err := db.FirstOrCreate(&tag, model.Tag{Slug: tag.Slug}).Error; err != nil {
			return err
		}
	}

	return nil
}