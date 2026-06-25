package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/yourorg/go-mvc-app/src/config"
	"github.com/yourorg/go-mvc-app/src/controllers"
	"github.com/yourorg/go-mvc-app/src/middleware"
	"github.com/yourorg/go-mvc-app/src/repositories"
	"github.com/yourorg/go-mvc-app/src/services"
	"github.com/yourorg/go-mvc-app/src/utils/database"
)

type healthResponse struct {
	Status   string `json:"status"    example:"ok"`
	Env      string `json:"env"       example:"dev"`
	Database string `json:"database"  example:"ok"`
}

// Health godoc
// @Summary      Health check
// @Description  Returns service liveness and database connectivity status
// @Tags         Health
// @Produce      json
// @Success      200  {object}  healthResponse
// @Success      503  {object}  healthResponse
// @Router       /health [get]
func Health(c *gin.Context) {
	dbStatus := "ok"
	httpStatus := http.StatusOK

	if err := database.Ping(); err != nil {
		dbStatus = "unreachable"
		httpStatus = http.StatusServiceUnavailable
	}

	c.JSON(httpStatus, healthResponse{
		Status:   "ok",
		Env:      config.App.App.Env,
		Database: dbStatus,
	})
}

func Setup(r *gin.Engine, drainer *middleware.Drainer) {
	r.Use(drainer.Handler())
	r.Use(middleware.Recovery())
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(middleware.RateLimit(config.DefaultRateLimitPerMin))

	// Swagger UI — disabled in prod
	if !config.App.App.IsProd {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// wire up dependencies
	userRepo := repositories.NewUserRepository()
	userSvc := services.NewUserService(userRepo)
	userCtrl := controllers.NewUserController(userSvc)

	authSvc := services.NewAuthService(userRepo)
	authCtrl := controllers.NewAuthController(authSvc, userSvc)

	v1 := r.Group(config.APIV1)
	{
		// Public endpoints
		v1.GET("/health", Health)

		auth := v1.Group("/auth")
		{
			auth.POST("/signup", authCtrl.Signup)
			auth.POST("/login", authCtrl.Login)
		}

		// Protected endpoints — require valid JWT
		protected := v1.Group("")
		protected.Use(middleware.Authenticate())
		{
			users := protected.Group("/users")
			{
				users.GET("", userCtrl.List)
				users.GET("/:id", userCtrl.GetByID)
				users.PATCH("/:id", userCtrl.Update)
				// DELETE requires admin role
				users.DELETE("/:id", middleware.RequireRole(config.RoleAdmin), userCtrl.Delete)
			}
		}
	}
}
