package service

import (
	"final-whistle/backend/internal/dto"
	"final-whistle/backend/internal/model"
)

func toTeamSummaryDTO(team model.Team) dto.TeamSummaryDTO {
	return dto.TeamSummaryDTO{
		ID:        team.ID,
		Name:      team.Name,
		ShortName: team.ShortName,
		Slug:      team.Slug,
		LogoURL:   team.LogoURL,
	}
}
