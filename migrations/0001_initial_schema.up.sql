-- Create extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

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

CREATE INDEX idx_fraud_rules_user ON fraud_rules(user_id);
CREATE INDEX idx_blocked_entities_user ON blocked_entities(user_id);
CREATE INDEX idx_blocked_entities_type_value ON blocked_entities(entity_type, entity_value);
CREATE INDEX idx_fraud_metrics_user_date ON fraud_metrics(user_id, date);

-- Indexes for performance
CREATE INDEX idx_bid_events_campaign_id ON bid_events(campaign_id);
CREATE INDEX idx_bid_events_timestamp ON bid_events(timestamp);
CREATE INDEX idx_bid_events_user_id ON bid_events(user_id);
CREATE INDEX idx_campaign_metrics_campaign_date ON campaign_metrics(campaign_id, date);
CREATE INDEX idx_predictions_bid_event_id ON predictions(bid_event_id);
CREATE INDEX idx_fraud_alerts_campaign_id ON fraud_alerts(campaign_id);
CREATE INDEX idx_fraud_alerts_detected_at ON fraud_alerts(detected_at);
CREATE INDEX idx_bid_events_user_timestamp ON bid_events(user_id, timestamp);
CREATE INDEX idx_bid_events_campaign ON bid_events(campaign_id);
CREATE INDEX idx_bid_events_keywords ON bid_events USING gin(keywords);

-- Triggers for updated_at timestamps
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_campaigns_updated_at BEFORE UPDATE ON campaigns 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
    