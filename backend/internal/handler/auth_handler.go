package handler

import (
	"net/http"
	"net/mail"
	"strings"
	"time"

	"final-whistle/backend/internal/dto"
	"final-whistle/backend/internal/middleware"
	"final-whistle/backend/internal/service"
	"final-whistle/backend/internal/utils"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service service.AuthService
	env     string
}

func NewAuthHandler(service service.AuthService, env string) *AuthHandler {
	return &AuthHandler{service: service, env: env}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "invalid login payload", nil)
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	req.Name = strings.TrimSpace(req.Name)
	if req.Email == "" || req.Name == "" {
		utils.ValidationErrorResponse(c, "email and name are required", nil)
		return
	}
	if _, err := mail.ParseAddress(req.Email); err != nil {
		utils.ValidationErrorResponse(c, "invalid email address", nil)
		return
	}

	result, err := h.service.Login(req.Email, req.Name)
	if err != nil {
		if err == service.ErrUnauthorized {
			utils.UnauthorizedResponse(c, "Login is not allowed for this user")
			return
		}
		utils.InternalErrorResponse(c, "Failed to login", nil)
		return
	}

	h.setSessionCookie(c, result.Session.Token, result.ExpiresAt)
	utils.OKResponse(c, dto.AuthUserResponseDTO{
		User: dto.UserSummaryDTO{
			ID:        result.User.ID,
			Name:      result.User.Name,
			AvatarURL: result.User.AvatarURL,
		},
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	token, _ := c.Cookie(service.SessionCookieName)
	if err := h.service.Logout(token); err != nil {
		utils.InternalErrorResponse(c, "Failed to logout", nil)
		return
	}

	h.clearSessionCookie(c)
	utils.OKResponse(c, dto.LogoutResponseDTO{OK: true})
}

func (h *AuthHandler) Me(c *gin.Context) {
	user, ok := middleware.CurrentUser(c)
	if !ok {
		if _, err := c.Cookie(service.SessionCookieName); err == nil {
			h.clearSessionCookie(c)
		}
		utils.UnauthorizedResponse(c, "Authentication required")
		return
	}

	utils.OKResponse(c, dto.AuthUserResponseDTO{
		User: dto.UserSummaryDTO{
			ID:        user.ID,
			Name:      user.Name,
			AvatarURL: user.AvatarURL,
		},
	})
}

func (h *AuthHandler) setSessionCookie(c *gin.Context, token string, expiresAt time.Time) {
	maxAge := int(time.Until(expiresAt).Seconds())
	if maxAge < 0 {
		maxAge = 0
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		service.SessionCookieName,
		token,
		maxAge,
		"/",
		"",
		h.env == "production",
		true,
	)
}

func (h *AuthHandler) clearSessionCookie(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		service.SessionCookieName,
		"",
		-1,
		"/",
		"",
		h.env == "production",
		true,
	)
}
