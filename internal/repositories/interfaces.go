package repositories

import (
	"context"
	"bagr-backend/internal/models"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id int) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	Update(ctx context.Context, id int, updates map[string]interface{}) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, limit, offset int) ([]*models.User, error)
}

// AuctionRepository defines the interface for auction data access
type AuctionRepository interface {
	Create(ctx context.Context, auction *models.Auction) error
	GetByID(ctx context.Context, id int) (*models.Auction, error)
	Update(ctx context.Context, id int, updates map[string]interface{}) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, limit, offset int) ([]*models.Auction, error)
	GetBySellerID(ctx context.Context, sellerID int, limit, offset int) ([]*models.Auction, error)
	GetActiveAuctions(ctx context.Context, limit, offset int) ([]*models.Auction, error)
	UpdateCurrentBid(ctx context.Context, auctionID int, bidAmount float64) error
}

// BidRepository defines the interface for bid data access
type BidRepository interface {
	Create(ctx context.Context, bid *models.Bid) error
	GetByID(ctx context.Context, id int) (*models.Bid, error)
	Update(ctx context.Context, id int, updates map[string]interface{}) error
	Delete(ctx context.Context, id int) error
	GetByAuctionID(ctx context.Context, auctionID int, limit, offset int) ([]*models.Bid, error)
	GetByBidderID(ctx context.Context, bidderID int, limit, offset int) ([]*models.Bid, error)
	GetHighestBidForAuction(ctx context.Context, auctionID int) (*models.Bid, error)
	GetBidHistory(ctx context.Context, auctionID int) ([]*models.Bid, error)
}

// TrackRepository defines the interface for track data access
type TrackRepository interface {
	Create(ctx context.Context, track *models.Track) error
	GetByID(ctx context.Context, id int) (*models.Track, error)
	Update(ctx context.Context, id int, updates map[string]interface{}) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, limit, offset int) ([]*models.Track, error)
	GetByArtistID(ctx context.Context, artistID int, limit, offset int) ([]*models.Track, error)
	Search(ctx context.Context, query string, limit, offset int) ([]*models.Track, error)
}

// Repositories holds all repository interfaces
type Repositories struct {
	User    UserRepository
	Auction AuctionRepository
	Bid     BidRepository
	Track   TrackRepository
}
