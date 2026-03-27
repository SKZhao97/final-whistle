package runner

import (
	"context"
	"errors"
	"log"
	"time"

	"final-whistle/backend/internal/config"
	"final-whistle/backend/internal/model"
	syncrepo "final-whistle/backend/internal/sync/repository"
	syncservice "final-whistle/backend/internal/sync/service"
)

type Runner struct {
	repo     *syncrepo.Repository
	service  *syncservice.Service
	interval time.Duration
}

func New(repo *syncrepo.Repository, service *syncservice.Service, cfg config.SyncConfig) *Runner {
	return &Runner{
		repo:     repo,
		service:  service,
		interval: time.Duration(cfg.AcquireIntervalSeconds) * time.Second,
	}
}

func (r *Runner) Run(ctx context.Context) {
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			r.runOnce(ctx)
		}
	}
}

func (r *Runner) RunOnce(ctx context.Context) error {
	return r.runOnce(ctx)
}

func (r *Runner) runOnce(ctx context.Context) error {
	job, err := r.repo.AcquireNextPendingJob(ctx)
	if err != nil {
		if errors.Is(err, syncrepo.ErrNoPendingJobs) {
			return nil
		}
		return err
	}

	if err := r.service.Execute(ctx, job); err != nil {
		retryAt := time.Now().UTC().Add(syncrepo.RetryBackoff(job.Attempt))
		if job.Attempt >= job.MaxAttempts {
			retryAt = time.Time{}
		}
		var retryPtr *time.Time
		if !retryAt.IsZero() {
			retryPtr = &retryAt
		}
		if markErr := r.repo.MarkJobFailed(ctx, job, err, retryPtr); markErr != nil {
			log.Printf("Failed to mark sync job %d failed: %v", job.ID, markErr)
		}
		return err
	}
	return r.repo.MarkJobSucceeded(ctx, job.ID)
}

func IsEnabledForRole(role string, component string) bool {
	switch role {
	case "all":
		return true
	case "scheduler":
		return component == "scheduler"
	case "runner":
		return component == "runner"
	case "api":
		return false
	default:
		return false
	}
}

func CanRunInProcess(role string, component string) bool {
	if role == "all" {
		return true
	}
	return IsEnabledForRole(role, component)
}

var _ = model.SyncJob{}
