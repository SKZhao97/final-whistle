package seed

import (
	"final-whistle/backend/internal/model"
	"gorm.io/gorm"
	"time"
)

// SeedMatches creates initial match data
func SeedMatches(db *gorm.DB) error {
	// Get team IDs by slug
	teams := make(map[string]uint)
	var teamList []model.Team
	if err := db.Find(&teamList).Error; err != nil {
		return err
	}
	for _, team := range teamList {
		teams[team.Slug] = team.ID
	}

	matches := []model.Match{
		{
			Competition: "Premier League",
			Season:      "2024-2025",
			Round:       stringPtr("Matchday 1"),
			Status:      model.MatchStatusFinished,
			KickoffAt:   time.Date(2024, 8, 10, 15, 0, 0, 0, time.UTC),
			HomeTeamID:  teams["manchester-city"],
			AwayTeamID:  teams["liverpool"],
			HomeScore:   intPtr(2),
			AwayScore:   intPtr(2),
			Venue:       stringPtr("Etihad Stadium"),
		},
		{
			Competition: "Premier League",
			Season:      "2024-2025",
			Round:       stringPtr("Matchday 1"),
			Status:      model.MatchStatusFinished,
			KickoffAt:   time.Date(2024, 8, 10, 12, 30, 0, 0, time.UTC),
			HomeTeamID:  teams["arsenal"],
			AwayTeamID:  teams["chelsea"],
			HomeScore:   intPtr(3),
			AwayScore:   intPtr(1),
			Venue:       stringPtr("Emirates Stadium"),
		},
		{
			Competition: "Premier League",
			Season:      "2024-2025",
			Round:       stringPtr("Matchday 2"),
			Status:      model.MatchStatusFinished,
			KickoffAt:   time.Date(2024, 8, 17, 17, 30, 0, 0, time.UTC),
			HomeTeamID:  teams["manchester-united"],
			AwayTeamID:  teams["tottenham-hotspur"],
			HomeScore:   intPtr(1),
			AwayScore:   intPtr(0),
			Venue:       stringPtr("Old Trafford"),
		},
		{
			Competition: "Premier League",
			Season:      "2024-2025",
			Round:       stringPtr("Matchday 2"),
			Status:      model.MatchStatusFinished,
			KickoffAt:   time.Date(2024, 8, 17, 15, 0, 0, 0, time.UTC),
			HomeTeamID:  teams["liverpool"],
			AwayTeamID:  teams["arsenal"],
			HomeScore:   intPtr(2),
			AwayScore:   intPtr(1),
			Venue:       stringPtr("Anfield"),
		},
		{
			Competition: "Premier League",
			Season:      "2024-2025",
			Round:       stringPtr("Matchday 3"),
			Status:      model.MatchStatusScheduled,
			KickoffAt:   time.Date(2024, 8, 24, 15, 0, 0, 0, time.UTC),
			HomeTeamID:  teams["chelsea"],
			AwayTeamID:  teams["manchester-city"],
			HomeScore:   nil,
			AwayScore:   nil,
			Venue:       stringPtr("Stamford Bridge"),
		},
		{
			Competition: "Premier League",
			Season:      "2024-2025",
			Round:       stringPtr("Matchday 3"),
			Status:      model.MatchStatusFinished,
			KickoffAt:   time.Date(2024, 8, 24, 17, 30, 0, 0, time.UTC),
			HomeTeamID:  teams["tottenham-hotspur"],
			AwayTeamID:  teams["chelsea"],
			HomeScore:   intPtr(2),
			AwayScore:   intPtr(2),
			Venue:       stringPtr("Tottenham Hotspur Stadium"),
		},
		{
			Competition: "Premier League",
			Season:      "2024-2025",
			Round:       stringPtr("Matchday 4"),
			Status:      model.MatchStatusFinished,
			KickoffAt:   time.Date(2024, 8, 31, 15, 0, 0, 0, time.UTC),
			HomeTeamID:  teams["manchester-city"],
			AwayTeamID:  teams["arsenal"],
			HomeScore:   intPtr(1),
			AwayScore:   intPtr(1),
			Venue:       stringPtr("Etihad Stadium"),
		},
		{
			Competition: "Premier League",
			Season:      "2024-2025",
			Round:       stringPtr("Matchday 4"),
			Status:      model.MatchStatusFinished,
			KickoffAt:   time.Date(2024, 8, 31, 12, 30, 0, 0, time.UTC),
			HomeTeamID:  teams["liverpool"],
			AwayTeamID:  teams["manchester-united"],
			HomeScore:   intPtr(3),
			AwayScore:   intPtr(0),
			Venue:       stringPtr("Anfield"),
		},
		{
			Competition: "Premier League",
			Season:      "2024-2025",
			Round:       stringPtr("Matchday 5"),
			Status:      model.MatchStatusScheduled,
			KickoffAt:   time.Date(2024, 9, 14, 15, 0, 0, 0, time.UTC),
			HomeTeamID:  teams["arsenal"],
			AwayTeamID:  teams["tottenham-hotspur"],
			HomeScore:   nil,
			AwayScore:   nil,
			Venue:       stringPtr("Emirates Stadium"),
		},
	}

	for _, match := range matches {
		// Create match if not exists based on competition, season, home team, away team, and kickoff time
		if err := db.Where(
			"competition = ? AND season = ? AND home_team_id = ? AND away_team_id = ? AND kickoff_at = ?",
			match.Competition, match.Season, match.HomeTeamID, match.AwayTeamID, match.KickoffAt,
		).FirstOrCreate(&match).Error; err != nil {
			return err
		}
	}

	return nil
}

func intPtr(i int) *int {
	return &i
}
