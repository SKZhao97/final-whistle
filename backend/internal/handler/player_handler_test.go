package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"final-whistle/backend/internal/dto"
	"final-whistle/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type fakePlayerService struct {
	result *dto.PlayerDetailDTO
	err    error
}

func (f *fakePlayerService) GetPlayerDetail(id uint, locale string) (*dto.PlayerDetailDTO, error) {
	return f.result, f.err
}

func TestPlayerHandlerDetailSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := NewPlayerHandler(&fakePlayerService{result: &dto.PlayerDetailDTO{ID: 1, Name: "Player", Slug: "player"}})
	router.GET("/players/:id", handler.Detail)

	req := httptest.NewRequest(http.MethodGet, "/players/1", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}
}

func TestPlayerHandlerDetailNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := NewPlayerHandler(&fakePlayerService{err: service.ErrNotFound})
	router.GET("/players/:id", handler.Detail)

	req := httptest.NewRequest(http.MethodGet, "/players/999", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.Code)
	}
}
