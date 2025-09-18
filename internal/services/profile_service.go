package services

import (
	"database/sql"
	"fmt"
	"time"

	"bagr-backend/internal/models"

	"github.com/sirupsen/logrus"
)

// ProfileService handles profile-related business logic
type ProfileService struct {
	db     *sql.DB
	logger *logrus.Logger
}

// NewProfileService creates a new profile service
func NewProfileService(db *sql.DB, logger *logrus.Logger) *ProfileService {
	return &ProfileService{
		db:     db,
		logger: logger,
	}
}

// GetProfileByUserID retrieves a profile by user ID
func (s *ProfileService) GetProfileByUserID(userID int) (*models.Profile, error) {
	query := `
		SELECT id, user_id, display_name, bio, location, profile_image_url, 
		       website_url, youtube_handle, tiktok_handle, instagram_handle, 
		       twitter_handle, created_at, updated_at
		FROM profiles 
		WHERE user_id = $1
	`

	var profile models.Profile
	err := s.db.QueryRow(query, userID).Scan(
		&profile.ID,
		&profile.UserID,
		&profile.DisplayName,
		&profile.Bio,
		&profile.Location,
		&profile.ProfileImageURL,
		&profile.WebsiteURL,
		&profile.YouTubeHandle,
		&profile.TikTokHandle,
		&profile.InstagramHandle,
		&profile.TwitterHandle,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("profile not found for user %d", userID)
		}
		s.logger.WithError(err).WithField("user_id", userID).Error("Failed to get profile")
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}

	return &profile, nil
}

// CreateProfile creates a new profile for a user
func (s *ProfileService) CreateProfile(userID int, req *models.CreateProfileRequest) (*models.Profile, error) {
	query := `
		INSERT INTO profiles (user_id, display_name, bio, location, website_url, 
		                     youtube_handle, tiktok_handle, instagram_handle, twitter_handle, 
		                     created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, user_id, display_name, bio, location, profile_image_url, 
		          website_url, youtube_handle, tiktok_handle, instagram_handle, 
		          twitter_handle, created_at, updated_at
	`

	now := time.Now()
	var profile models.Profile
	err := s.db.QueryRow(query,
		userID,
		req.DisplayName,
		getNullableString(req.Bio),
		getNullableString(req.Location),
		getNullableString(req.WebsiteURL),
		getNullableString(req.YouTubeHandle),
		getNullableString(req.TikTokHandle),
		getNullableString(req.InstagramHandle),
		getNullableString(req.TwitterHandle),
		now,
		now,
	).Scan(
		&profile.ID,
		&profile.UserID,
		&profile.DisplayName,
		&profile.Bio,
		&profile.Location,
		&profile.ProfileImageURL,
		&profile.WebsiteURL,
		&profile.YouTubeHandle,
		&profile.TikTokHandle,
		&profile.InstagramHandle,
		&profile.TwitterHandle,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)

	if err != nil {
		s.logger.WithError(err).WithField("user_id", userID).Error("Failed to create profile")
		return nil, fmt.Errorf("failed to create profile: %w", err)
	}

	s.logger.WithField("user_id", userID).Info("Profile created successfully")
	return &profile, nil
}

// UpdateProfile updates an existing profile
func (s *ProfileService) UpdateProfile(userID int, req *models.UpdateProfileRequest) (*models.Profile, error) {
	// Build dynamic update query
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.DisplayName != nil {
		setParts = append(setParts, fmt.Sprintf("display_name = $%d", argIndex))
		args = append(args, *req.DisplayName)
		argIndex++
	}
	if req.Bio != nil {
		setParts = append(setParts, fmt.Sprintf("bio = $%d", argIndex))
		args = append(args, *req.Bio)
		argIndex++
	}
	if req.Location != nil {
		setParts = append(setParts, fmt.Sprintf("location = $%d", argIndex))
		args = append(args, *req.Location)
		argIndex++
	}
	if req.WebsiteURL != nil {
		setParts = append(setParts, fmt.Sprintf("website_url = $%d", argIndex))
		args = append(args, *req.WebsiteURL)
		argIndex++
	}
	if req.YouTubeHandle != nil {
		setParts = append(setParts, fmt.Sprintf("youtube_handle = $%d", argIndex))
		args = append(args, *req.YouTubeHandle)
		argIndex++
	}
	if req.TikTokHandle != nil {
		setParts = append(setParts, fmt.Sprintf("tiktok_handle = $%d", argIndex))
		args = append(args, *req.TikTokHandle)
		argIndex++
	}
	if req.InstagramHandle != nil {
		setParts = append(setParts, fmt.Sprintf("instagram_handle = $%d", argIndex))
		args = append(args, *req.InstagramHandle)
		argIndex++
	}
	if req.TwitterHandle != nil {
		setParts = append(setParts, fmt.Sprintf("twitter_handle = $%d", argIndex))
		args = append(args, *req.TwitterHandle)
		argIndex++
	}

	if len(setParts) == 0 {
		return s.GetProfileByUserID(userID)
	}

	// Add updated_at
	setParts = append(setParts, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, time.Now())
	argIndex++

	// Add WHERE clause
	args = append(args, userID)

	query := fmt.Sprintf(`
		UPDATE profiles 
		SET %s
		WHERE user_id = $%d
		RETURNING id, user_id, display_name, bio, location, profile_image_url, 
		          website_url, youtube_handle, tiktok_handle, instagram_handle, 
		          twitter_handle, created_at, updated_at
	`, fmt.Sprintf("%s", setParts[0]), argIndex)

	// Fix the query building
	query = fmt.Sprintf(`
		UPDATE profiles 
		SET %s
		WHERE user_id = $%d
		RETURNING id, user_id, display_name, bio, location, profile_image_url, 
		          website_url, youtube_handle, tiktok_handle, instagram_handle, 
		          twitter_handle, created_at, updated_at
	`, fmt.Sprintf("%s", setParts[0]), argIndex)

	// Actually, let's build this properly
	query = fmt.Sprintf(`
		UPDATE profiles 
		SET %s
		WHERE user_id = $%d
		RETURNING id, user_id, display_name, bio, location, profile_image_url, 
		          website_url, youtube_handle, tiktok_handle, instagram_handle, 
		          twitter_handle, created_at, updated_at
	`, fmt.Sprintf("%s", setParts[0]), argIndex)

	// Let me fix this properly
	setClause := ""
	for i, part := range setParts {
		if i > 0 {
			setClause += ", "
		}
		setClause += part
	}

	query = fmt.Sprintf(`
		UPDATE profiles 
		SET %s
		WHERE user_id = $%d
		RETURNING id, user_id, display_name, bio, location, profile_image_url, 
		          website_url, youtube_handle, tiktok_handle, instagram_handle, 
		          twitter_handle, created_at, updated_at
	`, setClause, argIndex)

	var profile models.Profile
	err := s.db.QueryRow(query, args...).Scan(
		&profile.ID,
		&profile.UserID,
		&profile.DisplayName,
		&profile.Bio,
		&profile.Location,
		&profile.ProfileImageURL,
		&profile.WebsiteURL,
		&profile.YouTubeHandle,
		&profile.TikTokHandle,
		&profile.InstagramHandle,
		&profile.TwitterHandle,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)

	if err != nil {
		s.logger.WithError(err).WithField("user_id", userID).Error("Failed to update profile")
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	s.logger.WithField("user_id", userID).Info("Profile updated successfully")
	return &profile, nil
}

// UpdateProfileImage updates the profile image URL
func (s *ProfileService) UpdateProfileImage(userID int, imageURL string) error {
	query := `
		UPDATE profiles 
		SET profile_image_url = $1, updated_at = $2
		WHERE user_id = $3
	`

	_, err := s.db.Exec(query, imageURL, time.Now(), userID)
	if err != nil {
		s.logger.WithError(err).WithField("user_id", userID).Error("Failed to update profile image")
		return fmt.Errorf("failed to update profile image: %w", err)
	}

	s.logger.WithField("user_id", userID).Info("Profile image updated successfully")
	return nil
}

// ProfileExists checks if a profile exists for a user
func (s *ProfileService) ProfileExists(userID int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM profiles WHERE user_id = $1)`

	var exists bool
	err := s.db.QueryRow(query, userID).Scan(&exists)
	if err != nil {
		s.logger.WithError(err).WithField("user_id", userID).Error("Failed to check profile existence")
		return false, fmt.Errorf("failed to check profile existence: %w", err)
	}

	return exists, nil
}

// Helper function to convert empty string to NULL for database
func getNullableString(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
