package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yourorg/go-mvc-app/src/config"
	"github.com/yourorg/go-mvc-app/src/helpers"
	"github.com/yourorg/go-mvc-app/src/services"
)

// Authenticate validates the Bearer JWT token and injects userID + role into context.
func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader(config.HeaderAuthorization)
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			helpers.Unauthorized(c)
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims, err := services.ValidateJWT(tokenStr)
		if err != nil {
			helpers.Unauthorized(c)
			c.Abort()
			return
		}

		c.Set(config.CtxUserID, claims.UserID)
		c.Set(config.CtxUserRole, claims.Role)
		c.Next()
	}
}

// RequireRole aborts with 403 if the authenticated user's role is not in the allowed list.
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString(config.CtxUserRole)
		if !helpers.ContainsString(roles, role) {
			helpers.Forbidden(c)
			c.Abort()
			return
		}
		c.Next()
	}
}
