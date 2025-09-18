package models

import (
	"time"
)

// Profile represents a user's profile information
type Profile struct {
	ID              int        `json:"id" db:"id"`
	UserID          int        `json:"user_id" db:"user_id"`
	DisplayName     string     `json:"display_name" db:"display_name"`
	Bio             *string    `json:"bio" db:"bio"`
	Location        *string    `json:"location" db:"location"`
	ProfileImageURL *string    `json:"profile_image_url" db:"profile_image_url"`
	WebsiteURL      *string    `json:"website_url" db:"website_url"`
	YouTubeHandle   *string    `json:"youtube_handle" db:"youtube_handle"`
	TikTokHandle    *string    `json:"tiktok_handle" db:"tiktok_handle"`
	InstagramHandle *string    `json:"instagram_handle" db:"instagram_handle"`
	TwitterHandle   *string    `json:"twitter_handle" db:"twitter_handle"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
}

// CreateProfileRequest represents the request payload for creating a profile
type CreateProfileRequest struct {
	DisplayName     string `json:"display_name" binding:"required,min=1,max=100"`
	Bio             string `json:"bio" binding:"max=500"`
	Location        string `json:"location" binding:"max=100"`
	WebsiteURL      string `json:"website_url" binding:"omitempty,url"`
	YouTubeHandle   string `json:"youtube_handle" binding:"max=50"`
	TikTokHandle    string `json:"tiktok_handle" binding:"max=50"`
	InstagramHandle string `json:"instagram_handle" binding:"max=50"`
	TwitterHandle   string `json:"twitter_handle" binding:"max=50"`
}

// UpdateProfileRequest represents the request payload for updating a profile
type UpdateProfileRequest struct {
	DisplayName     *string `json:"display_name,omitempty" binding:"omitempty,min=1,max=100"`
	Bio             *string `json:"bio,omitempty" binding:"omitempty,max=500"`
	Location        *string `json:"location,omitempty" binding:"omitempty,max=100"`
	WebsiteURL      *string `json:"website_url,omitempty" binding:"omitempty,max=255"`
	YouTubeHandle   *string `json:"youtube_handle,omitempty" binding:"omitempty,max=50"`
	TikTokHandle    *string `json:"tiktok_handle,omitempty" binding:"omitempty,max=50"`
	InstagramHandle *string `json:"instagram_handle,omitempty" binding:"omitempty,max=50"`
	TwitterHandle   *string `json:"twitter_handle,omitempty" binding:"omitempty,max=50"`
}

// ProfileResponse represents the response payload for profile data
type ProfileResponse struct {
	ID              int    `json:"id"`
	UserID          int    `json:"user_id"`
	DisplayName     string `json:"display_name"`
	Bio             string `json:"bio"`
	Location        string `json:"location"`
	ProfileImageURL string `json:"profile_image_url"`
	WebsiteURL      string `json:"website_url"`
	YouTubeHandle   string `json:"youtube_handle"`
	TikTokHandle    string `json:"tiktok_handle"`
	InstagramHandle string `json:"instagram_handle"`
	TwitterHandle   string `json:"twitter_handle"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

// ToResponse converts Profile to ProfileResponse
func (p *Profile) ToResponse() *ProfileResponse {
	return &ProfileResponse{
		ID:              p.ID,
		UserID:          p.UserID,
		DisplayName:     p.DisplayName,
		Bio:             getStringValue(p.Bio),
		Location:        getStringValue(p.Location),
		ProfileImageURL: getStringValue(p.ProfileImageURL),
		WebsiteURL:      getStringValue(p.WebsiteURL),
		YouTubeHandle:   getStringValue(p.YouTubeHandle),
		TikTokHandle:    getStringValue(p.TikTokHandle),
		InstagramHandle: getStringValue(p.InstagramHandle),
		TwitterHandle:   getStringValue(p.TwitterHandle),
		CreatedAt:       p.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       p.UpdatedAt.Format(time.RFC3339),
	}
}

// Helper function to get string value from pointer
func getStringValue(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}
