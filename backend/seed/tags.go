package seed

import (
	"final-whistle/backend/internal/model"
	"gorm.io/gorm"
)

// SeedTags creates predefined tag dictionary
func SeedTags(db *gorm.DB) error {
	tags := []model.Tag{
		{
			Name:      "Intense",
			NameEn:    "Intense",
			NameZh:    "热血",
			Slug:      "hot-blooded",
			SortOrder: 1,
			IsActive:  true,
		},
		{
			Name:      "Boring",
			NameEn:    "Boring",
			NameZh:    "无聊",
			Slug:      "boring",
			SortOrder: 2,
			IsActive:  true,
		},
		{
			Name:      "Tense",
			NameEn:    "Tense",
			NameZh:    "窒息",
			Slug:      "suffocating",
			SortOrder: 3,
			IsActive:  true,
		},
		{
			Name:      "Classic",
			NameEn:    "Classic",
			NameZh:    "经典",
			Slug:      "classic",
			SortOrder: 4,
			IsActive:  true,
		},
		{
			Name:      "Wild",
			NameEn:    "Wild",
			NameZh:    "离谱",
			Slug:      "unbelievable",
			SortOrder: 5,
			IsActive:  true,
		},
		{
			Name:      "Heartbreaking",
			NameEn:    "Heartbreaking",
			NameZh:    "遗憾",
			Slug:      "regrettable",
			SortOrder: 6,
			IsActive:  true,
		},
		{
			Name:      "Dominant",
			NameEn:    "Dominant",
			NameZh:    "统治",
			Slug:      "dominance",
			SortOrder: 7,
			IsActive:  true,
		},
		{
			Name:      "Painful",
			NameEn:    "Painful",
			NameZh:    "折磨",
			Slug:      "torture",
			SortOrder: 8,
			IsActive:  true,
		},
		{
			Name:      "Comeback",
			NameEn:    "Comeback",
			NameZh:    "逆转",
			Slug:      "comeback",
			SortOrder: 9,
			IsActive:  true,
		},
		{
			Name:      "Destiny",
			NameEn:    "Destiny",
			NameZh:    "宿命",
			Slug:      "destiny",
			SortOrder: 10,
			IsActive:  true,
		},
	}

	for _, tag := range tags {
		if err := db.Where("slug = ?", tag.Slug).Assign(model.Tag{
			Name:      tag.Name,
			NameEn:    tag.NameEn,
			NameZh:    tag.NameZh,
			SortOrder: tag.SortOrder,
			IsActive:  tag.IsActive,
		}).FirstOrCreate(&tag).Error; err != nil {
			return err
		}
	}

	return nil
}
