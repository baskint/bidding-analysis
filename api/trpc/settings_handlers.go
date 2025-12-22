package trpc

import (
	"context"
	"fmt"

	"github.com/baskint/bidding-analysis/internal/models"
	"github.com/google/uuid"
)

// Request types
type UpdateUserSettingsRequest struct {
	models.UserSettingsUpdate
}

type GetIntegrationRequest struct {
	ID string `json:"id"`
}

type CreateIntegrationRequest struct {
	models.IntegrationCreate
}

type UpdateIntegrationRequest struct {
	ID string `json:"id"`
	models.IntegrationUpdate
}

type DeleteIntegrationRequest struct {
	ID string `json:"id"`
}

type TestIntegrationRequest struct {
	ID string `json:"id"`
}

// ============================================================================
// REFACTORED HANDLERS
// ============================================================================

// GetUserSettings retrieves user settings
func (h *Handler) GetUserSettings(ctx context.Context, userID uuid.UUID, req interface{}) (interface{}, error) {
	settings, err := h.settingsStore.GetUserSettings(ctx, userID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get settings: %w", err)
	}

	return settings, nil
}

// UpdateUserSettings updates user settings
func (h *Handler) UpdateUserSettings(ctx context.Context, userID uuid.UUID, req interface{}) (interface{}, error) {
	params := req.(*UpdateUserSettingsRequest)

	settings, err := h.settingsStore.UpdateUserSettings(ctx, userID.String(), &params.UserSettingsUpdate)
	if err != nil {
		return nil, fmt.Errorf("failed to update settings: %w", err)
	}

	return settings, nil
}

// RegenerateAPIKey generates a new API key
func (h *Handler) RegenerateAPIKey(ctx context.Context, userID uuid.UUID, req interface{}) (interface{}, error) {
	apiKey, err := h.settingsStore.RegenerateAPIKey(ctx, userID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to regenerate API key: %w", err)
	}

	return map[string]string{"api_key": apiKey}, nil
}

// ListIntegrations retrieves all user integrations
func (h *Handler) ListIntegrations(ctx context.Context, userID uuid.UUID, req interface{}) (interface{}, error) {
	integrations, err := h.settingsStore.ListIntegrations(ctx, userID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to list integrations: %w", err)
	}

	return map[string]interface{}{
		"integrations": integrations,
	}, nil
}

// GetIntegration retrieves a specific integration
// Note: This needs to extract ID from request context or query params
func (h *Handler) GetIntegration(ctx context.Context, userID uuid.UUID, req interface{}) (interface{}, error) {
	// For query param handlers, req will be the *http.Request
	// Extract ID from the request passed in
	r, ok := req.(*GetIntegrationRequest)
	if !ok {
		return nil, fmt.Errorf("invalid request type")
	}

	if r.ID == "" {
		return nil, fmt.Errorf("integration ID is required")
	}

	integration, err := h.settingsStore.GetIntegration(ctx, userID.String(), r.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get integration: %w", err)
	}

	return integration, nil
}

// CreateIntegration creates a new integration
func (h *Handler) CreateIntegration(ctx context.Context, userID uuid.UUID, req interface{}) (interface{}, error) {
	params := req.(*CreateIntegrationRequest)

	integration, err := h.settingsStore.CreateIntegration(ctx, userID.String(), &params.IntegrationCreate)
	if err != nil {
		return nil, fmt.Errorf("failed to create integration: %w", err)
	}

	return integration, nil
}

// UpdateIntegration updates an existing integration
func (h *Handler) UpdateIntegration(ctx context.Context, userID uuid.UUID, req interface{}) (interface{}, error) {
	params := req.(*UpdateIntegrationRequest)

	if params.ID == "" {
		return nil, fmt.Errorf("integration ID is required")
	}

	integration, err := h.settingsStore.UpdateIntegration(ctx, userID.String(), params.ID, &params.IntegrationUpdate)
	if err != nil {
		return nil, fmt.Errorf("failed to update integration: %w", err)
	}

	return integration, nil
}

// DeleteIntegration deletes an integration
func (h *Handler) DeleteIntegration(ctx context.Context, userID uuid.UUID, req interface{}) (interface{}, error) {
	params := req.(*DeleteIntegrationRequest)

	if params.ID == "" {
		return nil, fmt.Errorf("integration ID is required")
	}

	if err := h.settingsStore.DeleteIntegration(ctx, userID.String(), params.ID); err != nil {
		return nil, fmt.Errorf("failed to delete integration: %w", err)
	}

	return map[string]bool{"success": true}, nil
}

// TestIntegration tests an integration connection
func (h *Handler) TestIntegration(ctx context.Context, userID uuid.UUID, req interface{}) (interface{}, error) {
	params := req.(*TestIntegrationRequest)

	if params.ID == "" {
		return nil, fmt.Errorf("integration ID is required")
	}

	integration, err := h.settingsStore.GetIntegration(ctx, userID.String(), params.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get integration: %w", err)
	}

	if err := h.testIntegrationConnection(ctx, integration); err != nil {
		return map[string]string{
			"message": fmt.Sprintf("Connection test failed: %v", err),
		}, nil
	}

	return map[string]string{
		"message": "Integration connection successful",
	}, nil
}

// GetBillingInfo retrieves billing information
func (h *Handler) GetBillingInfo(ctx context.Context, userID uuid.UUID, req interface{}) (interface{}, error) {
	billing, err := h.settingsStore.GetBillingInfo(ctx, userID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get billing info: %w", err)
	}

	return billing, nil
}

// testIntegrationConnection tests the connection to an integration
func (h *Handler) testIntegrationConnection(ctx context.Context, integration *models.Integration) error {
	// Placeholder for actual integration testing logic
	// This would vary based on the integration provider
	switch integration.Provider {
	case "google_ads":
		// Test Google Ads connection
		return nil
	case "facebook_ads":
		// Test Facebook Ads connection
		return nil
	case "slack":
		// Test Slack connection
		return nil
	default:
		return fmt.Errorf("unsupported integration provider: %s", integration.Provider)
	}
}
