package models

import (
	"time"

	"github.com/google/uuid"
)

// Campaign represents an advertising campaign
type Campaign struct {
	ID          uuid.UUID `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	UserID      uuid.UUID `db:"user_id" json:"user_id"`
	Status      string    `db:"status" json:"status"`
	Budget      *float64  `db:"budget" json:"budget,omitempty"`
	DailyBudget *float64  `db:"daily_budget" json:"daily_budget,omitempty"`
	TargetCPA   *float64  `db:"target_cpa" json:"target_cpa,omitempty"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// CampaignMetrics represents aggregated campaign performance data
type CampaignMetrics struct {
	ID         uuid.UUID `db:"id" json:"id"`
	CampaignID uuid.UUID `db:"campaign_id" json:"campaign_id"`
	Date       time.Time `db:"date" json:"date"`
	Hour       *int      `db:"hour" json:"hour,omitempty"`

	// Raw metrics
	TotalBids   int     `db:"total_bids" json:"total_bids"`
	WonBids     int     `db:"won_bids" json:"won_bids"`
	Conversions int     `db:"conversions" json:"conversions"`
	TotalSpend  float64 `db:"total_spend" json:"total_spend"`
	Impressions int     `db:"impressions" json:"impressions"`
	Clicks      int     `db:"clicks" json:"clicks"`

	// Calculated metrics
	WinRate           *float64 `db:"win_rate" json:"win_rate,omitempty"`
	ConversionRate    *float64 `db:"conversion_rate" json:"conversion_rate,omitempty"`
	AverageBid        *float64 `db:"average_bid" json:"average_bid,omitempty"`
	CostPerConversion *float64 `db:"cost_per_conversion" json:"cost_per_conversion,omitempty"`
	ReturnOnAdSpend   *float64 `db:"return_on_ad_spend" json:"return_on_ad_spend,omitempty"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// FraudAlert represents a detected fraud pattern
type FraudAlert struct {
	ID              uuid.UUID  `db:"id" json:"id"`
	CampaignID      uuid.UUID  `db:"campaign_id" json:"campaign_id"`
	AlertType       string     `db:"alert_type" json:"alert_type"`
	Severity        int        `db:"severity" json:"severity"`
	Description     string     `db:"description" json:"description"`
	AffectedUserIDs []string   `db:"affected_user_ids" json:"affected_user_ids"`
	DetectedAt      time.Time  `db:"detected_at" json:"detected_at"`
	ResolvedAt      *time.Time `db:"resolved_at" json:"resolved_at,omitempty"`
	Status          string     `db:"status" json:"status"`
}

// ModelMetrics represents ML model performance tracking
type ModelMetrics struct {
	ID                  uuid.UUID `db:"id" json:"id"`
	ModelVersion        string    `db:"model_version" json:"model_version"`
	Date                time.Time `db:"date" json:"date"`
	PredictionAccuracy  *float64  `db:"prediction_accuracy" json:"prediction_accuracy,omitempty"`
	MeanAbsoluteError   *float64  `db:"mean_absolute_error" json:"mean_absolute_error,omitempty"`
	RootMeanSquareError *float64  `db:"root_mean_square_error" json:"root_mean_square_error,omitempty"`
	TotalPredictions    int       `db:"total_predictions" json:"total_predictions"`
	CreatedAt           time.Time `db:"created_at" json:"created_at"`
}

// CampaignStats represents comprehensive campaign statistics
type CampaignStats struct {
	CampaignID        uuid.UUID `json:"campaign_id"`
	TotalBids         int64     `json:"total_bids"`
	WonBids           int64     `json:"won_bids"`
	Conversions       int64     `json:"conversions"`
	TotalSpend        float64   `json:"total_spend"`
	AverageBid        float64   `json:"average_bid"`
	WinRate           float64   `json:"win_rate"`
	ConversionRate    float64   `json:"conversion_rate"`
	CostPerConversion float64   `json:"cost_per_conversion"`
	ReturnOnAdSpend   float64   `json:"return_on_ad_spend"`
}

// User represents a user in the system
type User struct {
	ID           uuid.UUID `db:"id" json:"id"`
	Username     string    `db:"username" json:"username"`
	PasswordHash string    `db:"password_hash" json:"-"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}
