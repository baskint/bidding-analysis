package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// UserSettings represents user preferences and configuration
type UserSettings struct {
	ID     string `json:"id" db:"id"`
	UserID string `json:"user_id" db:"user_id"`

	// Profile Settings
	FullName *string `json:"full_name" db:"full_name"`
	Email    *string `json:"email" db:"email"`
	Phone    *string `json:"phone" db:"phone"`
	Timezone string  `json:"timezone" db:"timezone"`
	Language string  `json:"language" db:"language"`

	// Notification Preferences
	EmailNotifications   bool   `json:"email_notifications" db:"email_notifications"`
	SlackNotifications   bool   `json:"slack_notifications" db:"slack_notifications"`
	WebhookNotifications bool   `json:"webhook_notifications" db:"webhook_notifications"`
	AlertFrequency       string `json:"alert_frequency" db:"alert_frequency"`

	// Alert Thresholds
	FraudAlertThreshold       float64 `json:"fraud_alert_threshold" db:"fraud_alert_threshold"`
	BudgetAlertThreshold      float64 `json:"budget_alert_threshold" db:"budget_alert_threshold"`
	PerformanceAlertThreshold float64 `json:"performance_alert_threshold" db:"performance_alert_threshold"`

	// Dashboard Preferences
	DefaultDashboardView string `json:"default_dashboard_view" db:"default_dashboard_view"`
	DefaultDateRange     string `json:"default_date_range" db:"default_date_range"`
	DarkMode             bool   `json:"dark_mode" db:"dark_mode"`

	// API Access
	APIKey       *string `json:"api_key,omitempty" db:"api_key"`
	APIRateLimit int     `json:"api_rate_limit" db:"api_rate_limit"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Integration represents a third-party service integration
type Integration struct {
	ID     string `json:"id" db:"id"`
	UserID string `json:"user_id" db:"user_id"`

	// Integration Type
	Provider        string `json:"provider" db:"provider"`
	IntegrationName string `json:"integration_name" db:"integration_name"`

	// Authentication
	AuthType     string  `json:"auth_type" db:"auth_type"`
	AccessToken  *string `json:"access_token,omitempty" db:"access_token"`
	RefreshToken *string `json:"refresh_token,omitempty" db:"refresh_token"`
	APIKey       *string `json:"api_key,omitempty" db:"api_key"`
	APISecret    *string `json:"api_secret,omitempty" db:"api_secret"`
	WebhookURL   *string `json:"webhook_url,omitempty" db:"webhook_url"`

	// Token Expiry
	TokenExpiresAt *time.Time `json:"token_expires_at,omitempty" db:"token_expires_at"`

	// Configuration
	Config IntegrationConfig `json:"config" db:"config"`

	// Status
	Status     string     `json:"status" db:"status"`
	LastSyncAt *time.Time `json:"last_sync_at,omitempty" db:"last_sync_at"`
	LastError  *string    `json:"last_error,omitempty" db:"last_error"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// IntegrationConfig holds provider-specific configuration as JSON
type IntegrationConfig map[string]interface{}

// Value implements the driver.Valuer interface
func (c IntegrationConfig) Value() (driver.Value, error) {
	if c == nil {
		return "{}", nil
	}
	return json.Marshal(c)
}

// Scan implements the sql.Scanner interface
func (c *IntegrationConfig) Scan(value interface{}) error {
	if value == nil {
		*c = make(IntegrationConfig)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return json.Unmarshal(value.([]byte), c)
	}

	return json.Unmarshal(bytes, c)
}

// IntegrationSyncHistory tracks sync operations
type IntegrationSyncHistory struct {
	ID            string `json:"id" db:"id"`
	IntegrationID string `json:"integration_id" db:"integration_id"`

	SyncType string `json:"sync_type" db:"sync_type"`
	Status   string `json:"status" db:"status"`

	RecordsProcessed int `json:"records_processed" db:"records_processed"`
	RecordsFailed    int `json:"records_failed" db:"records_failed"`

	ErrorMessage *string           `json:"error_message,omitempty" db:"error_message"`
	SyncDetails  IntegrationConfig `json:"sync_details" db:"sync_details"`

	StartedAt   time.Time  `json:"started_at" db:"started_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty" db:"completed_at"`
}

// WebhookEvent represents an outgoing webhook event
type WebhookEvent struct {
	ID            string `json:"id" db:"id"`
	IntegrationID string `json:"integration_id" db:"integration_id"`

	EventType string            `json:"event_type" db:"event_type"`
	Payload   IntegrationConfig `json:"payload" db:"payload"`

	Status      string `json:"status" db:"status"`
	Attempts    int    `json:"attempts" db:"attempts"`
	MaxAttempts int    `json:"max_attempts" db:"max_attempts"`

	ResponseStatus *int    `json:"response_status,omitempty" db:"response_status"`
	ResponseBody   *string `json:"response_body,omitempty" db:"response_body"`

	NextRetryAt *time.Time `json:"next_retry_at,omitempty" db:"next_retry_at"`

	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	SentAt    *time.Time `json:"sent_at,omitempty" db:"sent_at"`
}

// BillingInfo represents user billing and subscription information
type BillingInfo struct {
	ID     string `json:"id" db:"id"`
	UserID string `json:"user_id" db:"user_id"`

	// Stripe Integration
	StripeCustomerID     *string `json:"stripe_customer_id,omitempty" db:"stripe_customer_id"`
	StripeSubscriptionID *string `json:"stripe_subscription_id,omitempty" db:"stripe_subscription_id"`

	// Plan Details
	PlanType     string `json:"plan_type" db:"plan_type"`
	BillingCycle string `json:"billing_cycle" db:"billing_cycle"`

	// Limits
	MonthlyBidLimit *int `json:"monthly_bid_limit,omitempty" db:"monthly_bid_limit"`
	CampaignsLimit  *int `json:"campaigns_limit,omitempty" db:"campaigns_limit"`
	MLModelsLimit   *int `json:"ml_models_limit,omitempty" db:"ml_models_limit"`

	// Status
	SubscriptionStatus string     `json:"subscription_status" db:"subscription_status"`
	TrialEndsAt        *time.Time `json:"trial_ends_at,omitempty" db:"trial_ends_at"`
	CurrentPeriodStart *time.Time `json:"current_period_start,omitempty" db:"current_period_start"`
	CurrentPeriodEnd   *time.Time `json:"current_period_end,omitempty" db:"current_period_end"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// UserSettingsUpdate represents fields that can be updated
type UserSettingsUpdate struct {
	FullName                  *string  `json:"full_name,omitempty"`
	Email                     *string  `json:"email,omitempty"`
	Phone                     *string  `json:"phone,omitempty"`
	Timezone                  *string  `json:"timezone,omitempty"`
	Language                  *string  `json:"language,omitempty"`
	EmailNotifications        *bool    `json:"email_notifications,omitempty"`
	SlackNotifications        *bool    `json:"slack_notifications,omitempty"`
	WebhookNotifications      *bool    `json:"webhook_notifications,omitempty"`
	AlertFrequency            *string  `json:"alert_frequency,omitempty"`
	FraudAlertThreshold       *float64 `json:"fraud_alert_threshold,omitempty"`
	BudgetAlertThreshold      *float64 `json:"budget_alert_threshold,omitempty"`
	PerformanceAlertThreshold *float64 `json:"performance_alert_threshold,omitempty"`
	DefaultDashboardView      *string  `json:"default_dashboard_view,omitempty"`
	DefaultDateRange          *string  `json:"default_date_range,omitempty"`
	DarkMode                  *bool    `json:"dark_mode,omitempty"`
}

// IntegrationCreate represents the data needed to create a new integration
type IntegrationCreate struct {
	Provider        string            `json:"provider"`
	IntegrationName string            `json:"integration_name"`
	AuthType        string            `json:"auth_type"`
	AccessToken     *string           `json:"access_token,omitempty"`
	RefreshToken    *string           `json:"refresh_token,omitempty"`
	APIKey          *string           `json:"api_key,omitempty"`
	APISecret       *string           `json:"api_secret,omitempty"`
	WebhookURL      *string           `json:"webhook_url,omitempty"`
	Config          IntegrationConfig `json:"config,omitempty"`
}

// IntegrationUpdate represents fields that can be updated
type IntegrationUpdate struct {
	IntegrationName *string            `json:"integration_name,omitempty"`
	AccessToken     *string            `json:"access_token,omitempty"`
	RefreshToken    *string            `json:"refresh_token,omitempty"`
	APIKey          *string            `json:"api_key,omitempty"`
	APISecret       *string            `json:"api_secret,omitempty"`
	WebhookURL      *string            `json:"webhook_url,omitempty"`
	Config          *IntegrationConfig `json:"config,omitempty"`
	Status          *string            `json:"status,omitempty"`
}
