package auth

import (
	"net/http"
	"strconv"

	"bagr-backend/internal/models"
	"bagr-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

// AuthHandlers handles authentication HTTP requests
type AuthHandlers struct {
	authService *AuthService
}

// NewAuthHandlers creates new authentication handlers
func NewAuthHandlers(authService *AuthService) *AuthHandlers {
	return &AuthHandlers{
		authService: authService,
	}
}

// Register handles user registration
// POST /api/v1/auth/register
func (h *AuthHandlers) Register(c *gin.Context) {
	logger := utils.GetLogger()

	// Log incoming request
	logger.WithFields(map[string]interface{}{
		"endpoint": "/api/v1/auth/register",
		"method":   "POST",
		"ip":       c.ClientIP(),
	}).Info("Registration request received")

	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.WithError(err).Error("Failed to bind registration request JSON")
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request data", err.Error())
		return
	}

	// Log request data (sanitized)
	logger.WithFields(map[string]interface{}{
		"email":           req.Email,
		"username":        req.Username,
		"first_name":      req.FirstName,
		"last_name":       req.LastName,
		"role":            req.Role,
		"password_length": len(req.Password),
	}).Info("Processing registration request")

	// Validate role
	if !isValidRole(req.Role) {
		logger.WithField("role", req.Role).Error("Invalid role provided")
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ROLE", "Invalid role", "Role must be one of: admin, moderator, producer, artist, fan")
		return
	}

	// Register user
	logger.Info("Attempting to register user")
	response, err := h.authService.RegisterUser(&req)
	if err != nil {
		logger.WithError(err).Error("User registration failed")
		utils.ErrorResponse(c, http.StatusBadRequest, "REGISTRATION_FAILED", "Registration failed", err.Error())
		return
	}

	logger.WithFields(map[string]interface{}{
		"user_id":  response.User.ID,
		"email":    response.User.Email,
		"username": response.User.Username,
	}).Info("User registration successful")

	utils.SuccessResponse(c, http.StatusCreated, "User registered successfully. Please check your email for verification.", response)
}

// Login handles user login
// POST /api/v1/auth/login
func (h *AuthHandlers) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request data", err.Error())
		return
	}

	// Login user
	response, err := h.authService.LoginUser(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "LOGIN_FAILED", "Login failed", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Login successful", response)
}

// VerifyEmail handles email verification
// GET /api/v1/auth/verify?token=xxx
func (h *AuthHandlers) VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "MISSING_TOKEN", "Missing token", "Verification token is required")
		return
	}

	// Verify email
	user, err := h.authService.VerifyEmail(token)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "VERIFICATION_FAILED", "Verification failed", err.Error())
		return
	}

	// Return JSON response for API consistency
	utils.SuccessResponse(c, http.StatusOK, "Email verified successfully", gin.H{
		"message": "Your email has been successfully verified! You can now log in to your account.",
		"user_id": user.ID,
		"email":   user.Email,
	})
}

// ForgotPassword handles password reset request
// POST /api/v1/auth/forgot-password
func (h *AuthHandlers) ForgotPassword(c *gin.Context) {
	var req models.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request data", err.Error())
		return
	}

	// Send reset email
	err := h.authService.ForgotPassword(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "EMAIL_SEND_FAILED", "Failed to send reset email", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Password reset email sent", gin.H{
		"message": "If an account with this email exists, you will receive a password reset link.",
	})
}

// ResetPasswordPage handles password reset page display
// GET /api/v1/auth/reset-password?token=xxx
func (h *AuthHandlers) ResetPasswordPage(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		// Return HTML page for password reset
		c.HTML(http.StatusOK, "reset_password.html", gin.H{
			"title": "Reset Password - BAGR Auction System",
			"error": "Missing reset token",
		})
		return
	}

	// Return HTML page with token - validation will happen on form submission
	c.HTML(http.StatusOK, "reset_password.html", gin.H{
		"title": "Reset Password - BAGR Auction System",
		"token": token,
	})
}

// ResetPassword handles password reset
// POST /api/v1/auth/reset-password
func (h *AuthHandlers) ResetPassword(c *gin.Context) {
	logger := utils.GetLogger()
	
	var req models.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.WithError(err).Error("Failed to bind reset password request JSON")
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request data", err.Error())
		return
	}

	logger.WithFields(map[string]interface{}{
		"token_length": len(req.Token),
		"has_password": len(req.Password) > 0,
		"has_confirm":  len(req.ConfirmPassword) > 0,
	}).Info("Processing password reset request")

	// Reset password
	err := h.authService.ResetPassword(&req)
	if err != nil {
		logger.WithError(err).WithField("token", req.Token).Error("Password reset failed")
		utils.ErrorResponse(c, http.StatusBadRequest, "PASSWORD_RESET_FAILED", "Password reset failed", err.Error())
		return
	}

	logger.WithField("token", req.Token).Info("Password reset successful")
	utils.SuccessResponse(c, http.StatusOK, "Password reset successful", gin.H{
		"message": "Your password has been successfully reset. You can now log in with your new password.",
	})
}

// RefreshToken handles token refresh
// POST /api/v1/auth/refresh
func (h *AuthHandlers) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request data", err.Error())
		return
	}

	// Refresh token
	response, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "TOKEN_REFRESH_FAILED", "Token refresh failed", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Token refreshed successfully", response)
}

// GetProfile handles getting user profile
// GET /api/v1/auth/profile
func (h *AuthHandlers) GetProfile(c *gin.Context) {
	// Get user ID from JWT middleware
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized", "User ID not found in token")
		return
	}

	// Convert to int
	uid, ok := userID.(int)
	if !ok {
		utils.ErrorResponse(c, http.StatusInternalServerError, "INVALID_USER_ID", "Invalid user ID", "User ID is not a valid integer")
		return
	}

	// Get user from database
	user, err := h.authService.getUserByID(uid)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "USER_NOT_FOUND", "User not found", "User profile not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Profile retrieved successfully", user.ToResponse())
}

// UpdateProfile handles updating user profile
// PUT /api/v1/auth/profile
func (h *AuthHandlers) UpdateProfile(c *gin.Context) {
	// Get user ID from JWT middleware
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized", "User ID not found in token")
		return
	}

	// Convert to int
	uid, ok := userID.(int)
	if !ok {
		utils.ErrorResponse(c, http.StatusInternalServerError, "INVALID_USER_ID", "Invalid user ID", "User ID is not a valid integer")
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request data", err.Error())
		return
	}

	// Update user profile
	err := h.authService.updateUserProfile(uid, &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "PROFILE_UPDATE_FAILED", "Profile update failed", err.Error())
		return
	}

	// Get updated user
	user, err := h.authService.getUserByID(uid)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "PROFILE_RETRIEVAL_FAILED", "Failed to retrieve updated profile", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Profile updated successfully", user.ToResponse())
}

// Logout handles user logout
// POST /api/v1/auth/logout
func (h *AuthHandlers) Logout(c *gin.Context) {
	// In a stateless JWT system, logout is handled client-side
	// by removing the token from storage
	// We could implement a token blacklist here if needed

	utils.SuccessResponse(c, http.StatusOK, "Logout successful", gin.H{
		"message": "You have been successfully logged out. Please remove your tokens from client storage.",
	})
}

// GetRoles handles getting available user roles
// GET /api/v1/auth/roles
func (h *AuthHandlers) GetRoles(c *gin.Context) {
	roles := []gin.H{
		{"value": "producer", "label": "Producer", "description": "Music creators who sell beats"},
		{"value": "artist", "label": "Artist", "description": "Music creators who buy beats"},
		{"value": "fan", "label": "Fan", "description": "General users who participate in auctions"},
		{"value": "moderator", "label": "Moderator", "description": "Platform moderators"},
		{"value": "admin", "label": "Admin", "description": "Platform administrators"},
	}

	utils.SuccessResponse(c, http.StatusOK, "Roles retrieved successfully", roles)
}

// Helper functions

// isValidRole checks if the role is valid
func isValidRole(role models.UserRole) bool {
	validRoles := []models.UserRole{
		models.UserRoleAdmin,
		models.UserRoleModerator,
		models.UserRoleProducer,
		models.UserRoleArtist,
		models.UserRoleFan,
	}

	for _, validRole := range validRoles {
		if role == validRole {
			return true
		}
	}
	return false
}

// Add updateUserProfile method to AuthService
func (a *AuthService) updateUserProfile(userID int, req *models.UpdateUserRequest) error {
	// Build dynamic update query
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.Email != nil {
		setParts = append(setParts, "email = $"+strconv.Itoa(argIndex))
		args = append(args, *req.Email)
		argIndex++
	}

	if req.Username != nil {
		setParts = append(setParts, "username = $"+strconv.Itoa(argIndex))
		args = append(args, *req.Username)
		argIndex++
	}

	if req.FirstName != nil {
		setParts = append(setParts, "first_name = $"+strconv.Itoa(argIndex))
		args = append(args, *req.FirstName)
		argIndex++
	}

	if req.LastName != nil {
		setParts = append(setParts, "last_name = $"+strconv.Itoa(argIndex))
		args = append(args, *req.LastName)
		argIndex++
	}

	if req.Role != nil {
		setParts = append(setParts, "role = $"+strconv.Itoa(argIndex))
		args = append(args, *req.Role)
		argIndex++
	}

	if req.Status != nil {
		setParts = append(setParts, "status = $"+strconv.Itoa(argIndex))
		args = append(args, *req.Status)
		argIndex++
	}

	if len(setParts) == 0 {
		return nil // No updates to make
	}

	// Add updated_at
	setParts = append(setParts, "updated_at = $"+strconv.Itoa(argIndex))
	args = append(args, "NOW()")
	argIndex++

	// Add WHERE clause
	args = append(args, userID)

	query := "UPDATE users SET " + joinStrings(setParts, ", ") + " WHERE id = $" + strconv.Itoa(len(args))

	_, err := a.db.Exec(query, args...)
	return err
}

// joinStrings joins a slice of strings with a separator
func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	if len(strs) == 1 {
		return strs[0]
	}

	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}
