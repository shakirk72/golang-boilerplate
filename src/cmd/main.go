// @title           Go MVC API
// @version         1.0
// @description     Structured Go MVC REST API with PostgreSQL and GORM.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  mohit.chack@rosstechnologies.in

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and your JWT token.

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/yourorg/go-mvc-app/docs"
	"github.com/yourorg/go-mvc-app/src/config"
	"github.com/yourorg/go-mvc-app/src/routes"
	"github.com/yourorg/go-mvc-app/src/utils/database"
	"github.com/yourorg/go-mvc-app/src/utils/logger"
	"go.uber.org/zap"
)

func main() {
	config.Load(".env")

	logger.Init(config.App.App.Env)
	defer logger.Sync()

	if err := config.App.Validate(); err != nil {
		logger.Fatal("invalid configuration", zap.Error(err))
	}

	if config.App.App.IsProd {
		gin.SetMode(gin.ReleaseMode)
	}

	database.Connect()
	database.Migrate()

	r := gin.New()
	routes.Setup(r)

	addr := fmt.Sprintf(":%s", config.App.App.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Info("server starting", zap.String("addr", addr), zap.String("env", config.App.App.Env))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server failed to start", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutdown signal received, draining connections...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("server forced to shutdown", zap.Error(err))
	}

	logger.Info("server exited cleanly")
}
