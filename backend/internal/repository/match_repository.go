package repository

import (
	"time"

	"final-whistle/backend/internal/model"
	"gorm.io/gorm"
)

type MatchListParams struct {
	Competition string
	Season      string
	Page        int
	PageSize    int
}

type MatchAggregateRecord struct {
	MatchID           uint
	MatchRatingAvg    *float64
	HomeTeamRatingAvg *float64
	AwayTeamRatingAvg *float64
	CheckInCount      int64
}

type MatchPlayerRatingRecord struct {
	PlayerID      uint
	PlayerName    string
	PlayerSlug    string
	Position      *string
	AvatarURL     *string
	TeamID        uint
	TeamName      string
	TeamShortName *string
	TeamSlug      string
	TeamLogoURL   *string
	AvgRating     *float64
	RatingCount   int64
}

type MatchRecentReviewRecord struct {
	CheckInID     uint
	UserID        uint
	UserName      string
	UserAvatarURL *string
	MatchRating   int
	ShortReview   string
	CreatedAt     time.Time
	Tags          []model.Tag
}

type MatchRepository interface {
	ListMatches(params MatchListParams) ([]model.Match, int64, error)
	GetMatchAggregates(matchIDs []uint) (map[uint]MatchAggregateRecord, error)
	FindMatchByID(id uint) (*model.Match, error)
	GetPlayerRatingSummary(matchID uint, limit int) ([]MatchPlayerRatingRecord, error)
	GetRecentReviews(matchID uint, limit int) ([]MatchRecentReviewRecord, error)
}

type GormMatchRepository struct {
	*BaseRepository
}

func NewMatchRepository(db *gorm.DB) *GormMatchRepository {
	return &GormMatchRepository{BaseRepository: NewBaseRepository(db)}
}

func (r *GormMatchRepository) ListMatches(params MatchListParams) ([]model.Match, int64, error) {
	var matches []model.Match
	var total int64

	query := r.DB.Model(&model.Match{})
	if params.Competition != "" {
		query = query.Where("competition = ?", params.Competition)
	}
	if params.Season != "" {
		query = query.Where("season = ?", params.Season)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.PageSize
	err := query.
		Preload("HomeTeam").
		Preload("AwayTeam").
		Order("kickoff_at DESC").
		Limit(params.PageSize).
		Offset(offset).
		Find(&matches).Error
	if err != nil {
		return nil, 0, err
	}

	return matches, total, nil
}

func (r *GormMatchRepository) GetMatchAggregates(matchIDs []uint) (map[uint]MatchAggregateRecord, error) {
	result := make(map[uint]MatchAggregateRecord, len(matchIDs))
	if len(matchIDs) == 0 {
		return result, nil
	}

	var rows []MatchAggregateRecord
	if err := r.DB.
		Table("check_ins").
		Select("match_id, ROUND(AVG(match_rating)::numeric, 1) AS match_rating_avg, ROUND(AVG(home_team_rating)::numeric, 1) AS home_team_rating_avg, ROUND(AVG(away_team_rating)::numeric, 1) AS away_team_rating_avg, COUNT(*) AS check_in_count").
		Where("match_id IN ?", matchIDs).
		Group("match_id").
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	for _, row := range rows {
		result[row.MatchID] = row
	}

	for _, id := range matchIDs {
		if _, ok := result[id]; !ok {
			result[id] = MatchAggregateRecord{MatchID: id, CheckInCount: 0}
		}
	}

	return result, nil
}

func (r *GormMatchRepository) FindMatchByID(id uint) (*model.Match, error) {
	var match model.Match
	if err := r.DB.Preload("HomeTeam").Preload("AwayTeam").First(&match, id).Error; err != nil {
		return nil, err
	}
	return &match, nil
}

func (r *GormMatchRepository) GetPlayerRatingSummary(matchID uint, limit int) ([]MatchPlayerRatingRecord, error) {
	var rows []MatchPlayerRatingRecord
	err := r.DB.
		Table("player_ratings AS pr").
		Select(`pr.player_id,
			p.name AS player_name,
			p.slug AS player_slug,
			p.position,
			p.avatar_url,
			t.id AS team_id,
			t.name AS team_name,
			t.short_name AS team_short_name,
			t.slug AS team_slug,
			t.logo_url AS team_logo_url,
			ROUND(AVG(pr.rating)::numeric, 1) AS avg_rating,
			COUNT(*) AS rating_count`).
		Joins("JOIN check_ins ci ON ci.id = pr.check_in_id").
		Joins("JOIN players p ON p.id = pr.player_id").
		Joins("JOIN teams t ON t.id = p.team_id").
		Where("ci.match_id = ?", matchID).
		Group("pr.player_id, p.id, t.id").
		Order("avg_rating DESC NULLS LAST, rating_count DESC, p.name ASC").
		Limit(limit).
		Scan(&rows).Error
	return rows, err
}

func (r *GormMatchRepository) GetRecentReviews(matchID uint, limit int) ([]MatchRecentReviewRecord, error) {
	type reviewRow struct {
		CheckInID     uint
		UserID        uint
		UserName      string
		UserAvatarURL *string
		MatchRating   int
		ShortReview   string
		CreatedAt     time.Time
	}

	var rows []reviewRow
	if err := r.DB.
		Table("check_ins AS ci").
		Select("ci.id AS check_in_id, u.id AS user_id, u.name AS user_name, u.avatar_url AS user_avatar_url, ci.match_rating, ci.short_review, ci.created_at").
		Joins("JOIN users u ON u.id = ci.user_id").
		Where("ci.match_id = ? AND ci.short_review IS NOT NULL AND ci.short_review <> ''", matchID).
		Order("ci.created_at DESC").
		Limit(limit).
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return []MatchRecentReviewRecord{}, nil
	}

	checkInIDs := make([]uint, 0, len(rows))
	for _, row := range rows {
		checkInIDs = append(checkInIDs, row.CheckInID)
	}

	var checkIns []model.CheckIn
	if err := r.DB.
		Preload("Tags").
		Where("id IN ?", checkInIDs).
		Find(&checkIns).Error; err != nil {
		return nil, err
	}

	tagsByCheckInID := make(map[uint][]model.Tag, len(checkIns))
	for _, checkIn := range checkIns {
		tagsByCheckInID[checkIn.ID] = checkIn.Tags
	}

	result := make([]MatchRecentReviewRecord, 0, len(rows))
	for _, row := range rows {
		result = append(result, MatchRecentReviewRecord{
			CheckInID:     row.CheckInID,
			UserID:        row.UserID,
			UserName:      row.UserName,
			UserAvatarURL: row.UserAvatarURL,
			MatchRating:   row.MatchRating,
			ShortReview:   row.ShortReview,
			CreatedAt:     row.CreatedAt,
			Tags:          tagsByCheckInID[row.CheckInID],
		})
	}

	return result, nil
}
