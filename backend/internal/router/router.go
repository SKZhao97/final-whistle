package router

import (
	"net/http"
	"time"

	"final-whistle/backend/internal/config"
	"final-whistle/backend/internal/db"
	"final-whistle/backend/internal/handler"
	"final-whistle/backend/internal/middleware"
	"final-whistle/backend/internal/repository"
	"final-whistle/backend/internal/service"
	syncapp "final-whistle/backend/internal/sync/app"
	synchandler "final-whistle/backend/internal/sync/handler"
	"github.com/gin-gonic/gin"
)

func New(database *db.Database, env string, syncRuntime *syncapp.App, syncCfg config.SyncConfig) *gin.Engine {
	authRepository := repository.NewAuthRepository(database.DB)
	checkInRepository := repository.NewCheckInRepository(database.DB)
	matchRepository := repository.NewMatchRepository(database.DB)
	teamRepository := repository.NewTeamRepository(database.DB)
	playerRepository := repository.NewPlayerRepository(database.DB)
	userRepository := repository.NewUserRepository(database.DB)

	authService := service.NewAuthService(authRepository, env == "development")
	authHandler := handler.NewAuthHandler(authService, env)
	checkInHandler := handler.NewCheckInHandler(service.NewCheckInService(checkInRepository))
	matchHandler := handler.NewMatchHandler(service.NewMatchService(matchRepository))
	teamHandler := handler.NewTeamHandler(service.NewTeamService(teamRepository, matchRepository))
	playerHandler := handler.NewPlayerHandler(service.NewPlayerService(playerRepository))
	userHandler := handler.NewUserHandler(service.NewUserService(userRepository))
	adminSyncHandler := synchandler.New(syncRuntime.Repo, syncCfg)

	router := gin.New()
	router.Use(middleware.RequestLogger())
	router.Use(middleware.ErrorRecovery())
	router.Use(middleware.CORS())
	router.Use(middleware.ResolveLocale())

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

	protected := router.Group("")
	protected.Use(middleware.RequireAuth())
	protected.GET("/me/profile", userHandler.GetProfile)
	protected.GET("/me/checkins", userHandler.GetCheckInHistory)
	protected.GET("/matches/:id/my-checkin", checkInHandler.GetMyCheckIn)
	protected.POST("/matches/:id/checkin", checkInHandler.Create)
	protected.PUT("/matches/:id/checkin", checkInHandler.Update)

	admin := router.Group("/admin/sync")
	admin.Use(adminSyncHandler.RequireAdminToken())
	adminSyncHandler.RegisterRoutes(admin)

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"name":    "Final Whistle API",
			"version": "v1",
			"status":  "running",
		})
	})

	return router
}
