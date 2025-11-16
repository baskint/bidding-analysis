package trpc

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/baskint/bidding-analysis/internal/models"
)

// GetUserSettings retrieves user settings
func (h *Handler) GetUserSettings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := GetUserIDFromContext(ctx)

	settings, err := h.settingsStore.GetUserSettings(ctx, userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get settings", err)
		return
	}

	writeSuccess(w, settings)
}

// UpdateUserSettings updates user settings
func (h *Handler) UpdateUserSettings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := GetUserIDFromContext(ctx)

	var update models.UserSettingsUpdate
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	settings, err := h.settingsStore.UpdateUserSettings(ctx, userID, &update)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to update settings", err)
		return
	}

	writeSuccess(w, settings)
}

// RegenerateAPIKey generates a new API key
func (h *Handler) RegenerateAPIKey(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := GetUserIDFromContext(ctx)

	apiKey, err := h.settingsStore.RegenerateAPIKey(ctx, userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to regenerate API key", err)
		return
	}

	writeSuccess(w, map[string]string{"api_key": apiKey})
}

// ListIntegrations retrieves all user integrations
func (h *Handler) ListIntegrations(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := GetUserIDFromContext(ctx)

	integrations, err := h.settingsStore.ListIntegrations(ctx, userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to list integrations", err)
		return
	}

	writeSuccess(w, map[string]interface{}{
		"integrations": integrations,
	})
}

// GetIntegration retrieves a specific integration
func (h *Handler) GetIntegration(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := GetUserIDFromContext(ctx)
	id := r.URL.Query().Get("id")

	if id == "" {
		writeError(w, http.StatusBadRequest, "Integration ID is required", nil)
		return
	}

	integration, err := h.settingsStore.GetIntegration(ctx, id, userID)
	if err != nil {
		writeError(w, http.StatusNotFound, "Integration not found", err)
		return
	}

	writeSuccess(w, integration)
}

// CreateIntegration creates a new integration
func (h *Handler) CreateIntegration(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := GetUserIDFromContext(ctx)

	var create models.IntegrationCreate
	if err := json.NewDecoder(r.Body).Decode(&create); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Validate provider
	validProviders := map[string]bool{
		"google_ads":       true,
		"facebook_ads":     true,
		"microsoft_ads":    true,
		"slack":            true,
		"webhook":          true,
		"google_analytics": true,
		"segment":          true,
		"stripe":           true,
		"sendgrid":         true,
		"aws_ses":          true,
	}

	if !validProviders[create.Provider] {
		writeError(w, http.StatusBadRequest, "Invalid provider", nil)
		return
	}

	integration, err := h.settingsStore.CreateIntegration(ctx, userID, &create)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create integration", err)
		return
	}

	writeSuccess(w, integration)
}

// UpdateIntegration updates an existing integration
func (h *Handler) UpdateIntegration(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := GetUserIDFromContext(ctx)
	id := r.URL.Query().Get("id")

	if id == "" {
		writeError(w, http.StatusBadRequest, "Integration ID is required", nil)
		return
	}

	var update models.IntegrationUpdate
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	integration, err := h.settingsStore.UpdateIntegration(ctx, id, userID, &update)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to update integration", err)
		return
	}

	writeSuccess(w, integration)
}

// DeleteIntegration deletes an integration
func (h *Handler) DeleteIntegration(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := GetUserIDFromContext(ctx)
	id := r.URL.Query().Get("id")

	if id == "" {
		writeError(w, http.StatusBadRequest, "Integration ID is required", nil)
		return
	}

	err := h.settingsStore.DeleteIntegration(ctx, id, userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to delete integration", err)
		return
	}

	writeSuccess(w, map[string]string{"message": "Integration deleted successfully"})
}

// TestIntegration tests an integration connection
func (h *Handler) TestIntegration(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := GetUserIDFromContext(ctx)
	id := r.URL.Query().Get("id")

	if id == "" {
		writeError(w, http.StatusBadRequest, "Integration ID is required", nil)
		return
	}

	integration, err := h.settingsStore.GetIntegration(ctx, id, userID)
	if err != nil {
		writeError(w, http.StatusNotFound, "Integration not found", err)
		return
	}

	// Test the integration based on provider
	err = h.testIntegrationConnection(ctx, integration)
	if err != nil {
		// FIX: Store error message in a variable first
		errorMsg := err.Error()
		h.settingsStore.UpdateIntegrationStatus(ctx, id, userID, "error", &errorMsg)
		writeError(w, http.StatusBadGateway, "Integration test failed", err)
		return
	}

	h.settingsStore.UpdateIntegrationStatus(ctx, id, userID, "active", nil)
	writeSuccess(w, map[string]string{"message": "Integration test successful"})
}

// GetBillingInfo retrieves billing information
func (h *Handler) GetBillingInfo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := GetUserIDFromContext(ctx)

	billing, err := h.settingsStore.GetBillingInfo(ctx, userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get billing info", err)
		return
	}

	writeSuccess(w, billing)
}

// Helper function to test integration connections
func (h *Handler) testIntegrationConnection(ctx context.Context, integration *models.Integration) error {
	// Implement provider-specific testing logic here
	// For now, just check if credentials exist
	switch integration.Provider {
	case "google_ads", "facebook_ads", "microsoft_ads":
		if integration.AccessToken == nil || *integration.AccessToken == "" {
			return fmt.Errorf("access token is required")
		}
	case "slack", "webhook":
		if integration.WebhookURL == nil || *integration.WebhookURL == "" {
			return fmt.Errorf("webhook URL is required")
		}
	case "sendgrid", "aws_ses":
		if integration.APIKey == nil || *integration.APIKey == "" {
			return fmt.Errorf("API key is required")
		}
	case "stripe":
		if integration.APIKey == nil || integration.APISecret == nil {
			return fmt.Errorf("API key and secret are required")
		}
	}

	// In production, you would make actual API calls to test the connection
	return nil
}
