package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/yourorg/go-mvc-app/src/config"
	"github.com/yourorg/go-mvc-app/src/helpers"
	"github.com/yourorg/go-mvc-app/src/repositories"
	"github.com/yourorg/go-mvc-app/src/utils/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type LoginInput struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type TokenResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type AuthService interface {
	Login(input LoginInput) (*TokenResponse, error)
	IssueToken(userID, role string) (*TokenResponse, error)
}

type authService struct {
	repo repositories.UserRepository
}

func NewAuthService(repo repositories.UserRepository) AuthService {
	return &authService{repo: repo}
}

func (s *authService) Login(input LoginInput) (*TokenResponse, error) {
	user, err := s.repo.FindByEmail(helpers.SanitizeString(input.Email))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(config.MsgInvalidCredentials)
		}
		logger.Error("AuthService.Login db error", zap.Error(err))
		return nil, err
	}

	if !user.IsActive {
		return nil, errors.New(config.MsgInvalidCredentials)
	}

	if !helpers.CheckPassword(user.Password, input.Password) {
		return nil, errors.New(config.MsgInvalidCredentials)
	}

	expiresAt := time.Now().Add(time.Duration(config.App.JWT.ExpiryHours) * time.Hour)
	claims := &Claims{
		UserID: user.ID.String(),
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.ID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(config.App.JWT.Secret))
	if err != nil {
		logger.Error("AuthService.Login token signing failed", zap.Error(err))
		return nil, err
	}

	logger.Info("AuthService.Login success", zap.String("email", input.Email))
	return &TokenResponse{Token: tokenStr, ExpiresAt: expiresAt}, nil
}

// IssueToken generates a signed JWT for an already-authenticated user (e.g. post-signup).
func (s *authService) IssueToken(userID, role string) (*TokenResponse, error) {
	expiresAt := time.Now().Add(time.Duration(config.App.JWT.ExpiryHours) * time.Hour)
	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   userID,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(config.App.JWT.Secret))
	if err != nil {
		return nil, err
	}
	return &TokenResponse{Token: tokenStr, ExpiresAt: expiresAt}, nil
}

// ValidateJWT parses and validates a JWT token string, returning the claims.
func ValidateJWT(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(config.App.JWT.Secret), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New(config.MsgTokenInvalid)
	}
	return claims, nil
}
