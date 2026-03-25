package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"final-whistle/backend/internal/middleware"
	"final-whistle/backend/internal/model"
	"final-whistle/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type fakeAuthService struct {
	loginResult *service.AuthLoginResult
	loginErr    error
	logoutErr   error
	currentUser *model.User
	currentErr  error
}

func (f *fakeAuthService) Login(email, name string) (*service.AuthLoginResult, error) {
	return f.loginResult, f.loginErr
}

func (f *fakeAuthService) Logout(token string) error {
	return f.logoutErr
}

func (f *fakeAuthService) GetCurrentUser(token string) (*model.User, error) {
	return f.currentUser, f.currentErr
}

func TestAuthHandlerLoginSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := NewAuthHandler(&fakeAuthService{
		loginResult: &service.AuthLoginResult{
			User:      &model.User{ID: 1, Name: "Demo User"},
			Session:   &model.Session{Token: "token"},
			ExpiresAt: time.Now().Add(time.Hour),
		},
	}, "development")
	router.POST("/auth/login", handler.Login)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(`{"email":"demo@final-whistle.test","name":"Demo User"}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}
	if len(resp.Result().Cookies()) == 0 {
		t.Fatalf("expected session cookie")
	}
}

func TestAuthHandlerLoginInvalidPayload(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := NewAuthHandler(&fakeAuthService{}, "development")
	router.POST("/auth/login", handler.Login)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(`{"email":"","name":""}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.Code)
	}
}

func TestAuthHandlerMeUnauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := NewAuthHandler(&fakeAuthService{}, "development")
	router.GET("/auth/me", handler.Me)

	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.Code)
	}
}

func TestAuthHandlerLogoutSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := NewAuthHandler(&fakeAuthService{}, "development")
	router.POST("/auth/logout", handler.Logout)

	req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
	req.AddCookie(&http.Cookie{Name: service.SessionCookieName, Value: "token"})
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}
}

func TestAuthHandlerMeSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := NewAuthHandler(&fakeAuthService{}, "development")
	router.GET("/auth/me", func(c *gin.Context) {
		c.Set("currentUser", &model.User{ID: 1, Name: "Demo"})
		handler.Me(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.Code)
	}
}

var _ = middleware.CurrentUser
