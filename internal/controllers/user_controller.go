package controllers

import (
	"net/http"
	"strconv"

	"bagr-backend/internal/models"
	"bagr-backend/internal/services"
	"bagr-backend/internal/utils"
	"github.com/gin-gonic/gin"
)

// UserController handles user-related endpoints
type UserController struct {
	userService *services.UserService
}

// NewUserController creates a new user controller
func NewUserController(userService *services.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// CreateUser handles user creation
// @Summary Create a new user
// @Description Create a new user account
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.CreateUserRequest true "User creation data"
// @Success 201 {object} models.UserResponse
// @Failure 400 {object} utils.APIResponse
// @Failure 409 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /users [post]
func (uc *UserController) CreateUser(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	user, err := uc.userService.CreateUser(c.Request.Context(), &req)
	if err != nil {
		if err.Error() == "user with email "+req.Email+" already exists" ||
		   err.Error() == "user with username "+req.Username+" already exists" {
			utils.ErrorResponse(c, http.StatusConflict, "CONFLICT", err.Error(), "")
			return
		}
		utils.InternalErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "User created successfully", user.ToResponse())
}

// GetUser handles getting a user by ID
// @Summary Get user by ID
// @Description Get a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} utils.APIResponse
// @Failure 404 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /users/{id} [get]
func (uc *UserController) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid user ID", "ID must be a valid integer")
		return
	}

	user, err := uc.userService.GetUserByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "user not found" {
			utils.NotFoundResponse(c, "User")
			return
		}
		utils.InternalErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User retrieved successfully", user.ToResponse())
}

// UpdateUser handles user updates
// @Summary Update user
// @Description Update user information
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body models.UpdateUserRequest true "User update data"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} utils.APIResponse
// @Failure 404 {object} utils.APIResponse
// @Failure 409 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /users/{id} [put]
func (uc *UserController) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid user ID", "ID must be a valid integer")
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	user, err := uc.userService.UpdateUser(c.Request.Context(), id, &req)
	if err != nil {
		if err.Error() == "user not found" {
			utils.NotFoundResponse(c, "User")
			return
		}
		if err.Error() == "email "+*req.Email+" is already taken" ||
		   err.Error() == "username "+*req.Username+" is already taken" {
			utils.ErrorResponse(c, http.StatusConflict, "CONFLICT", err.Error(), "")
			return
		}
		utils.InternalErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User updated successfully", user.ToResponse())
}

// DeleteUser handles user deletion
// @Summary Delete user
// @Description Delete a user account
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} utils.APIResponse
// @Failure 400 {object} utils.APIResponse
// @Failure 404 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /users/{id} [delete]
func (uc *UserController) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "Invalid user ID", "ID must be a valid integer")
		return
	}

	err = uc.userService.DeleteUser(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "user not found" {
			utils.NotFoundResponse(c, "User")
			return
		}
		utils.InternalErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User deleted successfully", nil)
}

// ListUsers handles listing users with pagination
// @Summary List users
// @Description Get a paginated list of users
// @Tags users
// @Accept json
// @Produce json
// @Param limit query int false "Number of users to return (default: 10, max: 100)"
// @Param offset query int false "Number of users to skip (default: 0)"
// @Success 200 {array} models.UserResponse
// @Failure 400 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /users [get]
func (uc *UserController) ListUsers(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_LIMIT", "Invalid limit parameter", "Limit must be a positive integer")
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "INVALID_OFFSET", "Invalid offset parameter", "Offset must be a non-negative integer")
		return
	}

	users, err := uc.userService.ListUsers(c.Request.Context(), limit, offset)
	if err != nil {
		utils.InternalErrorResponse(c, err)
		return
	}

	// Convert to response format
	userResponses := make([]*models.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = user.ToResponse()
	}

	utils.SuccessResponse(c, http.StatusOK, "Users retrieved successfully", userResponses)
}
