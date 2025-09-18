package handlers

import (
	"net/http"
	"strconv"

	"bagr-backend/internal/models"
	"bagr-backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// ProfileHandlers handles profile-related HTTP requests
type ProfileHandlers struct {
	profileService *services.ProfileService
	s3Service      *services.S3Service
	logger         *logrus.Logger
}

// NewProfileHandlers creates a new profile handlers instance
func NewProfileHandlers(profileService *services.ProfileService, s3Service *services.S3Service, logger *logrus.Logger) *ProfileHandlers {
	return &ProfileHandlers{
		profileService: profileService,
		s3Service:      s3Service,
		logger:         logger,
	}
}

// GetProfile retrieves the current user's profile
func (h *ProfileHandlers) GetProfile(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.Error("User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized",
		})
		return
	}

	userIDInt, ok := userID.(int)
	if !ok {
		h.logger.Error("Invalid user ID type in context")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal server error",
		})
		return
	}

	// Get profile from service
	profile, err := h.profileService.GetProfileByUserID(userIDInt)
	if err != nil {
		h.logger.WithError(err).WithField("user_id", userIDInt).Error("Failed to get profile")
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Profile not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Profile retrieved successfully",
		"data":    profile.ToResponse(),
	})
}

// UpdateProfile updates the current user's profile
func (h *ProfileHandlers) UpdateProfile(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.Error("User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized",
		})
		return
	}

	userIDInt, ok := userID.(int)
	if !ok {
		h.logger.Error("Invalid user ID type in context")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal server error",
		})
		return
	}

	// Parse request body
	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Failed to bind update profile request")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data",
			"error":   err.Error(),
		})
		return
	}

	// Check if profile exists
	exists, err := h.profileService.ProfileExists(userIDInt)
	if err != nil {
		h.logger.WithError(err).WithField("user_id", userIDInt).Error("Failed to check profile existence")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal server error",
		})
		return
	}

	var profile *models.Profile
	if !exists {
		// Create new profile
		createReq := models.CreateProfileRequest{
			DisplayName:     *req.DisplayName,
			Bio:             getStringValue(req.Bio),
			Location:        getStringValue(req.Location),
			WebsiteURL:      getStringValue(req.WebsiteURL),
			YouTubeHandle:   getStringValue(req.YouTubeHandle),
			TikTokHandle:    getStringValue(req.TikTokHandle),
			InstagramHandle: getStringValue(req.InstagramHandle),
			TwitterHandle:   getStringValue(req.TwitterHandle),
		}
		profile, err = h.profileService.CreateProfile(userIDInt, &createReq)
	} else {
		// Update existing profile
		profile, err = h.profileService.UpdateProfile(userIDInt, &req)
	}

	if err != nil {
		h.logger.WithError(err).WithField("user_id", userIDInt).Error("Failed to update profile")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update profile",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Profile updated successfully",
		"data":    profile.ToResponse(),
	})
}

// UploadProfileImage uploads a profile image
func (h *ProfileHandlers) UploadProfileImage(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.Error("User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized",
		})
		return
	}

	userIDInt, ok := userID.(int)
	if !ok {
		h.logger.Error("Invalid user ID type in context")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal server error",
		})
		return
	}

	// Get the uploaded file
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		h.logger.WithError(err).Error("Failed to get uploaded file")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "No image file provided",
		})
		return
	}
	defer file.Close()

	// Validate file type
	contentType := header.Header.Get("Content-Type")
	if !h.s3Service.ValidateImageType(contentType) {
		h.logger.WithField("content_type", contentType).Error("Invalid image type")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid image type. Only JPEG, PNG, GIF, and WebP are allowed",
		})
		return
	}

	// Upload to S3
	imageURL, err := h.s3Service.UploadProfileImage(c.Request.Context(), userIDInt, file, contentType)
	if err != nil {
		h.logger.WithError(err).WithField("user_id", userIDInt).Error("Failed to upload profile image")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to upload image",
			"error":   err.Error(),
		})
		return
	}

	// Update profile with new image URL
	err = h.profileService.UpdateProfileImage(userIDInt, imageURL)
	if err != nil {
		h.logger.WithError(err).WithField("user_id", userIDInt).Error("Failed to update profile image URL")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update profile with new image",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Profile image uploaded successfully",
		"data": gin.H{
			"image_url": imageURL,
		},
	})
}

// GetProfileByID retrieves a profile by user ID (public endpoint)
func (h *ProfileHandlers) GetProfileByID(c *gin.Context) {
	// Get user ID from URL parameter
	userIDStr := c.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		h.logger.WithError(err).WithField("user_id_str", userIDStr).Error("Invalid user ID")
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid user ID",
		})
		return
	}

	// Get profile from service
	profile, err := h.profileService.GetProfileByUserID(userID)
	if err != nil {
		h.logger.WithError(err).WithField("user_id", userID).Error("Failed to get profile")
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Profile not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Profile retrieved successfully",
		"data":    profile.ToResponse(),
	})
}

// Helper function to get string value from pointer
func getStringValue(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}
