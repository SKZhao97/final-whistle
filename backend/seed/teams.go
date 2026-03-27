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
			NameZh:    stringPtr("曼城"),
			ShortName: stringPtr("MCI"),
			Slug:      "manchester-city",
			LogoURL:   stringPtr("https://example.com/logos/mancity.png"),
		},
		{
			Name:      "Liverpool",
			NameZh:    stringPtr("利物浦"),
			ShortName: stringPtr("LIV"),
			Slug:      "liverpool",
			LogoURL:   stringPtr("https://example.com/logos/liverpool.png"),
		},
		{
			Name:      "Arsenal",
			NameZh:    stringPtr("阿森纳"),
			ShortName: stringPtr("ARS"),
			Slug:      "arsenal",
			LogoURL:   stringPtr("https://example.com/logos/arsenal.png"),
		},
		{
			Name:      "Chelsea",
			NameZh:    stringPtr("切尔西"),
			ShortName: stringPtr("CHE"),
			Slug:      "chelsea",
			LogoURL:   stringPtr("https://example.com/logos/chelsea.png"),
		},
		{
			Name:      "Manchester United",
			NameZh:    stringPtr("曼联"),
			ShortName: stringPtr("MUN"),
			Slug:      "manchester-united",
			LogoURL:   stringPtr("https://example.com/logos/manutd.png"),
		},
		{
			Name:      "Tottenham Hotspur",
			NameZh:    stringPtr("热刺"),
			ShortName: stringPtr("TOT"),
			Slug:      "tottenham-hotspur",
			LogoURL:   stringPtr("https://example.com/logos/tottenham.png"),
		},
	}

	for _, team := range teams {
		existing := model.Team{}
		if err := db.Where("slug = ?", team.Slug).
			Assign(model.Team{
				Name:      team.Name,
				NameZh:    team.NameZh,
				ShortName: team.ShortName,
				LogoURL:   team.LogoURL,
			}).
			FirstOrCreate(&existing, model.Team{Slug: team.Slug}).Error; err != nil {
			return err
		}
	}

	return nil
}

func stringPtr(s string) *string {
	return &s
}
