package trpc

import (
	"context"
	"fmt"
	"time"

	"github.com/baskint/bidding-analysis/internal/fraud"
	"github.com/baskint/bidding-analysis/internal/models"
	"github.com/google/uuid"
)

// Request types
type FraudOverviewRequest struct {
	Days int `json:"days"`
}

type FraudAlertsRequest struct {
	Status      string `json:"status"`
	MinSeverity int    `json:"min_severity"`
	AlertType   string `json:"alert_type"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	Limit       int    `json:"limit"`
}

type UpdateFraudAlertRequest struct {
	AlertID string `json:"alert_id"`
	Status  string `json:"status"`
	Notes   string `json:"notes"`
}

type FraudTrendsRequest struct {
	Days int `json:"days"`
}

type DeviceFraudRequest struct {
	Days int `json:"days"`
}

type GeoFraudRequest struct {
	Days int `json:"days"`
}

type CreateFraudAlertRequest struct {
	CampaignID      string   `json:"campaign_id"`
	AlertType       string   `json:"alert_type"`
	Severity        int      `json:"severity"`
	Description     string   `json:"description"`
	AffectedUserIDs []string `json:"affected_user_ids"`
}

// ============================================================================
// REFACTORED HANDLERS
// ============================================================================

// getFraudOverview returns high-level fraud metrics and statistics
func (h *Handler) getFraudOverview(ctx context.Context, userID uuid.UUID, req interface{}) (interface{}, error) {
	params := req.(*FraudOverviewRequest)

	days := params.Days
	if days <= 0 || days > 365 {
		days = 30
	}

	fraudDetector := fraud.NewFraudDetector(h.bidStore.DB())
	overview, err := fraudDetector.GetFraudOverview(ctx, userID, days)
	if err != nil {
		return nil, fmt.Errorf("failed to get fraud overview: %w", err)
	}

	return overview, nil
}

// getRealFraudAlerts returns fraud alerts with filtering
func (h *Handler) getRealFraudAlerts(ctx context.Context, userID uuid.UUID, req interface{}) (interface{}, error) {
	params := req.(*FraudAlertsRequest)

	// Build filter
	filter := fraud.FraudAlertFilter{
		Status:      params.Status,
		MinSeverity: params.MinSeverity,
		AlertType:   params.AlertType,
		Limit:       params.Limit,
	}

	// Parse dates
	if params.StartDate != "" {
		if startDate, err := time.Parse("2006-01-02", params.StartDate); err == nil {
			filter.StartDate = startDate
		}
	}
	if params.EndDate != "" {
		if endDate, err := time.Parse("2006-01-02", params.EndDate); err == nil {
			filter.EndDate = endDate
		}
	}

	if filter.Limit <= 0 || filter.Limit > 100 {
		filter.Limit = 50
	}

	fraudDetector := fraud.NewFraudDetector(h.bidStore.DB())
	alerts, err := fraudDetector.GetFraudAlerts(ctx, userID, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get fraud alerts: %w", err)
	}

	// Convert to response format (frontend expects this structure)
	type alertResponse struct {
		ID              string     `json:"id"`
		CampaignID      string     `json:"campaign_id"`
		AlertType       string     `json:"alert_type"`
		Severity        int        `json:"severity"`
		Description     string     `json:"description"`
		AffectedUserIDs []string   `json:"affected_user_ids"`
		DetectedAt      time.Time  `json:"detected_at"`
		ResolvedAt      *time.Time `json:"resolved_at,omitempty"`
		Status          string     `json:"status"`
	}

	response := make([]alertResponse, 0, len(alerts))
	for _, alert := range alerts {
		response = append(response, alertResponse{
			ID:              alert.ID.String(),
			CampaignID:      alert.CampaignID.String(),
			AlertType:       alert.AlertType,
			Severity:        alert.Severity,
			Description:     alert.Description,
			AffectedUserIDs: alert.AffectedUserIDs,
			DetectedAt:      alert.DetectedAt,
			ResolvedAt:      alert.ResolvedAt,
			Status:          alert.Status,
		})
	}

	return response, nil
}

// updateFraudAlert updates the status of a fraud alert
func (h *Handler) updateFraudAlert(ctx context.Context, userID uuid.UUID, req interface{}) (interface{}, error) {
	params := req.(*UpdateFraudAlertRequest)

	if params.AlertID == "" {
		return nil, fmt.Errorf("alert_id is required")
	}

	alertUUID, err := uuid.Parse(params.AlertID)
	if err != nil {
		return nil, fmt.Errorf("invalid alert ID format")
	}

	validStatuses := map[string]bool{
		"active":         true,
		"investigating":  true,
		"resolved":       true,
		"false_positive": true,
	}

	if !validStatuses[params.Status] {
		return nil, fmt.Errorf("invalid status: must be one of active, investigating, resolved, false_positive")
	}

	// Create fraud detector
	fraudDetector := fraud.NewFraudDetector(h.bidStore.DB())

	// Update the alert (method is UpdateAlertStatus, not UpdateFraudAlert)
	err = fraudDetector.UpdateAlertStatus(ctx, alertUUID, userID, params.Status, params.Notes)
	if err != nil {
		return nil, fmt.Errorf("failed to update fraud alert: %w", err)
	}

	return map[string]interface{}{
		"success": true,
		"message": "Fraud alert updated successfully",
	}, nil
}

// getFraudTrends returns fraud trends over time
func (h *Handler) getFraudTrends(ctx context.Context, userID uuid.UUID, req interface{}) (interface{}, error) {
	params := req.(*FraudTrendsRequest)

	days := params.Days
	if days <= 0 || days > 365 {
		days = 30
	}

	fraudDetector := fraud.NewFraudDetector(h.bidStore.DB())
	trends, err := fraudDetector.GetFraudTrends(ctx, userID, days)
	if err != nil {
		return nil, fmt.Errorf("failed to get fraud trends: %w", err)
	}

	return trends, nil
}

// getDeviceFraudAnalysis returns device-specific fraud analysis
// NOTE: This method is not implemented in the fraud detector yet
// Returning mock data for now
func (h *Handler) getDeviceFraudAnalysis(ctx context.Context, userID uuid.UUID, req interface{}) (interface{}, error) {
	params := req.(*DeviceFraudRequest)

	days := params.Days
	if days <= 0 || days > 365 {
		days = 30
	}

	// TODO: Implement device fraud analysis in fraud detector
	// For now, return empty array
	return []interface{}{}, nil
}

// getGeoFraudAnalysis returns geographic fraud analysis
// NOTE: This method is not implemented in the fraud detector yet
// Returning mock data for now
func (h *Handler) getGeoFraudAnalysis(ctx context.Context, userID uuid.UUID, req interface{}) (interface{}, error) {
	params := req.(*GeoFraudRequest)

	days := params.Days
	if days <= 0 || days > 365 {
		days = 30
	}

	// TODO: Implement geo fraud analysis in fraud detector
	// For now, return empty array
	return []interface{}{}, nil
}

// createFraudAlert creates a new fraud alert
func (h *Handler) createFraudAlert(ctx context.Context, userID uuid.UUID, req interface{}) (interface{}, error) {
	params := req.(*CreateFraudAlertRequest)

	// Validate required fields
	if params.CampaignID == "" {
		return nil, fmt.Errorf("campaign_id is required")
	}
	if params.AlertType == "" {
		return nil, fmt.Errorf("alert_type is required")
	}
	if params.Severity < 1 || params.Severity > 10 {
		return nil, fmt.Errorf("severity must be between 1 and 10")
	}
	if params.Description == "" {
		return nil, fmt.Errorf("description is required")
	}

	campaignUUID, err := uuid.Parse(params.CampaignID)
	if err != nil {
		return nil, fmt.Errorf("invalid campaign_id format")
	}

	// Affected user IDs are already strings - use them directly
	affectedUserIDs := params.AffectedUserIDs
	if affectedUserIDs == nil {
		affectedUserIDs = []string{}
	}

	// Create the alert
	alert := &models.FraudAlert{
		ID:              uuid.New(),
		CampaignID:      campaignUUID,
		AlertType:       params.AlertType,
		Severity:        params.Severity,
		Description:     params.Description,
		AffectedUserIDs: affectedUserIDs,
		DetectedAt:      time.Now(),
		Status:          "active",
	}

	fraudDetector := fraud.NewFraudDetector(h.bidStore.DB())
	if err := fraudDetector.CreateFraudAlert(ctx, alert); err != nil {
		return nil, fmt.Errorf("failed to create fraud alert: %w", err)
	}

	return map[string]interface{}{
		"success":  true,
		"alert_id": alert.ID.String(),
		"message":  "Fraud alert created successfully",
	}, nil
}
