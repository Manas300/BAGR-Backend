package server

import (
	"bagr-backend/internal/auth"
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
		// Authentication routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", controllers.Auth.Register)
			auth.POST("/login", controllers.Auth.Login)
			auth.GET("/verify", controllers.Auth.VerifyEmail)
			auth.POST("/forgot-password", controllers.Auth.ForgotPassword)
			auth.GET("/reset-password", controllers.Auth.ResetPasswordPage)
			auth.POST("/reset-password", controllers.Auth.ResetPassword)
			auth.POST("/refresh", controllers.Auth.RefreshToken)
			auth.GET("/roles", controllers.Auth.GetRoles)
		}

		// Protected routes (require authentication)
		protected := v1.Group("/")
		protected.Use(JWTMiddleware())
		{
			// Auth protected routes
			authProtected := protected.Group("/auth")
			{
				authProtected.GET("/profile", controllers.Auth.GetProfile)
				authProtected.PUT("/profile", controllers.Auth.UpdateProfile)
				authProtected.POST("/logout", controllers.Auth.Logout)
			}

			// User routes (protected)
			users := protected.Group("/users")
			{
				users.POST("", controllers.User.CreateUser)
				users.GET("", controllers.User.ListUsers)
				users.GET("/:id", controllers.User.GetUser)
				users.PUT("/:id", controllers.User.UpdateUser)
				users.DELETE("/:id", controllers.User.DeleteUser)
			}

			// Future protected routes can be added here:
			// auctions := protected.Group("/auctions")
			// bids := protected.Group("/bids")
			// tracks := protected.Group("/tracks")
		}
	}
}

// Controllers holds all controller instances
type Controllers struct {
	Health *controllers.HealthController
	User   *controllers.UserController
	Auth   *auth.AuthHandlers
}

// NewControllers creates and returns all controller instances
func NewControllers(services *Services) *Controllers {
	return &Controllers{
		Health: controllers.NewHealthController(),
		User:   controllers.NewUserController(services.User),
		Auth:   auth.NewAuthHandlers(services.Auth),
	}
}
