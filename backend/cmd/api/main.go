package main

import (
	"context"
	"final-whistle/backend/internal/handler"
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
	"final-whistle/backend/internal/repository"
	"final-whistle/backend/internal/service"
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
	setupRoutes(router, database, cfg.Server.Env)

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

func setupRoutes(router *gin.Engine, database *db.Database, env string) {
	authRepository := repository.NewAuthRepository(database.DB)
	matchRepository := repository.NewMatchRepository(database.DB)
	teamRepository := repository.NewTeamRepository(database.DB)
	playerRepository := repository.NewPlayerRepository(database.DB)

	authService := service.NewAuthService(authRepository, env == "development")
	authHandler := handler.NewAuthHandler(authService, env)
	matchHandler := handler.NewMatchHandler(service.NewMatchService(matchRepository))
	teamHandler := handler.NewTeamHandler(service.NewTeamService(teamRepository, matchRepository))
	playerHandler := handler.NewPlayerHandler(service.NewPlayerService(playerRepository))

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().UTC(),
		})
	})

	router.Use(middleware.ResolveCurrentUser(authService))

	router.POST("/auth/login", authHandler.Login)
	router.POST("/auth/logout", authHandler.Logout)
	router.GET("/auth/me", authHandler.Me)

	router.GET("/matches", matchHandler.List)
	router.GET("/matches/:id", matchHandler.Detail)
	router.GET("/teams/:id", teamHandler.Detail)
	router.GET("/players/:id", playerHandler.Detail)

	// Root endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"name":    "Final Whistle API",
			"version": "v1",
			"status":  "running",
		})
	})
}
