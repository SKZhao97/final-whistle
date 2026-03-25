package repository

import (
	"final-whistle/backend/internal/model"
	"gorm.io/gorm"
)

type TeamDetailRatingSummary struct {
	AvgRating   *float64
	RatingCount int64
}

type TeamRepository interface {
	FindByID(id uint) (*model.Team, error)
	ListRecentMatches(teamID uint, limit int) ([]model.Match, error)
	GetRatingSummary(teamID uint) (TeamDetailRatingSummary, error)
}

type GormTeamRepository struct {
	*BaseRepository
}

func NewTeamRepository(db *gorm.DB) *GormTeamRepository {
	return &GormTeamRepository{BaseRepository: NewBaseRepository(db)}
}

func (r *GormTeamRepository) FindByID(id uint) (*model.Team, error) {
	var team model.Team
	if err := r.DB.First(&team, id).Error; err != nil {
		return nil, err
	}
	return &team, nil
}

func (r *GormTeamRepository) ListRecentMatches(teamID uint, limit int) ([]model.Match, error) {
	var matches []model.Match
	err := r.DB.
		Preload("HomeTeam").
		Preload("AwayTeam").
		Where("home_team_id = ? OR away_team_id = ?", teamID, teamID).
		Order("kickoff_at DESC").
		Limit(limit).
		Find(&matches).Error
	return matches, err
}

func (r *GormTeamRepository) GetRatingSummary(teamID uint) (TeamDetailRatingSummary, error) {
	var summary TeamDetailRatingSummary
	err := r.DB.
		Table("check_ins AS ci").
		Select(`ROUND(AVG(CASE
			WHEN m.home_team_id = ? THEN ci.home_team_rating
			WHEN m.away_team_id = ? THEN ci.away_team_rating
		END)::numeric, 1) AS avg_rating,
		COUNT(*) AS rating_count`, teamID, teamID).
		Joins("JOIN matches m ON m.id = ci.match_id").
		Where("m.home_team_id = ? OR m.away_team_id = ?", teamID, teamID).
		Scan(&summary).Error
	return summary, err
}
