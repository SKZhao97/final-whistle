package handler

import (
	"final-whistle/backend/internal/service"
	"final-whistle/backend/internal/utils"
	"github.com/gin-gonic/gin"
)

type TeamHandler struct {
	service service.TeamService
}

func NewTeamHandler(service service.TeamService) *TeamHandler {
	return &TeamHandler{service: service}
}

func (h *TeamHandler) Detail(c *gin.Context) {
	id := parseInt(c.Param("id"), 0)
	if id <= 0 {
		utils.ValidationErrorResponse(c, "invalid team id", nil)
		return
	}
	result, err := h.service.GetTeamDetail(uint(id))
	if err != nil {
		if err == service.ErrNotFound {
			utils.NotFoundResponse(c, "Team not found")
			return
		}
		utils.InternalErrorResponse(c, "Failed to load team", nil)
		return
	}
	utils.OKResponse(c, result)
}
