package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"final-whistle/backend/internal/config"
	"final-whistle/backend/internal/db"
	"final-whistle/backend/internal/model"
	syncapp "final-whistle/backend/internal/sync/app"
	syncrepo "final-whistle/backend/internal/sync/repository"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("usage: sync <daemon|bootstrap|run-once>")
	}

	switch os.Args[1] {
	case "daemon":
		runDaemon()
	case "bootstrap":
		runBootstrap()
	case "run-once":
		runOnce()
	default:
		log.Fatalf("unknown sync command: %s", os.Args[1])
	}
}

func runDaemon() {
	cfg, database := mustLoad()
	defer database.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	app := syncapp.New(database.DB, cfg.Sync)
	app.StartBackground(ctx)
	<-ctx.Done()
}

func runBootstrap() {
	fs := flag.NewFlagSet("bootstrap", flag.ExitOnError)
	competition := fs.String("competition", "PL", "competition code")
	fs.Parse(os.Args[2:])

	cfg, database := mustLoad()
	defer database.Close()

	repo := syncrepo.New(database.DB)
	now := time.Now().UTC()
	jobs := []syncrepo.EnqueueJobParams{
		{
			JobType:     "sync_teams",
			ScopeType:   "competition",
			ScopeKey:    "teams:" + *competition,
			DedupeKey:   "teams:" + *competition,
			TriggerMode: model.SyncTriggerModeManual,
			Priority:    10,
			ScheduledAt: now,
			MaxAttempts: 3,
			Payload:     map[string]any{"competitionCode": *competition},
		},
		{
			JobType:     "sync_players",
			ScopeType:   "competition",
			ScopeKey:    "players:" + *competition,
			DedupeKey:   "players:" + *competition,
			TriggerMode: model.SyncTriggerModeManual,
			Priority:    20,
			ScheduledAt: now,
			MaxAttempts: 3,
			Payload:     map[string]any{"competitionCode": *competition},
		},
	}
	for _, job := range jobs {
		if _, err := repo.EnqueueJob(context.Background(), job); err != nil {
			log.Fatalf("enqueue bootstrap job failed: %v", err)
		}
	}
	log.Printf("bootstrap jobs enqueued for %s with provider %s", *competition, cfg.Sync.Provider)
}

func runOnce() {
	cfg, database := mustLoad()
	defer database.Close()

	app := syncapp.New(database.DB, cfg.Sync)
	if err := app.Runner.RunOnce(context.Background()); err != nil {
		log.Fatalf("run-once failed: %v", err)
	}
}

func mustLoad() (*config.Config, *db.Database) {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	database, err := db.NewConnection(cfg.Database.URL)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	return cfg, database
}
