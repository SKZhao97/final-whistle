package service

import (
	"errors"
	"testing"
	"time"

	"final-whistle/backend/internal/model"
	"gorm.io/gorm"
)

type fakeAuthRepository struct {
	userByEmail      map[string]*model.User
	sessionByToken   map[string]*model.Session
	findUserErr      error
	createUserErr    error
	createSessionErr error
	findSessionErr   error
	deleteTokenErr   error
	deleteIDErr      error
	deletedSessionID uint
	deletedToken     string
	createdUsers     []*model.User
	createdSessions  []*model.Session
	nextUserID       uint
}

func (f *fakeAuthRepository) FindUserByEmail(email string) (*model.User, error) {
	if f.findUserErr != nil {
		return nil, f.findUserErr
	}
	user, ok := f.userByEmail[email]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return user, nil
}

func (f *fakeAuthRepository) CreateUser(user *model.User) error {
	if f.createUserErr != nil {
		return f.createUserErr
	}
	if f.nextUserID == 0 {
		f.nextUserID = 100
	}
	user.ID = f.nextUserID
	f.nextUserID++
	f.createdUsers = append(f.createdUsers, user)
	if f.userByEmail == nil {
		f.userByEmail = map[string]*model.User{}
	}
	f.userByEmail[user.Email] = user
	return nil
}

func (f *fakeAuthRepository) CreateSession(session *model.Session) error {
	if f.createSessionErr != nil {
		return f.createSessionErr
	}
	f.createdSessions = append(f.createdSessions, session)
	if f.sessionByToken == nil {
		f.sessionByToken = map[string]*model.Session{}
	}
	f.sessionByToken[session.Token] = session
	return nil
}

func (f *fakeAuthRepository) FindSessionByToken(token string) (*model.Session, error) {
	if f.findSessionErr != nil {
		return nil, f.findSessionErr
	}
	session, ok := f.sessionByToken[token]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return session, nil
}

func (f *fakeAuthRepository) DeleteSessionByToken(token string) error {
	if f.deleteTokenErr != nil {
		return f.deleteTokenErr
	}
	f.deletedToken = token
	delete(f.sessionByToken, token)
	return nil
}

func (f *fakeAuthRepository) DeleteSessionByID(id uint) error {
	if f.deleteIDErr != nil {
		return f.deleteIDErr
	}
	f.deletedSessionID = id
	return nil
}

func TestAuthServiceLoginExistingUser(t *testing.T) {
	repo := &fakeAuthRepository{
		userByEmail: map[string]*model.User{
			"demo@final-whistle.test": {ID: 1, Name: "Demo User", Email: "demo@final-whistle.test"},
		},
	}
	svc := NewAuthService(repo, true)

	result, err := svc.Login("demo@final-whistle.test", "Demo User")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.User.ID != 1 {
		t.Fatalf("expected existing user, got %#v", result.User)
	}
	if len(repo.createdUsers) != 0 {
		t.Fatalf("did not expect user creation")
	}
	if len(repo.createdSessions) != 1 {
		t.Fatalf("expected session creation")
	}
}

func TestAuthServiceLoginAutoCreatesUserInDevelopment(t *testing.T) {
	repo := &fakeAuthRepository{}
	svc := NewAuthService(repo, true)

	result, err := svc.Login("new@final-whistle.test", "New User")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.User.ID == 0 || len(repo.createdUsers) != 1 {
		t.Fatalf("expected user creation, got %#v", result.User)
	}
}

func TestAuthServiceLoginRejectsUnknownUserWhenAutoCreateDisabled(t *testing.T) {
	repo := &fakeAuthRepository{}
	svc := NewAuthService(repo, false)

	_, err := svc.Login("new@final-whistle.test", "New User")
	if !errors.Is(err, ErrUnauthorized) {
		t.Fatalf("expected ErrUnauthorized, got %v", err)
	}
}

func TestAuthServiceGetCurrentUserUnauthorizedWithoutToken(t *testing.T) {
	svc := NewAuthService(&fakeAuthRepository{}, true)
	_, err := svc.GetCurrentUser("")
	if !errors.Is(err, ErrUnauthorized) {
		t.Fatalf("expected ErrUnauthorized, got %v", err)
	}
}

func TestAuthServiceGetCurrentUserDeletesExpiredSession(t *testing.T) {
	repo := &fakeAuthRepository{
		sessionByToken: map[string]*model.Session{
			"expired": {
				ID:        5,
				Token:     "expired",
				ExpiredAt: time.Now().Add(-time.Hour),
				User:      model.User{ID: 1, Name: "Demo"},
			},
		},
	}
	svc := NewAuthService(repo, true)

	_, err := svc.GetCurrentUser("expired")
	if !errors.Is(err, ErrUnauthorized) {
		t.Fatalf("expected ErrUnauthorized, got %v", err)
	}
	if repo.deletedSessionID != 5 {
		t.Fatalf("expected expired session cleanup, got %d", repo.deletedSessionID)
	}
}

func TestAuthServiceLogoutIgnoresMissingToken(t *testing.T) {
	repo := &fakeAuthRepository{}
	svc := NewAuthService(repo, true)

	if err := svc.Logout(""); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
