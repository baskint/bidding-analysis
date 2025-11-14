// internal/models/fraud.go
package models

import (
	"time"

	"github.com/google/uuid"
)

// FraudOverview contains high-level fraud metrics
type FraudOverview struct {
	TotalAlerts          int            `json:"total_alerts"`
	ActiveAlerts         int            `json:"active_alerts"`
	BlockedBids          int64          `json:"blocked_bids"`
	AmountSaved          float64        `json:"amount_saved"`
	ThreatLevel          string         `json:"threat_level"` // low, medium, high, critical
	AlertsByType         map[string]int `json:"alerts_by_type"`
	TopAffectedCampaigns []CampaignRisk `json:"top_affected_campaigns"`
}

// CampaignRisk represents fraud risk assessment for a campaign
type CampaignRisk struct {
	CampaignID    string  `json:"campaign_id"`
	CampaignName  string  `json:"campaign_name"`
	RiskScore     float64 `json:"risk_score"`
	FraudAttempts int     `json:"fraud_attempts"`
	ThreatLevel   string  `json:"threat_level"` // low, medium, high, critical
}

// FraudTrend represents fraud metrics over time
type FraudTrend struct {
	Date          string  `json:"date"`
	FraudAttempts int     `json:"fraud_attempts"`
	BlockedBids   int     `json:"blocked_bids"`
	AmountSaved   float64 `json:"amount_saved"`
	AlertType     string  `json:"alert_type"`
}

// DeviceFraudAnalysis represents device-specific fraud metrics
type DeviceFraudAnalysis struct {
	DeviceType string  `json:"device_type"`
	Browser    string  `json:"browser"`
	OS         string  `json:"os"`
	TotalBids  int     `json:"total_bids"`
	FraudBids  int     `json:"fraud_bids"`
	FraudRate  float64 `json:"fraud_rate"`
}

// GeoFraudAnalysis represents geographic fraud patterns
type GeoFraudAnalysis struct {
	Country   string  `json:"country"`
	Region    string  `json:"region"`
	City      string  `json:"city"`
	TotalBids int     `json:"total_bids"`
	FraudBids int     `json:"fraud_bids"`
	FraudRate float64 `json:"fraud_rate"`
}

// FraudAlertFilter contains filtering options for fraud alerts
type FraudAlertFilter struct {
	Status      string
	MinSeverity int
	AlertType   string
	StartDate   time.Time
	EndDate     time.Time
	Limit       int
}

// FraudAlertCreate represents the input for creating a new fraud alert
type FraudAlertCreate struct {
	CampaignID      uuid.UUID `json:"campaign_id" binding:"required"`
	AlertType       string    `json:"alert_type" binding:"required"`
	Severity        int       `json:"severity" binding:"required"`
	Description     string    `json:"description" binding:"required"`
	AffectedUserIDs []string  `json:"affected_user_ids,omitempty"`
}

// FraudAlertUpdate represents the input for updating a fraud alert
type FraudAlertUpdate struct {
	Status      *string    `json:"status,omitempty"`
	ResolvedAt  *time.Time `json:"resolved_at,omitempty"`
	Description *string    `json:"description,omitempty"`
}

// FraudRule represents a custom fraud detection rule
type FraudRule struct {
	ID         uuid.UUID              `db:"id" json:"id"`
	UserID     uuid.UUID              `db:"user_id" json:"user_id"`
	Name       string                 `db:"name" json:"name"`
	RuleType   string                 `db:"rule_type" json:"rule_type"`   // click_velocity, geo_anomaly, etc.
	Conditions map[string]interface{} `db:"conditions" json:"conditions"` // JSONB field
	Threshold  float64                `db:"threshold" json:"threshold"`
	Severity   int                    `db:"severity" json:"severity"`
	Enabled    bool                   `db:"enabled" json:"enabled"`
	AutoBlock  bool                   `db:"auto_block" json:"auto_block"`
	CreatedAt  time.Time              `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time              `db:"updated_at" json:"updated_at"`
}

// FraudRuleCreate represents the input for creating a fraud rule
type FraudRuleCreate struct {
	Name       string                 `json:"name" binding:"required"`
	RuleType   string                 `json:"rule_type" binding:"required"`
	Conditions map[string]interface{} `json:"conditions" binding:"required"`
	Threshold  float64                `json:"threshold" binding:"required"`
	Severity   int                    `json:"severity" binding:"required"`
	AutoBlock  bool                   `json:"auto_block"`
}

// FraudRuleUpdate represents the input for updating a fraud rule
type FraudRuleUpdate struct {
	Name       *string                 `json:"name,omitempty"`
	Conditions *map[string]interface{} `json:"conditions,omitempty"`
	Threshold  *float64                `json:"threshold,omitempty"`
	Severity   *int                    `json:"severity,omitempty"`
	Enabled    *bool                   `json:"enabled,omitempty"`
	AutoBlock  *bool                   `json:"auto_block,omitempty"`
}

// BlockedEntity represents a blocked IP, device, or user
type BlockedEntity struct {
	ID            uuid.UUID  `db:"id" json:"id"`
	UserID        uuid.UUID  `db:"user_id" json:"user_id"`
	EntityType    string     `db:"entity_type" json:"entity_type"` // ip, device, user
	EntityValue   string     `db:"entity_value" json:"entity_value"`
	Reason        string     `db:"reason" json:"reason"`
	BlockedByRule *uuid.UUID `db:"blocked_by_rule_id" json:"blocked_by_rule_id,omitempty"`
	BlockedAt     time.Time  `db:"blocked_at" json:"blocked_at"`
	ExpiresAt     *time.Time `db:"expires_at" json:"expires_at,omitempty"`
	Permanent     bool       `db:"permanent" json:"permanent"`
	CreatedBy     *uuid.UUID `db:"created_by" json:"created_by,omitempty"`
}

// BlockedEntityCreate represents the input for creating a blocked entity
type BlockedEntityCreate struct {
	EntityType  string     `json:"entity_type" binding:"required"` // ip, device, user
	EntityValue string     `json:"entity_value" binding:"required"`
	Reason      string     `json:"reason" binding:"required"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	Permanent   bool       `json:"permanent"`
}

// FraudMetrics represents aggregated fraud metrics for reporting
type FraudMetrics struct {
	ID             uuid.UUID  `db:"id" json:"id"`
	UserID         uuid.UUID  `db:"user_id" json:"user_id"`
	Date           time.Time  `db:"date" json:"date"`
	Hour           *int       `db:"hour" json:"hour,omitempty"`
	FraudAttempts  int        `db:"fraud_attempts" json:"fraud_attempts"`
	BlockedBids    int        `db:"blocked_bids" json:"blocked_bids"`
	AmountSaved    float64    `db:"amount_saved" json:"amount_saved"`
	FalsePositives int        `db:"false_positives" json:"false_positives"`
	AlertType      string     `db:"alert_type" json:"alert_type"`
	CampaignID     *uuid.UUID `db:"campaign_id" json:"campaign_id,omitempty"`
	CreatedAt      time.Time  `db:"created_at" json:"created_at"`
}

// DetectionType represents different fraud detection patterns
type DetectionType string

const (
	DetectionClickVelocity     DetectionType = "click_velocity"
	DetectionIPAnomaly         DetectionType = "ip_anomaly"
	DetectionGeoAnomaly        DetectionType = "geo_anomaly"
	DetectionDeviceAnomaly     DetectionType = "device_anomaly"
	DetectionBidPriceAnomaly   DetectionType = "bid_price_anomaly"
	DetectionConversionAnomaly DetectionType = "conversion_anomaly"
	DetectionBotDetection      DetectionType = "bot_detection"
)

// FraudDetectionResult represents the result of fraud detection on a bid
type FraudDetectionResult struct {
	IsFraud       bool
	DetectionType DetectionType
	Reason        string
	Confidence    float64
}

// FraudAlertStatus represents possible alert statuses
type FraudAlertStatus string

const (
	AlertStatusActive        FraudAlertStatus = "active"
	AlertStatusInvestigating FraudAlertStatus = "investigating"
	AlertStatusResolved      FraudAlertStatus = "resolved"
	AlertStatusFalsePositive FraudAlertStatus = "false_positive"
)

// ThreatLevel represents the overall threat assessment
type ThreatLevel string

const (
	ThreatLevelLow      ThreatLevel = "low"
	ThreatLevelMedium   ThreatLevel = "medium"
	ThreatLevelHigh     ThreatLevel = "high"
	ThreatLevelCritical ThreatLevel = "critical"
)
