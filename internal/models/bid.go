package models

import (
	"time"
)

// Bid represents a bid in an auction
type Bid struct {
	ID        int       `json:"id" db:"id"`
	AuctionID int       `json:"auction_id" db:"auction_id"`
	BidderID  int       `json:"bidder_id" db:"bidder_id"`
	Amount    float64   `json:"amount" db:"amount"`
	Status    BidStatus `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	
	// Related entities (loaded via joins)
	Auction *Auction `json:"auction,omitempty"`
	Bidder  *User    `json:"bidder,omitempty"`
}

// BidStatus represents bid status
type BidStatus string

const (
	BidStatusActive    BidStatus = "active"
	BidStatusOutbid    BidStatus = "outbid"
	BidStatusWinning   BidStatus = "winning"
	BidStatusCancelled BidStatus = "cancelled"
)

// CreateBidRequest represents the request payload for creating a bid
type CreateBidRequest struct {
	AuctionID int     `json:"auction_id" binding:"required"`
	Amount    float64 `json:"amount" binding:"required,min=0"`
}

// BidResponse represents the response payload for bid data
type BidResponse struct {
	ID        int       `json:"id"`
	AuctionID int       `json:"auction_id"`
	BidderID  int       `json:"bidder_id"`
	Amount    float64   `json:"amount"`
	Status    BidStatus `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToResponse converts Bid to BidResponse
func (b *Bid) ToResponse() *BidResponse {
	return &BidResponse{
		ID:        b.ID,
		AuctionID: b.AuctionID,
		BidderID:  b.BidderID,
		Amount:    b.Amount,
		Status:    b.Status,
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
	}
}
