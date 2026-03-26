package repository

import (
	"final-whistle/backend/internal/model"
	"time"

	"gorm.io/gorm"
)

type UserProfileSummaryRecord struct {
	User               model.User
	CheckInCount       int64
	AvgMatchRating     *float64
	FavoriteTeamID     *uint
	FavoriteTeam       *model.Team
	MostUsedTagID      *uint
	MostUsedTag        *model.Tag
	RecentCheckInCount int64
}

type UserCheckInHistoryParams struct {
	Page     int
	PageSize int
}

type UserRepository interface {
	FindUserByID(id uint) (*model.User, error)
	GetUserProfileSummary(userID uint, recentSince time.Time) (*UserProfileSummaryRecord, error)
	GetUserCheckInHistory(userID uint, params UserCheckInHistoryParams) ([]model.CheckIn, int64, error)
}

type GormUserRepository struct {
	*BaseRepository
}

func NewUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{BaseRepository: NewBaseRepository(db)}
}

func (r *GormUserRepository) FindUserByID(id uint) (*model.User, error) {
	var user model.User
	if err := r.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GormUserRepository) GetUserProfileSummary(userID uint, recentSince time.Time) (*UserProfileSummaryRecord, error) {
	user, err := r.FindUserByID(userID)
	if err != nil {
		return nil, err
	}

	record := &UserProfileSummaryRecord{User: *user}

	var aggregate struct {
		CheckInCount   int64
		AvgMatchRating *float64
	}
	if err := r.DB.
		Model(&model.CheckIn{}).
		Select("COUNT(*) AS check_in_count, AVG(match_rating) AS avg_match_rating").
		Where("user_id = ?", userID).
		Scan(&aggregate).Error; err != nil {
		return nil, err
	}
	record.CheckInCount = aggregate.CheckInCount
	record.AvgMatchRating = aggregate.AvgMatchRating

	if err := r.DB.
		Model(&model.CheckIn{}).
		Where("user_id = ? AND watched_at >= ?", userID, recentSince).
		Count(&record.RecentCheckInCount).Error; err != nil {
		return nil, err
	}

	var favoriteTeamRow struct {
		TeamID uint
	}
	if err := r.DB.
		Table("check_ins").
		Select(`
			CASE
				WHEN supporter_side = 'HOME' THEN matches.home_team_id
				WHEN supporter_side = 'AWAY' THEN matches.away_team_id
			END AS team_id`,
		).
		Joins("JOIN matches ON matches.id = check_ins.match_id").
		Where("check_ins.user_id = ? AND check_ins.supporter_side IN ?", userID, []string{"HOME", "AWAY"}).
		Group("team_id").
		Order("COUNT(*) DESC, team_id ASC").
		Limit(1).
		Scan(&favoriteTeamRow).Error; err != nil {
		return nil, err
	}
	if favoriteTeamRow.TeamID != 0 {
		record.FavoriteTeamID = &favoriteTeamRow.TeamID
		var favoriteTeam model.Team
		if err := r.DB.First(&favoriteTeam, favoriteTeamRow.TeamID).Error; err != nil {
			return nil, err
		}
		record.FavoriteTeam = &favoriteTeam
	}

	var mostUsedTagRow struct {
		TagID uint
	}
	if err := r.DB.
		Table("checkin_tags").
		Select("checkin_tags.tag_id").
		Joins("JOIN check_ins ON check_ins.id = checkin_tags.check_in_id").
		Where("check_ins.user_id = ?", userID).
		Group("checkin_tags.tag_id").
		Order("COUNT(*) DESC, checkin_tags.tag_id ASC").
		Limit(1).
		Scan(&mostUsedTagRow).Error; err != nil {
		return nil, err
	}
	if mostUsedTagRow.TagID != 0 {
		record.MostUsedTagID = &mostUsedTagRow.TagID
		var tag model.Tag
		if err := r.DB.First(&tag, mostUsedTagRow.TagID).Error; err != nil {
			return nil, err
		}
		record.MostUsedTag = &tag
	}

	return record, nil
}

func (r *GormUserRepository) GetUserCheckInHistory(userID uint, params UserCheckInHistoryParams) ([]model.CheckIn, int64, error) {
	var total int64
	baseQuery := r.DB.Model(&model.CheckIn{}).Where("user_id = ?", userID)
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var checkIns []model.CheckIn
	if err := r.DB.
		Preload("Tags", func(db *gorm.DB) *gorm.DB {
			return db.Order("tags.sort_order ASC, tags.id ASC")
		}).
		Preload("Match.HomeTeam").
		Preload("Match.AwayTeam").
		Where("user_id = ?", userID).
		Order("watched_at DESC, id DESC").
		Limit(params.PageSize).
		Offset((params.Page - 1) * params.PageSize).
		Find(&checkIns).Error; err != nil {
		return nil, 0, err
	}

	return checkIns, total, nil
}
