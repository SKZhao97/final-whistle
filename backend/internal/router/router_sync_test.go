package router_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"final-whistle/backend/internal/config"
	"final-whistle/backend/internal/db"
	"final-whistle/backend/internal/model"
	"final-whistle/backend/internal/router"
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

func TestPublicAPICanReadSyncedData(t *testing.T) {
	_, dbURL, cleanup := testutil.CreateTestDatabase(t)
	defer cleanup()

	database, err := db.NewConnection(dbURL)
	if err != nil {
		t.Fatalf("connect test db: %v", err)
	}
	defer database.Close()

	cfg := testSyncConfig()
	app := syncapp.NewWithProvider(database.DB, cfg, testutil.NewFakeProvider())
	ctx := context.Background()
	now := time.Now().UTC()

	jobs := []syncrepo.EnqueueJobParams{
		{
			JobType:     "sync_teams",
			ScopeType:   "competition",
			ScopeKey:    "teams:PL",
			DedupeKey:   "teams-api-test:" + now.Format(time.RFC3339Nano),
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
			DedupeKey:   "players-api-test:" + now.Add(time.Second).Format(time.RFC3339Nano),
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
			DedupeKey:   "matches-api-test:" + now.Add(2*time.Second).Format(time.RFC3339Nano),
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

	var teamID uint
	if err := database.DB.Table("teams").Select("id").Where("name = ?", "Arsenal").Scan(&teamID).Error; err != nil {
		t.Fatalf("load team id: %v", err)
	}
	var playerID uint
	if err := database.DB.Table("players").Select("id").Where("name = ?", "Bukayo Saka").Scan(&playerID).Error; err != nil {
		t.Fatalf("load player id: %v", err)
	}

	engine := router.New(database, "development", app, cfg)

	req := httptest.NewRequest(http.MethodGet, "/matches?page=1&pageSize=20", nil)
	res := httptest.NewRecorder()
	engine.ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("matches status code = %d body=%s", res.Code, res.Body.String())
	}

	var listResp struct {
		Success bool `json:"success"`
		Data    struct {
			Items []struct {
				ID          uint   `json:"id"`
				Competition string `json:"competition"`
				HomeTeam    struct {
					Name string `json:"name"`
				} `json:"homeTeam"`
				AwayTeam struct {
					Name string `json:"name"`
				} `json:"awayTeam"`
			} `json:"items"`
		} `json:"data"`
	}
	if err := json.Unmarshal(res.Body.Bytes(), &listResp); err != nil {
		t.Fatalf("decode match list: %v", err)
	}
	if len(listResp.Data.Items) != 1 {
		t.Fatalf("expected 1 match list item, got %d", len(listResp.Data.Items))
	}
	if listResp.Data.Items[0].HomeTeam.Name != "Arsenal" || listResp.Data.Items[0].AwayTeam.Name != "Liverpool" {
		t.Fatalf("unexpected teams in match list: %+v", listResp.Data.Items[0])
	}

	matchID := listResp.Data.Items[0].ID
	req = httptest.NewRequest(http.MethodGet, "/matches/"+strconv.FormatUint(uint64(matchID), 10), nil)
	res = httptest.NewRecorder()
	engine.ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("match detail status code = %d body=%s", res.Code, res.Body.String())
	}

	req = httptest.NewRequest(http.MethodGet, "/teams/"+strconv.FormatUint(uint64(teamID), 10), nil)
	res = httptest.NewRecorder()
	engine.ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("team detail status code = %d body=%s", res.Code, res.Body.String())
	}

	req = httptest.NewRequest(http.MethodGet, "/players/"+strconv.FormatUint(uint64(playerID), 10), nil)
	res = httptest.NewRecorder()
	engine.ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("player detail status code = %d body=%s", res.Code, res.Body.String())
	}
}
