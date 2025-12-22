package trpc

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/baskint/bidding-analysis/internal/models"
	"github.com/google/uuid"
)

// Request types for campaign operations
type CreateCampaignRequest struct {
	Name        string   `json:"name"`
	Budget      *float64 `json:"budget,omitempty"`
	DailyBudget *float64 `json:"daily_budget,omitempty"`
	TargetCPA   *float64 `json:"target_cpa,omitempty"`
}

type GetCampaignRequest struct {
	ID string `json:"id"`
}

type UpdateCampaignRequest struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Status      string   `json:"status,omitempty"`
	Budget      *float64 `json:"budget,omitempty"`
	DailyBudget *float64 `json:"daily_budget,omitempty"`
	TargetCPA   *float64 `json:"target_cpa,omitempty"`
}

type DeleteCampaignRequest struct {
	ID string `json:"id"`
}

type PauseCampaignRequest struct {
	ID string `json:"id"`
}

type ActivateCampaignRequest struct {
	ID string `json:"id"`
}

type GetDailyMetricsRequest struct {
	ID string `json:"id"`
}

// ============================================================================
// CAMPAIGN HANDLERS
// ============================================================================

// listCampaigns retrieves all campaigns for the user
func (h *Handler) listCampaigns(ctx context.Context, userID uuid.UUID, req interface{}) (interface{}, error) {
	campaigns, err := h.campaignStore.GetUserCampaigns(ctx, userID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve campaigns: %w", err)
	}

	return campaigns, nil
}

// createCampaign creates a new campaign
func (h *Handler) createCampaign(ctx context.Context, userID uuid.UUID, req interface{}) (interface{}, error) {
	params := req.(*CreateCampaignRequest)

	campaign := &models.Campaign{
		Name:        params.Name,
		UserID:      userID,
		Status:      "active",
		Budget:      params.Budget,
		DailyBudget: params.DailyBudget,
		TargetCPA:   params.TargetCPA,
	}

	if err := h.campaignStore.CreateCampaign(campaign); err != nil {
		return nil, fmt.Errorf("failed to create campaign: %w", err)
	}

	return campaign, nil
}

// getCampaign retrieves a single campaign with detailed metrics
func (h *Handler) getCampaign(ctx context.Context, userID uuid.UUID, req interface{}) (interface{}, error) {
	var campaignID string

	switch r := req.(type) {
	case *http.Request:
		campaignID = r.URL.Query().Get("id")
	case *GetCampaignRequest:
		campaignID = r.ID
	default:
		return nil, fmt.Errorf("invalid request type")
	}

	if campaignID == "" {
		return nil, fmt.Errorf("campaign ID is required")
	}

	id, err := uuid.Parse(campaignID)
	if err != nil {
		return nil, fmt.Errorf("invalid campaign ID format")
	}

	campaign, err := h.campaignStore.GetCampaignWithMetrics(ctx, id, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve campaign: %w", err)
	}

	return campaign, nil
}

// updateCampaign updates an existing campaign
func (h *Handler) updateCampaign(ctx context.Context, userID uuid.UUID, req interface{}) (interface{}, error) {
	params := req.(*UpdateCampaignRequest)

	if params.ID == "" {
		return nil, fmt.Errorf("campaign ID is required")
	}

	id, err := uuid.Parse(params.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid campaign ID format")
	}

	// Validate inputs
	if params.Name != "" && (len(params.Name) < 3 || len(params.Name) > 255) {
		return nil, fmt.Errorf("campaign name must be between 3 and 255 characters")
	}

	if params.Status != "" && params.Status != "active" && params.Status != "paused" && params.Status != "archived" {
		return nil, fmt.Errorf("invalid status. Must be 'active', 'paused', or 'archived'")
	}

	if params.Budget != nil && *params.Budget < 0 {
		return nil, fmt.Errorf("budget must be positive")
	}

	if params.DailyBudget != nil && *params.DailyBudget < 0 {
		return nil, fmt.Errorf("daily budget must be positive")
	}

	if params.Budget != nil && params.DailyBudget != nil && *params.DailyBudget > *params.Budget {
		return nil, fmt.Errorf("daily budget cannot exceed total budget")
	}

	// Get existing campaign first
	existing, err := h.campaignStore.GetCampaign(id)
	if err != nil {
		return nil, fmt.Errorf("campaign not found")
	}

	// Verify ownership
	if existing.UserID != userID {
		return nil, fmt.Errorf("unauthorized")
	}

	// Update fields
	if params.Name != "" {
		existing.Name = params.Name
	}
	if params.Status != "" {
		existing.Status = params.Status
	}
	if params.Budget != nil {
		existing.Budget = params.Budget
	}
	if params.DailyBudget != nil {
		existing.DailyBudget = params.DailyBudget
	}
	if params.TargetCPA != nil {
		existing.TargetCPA = params.TargetCPA
	}

	if err := h.campaignStore.UpdateCampaign(existing); err != nil {
		return nil, fmt.Errorf("failed to update campaign: %w", err)
	}

	return existing, nil
}

// deleteCampaign soft deletes a campaign
func (h *Handler) deleteCampaign(ctx context.Context, userID uuid.UUID, req interface{}) (interface{}, error) {
	params := req.(*DeleteCampaignRequest)

	if params.ID == "" {
		return nil, fmt.Errorf("campaign ID is required")
	}

	id, err := uuid.Parse(params.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid campaign ID format")
	}

	if err := h.campaignStore.DeleteCampaign(id, userID); err != nil {
		return nil, fmt.Errorf("failed to delete campaign: %w", err)
	}

	return map[string]interface{}{
		"success": true,
		"message": "Campaign archived successfully",
	}, nil
}

// pauseCampaign pauses an active campaign
func (h *Handler) pauseCampaign(ctx context.Context, userID uuid.UUID, req interface{}) (interface{}, error) {
	params := req.(*PauseCampaignRequest)

	if params.ID == "" {
		return nil, fmt.Errorf("campaign ID is required")
	}

	id, err := uuid.Parse(params.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid campaign ID format")
	}

	if err := h.campaignStore.PauseCampaign(id, userID); err != nil {
		return nil, fmt.Errorf("failed to pause campaign: %w", err)
	}

	// Return updated campaign
	campaign, err := h.campaignStore.GetCampaign(id)
	if err != nil {
		return nil, fmt.Errorf("campaign paused but failed to retrieve")
	}

	return campaign, nil
}

// activateCampaign activates a paused campaign
func (h *Handler) activateCampaign(ctx context.Context, userID uuid.UUID, req interface{}) (interface{}, error) {
	params := req.(*ActivateCampaignRequest)

	if params.ID == "" {
		return nil, fmt.Errorf("campaign ID is required")
	}

	id, err := uuid.Parse(params.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid campaign ID format")
	}

	if err := h.campaignStore.ActivateCampaign(id, userID); err != nil {
		return nil, fmt.Errorf("failed to activate campaign: %w", err)
	}

	// Return updated campaign
	campaign, err := h.campaignStore.GetCampaign(id)
	if err != nil {
		return nil, fmt.Errorf("campaign activated but failed to retrieve")
	}

	return campaign, nil
}

// listCampaignsEnhanced lists campaigns with metrics
func (h *Handler) listCampaignsEnhanced(ctx context.Context, userID uuid.UUID, req interface{}) (interface{}, error) {
	// Parse pagination parameters (optional)
	limit := 100
	offset := 0

	campaigns, err := h.campaignStore.ListCampaignsWithMetrics(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve campaigns: %w", err)
	}

	return campaigns, nil
}

// getDailyMetrics retrieves daily metrics for a campaign
func (h *Handler) getDailyMetrics(ctx context.Context, userID uuid.UUID, req interface{}) (interface{}, error) {
	var campaignID string

	switch r := req.(type) {
	case *http.Request:
		campaignID = r.URL.Query().Get("id")
	case *GetDailyMetricsRequest:
		campaignID = r.ID
	default:
		return nil, fmt.Errorf("invalid request type")
	}

	if campaignID == "" {
		return nil, fmt.Errorf("campaign ID is required")
	}

	id, err := uuid.Parse(campaignID)
	if err != nil {
		return nil, fmt.Errorf("invalid campaign ID format")
	}

	// Verify ownership
	campaign, err := h.campaignStore.GetCampaign(id)
	if err != nil || campaign.UserID != userID {
		return nil, fmt.Errorf("campaign not found or unauthorized")
	}

	// Get metrics for last 30 days by default
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -30)

	metrics, err := h.campaignStore.GetCampaignDailyMetrics(id, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve metrics: %w", err)
	}

	return metrics, nil
}

// ============================================================================
// DASHBOARD & METRICS HANDLERS (Mock implementations - TODO: implement properly)
// ============================================================================

// getDashboardMetrics returns dashboard overview metrics
func (h *Handler) getDashboardMetrics(ctx context.Context, userID uuid.UUID, req interface{}) (interface{}, error) {
	// Mock dashboard metrics - replace with actual database queries
	metrics := map[string]interface{}{
		"total_campaigns": 8,
		"active_bids":     1247,
		"win_rate":        0.348,
		"avg_bid":         2.34,
		"total_spend":     12543.67,
		"conversions":     89,
		"fraud_alerts":    2,
		"model_accuracy":  0.92,
		"last_updated":    time.Now(),
	}

	return metrics, nil
}

// getCampaignStats returns campaign statistics (mock)
func (h *Handler) getCampaignStats(ctx context.Context, userID uuid.UUID, req interface{}) (interface{}, error) {
	stats := map[string]interface{}{
		"total_bids":  1500,
		"won_bids":    522,
		"win_rate":    0.348,
		"total_spend": 4567.89,
		"conversions": 45,
		"avg_cpa":     101.51,
	}
	return stats, nil
}

// getBidHistory returns recent bid history
func (h *Handler) getBidHistory(ctx context.Context, userID uuid.UUID, req interface{}) (interface{}, error) {
	bids, err := h.bidStore.GetRecentBids(ctx, 20)
	if err != nil {
		return nil, fmt.Errorf("failed to get bid history: %w", err)
	}
	return bids, nil
}

// getFraudAlerts returns fraud alerts (mock)
func (h *Handler) getFraudAlerts(ctx context.Context, userID uuid.UUID, req interface{}) (interface{}, error) {
	// Mock fraud alerts
	alerts := []map[string]interface{}{
		{
			"id":          "alert-1",
			"type":        "suspicious_click_velocity",
			"severity":    "medium",
			"campaign_id": "campaign-123",
			"detected_at": time.Now().Add(-2 * time.Hour),
			"status":      "active",
		},
		{
			"id":          "alert-2",
			"type":        "geographic_anomaly",
			"severity":    "high",
			"campaign_id": "campaign-456",
			"detected_at": time.Now().Add(-1 * time.Hour),
			"status":      "investigating",
		},
	}
	return alerts, nil
}

// getModelAccuracy returns ML model accuracy metrics (mock)
func (h *Handler) getModelAccuracy(ctx context.Context, userID uuid.UUID, req interface{}) (interface{}, error) {
	accuracy := map[string]interface{}{
		"current_accuracy":    0.924,
		"last_week_accuracy":  0.918,
		"trend":               "improving",
		"total_predictions":   15420,
		"correct_predictions": 14248,
		"last_updated":        time.Now(),
	}
	return accuracy, nil
}
