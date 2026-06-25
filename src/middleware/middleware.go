package middleware

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourorg/go-mvc-app/src/config"
	"github.com/yourorg/go-mvc-app/src/helpers"
	"github.com/yourorg/go-mvc-app/src/utils/logger"
	"go.uber.org/zap"
)

// RequestID injects a unique X-Request-ID into every request context.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.GetHeader(config.HeaderRequestID)
		if rid == "" {
			token, _ := helpers.GenerateToken(8)
			rid = fmt.Sprintf("req-%s", token)
		}
		c.Set(config.CtxRequestID, rid)
		c.Header(config.HeaderRequestID, rid)
		c.Next()
	}
}

// Logger logs method, path, status, latency, and request ID for every request.
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)

		log := logger.WithRequestID(c.GetString(config.CtxRequestID))
		log.Info("request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", latency),
			zap.String("ip", c.ClientIP()),
		)
	}
}

// Recovery catches panics, logs the stack trace, and returns a 500.
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				stack := string(debug.Stack())
				logger.Error("panic recovered",
					zap.Any("panic", r),
					zap.String("path", c.Request.URL.Path),
					zap.String("stack", stack),
				)
				helpers.InternalError(c, fmt.Sprintf("%v", r))
				c.Abort()
			}
		}()
		c.Next()
	}
}

// CORS sets headers based on CORS_ORIGINS config. Use "*" in config to allow all origins.
func CORS() gin.HandlerFunc {
	origins := config.App.App.CORSOrigins
	allowAll := len(origins) == 1 && origins[0] == "*"

	allowed := make(map[string]bool, len(origins))
	for _, o := range origins {
		allowed[o] = true
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		if allowAll {
			c.Header("Access-Control-Allow-Origin", "*")
		} else if origin != "" && allowed[origin] {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Vary", "Origin")
		}

		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin,Content-Type,Authorization,X-Request-ID")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

// Timeout attaches a deadline to every request context. Handlers using
// c.Request.Context() (e.g. database queries, API calls) will respect it.
func Timeout(d time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), d)
		defer cancel()
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// ipBucket is a fixed-window counter for one IP address.
type ipBucket struct {
	mu      sync.Mutex
	count   int
	resetAt time.Time
}

var rlMap sync.Map

// RateLimit enforces a per-IP fixed-window rate limit of rpm requests per minute.
func RateLimit(rpm int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()

		val, _ := rlMap.LoadOrStore(ip, &ipBucket{resetAt: now.Add(time.Minute)})
		bucket := val.(*ipBucket)

		bucket.mu.Lock()
		if now.After(bucket.resetAt) {
			bucket.count = 0
			bucket.resetAt = now.Add(time.Minute)
		}
		bucket.count++
		over := bucket.count > rpm
		bucket.mu.Unlock()

		if over {
			c.JSON(http.StatusTooManyRequests, helpers.Response{
				Success: false,
				Message: config.MsgTooManyRequests,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
