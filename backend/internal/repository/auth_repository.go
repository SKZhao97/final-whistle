package repository

import (
	"errors"

	"final-whistle/backend/internal/model"
	"gorm.io/gorm"
)

type AuthRepository interface {
	FindUserByEmail(email string) (*model.User, error)
	CreateUser(user *model.User) error
	CreateSession(session *model.Session) error
	FindSessionByToken(token string) (*model.Session, error)
	DeleteSessionByToken(token string) error
	DeleteSessionByID(id uint) error
}

type GormAuthRepository struct {
	*BaseRepository
}

func NewAuthRepository(db *gorm.DB) *GormAuthRepository {
	return &GormAuthRepository{BaseRepository: NewBaseRepository(db)}
}

func (r *GormAuthRepository) FindUserByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GormAuthRepository) CreateUser(user *model.User) error {
	return r.DB.Create(user).Error
}

func (r *GormAuthRepository) CreateSession(session *model.Session) error {
	return r.DB.Create(session).Error
}

func (r *GormAuthRepository) FindSessionByToken(token string) (*model.Session, error) {
	var session model.Session
	if err := r.DB.Preload("User").Where("token = ?", token).First(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *GormAuthRepository) DeleteSessionByToken(token string) error {
	result := r.DB.Where("token = ?", token).Delete(&model.Session{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *GormAuthRepository) DeleteSessionByID(id uint) error {
	result := r.DB.Delete(&model.Session{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func IsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}
