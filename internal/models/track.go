package models

import (
	"fmt"
	"time"
)

// Track represents a music track in the system
type Track struct {
	ID          int         `json:"id" db:"id"`
	ArtistID    int         `json:"artist_id" db:"artist_id"`
	Title       string      `json:"title" db:"title"`
	Genre       string      `json:"genre" db:"genre"`
	Duration    int         `json:"duration" db:"duration"` // Duration in seconds
	FileURL     string      `json:"file_url" db:"file_url"`
	CoverArtURL *string     `json:"cover_art_url,omitempty" db:"cover_art_url"`
	Description *string     `json:"description,omitempty" db:"description"`
	Status      TrackStatus `json:"status" db:"status"`
	CreatedAt   time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at" db:"updated_at"`
	
	// Related entities (loaded via joins)
	Artist   *User     `json:"artist,omitempty"`
	Auctions []Auction `json:"auctions,omitempty"`
}

// TrackStatus represents track status
type TrackStatus string

const (
	TrackStatusDraft     TrackStatus = "draft"
	TrackStatusActive    TrackStatus = "active"
	TrackStatusInactive  TrackStatus = "inactive"
	TrackStatusDeleted   TrackStatus = "deleted"
)

// CreateTrackRequest represents the request payload for creating a track
type CreateTrackRequest struct {
	Title       string      `json:"title" binding:"required,min=1,max=200"`
	Genre       string      `json:"genre" binding:"required,min=1,max=100"`
	Duration    int         `json:"duration" binding:"required,min=1"`
	FileURL     string      `json:"file_url" binding:"required,url"`
	CoverArtURL *string     `json:"cover_art_url,omitempty" binding:"omitempty,url"`
	Description *string     `json:"description,omitempty" binding:"omitempty,max=1000"`
}

// UpdateTrackRequest represents the request payload for updating a track
type UpdateTrackRequest struct {
	Title       *string      `json:"title,omitempty" binding:"omitempty,min=1,max=200"`
	Genre       *string      `json:"genre,omitempty" binding:"omitempty,min=1,max=100"`
	Duration    *int         `json:"duration,omitempty" binding:"omitempty,min=1"`
	FileURL     *string      `json:"file_url,omitempty" binding:"omitempty,url"`
	CoverArtURL *string      `json:"cover_art_url,omitempty" binding:"omitempty,url"`
	Description *string      `json:"description,omitempty" binding:"omitempty,max=1000"`
	Status      *TrackStatus `json:"status,omitempty" binding:"omitempty,oneof=draft active inactive deleted"`
}

// TrackResponse represents the response payload for track data
type TrackResponse struct {
	ID          int         `json:"id"`
	ArtistID    int         `json:"artist_id"`
	Title       string      `json:"title"`
	Genre       string      `json:"genre"`
	Duration    int         `json:"duration"`
	FileURL     string      `json:"file_url"`
	CoverArtURL *string     `json:"cover_art_url,omitempty"`
	Description *string     `json:"description,omitempty"`
	Status      TrackStatus `json:"status"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// ToResponse converts Track to TrackResponse
func (t *Track) ToResponse() *TrackResponse {
	return &TrackResponse{
		ID:          t.ID,
		ArtistID:    t.ArtistID,
		Title:       t.Title,
		Genre:       t.Genre,
		Duration:    t.Duration,
		FileURL:     t.FileURL,
		CoverArtURL: t.CoverArtURL,
		Description: t.Description,
		Status:      t.Status,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

// GetDurationFormatted returns the duration in MM:SS format
func (t *Track) GetDurationFormatted() string {
	minutes := t.Duration / 60
	seconds := t.Duration % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}
