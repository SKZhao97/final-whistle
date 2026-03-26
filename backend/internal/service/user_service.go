package service

import (
	"errors"
	"time"

	"final-whistle/backend/internal/dto"
	"final-whistle/backend/internal/repository"
	"gorm.io/gorm"
)

type UserService interface {
	GetProfileSummary(userID uint) (*dto.UserProfileSummaryDTO, error)
	GetCheckInHistory(userID uint, page, pageSize int) (*dto.UserCheckInHistoryResponseDTO, error)
}

type userService struct {
	repo repository.UserRepository
	now  func() time.Time
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
		now:  time.Now,
	}
}

func (s *userService) GetProfileSummary(userID uint) (*dto.UserProfileSummaryDTO, error) {
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

func (s *userService) GetCheckInHistory(userID uint, page, pageSize int) (*dto.UserCheckInHistoryResponseDTO, error) {
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
