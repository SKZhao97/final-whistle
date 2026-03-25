package seed

import (
	"final-whistle/backend/internal/model"
	"fmt"
	"gorm.io/gorm"
)

// ValidateSeedData performs basic validation of seeded data integrity
func ValidateSeedData(db *gorm.DB) error {
	// Check each table has data
	tables := []struct {
		name   string
		model  interface{}
		minCount int
	}{
		{"teams", &model.Team{}, 3},
		{"players", &model.Player{}, 8},
		{"matches", &model.Match{}, 3},
		{"match_players", &model.MatchPlayer{}, 10},
		{"tags", &model.Tag{}, 5},
		{"users", &model.User{}, 0}, // Users may not be seeded in base data
	}

	for _, table := range tables {
		var count int64
		if err := db.Model(table.model).Count(&count).Error; err != nil {
			return fmt.Errorf("failed to count %s: %v", table.name, err)
		}

		if count < int64(table.minCount) {
			return fmt.Errorf("%s has insufficient data: got %d, expected at least %d",
				table.name, count, table.minCount)
		}
	}

	// Validate foreign key relationships
	// Check that all match_players have valid match and player references
	var invalidMatchPlayers []model.MatchPlayer
	if err := db.Joins("LEFT JOIN matches ON match_players.match_id = matches.id").
		Joins("LEFT JOIN players ON match_players.player_id = players.id").
		Where("matches.id IS NULL OR players.id IS NULL").
		Find(&invalidMatchPlayers).Error; err != nil {
		return fmt.Errorf("failed to validate match_players foreign keys: %v", err)
	}

	if len(invalidMatchPlayers) > 0 {
		return fmt.Errorf("found %d match_players with invalid foreign key references",
			len(invalidMatchPlayers))
	}

	// Check that all players have valid team references
	var invalidPlayers []model.Player
	if err := db.Joins("LEFT JOIN teams ON players.team_id = teams.id").
		Where("teams.id IS NULL").
		Find(&invalidPlayers).Error; err != nil {
		return fmt.Errorf("failed to validate players foreign keys: %v", err)
	}

	if len(invalidPlayers) > 0 {
		return fmt.Errorf("found %d players with invalid team references",
			len(invalidPlayers))
	}

	// Check that all matches have valid home and away team references
	var invalidMatches []model.Match
	if err := db.Joins("LEFT JOIN teams AS home_teams ON matches.home_team_id = home_teams.id").
		Joins("LEFT JOIN teams AS away_teams ON matches.away_team_id = away_teams.id").
		Where("home_teams.id IS NULL OR away_teams.id IS NULL").
		Find(&invalidMatches).Error; err != nil {
		return fmt.Errorf("failed to validate matches foreign keys: %v", err)
	}

	if len(invalidMatches) > 0 {
		return fmt.Errorf("found %d matches with invalid team references",
			len(invalidMatches))
	}

	return nil
}