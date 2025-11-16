package store

import (
	"context"
	"time"

	"github.com/baskint/bidding-analysis/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// BidStore handles database operations for bid data
type BidStore struct {
	db *sqlx.DB
}

// NewBidStore creates a new BidStore instance
func NewBidStore(db *sqlx.DB) *BidStore {
	return &BidStore{db: db}
}

// DB returns the underlying database connection for use by other stores
func (s *BidStore) DB() *sqlx.DB {
	return s.db
}

// StoreBidEvent stores a new bid event in the database
func (s *BidStore) StoreBidEvent(bid *models.BidEvent) error {
	// Generate UUID if not set
	if bid.ID == uuid.Nil {
		bid.ID = uuid.New()
	}

	// Set created_at if not set
	if bid.CreatedAt.IsZero() {
		bid.CreatedAt = time.Now()
	}

	query := `
		INSERT INTO bid_events (
			id, campaign_id, user_id, bid_price, win_price, floor_price, won, converted,
			segment_id, segment_category, engagement_score, conversion_probability,
			country, region, city, latitude, longitude,
			device_type, os, browser, is_mobile, keywords, timestamp, created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24
		)`

	_, err := s.db.Exec(
		query,
		bid.ID, bid.CampaignID, bid.UserID, bid.BidPrice, bid.WinPrice, bid.FloorPrice,
		bid.Won, bid.Converted, bid.SegmentID, bid.SegmentCategory,
		bid.EngagementScore, bid.ConversionProbability, bid.Country, bid.Region,
		bid.City, bid.Latitude, bid.Longitude, bid.DeviceType, bid.OS,
		bid.Browser, bid.IsMobile, bid.Keywords, bid.Timestamp, bid.CreatedAt,
	)

	return err
}

// GetBidHistory retrieves bid history for a campaign
func (s *BidStore) GetBidHistory(campaignID string, startTime, endTime time.Time, limit, offset int) ([]*models.BidEvent, error) {
	query := `
		SELECT id, campaign_id, user_id, bid_price, win_price, floor_price, won, converted, 
			   segment_id, segment_category, engagement_score, conversion_probability,
			   country, region, city, latitude, longitude, device_type, os, browser, 
			   is_mobile, keywords, timestamp, created_at
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
			&bid.WinPrice, &bid.FloorPrice, &bid.Won, &bid.Converted,
			&bid.SegmentID, &bid.SegmentCategory, &bid.EngagementScore, &bid.ConversionProbability,
			&bid.Country, &bid.Region, &bid.City, &bid.Latitude, &bid.Longitude,
			&bid.DeviceType, &bid.OS, &bid.Browser, &bid.IsMobile, &bid.Keywords,
			&bid.Timestamp, &bid.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		bids = append(bids, bid)
	}

	return bids, rows.Err()
}

// GetRecentBids gets recent bids for ML training (updated signature to match tRPC handler)
func (s *BidStore) GetRecentBids(ctx context.Context, limit int) ([]*models.BidEvent, error) {
	query := `
		SELECT id, campaign_id, user_id, bid_price, win_price, floor_price, won, converted,
			   segment_id, segment_category, engagement_score, conversion_probability,
			   country, region, city, latitude, longitude, device_type, os, browser,
			   is_mobile, keywords, timestamp, created_at
		FROM bid_events 
		ORDER BY timestamp DESC 
		LIMIT $1`

	rows, err := s.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bids []*models.BidEvent
	for rows.Next() {
		bid := &models.BidEvent{}
		err := rows.Scan(
			&bid.ID, &bid.CampaignID, &bid.UserID, &bid.BidPrice,
			&bid.WinPrice, &bid.FloorPrice, &bid.Won, &bid.Converted,
			&bid.SegmentID, &bid.SegmentCategory, &bid.EngagementScore, &bid.ConversionProbability,
			&bid.Country, &bid.Region, &bid.City, &bid.Latitude, &bid.Longitude,
			&bid.DeviceType, &bid.OS, &bid.Browser, &bid.IsMobile, &bid.Keywords,
			&bid.Timestamp, &bid.CreatedAt,
		)
		if err != nil {
			continue // Skip invalid rows
		}
		bids = append(bids, bid)
	}

	return bids, rows.Err()
}

// GetSimilarBids retrieves bids similar to the current request for ML context
func (s *BidStore) GetSimilarBids(ctx context.Context, query *BidQuery) ([]*models.BidEvent, error) {
	sqlQuery := `
		SELECT id, campaign_id, user_id, bid_price, win_price, floor_price, won, converted,
			   segment_id, segment_category, engagement_score, conversion_probability,
			   country, region, city, latitude, longitude, device_type, os, browser,
			   is_mobile, keywords, timestamp, created_at
		FROM bid_events 
		WHERE 1=1`

	args := []interface{}{}
	argCount := 0

	// Add campaign filter
	if query.CampaignID != uuid.Nil {
		argCount++
		sqlQuery += ` AND campaign_id = $` + string(rune(argCount+'0'))
		args = append(args, query.CampaignID)
	}

	// Add segment category filter
	if query.SegmentCategory != "" {
		argCount++
		sqlQuery += ` AND segment_category = $` + string(rune(argCount+'0'))
		args = append(args, query.SegmentCategory)
	}

	// Add country filter
	if query.Country != "" {
		argCount++
		sqlQuery += ` AND country = $` + string(rune(argCount+'0'))
		args = append(args, query.Country)
	}

	// Add device type filter
	if query.DeviceType != "" {
		argCount++
		sqlQuery += ` AND device_type = $` + string(rune(argCount+'0'))
		args = append(args, query.DeviceType)
	}

	// Add time range filter
	if !query.Since.IsZero() {
		argCount++
		sqlQuery += ` AND timestamp > $` + string(rune(argCount+'0'))
		args = append(args, query.Since)
	}

	sqlQuery += ` ORDER BY timestamp DESC`

	// Add limit
	if query.Limit > 0 {
		argCount++
		sqlQuery += ` LIMIT $` + string(rune(argCount+'0'))
		args = append(args, query.Limit)
	}

	rows, err := s.db.QueryContext(ctx, sqlQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bids []*models.BidEvent
	for rows.Next() {
		bid := &models.BidEvent{}
		err := rows.Scan(
			&bid.ID, &bid.CampaignID, &bid.UserID, &bid.BidPrice,
			&bid.WinPrice, &bid.FloorPrice, &bid.Won, &bid.Converted,
			&bid.SegmentID, &bid.SegmentCategory, &bid.EngagementScore, &bid.ConversionProbability,
			&bid.Country, &bid.Region, &bid.City, &bid.Latitude, &bid.Longitude,
			&bid.DeviceType, &bid.OS, &bid.Browser, &bid.IsMobile, &bid.Keywords,
			&bid.Timestamp, &bid.CreatedAt,
		)
		if err != nil {
			continue
		}
		bids = append(bids, bid)
	}

	return bids, rows.Err()
}

// GetUserBids retrieves bid history for a specific user
func (s *BidStore) GetUserBids(ctx context.Context, userID string, since time.Time, limit int) ([]*models.BidEvent, error) {
	query := `
		SELECT id, campaign_id, user_id, bid_price, win_price, floor_price, won, converted,
			   segment_id, segment_category, engagement_score, conversion_probability,
			   country, region, city, latitude, longitude, device_type, os, browser,
			   is_mobile, keywords, timestamp, created_at
		FROM bid_events 
		WHERE user_id = $1 AND timestamp > $2
		ORDER BY timestamp DESC 
		LIMIT $3`

	rows, err := s.db.QueryContext(ctx, query, userID, since, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bids []*models.BidEvent
	for rows.Next() {
		bid := &models.BidEvent{}
		err := rows.Scan(
			&bid.ID, &bid.CampaignID, &bid.UserID, &bid.BidPrice,
			&bid.WinPrice, &bid.FloorPrice, &bid.Won, &bid.Converted,
			&bid.SegmentID, &bid.SegmentCategory, &bid.EngagementScore, &bid.ConversionProbability,
			&bid.Country, &bid.Region, &bid.City, &bid.Latitude, &bid.Longitude,
			&bid.DeviceType, &bid.OS, &bid.Browser, &bid.IsMobile, &bid.Keywords,
			&bid.Timestamp, &bid.CreatedAt,
		)
		if err != nil {
			continue
		}
		bids = append(bids, bid)
	}

	return bids, rows.Err()
}

// GetCampaignBids retrieves all bids for a campaign since a specific time
func (s *BidStore) GetCampaignBids(ctx context.Context, campaignID uuid.UUID, since time.Time) ([]*models.BidEvent, error) {
	query := `
		SELECT id, campaign_id, user_id, bid_price, win_price, floor_price, won, converted,
			   segment_id, segment_category, engagement_score, conversion_probability,
			   country, region, city, latitude, longitude, device_type, os, browser,
			   is_mobile, keywords, timestamp, created_at
		FROM bid_events 
		WHERE campaign_id = $1 AND timestamp > $2
		ORDER BY timestamp DESC`

	rows, err := s.db.QueryContext(ctx, query, campaignID, since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bids []*models.BidEvent
	for rows.Next() {
		bid := &models.BidEvent{}
		err := rows.Scan(
			&bid.ID, &bid.CampaignID, &bid.UserID, &bid.BidPrice,
			&bid.WinPrice, &bid.FloorPrice, &bid.Won, &bid.Converted,
			&bid.SegmentID, &bid.SegmentCategory, &bid.EngagementScore, &bid.ConversionProbability,
			&bid.Country, &bid.Region, &bid.City, &bid.Latitude, &bid.Longitude,
			&bid.DeviceType, &bid.OS, &bid.Browser, &bid.IsMobile, &bid.Keywords,
			&bid.Timestamp, &bid.CreatedAt,
		)
		if err != nil {
			continue
		}
		bids = append(bids, bid)
	}

	return bids, rows.Err()
}

// BidQuery represents parameters for querying similar bids
type BidQuery struct {
	CampaignID      uuid.UUID
	SegmentCategory string
	Country         string
	DeviceType      string
	Since           time.Time
	Limit           int
}
