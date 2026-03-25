package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"final-whistle/backend/internal/model"
	"final-whistle/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type fakeAuthResolver struct {
	user *model.User
	err  error
}

func (f *fakeAuthResolver) Login(email, name string) (*service.AuthLoginResult, error) {
	return nil, nil
}

func (f *fakeAuthResolver) Logout(token string) error {
	return nil
}

func (f *fakeAuthResolver) GetCurrentUser(token string) (*model.User, error) {
	return f.user, f.err
}

func TestResolveCurrentUserSetsContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(ResolveCurrentUser(&fakeAuthResolver{user: &model.User{ID: 1, Name: "Demo"}}))
	router.GET("/test", func(c *gin.Context) {
		if user, ok := CurrentUser(c); !ok || user.ID != 1 {
			c.Status(http.StatusUnauthorized)
			return
		}
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.AddCookie(&http.Cookie{Name: service.SessionCookieName, Value: "token"})
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}
}

func TestRequireAuthRejectsMissingUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/protected", RequireAuth(), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.Code)
	}
}

func TestResolveCurrentUserIgnoresUnauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(ResolveCurrentUser(&fakeAuthResolver{err: errors.New("unauthorized")}))
	router.GET("/test", func(c *gin.Context) {
		_, ok := CurrentUser(c)
		if ok {
			c.Status(http.StatusBadRequest)
			return
		}
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.AddCookie(&http.Cookie{Name: service.SessionCookieName, Value: "token"})
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}
}
