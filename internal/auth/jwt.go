package auth

import (
	"errors"
	"fmt"
	"time"

	"bagr-backend/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

// JWTService handles JWT token operations
type JWTService struct {
	accessSecret  []byte
	refreshSecret []byte
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

// Claims represents the JWT claims
type Claims struct {
	UserID    int             `json:"user_id"`
	Email     string          `json:"email"`
	Role      models.UserRole `json:"role"`
	TokenType string          `json:"token_type"` // "access" or "refresh"
	jwt.RegisteredClaims
}

// NewJWTService creates a new JWT service
func NewJWTService(accessSecret, refreshSecret string) *JWTService {
	return &JWTService{
		accessSecret:  []byte(accessSecret),
		refreshSecret: []byte(refreshSecret),
		accessExpiry:  24 * time.Hour,     // 24 hours as requested
		refreshExpiry: 7 * 24 * time.Hour, // 7 days for refresh tokens
	}
}

// GenerateTokenPair generates both access and refresh tokens
func (j *JWTService) GenerateTokenPair(user *models.User) (string, string, time.Time, error) {
	now := time.Now()
	accessExpiry := now.Add(j.accessExpiry)
	refreshExpiry := now.Add(j.refreshExpiry)

	// Generate access token
	accessClaims := &Claims{
		UserID:    user.ID,
		Email:     user.Email,
		Role:      user.Role,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpiry),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "bagr-auction-system",
			Subject:   fmt.Sprintf("user:%d", user.ID),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(j.accessSecret)
	if err != nil {
		return "", "", time.Time{}, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token
	refreshClaims := &Claims{
		UserID:    user.ID,
		Email:     user.Email,
		Role:      user.Role,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpiry),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "bagr-auction-system",
			Subject:   fmt.Sprintf("user:%d", user.ID),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(j.refreshSecret)
	if err != nil {
		return "", "", time.Time{}, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return accessTokenString, refreshTokenString, accessExpiry, nil
}

// ValidateAccessToken validates an access token
func (j *JWTService) ValidateAccessToken(tokenString string) (*Claims, error) {
	return j.validateToken(tokenString, j.accessSecret, "access")
}

// ValidateRefreshToken validates a refresh token
func (j *JWTService) ValidateRefreshToken(tokenString string) (*Claims, error) {
	return j.validateToken(tokenString, j.refreshSecret, "refresh")
}

// validateToken validates a JWT token with the given secret and expected type
func (j *JWTService) validateToken(tokenString string, secret []byte, expectedType string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Check if token is expired
	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("token has expired")
	}

	// Check token type
	if claims.TokenType != expectedType {
		return nil, fmt.Errorf("invalid token type: expected %s, got %s", expectedType, claims.TokenType)
	}

	return claims, nil
}

// RefreshAccessToken generates a new access token from a refresh token
func (j *JWTService) RefreshAccessToken(refreshTokenString string) (string, time.Time, error) {
	claims, err := j.ValidateRefreshToken(refreshTokenString)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Create new access token with same user info
	now := time.Now()
	accessExpiry := now.Add(j.accessExpiry)

	accessClaims := &Claims{
		UserID:    claims.UserID,
		Email:     claims.Email,
		Role:      claims.Role,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpiry),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "bagr-auction-system",
			Subject:   fmt.Sprintf("user:%d", claims.UserID),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(j.accessSecret)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to generate new access token: %w", err)
	}

	return accessTokenString, accessExpiry, nil
}

// ExtractUserFromToken extracts user information from a token
func (j *JWTService) ExtractUserFromToken(tokenString string) (*models.User, error) {
	claims, err := j.ValidateAccessToken(tokenString)
	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:    claims.UserID,
		Email: claims.Email,
		Role:  claims.Role,
	}, nil
}

// IsTokenExpired checks if a token is expired
func (j *JWTService) IsTokenExpired(tokenString string) bool {
	_, err := j.ValidateAccessToken(tokenString)
	return err != nil
}

