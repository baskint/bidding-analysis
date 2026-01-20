-- ============================================================================
-- CONSOLIDATED MIGRATION - UP
-- Combines: 0001_initial_schema, 0002_ml_models, 0003_user_settings_integrations
-- ============================================================================

-- ====================
-- EXTENSIONS
-- ====================
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ====================
-- CORE TABLES
-- ====================

-- Users table for authentication
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Campaigns table
CREATE TABLE campaigns (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    user_id UUID REFERENCES users(id),
    status VARCHAR(50) DEFAULT 'active',
    budget DECIMAL(10,2),
    daily_budget DECIMAL(10,2),
    target_cpa DECIMAL(10,2),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Bid events table (main table for real-time data)
CREATE TABLE bid_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    campaign_id UUID REFERENCES campaigns(id),
    user_id UUID REFERENCES users(id),
    bid_price DECIMAL(8,4) NOT NULL,
    win_price DECIMAL(8,4),
    floor_price DECIMAL(8,4),
    won BOOLEAN DEFAULT FALSE,
    converted BOOLEAN DEFAULT FALSE,

    -- User segment data
    segment_id VARCHAR(50),
    segment_category VARCHAR(50),
    engagement_score DECIMAL(3,2),
    conversion_probability DECIMAL(3,2),

    -- Geo data
    country VARCHAR(2),
    region VARCHAR(100),
    city VARCHAR(100),
    latitude DECIMAL(10,8),
    longitude DECIMAL(11,8),

    -- Device data
    device_type VARCHAR(20),
    os VARCHAR(50),
    browser VARCHAR(50),
    is_mobile BOOLEAN DEFAULT FALSE,

    -- Keywords
    keywords TEXT[],

    -- Timestamps
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- ====================
-- ML & PREDICTIONS
-- ====================

-- Predictions table (ML model outputs)
CREATE TABLE predictions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    bid_event_id UUID REFERENCES bid_events(id),
    predicted_bid_price DECIMAL(8,4) NOT NULL,
    confidence DECIMAL(3,2) NOT NULL,
    strategy VARCHAR(100),
    fraud_risk BOOLEAN DEFAULT FALSE,
    model_version VARCHAR(20),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- ML model metadata (training runs)
CREATE TABLE ml_model_metadata (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Model identification
    model_path VARCHAR(500) NOT NULL,
    model_type VARCHAR(100),

    -- Performance metrics
    train_rmse DECIMAL(10,6),
    val_rmse DECIMAL(10,6),
    train_r2 DECIMAL(10,6),
    val_r2 DECIMAL(10,6),
    train_mae DECIMAL(10,6),
    val_mae DECIMAL(10,6),

    -- Feature importance (JSON)
    feature_importance JSONB,

    -- Training metadata
    training_samples INTEGER,
    validation_samples INTEGER,
    hyperparameters JSONB,

    -- Deployment tracking
    deployed BOOLEAN DEFAULT FALSE,
    deployment_date TIMESTAMP WITH TIME ZONE,

    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- ML models table (user-configured models)
CREATE TABLE ml_models (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(100) NOT NULL,  -- e.g., 'bidding_optimizer', 'fraud_detector', 'conversion_predictor'
    version VARCHAR(50) NOT NULL,
    description TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'inactive',  -- 'active', 'inactive', 'training', 'testing'
    provider VARCHAR(100) NOT NULL,  -- e.g., 'openai', 'custom', 'tensorflow', 'pytorch'
    endpoint VARCHAR(500),  -- API endpoint if applicable
    api_key TEXT,  -- Encrypted API key if needed
    config JSONB NOT NULL DEFAULT '{}',  -- Flexible configuration storage
    metrics JSONB NOT NULL DEFAULT '{}',  -- Performance metrics storage
    is_default BOOLEAN NOT NULL DEFAULT false,  -- Whether this is the default model for its type
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Model performance tracking
CREATE TABLE model_metrics (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    model_version VARCHAR(20) NOT NULL,
    date DATE NOT NULL,

    -- Accuracy metrics
    prediction_accuracy DECIMAL(5,4),
    mean_absolute_error DECIMAL(8,4),
    root_mean_square_error DECIMAL(8,4),
    total_predictions INTEGER,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    UNIQUE(model_version, date)
);

-- ====================
-- ANALYTICS & METRICS
-- ====================

-- Campaign metrics (aggregated data for analytics)
CREATE TABLE campaign_metrics (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    campaign_id UUID REFERENCES campaigns(id),
    date DATE NOT NULL,
    hour INTEGER, -- 0-23 for hourly aggregation

    -- Metrics
    total_bids INTEGER DEFAULT 0,
    won_bids INTEGER DEFAULT 0,
    conversions INTEGER DEFAULT 0,
    total_spend DECIMAL(10,2) DEFAULT 0,
    impressions INTEGER DEFAULT 0,
    clicks INTEGER DEFAULT 0,

    -- Calculated metrics
    win_rate DECIMAL(5,4),
    conversion_rate DECIMAL(5,4),
    average_bid DECIMAL(8,4),
    cost_per_conversion DECIMAL(8,4),
    return_on_ad_spend DECIMAL(8,4),

    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    UNIQUE(campaign_id, date, hour)
);

-- ====================
-- FRAUD DETECTION
-- ====================

-- Fraud alerts table
CREATE TABLE fraud_alerts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    campaign_id UUID REFERENCES campaigns(id),
    alert_type VARCHAR(50) NOT NULL, -- click_fraud, impression_fraud, etc.
    severity INTEGER CHECK (severity BETWEEN 1 AND 10),
    description TEXT,
    affected_user_ids TEXT[],
    detected_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    resolved_at TIMESTAMP WITH TIME ZONE,
    status VARCHAR(20) DEFAULT 'active' -- active, resolved, false_positive
);

-- Fraud rules table
CREATE TABLE fraud_rules (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id),
    name VARCHAR(255) NOT NULL,
    rule_type VARCHAR(50) NOT NULL,
    conditions JSONB NOT NULL,
    threshold DECIMAL(5,2),
    severity INTEGER CHECK (severity BETWEEN 1 AND 10),
    enabled BOOLEAN DEFAULT TRUE,
    auto_block BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Blocklist table
CREATE TABLE blocked_entities (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id),
    entity_type VARCHAR(20) NOT NULL, -- 'ip', 'device', 'user'
    entity_value VARCHAR(255) NOT NULL,
    reason TEXT,
    blocked_by_rule_id UUID REFERENCES fraud_rules(id),
    blocked_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    expires_at TIMESTAMP WITH TIME ZONE,
    permanent BOOLEAN DEFAULT FALSE,
    created_by UUID REFERENCES users(id)
);

-- Fraud metrics aggregation table
CREATE TABLE fraud_metrics (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id),
    date DATE NOT NULL,
    hour INTEGER,
    fraud_attempts INTEGER DEFAULT 0,
    blocked_bids INTEGER DEFAULT 0,
    amount_saved DECIMAL(10,2) DEFAULT 0,
    false_positives INTEGER DEFAULT 0,
    alert_type VARCHAR(50),
    campaign_id UUID REFERENCES campaigns(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, date, hour, alert_type, campaign_id)
);

-- ====================
-- ALERTS SYSTEM
-- ====================

-- Alert types
CREATE TYPE alert_type AS ENUM ('fraud', 'budget', 'performance', 'model', 'system', 'campaign');
CREATE TYPE alert_severity AS ENUM ('low', 'medium', 'high', 'critical');
CREATE TYPE alert_status AS ENUM ('unread', 'read', 'acknowledged', 'resolved', 'dismissed');

-- Alerts table
CREATE TABLE alerts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) NOT NULL,
    type alert_type NOT NULL,
    severity alert_severity NOT NULL,
    status alert_status DEFAULT 'unread' NOT NULL,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    campaign_id UUID REFERENCES campaigns(id),
    metadata JSONB,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    acknowledged_at TIMESTAMP WITH TIME ZONE,
    resolved_at TIMESTAMP WITH TIME ZONE
);

-- ====================
-- USER SETTINGS & INTEGRATIONS
-- ====================

-- User Settings and Preferences
CREATE TABLE user_settings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) UNIQUE NOT NULL,

    -- Profile Settings
    full_name VARCHAR(255),
    email VARCHAR(255),
    phone VARCHAR(50),
    timezone VARCHAR(50) DEFAULT 'UTC',
    language VARCHAR(10) DEFAULT 'en',

    -- Notification Preferences
    email_notifications BOOLEAN DEFAULT true,
    slack_notifications BOOLEAN DEFAULT false,
    webhook_notifications BOOLEAN DEFAULT false,
    alert_frequency VARCHAR(20) DEFAULT 'realtime', -- realtime, hourly, daily

    -- Alert Thresholds
    fraud_alert_threshold DECIMAL(5,2) DEFAULT 0.7,
    budget_alert_threshold DECIMAL(5,2) DEFAULT 0.8,
    performance_alert_threshold DECIMAL(5,2) DEFAULT 0.5,

    -- Dashboard Preferences
    default_dashboard_view VARCHAR(50) DEFAULT 'overview',
    default_date_range VARCHAR(20) DEFAULT '7d',
    dark_mode BOOLEAN DEFAULT false,

    -- API Access
    api_key VARCHAR(255) UNIQUE,
    api_rate_limit INTEGER DEFAULT 1000,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Third-Party Integrations
CREATE TABLE integrations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) NOT NULL,

    -- Integration Type
    provider VARCHAR(50) NOT NULL, -- google_ads, facebook_ads, slack, etc.
    integration_name VARCHAR(255) NOT NULL,

    -- Authentication
    auth_type VARCHAR(20) NOT NULL, -- oauth, api_key, webhook
    access_token TEXT,
    refresh_token TEXT,
    api_key TEXT,
    api_secret TEXT,
    webhook_url TEXT,

    -- Token Expiry
    token_expires_at TIMESTAMP WITH TIME ZONE,

    -- Configuration
    config JSONB DEFAULT '{}',

    -- Status
    status VARCHAR(20) DEFAULT 'active', -- active, inactive, error, expired
    last_sync_at TIMESTAMP WITH TIME ZONE,
    last_error TEXT,

    -- Metadata
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    UNIQUE(user_id, provider, integration_name)
);

-- Integration Sync History
CREATE TABLE integration_sync_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    integration_id UUID REFERENCES integrations(id) ON DELETE CASCADE,

    sync_type VARCHAR(50) NOT NULL, -- import, export, sync
    status VARCHAR(20) NOT NULL, -- success, failed, partial

    records_processed INTEGER DEFAULT 0,
    records_failed INTEGER DEFAULT 0,

    error_message TEXT,
    sync_details JSONB DEFAULT '{}',

    started_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    completed_at TIMESTAMP WITH TIME ZONE
);

-- Webhook Events Log
CREATE TABLE webhook_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    integration_id UUID REFERENCES integrations(id) ON DELETE CASCADE,

    event_type VARCHAR(100) NOT NULL,
    payload JSONB NOT NULL,

    status VARCHAR(20) DEFAULT 'pending', -- pending, sent, failed, retry
    attempts INTEGER DEFAULT 0,
    max_attempts INTEGER DEFAULT 3,

    response_status INTEGER,
    response_body TEXT,

    next_retry_at TIMESTAMP WITH TIME ZONE,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    sent_at TIMESTAMP WITH TIME ZONE
);

-- Billing & Subscription
CREATE TABLE billing_info (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) UNIQUE NOT NULL,

    -- Stripe Integration
    stripe_customer_id VARCHAR(255),
    stripe_subscription_id VARCHAR(255),

    -- Plan Details
    plan_type VARCHAR(50) DEFAULT 'free', -- free, starter, professional, enterprise
    billing_cycle VARCHAR(20) DEFAULT 'monthly', -- monthly, yearly

    -- Limits
    monthly_bid_limit INTEGER,
    campaigns_limit INTEGER,
    ml_models_limit INTEGER,

    -- Status
    subscription_status VARCHAR(20) DEFAULT 'active', -- active, past_due, canceled, trial
    trial_ends_at TIMESTAMP WITH TIME ZONE,
    current_period_start TIMESTAMP WITH TIME ZONE,
    current_period_end TIMESTAMP WITH TIME ZONE,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- ====================
-- INDEXES
-- ====================

-- Core table indexes
CREATE INDEX idx_bid_events_campaign_id ON bid_events(campaign_id);
CREATE INDEX idx_bid_events_timestamp ON bid_events(timestamp);
CREATE INDEX idx_bid_events_user_id ON bid_events(user_id);
CREATE INDEX idx_bid_events_user_timestamp ON bid_events(user_id, timestamp);
CREATE INDEX idx_bid_events_campaign ON bid_events(campaign_id);
CREATE INDEX idx_bid_events_keywords ON bid_events USING gin(keywords);

CREATE INDEX idx_campaign_metrics_campaign_date ON campaign_metrics(campaign_id, date);

CREATE INDEX idx_predictions_bid_event_id ON predictions(bid_event_id);

-- ML model indexes
CREATE INDEX idx_ml_model_metadata_created ON ml_model_metadata(created_at DESC);
CREATE INDEX idx_ml_model_metadata_deployed ON ml_model_metadata(deployed, created_at DESC);
CREATE INDEX idx_ml_model_metadata_type ON ml_model_metadata(model_type);

CREATE UNIQUE INDEX unique_default_per_type ON ml_models (user_id, type) WHERE is_default = true AND deleted_at IS NULL;
CREATE INDEX idx_ml_models_user_id ON ml_models(user_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_ml_models_type ON ml_models(type) WHERE deleted_at IS NULL;
CREATE INDEX idx_ml_models_status ON ml_models(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_ml_models_is_default ON ml_models(user_id, type, is_default) WHERE is_default = true AND deleted_at IS NULL;
CREATE INDEX idx_ml_models_created_at ON ml_models(created_at DESC) WHERE deleted_at IS NULL;

-- Fraud indexes
CREATE INDEX idx_fraud_alerts_campaign_id ON fraud_alerts(campaign_id);
CREATE INDEX idx_fraud_alerts_detected_at ON fraud_alerts(detected_at);
CREATE INDEX idx_fraud_rules_user ON fraud_rules(user_id);
CREATE INDEX idx_blocked_entities_user ON blocked_entities(user_id);
CREATE INDEX idx_blocked_entities_type_value ON blocked_entities(entity_type, entity_value);
CREATE INDEX idx_fraud_metrics_user_date ON fraud_metrics(user_id, date);

-- Alert indexes
CREATE INDEX idx_alerts_user_id ON alerts(user_id);
CREATE INDEX idx_alerts_created_at ON alerts(created_at DESC);
CREATE INDEX idx_alerts_status ON alerts(status);
CREATE INDEX idx_alerts_severity ON alerts(severity);
CREATE INDEX idx_alerts_type ON alerts(type);
CREATE INDEX idx_alerts_campaign_id ON alerts(campaign_id) WHERE campaign_id IS NOT NULL;
CREATE INDEX idx_alerts_user_status ON alerts(user_id, status, created_at DESC);

-- Settings & Integrations indexes
CREATE INDEX idx_user_settings_user_id ON user_settings(user_id);
CREATE INDEX idx_integrations_user_id ON integrations(user_id);
CREATE INDEX idx_integrations_provider ON integrations(provider);
CREATE INDEX idx_integrations_status ON integrations(status);
CREATE INDEX idx_integration_sync_history_integration_id ON integration_sync_history(integration_id);
CREATE INDEX idx_webhook_events_integration_id ON webhook_events(integration_id);
CREATE INDEX idx_webhook_events_status ON webhook_events(status);
CREATE INDEX idx_billing_info_user_id ON billing_info(user_id);

-- ====================
-- FUNCTIONS & TRIGGERS
-- ====================

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply triggers to relevant tables
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_campaigns_updated_at BEFORE UPDATE ON campaigns
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_alerts_updated_at BEFORE UPDATE ON alerts
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_settings_updated_at BEFORE UPDATE ON user_settings
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_integrations_updated_at BEFORE UPDATE ON integrations
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_billing_info_updated_at BEFORE UPDATE ON billing_info
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ====================
-- COMMENTS
-- ====================

COMMENT ON TABLE ml_model_metadata IS 'Tracks ML model training runs and performance metrics';
COMMENT ON TABLE ml_models IS 'Stores ML model configurations for bidding analysis';
COMMENT ON COLUMN ml_models.type IS 'Type of model: bidding_optimizer, fraud_detector, conversion_predictor, etc.';
COMMENT ON COLUMN ml_models.config IS 'JSONB field for flexible model configuration parameters';
COMMENT ON COLUMN ml_models.metrics IS 'JSONB field for storing model performance metrics (accuracy, precision, recall, etc.)';
COMMENT ON COLUMN ml_models.is_default IS 'Indicates if this is the default model for its type for the user';
