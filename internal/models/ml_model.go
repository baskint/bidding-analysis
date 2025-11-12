// internal/models/ml_model.go
package models

import (
	"time"

	"github.com/google/uuid"
)

// MLModel represents a machine learning model configuration
type MLModel struct {
	ID          uuid.UUID              `json:"id" db:"id"`
	UserID      uuid.UUID              `json:"user_id" db:"user_id"`
	Name        string                 `json:"name" db:"name"`
	Type        string                 `json:"type" db:"type"` // e.g., "bidding_optimizer", "fraud_detector", "conversion_predictor"
	Version     string                 `json:"version" db:"version"`
	Description string                 `json:"description" db:"description"`
	Status      string                 `json:"status" db:"status"` // "active", "inactive", "training", "testing"
	Provider    string                 `json:"provider" db:"provider"` // e.g., "openai", "custom", "tensorflow", "pytorch"
	Endpoint    *string                `json:"endpoint,omitempty" db:"endpoint"` // API endpoint if applicable
	APIKey      *string                `json:"-" db:"api_key"` // Don't expose in JSON
	Config      map[string]interface{} `json:"config" db:"config"` // JSONB field for flexible configuration
	Metrics     map[string]interface{} `json:"metrics" db:"metrics"` // JSONB field for performance metrics
	IsDefault   bool                   `json:"is_default" db:"is_default"` // Whether this is the default model for its type
	CreatedAt   time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at" db:"updated_at"`
	DeletedAt   *time.Time             `json:"deleted_at,omitempty" db:"deleted_at"`
}

// MLModelCreate represents the input for creating a new model
type MLModelCreate struct {
	Name        string                 `json:"name" binding:"required"`
	Type        string                 `json:"type" binding:"required"`
	Version     string                 `json:"version" binding:"required"`
	Description string                 `json:"description"`
	Provider    string                 `json:"provider" binding:"required"`
	Endpoint    *string                `json:"endpoint,omitempty"`
	APIKey      *string                `json:"api_key,omitempty"`
	Config      map[string]interface{} `json:"config"`
	IsDefault   bool                   `json:"is_default"`
}

// MLModelUpdate represents the input for updating a model
type MLModelUpdate struct {
	Name        *string                 `json:"name,omitempty"`
	Version     *string                 `json:"version,omitempty"`
	Description *string                 `json:"description,omitempty"`
	Status      *string                 `json:"status,omitempty"`
	Provider    *string                 `json:"provider,omitempty"`
	Endpoint    *string                 `json:"endpoint,omitempty"`
	APIKey      *string                 `json:"api_key,omitempty"`
	Config      *map[string]interface{} `json:"config,omitempty"`
	Metrics     *map[string]interface{} `json:"metrics,omitempty"`
	IsDefault   *bool                   `json:"is_default,omitempty"`
}

// MLModelListResponse represents a list of models with pagination
type MLModelListResponse struct {
	Models     []MLModel `json:"models"`
	Total      int       `json:"total"`
	Page       int       `json:"page"`
	PageSize   int       `json:"page_size"`
	TotalPages int       `json:"total_pages"`
}
