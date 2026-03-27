package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"final-whistle/backend/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrNoPendingJobs = errors.New("no pending sync jobs")

type EnqueueJobParams struct {
	JobType     string
	ScopeType   string
	ScopeKey    string
	DedupeKey   string
	TriggerMode model.SyncTriggerMode
	Priority    int
	ScheduledAt time.Time
	MaxAttempts int
	Payload     any
}

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) EnqueueJob(ctx context.Context, params EnqueueJobParams) (*model.SyncJob, error) {
	payloadBytes, err := json.Marshal(params.Payload)
	if err != nil {
		return nil, err
	}
	var existing model.SyncJob
	err = r.db.WithContext(ctx).
		Where("dedupe_key = ? AND status IN ?", params.DedupeKey, []model.SyncJobStatus{model.SyncJobStatusPending, model.SyncJobStatusRunning}).
		Order("id DESC").
		First(&existing).Error
	if err == nil {
		return &existing, nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	job := &model.SyncJob{
		JobType:     params.JobType,
		ScopeType:   params.ScopeType,
		ScopeKey:    params.ScopeKey,
		DedupeKey:   params.DedupeKey,
		TriggerMode: params.TriggerMode,
		Priority:    params.Priority,
		Status:      model.SyncJobStatusPending,
		ScheduledAt: params.ScheduledAt,
		MaxAttempts: params.MaxAttempts,
		Payload:     payloadBytes,
	}
	err = r.db.WithContext(ctx).Create(job).Error
	if err != nil {
		return nil, err
	}
	return job, nil
}

func (r *Repository) AcquireNextPendingJob(ctx context.Context) (*model.SyncJob, error) {
	var job model.SyncJob
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
			Where("status = ? AND scheduled_at <= ?", model.SyncJobStatusPending, time.Now().UTC()).
			Order("priority ASC, scheduled_at ASC, id ASC").
			First(&job).Error; err != nil {
			return err
		}

		now := time.Now().UTC()
		return tx.Model(&job).Updates(map[string]any{
			"status":     model.SyncJobStatusRunning,
			"started_at": now,
			"attempt":    gorm.Expr("attempt + 1"),
			"updated_at": now,
		}).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNoPendingJobs
		}
		return nil, err
	}
	if err := r.db.WithContext(ctx).First(&job, job.ID).Error; err != nil {
		return nil, err
	}
	return &job, nil
}

func (r *Repository) MarkJobSucceeded(ctx context.Context, id uint) error {
	now := time.Now().UTC()
	return r.db.WithContext(ctx).Model(&model.SyncJob{}).Where("id = ?", id).Updates(map[string]any{
		"status":      model.SyncJobStatusSucceeded,
		"finished_at": now,
		"updated_at":  now,
		"last_error":  nil,
	}).Error
}

func (r *Repository) MarkJobFailed(ctx context.Context, job *model.SyncJob, err error, retryAt *time.Time) error {
	now := time.Now().UTC()
	updates := map[string]any{
		"last_error": err.Error(),
		"updated_at": now,
	}
	if retryAt != nil && job.Attempt < job.MaxAttempts {
		updates["status"] = model.SyncJobStatusPending
		updates["scheduled_at"] = *retryAt
		updates["started_at"] = nil
	} else {
		updates["status"] = model.SyncJobStatusFailed
		updates["finished_at"] = now
	}
	return r.db.WithContext(ctx).Model(&model.SyncJob{}).Where("id = ?", job.ID).Updates(updates).Error
}

func (r *Repository) UpsertCursor(ctx context.Context, providerName, resourceType, scopeKey string, success bool, errMsg *string) error {
	now := time.Now().UTC()
	updates := map[string]any{
		"last_attempt_at": now,
		"updated_at":      now,
	}
	if success {
		updates["last_success_at"] = now
		updates["last_error"] = nil
		updates["last_error_at"] = nil
	} else {
		updates["last_error"] = errMsg
		updates["last_error_at"] = now
	}

	cursor := model.SyncCursor{
		Provider:     providerName,
		ResourceType: resourceType,
		ScopeKey:     scopeKey,
		Metadata:     []byte(`{}`),
	}
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "provider"}, {Name: "resource_type"}, {Name: "scope_key"}},
		DoUpdates: clause.Assignments(updates),
	}).Create(&cursor).Error
}

func (r *Repository) GetCursor(ctx context.Context, providerName, resourceType, scopeKey string) (*model.SyncCursor, error) {
	var cursor model.SyncCursor
	if err := r.db.WithContext(ctx).
		Where("provider = ? AND resource_type = ? AND scope_key = ?", providerName, resourceType, scopeKey).
		First(&cursor).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &cursor, nil
}

func (r *Repository) ListRecentJobs(ctx context.Context, limit int) ([]model.SyncJob, error) {
	if limit <= 0 {
		limit = 20
	}
	var jobs []model.SyncJob
	err := r.db.WithContext(ctx).Order("id DESC").Limit(limit).Find(&jobs).Error
	return jobs, err
}

type MatchRecord struct {
	ID          uint
	Competition string
	KickoffAt   time.Time
}

func (r *Repository) ListRelevantMatches(ctx context.Context, competition string, from, to time.Time) ([]MatchRecord, error) {
	var matches []MatchRecord
	err := r.db.WithContext(ctx).
		Table("matches").
		Select("id, competition, kickoff_at").
		Where("competition = ? AND kickoff_at BETWEEN ? AND ?", competition, from, to).
		Order("kickoff_at ASC").
		Scan(&matches).Error
	return matches, err
}

func (r *Repository) MatchRosterExists(ctx context.Context, matchID uint) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Table("match_players").Where("match_id = ?", matchID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func RetryBackoff(attempt int) time.Duration {
	switch attempt {
	case 1:
		return 5 * time.Minute
	case 2:
		return 15 * time.Minute
	default:
		return time.Hour
	}
}

func MatchesRangeScopeKey(competitionCode, date string) string {
	return fmt.Sprintf("matches_range:%s:%s", competitionCode, date)
}

func MatchRosterScopeKey(matchID uint) string {
	return fmt.Sprintf("match_roster:%d", matchID)
}
