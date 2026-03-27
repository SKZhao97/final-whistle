// Package repository 提供数据访问层，负责与数据库交互，执行查询和事务操作。
// 本包包含用户相关的数据访问接口和实现。
package repository

import (
	"final-whistle/backend/internal/model"
	"time"

	"gorm.io/gorm"
)

// UserProfileSummaryRecord 用于存储用户资料摘要的聚合查询结果。
// 这个结构体是内部使用的，不直接暴露给API。
type UserProfileSummaryRecord struct {
	User               model.User   // 用户基本信息
	CheckInCount       int64        // 签到总数
	AvgMatchRating     *float64     // 平均比赛评分，可能为nil
	FavoriteTeamID     *uint        // 最爱球队ID，可能为nil
	FavoriteTeam       *model.Team  // 最爱球队详情，可能为nil
	MostUsedTagID      *uint        // 最常用标签ID，可能为nil
	MostUsedTag        *model.Tag   // 最常用标签详情，可能为nil
	RecentCheckInCount int64        // 最近30天签到数
}

// UserCheckInHistoryParams 包含用户签到历史查询的分页参数。
type UserCheckInHistoryParams struct {
	Page     int // 页码，从1开始
	PageSize int // 每页大小
}

// UserRepository 定义了用户数据访问接口。
type UserRepository interface {
	// FindUserByID 根据ID查找用户。
	FindUserByID(id uint) (*model.User, error)
	// GetUserProfileSummary 获取用户资料摘要，包括聚合统计信息。
	GetUserProfileSummary(userID uint, recentSince time.Time) (*UserProfileSummaryRecord, error)
	// GetUserCheckInHistory 获取用户签到历史，支持分页。
	GetUserCheckInHistory(userID uint, params UserCheckInHistoryParams) ([]model.CheckIn, int64, error)
}

// GormUserRepository 是 UserRepository 接口的 GORM 实现。
type GormUserRepository struct {
	*BaseRepository
}

// NewUserRepository 创建并返回一个新的 GormUserRepository 实例。
func NewUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{BaseRepository: NewBaseRepository(db)}
}

// FindUserByID 根据用户ID查找用户。
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
