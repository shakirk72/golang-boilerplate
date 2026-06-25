package controllers

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/yourorg/go-mvc-app/src/config"
	"github.com/yourorg/go-mvc-app/src/helpers"
	"github.com/yourorg/go-mvc-app/src/services"
	"github.com/yourorg/go-mvc-app/src/utils/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserController struct {
	service services.UserService
}

func NewUserController(svc services.UserService) *UserController {
	return &UserController{service: svc}
}

func requestLog(c *gin.Context) *logger.Logger {
	return logger.WithRequestID(c.GetString(config.CtxRequestID))
}

// Create godoc
// @Summary      Create a user
// @Description  Register a new user account
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        body  body      services.CreateUserInput  true  "User payload"
// @Success      201   {object}  helpers.Response{data=models.User}
// @Failure      400   {object}  helpers.Response
// @Failure      409   {object}  helpers.Response
// @Failure      500   {object}  helpers.Response
// @Router       /users [post]
func (uc *UserController) Create(c *gin.Context) {
	log := requestLog(c)
	var input services.CreateUserInput

	if err := helpers.Try(func() error {
		if bindErr := c.ShouldBindJSON(&input); bindErr != nil {
			helpers.BadRequest(c, bindErr.Error())
			return bindErr
		}
		user, svcErr := uc.service.Create(input)
		if svcErr != nil {
			if svcErr.Error() == "email already registered" {
				helpers.Conflict(c, svcErr.Error())
			} else {
				helpers.InternalError(c, svcErr.Error())
			}
			return svcErr
		}
		log.Info("user created", zap.String("id", user.ID.String()))
		helpers.Created(c, user)
		return nil
	}); err != nil {
		log.Error("UserController.Create", zap.Error(err))
	}
}

// GetByID godoc
// @Summary      Get user by ID
// @Description  Fetch a single user by UUID
// @Tags         Users
// @Produce      json
// @Param        id   path      string  true  "User UUID"
// @Success      200  {object}  helpers.Response{data=models.User}
// @Failure      404  {object}  helpers.Response
// @Failure      500  {object}  helpers.Response
// @Security     BearerAuth
// @Router       /users/{id} [get]
func (uc *UserController) GetByID(c *gin.Context) {
	log := requestLog(c)

	if err := helpers.Try(func() error {
		user, svcErr := uc.service.GetByID(c.Param("id"))
		if svcErr != nil {
			if errors.Is(svcErr, gorm.ErrRecordNotFound) {
				helpers.NotFound(c)
			} else {
				helpers.InternalError(c, svcErr.Error())
			}
			return svcErr
		}
		helpers.OK(c, "success", user)
		return nil
	}); err != nil {
		log.Error("UserController.GetByID", zap.Error(err), zap.String("id", c.Param("id")))
	}
}

// List godoc
// @Summary      List users
// @Description  Paginated list of all users
// @Tags         Users
// @Produce      json
// @Param        page       query     int  false  "Page number"    default(1)
// @Param        page_size  query     int  false  "Items per page" default(10)
// @Success      200        {object}  helpers.Response{data=[]models.User,meta=helpers.Meta}
// @Failure      500        {object}  helpers.Response
// @Security     BearerAuth
// @Router       /users [get]
func (uc *UserController) List(c *gin.Context) {
	log := requestLog(c)

	if err := helpers.Try(func() error {
		page, pageSize := helpers.ParsePagination(c)
		users, total, svcErr := uc.service.List(page, pageSize)
		if svcErr != nil {
			helpers.InternalError(c, svcErr.Error())
			return svcErr
		}
		helpers.OKPaginated(c, users, helpers.PaginationMeta(page, pageSize, total))
		return nil
	}); err != nil {
		log.Error("UserController.List", zap.Error(err))
	}
}

// Update godoc
// @Summary      Update user
// @Description  Partially update a user's fields
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id    path      string                   true  "User UUID"
// @Param        body  body      services.UpdateUserInput true  "Fields to update"
// @Success      200   {object}  helpers.Response{data=models.User}
// @Failure      400   {object}  helpers.Response
// @Failure      404   {object}  helpers.Response
// @Failure      500   {object}  helpers.Response
// @Security     BearerAuth
// @Router       /users/{id} [patch]
func (uc *UserController) Update(c *gin.Context) {
	log := requestLog(c)
	var input services.UpdateUserInput

	if err := helpers.Try(func() error {
		if bindErr := c.ShouldBindJSON(&input); bindErr != nil {
			helpers.BadRequest(c, bindErr.Error())
			return bindErr
		}
		user, svcErr := uc.service.Update(c.Param("id"), input)
		if svcErr != nil {
			if errors.Is(svcErr, gorm.ErrRecordNotFound) {
				helpers.NotFound(c)
			} else {
				helpers.InternalError(c, svcErr.Error())
			}
			return svcErr
		}
		helpers.OK(c, "updated", user)
		return nil
	}); err != nil {
		log.Error("UserController.Update", zap.Error(err), zap.String("id", c.Param("id")))
	}
}

// Delete godoc
// @Summary      Delete user
// @Description  Soft-delete a user by UUID
// @Tags         Users
// @Produce      json
// @Param        id   path      string  true  "User UUID"
// @Success      200  {object}  helpers.Response
// @Failure      404  {object}  helpers.Response
// @Failure      500  {object}  helpers.Response
// @Security     BearerAuth
// @Router       /users/{id} [delete]
func (uc *UserController) Delete(c *gin.Context) {
	log := requestLog(c)

	if err := helpers.Try(func() error {
		svcErr := uc.service.Delete(c.Param("id"))
		if svcErr != nil {
			if errors.Is(svcErr, gorm.ErrRecordNotFound) {
				helpers.NotFound(c)
			} else {
				helpers.InternalError(c, svcErr.Error())
			}
			return svcErr
		}
		helpers.OK(c, "deleted", nil)
		return nil
	}); err != nil {
		log.Error("UserController.Delete", zap.Error(err), zap.String("id", c.Param("id")))
	}
}
