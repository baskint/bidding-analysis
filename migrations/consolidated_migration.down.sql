-- ============================================================================
-- CONSOLIDATED MIGRATION - DOWN
-- Reverses: 0001_initial_schema, 0002_ml_models, 0003_user_settings_integrations
-- ============================================================================

-- ====================
-- DROP TRIGGERS
-- ====================

DROP TRIGGER IF EXISTS update_billing_info_updated_at ON billing_info;
DROP TRIGGER IF EXISTS update_integrations_updated_at ON integrations;
DROP TRIGGER IF EXISTS update_user_settings_updated_at ON user_settings;
DROP TRIGGER IF EXISTS update_alerts_updated_at ON alerts;
DROP TRIGGER IF EXISTS update_campaigns_updated_at ON campaigns;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- ====================
-- DROP FUNCTIONS
-- ====================

DROP FUNCTION IF EXISTS update_updated_at_column();

-- ====================
-- DROP INDEXES
-- ====================

-- Settings & Integrations indexes
DROP INDEX IF EXISTS idx_billing_info_user_id;
DROP INDEX IF EXISTS idx_webhook_events_status;
DROP INDEX IF EXISTS idx_webhook_events_integration_id;
DROP INDEX IF EXISTS idx_integration_sync_history_integration_id;
DROP INDEX IF EXISTS idx_integrations_status;
DROP INDEX IF EXISTS idx_integrations_provider;
DROP INDEX IF EXISTS idx_integrations_user_id;
DROP INDEX IF EXISTS idx_user_settings_user_id;

-- Alert indexes
DROP INDEX IF EXISTS idx_alerts_user_status;
DROP INDEX IF EXISTS idx_alerts_campaign_id;
DROP INDEX IF EXISTS idx_alerts_type;
DROP INDEX IF EXISTS idx_alerts_severity;
DROP INDEX IF EXISTS idx_alerts_status;
DROP INDEX IF EXISTS idx_alerts_created_at;
DROP INDEX IF EXISTS idx_alerts_user_id;

-- Fraud indexes
DROP INDEX IF EXISTS idx_fraud_metrics_user_date;
DROP INDEX IF EXISTS idx_blocked_entities_type_value;
DROP INDEX IF EXISTS idx_blocked_entities_user;
DROP INDEX IF EXISTS idx_fraud_rules_user;
DROP INDEX IF EXISTS idx_fraud_alerts_detected_at;
DROP INDEX IF EXISTS idx_fraud_alerts_campaign_id;

-- ML model indexes
DROP INDEX IF EXISTS idx_ml_models_created_at;
DROP INDEX IF EXISTS idx_ml_models_is_default;
DROP INDEX IF EXISTS idx_ml_models_status;
DROP INDEX IF EXISTS idx_ml_models_type;
DROP INDEX IF EXISTS idx_ml_models_user_id;
DROP INDEX IF EXISTS unique_default_per_type;
DROP INDEX IF EXISTS idx_ml_model_metadata_type;
DROP INDEX IF EXISTS idx_ml_model_metadata_deployed;
DROP INDEX IF EXISTS idx_ml_model_metadata_created;

-- Core table indexes
DROP INDEX IF EXISTS idx_predictions_bid_event_id;
DROP INDEX IF EXISTS idx_campaign_metrics_campaign_date;
DROP INDEX IF EXISTS idx_bid_events_keywords;
DROP INDEX IF EXISTS idx_bid_events_campaign;
DROP INDEX IF EXISTS idx_bid_events_user_timestamp;
DROP INDEX IF EXISTS idx_bid_events_user_id;
DROP INDEX IF EXISTS idx_bid_events_timestamp;
DROP INDEX IF EXISTS idx_bid_events_campaign_id;

-- ====================
-- DROP TABLES
-- ====================
-- Drop in reverse order of dependencies to avoid foreign key violations

-- Settings & Integrations tables
DROP TABLE IF EXISTS webhook_events CASCADE;
DROP TABLE IF EXISTS integration_sync_history CASCADE;
DROP TABLE IF EXISTS integrations CASCADE;
DROP TABLE IF EXISTS billing_info CASCADE;
DROP TABLE IF EXISTS user_settings CASCADE;

-- Alerts tables
DROP TABLE IF EXISTS alerts CASCADE;

-- Fraud tables
DROP TABLE IF EXISTS fraud_metrics CASCADE;
DROP TABLE IF EXISTS blocked_entities CASCADE;
DROP TABLE IF EXISTS fraud_rules CASCADE;
DROP TABLE IF EXISTS fraud_alerts CASCADE;

-- Analytics & Metrics tables
DROP TABLE IF EXISTS campaign_metrics CASCADE;

-- ML tables
DROP TABLE IF EXISTS model_metrics CASCADE;
DROP TABLE IF EXISTS ml_models CASCADE;
DROP TABLE IF EXISTS ml_model_metadata CASCADE;
DROP TABLE IF EXISTS predictions CASCADE;

-- Core tables
DROP TABLE IF EXISTS bid_events CASCADE;
DROP TABLE IF EXISTS campaigns CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- ====================
-- DROP TYPES
-- ====================

DROP TYPE IF EXISTS alert_status;
DROP TYPE IF EXISTS alert_severity;
DROP TYPE IF EXISTS alert_type;

-- ====================
-- DROP EXTENSIONS
-- ====================

DROP EXTENSION IF EXISTS "uuid-ossp";
