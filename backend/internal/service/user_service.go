// Package service 包含业务逻辑层，负责处理核心业务规则、数据验证和事务管理。
// 本包包含用户相关的业务服务。
package service

import (
	"errors"
	"time"

	"final-whistle/backend/internal/dto"
	"final-whistle/backend/internal/repository"
	"gorm.io/gorm"
)

// UserService 定义了用户相关的业务操作接口。
type UserService interface {
	// GetProfileSummary 获取用户资料摘要。
	GetProfileSummary(userID uint) (*dto.UserProfileSummaryDTO, error)
	// GetCheckInHistory 获取用户签到历史。
	GetCheckInHistory(userID uint, page, pageSize int) (*dto.UserCheckInHistoryResponseDTO, error)
}

// userService 是 UserService 接口的实现。
type userService struct {
	repo repository.UserRepository
	now  func() time.Time
}

// NewUserService 创建并返回一个新的 UserService 实例。
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
		now:  time.Now,
	}
}

// GetProfileSummary 获取用户资料摘要，包括签到统计、平均评分、最爱球队等信息。
// 参数:
//   - userID: 用户ID
// 返回:
//   - UserProfileSummaryDTO: 用户资料摘要
//   - error: 错误信息，如用户不存在或数据库错误
func (s *userService) GetProfileSummary(userID uint) (*dto.UserProfileSummaryDTO, error) {
	// 获取用户资料摘要记录，recentSince为30天前的时间点
	record, err := s.repo.GetUserProfileSummary(userID, s.now().AddDate(0, 0, -30))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	result := &dto.UserProfileSummaryDTO{
		User: dto.UserSummaryDTO{
			ID:        record.User.ID,
			Name:      record.User.Name,
			AvatarURL: record.User.AvatarURL,
		},
		CheckInCount:       int(record.CheckInCount),
		AvgMatchRating:     record.AvgMatchRating,
		FavoriteTeamID:     record.FavoriteTeamID,
		MostUsedTagID:      record.MostUsedTagID,
		RecentCheckInCount: int(record.RecentCheckInCount),
	}

	if record.FavoriteTeam != nil {
		result.FavoriteTeam = &dto.TeamSummaryDTO{
			ID:        record.FavoriteTeam.ID,
			Name:      record.FavoriteTeam.Name,
			ShortName: record.FavoriteTeam.ShortName,
			Slug:      record.FavoriteTeam.Slug,
			LogoURL:   record.FavoriteTeam.LogoURL,
		}
	}
	if record.MostUsedTag != nil {
		result.MostUsedTag = &dto.TagDTO{
			ID:   record.MostUsedTag.ID,
			Name: record.MostUsedTag.Name,
			Slug: record.MostUsedTag.Slug,
		}
	}

	return result, nil
}

// GetCheckInHistory 获取用户签到历史，支持分页。
// 参数:
//   - userID: 用户ID
//   - page: 页码，小于1时默认1
//   - pageSize: 每页大小，小于1时默认20，大于50时限制为50
// 返回:
//   - UserCheckInHistoryResponseDTO: 分页签到历史响应
//   - error: 错误信息，如数据库错误
func (s *userService) GetCheckInHistory(userID uint, page, pageSize int) (*dto.UserCheckInHistoryResponseDTO, error) {
	// 验证和规范化分页参数
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 50 {
		pageSize = 50
	}

	checkIns, total, err := s.repo.GetUserCheckInHistory(userID, repository.UserCheckInHistoryParams{
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		return nil, err
	}

	items := make([]dto.UserCheckInHistoryItemDTO, 0, len(checkIns))
	for _, checkIn := range checkIns {
		tags := make([]dto.TagDTO, 0, len(checkIn.Tags))
		for _, tag := range checkIn.Tags {
			tags = append(tags, dto.TagDTO{
				ID:   tag.ID,
				Name: tag.Name,
				Slug: tag.Slug,
			})
		}

		items = append(items, dto.UserCheckInHistoryItemDTO{
			ID:             checkIn.ID,
			MatchID:        checkIn.MatchID,
			WatchedType:    string(checkIn.WatchedType),
			SupporterSide:  string(checkIn.SupporterSide),
			MatchRating:    checkIn.MatchRating,
			HomeTeamRating: checkIn.HomeTeamRating,
			AwayTeamRating: checkIn.AwayTeamRating,
			ShortReview:    checkIn.ShortReview,
			WatchedAt:      checkIn.WatchedAt,
			CreatedAt:      checkIn.CreatedAt,
			UpdatedAt:      checkIn.UpdatedAt,
			Tags:           tags,
			Match: dto.MatchListItemDTO{
				ID:          checkIn.Match.ID,
				Competition: checkIn.Match.Competition,
				Season:      checkIn.Match.Season,
				Round:       checkIn.Match.Round,
				Status:      string(checkIn.Match.Status),
				KickoffAt:   checkIn.Match.KickoffAt,
				HomeTeam:    toTeamSummaryDTO(checkIn.Match.HomeTeam),
				AwayTeam:    toTeamSummaryDTO(checkIn.Match.AwayTeam),
				HomeScore:   checkIn.Match.HomeScore,
				AwayScore:   checkIn.Match.AwayScore,
				Aggregates:  dto.MatchAggregateSummaryDTO{},
			},
		})
	}

	return &dto.UserCheckInHistoryResponseDTO{
		Items:    items,
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}, nil
}
