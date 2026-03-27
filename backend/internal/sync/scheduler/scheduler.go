package scheduler

import (
	"context"
	"fmt"
	"log"
	"time"

	"final-whistle/backend/internal/config"
	"final-whistle/backend/internal/model"
	"final-whistle/backend/internal/sync/policy"
	syncrepo "final-whistle/backend/internal/sync/repository"
)

type Scheduler struct {
	repo     *syncrepo.Repository
	cfg      config.SyncConfig
	policy   policy.MatchSchedulePolicy
	interval time.Duration
	nowFn    func() time.Time
}

func New(repo *syncrepo.Repository, cfg config.SyncConfig) *Scheduler {
	return &Scheduler{
		repo:     repo,
		cfg:      cfg,
		policy:   policy.New(cfg),
		interval: time.Duration(cfg.ScanIntervalSeconds) * time.Second,
		nowFn:    func() time.Time { return time.Now().UTC() },
	}
}

func (s *Scheduler) Run(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := s.RunOnce(ctx); err != nil {
				log.Printf("sync scheduler run failed: %v", err)
			}
		}
	}
}

func (s *Scheduler) RunOnce(ctx context.Context) error {
	now := s.nowFn()
	if err := s.ensureStaticJobs(ctx, now); err != nil {
		return err
	}

	from := now.Add(-s.policy.Lookback)
	to := now.Add(s.policy.Lookahead)
	matches, err := s.repo.ListRelevantMatches(ctx, "Premier League", from, to)
	if err != nil {
		return err
	}

	if err := s.scheduleMatchRangeJobs(ctx, now, matches); err != nil {
		return err
	}
	return s.scheduleRosterJobs(ctx, now, matches)
}

func (s *Scheduler) ensureStaticJobs(ctx context.Context, now time.Time) error {
	staticJobs := []struct {
		jobType   string
		scopeKey  string
		interval  time.Duration
		scopeType string
		priority  int
	}{
		{"sync_teams", fmt.Sprintf("teams:%s", s.cfg.CompetitionCode), time.Duration(s.cfg.TeamSyncHours) * time.Hour, "competition", 50},
		{"sync_players", fmt.Sprintf("players:%s", s.cfg.CompetitionCode), time.Duration(s.cfg.PlayerSyncHours) * time.Hour, "competition", 60},
	}

	for _, item := range staticJobs {
		cursor, err := s.repo.GetCursor(ctx, "football-data", item.jobType, item.scopeKey)
		if err != nil {
			return err
		}
		if cursor != nil && cursor.LastSuccessAt != nil && now.Sub(*cursor.LastSuccessAt) < item.interval {
			continue
		}
		_, err = s.repo.EnqueueJob(ctx, syncrepo.EnqueueJobParams{
			JobType:     item.jobType,
			ScopeType:   item.scopeType,
			ScopeKey:    item.scopeKey,
			DedupeKey:   item.scopeKey,
			TriggerMode: model.SyncTriggerModeAutomatic,
			Priority:    item.priority,
			ScheduledAt: now,
			MaxAttempts: 3,
			Payload: map[string]any{
				"competitionCode": s.cfg.CompetitionCode,
			},
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Scheduler) scheduleMatchRangeJobs(ctx context.Context, now time.Time, matches []syncrepo.MatchRecord) error {
	type dayState struct {
		minInterval time.Duration
	}
	byDate := map[string]dayState{}
	for _, match := range matches {
		window := s.policy.ClassifyMatchWindow(now, match.KickoffAt)
		interval := s.policy.IntervalForWindow(window)
		if interval == 0 {
			continue
		}
		dateKey := match.KickoffAt.UTC().Format("2006-01-02")
		current, ok := byDate[dateKey]
		if !ok || interval < current.minInterval {
			byDate[dateKey] = dayState{minInterval: interval}
		}
	}

	for dateKey, state := range byDate {
		scopeKey := syncrepo.MatchesRangeScopeKey(s.cfg.CompetitionCode, dateKey)
		cursor, err := s.repo.GetCursor(ctx, "football-data", "sync_matches_range", scopeKey)
		if err != nil {
			return err
		}
		if cursor != nil && cursor.LastSuccessAt != nil && now.Sub(*cursor.LastSuccessAt) < state.minInterval {
			continue
		}
		_, err = s.repo.EnqueueJob(ctx, syncrepo.EnqueueJobParams{
			JobType:     "sync_matches_range",
			ScopeType:   "date",
			ScopeKey:    scopeKey,
			DedupeKey:   scopeKey,
			TriggerMode: model.SyncTriggerModeAutomatic,
			Priority:    10,
			ScheduledAt: now,
			MaxAttempts: 3,
			Payload: map[string]any{
				"competitionCode": s.cfg.CompetitionCode,
				"dateFrom":        dateKey,
				"dateTo":          dateKey,
			},
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Scheduler) scheduleRosterJobs(ctx context.Context, now time.Time, matches []syncrepo.MatchRecord) error {
	for _, match := range matches {
		if !s.policy.InRosterWindow(now, match.KickoffAt) {
			continue
		}
		exists, err := s.repo.MatchRosterExists(ctx, match.ID)
		if err != nil {
			return err
		}
		if exists {
			continue
		}

		scopeKey := syncrepo.MatchRosterScopeKey(match.ID)
		cursor, err := s.repo.GetCursor(ctx, "football-data", "sync_match_roster", scopeKey)
		if err != nil {
			return err
		}
		if cursor != nil && cursor.LastSuccessAt != nil && now.Sub(*cursor.LastSuccessAt) < s.policy.RosterEvery {
			continue
		}

		_, err = s.repo.EnqueueJob(ctx, syncrepo.EnqueueJobParams{
			JobType:     "sync_match_roster",
			ScopeType:   "match",
			ScopeKey:    scopeKey,
			DedupeKey:   scopeKey,
			TriggerMode: model.SyncTriggerModeAutomatic,
			Priority:    20,
			ScheduledAt: now,
			MaxAttempts: 3,
			Payload: map[string]any{
				"matchId": match.ID,
			},
		})
		if err != nil {
			return err
		}
	}
	return nil
}
