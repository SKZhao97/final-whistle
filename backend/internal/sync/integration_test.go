package sync_test

import (
	"context"
	"testing"
	"time"

	"final-whistle/backend/internal/config"
	"final-whistle/backend/internal/db"
	"final-whistle/backend/internal/model"
	syncapp "final-whistle/backend/internal/sync/app"
	syncrepo "final-whistle/backend/internal/sync/repository"
	"final-whistle/backend/internal/testutil"
)

func testSyncConfig() config.SyncConfig {
	return config.SyncConfig{
		Enabled:                          true,
		AutoStart:                        false,
		Role:                             "all",
		Provider:                         "football-data",
		CompetitionCode:                  "PL",
		ScanIntervalSeconds:              60,
		AcquireIntervalSeconds:           1,
		MaxWorkers:                       1,
		SafeRateLimitPerMinute:           8,
		MatchLookbackHours:               24,
		MatchLookaheadDays:               7,
		WindowFarMatchDays:               7,
		WindowPreMatchMinutes:            90,
		WindowLiveAfterKickoffMinutes:    180,
		WindowPostMatchMinutes:           360,
		ScheduleFarMatchEveryMinutes:     2880,
		SchedulePreMatchEveryMinutes:     30,
		ScheduleLiveEveryMinutes:         2,
		SchedulePostMatchEveryMinutes:    20,
		RosterWindowBeforeKickoffMinutes: 75,
		RosterWindowAfterKickoffMinutes:  120,
		RosterScheduleEveryMinutes:       10,
		TeamSyncHours:                    24,
		PlayerSyncHours:                  24,
		AdminToken:                       "test-token",
	}
}

func TestSyncWritesTeamsPlayersAndMatches(t *testing.T) {
	_, dbURL, cleanup := testutil.CreateTestDatabase(t)
	defer cleanup()

	database, err := db.NewConnection(dbURL)
	if err != nil {
		t.Fatalf("connect test db: %v", err)
	}
	defer database.Close()

	app := syncapp.NewWithProvider(database.DB, testSyncConfig(), testutil.NewFakeProvider())
	ctx := context.Background()
	now := time.Now().UTC()

	jobs := []syncrepo.EnqueueJobParams{
		{
			JobType:     "sync_teams",
			ScopeType:   "competition",
			ScopeKey:    "teams:PL",
			DedupeKey:   "teams:PL:" + now.Format(time.RFC3339Nano),
			TriggerMode: model.SyncTriggerModeManual,
			Priority:    10,
			ScheduledAt: now,
			MaxAttempts: 3,
			Payload:     map[string]any{"competitionCode": "PL"},
		},
		{
			JobType:     "sync_players",
			ScopeType:   "competition",
			ScopeKey:    "players:PL",
			DedupeKey:   "players:PL:" + now.Add(time.Second).Format(time.RFC3339Nano),
			TriggerMode: model.SyncTriggerModeManual,
			Priority:    20,
			ScheduledAt: now,
			MaxAttempts: 3,
			Payload:     map[string]any{"competitionCode": "PL"},
		},
		{
			JobType:     "sync_matches_range",
			ScopeType:   "date",
			ScopeKey:    "matches_range:PL:2026-03-27",
			DedupeKey:   "matches_range:PL:2026-03-27:" + now.Add(2*time.Second).Format(time.RFC3339Nano),
			TriggerMode: model.SyncTriggerModeManual,
			Priority:    30,
			ScheduledAt: now,
			MaxAttempts: 3,
			Payload: map[string]any{
				"competitionCode": "PL",
				"dateFrom":        "2026-03-27",
				"dateTo":          "2026-03-27",
			},
		},
	}
	for _, job := range jobs {
		if _, err := app.Repo.EnqueueJob(ctx, job); err != nil {
			t.Fatalf("enqueue job %s: %v", job.JobType, err)
		}
	}

	for i := 0; i < 3; i++ {
		if err := app.Runner.RunOnce(ctx); err != nil {
			t.Fatalf("runner run once %d: %v", i, err)
		}
	}

	var teamCount, playerCount, matchCount int64
	if err := database.DB.Table("teams").Count(&teamCount).Error; err != nil {
		t.Fatalf("count teams: %v", err)
	}
	if err := database.DB.Table("players").Count(&playerCount).Error; err != nil {
		t.Fatalf("count players: %v", err)
	}
	if err := database.DB.Table("matches").Count(&matchCount).Error; err != nil {
		t.Fatalf("count matches: %v", err)
	}

	if teamCount != 2 {
		t.Fatalf("expected 2 teams, got %d", teamCount)
	}
	if playerCount != 4 {
		t.Fatalf("expected 4 players, got %d", playerCount)
	}
	if matchCount != 1 {
		t.Fatalf("expected 1 match, got %d", matchCount)
	}
}
