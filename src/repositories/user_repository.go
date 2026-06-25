package repositories

import (
	"github.com/google/uuid"
	"github.com/yourorg/go-mvc-app/src/models"
	"github.com/yourorg/go-mvc-app/src/helpers"
	"github.com/yourorg/go-mvc-app/src/utils/database"
	"github.com/yourorg/go-mvc-app/src/utils/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	FindByID(id uuid.UUID) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	FindAll(offset, limit int) ([]models.User, int64, error)
	Update(user *models.User) error
	Delete(id uuid.UUID) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository() UserRepository {
	return &userRepository{db: database.DB}
}

func (r *userRepository) Create(user *models.User) (err error) {
	err = helpers.Try(func() error {
		return r.db.Create(user).Error
	})
	if err != nil {
		logger.Error("UserRepository.Create failed", zap.Error(err), zap.String("email", user.Email))
	}
	return
}

func (r *userRepository) FindByID(id uuid.UUID) (user *models.User, err error) {
	err = helpers.Try(func() error {
		user = &models.User{}
		return r.db.First(user, "id = ?", id).Error
	})
	if err != nil {
		logger.Error("UserRepository.FindByID failed", zap.Error(err), zap.String("id", id.String()))
		return nil, err
	}
	return
}

func (r *userRepository) FindByEmail(email string) (user *models.User, err error) {
	err = helpers.Try(func() error {
		user = &models.User{}
		return r.db.Where("email = ?", email).First(user).Error
	})
	if err != nil {
		logger.Debug("UserRepository.FindByEmail", zap.Error(err), zap.String("email", email))
		return nil, err
	}
	return
}

func (r *userRepository) FindAll(offset, limit int) (users []models.User, count int64, err error) {
	err = helpers.Try(func() error {
		if e := r.db.Model(&models.User{}).Count(&count).Error; e != nil {
			return e
		}
		return r.db.Offset(offset).Limit(limit).Find(&users).Error
	})
	if err != nil {
		logger.Error("UserRepository.FindAll failed", zap.Error(err))
	}
	return
}

func (r *userRepository) Update(user *models.User) (err error) {
	err = helpers.Try(func() error {
		return r.db.Save(user).Error
	})
	if err != nil {
		logger.Error("UserRepository.Update failed", zap.Error(err), zap.String("id", user.ID.String()))
	}
	return
}

func (r *userRepository) Delete(id uuid.UUID) (err error) {
	err = helpers.Try(func() error {
		return r.db.Delete(&models.User{}, "id = ?", id).Error
	})
	if err != nil {
		logger.Error("UserRepository.Delete failed", zap.Error(err), zap.String("id", id.String()))
	}
	return
}
