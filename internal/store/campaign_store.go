package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/baskint/bidding-analysis/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// CampaignStore handles database operations for campaign data
type CampaignStore struct {
	db *sqlx.DB
}

// NewCampaignStore creates a new CampaignStore instance
func NewCampaignStore(db *sqlx.DB) *CampaignStore {
	return &CampaignStore{db: db}
}

// GetCampaignStats retrieves campaign statistics
func (s *CampaignStore) GetCampaignStats(campaignID uuid.UUID, startTime, endTime time.Time) (*models.CampaignStats, error) {
	query := `
		SELECT
			campaign_id,
			COALESCE(SUM(total_bids), 0) as total_bids,
			COALESCE(SUM(won_bids), 0) as won_bids,
			COALESCE(SUM(conversions), 0) as conversions,
			COALESCE(SUM(total_spend), 0) as total_spend,
			CASE
				WHEN SUM(total_bids) > 0 THEN SUM(total_spend) / SUM(total_bids)
				ELSE 0
			END as average_bid,
			CASE
				WHEN SUM(total_bids) > 0 THEN CAST(SUM(won_bids) AS FLOAT) / SUM(total_bids)
				ELSE 0
			END as win_rate,
			CASE
				WHEN SUM(won_bids) > 0 THEN CAST(SUM(conversions) AS FLOAT) / SUM(won_bids)
				ELSE 0
			END as conversion_rate,
			CASE
				WHEN SUM(conversions) > 0 THEN SUM(total_spend) / SUM(conversions)
				ELSE 0
			END as cost_per_conversion,
			CASE
				WHEN SUM(total_spend) > 0 THEN CAST(SUM(conversions) AS FLOAT) / SUM(total_spend)
				ELSE 0
			END as return_on_ad_spend
		FROM campaign_metrics
		WHERE campaign_id = $1 AND date BETWEEN $2 AND $3
		GROUP BY campaign_id`

	stats := &models.CampaignStats{}
	err := s.db.QueryRow(query, campaignID, startTime.Format("2006-01-02"), endTime.Format("2006-01-02")).Scan(
		&stats.CampaignID,
		&stats.TotalBids,
		&stats.WonBids,
		&stats.Conversions,
		&stats.TotalSpend,
		&stats.AverageBid,
		&stats.WinRate,
		&stats.ConversionRate,
		&stats.CostPerConversion,
		&stats.ReturnOnAdSpend,
	)

	if err != nil {
		return nil, err
	}

	return stats, nil
}

// GetFraudAlerts retrieves fraud alerts for a time period
func (s *CampaignStore) GetFraudAlerts(startTime, endTime time.Time, severityThreshold int) ([]*models.FraudAlert, error) {
	query := `
		SELECT id, campaign_id, alert_type, severity, description, detected_at, status
		FROM fraud_alerts
		WHERE detected_at BETWEEN $1 AND $2 AND severity >= $3
		ORDER BY detected_at DESC`

	rows, err := s.db.Query(query, startTime, endTime, severityThreshold)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []*models.FraudAlert
	for rows.Next() {
		alert := &models.FraudAlert{}
		err := rows.Scan(
			&alert.ID,
			&alert.CampaignID,
			&alert.AlertType,
			&alert.Severity,
			&alert.Description,
			&alert.DetectedAt,
			&alert.Status,
		)
		if err != nil {
			return nil, err
		}
		alerts = append(alerts, alert)
	}

	return alerts, rows.Err()
}

// GetModelAccuracy retrieves ML model performance metrics
func (s *CampaignStore) GetModelAccuracy(startTime, endTime time.Time, modelVersion string) (*models.ModelMetrics, error) {
	query := `
		SELECT
			model_version,
			AVG(prediction_accuracy) as avg_accuracy,
			AVG(mean_absolute_error) as avg_mae,
			AVG(root_mean_square_error) as avg_rmse,
			SUM(total_predictions) as total_predictions
		FROM model_metrics
		WHERE date BETWEEN $1 AND $2`

	args := []interface{}{startTime.Format("2006-01-02"), endTime.Format("2006-01-02")}

	if modelVersion != "" {
		query += " AND model_version = $3"
		args = append(args, modelVersion)
	}

	query += " GROUP BY model_version"

	metrics := &models.ModelMetrics{}
	err := s.db.QueryRow(query, args...).Scan(
		&metrics.ModelVersion,
		&metrics.PredictionAccuracy,
		&metrics.MeanAbsoluteError,
		&metrics.RootMeanSquareError,
		&metrics.TotalPredictions,
	)

	if err != nil {
		return nil, err
	}

	return metrics, nil
}

// CreateCampaign creates a new campaign
func (s *CampaignStore) CreateCampaign(campaign *models.Campaign) error {
	query := `
		INSERT INTO campaigns (name, user_id, status, budget, daily_budget, target_cpa)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`

	err := s.db.QueryRow(
		query,
		campaign.Name,
		campaign.UserID,
		campaign.Status,
		campaign.Budget,
		campaign.DailyBudget,
		campaign.TargetCPA,
	).Scan(&campaign.ID, &campaign.CreatedAt, &campaign.UpdatedAt)

	return err
}

// GetCampaign retrieves a campaign by ID
func (s *CampaignStore) GetCampaign(id uuid.UUID) (*models.Campaign, error) {
	query := `
		SELECT id, name, user_id, status, budget, daily_budget, target_cpa, created_at, updated_at
		FROM campaigns
		WHERE id = $1`

	campaign := &models.Campaign{}
	err := s.db.QueryRow(query, id).Scan(
		&campaign.ID,
		&campaign.Name,
		&campaign.UserID,
		&campaign.Status,
		&campaign.Budget,
		&campaign.DailyBudget,
		&campaign.TargetCPA,
		&campaign.CreatedAt,
		&campaign.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return campaign, nil
}

// GetDB returns the database connection (for internal use)
func (s *CampaignStore) GetDB() *sqlx.DB {
	return s.db
}

// ListCampaigns retrieves campaigns for a user
func (s *CampaignStore) ListCampaigns(userID uuid.UUID) ([]*models.Campaign, error) {
	query := `
		SELECT id, name, user_id, status, budget, daily_budget, target_cpa, created_at, updated_at
		FROM campaigns
		WHERE user_id = $1
		ORDER BY created_at DESC`

	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var campaigns []*models.Campaign
	for rows.Next() {
		campaign := &models.Campaign{}
		err := rows.Scan(
			&campaign.ID,
			&campaign.Name,
			&campaign.UserID,
			&campaign.Status,
			&campaign.Budget,
			&campaign.DailyBudget,
			&campaign.TargetCPA,
			&campaign.CreatedAt,
			&campaign.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		campaigns = append(campaigns, campaign)
	}

	return campaigns, rows.Err()
}

// GetUserCampaigns retrieves campaigns for a user with context support
func (s *CampaignStore) GetUserCampaigns(ctx context.Context, userID string) ([]*models.Campaign, error) {
	query := `
		SELECT id, name, user_id, status, budget, daily_budget, target_cpa, created_at, updated_at
		FROM campaigns
		WHERE user_id = $1
		ORDER BY created_at DESC`

	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var campaigns []*models.Campaign
	for rows.Next() {
		campaign := &models.Campaign{}
		err := rows.Scan(
			&campaign.ID,
			&campaign.Name,
			&campaign.UserID,
			&campaign.Status,
			&campaign.Budget,
			&campaign.DailyBudget,
			&campaign.TargetCPA,
			&campaign.CreatedAt,
			&campaign.UpdatedAt,
		)
		if err != nil {
			continue // Skip invalid rows
		}
		campaigns = append(campaigns, campaign)
	}

	return campaigns, rows.Err()
}

// UpdateCampaign updates an existing campaign
func (s *CampaignStore) UpdateCampaign(campaign *models.Campaign) error {
	query := `
		UPDATE campaigns
		SET name = $1, status = $2, budget = $3, daily_budget = $4, target_cpa = $5, updated_at = NOW()
		WHERE id = $6 AND user_id = $7
		RETURNING updated_at`

	err := s.db.QueryRow(
		query,
		campaign.Name,
		campaign.Status,
		campaign.Budget,
		campaign.DailyBudget,
		campaign.TargetCPA,
		campaign.ID,
		campaign.UserID,
	).Scan(&campaign.UpdatedAt)

	if err == sql.ErrNoRows {
		return fmt.Errorf("campaign not found or unauthorized")
	}

	return err
}

// DeleteCampaign soft deletes a campaign by setting status to 'archived'
func (s *CampaignStore) DeleteCampaign(id uuid.UUID, userID uuid.UUID) error {
	query := `
		UPDATE campaigns
		SET status = 'archived', updated_at = NOW()
		WHERE id = $1 AND user_id = $2`

	result, err := s.db.Exec(query, id, userID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("campaign not found or unauthorized")
	}

	return nil
}

// PauseCampaign pauses an active campaign
func (s *CampaignStore) PauseCampaign(id uuid.UUID, userID uuid.UUID) error {
	query := `
		UPDATE campaigns
		SET status = 'paused', updated_at = NOW()
		WHERE id = $1 AND user_id = $2 AND status = 'active'`

	result, err := s.db.Exec(query, id, userID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("campaign not found, unauthorized, or not active")
	}

	return nil
}

// ActivateCampaign activates a paused campaign
func (s *CampaignStore) ActivateCampaign(id uuid.UUID, userID uuid.UUID) error {
	query := `
		UPDATE campaigns
		SET status = 'active', updated_at = NOW()
		WHERE id = $1 AND user_id = $2 AND status = 'paused'`

	result, err := s.db.Exec(query, id, userID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("campaign not found, unauthorized, or not paused")
	}

	return nil
}

// GetCampaignWithMetrics retrieves a campaign with aggregated performance metrics
func (s *CampaignStore) GetCampaignWithMetrics(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*models.CampaignDetail, error) {
	// Get base campaign
	campaign, err := s.GetCampaign(id)
	if err != nil {
		return nil, err
	}

	// Verify ownership
	if campaign.UserID != userID {
		return nil, fmt.Errorf("unauthorized")
	}

	// Get campaign stats
	stats, err := s.GetCampaignStats(id, time.Now().AddDate(0, 0, -30), time.Now())
	if err != nil {
		// If no stats yet, use zero values
		stats = &models.CampaignStats{
			CampaignID: id,
		}
	}

	detail := &models.CampaignDetail{
		Campaign:      *campaign,
		CampaignStats: *stats,
	}

	// Get daily metrics for the last 30 days
	dailyMetrics, err := s.GetCampaignDailyMetrics(id, time.Now().AddDate(0, 0, -30), time.Now())
	if err == nil {
		detail.DailyMetrics = dailyMetrics
	}

	// Get recent bids (last 20)
	recentBids, err := s.getRecentCampaignBids(id, 20)
	if err == nil {
		detail.RecentBids = recentBids
	}

	// Get device breakdown
	deviceBreakdown, err := s.getDeviceBreakdown(id)
	if err == nil {
		detail.DeviceBreakdown = deviceBreakdown
	}

	// Get geo breakdown
	geoBreakdown, err := s.getGeoBreakdown(id)
	if err == nil {
		detail.GeoBreakdown = geoBreakdown
	}

	// Get top keywords
	topKeywords, err := s.getTopKeywords(id)
	if err == nil {
		detail.TopKeywords = topKeywords
	}

	return detail, nil
}

// ListCampaignsWithMetrics lists campaigns with summary metrics
func (s *CampaignStore) ListCampaignsWithMetrics(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.CampaignSummary, error) {
	query := `
		SELECT
			c.id, c.name, c.user_id, c.status, c.budget, c.daily_budget, c.target_cpa, c.created_at, c.updated_at,
			COALESCE(SUM(cm.total_bids), 0) as total_bids,
			COALESCE(SUM(cm.won_bids), 0) as won_bids,
			COALESCE(SUM(cm.conversions), 0) as conversions,
			COALESCE(SUM(cm.total_spend), 0) as total_spend,
			CASE
				WHEN SUM(cm.total_bids) > 0 THEN CAST(SUM(cm.won_bids) AS FLOAT) / SUM(cm.total_bids)
				ELSE 0
			END as win_rate,
			CASE
				WHEN SUM(cm.won_bids) > 0 THEN CAST(SUM(cm.conversions) AS FLOAT) / SUM(cm.won_bids)
				ELSE 0
			END as conversion_rate,
			CASE
				WHEN SUM(cm.total_bids) > 0 THEN SUM(cm.total_spend) / SUM(cm.total_bids)
				ELSE 0
			END as average_bid,
			CASE
				WHEN SUM(cm.conversions) > 0 THEN SUM(cm.total_spend) / SUM(cm.conversions)
				ELSE 0
			END as cost_per_conversion,
			MAX(be.timestamp) as last_activity_at
		FROM campaigns c
		LEFT JOIN campaign_metrics cm ON c.id = cm.campaign_id AND cm.date >= CURRENT_DATE - INTERVAL '30 days'
		LEFT JOIN bid_events be ON c.id = be.campaign_id
		WHERE c.user_id = $1
		GROUP BY c.id, c.name, c.user_id, c.status, c.budget, c.daily_budget, c.target_cpa, c.created_at, c.updated_at
		ORDER BY c.created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := s.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var campaigns []*models.CampaignSummary
	for rows.Next() {
		summary := &models.CampaignSummary{}
		err := rows.Scan(
			&summary.ID,
			&summary.Name,
			&summary.UserID,
			&summary.Status,
			&summary.Budget,
			&summary.DailyBudget,
			&summary.TargetCPA,
			&summary.CreatedAt,
			&summary.UpdatedAt,
			&summary.TotalBids,
			&summary.WonBids,
			&summary.Conversions,
			&summary.TotalSpend,
			&summary.WinRate,
			&summary.ConversionRate,
			&summary.AverageBid,
			&summary.CostPerConversion,
			&summary.LastActivityAt,
		)
		if err != nil {
			continue
		}
		campaigns = append(campaigns, summary)
	}

	return campaigns, rows.Err()
}

// GetCampaignDailyMetrics retrieves daily performance metrics for a campaign
func (s *CampaignStore) GetCampaignDailyMetrics(campaignID uuid.UUID, startDate, endDate time.Time) ([]*models.DailyMetric, error) {
	query := `
		SELECT
			date,
			COALESCE(SUM(total_bids), 0) as total_bids,
			COALESCE(SUM(won_bids), 0) as won_bids,
			COALESCE(SUM(conversions), 0) as conversions,
			COALESCE(SUM(total_spend), 0) as total_spend,
			CASE
				WHEN SUM(total_bids) > 0 THEN CAST(SUM(won_bids) AS FLOAT) / SUM(total_bids)
				ELSE 0
			END as win_rate,
			CASE
				WHEN SUM(won_bids) > 0 THEN CAST(SUM(conversions) AS FLOAT) / SUM(won_bids)
				ELSE 0
			END as conversion_rate
		FROM campaign_metrics
		WHERE campaign_id = $1 AND date BETWEEN $2 AND $3
		GROUP BY date
		ORDER BY date ASC`

	rows, err := s.db.Query(query, campaignID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []*models.DailyMetric
	for rows.Next() {
		metric := &models.DailyMetric{}
		err := rows.Scan(
			&metric.Date,
			&metric.TotalBids,
			&metric.WonBids,
			&metric.Conversions,
			&metric.TotalSpend,
			&metric.WinRate,
			&metric.ConversionRate,
		)
		if err != nil {
			continue
		}
		metrics = append(metrics, metric)
	}

	return metrics, rows.Err()
}

// Helper methods for detailed breakdown

func (s *CampaignStore) getRecentCampaignBids(campaignID uuid.UUID, limit int) ([]*models.BidEvent, error) {
	query := `
		SELECT id, campaign_id, user_id, bid_price, win_price, floor_price, won, converted,
			   segment_id, segment_category, country, region, city,
			   device_type, os, browser, is_mobile, timestamp
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
			&bid.ID, &bid.CampaignID, &bid.UserID, &bid.BidPrice, &bid.WinPrice, &bid.FloorPrice,
			&bid.Won, &bid.Converted, &bid.SegmentID, &bid.SegmentCategory,
			&bid.Country, &bid.Region, &bid.City, &bid.DeviceType, &bid.OS, &bid.Browser,
			&bid.IsMobile, &bid.Timestamp,
		)
		if err != nil {
			continue
		}
		bids = append(bids, bid)
	}

	return bids, rows.Err()
}

func (s *CampaignStore) getDeviceBreakdown(campaignID uuid.UUID) ([]models.DeviceStat, error) {
	query := `
		SELECT
			device_type,
			COUNT(*) as bids,
			SUM(CASE WHEN won THEN 1 ELSE 0 END) as won_bids,
			SUM(CASE WHEN converted THEN 1 ELSE 0 END) as conversions,
			CASE
				WHEN COUNT(*) > 0 THEN CAST(SUM(CASE WHEN won THEN 1 ELSE 0 END) AS FLOAT) / COUNT(*)
				ELSE 0
			END as win_rate
		FROM bid_events
		WHERE campaign_id = $1 AND device_type IS NOT NULL
		GROUP BY device_type
		ORDER BY bids DESC`

	rows, err := s.db.Query(query, campaignID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []models.DeviceStat
	for rows.Next() {
		stat := models.DeviceStat{}
		err := rows.Scan(&stat.DeviceType, &stat.Bids, &stat.WonBids, &stat.Conversions, &stat.WinRate)
		if err != nil {
			continue
		}
		stats = append(stats, stat)
	}

	return stats, rows.Err()
}

func (s *CampaignStore) getGeoBreakdown(campaignID uuid.UUID) ([]models.GeoStat, error) {
	query := `
		SELECT
			country,
			COUNT(*) as bids,
			SUM(CASE WHEN won THEN 1 ELSE 0 END) as won_bids,
			SUM(CASE WHEN converted THEN 1 ELSE 0 END) as conversions,
			CASE
				WHEN COUNT(*) > 0 THEN CAST(SUM(CASE WHEN won THEN 1 ELSE 0 END) AS FLOAT) / COUNT(*)
				ELSE 0
			END as win_rate
		FROM bid_events
		WHERE campaign_id = $1 AND country IS NOT NULL
		GROUP BY country
		ORDER BY bids DESC
		LIMIT 10`

	rows, err := s.db.Query(query, campaignID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []models.GeoStat
	for rows.Next() {
		stat := models.GeoStat{}
		err := rows.Scan(&stat.Country, &stat.Bids, &stat.WonBids, &stat.Conversions, &stat.WinRate)
		if err != nil {
			continue
		}
		stats = append(stats, stat)
	}

	return stats, rows.Err()
}

func (s *CampaignStore) getTopKeywords(campaignID uuid.UUID) ([]models.KeywordStat, error) {
	query := `
		SELECT
			unnest(keywords) as keyword,
			COUNT(*) as bids,
			SUM(CASE WHEN won THEN 1 ELSE 0 END) as won_bids,
			SUM(CASE WHEN converted THEN 1 ELSE 0 END) as conversions,
			COALESCE(SUM(CASE WHEN won THEN win_price ELSE 0 END), 0) as spend,
			CASE
				WHEN COUNT(*) > 0 THEN CAST(SUM(CASE WHEN won THEN 1 ELSE 0 END) AS FLOAT) / COUNT(*)
				ELSE 0
			END as win_rate
		FROM bid_events
		WHERE campaign_id = $1 AND keywords IS NOT NULL AND array_length(keywords, 1) > 0
		GROUP BY keyword
		ORDER BY bids DESC
		LIMIT 10`

	rows, err := s.db.Query(query, campaignID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []models.KeywordStat
	for rows.Next() {
		stat := models.KeywordStat{}
		err := rows.Scan(&stat.Keyword, &stat.Bids, &stat.WonBids, &stat.Conversions, &stat.Spend, &stat.WinRate)
		if err != nil {
			continue
		}
		stats = append(stats, stat)
	}

	return stats, rows.Err()
}
