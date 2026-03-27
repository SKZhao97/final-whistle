package service

import (
	"errors"

	"final-whistle/backend/internal/dto"
	"final-whistle/backend/internal/repository"
	"gorm.io/gorm"
)

type TeamService interface {
	GetTeamDetail(id uint, locale string) (*dto.TeamDetailDTO, error)
}

type teamService struct {
	repo      repository.TeamRepository
	matchRepo repository.MatchRepository
}

func NewTeamService(repo repository.TeamRepository, matchRepo repository.MatchRepository) TeamService {
	return &teamService{repo: repo, matchRepo: matchRepo}
}

func (s *teamService) GetTeamDetail(id uint, locale string) (*dto.TeamDetailDTO, error) {
	team, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	recentMatches, err := s.repo.ListRecentMatches(id, 5)
	if err != nil {
		return nil, err
	}
	matchIDs := make([]uint, 0, len(recentMatches))
	for _, match := range recentMatches {
		matchIDs = append(matchIDs, match.ID)
	}
	aggregates, err := s.matchRepo.GetMatchAggregates(matchIDs)
	if err != nil {
		return nil, err
	}
	ratingSummary, err := s.repo.GetRatingSummary(id)
	if err != nil {
		return nil, err
	}

	items := make([]dto.MatchListItemDTO, 0, len(recentMatches))
	for _, match := range recentMatches {
		items = append(items, dto.MatchListItemDTO{
			ID:          match.ID,
			Competition: localizedCompetitionName(match.Competition, locale),
			Season:      match.Season,
			Round:       localizedRound(match.Round, locale),
			Status:      string(match.Status),
			KickoffAt:   match.KickoffAt,
			HomeTeam:    toTeamSummaryDTO(match.HomeTeam, locale),
			AwayTeam:    toTeamSummaryDTO(match.AwayTeam, locale),
			HomeScore:   match.HomeScore,
			AwayScore:   match.AwayScore,
			Aggregates:  toAggregateDTO(aggregates[match.ID]),
		})
	}

	return &dto.TeamDetailDTO{
		ID:            team.ID,
		Name:          localizedTeamName(*team, locale),
		ShortName:     team.ShortName,
		Slug:          team.Slug,
		LogoURL:       team.LogoURL,
		RecentMatches: items,
		RatingSummary: dto.TeamRatingSummaryDTO{
			AvgRating:   ratingSummary.AvgRating,
			RatingCount: ratingSummary.RatingCount,
		},
	}, nil
}
