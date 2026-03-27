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
			Slug:      "hot-blooded",
			SortOrder: 1,
			IsActive:  true,
		},
		{
			Name:      "Boring",
			Slug:      "boring",
			SortOrder: 2,
			IsActive:  true,
		},
		{
			Name:      "Tense",
			Slug:      "suffocating",
			SortOrder: 3,
			IsActive:  true,
		},
		{
			Name:      "Classic",
			Slug:      "classic",
			SortOrder: 4,
			IsActive:  true,
		},
		{
			Name:      "Wild",
			Slug:      "unbelievable",
			SortOrder: 5,
			IsActive:  true,
		},
		{
			Name:      "Heartbreaking",
			Slug:      "regrettable",
			SortOrder: 6,
			IsActive:  true,
		},
		{
			Name:      "Dominant",
			Slug:      "dominance",
			SortOrder: 7,
			IsActive:  true,
		},
		{
			Name:      "Painful",
			Slug:      "torture",
			SortOrder: 8,
			IsActive:  true,
		},
		{
			Name:      "Comeback",
			Slug:      "comeback",
			SortOrder: 9,
			IsActive:  true,
		},
		{
			Name:      "Destiny",
			Slug:      "destiny",
			SortOrder: 10,
			IsActive:  true,
		},
	}

	for _, tag := range tags {
		if err := db.Where("slug = ?", tag.Slug).Assign(model.Tag{
			Name:      tag.Name,
			SortOrder: tag.SortOrder,
			IsActive:  tag.IsActive,
		}).FirstOrCreate(&tag).Error; err != nil {
			return err
		}
	}

	return nil
}
