package models

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID                  int        `json:"id" db:"id"`
	Email               string     `json:"email" db:"email"`
	Username            string     `json:"username" db:"username"`
	FirstName           string     `json:"first_name" db:"first_name"`
	LastName            string     `json:"last_name" db:"last_name"`
	Password            string     `json:"-" db:"password"` // Never expose password in JSON
	PasswordHash        string     `json:"-" db:"password_hash"`
	Role                UserRole   `json:"role" db:"role"`
	Status              UserStatus `json:"status" db:"status"`
	EmailVerified       bool       `json:"email_verified" db:"email_verified"`
	VerificationToken   *string    `json:"-" db:"verification_token"`
	ResetToken          *string    `json:"-" db:"reset_token"`
	ResetTokenExpires   *time.Time `json:"-" db:"reset_token_expires"`
	LastLoginAt         *time.Time `json:"last_login_at" db:"last_login_at"`
	CreatedAt           time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at" db:"updated_at"`
}

// UserRole represents user roles in the system
type UserRole string

const (
	UserRoleAdmin     UserRole = "admin"
	UserRoleArtist    UserRole = "artist"
	UserRoleBuyer     UserRole = "buyer"
	UserRoleModerator UserRole = "moderator"
	UserRoleProducer  UserRole = "producer"
	UserRoleFan       UserRole = "fan"
)

// UserStatus represents user account status
type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusInactive  UserStatus = "inactive"
	UserStatusSuspended UserStatus = "suspended"
)

// CreateUserRequest represents the request payload for creating a user
type CreateUserRequest struct {
	Email           string   `json:"email" binding:"required,email"`
	Username        string   `json:"username" binding:"required,min=3,max=50"`
	FirstName       string   `json:"first_name" binding:"required,min=1,max=100"`
	LastName        string   `json:"last_name" binding:"required,min=1,max=100"`
	Password        string   `json:"password" binding:"required,min=8"`
	ConfirmPassword string   `json:"confirm_password" binding:"required,min=8"`
	Role            UserRole `json:"role" binding:"required,oneof=admin artist buyer moderator producer fan"`
}

// UpdateUserRequest represents the request payload for updating a user
type UpdateUserRequest struct {
	Email     *string     `json:"email,omitempty" binding:"omitempty,email"`
	Username  *string     `json:"username,omitempty" binding:"omitempty,min=3,max=50"`
	FirstName *string     `json:"first_name,omitempty" binding:"omitempty,min=1,max=100"`
	LastName  *string     `json:"last_name,omitempty" binding:"omitempty,min=1,max=100"`
	Role      *UserRole   `json:"role,omitempty" binding:"omitempty,oneof=admin artist buyer moderator producer fan"`
	Status    *UserStatus `json:"status,omitempty" binding:"omitempty,oneof=active inactive suspended"`
}

// UserResponse represents the response payload for user data
type UserResponse struct {
	ID            int         `json:"id"`
	Email         string      `json:"email"`
	Username      string      `json:"username"`
	FirstName     string      `json:"first_name"`
	LastName      string      `json:"last_name"`
	Role          UserRole    `json:"role"`
	Status        UserStatus  `json:"status"`
	EmailVerified bool        `json:"email_verified"`
	LastLoginAt   *time.Time  `json:"last_login_at"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
}

// ToResponse converts User to UserResponse
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:            u.ID,
		Email:         u.Email,
		Username:      u.Username,
		FirstName:     u.FirstName,
		LastName:      u.LastName,
		Role:          u.Role,
		Status:        u.Status,
		EmailVerified: u.EmailVerified,
		LastLoginAt:   u.LastLoginAt,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
	}
}

// AuthResponse represents the response payload for authentication
type AuthResponse struct {
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	ExpiresAt    time.Time     `json:"expires_at"`
	User         *UserResponse `json:"user"`
}

// LoginRequest represents the request payload for user login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// ForgotPasswordRequest represents the request payload for forgot password
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest represents the request payload for reset password
type ResetPasswordRequest struct {
	Token           string `json:"token" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=8"`
}
