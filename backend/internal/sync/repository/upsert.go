package repository

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"final-whistle/backend/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const ProviderFootballData = "football-data"

type UpsertTeamParams struct {
	Name              string
	ShortName         *string
	Slug              string
	LogoURL           *string
	ExternalSource    string
	ExternalID        string
	ExternalUpdatedAt *time.Time
}

type UpsertPlayerParams struct {
	TeamID            uint
	Name              string
	Slug              string
	Position          *string
	AvatarURL         *string
	ExternalSource    string
	ExternalID        string
	ExternalUpdatedAt *time.Time
}

type UpsertMatchParams struct {
	Competition       string
	Season            string
	Round             *string
	Status            model.MatchStatus
	KickoffAt         time.Time
	HomeTeamID        uint
	AwayTeamID        uint
	HomeScore         *int
	AwayScore         *int
	Venue             *string
	ExternalSource    string
	ExternalID        string
	ExternalUpdatedAt *time.Time
}

var slugSanitizer = regexp.MustCompile(`[^a-z0-9]+`)

func NormalizeSlug(value string) string {
	lowered := strings.ToLower(strings.TrimSpace(value))
	lowered = slugSanitizer.ReplaceAllString(lowered, "-")
	lowered = strings.Trim(lowered, "-")
	if lowered == "" {
		return "item"
	}
	return lowered
}

func (r *Repository) UpsertTeam(ctx context.Context, params UpsertTeamParams) (*model.Team, error) {
	now := time.Now().UTC()
	team := model.Team{
		Name:              params.Name,
		ShortName:         params.ShortName,
		Slug:              params.Slug,
		LogoURL:           params.LogoURL,
		ExternalSource:    stringPtr(params.ExternalSource),
		ExternalID:        stringPtr(params.ExternalID),
		ExternalUpdatedAt: params.ExternalUpdatedAt,
	}
	if err := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "external_source"}, {Name: "external_id"}},
		DoUpdates: clause.Assignments(map[string]any{
			"name":                params.Name,
			"short_name":          params.ShortName,
			"slug":                params.Slug,
			"logo_url":            params.LogoURL,
			"external_updated_at": params.ExternalUpdatedAt,
			"updated_at":          now,
		}),
	}).Create(&team).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).
		Where("external_source = ? AND external_id = ?", params.ExternalSource, params.ExternalID).
		First(&team).Error; err != nil {
		return nil, err
	}
	return &team, nil
}

func (r *Repository) UpsertPlayer(ctx context.Context, params UpsertPlayerParams) (*model.Player, error) {
	now := time.Now().UTC()
	player := model.Player{
		TeamID:            params.TeamID,
		Name:              params.Name,
		Slug:              params.Slug,
		Position:          params.Position,
		AvatarURL:         params.AvatarURL,
		ExternalSource:    stringPtr(params.ExternalSource),
		ExternalID:        stringPtr(params.ExternalID),
		ExternalUpdatedAt: params.ExternalUpdatedAt,
	}
	if err := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "external_source"}, {Name: "external_id"}},
		DoUpdates: clause.Assignments(map[string]any{
			"team_id":             params.TeamID,
			"name":                params.Name,
			"slug":                params.Slug,
			"position":            params.Position,
			"avatar_url":          params.AvatarURL,
			"external_updated_at": params.ExternalUpdatedAt,
			"updated_at":          now,
		}),
	}).Create(&player).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).
		Where("external_source = ? AND external_id = ?", params.ExternalSource, params.ExternalID).
		First(&player).Error; err != nil {
		return nil, err
	}
	return &player, nil
}

func (r *Repository) UpsertMatch(ctx context.Context, params UpsertMatchParams) (*model.Match, error) {
	now := time.Now().UTC()
	match := model.Match{
		Competition:       params.Competition,
		Season:            params.Season,
		Round:             params.Round,
		Status:            params.Status,
		KickoffAt:         params.KickoffAt,
		HomeTeamID:        params.HomeTeamID,
		AwayTeamID:        params.AwayTeamID,
		HomeScore:         params.HomeScore,
		AwayScore:         params.AwayScore,
		Venue:             params.Venue,
		ExternalSource:    stringPtr(params.ExternalSource),
		ExternalID:        stringPtr(params.ExternalID),
		ExternalUpdatedAt: params.ExternalUpdatedAt,
	}
	if err := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "external_source"}, {Name: "external_id"}},
		DoUpdates: clause.Assignments(map[string]any{
			"competition":         params.Competition,
			"season":              params.Season,
			"round":               params.Round,
			"status":              params.Status,
			"kickoff_at":          params.KickoffAt,
			"home_team_id":        params.HomeTeamID,
			"away_team_id":        params.AwayTeamID,
			"home_score":          params.HomeScore,
			"away_score":          params.AwayScore,
			"venue":               params.Venue,
			"external_updated_at": params.ExternalUpdatedAt,
			"updated_at":          now,
		}),
	}).Create(&match).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).
		Where("external_source = ? AND external_id = ?", params.ExternalSource, params.ExternalID).
		First(&match).Error; err != nil {
		return nil, err
	}
	return &match, nil
}

func (r *Repository) FindTeamByExternalID(ctx context.Context, source, externalID string) (*model.Team, error) {
	var team model.Team
	err := r.db.WithContext(ctx).Where("external_source = ? AND external_id = ?", source, externalID).First(&team).Error
	if err != nil {
		return nil, err
	}
	return &team, nil
}

func (r *Repository) FindPlayerByExternalID(ctx context.Context, source, externalID string) (*model.Player, error) {
	var player model.Player
	err := r.db.WithContext(ctx).Where("external_source = ? AND external_id = ?", source, externalID).First(&player).Error
	if err != nil {
		return nil, err
	}
	return &player, nil
}

func (r *Repository) ListTeamsByExternalSource(ctx context.Context, source string) ([]model.Team, error) {
	var teams []model.Team
	err := r.db.WithContext(ctx).Where("external_source = ?", source).Order("id ASC").Find(&teams).Error
	return teams, err
}

func (r *Repository) GetMatchByID(ctx context.Context, id uint) (*model.Match, error) {
	var match model.Match
	if err := r.db.WithContext(ctx).First(&match, id).Error; err != nil {
		return nil, err
	}
	return &match, nil
}

func (r *Repository) ReplaceMatchRoster(ctx context.Context, matchID uint, entries []model.MatchPlayer) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("match_id = ?", matchID).Delete(&model.MatchPlayer{}).Error; err != nil {
			return err
		}
		if len(entries) == 0 {
			return nil
		}
		return tx.Create(&entries).Error
	})
}

func MapProviderMatchStatus(status string) model.MatchStatus {
	switch strings.ToUpper(strings.TrimSpace(status)) {
	case "FINISHED", "AWARDED":
		return model.MatchStatusFinished
	default:
		return model.MatchStatusScheduled
	}
}

func SeasonFromKickoff(kickoff time.Time) string {
	year := kickoff.UTC().Year()
	return fmt.Sprintf("%d", year)
}

func stringPtr(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}
