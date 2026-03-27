package service

import (
	"final-whistle/backend/internal/dto"
	"final-whistle/backend/internal/model"
	"fmt"
	"strconv"
	"strings"
)

func toTeamSummaryDTO(team model.Team, locale string) dto.TeamSummaryDTO {
	return dto.TeamSummaryDTO{
		ID:        team.ID,
		Name:      localizedTeamName(team, locale),
		ShortName: team.ShortName,
		Slug:      team.Slug,
		LogoURL:   team.LogoURL,
	}
}

func localizedTeamName(team model.Team, locale string) string {
	return localizedTeamNameValue(team.Name, team.NameZh, locale)
}

func localizedTeamNameValue(name string, nameZh *string, locale string) string {
	if locale == "zh" && nameZh != nil && *nameZh != "" {
		return *nameZh
	}
	return name
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

func localizedCompetitionName(raw string, locale string) string {
	if locale != "zh" {
		return raw
	}

	switch raw {
	case "Premier League":
		return "英超"
	default:
		return raw
	}
}

func localizedRound(round *string, locale string) *string {
	if round == nil {
		return nil
	}
	if locale != "zh" {
		return round
	}

	raw := strings.TrimSpace(*round)
	if strings.HasPrefix(raw, "Matchday ") {
		number := strings.TrimPrefix(raw, "Matchday ")
		if _, err := strconv.Atoi(number); err == nil {
			localized := fmt.Sprintf("第%s轮", number)
			return &localized
		}
	}

	return round
}
