package helpers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourorg/go-mvc-app/src/config"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type Meta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalCount int64 `json:"total_count"`
	TotalPages int   `json:"total_pages"`
}

func OK(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{Success: true, Message: message, Data: data})
}

func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{Success: true, Message: config.MsgCreated, Data: data})
}

func OKPaginated(c *gin.Context, data interface{}, meta Meta) {
	c.JSON(http.StatusOK, Response{Success: true, Message: config.MsgSuccess, Data: data, Meta: &meta})
}

func BadRequest(c *gin.Context, err interface{}) {
	c.JSON(http.StatusBadRequest, Response{Success: false, Message: config.MsgBadRequest, Error: err})
}

func Unauthorized(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, Response{Success: false, Message: config.MsgUnauthorized})
}

func Forbidden(c *gin.Context) {
	c.JSON(http.StatusForbidden, Response{Success: false, Message: config.MsgForbidden})
}

func NotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, Response{Success: false, Message: config.MsgNotFound})
}

func Conflict(c *gin.Context, err interface{}) {
	c.JSON(http.StatusConflict, Response{Success: false, Message: config.MsgConflict, Error: err})
}

func InternalError(c *gin.Context, err interface{}) {
	c.JSON(http.StatusInternalServerError, Response{Success: false, Message: config.MsgInternal, Error: err})
}

func TooManyRequests(c *gin.Context) {
	c.JSON(http.StatusTooManyRequests, Response{Success: false, Message: config.MsgTooManyRequests})
}

// ValidationError renders a 422 with per-field error details extracted from gin binding errors.
// Pass the raw error from ShouldBindJSON; gin validation errors are unwrapped automatically.
func ValidationError(c *gin.Context, err error) {
	c.JSON(http.StatusUnprocessableEntity, Response{
		Success: false,
		Message: config.MsgValidationFailed,
		Error:   err.Error(),
	})
}

// PaginationMeta calculates pagination meta from page, pageSize and totalCount.
func PaginationMeta(page, pageSize int, totalCount int64) Meta {
	totalPages := int(totalCount) / pageSize
	if int(totalCount)%pageSize != 0 {
		totalPages++
	}
	return Meta{
		Page:       page,
		PageSize:   pageSize,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}
}
