package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"final-whistle/backend/internal/dto"
	"final-whistle/backend/internal/model"
	"final-whistle/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type fakeUserService struct {
	profile    *dto.UserProfileSummaryDTO
	profileErr error
	history    *dto.UserCheckInHistoryResponseDTO
	historyErr error
}

func (f *fakeUserService) GetProfileSummary(userID uint) (*dto.UserProfileSummaryDTO, error) {
	return f.profile, f.profileErr
}

func (f *fakeUserService) GetCheckInHistory(userID uint, page, pageSize int) (*dto.UserCheckInHistoryResponseDTO, error) {
	return f.history, f.historyErr
}

func TestUserHandlerGetProfileUnauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := NewUserHandler(&fakeUserService{})
	router.GET("/me/profile", handler.GetProfile)

	req := httptest.NewRequest(http.MethodGet, "/me/profile", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.Code)
	}
}

func TestUserHandlerGetProfileSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := NewUserHandler(&fakeUserService{
		profile: &dto.UserProfileSummaryDTO{
			User:         dto.UserSummaryDTO{ID: 1, Name: "Demo User"},
			CheckInCount: 3,
		},
	})
	router.GET("/me/profile", func(c *gin.Context) {
		c.Set("currentUser", &model.User{ID: 1, Name: "Demo User"})
		handler.GetProfile(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/me/profile", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}
}

func TestUserHandlerGetCheckInHistorySuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := NewUserHandler(&fakeUserService{
		history: &dto.UserCheckInHistoryResponseDTO{
			Items:    []dto.UserCheckInHistoryItemDTO{},
			Page:     2,
			PageSize: 10,
			Total:    0,
		},
	})
	router.GET("/me/checkins", func(c *gin.Context) {
		c.Set("currentUser", &model.User{ID: 1, Name: "Demo User"})
		handler.GetCheckInHistory(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/me/checkins?page=2&pageSize=10", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}
}

var _ service.UserService = (*fakeUserService)(nil)
