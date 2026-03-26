package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"final-whistle/backend/internal/dto"
	"final-whistle/backend/internal/middleware"
	"final-whistle/backend/internal/model"
	"final-whistle/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type fakeCheckInService struct {
	getResult    *dto.CheckInDetailDTO
	getErr       error
	createResult *dto.CheckInDetailDTO
	createErr    error
	updateResult *dto.CheckInDetailDTO
	updateErr    error
}

func (f *fakeCheckInService) GetMyCheckIn(matchID, userID uint) (*dto.CheckInDetailDTO, error) {
	return f.getResult, f.getErr
}

func (f *fakeCheckInService) CreateCheckIn(matchID, userID uint, req dto.UpsertCheckInRequestDTO) (*dto.CheckInDetailDTO, error) {
	return f.createResult, f.createErr
}

func (f *fakeCheckInService) UpdateCheckIn(matchID, userID uint, req dto.UpsertCheckInRequestDTO) (*dto.CheckInDetailDTO, error) {
	return f.updateResult, f.updateErr
}

func TestCheckInHandlerGetMyCheckInSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := NewCheckInHandler(&fakeCheckInService{
		getResult: &dto.CheckInDetailDTO{ID: 1, MatchID: 3, WatchedType: "FULL", SupporterSide: "NEUTRAL"},
	})

	protected := router.Group("")
	protected.Use(withCurrentUser(), middleware.RequireAuth())
	protected.GET("/matches/:id/my-checkin", handler.GetMyCheckIn)

	req := httptest.NewRequest(http.MethodGet, "/matches/3/my-checkin", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}
}

func TestCheckInHandlerCreateSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := NewCheckInHandler(&fakeCheckInService{
		createResult: &dto.CheckInDetailDTO{ID: 2, MatchID: 3, WatchedType: "FULL", SupporterSide: "HOME"},
	})

	protected := router.Group("")
	protected.Use(withCurrentUser(), middleware.RequireAuth())
	protected.POST("/matches/:id/checkin", handler.Create)

	body := `{"watchedType":"FULL","supporterSide":"HOME","matchRating":8,"homeTeamRating":8,"awayTeamRating":7,"watchedAt":"2026-03-26T10:00:00Z","tags":[],"playerRatings":[]}`
	req := httptest.NewRequest(http.MethodPost, "/matches/3/checkin", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.Code)
	}
}

func TestCheckInHandlerUpdateSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := NewCheckInHandler(&fakeCheckInService{
		updateResult: &dto.CheckInDetailDTO{ID: 2, MatchID: 3, WatchedType: "PARTIAL", SupporterSide: "AWAY"},
	})

	protected := router.Group("")
	protected.Use(withCurrentUser(), middleware.RequireAuth())
	protected.PUT("/matches/:id/checkin", handler.Update)

	body := `{"watchedType":"PARTIAL","supporterSide":"AWAY","matchRating":7,"homeTeamRating":6,"awayTeamRating":8,"watchedAt":"2026-03-26T11:00:00Z","tags":[1],"playerRatings":[]}`
	req := httptest.NewRequest(http.MethodPut, "/matches/3/checkin", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}
}

func TestCheckInHandlerRejectsUnauthenticatedAccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := NewCheckInHandler(&fakeCheckInService{})

	protected := router.Group("")
	protected.Use(middleware.RequireAuth())
	protected.GET("/matches/:id/my-checkin", handler.GetMyCheckIn)

	req := httptest.NewRequest(http.MethodGet, "/matches/3/my-checkin", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.Code)
	}
}

func TestCheckInHandlerMapsValidationAndConflictErrors(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("validation", func(t *testing.T) {
		router := gin.New()
		handler := NewCheckInHandler(&fakeCheckInService{
			createErr: &service.CheckInValidationError{Message: "invalid payload"},
		})
		protected := router.Group("")
		protected.Use(withCurrentUser(), middleware.RequireAuth())
		protected.POST("/matches/:id/checkin", handler.Create)

		body := `{"watchedType":"FULL","supporterSide":"HOME","matchRating":8,"homeTeamRating":8,"awayTeamRating":7,"watchedAt":"2026-03-26T10:00:00Z","tags":[],"playerRatings":[]}`
		req := httptest.NewRequest(http.MethodPost, "/matches/3/checkin", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", resp.Code)
		}
	})

	t.Run("conflict", func(t *testing.T) {
		router := gin.New()
		handler := NewCheckInHandler(&fakeCheckInService{
			createErr: service.ErrCheckInAlreadyExists,
		})
		protected := router.Group("")
		protected.Use(withCurrentUser(), middleware.RequireAuth())
		protected.POST("/matches/:id/checkin", handler.Create)

		body := `{"watchedType":"FULL","supporterSide":"HOME","matchRating":8,"homeTeamRating":8,"awayTeamRating":7,"watchedAt":"2026-03-26T10:00:00Z","tags":[],"playerRatings":[]}`
		req := httptest.NewRequest(http.MethodPost, "/matches/3/checkin", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusConflict {
			t.Fatalf("expected 409, got %d", resp.Code)
		}
	})
}

func withCurrentUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("currentUser", &model.User{
			ID:    10,
			Name:  "Demo User",
			Email: "demo@final-whistle.test",
		})
		c.Next()
	}
}
