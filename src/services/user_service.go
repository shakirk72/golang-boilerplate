package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/yourorg/go-mvc-app/src/helpers"
	"github.com/yourorg/go-mvc-app/src/models"
	"github.com/yourorg/go-mvc-app/src/repositories"
	"github.com/yourorg/go-mvc-app/src/utils/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CreateUserInput struct {
	Name     string `json:"name"     binding:"required,min=2,max=100"`
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Role     string `json:"role"`
}

type UpdateUserInput struct {
	Name     string `json:"name"      binding:"omitempty,min=2,max=100"`
	Role     string `json:"role"      binding:"omitempty"`
	IsActive *bool  `json:"is_active"`
}

type UserService interface {
	Create(input CreateUserInput) (*models.User, error)
	GetByID(id string) (*models.User, error)
	List(page, pageSize int) ([]models.User, int64, error)
	Update(id string, input UpdateUserInput) (*models.User, error)
	Delete(id string) error
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Create(input CreateUserInput) (user *models.User, err error) {
	err = helpers.Try(func() error {
		existing, findErr := s.repo.FindByEmail(input.Email)
		if findErr != nil && !errors.Is(findErr, gorm.ErrRecordNotFound) {
			return findErr
		}
		if existing != nil {
			return errors.New("email already registered")
		}

		hashed, hashErr := helpers.HashPassword(input.Password)
		if hashErr != nil {
			return hashErr
		}

		role := input.Role
		if role == "" {
			role = "user"
		}

		user = &models.User{
			Name:     input.Name,
			Email:    helpers.SanitizeString(input.Email),
			Password: hashed,
			Role:     role,
		}
		return s.repo.Create(user)
	})
	if err != nil {
		logger.Error("UserService.Create failed", zap.Error(err), zap.String("email", input.Email))
		return nil, err
	}
	logger.Info("UserService.Create success", zap.String("email", user.Email), zap.String("id", user.ID.String()))
	return
}

func (s *userService) GetByID(id string) (user *models.User, err error) {
	err = helpers.Try(func() error {
		uid, parseErr := uuid.Parse(id)
		if parseErr != nil {
			return errors.New("invalid user id")
		}
		var fetchErr error
		user, fetchErr = s.repo.FindByID(uid)
		return fetchErr
	})
	if err != nil {
		logger.Warn("UserService.GetByID failed", zap.Error(err), zap.String("id", id))
		return nil, err
	}
	return
}

func (s *userService) List(page, pageSize int) (users []models.User, total int64, err error) {
	err = helpers.Try(func() error {
		offset := helpers.Offset(page, pageSize)
		var fetchErr error
		users, total, fetchErr = s.repo.FindAll(offset, pageSize)
		return fetchErr
	})
	if err != nil {
		logger.Error("UserService.List failed", zap.Error(err), zap.Int("page", page))
	}
	return
}

func (s *userService) Update(id string, input UpdateUserInput) (user *models.User, err error) {
	err = helpers.Try(func() error {
		uid, parseErr := uuid.Parse(id)
		if parseErr != nil {
			return errors.New("invalid user id")
		}

		var fetchErr error
		user, fetchErr = s.repo.FindByID(uid)
		if fetchErr != nil {
			return fetchErr
		}

		if input.Name != "" {
			user.Name = input.Name
		}
		if input.Role != "" {
			user.Role = input.Role
		}
		if input.IsActive != nil {
			user.IsActive = *input.IsActive
		}
		return s.repo.Update(user)
	})
	if err != nil {
		logger.Error("UserService.Update failed", zap.Error(err), zap.String("id", id))
		return nil, err
	}
	logger.Info("UserService.Update success", zap.String("id", id))
	return
}

func (s *userService) Delete(id string) (err error) {
	err = helpers.Try(func() error {
		uid, parseErr := uuid.Parse(id)
		if parseErr != nil {
			return errors.New("invalid user id")
		}
		return s.repo.Delete(uid)
	})
	if err != nil {
		logger.Error("UserService.Delete failed", zap.Error(err), zap.String("id", id))
	} else {
		logger.Info("UserService.Delete success", zap.String("id", id))
	}
	return
}
