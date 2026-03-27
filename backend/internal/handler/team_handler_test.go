package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"final-whistle/backend/internal/dto"
	"final-whistle/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type fakeTeamService struct {
	result *dto.TeamDetailDTO
	err    error
}

func (f *fakeTeamService) GetTeamDetail(id uint, locale string) (*dto.TeamDetailDTO, error) {
	return f.result, f.err
}

func TestTeamHandlerDetailSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := NewTeamHandler(&fakeTeamService{result: &dto.TeamDetailDTO{ID: 1, Name: "Arsenal", Slug: "arsenal"}})
	router.GET("/teams/:id", handler.Detail)

	req := httptest.NewRequest(http.MethodGet, "/teams/1", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}
}

func TestTeamHandlerDetailNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := NewTeamHandler(&fakeTeamService{err: service.ErrNotFound})
	router.GET("/teams/:id", handler.Detail)

	req := httptest.NewRequest(http.MethodGet, "/teams/999", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.Code)
	}
}
