package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"final-whistle/backend/internal/config"
	"final-whistle/backend/internal/db"
	"final-whistle/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Set Gin mode based on environment
	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Initialize database connection
	database, err := db.NewConnection(cfg.Database.URL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	log.Printf("Database connection established")

	// Create Gin router with middleware
	router := gin.New()

	// Apply middleware
	router.Use(middleware.RequestLogger())
	router.Use(middleware.ErrorRecovery())
	router.Use(middleware.CORS())

	// Set up routes
	setupRoutes(router, database)

	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Starting server on port %d in %s mode", cfg.Server.Port, cfg.Server.Env)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Give server time to finish existing requests
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}

func setupRoutes(router *gin.Engine, db *db.Database) {
	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().UTC(),
		})
	})

	// API route groups (to be extended by later specs)
	// These groups are intentionally empty - handlers will be added by subsequent specs
	api := router.Group("/api")
	{
		// Auth routes will be added by spec2
		_ = api.Group("/auth")

		// Match routes will be added by spec3
		_ = api.Group("/matches")

		// Check-in routes will be added by spec4
		_ = api.Group("/checkins")

		// Team routes will be added by spec3
		_ = api.Group("/teams")

		// Player routes will be added by spec3
		_ = api.Group("/players")

		// User routes will be added by spec5
		_ = api.Group("/users")
	}

	// Root endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"name":    "Final Whistle API",
			"version": "v1",
			"status":  "running",
		})
	})
}
