package store

import (
	"database/sql"
	"time"

	"github.com/baskint/bidding-analysis/internal/models"
)

// BidStore handles database operations for bid data
type BidStore struct {
	db *sql.DB
}

// NewBidStore creates a new BidStore instance
func NewBidStore(db *sql.DB) *BidStore {
	return &BidStore{db: db}
}

// StoreBidEvent stores a new bid event in the database
func (s *BidStore) StoreBidEvent(bid *models.BidEvent) error {
	query := `
		INSERT INTO bid_events (
			campaign_id, user_id, bid_price, win_price, floor_price, won, converted,
			segment_id, segment_category, engagement_score, conversion_probability,
			country, region, city, latitude, longitude,
			device_type, os, browser, is_mobile, keywords, timestamp
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22
		) RETURNING id`

	err := s.db.QueryRow(
		query,
		bid.CampaignID, bid.UserID, bid.BidPrice, bid.WinPrice, bid.FloorPrice,
		bid.Won, bid.Converted, bid.SegmentID, bid.SegmentCategory,
		bid.EngagementScore, bid.ConversionProbability, bid.Country, bid.Region,
		bid.City, bid.Latitude, bid.Longitude, bid.DeviceType, bid.OS,
		bid.Browser, bid.IsMobile, bid.Keywords, bid.Timestamp,
	).Scan(&bid.ID)

	return err
}

// GetBidHistory retrieves bid history for a campaign
func (s *BidStore) GetBidHistory(campaignID string, startTime, endTime time.Time, limit, offset int) ([]*models.BidEvent, error) {
	query := `
		SELECT id, campaign_id, user_id, bid_price, win_price, won, converted, timestamp
		FROM bid_events 
		WHERE campaign_id = $1 AND timestamp BETWEEN $2 AND $3
		ORDER BY timestamp DESC
		LIMIT $4 OFFSET $5`

	rows, err := s.db.Query(query, campaignID, startTime, endTime, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bids []*models.BidEvent
	for rows.Next() {
		bid := &models.BidEvent{}
		err := rows.Scan(
			&bid.ID, &bid.CampaignID, &bid.UserID, &bid.BidPrice,
			&bid.WinPrice, &bid.Won, &bid.Converted, &bid.Timestamp,
		)
		if err != nil {
			return nil, err
		}
		bids = append(bids, bid)
	}

	return bids, rows.Err()
}

// GetRecentBids gets recent bids for ML training
func (s *BidStore) GetRecentBids(campaignID string, limit int) ([]*models.BidEvent, error) {
	query := `
		SELECT id, campaign_id, bid_price, win_price, won, converted, timestamp
		FROM bid_events 
		WHERE campaign_id = $1 
		ORDER BY timestamp DESC 
		LIMIT $2`

	rows, err := s.db.Query(query, campaignID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bids []*models.BidEvent
	for rows.Next() {
		bid := &models.BidEvent{}
		err := rows.Scan(
			&bid.ID, &bid.CampaignID, &bid.BidPrice,
			&bid.WinPrice, &bid.Won, &bid.Converted, &bid.Timestamp,
		)
		if err != nil {
			return nil, err
		}
		bids = append(bids, bid)
	}

	return bids, rows.Err()
}
