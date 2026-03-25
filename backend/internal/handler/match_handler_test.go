package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"final-whistle/backend/internal/dto"
	"final-whistle/backend/internal/repository"
	"final-whistle/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type fakeMatchService struct {
	listResult   dto.MatchListResponseDTO
	listErr      error
	detailResult *dto.MatchDetailDTO
	detailErr    error
}

func (f *fakeMatchService) ListMatches(params repository.MatchListParams) (dto.MatchListResponseDTO, error) {
	return f.listResult, f.listErr
}
func (f *fakeMatchService) GetMatchDetail(id uint) (*dto.MatchDetailDTO, error) {
	return f.detailResult, f.detailErr
}

func TestMatchHandlerListSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := NewMatchHandler(&fakeMatchService{listResult: dto.MatchListResponseDTO{Items: []dto.MatchListItemDTO{}, Page: 1, PageSize: 20, Total: 0}})
	router.GET("/matches", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/matches", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}
}

func TestMatchHandlerDetailNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := NewMatchHandler(&fakeMatchService{detailErr: service.ErrNotFound})
	router.GET("/matches/:id", handler.Detail)

	req := httptest.NewRequest(http.MethodGet, "/matches/1", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.Code)
	}
}
