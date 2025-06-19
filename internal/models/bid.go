package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// BidEvent represents a bid request/response event
type BidEvent struct {
	ID         uuid.UUID `db:"id" json:"id"`
	CampaignID uuid.UUID `db:"campaign_id" json:"campaign_id"`
	UserID     string    `db:"user_id" json:"user_id"`

	// Pricing
	BidPrice   float64  `db:"bid_price" json:"bid_price"`
	WinPrice   *float64 `db:"win_price" json:"win_price,omitempty"`
	FloorPrice float64  `db:"floor_price" json:"floor_price"`

	// Outcomes
	Won       bool `db:"won" json:"won"`
	Converted bool `db:"converted" json:"converted"`

	// User Segment
	SegmentID             string   `db:"segment_id" json:"segment_id"`
	SegmentCategory       string   `db:"segment_category" json:"segment_category"`
	EngagementScore       *float64 `db:"engagement_score" json:"engagement_score,omitempty"`
	ConversionProbability *float64 `db:"conversion_probability" json:"conversion_probability,omitempty"`

	// Geography
	Country   string   `db:"country" json:"country"`
	Region    string   `db:"region" json:"region"`
	City      string   `db:"city" json:"city"`
	Latitude  *float64 `db:"latitude" json:"latitude,omitempty"`
	Longitude *float64 `db:"longitude" json:"longitude,omitempty"`

	// Device
	DeviceType string `db:"device_type" json:"device_type"`
	OS         string `db:"os" json:"os"`
	Browser    string `db:"browser" json:"browser"`
	IsMobile   bool   `db:"is_mobile" json:"is_mobile"`

	// Keywords
	Keywords pq.StringArray `db:"keywords" json:"keywords"`

	// Timestamps
	Timestamp time.Time `db:"timestamp" json:"timestamp"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// Prediction represents an ML prediction for a bid
type Prediction struct {
	ID                uuid.UUID `db:"id" json:"id"`
	BidEventID        uuid.UUID `db:"bid_event_id" json:"bid_event_id"`
	PredictedBidPrice float64   `db:"predicted_bid_price" json:"predicted_bid_price"`
	Confidence        float64   `db:"confidence" json:"confidence"`
	Strategy          string    `db:"strategy" json:"strategy"`
	FraudRisk         bool      `db:"fraud_risk" json:"fraud_risk"`
	ModelVersion      string    `db:"model_version" json:"model_version"`
	CreatedAt         time.Time `db:"created_at" json:"created_at"`
}

// UserSegment represents audience segmentation data
type UserSegment struct {
	SegmentID             string  `json:"segment_id"`
	Category              string  `json:"category"`
	EngagementScore       float64 `json:"engagement_score"`
	ConversionProbability float64 `json:"conversion_probability"`
}

// GeoLocation represents geographical data
type GeoLocation struct {
	Country   string  `json:"country"`
	Region    string  `json:"region"`
	City      string  `json:"city"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// DeviceInfo represents device and browser information
type DeviceInfo struct {
	DeviceType string `json:"device_type"`
	OS         string `json:"os"`
	Browser    string `json:"browser"`
	IsMobile   bool   `json:"is_mobile"`
}

// BidRequest represents an incoming bid request
type BidRequest struct {
	CampaignID  uuid.UUID   `json:"campaign_id"`
	UserID      string      `json:"user_id"`
	UserSegment UserSegment `json:"user_segment"`
	GeoLocation GeoLocation `json:"geo_location"`
	DeviceInfo  DeviceInfo  `json:"device_info"`
	FloorPrice  float64     `json:"floor_price"`
	Keywords    []string    `json:"keywords"`
	Timestamp   time.Time   `json:"timestamp"`
}

// BidResponse represents a bid response with prediction
type BidResponse struct {
	BidPrice     float64 `json:"bid_price"`
	Confidence   float64 `json:"confidence"`
	Strategy     string  `json:"strategy"`
	FraudRisk    bool    `json:"fraud_risk"`
	PredictionID string  `json:"prediction_id"`
}
