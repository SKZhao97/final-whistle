package seed

import (
	"final-whistle/backend/internal/model"
	"gorm.io/gorm"
)

// SeedMatchPlayers creates match-player relationships for seeded matches
func SeedMatchPlayers(db *gorm.DB) error {
	// Get all seeded matches
	var matches []model.Match
	if err := db.Find(&matches).Error; err != nil {
		return err
	}

	// Get all players grouped by team
	var players []model.Player
	if err := db.Find(&players).Error; err != nil {
		return err
	}

	// Group players by team ID
	playersByTeam := make(map[uint][]model.Player)
	for _, player := range players {
		playersByTeam[player.TeamID] = append(playersByTeam[player.TeamID], player)
	}

	// Create match-player relationships
	for _, match := range matches {
		// Get home team players (first 3-4 players)
		homePlayers := playersByTeam[match.HomeTeamID]
		if len(homePlayers) > 4 {
			homePlayers = homePlayers[:4]
		}

		// Get away team players (first 3-4 players)
		awayPlayers := playersByTeam[match.AwayTeamID]
		if len(awayPlayers) > 4 {
			awayPlayers = awayPlayers[:4]
		}

		// Create match-player entries for home team
		for _, player := range homePlayers {
			matchPlayer := model.MatchPlayer{
				MatchID:  match.ID,
				PlayerID: player.ID,
				TeamID:   match.HomeTeamID,
			}
			if err := db.FirstOrCreate(&matchPlayer, model.MatchPlayer{
				MatchID:  match.ID,
				PlayerID: player.ID,
			}).Error; err != nil {
				return err
			}
		}

		// Create match-player entries for away team
		for _, player := range awayPlayers {
			matchPlayer := model.MatchPlayer{
				MatchID:  match.ID,
				PlayerID: player.ID,
				TeamID:   match.AwayTeamID,
			}
			if err := db.FirstOrCreate(&matchPlayer, model.MatchPlayer{
				MatchID:  match.ID,
				PlayerID: player.ID,
			}).Error; err != nil {
				return err
			}
		}
	}

	return nil
}