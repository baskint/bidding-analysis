package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/baskint/bidding-analysis/internal/models"
	"github.com/google/uuid"
)

// CampaignStore handles database operations for campaign data
type CampaignStore struct {
	db *sql.DB
}

// NewCampaignStore creates a new CampaignStore instance
func NewCampaignStore(db *sql.DB) *CampaignStore {
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
func (s *CampaignStore) GetDB() *sql.DB {
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
