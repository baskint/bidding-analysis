-- migrations/0002_ml_models.down.sql

-- Drop indexes
DROP INDEX IF EXISTS idx_ml_models_created_at;
DROP INDEX IF EXISTS idx_ml_models_is_default;
DROP INDEX IF EXISTS idx_ml_models_status;
DROP INDEX IF EXISTS idx_ml_models_type;
DROP INDEX IF EXISTS idx_ml_models_user_id;

-- Drop the ml_models table
DROP TABLE IF EXISTS ml_models;
