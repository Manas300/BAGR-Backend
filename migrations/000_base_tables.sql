-- Migration: Create base tables
-- Created: 2024-01-01
-- Description: Creates the base users table and other core tables

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(50) UNIQUE NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    password VARCHAR(255), -- Legacy field, will be replaced by password_hash
    role VARCHAR(20) NOT NULL DEFAULT 'fan',
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Create auctions table
CREATE TABLE IF NOT EXISTS auctions (
    id SERIAL PRIMARY KEY,
    track_id INTEGER,
    seller_id INTEGER NOT NULL REFERENCES users(id),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    start_price DECIMAL(10,2) NOT NULL,
    reserve_price DECIMAL(10,2),
    current_bid DECIMAL(10,2) DEFAULT 0,
    bid_count INTEGER DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    start_time TIMESTAMP DEFAULT NOW(),
    end_time TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Create bids table
CREATE TABLE IF NOT EXISTS bids (
    id SERIAL PRIMARY KEY,
    auction_id INTEGER NOT NULL REFERENCES auctions(id),
    bidder_id INTEGER NOT NULL REFERENCES users(id),
    amount DECIMAL(10,2) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP DEFAULT NOW()
);

-- Create tracks table
CREATE TABLE IF NOT EXISTS tracks (
    id SERIAL PRIMARY KEY,
    artist_id INTEGER NOT NULL REFERENCES users(id),
    title VARCHAR(255) NOT NULL,
    genre VARCHAR(100),
    duration INTEGER, -- in seconds
    file_url VARCHAR(500),
    cover_art_url VARCHAR(500),
    description TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_auctions_seller_id ON auctions(seller_id);
CREATE INDEX IF NOT EXISTS idx_auctions_status ON auctions(status);
CREATE INDEX IF NOT EXISTS idx_bids_auction_id ON bids(auction_id);
CREATE INDEX IF NOT EXISTS idx_bids_bidder_id ON bids(bidder_id);
CREATE INDEX IF NOT EXISTS idx_tracks_artist_id ON tracks(artist_id);

-- Add comments for documentation
COMMENT ON TABLE users IS 'Users table for the auction system';
COMMENT ON TABLE auctions IS 'Auctions table for music auctions';
COMMENT ON TABLE bids IS 'Bids table for auction bids';
COMMENT ON TABLE tracks IS 'Tracks table for music tracks';
