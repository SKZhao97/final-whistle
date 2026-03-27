package handler

import (
	"errors"
	"net/http"

	"final-whistle/backend/internal/dto"
	"final-whistle/backend/internal/middleware"
	"final-whistle/backend/internal/service"
	"final-whistle/backend/internal/utils"
	"github.com/gin-gonic/gin"
)

type CheckInHandler struct {
	service service.CheckInService
}

func NewCheckInHandler(service service.CheckInService) *CheckInHandler {
	return &CheckInHandler{service: service}
}

func (h *CheckInHandler) GetMyCheckIn(c *gin.Context) {
	user, ok := middleware.CurrentUser(c)
	if !ok {
		utils.UnauthorizedResponse(c, "Authentication required")
		return
	}

	matchID := parseInt(c.Param("id"), 0)
	if matchID <= 0 {
		utils.ValidationErrorResponse(c, "invalid match id", nil)
		return
	}

	result, err := h.service.GetMyCheckIn(uint(matchID), user.ID, middleware.CurrentLocale(c))
	if err != nil {
		h.writeError(c, err)
		return
	}

	utils.OKResponse(c, result)
}

func (h *CheckInHandler) Create(c *gin.Context) {
	h.upsert(c, http.StatusCreated, h.service.CreateCheckIn)
}

func (h *CheckInHandler) Update(c *gin.Context) {
	h.upsert(c, http.StatusOK, h.service.UpdateCheckIn)
}

func (h *CheckInHandler) upsert(c *gin.Context, successStatus int, fn func(matchID, userID uint, req dto.UpsertCheckInRequestDTO, locale string) (*dto.CheckInDetailDTO, error)) {
	user, ok := middleware.CurrentUser(c)
	if !ok {
		utils.UnauthorizedResponse(c, "Authentication required")
		return
	}

	matchID := parseInt(c.Param("id"), 0)
	if matchID <= 0 {
		utils.ValidationErrorResponse(c, "invalid match id", nil)
		return
	}

	var req dto.UpsertCheckInRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "invalid check-in payload", nil)
		return
	}

	result, err := fn(uint(matchID), user.ID, req, middleware.CurrentLocale(c))
	if err != nil {
		h.writeError(c, err)
		return
	}

	utils.SuccessResponse(c, successStatus, result)
}

func (h *CheckInHandler) writeError(c *gin.Context, err error) {
	var validationErr *service.CheckInValidationError
	switch {
	case errors.As(err, &validationErr):
		utils.ValidationErrorResponse(c, validationErr.Message, validationErr.Details)
	case errors.Is(err, service.ErrNotFound):
		utils.NotFoundResponse(c, "Match not found")
	case errors.Is(err, service.ErrCheckInAlreadyExists):
		utils.ConflictResponse(c, "Check-in already exists for this match", nil)
	case errors.Is(err, service.ErrCheckInMissing):
		utils.NotFoundResponse(c, "Check-in not found for this match")
	default:
		utils.InternalErrorResponse(c, "Failed to process check-in", nil)
	}
}
