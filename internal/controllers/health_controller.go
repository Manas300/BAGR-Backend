package controllers

import (
	"net/http"
	"time"

	"bagr-backend/internal/utils"
	"github.com/gin-gonic/gin"
)

// HealthController handles health check endpoints
type HealthController struct{}

// NewHealthController creates a new health controller
func NewHealthController() *HealthController {
	return &HealthController{}
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
	Service   string    `json:"service"`
}

// Health handles the health check endpoint
// @Summary Health check
// @Description Returns the health status of the service
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func (h *HealthController) Health(c *gin.Context) {
	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Version:   "1.0.0", // TODO: Get from build info
		Service:   "bagr-backend",
	}

	utils.SuccessResponse(c, http.StatusOK, "Service is healthy", response)
}

// Ready handles the readiness check endpoint
// @Summary Readiness check
// @Description Returns whether the service is ready to handle requests
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /ready [get]
func (h *HealthController) Ready(c *gin.Context) {
	// TODO: Add actual readiness checks (database connectivity, etc.)
	response := HealthResponse{
		Status:    "ready",
		Timestamp: time.Now(),
		Version:   "1.0.0",
		Service:   "bagr-backend",
	}

	utils.SuccessResponse(c, http.StatusOK, "Service is ready", response)
}
