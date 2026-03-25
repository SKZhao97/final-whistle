package seed

import (
	"final-whistle/backend/internal/model"
	"gorm.io/gorm"
)

// SeedPlayers creates initial player data
func SeedPlayers(db *gorm.DB) error {
	// Get team IDs by slug
	teams := make(map[string]uint)
	var teamList []model.Team
	if err := db.Find(&teamList).Error; err != nil {
		return err
	}
	for _, team := range teamList {
		teams[team.Slug] = team.ID
	}

	players := []model.Player{
		// Manchester City players
		{
			TeamID:   teams["manchester-city"],
			Name:     "Kevin De Bruyne",
			Slug:     "kevin-de-bruyne",
			Position: stringPtr("Midfielder"),
		},
		{
			TeamID:   teams["manchester-city"],
			Name:     "Erling Haaland",
			Slug:     "erling-haaland",
			Position: stringPtr("Forward"),
		},
		{
			TeamID:   teams["manchester-city"],
			Name:     "Phil Foden",
			Slug:     "phil-foden",
			Position: stringPtr("Midfielder"),
		},
		// Liverpool players
		{
			TeamID:   teams["liverpool"],
			Name:     "Mohamed Salah",
			Slug:     "mohamed-salah",
			Position: stringPtr("Forward"),
		},
		{
			TeamID:   teams["liverpool"],
			Name:     "Virgil van Dijk",
			Slug:     "virgil-van-dijk",
			Position: stringPtr("Defender"),
		},
		{
			TeamID:   teams["liverpool"],
			Name:     "Trent Alexander-Arnold",
			Slug:     "trent-alexander-arnold",
			Position: stringPtr("Defender"),
		},
		// Arsenal players
		{
			TeamID:   teams["arsenal"],
			Name:     "Bukayo Saka",
			Slug:     "bukayo-saka",
			Position: stringPtr("Forward"),
		},
		{
			TeamID:   teams["arsenal"],
			Name:     "Martin Ødegaard",
			Slug:     "martin-odegaard",
			Position: stringPtr("Midfielder"),
		},
		{
			TeamID:   teams["arsenal"],
			Name:     "William Saliba",
			Slug:     "william-saliba",
			Position: stringPtr("Defender"),
		},
	}

	for _, player := range players {
		if err := db.FirstOrCreate(&player, model.Player{Slug: player.Slug}).Error; err != nil {
			return err
		}
	}

	return nil
}