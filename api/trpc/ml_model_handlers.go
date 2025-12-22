package trpc

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/baskint/bidding-analysis/internal/models"
	"github.com/google/uuid"
)

// Request types
type ListMLModelsRequest struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

type GetMLModelRequest struct {
	ID string `json:"id"`
}

type CreateMLModelRequest struct {
	models.MLModelCreate
}

type UpdateMLModelRequest struct {
	ID string `json:"id"`
	models.MLModelUpdate
}

type DeleteMLModelRequest struct {
	ID string `json:"id"`
}

type SetDefaultMLModelRequest struct {
	ID string `json:"id"`
}

type GetDefaultMLModelRequest struct {
	Type string `json:"type"`
}

// ============================================================================
// REFACTORED HANDLERS
// ============================================================================

// listMLModels retrieves all ML models for the current user with pagination
// This handler needs special handling for query params
func (h *Handler) listMLModels(ctx context.Context, userID uuid.UUID, req interface{}) (interface{}, error) {
	// req will be *http.Request for query param handlers
	// We'll need to handle this specially
	page := 1
	pageSize := 20

	// Try to extract from request if it's the right type
	if r, ok := req.(*http.Request); ok {
		if pageStr := r.URL.Query().Get("page"); pageStr != "" {
			if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
				page = p
			}
		}
		if pageSizeStr := r.URL.Query().Get("pageSize"); pageSizeStr != "" {
			if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 100 {
				pageSize = ps
			}
		}
	}

	response, err := h.mlModelStore.List(ctx, userID, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to list ML models: %w", err)
	}

	return response, nil
}

// getMLModel retrieves a specific ML model by ID
func (h *Handler) getMLModel(ctx context.Context, userID uuid.UUID, req interface{}) (interface{}, error) {
	var modelID string

	// Handle both query param and body request
	switch r := req.(type) {
	case *http.Request:
		modelID = r.URL.Query().Get("id")
	case *GetMLModelRequest:
		modelID = r.ID
	default:
		return nil, fmt.Errorf("invalid request type")
	}

	if modelID == "" {
		return nil, fmt.Errorf("model ID is required")
	}

	modelUUID, err := uuid.Parse(modelID)
	if err != nil {
		return nil, fmt.Errorf("invalid model ID format")
	}

	model, err := h.mlModelStore.GetByID(ctx, modelUUID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get ML model: %w", err)
	}

	return model, nil
}

// createMLModel creates a new ML model
func (h *Handler) createMLModel(ctx context.Context, userID uuid.UUID, req interface{}) (interface{}, error) {
	params := req.(*CreateMLModelRequest)

	// Validate required fields
	if params.Name == "" {
		return nil, fmt.Errorf("model name is required")
	}
	if params.Type == "" {
		return nil, fmt.Errorf("model type is required")
	}
	if params.Provider == "" {
		return nil, fmt.Errorf("model provider is required")
	}

	model, err := h.mlModelStore.Create(ctx, userID, &params.MLModelCreate)
	if err != nil {
		return nil, fmt.Errorf("failed to create ML model: %w", err)
	}

	return model, nil
}

// updateMLModel updates an existing ML model
func (h *Handler) updateMLModel(ctx context.Context, userID uuid.UUID, req interface{}) (interface{}, error) {
	params := req.(*UpdateMLModelRequest)

	if params.ID == "" {
		return nil, fmt.Errorf("model ID is required")
	}

	modelUUID, err := uuid.Parse(params.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid model ID format")
	}

	model, err := h.mlModelStore.Update(ctx, userID, modelUUID, &params.MLModelUpdate)
	if err != nil {
		return nil, fmt.Errorf("failed to update ML model: %w", err)
	}

	return model, nil
}

// deleteMLModel deletes an ML model
func (h *Handler) deleteMLModel(ctx context.Context, userID uuid.UUID, req interface{}) (interface{}, error) {
	params := req.(*DeleteMLModelRequest)

	if params.ID == "" {
		return nil, fmt.Errorf("model ID is required")
	}

	modelUUID, err := uuid.Parse(params.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid model ID format")
	}

	if err := h.mlModelStore.Delete(ctx, userID, modelUUID); err != nil {
		return nil, fmt.Errorf("failed to delete ML model: %w", err)
	}

	return map[string]interface{}{
		"success": true,
		"message": "ML model deleted successfully",
	}, nil
}

// setDefaultMLModel sets an ML model as the default for its type
func (h *Handler) setDefaultMLModel(ctx context.Context, userID uuid.UUID, req interface{}) (interface{}, error) {
	params := req.(*SetDefaultMLModelRequest)

	if params.ID == "" {
		return nil, fmt.Errorf("model ID is required")
	}

	modelUUID, err := uuid.Parse(params.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid model ID format")
	}

	if err := h.mlModelStore.SetAsDefault(ctx, modelUUID, userID); err != nil {
		return nil, fmt.Errorf("failed to set default ML model: %w", err)
	}

	// Get the updated model to return
	model, err := h.mlModelStore.GetByID(ctx, modelUUID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated ML model: %w", err)
	}

	return model, nil
}

// getDefaultMLModel retrieves the default ML model for a specific type
func (h *Handler) getDefaultMLModel(ctx context.Context, userID uuid.UUID, req interface{}) (interface{}, error) {
	var modelType string

	// Handle both query param and body request
	switch r := req.(type) {
	case *http.Request:
		modelType = r.URL.Query().Get("type")
	case *GetDefaultMLModelRequest:
		modelType = r.Type
	default:
		return nil, fmt.Errorf("invalid request type")
	}

	if modelType == "" {
		return nil, fmt.Errorf("model type is required")
	}

	model, err := h.mlModelStore.GetDefaultForType(ctx, userID, modelType)
	if err != nil {
		return nil, fmt.Errorf("failed to get default ML model: %w", err)
	}

	return model, nil
}
