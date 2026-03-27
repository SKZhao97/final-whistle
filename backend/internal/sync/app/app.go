package app

import (
	"context"

	"final-whistle/backend/internal/config"
	syncprovider "final-whistle/backend/internal/sync/provider"
	syncrepo "final-whistle/backend/internal/sync/repository"
	syncrunner "final-whistle/backend/internal/sync/runner"
	syncscheduler "final-whistle/backend/internal/sync/scheduler"
	syncservice "final-whistle/backend/internal/sync/service"
	"gorm.io/gorm"
)

type App struct {
	Repo      *syncrepo.Repository
	Service   *syncservice.Service
	Scheduler *syncscheduler.Scheduler
	Runner    *syncrunner.Runner
	Provider  syncprovider.Client
	Config    config.SyncConfig
}

func New(db *gorm.DB, cfg config.SyncConfig) *App {
	return NewWithProvider(db, cfg, syncprovider.NewFootballDataClient(cfg))
}

func NewWithProvider(db *gorm.DB, cfg config.SyncConfig, provider syncprovider.Client) *App {
	repo := syncrepo.New(db)
	service := syncservice.New(repo, provider, cfg)
	return &App{
		Repo:      repo,
		Service:   service,
		Scheduler: syncscheduler.New(repo, cfg),
		Runner:    syncrunner.New(repo, service, cfg),
		Provider:  provider,
		Config:    cfg,
	}
}

func (a *App) StartBackground(ctx context.Context) {
	if !a.Config.Enabled || !a.Config.AutoStart {
		return
	}
	if syncrunner.CanRunInProcess(a.Config.Role, "scheduler") {
		go a.Scheduler.Run(ctx)
	}
	if syncrunner.CanRunInProcess(a.Config.Role, "runner") {
		go a.Runner.Run(ctx)
	}
}
