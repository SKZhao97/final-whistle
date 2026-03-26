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
		// Chelsea players
		{
			TeamID:   teams["chelsea"],
			Name:     "Cole Palmer",
			Slug:     "cole-palmer",
			Position: stringPtr("Forward"),
		},
		{
			TeamID:   teams["chelsea"],
			Name:     "Enzo Fernandez",
			Slug:     "enzo-fernandez",
			Position: stringPtr("Midfielder"),
		},
		{
			TeamID:   teams["chelsea"],
			Name:     "Reece James",
			Slug:     "reece-james",
			Position: stringPtr("Defender"),
		},
		// Manchester United players
		{
			TeamID:   teams["manchester-united"],
			Name:     "Bruno Fernandes",
			Slug:     "bruno-fernandes",
			Position: stringPtr("Midfielder"),
		},
		{
			TeamID:   teams["manchester-united"],
			Name:     "Marcus Rashford",
			Slug:     "marcus-rashford",
			Position: stringPtr("Forward"),
		},
		{
			TeamID:   teams["manchester-united"],
			Name:     "Lisandro Martinez",
			Slug:     "lisandro-martinez",
			Position: stringPtr("Defender"),
		},
		// Tottenham players
		{
			TeamID:   teams["tottenham-hotspur"],
			Name:     "Son Heung-min",
			Slug:     "son-heung-min",
			Position: stringPtr("Forward"),
		},
		{
			TeamID:   teams["tottenham-hotspur"],
			Name:     "James Maddison",
			Slug:     "james-maddison",
			Position: stringPtr("Midfielder"),
		},
		{
			TeamID:   teams["tottenham-hotspur"],
			Name:     "Cristian Romero",
			Slug:     "cristian-romero",
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
