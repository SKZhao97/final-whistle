package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	App      AppConfig
	Sync     SyncConfig
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Port int
	Env  string
}

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	URL string
}

// AppConfig holds application-level configuration
type AppConfig struct {
	LogLevel string
}

type SyncConfig struct {
	Enabled                          bool
	AutoStart                        bool
	Role                             string
	Provider                         string
	CompetitionCode                  string
	AdminToken                       string
	ScanIntervalSeconds              int
	AcquireIntervalSeconds           int
	MaxWorkers                       int
	SafeRateLimitPerMinute           int
	MatchLookbackHours               int
	MatchLookaheadDays               int
	WindowFarMatchDays               int
	WindowPreMatchMinutes            int
	WindowLiveAfterKickoffMinutes    int
	WindowPostMatchMinutes           int
	ScheduleFarMatchEveryMinutes     int
	SchedulePreMatchEveryMinutes     int
	ScheduleLiveEveryMinutes         int
	SchedulePostMatchEveryMinutes    int
	RosterWindowBeforeKickoffMinutes int
	RosterWindowAfterKickoffMinutes  int
	RosterScheduleEveryMinutes       int
	TeamSyncHours                    int
	PlayerSyncHours                  int
	FootballDataAPIToken             string
}

// Load loads configuration from environment variables and defaults
func Load() (*Config, error) {
	// Set up Viper
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	// Set defaults
	setDefaults()

	// Read config file (optional)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Printf("Warning reading config file: %v", err)
		}
	}

	// Bind environment variables
	bindEnvVars()

	// Create config struct
	cfg := &Config{
		Server: ServerConfig{
			Port: viper.GetInt("server.port"),
			Env:  viper.GetString("server.env"),
		},
		Database: DatabaseConfig{
			URL: viper.GetString("database.url"),
		},
		App: AppConfig{
			LogLevel: viper.GetString("app.log_level"),
		},
		Sync: SyncConfig{
			Enabled:                          viper.GetBool("sync.enabled"),
			AutoStart:                        viper.GetBool("sync.auto_start"),
			Role:                             viper.GetString("sync.role"),
			Provider:                         viper.GetString("sync.provider"),
			CompetitionCode:                  viper.GetString("sync.competition_code"),
			AdminToken:                       viper.GetString("sync.admin_token"),
			ScanIntervalSeconds:              viper.GetInt("sync.scan_interval_seconds"),
			AcquireIntervalSeconds:           viper.GetInt("sync.acquire_interval_seconds"),
			MaxWorkers:                       viper.GetInt("sync.max_workers"),
			SafeRateLimitPerMinute:           viper.GetInt("sync.safe_rate_limit_per_minute"),
			MatchLookbackHours:               viper.GetInt("sync.match_lookback_hours"),
			MatchLookaheadDays:               viper.GetInt("sync.match_lookahead_days"),
			WindowFarMatchDays:               viper.GetInt("sync.window_far_match_days"),
			WindowPreMatchMinutes:            viper.GetInt("sync.window_pre_match_minutes"),
			WindowLiveAfterKickoffMinutes:    viper.GetInt("sync.window_live_after_kickoff_minutes"),
			WindowPostMatchMinutes:           viper.GetInt("sync.window_post_match_minutes"),
			ScheduleFarMatchEveryMinutes:     viper.GetInt("sync.schedule_far_match_every_minutes"),
			SchedulePreMatchEveryMinutes:     viper.GetInt("sync.schedule_pre_match_every_minutes"),
			ScheduleLiveEveryMinutes:         viper.GetInt("sync.schedule_live_every_minutes"),
			SchedulePostMatchEveryMinutes:    viper.GetInt("sync.schedule_post_match_every_minutes"),
			RosterWindowBeforeKickoffMinutes: viper.GetInt("sync.roster_window_before_kickoff_minutes"),
			RosterWindowAfterKickoffMinutes:  viper.GetInt("sync.roster_window_after_kickoff_minutes"),
			RosterScheduleEveryMinutes:       viper.GetInt("sync.roster_schedule_every_minutes"),
			TeamSyncHours:                    viper.GetInt("sync.team_sync_hours"),
			PlayerSyncHours:                  viper.GetInt("sync.player_sync_hours"),
			FootballDataAPIToken:             viper.GetString("sync.football_data_api_token"),
		},
	}

	// Validate required configuration
	if err := validate(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func setDefaults() {
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.env", "development")
	viper.SetDefault("database.url", "postgres://postgres:postgres@localhost:5432/final_whistle?sslmode=disable")
	viper.SetDefault("app.log_level", "info")
	viper.SetDefault("sync.enabled", false)
	viper.SetDefault("sync.auto_start", false)
	viper.SetDefault("sync.role", "api")
	viper.SetDefault("sync.provider", "football-data")
	viper.SetDefault("sync.competition_code", "PL")
	viper.SetDefault("sync.scan_interval_seconds", 60)
	viper.SetDefault("sync.acquire_interval_seconds", 15)
	viper.SetDefault("sync.max_workers", 1)
	viper.SetDefault("sync.safe_rate_limit_per_minute", 8)
	viper.SetDefault("sync.match_lookback_hours", 24)
	viper.SetDefault("sync.match_lookahead_days", 7)
	viper.SetDefault("sync.window_far_match_days", 7)
	viper.SetDefault("sync.window_pre_match_minutes", 90)
	viper.SetDefault("sync.window_live_after_kickoff_minutes", 180)
	viper.SetDefault("sync.window_post_match_minutes", 360)
	viper.SetDefault("sync.schedule_far_match_every_minutes", 2880)
	viper.SetDefault("sync.schedule_pre_match_every_minutes", 30)
	viper.SetDefault("sync.schedule_live_every_minutes", 2)
	viper.SetDefault("sync.schedule_post_match_every_minutes", 20)
	viper.SetDefault("sync.roster_window_before_kickoff_minutes", 75)
	viper.SetDefault("sync.roster_window_after_kickoff_minutes", 120)
	viper.SetDefault("sync.roster_schedule_every_minutes", 10)
	viper.SetDefault("sync.team_sync_hours", 24)
	viper.SetDefault("sync.player_sync_hours", 24)
}

func bindEnvVars() {
	viper.BindEnv("server.port", "PORT")
	viper.BindEnv("server.env", "ENV")
	viper.BindEnv("database.url", "DATABASE_URL")
	viper.BindEnv("app.log_level", "LOG_LEVEL")
	viper.BindEnv("sync.enabled", "SYNC_ENABLED")
	viper.BindEnv("sync.auto_start", "SYNC_AUTO_START")
	viper.BindEnv("sync.role", "SYNC_ROLE")
	viper.BindEnv("sync.provider", "SYNC_PROVIDER")
	viper.BindEnv("sync.competition_code", "SYNC_COMPETITION_CODE")
	viper.BindEnv("sync.admin_token", "SYNC_ADMIN_TOKEN")
	viper.BindEnv("sync.scan_interval_seconds", "SYNC_SCAN_INTERVAL_SECONDS")
	viper.BindEnv("sync.acquire_interval_seconds", "SYNC_ACQUIRE_INTERVAL_SECONDS")
	viper.BindEnv("sync.max_workers", "SYNC_MAX_WORKERS")
	viper.BindEnv("sync.safe_rate_limit_per_minute", "SYNC_SAFE_RATE_LIMIT_PER_MINUTE")
	viper.BindEnv("sync.match_lookback_hours", "SYNC_MATCH_LOOKBACK_HOURS")
	viper.BindEnv("sync.match_lookahead_days", "SYNC_MATCH_LOOKAHEAD_DAYS")
	viper.BindEnv("sync.window_far_match_days", "SYNC_WINDOW_FAR_MATCH_DAYS")
	viper.BindEnv("sync.window_pre_match_minutes", "SYNC_WINDOW_PRE_MATCH_MINUTES")
	viper.BindEnv("sync.window_live_after_kickoff_minutes", "SYNC_WINDOW_LIVE_AFTER_KICKOFF_MINUTES")
	viper.BindEnv("sync.window_post_match_minutes", "SYNC_WINDOW_POST_MATCH_MINUTES")
	viper.BindEnv("sync.schedule_far_match_every_minutes", "SYNC_SCHEDULE_FAR_MATCH_EVERY_MINUTES")
	viper.BindEnv("sync.schedule_pre_match_every_minutes", "SYNC_SCHEDULE_PRE_MATCH_EVERY_MINUTES")
	viper.BindEnv("sync.schedule_live_every_minutes", "SYNC_SCHEDULE_LIVE_EVERY_MINUTES")
	viper.BindEnv("sync.schedule_post_match_every_minutes", "SYNC_SCHEDULE_POST_MATCH_EVERY_MINUTES")
	viper.BindEnv("sync.roster_window_before_kickoff_minutes", "SYNC_ROSTER_WINDOW_BEFORE_KICKOFF_MINUTES")
	viper.BindEnv("sync.roster_window_after_kickoff_minutes", "SYNC_ROSTER_WINDOW_AFTER_KICKOFF_MINUTES")
	viper.BindEnv("sync.roster_schedule_every_minutes", "SYNC_ROSTER_SCHEDULE_EVERY_MINUTES")
	viper.BindEnv("sync.team_sync_hours", "SYNC_TEAM_SYNC_HOURS")
	viper.BindEnv("sync.player_sync_hours", "SYNC_PLAYER_SYNC_HOURS")
	viper.BindEnv("sync.football_data_api_token", "FOOTBALL_DATA_API_TOKEN")
}

func validate(cfg *Config) error {
	if cfg.Server.Port < 1 || cfg.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", cfg.Server.Port)
	}

	if cfg.Database.URL == "" {
		return fmt.Errorf("database URL is required")
	}

	validEnvs := map[string]bool{"development": true, "staging": true, "production": true}
	if !validEnvs[cfg.Server.Env] {
		return fmt.Errorf("invalid environment: %s", cfg.Server.Env)
	}

	validLogLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if !validLogLevels[cfg.App.LogLevel] {
		return fmt.Errorf("invalid log level: %s", cfg.App.LogLevel)
	}

	validSyncRoles := map[string]bool{"all": true, "api": true, "scheduler": true, "runner": true}
	if !validSyncRoles[cfg.Sync.Role] {
		return fmt.Errorf("invalid sync role: %s", cfg.Sync.Role)
	}

	if cfg.Sync.SafeRateLimitPerMinute < 1 || cfg.Sync.SafeRateLimitPerMinute > 10 {
		return fmt.Errorf("invalid sync safe rate limit per minute: %d", cfg.Sync.SafeRateLimitPerMinute)
	}

	intsMustBePositive := map[string]int{
		"sync.scan_interval_seconds":                cfg.Sync.ScanIntervalSeconds,
		"sync.acquire_interval_seconds":             cfg.Sync.AcquireIntervalSeconds,
		"sync.max_workers":                          cfg.Sync.MaxWorkers,
		"sync.match_lookback_hours":                 cfg.Sync.MatchLookbackHours,
		"sync.match_lookahead_days":                 cfg.Sync.MatchLookaheadDays,
		"sync.window_far_match_days":                cfg.Sync.WindowFarMatchDays,
		"sync.window_pre_match_minutes":             cfg.Sync.WindowPreMatchMinutes,
		"sync.window_live_after_kickoff_minutes":    cfg.Sync.WindowLiveAfterKickoffMinutes,
		"sync.window_post_match_minutes":            cfg.Sync.WindowPostMatchMinutes,
		"sync.schedule_far_match_every_minutes":     cfg.Sync.ScheduleFarMatchEveryMinutes,
		"sync.schedule_pre_match_every_minutes":     cfg.Sync.SchedulePreMatchEveryMinutes,
		"sync.schedule_live_every_minutes":          cfg.Sync.ScheduleLiveEveryMinutes,
		"sync.schedule_post_match_every_minutes":    cfg.Sync.SchedulePostMatchEveryMinutes,
		"sync.roster_window_before_kickoff_minutes": cfg.Sync.RosterWindowBeforeKickoffMinutes,
		"sync.roster_window_after_kickoff_minutes":  cfg.Sync.RosterWindowAfterKickoffMinutes,
		"sync.roster_schedule_every_minutes":        cfg.Sync.RosterScheduleEveryMinutes,
		"sync.team_sync_hours":                      cfg.Sync.TeamSyncHours,
		"sync.player_sync_hours":                    cfg.Sync.PlayerSyncHours,
	}
	for key, value := range intsMustBePositive {
		if value <= 0 {
			return fmt.Errorf("%s must be > 0", key)
		}
	}

	if cfg.Sync.WindowFarMatchDays > cfg.Sync.MatchLookaheadDays {
		return fmt.Errorf("sync.window_far_match_days cannot exceed sync.match_lookahead_days")
	}

	if cfg.Sync.WindowPreMatchMinutes < cfg.Sync.RosterWindowBeforeKickoffMinutes {
		return fmt.Errorf("sync.window_pre_match_minutes cannot be smaller than sync.roster_window_before_kickoff_minutes")
	}

	return nil
}

// GetEnv returns environment variable with fallback
func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetEnvInt returns environment variable as int with fallback
func GetEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
