package repository

import (
	"time"

	"final-whistle/backend/internal/model"
	"gorm.io/gorm"
)

type PlayerDetailRatingSummary struct {
	AvgRating   *float64
	RatingCount int64
}

type PlayerRecentMatchRecord struct {
	Match       model.Match
	AvgRating   *float64
	RatingCount int64
}

type PlayerRepository interface {
	FindByID(id uint) (*model.Player, error)
	ListRecentRatedMatches(playerID uint, limit int) ([]PlayerRecentMatchRecord, error)
	GetRatingSummary(playerID uint) (PlayerDetailRatingSummary, error)
}

type GormPlayerRepository struct {
	*BaseRepository
}

func NewPlayerRepository(db *gorm.DB) *GormPlayerRepository {
	return &GormPlayerRepository{BaseRepository: NewBaseRepository(db)}
}

func (r *GormPlayerRepository) FindByID(id uint) (*model.Player, error) {
	var player model.Player
	if err := r.DB.Preload("Team").First(&player, id).Error; err != nil {
		return nil, err
	}
	return &player, nil
}

func (r *GormPlayerRepository) ListRecentRatedMatches(playerID uint, limit int) ([]PlayerRecentMatchRecord, error) {
	type row struct {
		MatchID     uint
		AvgRating   *float64
		RatingCount int64
		LastRatedAt time.Time
	}

	var rows []row
	if err := r.DB.
		Table("player_ratings AS pr").
		Select("ci.match_id, ROUND(AVG(pr.rating)::numeric, 1) AS avg_rating, COUNT(*) AS rating_count, MAX(ci.created_at) AS last_rated_at").
		Joins("JOIN check_ins ci ON ci.id = pr.check_in_id").
		Where("pr.player_id = ?", playerID).
		Group("ci.match_id").
		Order("last_rated_at DESC").
		Limit(limit).
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	result := make([]PlayerRecentMatchRecord, 0, len(rows))
	for _, row := range rows {
		var match model.Match
		if err := r.DB.Preload("HomeTeam").Preload("AwayTeam").First(&match, row.MatchID).Error; err != nil {
			return nil, err
		}
		result = append(result, PlayerRecentMatchRecord{
			Match:       match,
			AvgRating:   row.AvgRating,
			RatingCount: row.RatingCount,
		})
	}
	return result, nil
}

func (r *GormPlayerRepository) GetRatingSummary(playerID uint) (PlayerDetailRatingSummary, error) {
	var summary PlayerDetailRatingSummary
	err := r.DB.
		Table("player_ratings").
		Select("ROUND(AVG(rating)::numeric, 1) AS avg_rating, COUNT(*) AS rating_count").
		Where("player_id = ?", playerID).
		Scan(&summary).Error
	return summary, err
}
