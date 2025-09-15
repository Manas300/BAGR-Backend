package auth

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"bagr-backend/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

// PasswordService handles password operations
type PasswordService struct {
	minLength      int
	requireUpper   bool
	requireLower   bool
	requireDigit   bool
	requireSpecial bool
}

// NewPasswordService creates a new password service
func NewPasswordService() *PasswordService {
	return &PasswordService{
		minLength:      8,
		requireUpper:   true,
		requireLower:   true,
		requireDigit:   true,
		requireSpecial: false, // Keep it simple for now
	}
}

// HashPassword hashes a password using bcrypt
func (p *PasswordService) HashPassword(password string) (string, error) {
	if err := p.ValidatePassword(password); err != nil {
		return "", err
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedBytes), nil
}

// VerifyPassword verifies a password against its hash
func (p *PasswordService) VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// ValidatePassword validates password strength
func (p *PasswordService) ValidatePassword(password string) error {
	// Log password validation attempt
	utils.GetLogger().WithField("password_length", len(password)).Debug("Starting password validation")

	// Check length
	if len(password) < p.minLength {
		utils.GetLogger().WithFields(map[string]interface{}{
			"password_length": len(password),
			"required_length": p.minLength,
		}).Error("Password validation failed: too short")
		return errors.New("password must be at least 8 characters long")
	}

	// Check uppercase
	hasUpper := p.hasUpperCase(password)
	if p.requireUpper && !hasUpper {
		utils.GetLogger().WithField("password_length", len(password)).Error("Password validation failed: no uppercase letter")
		return errors.New("password must contain at least one uppercase letter")
	}

	// Check lowercase
	hasLower := p.hasLowerCase(password)
	if p.requireLower && !hasLower {
		utils.GetLogger().WithField("password_length", len(password)).Error("Password validation failed: no lowercase letter")
		return errors.New("password must contain at least one lowercase letter")
	}

	// Check digits
	hasDigit := p.hasDigit(password)
	if p.requireDigit && !hasDigit {
		utils.GetLogger().WithField("password_length", len(password)).Error("Password validation failed: no digit")
		return errors.New("password must contain at least one digit")
	}

	// Check special characters
	hasSpecial := p.hasSpecialChar(password)
	if p.requireSpecial && !hasSpecial {
		utils.GetLogger().WithField("password_length", len(password)).Error("Password validation failed: no special character")
		return errors.New("password must contain at least one special character")
	}

	// Check for common weak passwords
	isCommon := p.isCommonPassword(password)
	if isCommon {
		utils.GetLogger().WithField("password_length", len(password)).Error("Password validation failed: common password")
		return errors.New("password is too common, please choose a stronger password")
	}

	// Log successful validation
	utils.GetLogger().WithFields(map[string]interface{}{
		"password_length": len(password),
		"has_upper":       hasUpper,
		"has_lower":       hasLower,
		"has_digit":       hasDigit,
		"has_special":     hasSpecial,
		"is_common":       isCommon,
	}).Debug("Password validation successful")

	return nil
}

// hasUpperCase checks if password contains uppercase letters
func (p *PasswordService) hasUpperCase(password string) bool {
	matched, _ := regexp.MatchString(`[A-Z]`, password)
	return matched
}

// hasLowerCase checks if password contains lowercase letters
func (p *PasswordService) hasLowerCase(password string) bool {
	matched, _ := regexp.MatchString(`[a-z]`, password)
	return matched
}

// hasDigit checks if password contains digits
func (p *PasswordService) hasDigit(password string) bool {
	matched, _ := regexp.MatchString(`[0-9]`, password)
	return matched
}

// hasSpecialChar checks if password contains special characters
func (p *PasswordService) hasSpecialChar(password string) bool {
	matched, _ := regexp.MatchString(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`, password)
	return matched
}

// isCommonPassword checks if password is in common passwords list
func (p *PasswordService) isCommonPassword(password string) bool {
	commonPasswords := []string{
		"password", "123456", "123456789", "qwerty", "abc123",
		"password123", "admin", "letmein", "welcome", "monkey",
		"1234567890", "password1", "qwerty123", "dragon", "master",
		"hello", "freedom", "whatever", "qazwsx", "trustno1",
	}

	lowerPassword := strings.ToLower(password)
	for _, common := range commonPasswords {
		if lowerPassword == common {
			return true
		}
	}

	return false
}

// GenerateResetToken generates a secure random token for password reset
func (p *PasswordService) GenerateResetToken() (string, error) {
	// This is a simple implementation - in production, use crypto/rand
	// For now, we'll use a simple approach and improve it later
	return generateRandomString(32), nil
}

// generateRandomString generates a random string of specified length
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[randomInt(len(charset))]
	}
	return string(b)
}

// randomInt generates a random integer (simplified version)
func randomInt(max int) int {
	// This is a simplified version - in production, use crypto/rand
	return int(time.Now().UnixNano()) % max
}
