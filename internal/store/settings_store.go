package store

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/baskint/bidding-analysis/internal/models"
	"github.com/jmoiron/sqlx"
)

// SettingsStore handles database operations for user settings
type SettingsStore struct {
	db *sqlx.DB
}

// NewSettingsStore creates a new settings store
func NewSettingsStore(db *sqlx.DB) *SettingsStore {
	return &SettingsStore{db: db}
}

// GetUserSettings retrieves user settings, creating default if not exists
func (s *SettingsStore) GetUserSettings(ctx context.Context, userID string) (*models.UserSettings, error) {
	var settings models.UserSettings
	query := `SELECT * FROM user_settings WHERE user_id = $1`

	err := s.db.GetContext(ctx, &settings, query, userID)
	if err == sql.ErrNoRows {
		// Create default settings
		return s.CreateDefaultSettings(ctx, userID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user settings: %w", err)
	}

	return &settings, nil
}

// CreateDefaultSettings creates default settings for a new user
func (s *SettingsStore) CreateDefaultSettings(ctx context.Context, userID string) (*models.UserSettings, error) {
	query := `
		INSERT INTO user_settings (
			user_id, timezone, language, email_notifications, 
			alert_frequency, fraud_alert_threshold, budget_alert_threshold,
			performance_alert_threshold, default_dashboard_view, 
			default_date_range, dark_mode, api_rate_limit
		) VALUES (
			$1, 'UTC', 'en', true, 'realtime', 0.7, 0.8, 0.5, 
			'overview', '7d', false, 1000
		)
		RETURNING *
	`

	var settings models.UserSettings
	err := s.db.GetContext(ctx, &settings, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to create default settings: %w", err)
	}

	return &settings, nil
}

// UpdateUserSettings updates user settings
func (s *SettingsStore) UpdateUserSettings(ctx context.Context, userID string, update *models.UserSettingsUpdate) (*models.UserSettings, error) {
	// Build dynamic update query
	query := `UPDATE user_settings SET updated_at = NOW()`
	args := []interface{}{}
	argPos := 1

	if update.FullName != nil {
		argPos++
		query += fmt.Sprintf(", full_name = $%d", argPos)
		args = append(args, *update.FullName)
	}
	if update.Email != nil {
		argPos++
		query += fmt.Sprintf(", email = $%d", argPos)
		args = append(args, *update.Email)
	}
	if update.Phone != nil {
		argPos++
		query += fmt.Sprintf(", phone = $%d", argPos)
		args = append(args, *update.Phone)
	}
	if update.Timezone != nil {
		argPos++
		query += fmt.Sprintf(", timezone = $%d", argPos)
		args = append(args, *update.Timezone)
	}
	if update.Language != nil {
		argPos++
		query += fmt.Sprintf(", language = $%d", argPos)
		args = append(args, *update.Language)
	}
	if update.EmailNotifications != nil {
		argPos++
		query += fmt.Sprintf(", email_notifications = $%d", argPos)
		args = append(args, *update.EmailNotifications)
	}
	if update.SlackNotifications != nil {
		argPos++
		query += fmt.Sprintf(", slack_notifications = $%d", argPos)
		args = append(args, *update.SlackNotifications)
	}
	if update.WebhookNotifications != nil {
		argPos++
		query += fmt.Sprintf(", webhook_notifications = $%d", argPos)
		args = append(args, *update.WebhookNotifications)
	}
	if update.AlertFrequency != nil {
		argPos++
		query += fmt.Sprintf(", alert_frequency = $%d", argPos)
		args = append(args, *update.AlertFrequency)
	}
	if update.FraudAlertThreshold != nil {
		argPos++
		query += fmt.Sprintf(", fraud_alert_threshold = $%d", argPos)
		args = append(args, *update.FraudAlertThreshold)
	}
	if update.BudgetAlertThreshold != nil {
		argPos++
		query += fmt.Sprintf(", budget_alert_threshold = $%d", argPos)
		args = append(args, *update.BudgetAlertThreshold)
	}
	if update.PerformanceAlertThreshold != nil {
		argPos++
		query += fmt.Sprintf(", performance_alert_threshold = $%d", argPos)
		args = append(args, *update.PerformanceAlertThreshold)
	}
	if update.DefaultDashboardView != nil {
		argPos++
		query += fmt.Sprintf(", default_dashboard_view = $%d", argPos)
		args = append(args, *update.DefaultDashboardView)
	}
	if update.DefaultDateRange != nil {
		argPos++
		query += fmt.Sprintf(", default_date_range = $%d", argPos)
		args = append(args, *update.DefaultDateRange)
	}
	if update.DarkMode != nil {
		argPos++
		query += fmt.Sprintf(", dark_mode = $%d", argPos)
		args = append(args, *update.DarkMode)
	}

	query += " WHERE user_id = $1 RETURNING *"
	args = append([]interface{}{userID}, args...)

	var settings models.UserSettings
	err := s.db.GetContext(ctx, &settings, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update user settings: %w", err)
	}

	return &settings, nil
}

// RegenerateAPIKey generates a new API key for the user
func (s *SettingsStore) RegenerateAPIKey(ctx context.Context, userID string) (string, error) {
	apiKey := generateAPIKey()

	query := `UPDATE user_settings SET api_key = $1, updated_at = NOW() WHERE user_id = $2`
	_, err := s.db.ExecContext(ctx, query, apiKey, userID)
	if err != nil {
		return "", fmt.Errorf("failed to regenerate API key: %w", err)
	}

	return apiKey, nil
}

// Integration Store Methods

// ListIntegrations retrieves all integrations for a user
func (s *SettingsStore) ListIntegrations(ctx context.Context, userID string) ([]models.Integration, error) {
	var integrations []models.Integration
	query := `SELECT * FROM integrations WHERE user_id = $1 ORDER BY created_at DESC`

	err := s.db.SelectContext(ctx, &integrations, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list integrations: %w", err)
	}

	return integrations, nil
}

// GetIntegration retrieves a specific integration
func (s *SettingsStore) GetIntegration(ctx context.Context, id, userID string) (*models.Integration, error) {
	var integration models.Integration
	query := `SELECT * FROM integrations WHERE id = $1 AND user_id = $2`

	err := s.db.GetContext(ctx, &integration, query, id, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get integration: %w", err)
	}

	return &integration, nil
}

// CreateIntegration creates a new integration
func (s *SettingsStore) CreateIntegration(ctx context.Context, userID string, create *models.IntegrationCreate) (*models.Integration, error) {
	query := `
		INSERT INTO integrations (
			user_id, provider, integration_name, auth_type,
			access_token, refresh_token, api_key, api_secret,
			webhook_url, config, status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, 'active')
		RETURNING *
	`

	var integration models.Integration
	err := s.db.GetContext(ctx, &integration, query,
		userID,
		create.Provider,
		create.IntegrationName,
		create.AuthType,
		create.AccessToken,
		create.RefreshToken,
		create.APIKey,
		create.APISecret,
		create.WebhookURL,
		create.Config,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create integration: %w", err)
	}

	return &integration, nil
}

// UpdateIntegration updates an existing integration
func (s *SettingsStore) UpdateIntegration(ctx context.Context, id, userID string, update *models.IntegrationUpdate) (*models.Integration, error) {
	query := `UPDATE integrations SET updated_at = NOW()`
	args := []interface{}{}
	argPos := 1

	if update.IntegrationName != nil {
		argPos++
		query += fmt.Sprintf(", integration_name = $%d", argPos)
		args = append(args, *update.IntegrationName)
	}
	if update.AccessToken != nil {
		argPos++
		query += fmt.Sprintf(", access_token = $%d", argPos)
		args = append(args, *update.AccessToken)
	}
	if update.RefreshToken != nil {
		argPos++
		query += fmt.Sprintf(", refresh_token = $%d", argPos)
		args = append(args, *update.RefreshToken)
	}
	if update.APIKey != nil {
		argPos++
		query += fmt.Sprintf(", api_key = $%d", argPos)
		args = append(args, *update.APIKey)
	}
	if update.APISecret != nil {
		argPos++
		query += fmt.Sprintf(", api_secret = $%d", argPos)
		args = append(args, *update.APISecret)
	}
	if update.WebhookURL != nil {
		argPos++
		query += fmt.Sprintf(", webhook_url = $%d", argPos)
		args = append(args, *update.WebhookURL)
	}
	if update.Config != nil {
		argPos++
		query += fmt.Sprintf(", config = $%d", argPos)
		args = append(args, *update.Config)
	}
	if update.Status != nil {
		argPos++
		query += fmt.Sprintf(", status = $%d", argPos)
		args = append(args, *update.Status)
	}

	query += " WHERE id = $1 AND user_id = $2 RETURNING *"
	args = append([]interface{}{id, userID}, args...)

	var integration models.Integration
	err := s.db.GetContext(ctx, &integration, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update integration: %w", err)
	}

	return &integration, nil
}

// DeleteIntegration deletes an integration
func (s *SettingsStore) DeleteIntegration(ctx context.Context, id, userID string) error {
	query := `DELETE FROM integrations WHERE id = $1 AND user_id = $2`

	result, err := s.db.ExecContext(ctx, query, id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete integration: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("integration not found")
	}

	return nil
}

// TestIntegration tests an integration connection
func (s *SettingsStore) TestIntegration(ctx context.Context, id, userID string) error {
	// This would implement actual testing logic based on provider
	// For now, just update last_sync_at
	query := `UPDATE integrations SET last_sync_at = NOW() WHERE id = $1 AND user_id = $2`
	_, err := s.db.ExecContext(ctx, query, id, userID)
	return err
}

// UpdateIntegrationStatus updates the status of an integration
func (s *SettingsStore) UpdateIntegrationStatus(ctx context.Context, id, userID, status string, lastError *string) error {
	query := `
		UPDATE integrations 
		SET status = $1, last_error = $2, updated_at = NOW() 
		WHERE id = $3 AND user_id = $4
	`
	_, err := s.db.ExecContext(ctx, query, status, lastError, id, userID)
	return err
}

// GetBillingInfo retrieves billing information
func (s *SettingsStore) GetBillingInfo(ctx context.Context, userID string) (*models.BillingInfo, error) {
	var billing models.BillingInfo
	query := `SELECT * FROM billing_info WHERE user_id = $1`

	err := s.db.GetContext(ctx, &billing, query, userID)
	if err == sql.ErrNoRows {
		// Create default billing info
		return s.CreateDefaultBillingInfo(ctx, userID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get billing info: %w", err)
	}

	return &billing, nil
}

// CreateDefaultBillingInfo creates default billing info for new user
func (s *SettingsStore) CreateDefaultBillingInfo(ctx context.Context, userID string) (*models.BillingInfo, error) {
	trialEnds := time.Now().Add(14 * 24 * time.Hour) // 14-day trial

	query := `
		INSERT INTO billing_info (
			user_id, plan_type, subscription_status, trial_ends_at
		) VALUES ($1, 'free', 'trial', $2)
		RETURNING *
	`

	var billing models.BillingInfo
	err := s.db.GetContext(ctx, &billing, query, userID, trialEnds)
	if err != nil {
		return nil, fmt.Errorf("failed to create default billing info: %w", err)
	}

	return &billing, nil
}

// Helper function to generate API keys
func generateAPIKey() string {
	// Generate 32 random bytes
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to timestamp-based key (less secure but functional)
		return fmt.Sprintf("ba_%d_%d", time.Now().Unix(), time.Now().UnixNano())
	}
	// Return as hex string with prefix
	return fmt.Sprintf("ba_%s", hex.EncodeToString(bytes))
}
