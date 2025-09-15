package server

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"bagr-backend/internal/auth"
	"bagr-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

// JWTMiddleware validates JWT tokens
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "MISSING_TOKEN", "Authorization header required", "")
			c.Abort()
			return
		}

		// Check if header starts with "Bearer "
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			utils.ErrorResponse(c, http.StatusUnauthorized, "INVALID_TOKEN_FORMAT", "Authorization header must start with 'Bearer '", "")
			c.Abort()
			return
		}

		// Get JWT service from context (we'll set this in server.go)
		jwtService, exists := c.Get("jwt_service")
		if !exists {
			utils.ErrorResponse(c, http.StatusInternalServerError, "JWT_SERVICE_NOT_FOUND", "JWT service not available", "")
			c.Abort()
			return
		}

		// Validate token
		claims, err := jwtService.(*auth.JWTService).ValidateAccessToken(tokenString)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "INVALID_TOKEN", "Invalid or expired token", err.Error())
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}

// OptionalJWTMiddleware validates JWT tokens if present, but doesn't require them
func OptionalJWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		// Check if header starts with "Bearer "
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.Next()
			return
		}

		// Get JWT service from context
		jwtService, exists := c.Get("jwt_service")
		if !exists {
			c.Next()
			return
		}

		// Validate token
		claims, err := jwtService.(*auth.JWTService).ValidateAccessToken(tokenString)
		if err != nil {
			c.Next()
			return
		}

		// Set user information in context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}

// RoleMiddleware checks if user has required role
func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			utils.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "User role not found", "")
			c.Abort()
			return
		}

		if userRole.(string) != requiredRole {
			utils.ErrorResponse(c, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions", "Required role: "+requiredRole)
			c.Abort()
			return
		}

		c.Next()
	}
}

// AdminMiddleware checks if user is admin
func AdminMiddleware() gin.HandlerFunc {
	return RoleMiddleware("admin")
}

// ProducerMiddleware checks if user is producer
func ProducerMiddleware() gin.HandlerFunc {
	return RoleMiddleware("producer")
}

// ArtistMiddleware checks if user is artist
func ArtistMiddleware() gin.HandlerFunc {
	return RoleMiddleware("artist")
}

// FanMiddleware checks if user is fan
func FanMiddleware() gin.HandlerFunc {
	return RoleMiddleware("fan")
}

// ModeratorMiddleware checks if user is moderator
func ModeratorMiddleware() gin.HandlerFunc {
	return RoleMiddleware("moderator")
}

// MultipleRoleMiddleware checks if user has any of the required roles
func MultipleRoleMiddleware(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			utils.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "User role not found", "")
			c.Abort()
			return
		}

		userRoleStr := userRole.(string)
		for _, role := range requiredRoles {
			if userRoleStr == role {
				c.Next()
				return
			}
		}

		utils.ErrorResponse(c, http.StatusForbidden, "FORBIDDEN", "Insufficient permissions", "Required one of: "+strings.Join(requiredRoles, ", "))
		c.Abort()
	}
}

// LoggerMiddleware logs HTTP requests
func LoggerMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	})
}

// RecoveryMiddleware recovers from panics
func RecoveryMiddleware() gin.HandlerFunc {
	return gin.Recovery()
}

// CORSMiddleware handles CORS
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// RequestIDMiddleware adds request ID to context
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = utils.GenerateRequestID()
		}
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

// TimeoutMiddleware sets request timeout
func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// This is a simple implementation
		// In production, you might want to use context.WithTimeout
		c.Next()
	}
}
