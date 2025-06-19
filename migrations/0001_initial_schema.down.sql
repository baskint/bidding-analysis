-- Drop triggers
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TRIGGER IF EXISTS update_campaigns_updated_at ON campaigns;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS idx_bid_events_campaign_id;
DROP INDEX IF EXISTS idx_bid_events_timestamp;
DROP INDEX IF EXISTS idx_bid_events_user_id;
DROP INDEX IF EXISTS idx_campaign_metrics_campaign_date;
DROP INDEX IF EXISTS idx_predictions_bid_event_id;
DROP INDEX IF EXISTS idx_fraud_alerts_campaign_id;
DROP INDEX IF EXISTS idx_fraud_alerts_detected_at;

-- Drop tables (in reverse order of dependencies)
DROP TABLE IF EXISTS model_metrics;
DROP TABLE IF EXISTS fraud_alerts;
DROP TABLE IF EXISTS campaign_metrics;
DROP TABLE IF EXISTS predictions;
DROP TABLE IF EXISTS bid_events;
DROP TABLE IF EXISTS campaigns;
DROP TABLE IF EXISTS users;

-- Drop extensions
DROP EXTENSION IF EXISTS "uuid-ossp";
