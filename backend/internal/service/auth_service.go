package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"strings"
	"time"

	"final-whistle/backend/internal/model"
	"final-whistle/backend/internal/repository"
	"gorm.io/gorm"
)

const (
	SessionCookieName = "final_whistle_session"
	SessionDuration   = 7 * 24 * time.Hour
)

var ErrUnauthorized = errors.New("unauthorized")

type AuthLoginResult struct {
	User      *model.User
	Session   *model.Session
	Created   bool
	ExpiresAt time.Time
}

type AuthService interface {
	Login(email, name string) (*AuthLoginResult, error)
	Logout(token string) error
	GetCurrentUser(token string) (*model.User, error)
}

type authService struct {
	repo            repository.AuthRepository
	allowAutoCreate bool
	now             func() time.Time
}

func NewAuthService(repo repository.AuthRepository, allowAutoCreate bool) AuthService {
	return &authService{
		repo:            repo,
		allowAutoCreate: allowAutoCreate,
		now:             time.Now,
	}
}

func (s *authService) Login(email, name string) (*AuthLoginResult, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	name = strings.TrimSpace(name)

	user, err := s.repo.FindUserByEmail(email)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		if !s.allowAutoCreate {
			return nil, ErrUnauthorized
		}

		user = &model.User{
			Email: email,
			Name:  name,
		}
		if err := s.repo.CreateUser(user); err != nil {
			return nil, err
		}
	}

	token, err := generateSessionToken()
	if err != nil {
		return nil, err
	}

	expiresAt := s.now().Add(SessionDuration)
	session := &model.Session{
		UserID:    user.ID,
		Token:     token,
		ExpiredAt: expiresAt,
	}
	if err := s.repo.CreateSession(session); err != nil {
		return nil, err
	}

	return &AuthLoginResult{
		User:      user,
		Session:   session,
		ExpiresAt: expiresAt,
	}, nil
}

func (s *authService) Logout(token string) error {
	token = strings.TrimSpace(token)
	if token == "" {
		return nil
	}
	return s.repo.DeleteSessionByToken(token)
}

func (s *authService) GetCurrentUser(token string) (*model.User, error) {
	token = strings.TrimSpace(token)
	if token == "" {
		return nil, ErrUnauthorized
	}

	session, err := s.repo.FindSessionByToken(token)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUnauthorized
		}
		return nil, err
	}

	if session.ExpiredAt.Before(s.now()) {
		if err := s.repo.DeleteSessionByID(session.ID); err != nil {
			log.Printf("failed to delete expired session %d: %v", session.ID, err)
		}
		return nil, ErrUnauthorized
	}

	return &session.User, nil
}

func generateSessionToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}
