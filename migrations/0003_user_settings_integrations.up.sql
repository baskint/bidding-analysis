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

-- Indexes
CREATE INDEX idx_user_settings_user_id ON user_settings(user_id);
CREATE INDEX idx_integrations_user_id ON integrations(user_id);
CREATE INDEX idx_integrations_provider ON integrations(provider);
CREATE INDEX idx_integrations_status ON integrations(status);
CREATE INDEX idx_integration_sync_history_integration_id ON integration_sync_history(integration_id);
CREATE INDEX idx_webhook_events_integration_id ON webhook_events(integration_id);
CREATE INDEX idx_webhook_events_status ON webhook_events(status);
CREATE INDEX idx_billing_info_user_id ON billing_info(user_id);

-- Triggers
CREATE TRIGGER update_user_settings_updated_at BEFORE UPDATE ON user_settings 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_integrations_updated_at BEFORE UPDATE ON integrations 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_billing_info_updated_at BEFORE UPDATE ON billing_info 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
