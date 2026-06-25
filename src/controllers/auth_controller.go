package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/yourorg/go-mvc-app/src/config"
	"github.com/yourorg/go-mvc-app/src/helpers"
	"github.com/yourorg/go-mvc-app/src/services"
	"go.uber.org/zap"
)



type AuthController struct {
	service  services.AuthService
	userSvc  services.UserService
}

func NewAuthController(svc services.AuthService, userSvc services.UserService) *AuthController {
	return &AuthController{service: svc, userSvc: userSvc}
}

// Login godoc
// @Summary      Login
// @Description  Authenticate with email and password, receive a JWT bearer token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      services.LoginInput                        true  "Login credentials"
// @Success      200   {object}  helpers.Response{data=services.TokenResponse}
// @Failure      400   {object}  helpers.Response
// @Failure      401   {object}  helpers.Response
// @Failure      500   {object}  helpers.Response
// @Router       /auth/login [post]
func (ac *AuthController) Login(c *gin.Context) {
	log := requestLog(c)
	var input services.LoginInput

	if err := helpers.Try(func() error {
		if bindErr := c.ShouldBindJSON(&input); bindErr != nil {
			helpers.BadRequest(c, bindErr.Error())
			return bindErr
		}
		resp, svcErr := ac.service.Login(input)
		if svcErr != nil {
			if svcErr.Error() == config.MsgInvalidCredentials {
				helpers.Unauthorized(c)
			} else {
				helpers.InternalError(c, svcErr.Error())
			}
			return svcErr
		}
		log.Info("login success", zap.String("email", input.Email))
		helpers.OK(c, config.MsgLoginSuccess, resp)
		return nil
	}); err != nil {
		log.Error("AuthController.Login", zap.Error(err))
	}
}

// Signup godoc
// @Summary      Sign up
// @Description  Register a new account and receive a JWT token immediately
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      services.CreateUserInput  true  "Registration payload"
// @Success      201   {object}  helpers.Response
// @Failure      400   {object}  helpers.Response
// @Failure      409   {object}  helpers.Response
// @Failure      500   {object}  helpers.Response
// @Router       /auth/signup [post]
func (ac *AuthController) Signup(c *gin.Context) {
	log := requestLog(c)
	var input services.CreateUserInput

	if err := helpers.Try(func() error {
		if bindErr := c.ShouldBindJSON(&input); bindErr != nil {
			helpers.BadRequest(c, bindErr.Error())
			return bindErr
		}

		user, svcErr := ac.userSvc.Create(input)
		if svcErr != nil {
			if svcErr.Error() == "email already registered" {
				helpers.Conflict(c, svcErr.Error())
			} else {
				helpers.InternalError(c, svcErr.Error())
			}
			return svcErr
		}

		tokenResp, tokenErr := ac.service.IssueToken(user.ID.String(), user.Role)
		if tokenErr != nil {
			helpers.InternalError(c, tokenErr.Error())
			return tokenErr
		}

		log.Info("signup success", zap.String("email", user.Email), zap.String("id", user.ID.String()))
		helpers.Created(c, gin.H{
			"token":      tokenResp.Token,
			"expires_at": tokenResp.ExpiresAt,
			"user":       user,
		})
		return nil
	}); err != nil {
		log.Error("AuthController.Signup", zap.Error(err))
	}
}
