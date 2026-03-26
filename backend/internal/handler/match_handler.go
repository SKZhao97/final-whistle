package handler

import (
	"final-whistle/backend/internal/repository"
	"final-whistle/backend/internal/service"
	"final-whistle/backend/internal/utils"
	"github.com/gin-gonic/gin"
)

type MatchHandler struct {
	service service.MatchService
}

func NewMatchHandler(service service.MatchService) *MatchHandler {
	return &MatchHandler{service: service}
}

func (h *MatchHandler) List(c *gin.Context) {
	page := parseInt(c.Query("page"), 1)
	pageSize := parseInt(c.Query("pageSize"), 20)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 50 {
		pageSize = 50
	}

	result, err := h.service.ListMatches(repository.MatchListParams{
		Competition: c.Query("competition"),
		Season:      c.Query("season"),
		Page:        page,
		PageSize:    pageSize,
	})
	if err != nil {
		utils.InternalErrorResponse(c, "Failed to load matches", nil)
		return
	}
	utils.OKResponse(c, result)
}

func (h *MatchHandler) Detail(c *gin.Context) {
	id := parseInt(c.Param("id"), 0)
	if id <= 0 {
		utils.ValidationErrorResponse(c, "invalid match id", nil)
		return
	}
	result, err := h.service.GetMatchDetail(uint(id))
	if err != nil {
		if err == service.ErrNotFound {
			utils.NotFoundResponse(c, "Match not found")
			return
		}
		utils.InternalErrorResponse(c, "Failed to load match", nil)
		return
	}
	utils.OKResponse(c, result)
}
