package handler

import (
	"net/http"
	"time"

	"final-whistle/backend/internal/config"
	"final-whistle/backend/internal/model"
	syncrepo "final-whistle/backend/internal/sync/repository"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	repo *syncrepo.Repository
	cfg  config.SyncConfig
}

func New(repo *syncrepo.Repository, cfg config.SyncConfig) *Handler {
	return &Handler{repo: repo, cfg: cfg}
}

func (h *Handler) RegisterRoutes(router gin.IRoutes) {
	router.GET("/status", h.Status)
	router.GET("/jobs", h.ListJobs)
	router.POST("/jobs", h.EnqueueJob)
	router.POST("/bootstrap", h.Bootstrap)
}

func (h *Handler) RequireAdminToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		if h.cfg.AdminToken == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "sync admin token is not configured"})
			return
		}
		if c.GetHeader("Authorization") != "Bearer "+h.cfg.AdminToken {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid sync admin token"})
			return
		}
		c.Next()
	}
}

func (h *Handler) Status(c *gin.Context) {
	jobs, err := h.repo.ListRecentJobs(c.Request.Context(), 100)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	pending, running, failed := 0, 0, 0
	for _, job := range jobs {
		switch job.Status {
		case model.SyncJobStatusPending:
			pending++
		case model.SyncJobStatusRunning:
			running++
		case model.SyncJobStatusFailed:
			failed++
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"enabled":     h.cfg.Enabled,
		"autoStart":   h.cfg.AutoStart,
		"role":        h.cfg.Role,
		"provider":    h.cfg.Provider,
		"competition": h.cfg.CompetitionCode,
		"pendingJobs": pending,
		"runningJobs": running,
		"failedJobs":  failed,
	})
}

func (h *Handler) ListJobs(c *gin.Context) {
	jobs, err := h.repo.ListRecentJobs(c.Request.Context(), 50)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": jobs})
}

type enqueueJobRequest struct {
	JobType   string         `json:"jobType"`
	ScopeType string         `json:"scopeType"`
	ScopeKey  string         `json:"scopeKey"`
	DedupeKey string         `json:"dedupeKey"`
	Priority  int            `json:"priority"`
	Payload   map[string]any `json:"payload"`
}

func (h *Handler) EnqueueJob(c *gin.Context) {
	var req enqueueJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.JobType == "" || req.ScopeType == "" || req.ScopeKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "jobType, scopeType, and scopeKey are required"})
		return
	}
	dedupeKey := req.DedupeKey
	if dedupeKey == "" {
		dedupeKey = req.ScopeKey
	}
	job, err := h.repo.EnqueueJob(c.Request.Context(), syncrepo.EnqueueJobParams{
		JobType:     req.JobType,
		ScopeType:   req.ScopeType,
		ScopeKey:    req.ScopeKey,
		DedupeKey:   dedupeKey,
		TriggerMode: model.SyncTriggerModeManual,
		Priority:    req.Priority,
		ScheduledAt: time.Now().UTC(),
		MaxAttempts: 3,
		Payload:     req.Payload,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"item": job})
}

func (h *Handler) Bootstrap(c *gin.Context) {
	ctx := c.Request.Context()
	now := time.Now().UTC()
	jobs := []syncrepo.EnqueueJobParams{
		{
			JobType:     "sync_teams",
			ScopeType:   "competition",
			ScopeKey:    "teams:" + h.cfg.CompetitionCode,
			DedupeKey:   "teams:" + h.cfg.CompetitionCode,
			TriggerMode: model.SyncTriggerModeManual,
			Priority:    10,
			ScheduledAt: now,
			MaxAttempts: 3,
			Payload:     map[string]any{"competitionCode": h.cfg.CompetitionCode},
		},
		{
			JobType:     "sync_players",
			ScopeType:   "competition",
			ScopeKey:    "players:" + h.cfg.CompetitionCode,
			DedupeKey:   "players:" + h.cfg.CompetitionCode,
			TriggerMode: model.SyncTriggerModeManual,
			Priority:    20,
			ScheduledAt: now,
			MaxAttempts: 3,
			Payload:     map[string]any{"competitionCode": h.cfg.CompetitionCode},
		},
	}

	created := make([]*model.SyncJob, 0, len(jobs))
	for _, job := range jobs {
		item, err := h.repo.EnqueueJob(ctx, job)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		created = append(created, item)
	}
	c.JSON(http.StatusCreated, gin.H{"items": created})
}
