-- migrations/0002_ml_models.up.sql

-- Create ml_models table
CREATE TABLE IF NOT EXISTS ml_models (
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

CREATE UNIQUE INDEX unique_default_per_type
ON ml_models (user_id, type)
WHERE is_default = true AND deleted_at IS NULL;

-- Create indexes for better query performance
CREATE INDEX idx_ml_models_user_id ON ml_models(user_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_ml_models_type ON ml_models(type) WHERE deleted_at IS NULL;
CREATE INDEX idx_ml_models_status ON ml_models(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_ml_models_is_default ON ml_models(user_id, type, is_default) WHERE is_default = true AND deleted_at IS NULL;
CREATE INDEX idx_ml_models_created_at ON ml_models(created_at DESC) WHERE deleted_at IS NULL;

-- Add comments for documentation
COMMENT ON TABLE ml_models IS 'Stores ML model configurations for bidding analysis';
COMMENT ON COLUMN ml_models.type IS 'Type of model: bidding_optimizer, fraud_detector, conversion_predictor, etc.';
COMMENT ON COLUMN ml_models.config IS 'JSONB field for flexible model configuration parameters';
COMMENT ON COLUMN ml_models.metrics IS 'JSONB field for storing model performance metrics (accuracy, precision, recall, etc.)';
COMMENT ON COLUMN ml_models.is_default IS 'Indicates if this is the default model for its type for the user';
