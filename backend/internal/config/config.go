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
}

func bindEnvVars() {
	viper.BindEnv("server.port", "PORT")
	viper.BindEnv("server.env", "ENV")
	viper.BindEnv("database.url", "DATABASE_URL")
	viper.BindEnv("app.log_level", "LOG_LEVEL")
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