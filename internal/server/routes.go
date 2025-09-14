package server

import (
	"bagr-backend/internal/controllers"
	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(router *gin.Engine, controllers *Controllers) {
	// Health check routes
	router.GET("/health", controllers.Health.Health)
	router.GET("/ready", controllers.Health.Ready)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// User routes
		users := v1.Group("/users")
		{
			users.POST("", controllers.User.CreateUser)
			users.GET("", controllers.User.ListUsers)
			users.GET("/:id", controllers.User.GetUser)
			users.PUT("/:id", controllers.User.UpdateUser)
			users.DELETE("/:id", controllers.User.DeleteUser)
		}

		// Future routes can be added here:
		// auctions := v1.Group("/auctions")
		// bids := v1.Group("/bids")
		// tracks := v1.Group("/tracks")
	}
}

// Controllers holds all controller instances
type Controllers struct {
	Health *controllers.HealthController
	User   *controllers.UserController
}

// NewControllers creates and returns all controller instances
func NewControllers(services *Services) *Controllers {
	return &Controllers{
		Health: controllers.NewHealthController(),
		User:   controllers.NewUserController(services.User),
	}
}
