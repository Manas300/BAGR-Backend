package services

import (
	"context"
	"fmt"

	"bagr-backend/internal/models"
	"bagr-backend/internal/repositories"
	"bagr-backend/internal/utils"
)

// UserService handles user business logic
type UserService struct {
	userRepo repositories.UserRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo repositories.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.User, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing user by email: %w", err)
	}
	if existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	existingUser, err = s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing user by username: %w", err)
	}
	if existingUser != nil {
		return nil, fmt.Errorf("user with username %s already exists", req.Username)
	}

	// TODO: Hash password before storing
	// For now, we'll store the plain password (NOT RECOMMENDED FOR PRODUCTION)
	user := &models.User{
		Email:     req.Email,
		Username:  req.Username,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Password:  req.Password, // TODO: Hash this
		Role:      req.Role,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		utils.GetLogger().WithError(err).Error("Failed to create user")
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	utils.GetLogger().WithField("user_id", user.ID).Info("User created successfully")
	return user, nil
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

// GetUserByEmail retrieves a user by email
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

// UpdateUser updates a user
func (s *UserService) UpdateUser(ctx context.Context, id int, req *models.UpdateUserRequest) (*models.User, error) {
	// Check if user exists
	existingUser, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if existingUser == nil {
		return nil, fmt.Errorf("user not found")
	}

	// Build updates map
	updates := make(map[string]interface{})
	
	if req.Email != nil {
		// Check if email is already taken by another user
		userWithEmail, err := s.userRepo.GetByEmail(ctx, *req.Email)
		if err != nil {
			return nil, fmt.Errorf("failed to check email availability: %w", err)
		}
		if userWithEmail != nil && userWithEmail.ID != id {
			return nil, fmt.Errorf("email %s is already taken", *req.Email)
		}
		updates["email"] = *req.Email
	}

	if req.Username != nil {
		// Check if username is already taken by another user
		userWithUsername, err := s.userRepo.GetByUsername(ctx, *req.Username)
		if err != nil {
			return nil, fmt.Errorf("failed to check username availability: %w", err)
		}
		if userWithUsername != nil && userWithUsername.ID != id {
			return nil, fmt.Errorf("username %s is already taken", *req.Username)
		}
		updates["username"] = *req.Username
	}

	if req.FirstName != nil {
		updates["first_name"] = *req.FirstName
	}

	if req.LastName != nil {
		updates["last_name"] = *req.LastName
	}

	if req.Role != nil {
		updates["role"] = *req.Role
	}

	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if len(updates) > 0 {
		if err := s.userRepo.Update(ctx, id, updates); err != nil {
			utils.GetLogger().WithError(err).Error("Failed to update user")
			return nil, fmt.Errorf("failed to update user: %w", err)
		}
	}

	// Return updated user
	updatedUser, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated user: %w", err)
	}

	utils.GetLogger().WithField("user_id", id).Info("User updated successfully")
	return updatedUser, nil
}

// DeleteUser deletes a user
func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	// Check if user exists
	existingUser, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if existingUser == nil {
		return fmt.Errorf("user not found")
	}

	if err := s.userRepo.Delete(ctx, id); err != nil {
		utils.GetLogger().WithError(err).Error("Failed to delete user")
		return fmt.Errorf("failed to delete user: %w", err)
	}

	utils.GetLogger().WithField("user_id", id).Info("User deleted successfully")
	return nil
}

// ListUsers retrieves a list of users with pagination
func (s *UserService) ListUsers(ctx context.Context, limit, offset int) ([]*models.User, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	users, err := s.userRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return users, nil
}
