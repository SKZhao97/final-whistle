package middleware

import (
	"final-whistle/backend/internal/model"
	"final-whistle/backend/internal/service"
	"final-whistle/backend/internal/utils"
	"github.com/gin-gonic/gin"
)

const currentUserKey = "currentUser"

func ResolveCurrentUser(authService service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie(service.SessionCookieName)
		if err != nil || token == "" {
			c.Next()
			return
		}

		user, err := authService.GetCurrentUser(token)
		if err == nil && user != nil {
			c.Set(currentUserKey, user)
		}

		c.Next()
	}
}

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, ok := CurrentUser(c); !ok {
			utils.UnauthorizedResponse(c, "Authentication required")
			c.Abort()
			return
		}
		c.Next()
	}
}

func CurrentUser(c *gin.Context) (*model.User, bool) {
	value, exists := c.Get(currentUserKey)
	if !exists {
		return nil, false
	}

	user, ok := value.(*model.User)
	if !ok || user == nil {
		return nil, false
	}
	return user, true
}
