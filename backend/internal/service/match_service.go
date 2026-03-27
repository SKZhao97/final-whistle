package service

import (
	"errors"

	"final-whistle/backend/internal/dto"
	"final-whistle/backend/internal/repository"
	"gorm.io/gorm"
)

type MatchService interface {
	ListMatches(params repository.MatchListParams) (dto.MatchListResponseDTO, error)
	GetMatchDetail(id uint, locale string) (*dto.MatchDetailDTO, error)
}

type matchService struct {
	repo repository.MatchRepository
}

func NewMatchService(repo repository.MatchRepository) MatchService {
	return &matchService{repo: repo}
}

func (s *matchService) ListMatches(params repository.MatchListParams) (dto.MatchListResponseDTO, error) {
	matches, total, err := s.repo.ListMatches(params)
	if err != nil {
		return dto.MatchListResponseDTO{}, err
	}

	matchIDs := make([]uint, 0, len(matches))
	for _, match := range matches {
		matchIDs = append(matchIDs, match.ID)
	}
	aggregates, err := s.repo.GetMatchAggregates(matchIDs)
	if err != nil {
		return dto.MatchListResponseDTO{}, err
	}

	items := make([]dto.MatchListItemDTO, 0, len(matches))
	for _, match := range matches {
		items = append(items, dto.MatchListItemDTO{
			ID:          match.ID,
			Competition: match.Competition,
			Season:      match.Season,
			Round:       match.Round,
			Status:      string(match.Status),
			KickoffAt:   match.KickoffAt,
			HomeTeam:    toTeamSummaryDTO(match.HomeTeam),
			AwayTeam:    toTeamSummaryDTO(match.AwayTeam),
			HomeScore:   match.HomeScore,
			AwayScore:   match.AwayScore,
			Aggregates:  toAggregateDTO(aggregates[match.ID]),
		})
	}

	return dto.MatchListResponseDTO{
		Items:    items,
		Page:     params.Page,
		PageSize: params.PageSize,
		Total:    total,
	}, nil
}

func (s *matchService) GetMatchDetail(id uint, locale string) (*dto.MatchDetailDTO, error) {
	match, err := s.repo.FindMatchByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	aggregates, err := s.repo.GetMatchAggregates([]uint{id})
	if err != nil {
		return nil, err
	}
	playerRatings, err := s.repo.GetPlayerRatingSummary(id, 10)
	if err != nil {
		return nil, err
	}
	rosterPlayers, err := s.repo.GetMatchRoster(id)
	if err != nil {
		return nil, err
	}
	availableTags, err := s.repo.ListActiveTags()
	if err != nil {
		return nil, err
	}
	reviews, err := s.repo.GetRecentReviews(id, 10)
	if err != nil {
		return nil, err
	}

	tagItems := make([]dto.TagDTO, 0, len(availableTags))
	for _, tag := range availableTags {
		tagItems = append(tagItems, toTagDTO(tag, locale))
	}

	rosterItems := make([]dto.PlayerSummaryDTO, 0, len(rosterPlayers))
	for _, item := range rosterPlayers {
		rosterItems = append(rosterItems, dto.PlayerSummaryDTO{
			ID:        item.PlayerID,
			Name:      item.PlayerName,
			Slug:      item.PlayerSlug,
			Position:  item.Position,
			AvatarURL: item.AvatarURL,
			Team: dto.TeamSummaryDTO{
				ID:        item.TeamID,
				Name:      item.TeamName,
				ShortName: item.TeamShortName,
				Slug:      item.TeamSlug,
				LogoURL:   item.TeamLogoURL,
			},
		})
	}

	playerItems := make([]dto.MatchPlayerRatingSummaryDTO, 0, len(playerRatings))
	for _, item := range playerRatings {
		playerItems = append(playerItems, dto.MatchPlayerRatingSummaryDTO{
			Player: dto.PlayerSummaryDTO{
				ID:        item.PlayerID,
				Name:      item.PlayerName,
				Slug:      item.PlayerSlug,
				Position:  item.Position,
				AvatarURL: item.AvatarURL,
				Team: dto.TeamSummaryDTO{
					ID:        item.TeamID,
					Name:      item.TeamName,
					ShortName: item.TeamShortName,
					Slug:      item.TeamSlug,
					LogoURL:   item.TeamLogoURL,
				},
			},
			AvgRating:   item.AvgRating,
			RatingCount: item.RatingCount,
		})
	}

	reviewItems := make([]dto.MatchRecentReviewDTO, 0, len(reviews))
	for _, review := range reviews {
		tags := make([]dto.TagDTO, 0, len(review.Tags))
		for _, tag := range review.Tags {
			tags = append(tags, toTagDTO(tag, locale))
		}
		reviewItems = append(reviewItems, dto.MatchRecentReviewDTO{
			ID: review.CheckInID,
			User: dto.UserSummaryDTO{
				ID:        review.UserID,
				Name:      review.UserName,
				AvatarURL: review.UserAvatarURL,
			},
			MatchRating: review.MatchRating,
			ShortReview: review.ShortReview,
			Tags:        tags,
			CreatedAt:   review.CreatedAt,
		})
	}

	return &dto.MatchDetailDTO{
		ID:            match.ID,
		Competition:   match.Competition,
		Season:        match.Season,
		Round:         match.Round,
		Status:        string(match.Status),
		KickoffAt:     match.KickoffAt,
		HomeTeam:      toTeamSummaryDTO(match.HomeTeam),
		AwayTeam:      toTeamSummaryDTO(match.AwayTeam),
		HomeScore:     match.HomeScore,
		AwayScore:     match.AwayScore,
		Venue:         match.Venue,
		Aggregates:    toAggregateDTO(aggregates[id]),
		AvailableTags: tagItems,
		MatchPlayers:  rosterItems,
		PlayerRatings: playerItems,
		RecentReviews: reviewItems,
	}, nil
}

func toAggregateDTO(record repository.MatchAggregateRecord) dto.MatchAggregateSummaryDTO {
	return dto.MatchAggregateSummaryDTO{
		MatchRatingAvg:    record.MatchRatingAvg,
		HomeTeamRatingAvg: record.HomeTeamRatingAvg,
		AwayTeamRatingAvg: record.AwayTeamRatingAvg,
		CheckInCount:      record.CheckInCount,
	}
}
