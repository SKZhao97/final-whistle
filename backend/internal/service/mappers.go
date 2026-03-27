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

func toTagDTO(tag model.Tag, locale string) dto.TagDTO {
	return dto.TagDTO{
		ID:   tag.ID,
		Name: localizedTagName(tag, locale),
		Slug: tag.Slug,
	}
}

func localizedTagName(tag model.Tag, locale string) string {
	if locale == "zh" && tag.NameZh != "" {
		return tag.NameZh
	}
	if tag.NameEn != "" {
		return tag.NameEn
	}
	return tag.Name
}
