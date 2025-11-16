package fraud

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/baskint/bidding-analysis/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// DetectionType represents different fraud detection patterns
type DetectionType string

const (
	ClickVelocity     DetectionType = "click_velocity"
	IPAnomaly         DetectionType = "ip_anomaly"
	GeoAnomaly        DetectionType = "geo_anomaly"
	DeviceAnomaly     DetectionType = "device_anomaly"
	BidPriceAnomaly   DetectionType = "bid_price_anomaly"
	ConversionAnomaly DetectionType = "conversion_anomaly"
	BotDetection      DetectionType = "bot_detection"
)

// FraudDetector handles fraud detection logic
type FraudDetector struct {
	db *sqlx.DB
}

// NewFraudDetector creates a new fraud detector instance
func NewFraudDetector(db *sqlx.DB) *FraudDetector {
	return &FraudDetector{db: db}
}

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

// CampaignRisk represents fraud risk for a campaign
type CampaignRisk struct {
	CampaignID    string  `json:"campaign_id"`
	CampaignName  string  `json:"campaign_name"`
	RiskScore     float64 `json:"risk_score"`
	FraudAttempts int     `json:"fraud_attempts"`
	ThreatLevel   string  `json:"threat_level"`
}

// FraudTrend represents fraud metrics over time
type FraudTrend struct {
	Date          string  `json:"date"`
	FraudAttempts int     `json:"fraud_attempts"`
	BlockedBids   int     `json:"blocked_bids"`
	AmountSaved   float64 `json:"amount_saved"`
	AlertType     string  `json:"alert_type"`
}

// BlockedEntity represents a blocked IP, device, or user
type BlockedEntity struct {
	ID          string     `json:"id"`
	EntityType  string     `json:"entity_type"` // ip, device, user
	EntityValue string     `json:"entity_value"`
	Reason      string     `json:"reason"`
	BlockedAt   time.Time  `json:"blocked_at"`
	ExpiresAt   *time.Time `json:"expires_at"`
	Permanent   bool       `json:"permanent"`
}

// GetFraudOverview returns high-level fraud metrics
func (fd *FraudDetector) GetFraudOverview(ctx context.Context, userID uuid.UUID, days int) (*FraudOverview, error) {
	startDate := time.Now().AddDate(0, 0, -days)

	// Get total alerts count
	var totalAlerts, activeAlerts int
	err := fd.db.QueryRowContext(ctx, `
		SELECT 
			COUNT(*) as total,
			-- FIX: Use COALESCE to ensure the SUM is always an INT (0 if no rows match)
			COALESCE(SUM(CASE WHEN fa.status = 'active' THEN 1 ELSE 0 END), 0) as active
		FROM fraud_alerts fa
		JOIN campaigns c ON fa.campaign_id = c.id
		WHERE c.user_id = $1 AND fa.detected_at >= $2
	`, userID, startDate).Scan(&totalAlerts, &activeAlerts)

	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get alert counts: %w", err)
	}

	// Get blocked bids and amount saved (estimated from predictions)
	var blockedBids int64
	var amountSaved float64
	err = fd.db.QueryRowContext(ctx, `
		SELECT 
			COUNT(*) as blocked_bids,
			COALESCE(SUM(be.bid_price), 0) as amount_saved
		FROM predictions p
		JOIN bid_events be ON p.bid_event_id = be.id
		WHERE be.user_id = $1 
			AND p.fraud_risk = true 
			AND be.timestamp >= $2
	`, userID, startDate).Scan(&blockedBids, &amountSaved)

	if err != nil && err != sql.ErrNoRows {
		log.Printf("Warning: Could not get blocked bids: %v", err)
	}

	// Get alerts by type
	alertsByType := make(map[string]int)
	rows, err := fd.db.QueryContext(ctx, `
		SELECT alert_type, COUNT(*) as count
		FROM fraud_alerts fa
		JOIN campaigns c ON fa.campaign_id = c.id
		WHERE c.user_id = $1 AND fa.detected_at >= $2
		GROUP BY alert_type
	`, userID, startDate)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var alertType string
			var count int
			if err := rows.Scan(&alertType, &count); err == nil {
				alertsByType[alertType] = count
			}
		}
	}

	// Calculate threat level
	threatLevel := "low"
	if activeAlerts > 10 {
		threatLevel = "critical"
	} else if activeAlerts > 5 {
		threatLevel = "high"
	} else if activeAlerts > 2 {
		threatLevel = "medium"
	}

	// Get top affected campaigns
	topCampaigns := []CampaignRisk{}
	campaignRows, err := fd.db.QueryContext(ctx, `
		SELECT 
			c.id,
			c.name,
			COUNT(fa.id) as fraud_attempts,
			CASE 
				WHEN COUNT(fa.id) > 10 THEN 9.0
				WHEN COUNT(fa.id) > 5 THEN 7.0
				WHEN COUNT(fa.id) > 2 THEN 5.0
				ELSE 3.0
			END as risk_score
		FROM campaigns c
		LEFT JOIN fraud_alerts fa ON c.id = fa.campaign_id AND fa.detected_at >= $2
		WHERE c.user_id = $1
		GROUP BY c.id, c.name
		HAVING COUNT(fa.id) > 0
		ORDER BY fraud_attempts DESC
		LIMIT 5
	`, userID, startDate)

	if err == nil {
		defer campaignRows.Close()
		for campaignRows.Next() {
			var risk CampaignRisk
			var id uuid.UUID
			err := campaignRows.Scan(&id, &risk.CampaignName, &risk.FraudAttempts, &risk.RiskScore)
			if err == nil {
				risk.CampaignID = id.String()
				// Determine threat level based on risk score
				if risk.RiskScore >= 8 {
					risk.ThreatLevel = "critical"
				} else if risk.RiskScore >= 6 {
					risk.ThreatLevel = "high"
				} else if risk.RiskScore >= 4 {
					risk.ThreatLevel = "medium"
				} else {
					risk.ThreatLevel = "low"
				}
				topCampaigns = append(topCampaigns, risk)
			}
		}
	}

	return &FraudOverview{
		TotalAlerts:          totalAlerts,
		ActiveAlerts:         activeAlerts,
		BlockedBids:          blockedBids,
		AmountSaved:          amountSaved,
		ThreatLevel:          threatLevel,
		AlertsByType:         alertsByType,
		TopAffectedCampaigns: topCampaigns,
	}, nil
}

// GetFraudAlerts retrieves fraud alerts with filtering
func (fd *FraudDetector) GetFraudAlerts(ctx context.Context, userID uuid.UUID, filter FraudAlertFilter) ([]*models.FraudAlert, error) {
	query := `
		SELECT 
			fa.id, fa.campaign_id, fa.alert_type, fa.severity, 
			fa.description, fa.affected_user_ids, fa.detected_at, 
			fa.resolved_at, fa.status
		FROM fraud_alerts fa
		JOIN campaigns c ON fa.campaign_id = c.id
		WHERE c.user_id = $1
	`
	args := []interface{}{userID}
	argCount := 1

	// Add filters
	if filter.Status != "" {
		argCount++
		query += fmt.Sprintf(" AND fa.status = $%d", argCount)
		args = append(args, filter.Status)
	}

	if filter.MinSeverity > 0 {
		argCount++
		query += fmt.Sprintf(" AND fa.severity >= $%d", argCount)
		args = append(args, filter.MinSeverity)
	}

	if filter.AlertType != "" {
		argCount++
		query += fmt.Sprintf(" AND fa.alert_type = $%d", argCount)
		args = append(args, filter.AlertType)
	}

	if !filter.StartDate.IsZero() {
		argCount++
		query += fmt.Sprintf(" AND fa.detected_at >= $%d", argCount)
		args = append(args, filter.StartDate)
	}

	if !filter.EndDate.IsZero() {
		argCount++
		query += fmt.Sprintf(" AND fa.detected_at <= $%d", argCount)
		args = append(args, filter.EndDate)
	}

	query += " ORDER BY fa.detected_at DESC"

	if filter.Limit > 0 {
		argCount++
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, filter.Limit)
	} else {
		query += " LIMIT 100"
	}

	rows, err := fd.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query fraud alerts: %w", err)
	}
	defer rows.Close()

	var alerts []*models.FraudAlert
	for rows.Next() {
		alert := &models.FraudAlert{}
		// Note: affected_user_ids is TEXT[] in postgres, we need special handling
		var affectedUserIDs sql.NullString
		err := rows.Scan(
			&alert.ID,
			&alert.CampaignID,
			&alert.AlertType,
			&alert.Severity,
			&alert.Description,
			&affectedUserIDs,
			&alert.DetectedAt,
			&alert.ResolvedAt,
			&alert.Status,
		)
		if err != nil {
			log.Printf("Error scanning fraud alert: %v", err)
			continue
		}

		// Parse affected user IDs if present
		// In production, you'd properly parse the PostgreSQL array
		if affectedUserIDs.Valid {
			alert.AffectedUserIDs = []string{} // Simplified for now
		}

		alerts = append(alerts, alert)
	}

	return alerts, nil
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

// UpdateAlertStatus updates the status of a fraud alert
func (fd *FraudDetector) UpdateAlertStatus(ctx context.Context, alertID uuid.UUID, userID uuid.UUID, status string, notes string) error {
	// Verify ownership
	var campaignUserID uuid.UUID
	err := fd.db.QueryRowContext(ctx, `
		SELECT c.user_id
		FROM fraud_alerts fa
		JOIN campaigns c ON fa.campaign_id = c.id
		WHERE fa.id = $1
	`, alertID).Scan(&campaignUserID)

	if err != nil {
		return fmt.Errorf("alert not found: %w", err)
	}

	if campaignUserID != userID {
		return fmt.Errorf("unauthorized")
	}

	// Update status
	resolvedAt := sql.NullTime{}
	if status == "resolved" || status == "false_positive" {
		resolvedAt = sql.NullTime{Time: time.Now(), Valid: true}
	}

	_, err = fd.db.ExecContext(ctx, `
		UPDATE fraud_alerts 
		SET status = $1, resolved_at = $2, description = COALESCE($3, description)
		WHERE id = $4
	`, status, resolvedAt, notes, alertID)

	return err
}

// GetFraudTrends returns fraud metrics over time
func (fd *FraudDetector) GetFraudTrends(ctx context.Context, userID uuid.UUID, days int) ([]FraudTrend, error) {
	startDate := time.Now().AddDate(0, 0, -days)

	query := `
		SELECT 
			DATE(fa.detected_at) as date,
			fa.alert_type,
			COUNT(*) as fraud_attempts,
			0 as blocked_bids,
			0.0 as amount_saved
		FROM fraud_alerts fa
		JOIN campaigns c ON fa.campaign_id = c.id
		WHERE c.user_id = $1 AND fa.detected_at >= $2
		GROUP BY DATE(fa.detected_at), fa.alert_type
		ORDER BY date DESC, alert_type
	`

	rows, err := fd.db.QueryContext(ctx, query, userID, startDate)
	if err != nil {
		return nil, fmt.Errorf("failed to query fraud trends: %w", err)
	}
	defer rows.Close()

	var trends []FraudTrend
	for rows.Next() {
		var trend FraudTrend
		var date time.Time
		err := rows.Scan(
			&date,
			&trend.AlertType,
			&trend.FraudAttempts,
			&trend.BlockedBids,
			&trend.AmountSaved,
		)
		if err != nil {
			log.Printf("Error scanning fraud trend: %v", err)
			continue
		}
		trend.Date = date.Format("2006-01-02")
		trends = append(trends, trend)
	}

	return trends, nil
}

// DetectFraud runs fraud detection on a bid event
// This would be called in real-time during bid processing
func (fd *FraudDetector) DetectFraud(ctx context.Context, bidEvent *models.BidEvent) (bool, DetectionType, string) {
	// Simple click velocity check
	isFraud, detectionType := fd.checkClickVelocity(ctx, bidEvent)
	if isFraud {
		reason := "Abnormal click velocity detected"
		return true, detectionType, reason
	}

	// Add more detection logic here
	// - IP reputation check
	// - Device fingerprint analysis
	// - Geographic anomalies
	// - Bid price anomalies
	// etc.

	return false, "", ""
}

// checkClickVelocity checks for abnormal click patterns
func (fd *FraudDetector) checkClickVelocity(ctx context.Context, bidEvent *models.BidEvent) (bool, DetectionType) {
	// Check number of bids from same segment in last minute
	var recentBids int
	oneMinuteAgo := time.Now().Add(-1 * time.Minute)

	err := fd.db.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM bid_events
		WHERE segment_id = $1 
			AND campaign_id = $2 
			AND timestamp > $3
	`, bidEvent.SegmentID, bidEvent.CampaignID, oneMinuteAgo).Scan(&recentBids)

	if err != nil {
		log.Printf("Error checking click velocity: %v", err)
		return false, ""
	}

	// If more than 10 bids in one minute from same segment, flag as suspicious
	if recentBids > 10 {
		return true, ClickVelocity
	}

	return false, ""
}

// CreateFraudAlert creates a new fraud alert
func (fd *FraudDetector) CreateFraudAlert(ctx context.Context, alert *models.FraudAlert) error {
	_, err := fd.db.ExecContext(ctx, `
		INSERT INTO fraud_alerts (
			id, campaign_id, alert_type, severity, description, 
			affected_user_ids, detected_at, status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, alert.ID, alert.CampaignID, alert.AlertType, alert.Severity,
		alert.Description, alert.AffectedUserIDs, alert.DetectedAt, alert.Status)

	return err
}
