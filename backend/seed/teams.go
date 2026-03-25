package seed

import (
	"final-whistle/backend/internal/model"
	"gorm.io/gorm"
)

// SeedTeams creates initial team data
func SeedTeams(db *gorm.DB) error {
	teams := []model.Team{
		{
			Name:      "Manchester City",
			ShortName: stringPtr("MCI"),
			Slug:      "manchester-city",
			LogoURL:   stringPtr("https://example.com/logos/mancity.png"),
		},
		{
			Name:      "Liverpool",
			ShortName: stringPtr("LIV"),
			Slug:      "liverpool",
			LogoURL:   stringPtr("https://example.com/logos/liverpool.png"),
		},
		{
			Name:      "Arsenal",
			ShortName: stringPtr("ARS"),
			Slug:      "arsenal",
			LogoURL:   stringPtr("https://example.com/logos/arsenal.png"),
		},
		{
			Name:      "Chelsea",
			ShortName: stringPtr("CHE"),
			Slug:      "chelsea",
			LogoURL:   stringPtr("https://example.com/logos/chelsea.png"),
		},
		{
			Name:      "Manchester United",
			ShortName: stringPtr("MUN"),
			Slug:      "manchester-united",
			LogoURL:   stringPtr("https://example.com/logos/manutd.png"),
		},
		{
			Name:      "Tottenham Hotspur",
			ShortName: stringPtr("TOT"),
			Slug:      "tottenham-hotspur",
			LogoURL:   stringPtr("https://example.com/logos/tottenham.png"),
		},
	}

	for _, team := range teams {
		if err := db.FirstOrCreate(&team, model.Team{Slug: team.Slug}).Error; err != nil {
			return err
		}
	}

	return nil
}

func stringPtr(s string) *string {
	return &s
}