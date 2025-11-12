// internal/store/ml_model_store.go
package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/baskint/bidding-analysis/internal/models"
	"github.com/google/uuid"
)

// MLModelStore handles database operations for ML models
type MLModelStore struct {
	db *sql.DB
}

// NewMLModelStore creates a new ML model store
func NewMLModelStore(db *sql.DB) *MLModelStore {
	return &MLModelStore{db: db}
}

// Create creates a new ML model
func (s *MLModelStore) Create(ctx context.Context, userID uuid.UUID, input *models.MLModelCreate) (*models.MLModel, error) {
	model := &models.MLModel{
		ID:          uuid.New(),
		UserID:      userID,
		Name:        input.Name,
		Type:        input.Type,
		Version:     input.Version,
		Description: input.Description,
		Provider:    input.Provider,
		Endpoint:    input.Endpoint,
		APIKey:      input.APIKey,
		Config:      input.Config,
		Metrics:     make(map[string]interface{}),
		IsDefault:   input.IsDefault,
		Status:      "inactive", // Default status
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// If this is set as default, unset other defaults of the same type
	if input.IsDefault {
		if err := s.unsetDefaultsForType(ctx, userID, input.Type); err != nil {
			return nil, fmt.Errorf("failed to unset other defaults: %w", err)
		}
	}

	configJSON, err := json.Marshal(model.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal config: %w", err)
	}

	metricsJSON, err := json.Marshal(model.Metrics)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal metrics: %w", err)
	}

	query := `
		INSERT INTO ml_models (
			id, user_id, name, type, version, description, status, provider, 
			endpoint, api_key, config, metrics, is_default, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
		)
	`

	_, err = s.db.ExecContext(ctx, query,
		model.ID, model.UserID, model.Name, model.Type, model.Version, model.Description,
		model.Status, model.Provider, model.Endpoint, model.APIKey, configJSON, metricsJSON,
		model.IsDefault, model.CreatedAt, model.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create model: %w", err)
	}

	return model, nil
}

// GetByID retrieves a model by ID
func (s *MLModelStore) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*models.MLModel, error) {
	query := `
		SELECT id, user_id, name, type, version, description, status, provider, 
		       endpoint, config, metrics, is_default, created_at, updated_at
		FROM ml_models
		WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL
	`

	var model models.MLModel
	var configJSON, metricsJSON []byte

	err := s.db.QueryRowContext(ctx, query, id, userID).Scan(
		&model.ID, &model.UserID, &model.Name, &model.Type, &model.Version, &model.Description,
		&model.Status, &model.Provider, &model.Endpoint, &configJSON, &metricsJSON,
		&model.IsDefault, &model.CreatedAt, &model.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("model not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get model: %w", err)
	}

	if err := json.Unmarshal(configJSON, &model.Config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := json.Unmarshal(metricsJSON, &model.Metrics); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metrics: %w", err)
	}

	return &model, nil
}

// List retrieves all models for a user with pagination
func (s *MLModelStore) List(ctx context.Context, userID uuid.UUID, page, pageSize int) (*models.MLModelListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM ml_models WHERE user_id = $1 AND deleted_at IS NULL`
	if err := s.db.QueryRowContext(ctx, countQuery, userID).Scan(&total); err != nil {
		return nil, fmt.Errorf("failed to count models: %w", err)
	}

	// Get models
	query := `
		SELECT id, user_id, name, type, version, description, status, provider, 
		       endpoint, config, metrics, is_default, created_at, updated_at
		FROM ml_models
		WHERE user_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := s.db.QueryContext(ctx, query, userID, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list models: %w", err)
	}
	defer rows.Close()

	var modelsList []models.MLModel
	for rows.Next() {
		var model models.MLModel
		var configJSON, metricsJSON []byte

		err := rows.Scan(
			&model.ID, &model.UserID, &model.Name, &model.Type, &model.Version, &model.Description,
			&model.Status, &model.Provider, &model.Endpoint, &configJSON, &metricsJSON,
			&model.IsDefault, &model.CreatedAt, &model.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan model: %w", err)
		}

		if err := json.Unmarshal(configJSON, &model.Config); err != nil {
			return nil, fmt.Errorf("failed to unmarshal config: %w", err)
		}

		if err := json.Unmarshal(metricsJSON, &model.Metrics); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metrics: %w", err)
		}

		modelsList = append(modelsList, model)
	}

	totalPages := (total + pageSize - 1) / pageSize

	return &models.MLModelListResponse{
		Models:     modelsList,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// Update updates a model
func (s *MLModelStore) Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, input *models.MLModelUpdate) (*models.MLModel, error) {
	// Get existing model
	existing, err := s.GetByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	// Update fields
	if input.Name != nil {
		existing.Name = *input.Name
	}
	if input.Version != nil {
		existing.Version = *input.Version
	}
	if input.Description != nil {
		existing.Description = *input.Description
	}
	if input.Status != nil {
		existing.Status = *input.Status
	}
	if input.Provider != nil {
		existing.Provider = *input.Provider
	}
	if input.Endpoint != nil {
		existing.Endpoint = input.Endpoint
	}
	if input.APIKey != nil {
		existing.APIKey = input.APIKey
	}
	if input.Config != nil {
		existing.Config = *input.Config
	}
	if input.Metrics != nil {
		existing.Metrics = *input.Metrics
	}
	if input.IsDefault != nil && *input.IsDefault {
		// Unset other defaults of the same type
		if err := s.unsetDefaultsForType(ctx, userID, existing.Type); err != nil {
			return nil, fmt.Errorf("failed to unset other defaults: %w", err)
		}
		existing.IsDefault = true
	} else if input.IsDefault != nil {
		existing.IsDefault = *input.IsDefault
	}

	existing.UpdatedAt = time.Now()

	configJSON, err := json.Marshal(existing.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal config: %w", err)
	}

	metricsJSON, err := json.Marshal(existing.Metrics)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal metrics: %w", err)
	}

	query := `
		UPDATE ml_models
		SET name = $1, version = $2, description = $3, status = $4, provider = $5,
		    endpoint = $6, api_key = $7, config = $8, metrics = $9, is_default = $10, updated_at = $11
		WHERE id = $12 AND user_id = $13 AND deleted_at IS NULL
	`

	_, err = s.db.ExecContext(ctx, query,
		existing.Name, existing.Version, existing.Description, existing.Status, existing.Provider,
		existing.Endpoint, existing.APIKey, configJSON, metricsJSON, existing.IsDefault,
		existing.UpdatedAt, id, userID,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update model: %w", err)
	}

	return existing, nil
}

// Delete soft deletes a model
func (s *MLModelStore) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	query := `
		UPDATE ml_models
		SET deleted_at = $1
		WHERE id = $2 AND user_id = $3 AND deleted_at IS NULL
	`

	result, err := s.db.ExecContext(ctx, query, time.Now(), id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete model: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("model not found")
	}

	return nil
}

// GetDefaultForType retrieves the default model for a specific type
func (s *MLModelStore) GetDefaultForType(ctx context.Context, userID uuid.UUID, modelType string) (*models.MLModel, error) {
	query := `
		SELECT id, user_id, name, type, version, description, status, provider, 
		       endpoint, config, metrics, is_default, created_at, updated_at
		FROM ml_models
		WHERE user_id = $1 AND type = $2 AND is_default = true AND deleted_at IS NULL
		LIMIT 1
	`

	var model models.MLModel
	var configJSON, metricsJSON []byte

	err := s.db.QueryRowContext(ctx, query, userID, modelType).Scan(
		&model.ID, &model.UserID, &model.Name, &model.Type, &model.Version, &model.Description,
		&model.Status, &model.Provider, &model.Endpoint, &configJSON, &metricsJSON,
		&model.IsDefault, &model.CreatedAt, &model.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no default model found for type: %s", modelType)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get default model: %w", err)
	}

	if err := json.Unmarshal(configJSON, &model.Config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := json.Unmarshal(metricsJSON, &model.Metrics); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metrics: %w", err)
	}

	return &model, nil
}

// unsetDefaultsForType unsets all default flags for a specific model type
func (s *MLModelStore) unsetDefaultsForType(ctx context.Context, userID uuid.UUID, modelType string) error {
	query := `
		UPDATE ml_models
		SET is_default = false, updated_at = $1
		WHERE user_id = $2 AND type = $3 AND is_default = true AND deleted_at IS NULL
	`

	_, err := s.db.ExecContext(ctx, query, time.Now(), userID, modelType)
	if err != nil {
		return fmt.Errorf("failed to unset defaults: %w", err)
	}

	return nil
}

// SetAsDefault sets a model as the default for its type
func (s *MLModelStore) SetAsDefault(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	// First get the model to know its type
	model, err := s.GetByID(ctx, id, userID)
	if err != nil {
		return err
	}

	// Unset other defaults of the same type
	if err := s.unsetDefaultsForType(ctx, userID, model.Type); err != nil {
		return err
	}

	// Set this model as default
	query := `
		UPDATE ml_models
		SET is_default = true, updated_at = $1
		WHERE id = $2 AND user_id = $3 AND deleted_at IS NULL
	`

	_, err = s.db.ExecContext(ctx, query, time.Now(), id, userID)
	if err != nil {
		return fmt.Errorf("failed to set as default: %w", err)
	}

	return nil
}

// DB returns the database connection
func (s *MLModelStore) DB() *sql.DB {
	return s.db
}
