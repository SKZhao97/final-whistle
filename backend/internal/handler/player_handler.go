package handler

import (
	"final-whistle/backend/internal/service"
	"final-whistle/backend/internal/utils"
	"github.com/gin-gonic/gin"
)

type PlayerHandler struct {
	service service.PlayerService
}

func NewPlayerHandler(service service.PlayerService) *PlayerHandler {
	return &PlayerHandler{service: service}
}

func (h *PlayerHandler) Detail(c *gin.Context) {
	id := parseInt(c.Param("id"), 0)
	if id <= 0 {
		utils.ValidationErrorResponse(c, "invalid player id", nil)
		return
	}
	result, err := h.service.GetPlayerDetail(uint(id))
	if err != nil {
		if err == service.ErrNotFound {
			utils.NotFoundResponse(c, "Player not found")
			return
		}
		utils.InternalErrorResponse(c, "Failed to load player", nil)
		return
	}
	utils.OKResponse(c, result)
}
