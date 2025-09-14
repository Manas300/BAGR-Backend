package models

import (
	"time"
)

// Auction represents an auction in the system
type Auction struct {
	ID          int           `json:"id" db:"id"`
	TrackID     int           `json:"track_id" db:"track_id"`
	SellerID    int           `json:"seller_id" db:"seller_id"`
	Title       string        `json:"title" db:"title"`
	Description string        `json:"description" db:"description"`
	StartPrice  float64       `json:"start_price" db:"start_price"`
	ReservePrice *float64     `json:"reserve_price,omitempty" db:"reserve_price"`
	CurrentBid  *float64      `json:"current_bid,omitempty" db:"current_bid"`
	BidCount    int           `json:"bid_count" db:"bid_count"`
	Status      AuctionStatus `json:"status" db:"status"`
	StartTime   time.Time     `json:"start_time" db:"start_time"`
	EndTime     time.Time     `json:"end_time" db:"end_time"`
	CreatedAt   time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at" db:"updated_at"`
	
	// Related entities (loaded via joins)
	Track  *Track `json:"track,omitempty"`
	Seller *User  `json:"seller,omitempty"`
	Bids   []Bid  `json:"bids,omitempty"`
}

// AuctionStatus represents auction status
type AuctionStatus string

const (
	AuctionStatusDraft     AuctionStatus = "draft"
	AuctionStatusActive    AuctionStatus = "active"
	AuctionStatusCompleted AuctionStatus = "completed"
	AuctionStatusCancelled AuctionStatus = "cancelled"
	AuctionStatusExpired   AuctionStatus = "expired"
)

// CreateAuctionRequest represents the request payload for creating an auction
type CreateAuctionRequest struct {
	TrackID      int       `json:"track_id" binding:"required"`
	Title        string    `json:"title" binding:"required,min=1,max=200"`
	Description  string    `json:"description" binding:"required,min=1,max=1000"`
	StartPrice   float64   `json:"start_price" binding:"required,min=0"`
	ReservePrice *float64  `json:"reserve_price,omitempty" binding:"omitempty,min=0"`
	StartTime    time.Time `json:"start_time" binding:"required"`
	EndTime      time.Time `json:"end_time" binding:"required"`
}

// UpdateAuctionRequest represents the request payload for updating an auction
type UpdateAuctionRequest struct {
	Title        *string        `json:"title,omitempty" binding:"omitempty,min=1,max=200"`
	Description  *string        `json:"description,omitempty" binding:"omitempty,min=1,max=1000"`
	StartPrice   *float64       `json:"start_price,omitempty" binding:"omitempty,min=0"`
	ReservePrice *float64       `json:"reserve_price,omitempty" binding:"omitempty,min=0"`
	Status       *AuctionStatus `json:"status,omitempty" binding:"omitempty,oneof=draft active completed cancelled expired"`
	StartTime    *time.Time     `json:"start_time,omitempty"`
	EndTime      *time.Time     `json:"end_time,omitempty"`
}

// IsActive returns true if the auction is currently active
func (a *Auction) IsActive() bool {
	now := time.Now()
	return a.Status == AuctionStatusActive && 
		   now.After(a.StartTime) && 
		   now.Before(a.EndTime)
}

// IsExpired returns true if the auction has expired
func (a *Auction) IsExpired() bool {
	return time.Now().After(a.EndTime)
}

// HasReserveMet returns true if the current bid meets the reserve price
func (a *Auction) HasReserveMet() bool {
	if a.ReservePrice == nil {
		return true // No reserve price set
	}
	if a.CurrentBid == nil {
		return false // No bids yet
	}
	return *a.CurrentBid >= *a.ReservePrice
}
