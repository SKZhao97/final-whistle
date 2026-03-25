package service

import (
	"errors"

	"final-whistle/backend/internal/dto"
	"final-whistle/backend/internal/repository"
	"gorm.io/gorm"
)

type PlayerService interface {
	GetPlayerDetail(id uint) (*dto.PlayerDetailDTO, error)
}

type playerService struct {
	repo repository.PlayerRepository
}

func NewPlayerService(repo repository.PlayerRepository) PlayerService {
	return &playerService{repo: repo}
}

func (s *playerService) GetPlayerDetail(id uint) (*dto.PlayerDetailDTO, error) {
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
				Competition: match.Match.Competition,
				Season:      match.Match.Season,
				Round:       match.Match.Round,
				Status:      string(match.Match.Status),
				KickoffAt:   match.Match.KickoffAt,
				HomeTeam:    toTeamSummaryDTO(match.Match.HomeTeam),
				AwayTeam:    toTeamSummaryDTO(match.Match.AwayTeam),
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
		Team:          toTeamSummaryDTO(player.Team),
		RecentMatches: items,
		RatingSummary: dto.PlayerRatingSummaryDTO{
			AvgRating:   ratingSummary.AvgRating,
			RatingCount: ratingSummary.RatingCount,
		},
	}, nil
}
