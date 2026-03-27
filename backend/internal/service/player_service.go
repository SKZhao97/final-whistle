package service

import (
	"errors"

	"final-whistle/backend/internal/dto"
	"final-whistle/backend/internal/repository"
	"gorm.io/gorm"
)

type PlayerService interface {
	GetPlayerDetail(id uint, locale string) (*dto.PlayerDetailDTO, error)
}

type playerService struct {
	repo repository.PlayerRepository
}

func NewPlayerService(repo repository.PlayerRepository) PlayerService {
	return &playerService{repo: repo}
}

func (s *playerService) GetPlayerDetail(id uint, locale string) (*dto.PlayerDetailDTO, error) {
	player, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	recentMatches, err := s.repo.ListRecentRatedMatches(id, 5)
	if err != nil {
		return nil, err
	}
	ratingSummary, err := s.repo.GetRatingSummary(id)
	if err != nil {
		return nil, err
	}

	items := make([]dto.PlayerRecentMatchDTO, 0, len(recentMatches))
	for _, match := range recentMatches {
		items = append(items, dto.PlayerRecentMatchDTO{
			Match: dto.MatchListItemDTO{
				ID:          match.Match.ID,
				Competition: localizedCompetitionName(match.Match.Competition, locale),
				Season:      match.Match.Season,
				Round:       localizedRound(match.Match.Round, locale),
				Status:      string(match.Match.Status),
				KickoffAt:   match.Match.KickoffAt,
				HomeTeam:    toTeamSummaryDTO(match.Match.HomeTeam, locale),
				AwayTeam:    toTeamSummaryDTO(match.Match.AwayTeam, locale),
				HomeScore:   match.Match.HomeScore,
				AwayScore:   match.Match.AwayScore,
			},
			AvgRating:   match.AvgRating,
			RatingCount: match.RatingCount,
		})
	}

	return &dto.PlayerDetailDTO{
		ID:            player.ID,
		Name:          player.Name,
		Slug:          player.Slug,
		Position:      player.Position,
		AvatarURL:     player.AvatarURL,
		Team:          toTeamSummaryDTO(player.Team, locale),
		RecentMatches: items,
		RatingSummary: dto.PlayerRatingSummaryDTO{
			AvgRating:   ratingSummary.AvgRating,
			RatingCount: ratingSummary.RatingCount,
		},
	}, nil
}
