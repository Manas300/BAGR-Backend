package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"bagr-backend/internal/models"
	"bagr-backend/internal/utils"
)

// AuthService handles all authentication operations
type AuthService struct {
	db              *sql.DB
	jwtService      *JWTService
	passwordService *PasswordService
	emailService    *EmailService
}

// NewAuthService creates a new authentication service
func NewAuthService(db *sql.DB, jwtService *JWTService, passwordService *PasswordService, emailService *EmailService) *AuthService {
	return &AuthService{
		db:              db,
		jwtService:      jwtService,
		passwordService: passwordService,
		emailService:    emailService,
	}
}

// GetJWTService returns the JWT service instance
func (a *AuthService) GetJWTService() *JWTService {
	return a.jwtService
}

// RegisterUser handles user registration
func (a *AuthService) RegisterUser(req *models.CreateUserRequest) (*models.AuthResponse, error) {
	logger := utils.GetLogger()

	logger.WithFields(map[string]interface{}{
		"email":    req.Email,
		"username": req.Username,
		"role":     req.Role,
	}).Info("Starting user registration process")

	// Check if user already exists
	logger.Debug("Checking if email already exists in database")
	exists, err := a.userExistsByEmail(req.Email)
	if err != nil {
		logger.WithError(err).Error("Failed to check email existence in database")
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}
	if exists {
		logger.WithField("email", req.Email).Error("Registration failed: email already exists")
		return nil, errors.New("user with this email already exists")
	}
	logger.Debug("Email is available")

	// Check if username already exists
	logger.Debug("Checking if username already exists in database")
	exists, err = a.userExistsByUsername(req.Username)
	if err != nil {
		logger.WithError(err).Error("Failed to check username existence in database")
		return nil, fmt.Errorf("failed to check username existence: %w", err)
	}
	if exists {
		logger.WithField("username", req.Username).Error("Registration failed: username already taken")
		return nil, errors.New("username already taken")
	}
	logger.Debug("Username is available")

	// Hash password
	logger.Debug("Hashing password")
	hashedPassword, err := a.passwordService.HashPassword(req.Password)
	if err != nil {
		logger.WithError(err).Error("Password hashing failed")
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	logger.Debug("Password hashed successfully")

	// Generate verification token
	logger.Debug("Generating verification token")
	verificationToken, err := a.passwordService.GenerateResetToken()
	if err != nil {
		logger.WithError(err).Error("Failed to generate verification token")
		return nil, fmt.Errorf("failed to generate verification token: %w", err)
	}
	logger.WithField("token_length", len(verificationToken)).Debug("Verification token generated")

	// Create user
	logger.Debug("Creating user object")
	user := &models.User{
		Email:             req.Email,
		Username:          req.Username,
		FirstName:         req.FirstName,
		LastName:          req.LastName,
		PasswordHash:      hashedPassword,
		Role:              req.Role,
		Status:            models.UserStatusActive,
		EmailVerified:     false,
		VerificationToken: &verificationToken,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	// Insert user into database
	logger.Debug("Inserting user into database")
	userID, err := a.insertUser(user)
	if err != nil {
		logger.WithError(err).Error("Failed to insert user into database")
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	user.ID = userID
	logger.WithField("user_id", userID).Info("User created successfully in database")

	// Store verification token
	logger.Debug("Storing verification token in database")
	err = a.storeVerificationToken(userID, verificationToken)
	if err != nil {
		logger.WithError(err).Error("Failed to store verification token")
		return nil, fmt.Errorf("failed to store verification token: %w", err)
	}
	logger.Debug("Verification token stored successfully")

	// Send verification email
	logger.Debug("Sending verification email")
	err = a.emailService.SendVerificationEmail(user.Email, user.Username, verificationToken)
	if err != nil {
		logger.WithError(err).Error("Failed to send verification email")
		// Log error but don't fail registration
		fmt.Printf("Warning: Failed to send verification email: %v\n", err)
	} else {
		logger.Info("Verification email sent successfully")
	}

	// Generate tokens
	logger.Debug("Generating JWT tokens")
	accessToken, refreshToken, expiresAt, err := a.jwtService.GenerateTokenPair(user)
	if err != nil {
		logger.WithError(err).Error("Failed to generate JWT tokens")
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}
	logger.Debug("JWT tokens generated successfully")

	logger.WithFields(map[string]interface{}{
		"user_id": userID,
		"email":   user.Email,
		"role":    user.Role,
	}).Info("User registration completed successfully")

	return &models.AuthResponse{
		User:         user.ToResponse(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

// LoginUser handles user login
func (a *AuthService) LoginUser(req *models.LoginRequest) (*models.AuthResponse, error) {
	// Get user by email
	user, err := a.getUserByEmail(req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("invalid email or password")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Check if user is active
	if user.Status != models.UserStatusActive {
		return nil, errors.New("account is not active")
	}

	// Verify password
	err = a.passwordService.VerifyPassword(user.PasswordHash, req.Password)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Check if email is verified
	if !user.EmailVerified {
		return nil, errors.New("please verify your email before logging in")
	}

	// Update last login time
	err = a.updateLastLogin(user.ID)
	if err != nil {
		// Log error but don't fail login
		fmt.Printf("Warning: Failed to update last login time: %v\n", err)
	}

	// Generate tokens
	accessToken, refreshToken, expiresAt, err := a.jwtService.GenerateTokenPair(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return &models.AuthResponse{
		User:         user.ToResponse(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

// VerifyEmail handles email verification
func (a *AuthService) VerifyEmail(token string) (*models.User, error) {
	// Get verification record
	userID, err := a.getVerificationUserID(token)
	if err != nil {
		return nil, fmt.Errorf("invalid or expired verification token")
	}

	// Update user email verification status
	err = a.updateEmailVerification(userID, true)
	if err != nil {
		return nil, fmt.Errorf("failed to verify email: %w", err)
	}

	// Mark verification token as used
	err = a.markVerificationTokenUsed(token)
	if err != nil {
		// Log error but don't fail verification
		fmt.Printf("Warning: Failed to mark verification token as used: %v\n", err)
	}

	// Get user and send welcome email
	user, err := a.getUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Send welcome email
	a.emailService.SendWelcomeEmail(user.Email, user.Username, string(user.Role))

	return user, nil
}

// ForgotPassword handles password reset request
func (a *AuthService) ForgotPassword(req *models.ForgotPasswordRequest) error {
	// Get user by email
	user, err := a.getUserByEmail(req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			// Don't reveal if email exists or not
			return nil
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Generate reset token
	resetToken, err := a.passwordService.GenerateResetToken()
	if err != nil {
		return fmt.Errorf("failed to generate reset token: %w", err)
	}

	// Store reset token
	expiresAt := time.Now().Add(1 * time.Hour) // 1 hour expiry
	err = a.storeResetToken(user.ID, resetToken, expiresAt)
	if err != nil {
		return fmt.Errorf("failed to store reset token: %w", err)
	}

	// Send reset email
	err = a.emailService.SendPasswordResetEmail(user.Email, user.Username, resetToken)
	if err != nil {
		return fmt.Errorf("failed to send reset email: %w", err)
	}

	return nil
}

// ResetPassword handles password reset
func (a *AuthService) ResetPassword(req *models.ResetPasswordRequest) error {
	// Get user by reset token
	userID, err := a.getResetTokenUserID(req.Token)
	if err != nil {
		return fmt.Errorf("invalid or expired reset token")
	}

	// Hash new password
	hashedPassword, err := a.passwordService.HashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update password
	err = a.updatePassword(userID, hashedPassword)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Mark reset token as used
	err = a.markResetTokenUsed(req.Token)
	if err != nil {
		// Log error but don't fail reset
		fmt.Printf("Warning: Failed to mark reset token as used: %v\n", err)
	}

	return nil
}

// RefreshToken handles token refresh
func (a *AuthService) RefreshToken(refreshToken string) (*models.AuthResponse, error) {
	// Validate refresh token
	claims, err := a.jwtService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Get user
	user, err := a.getUserByID(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Check if user is still active
	if user.Status != models.UserStatusActive {
		return nil, errors.New("account is not active")
	}

	// Generate new access token
	accessToken, expiresAt, err := a.jwtService.RefreshAccessToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	return &models.AuthResponse{
		User:         user.ToResponse(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken, // Keep the same refresh token
		ExpiresAt:    expiresAt,
	}, nil
}

// Helper methods for database operations

func (a *AuthService) userExistsByEmail(email string) (bool, error) {
	var count int
	err := a.db.QueryRow("SELECT COUNT(*) FROM users WHERE email = $1", email).Scan(&count)
	return count > 0, err
}

func (a *AuthService) userExistsByUsername(username string) (bool, error) {
	var count int
	err := a.db.QueryRow("SELECT COUNT(*) FROM users WHERE username = $1", username).Scan(&count)
	return count > 0, err
}

func (a *AuthService) insertUser(user *models.User) (int, error) {
	query := `
		INSERT INTO users (email, username, first_name, last_name, password_hash, role, status, email_verified, verification_token, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id`

	var userID int
	err := a.db.QueryRow(query,
		user.Email, user.Username, user.FirstName, user.LastName,
		user.PasswordHash, user.Role, user.Status, user.EmailVerified,
		user.VerificationToken, user.CreatedAt, user.UpdatedAt,
	).Scan(&userID)

	return userID, err
}

func (a *AuthService) getUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, username, first_name, last_name, password_hash, role, status, 
		       email_verified, verification_token, reset_token, reset_token_expires, 
		       last_login_at, created_at, updated_at
		FROM users WHERE email = $1`

	err := a.db.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.Username, &user.FirstName, &user.LastName,
		&user.PasswordHash, &user.Role, &user.Status, &user.EmailVerified,
		&user.VerificationToken, &user.ResetToken, &user.ResetTokenExpires,
		&user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt,
	)

	return user, err
}

func (a *AuthService) getUserByID(id int) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, username, first_name, last_name, password_hash, role, status, 
		       email_verified, verification_token, reset_token, reset_token_expires, 
		       last_login_at, created_at, updated_at
		FROM users WHERE id = $1`

	err := a.db.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.Username, &user.FirstName, &user.LastName,
		&user.PasswordHash, &user.Role, &user.Status, &user.EmailVerified,
		&user.VerificationToken, &user.ResetToken, &user.ResetTokenExpires,
		&user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt,
	)

	return user, err
}

func (a *AuthService) updateLastLogin(userID int) error {
	query := "UPDATE users SET last_login_at = $1, updated_at = $2 WHERE id = $3"
	_, err := a.db.Exec(query, time.Now(), time.Now(), userID)
	return err
}

func (a *AuthService) updateEmailVerification(userID int, verified bool) error {
	query := "UPDATE users SET email_verified = $1, updated_at = $2 WHERE id = $3"
	_, err := a.db.Exec(query, verified, time.Now(), userID)
	return err
}

func (a *AuthService) updatePassword(userID int, hashedPassword string) error {
	query := "UPDATE users SET password_hash = $1, updated_at = $2 WHERE id = $3"
	_, err := a.db.Exec(query, hashedPassword, time.Now(), userID)
	return err
}

func (a *AuthService) storeVerificationToken(userID int, token string) error {
	query := `
		INSERT INTO email_verifications (user_id, token, expires_at)
		VALUES ($1, $2, $3)`

	expiresAt := time.Now().Add(24 * time.Hour) // 24 hours expiry
	_, err := a.db.Exec(query, userID, token, expiresAt)
	return err
}

func (a *AuthService) getVerificationUserID(token string) (int, error) {
	var userID int
	var expiresAt time.Time

	query := `
		SELECT user_id, expires_at 
		FROM email_verifications 
		WHERE token = $1 AND verified_at IS NULL`

	err := a.db.QueryRow(query, token).Scan(&userID, &expiresAt)
	if err != nil {
		return 0, err
	}

	// Check if token is expired
	if time.Now().After(expiresAt) {
		return 0, errors.New("token expired")
	}

	return userID, nil
}

func (a *AuthService) markVerificationTokenUsed(token string) error {
	query := "UPDATE email_verifications SET verified_at = $1 WHERE token = $2"
	_, err := a.db.Exec(query, time.Now(), token)
	return err
}

func (a *AuthService) storeResetToken(userID int, token string, expiresAt time.Time) error {
	query := `
		INSERT INTO password_resets (user_id, token, expires_at)
		VALUES ($1, $2, $3)`

	_, err := a.db.Exec(query, userID, token, expiresAt)
	return err
}

func (a *AuthService) getResetTokenUserID(token string) (int, error) {
	var userID int
	var expiresAt time.Time

	query := `
		SELECT user_id, expires_at 
		FROM password_resets 
		WHERE token = $1 AND used_at IS NULL`

	err := a.db.QueryRow(query, token).Scan(&userID, &expiresAt)
	if err != nil {
		return 0, err
	}

	// Check if token is expired
	if time.Now().After(expiresAt) {
		return 0, errors.New("token expired")
	}

	return userID, nil
}

func (a *AuthService) markResetTokenUsed(token string) error {
	query := "UPDATE password_resets SET used_at = $1 WHERE token = $2"
	_, err := a.db.Exec(query, time.Now(), token)
	return err
}
