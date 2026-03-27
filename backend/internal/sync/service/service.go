package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"final-whistle/backend/internal/config"
	"final-whistle/backend/internal/model"
	syncprovider "final-whistle/backend/internal/sync/provider"
	syncrepo "final-whistle/backend/internal/sync/repository"
)

type Executor interface {
	Execute(ctx context.Context, job *model.SyncJob) error
}

type Service struct {
	repo     *syncrepo.Repository
	provider syncprovider.Client
	cfg      config.SyncConfig
}

func New(repo *syncrepo.Repository, provider syncprovider.Client, cfg config.SyncConfig) *Service {
	return &Service{repo: repo, provider: provider, cfg: cfg}
}

func (s *Service) Execute(ctx context.Context, job *model.SyncJob) error {
	switch job.JobType {
	case "sync_teams":
		return s.executeSyncTeams(ctx, job)
	case "sync_players":
		return s.executeSyncPlayers(ctx, job)
	case "sync_matches_range":
		return s.executeSyncMatchesRange(ctx, job)
	case "sync_match_roster":
		return s.executeSyncMatchRoster(ctx, job)
	default:
		return fmt.Errorf("unsupported job type: %s", job.JobType)
	}
}

func (s *Service) executeSyncTeams(ctx context.Context, job *model.SyncJob) error {
	var payload struct {
		CompetitionCode string `json:"competitionCode"`
	}
	if err := json.Unmarshal(job.Payload, &payload); err != nil {
		return err
	}
	teams, err := s.provider.ListCompetitionTeams(ctx, payload.CompetitionCode)
	if err != nil {
		msg := err.Error()
		_ = s.repo.UpsertCursor(ctx, syncrepo.ProviderFootballData, job.JobType, job.ScopeKey, false, &msg)
		return err
	}
	now := time.Now().UTC()
	for _, item := range teams {
		if _, err := s.repo.UpsertTeam(ctx, syncrepo.UpsertTeamParams{
			Name:              item.Name,
			ShortName:         item.ShortName,
			Slug:              syncrepo.NormalizeSlug(item.Name),
			LogoURL:           item.CrestURL,
			ExternalSource:    syncrepo.ProviderFootballData,
			ExternalID:        item.ExternalID,
			ExternalUpdatedAt: &now,
		}); err != nil {
			msg := err.Error()
			_ = s.repo.UpsertCursor(ctx, syncrepo.ProviderFootballData, job.JobType, job.ScopeKey, false, &msg)
			return err
		}
	}
	return s.repo.UpsertCursor(ctx, syncrepo.ProviderFootballData, job.JobType, job.ScopeKey, true, nil)
}

func (s *Service) executeSyncPlayers(ctx context.Context, job *model.SyncJob) error {
	teams, err := s.repo.ListTeamsByExternalSource(ctx, syncrepo.ProviderFootballData)
	if err != nil {
		msg := err.Error()
		_ = s.repo.UpsertCursor(ctx, syncrepo.ProviderFootballData, job.JobType, job.ScopeKey, false, &msg)
		return err
	}
	now := time.Now().UTC()
	for _, team := range teams {
		if team.ExternalID == nil {
			continue
		}
		players, err := s.provider.ListTeamPlayers(ctx, *team.ExternalID)
		if err != nil {
			msg := err.Error()
			_ = s.repo.UpsertCursor(ctx, syncrepo.ProviderFootballData, job.JobType, job.ScopeKey, false, &msg)
			return err
		}
		for _, item := range players {
			if _, err := s.repo.UpsertPlayer(ctx, syncrepo.UpsertPlayerParams{
				TeamID:            team.ID,
				Name:              item.Name,
				Slug:              syncrepo.NormalizeSlug(item.Name),
				Position:          item.Position,
				AvatarURL:         nil,
				ExternalSource:    syncrepo.ProviderFootballData,
				ExternalID:        item.ExternalID,
				ExternalUpdatedAt: &now,
			}); err != nil {
				msg := err.Error()
				_ = s.repo.UpsertCursor(ctx, syncrepo.ProviderFootballData, job.JobType, job.ScopeKey, false, &msg)
				return err
			}
		}
	}
	return s.repo.UpsertCursor(ctx, syncrepo.ProviderFootballData, job.JobType, job.ScopeKey, true, nil)
}

func (s *Service) executeSyncMatchesRange(ctx context.Context, job *model.SyncJob) error {
	var payload struct {
		CompetitionCode string `json:"competitionCode"`
		DateFrom        string `json:"dateFrom"`
		DateTo          string `json:"dateTo"`
	}
	if err := json.Unmarshal(job.Payload, &payload); err != nil {
		return err
	}
	matches, err := s.provider.ListCompetitionMatches(ctx, payload.CompetitionCode, payload.DateFrom, payload.DateTo)
	if err != nil {
		msg := err.Error()
		_ = s.repo.UpsertCursor(ctx, syncrepo.ProviderFootballData, job.JobType, job.ScopeKey, false, &msg)
		return err
	}
	now := time.Now().UTC()
	for _, item := range matches {
		homeTeam, err := s.repo.UpsertTeam(ctx, syncrepo.UpsertTeamParams{
			Name:              item.HomeTeam.Name,
			ShortName:         item.HomeTeam.ShortName,
			Slug:              syncrepo.NormalizeSlug(item.HomeTeam.Name),
			LogoURL:           item.HomeTeam.CrestURL,
			ExternalSource:    syncrepo.ProviderFootballData,
			ExternalID:        item.HomeTeam.ExternalID,
			ExternalUpdatedAt: &now,
		})
		if err != nil {
			msg := err.Error()
			_ = s.repo.UpsertCursor(ctx, syncrepo.ProviderFootballData, job.JobType, job.ScopeKey, false, &msg)
			return err
		}
		awayTeam, err := s.repo.UpsertTeam(ctx, syncrepo.UpsertTeamParams{
			Name:              item.AwayTeam.Name,
			ShortName:         item.AwayTeam.ShortName,
			Slug:              syncrepo.NormalizeSlug(item.AwayTeam.Name),
			LogoURL:           item.AwayTeam.CrestURL,
			ExternalSource:    syncrepo.ProviderFootballData,
			ExternalID:        item.AwayTeam.ExternalID,
			ExternalUpdatedAt: &now,
		})
		if err != nil {
			msg := err.Error()
			_ = s.repo.UpsertCursor(ctx, syncrepo.ProviderFootballData, job.JobType, job.ScopeKey, false, &msg)
			return err
		}
		season := item.Season
		if season == "" {
			season = syncrepo.SeasonFromKickoff(item.KickoffAt)
		}
		_, err = s.repo.UpsertMatch(ctx, syncrepo.UpsertMatchParams{
			Competition:       item.Competition,
			Season:            season,
			Round:             item.Round,
			Status:            syncrepo.MapProviderMatchStatus(item.Status),
			KickoffAt:         item.KickoffAt,
			HomeTeamID:        homeTeam.ID,
			AwayTeamID:        awayTeam.ID,
			HomeScore:         item.HomeScore,
			AwayScore:         item.AwayScore,
			Venue:             item.Venue,
			ExternalSource:    syncrepo.ProviderFootballData,
			ExternalID:        item.ExternalID,
			ExternalUpdatedAt: &now,
		})
		if err != nil {
			msg := err.Error()
			_ = s.repo.UpsertCursor(ctx, syncrepo.ProviderFootballData, job.JobType, job.ScopeKey, false, &msg)
			return err
		}
	}
	return s.repo.UpsertCursor(ctx, syncrepo.ProviderFootballData, job.JobType, job.ScopeKey, true, nil)
}

func (s *Service) executeSyncMatchRoster(ctx context.Context, job *model.SyncJob) error {
	err := fmt.Errorf("sync_match_roster is not implemented yet")
	msg := err.Error()
	_ = s.repo.UpsertCursor(ctx, syncrepo.ProviderFootballData, job.JobType, job.ScopeKey, false, &msg)
	return err
}
