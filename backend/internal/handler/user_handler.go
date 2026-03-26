package handler

import (
	"final-whistle/backend/internal/middleware"
	"final-whistle/backend/internal/service"
	"final-whistle/backend/internal/utils"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	user, ok := middleware.CurrentUser(c)
	if !ok {
		utils.UnauthorizedResponse(c, "Authentication required")
		return
	}

	result, err := h.service.GetProfileSummary(user.ID)
	if err != nil {
		if err == service.ErrNotFound {
			utils.NotFoundResponse(c, "User not found")
			return
		}
		utils.InternalErrorResponse(c, "Failed to load profile", nil)
		return
	}

	utils.OKResponse(c, result)
}

func (h *UserHandler) GetCheckInHistory(c *gin.Context) {
	user, ok := middleware.CurrentUser(c)
	if !ok {
		utils.UnauthorizedResponse(c, "Authentication required")
		return
	}

	page := parseInt(c.Query("page"), 1)
	pageSize := parseInt(c.Query("pageSize"), 20)

	result, err := h.service.GetCheckInHistory(user.ID, page, pageSize)
	if err != nil {
		utils.InternalErrorResponse(c, "Failed to load check-in history", nil)
		return
	}

	utils.OKResponse(c, result)
}
