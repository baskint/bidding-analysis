package trpc

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/baskint/bidding-analysis/internal/fraud"
	"github.com/baskint/bidding-analysis/internal/models"
	"github.com/google/uuid"
)

// getFraudOverview returns high-level fraud metrics and statistics
func (h *Handler) getFraudOverview(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	if userID == "" {
		h.writeErrorResponse(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		h.writeErrorResponse(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Days int `json:"days"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Default to 30 days if not specified
		req.Days = 30
	}

	if req.Days <= 0 || req.Days > 365 {
		req.Days = 30
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	fraudDetector := fraud.NewFraudDetector(h.bidStore.DB())
	overview, err := fraudDetector.GetFraudOverview(ctx, userUUID, req.Days)
	if err != nil {
		log.Printf("Failed to get fraud overview: %v", err)
		h.writeErrorResponse(w, "Failed to retrieve fraud overview", http.StatusInternalServerError)
		return
	}

	h.writeTRPCResponse(w, overview)
}

// Updated getReadFraudAlerts handler with real data
func (h *Handler) getRealFraudAlerts(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	if userID == "" {
		h.writeErrorResponse(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		h.writeErrorResponse(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Parse request parameters
	var req struct {
		Status      string `json:"status"`
		MinSeverity int    `json:"min_severity"`
		AlertType   string `json:"alert_type"`
		StartDate   string `json:"start_date"`
		EndDate     string `json:"end_date"`
		Limit       int    `json:"limit"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Use defaults if decode fails
		log.Printf("Failed to decode request, using defaults: %v", err)
	}

	// Build filter
	filter := fraud.FraudAlertFilter{
		Status:      req.Status,
		MinSeverity: req.MinSeverity,
		AlertType:   req.AlertType,
		Limit:       req.Limit,
	}

	// Parse dates
	if req.StartDate != "" {
		if startDate, err := time.Parse("2006-01-02", req.StartDate); err == nil {
			filter.StartDate = startDate
		}
	}
	if req.EndDate != "" {
		if endDate, err := time.Parse("2006-01-02", req.EndDate); err == nil {
			filter.EndDate = endDate
		}
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	fraudDetector := fraud.NewFraudDetector(h.bidStore.DB())
	alerts, err := fraudDetector.GetFraudAlerts(ctx, userUUID, filter)
	if err != nil {
		log.Printf("Failed to get fraud alerts: %v", err)
		h.writeErrorResponse(w, "Failed to retrieve fraud alerts", http.StatusInternalServerError)
		return
	}

	// Convert to response format
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

	h.writeTRPCResponse(w, response)
}

// updateFraudAlert updates a fraud alert's status
func (h *Handler) updateFraudAlert(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	if userID == "" {
		h.writeErrorResponse(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		h.writeErrorResponse(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req struct {
		AlertID string `json:"alert_id"`
		Status  string `json:"status"`
		Notes   string `json:"notes"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.AlertID == "" {
		h.writeErrorResponse(w, "Alert ID is required", http.StatusBadRequest)
		return
	}

	// Validate status
	validStatuses := map[string]bool{
		"active":         true,
		"investigating":  true,
		"resolved":       true,
		"false_positive": true,
	}

	if !validStatuses[req.Status] {
		h.writeErrorResponse(w, "Invalid status", http.StatusBadRequest)
		return
	}

	alertID, err := uuid.Parse(req.AlertID)
	if err != nil {
		h.writeErrorResponse(w, "Invalid alert ID format", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	fraudDetector := fraud.NewFraudDetector(h.bidStore.DB())
	err = fraudDetector.UpdateAlertStatus(ctx, alertID, userUUID, req.Status, req.Notes)
	if err != nil {
		log.Printf("Failed to update fraud alert: %v", err)
		if err.Error() == "unauthorized" {
			h.writeErrorResponse(w, "Unauthorized", http.StatusForbidden)
		} else {
			h.writeErrorResponse(w, "Failed to update alert", http.StatusInternalServerError)
		}
		return
	}

	h.writeTRPCResponse(w, map[string]interface{}{
		"success": true,
		"message": "Alert updated successfully",
	})
}

// getFraudTrends returns fraud metrics over time
func (h *Handler) getFraudTrends(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	if userID == "" {
		h.writeErrorResponse(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		h.writeErrorResponse(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Days int `json:"days"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		req.Days = 30
	}

	if req.Days <= 0 || req.Days > 365 {
		req.Days = 30
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	fraudDetector := fraud.NewFraudDetector(h.bidStore.DB())
	trends, err := fraudDetector.GetFraudTrends(ctx, userUUID, req.Days)
	if err != nil {
		log.Printf("Failed to get fraud trends: %v", err)
		h.writeErrorResponse(w, "Failed to retrieve fraud trends", http.StatusInternalServerError)
		return
	}

	h.writeTRPCResponse(w, trends)
}

// getDeviceFraudAnalysis returns device-specific fraud metrics
func (h *Handler) getDeviceFraudAnalysis(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	if userID == "" {
		h.writeErrorResponse(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		h.writeErrorResponse(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Days int `json:"days"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		req.Days = 30
	}

	startDate := time.Now().AddDate(0, 0, -req.Days)

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	query := `
		SELECT 
			be.device_type,
			be.browser,
			be.os,
			COUNT(*) as total_bids,
			SUM(CASE WHEN p.fraud_risk = true THEN 1 ELSE 0 END) as fraud_bids,
			CASE 
				WHEN COUNT(*) > 0 
				THEN CAST(SUM(CASE WHEN p.fraud_risk = true THEN 1 ELSE 0 END) AS FLOAT) / COUNT(*)
				ELSE 0 
			END as fraud_rate
		FROM bid_events be
		LEFT JOIN predictions p ON be.id = p.bid_event_id
		WHERE be.user_id = $1 AND be.timestamp >= $2
		GROUP BY be.device_type, be.browser, be.os
		HAVING SUM(CASE WHEN p.fraud_risk = true THEN 1 ELSE 0 END) > 0
		ORDER BY fraud_rate DESC
		LIMIT 20
	`

	rows, err := h.bidStore.DB().QueryContext(ctx, query, userUUID, startDate)
	if err != nil {
		log.Printf("Failed to get device fraud analysis: %v", err)
		h.writeErrorResponse(w, "Failed to retrieve device fraud analysis", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type deviceFraud struct {
		DeviceType string  `json:"device_type"`
		Browser    string  `json:"browser"`
		OS         string  `json:"os"`
		TotalBids  int     `json:"total_bids"`
		FraudBids  int     `json:"fraud_bids"`
		FraudRate  float64 `json:"fraud_rate"`
	}

	var results []deviceFraud
	for rows.Next() {
		var df deviceFraud
		err := rows.Scan(&df.DeviceType, &df.Browser, &df.OS, &df.TotalBids, &df.FraudBids, &df.FraudRate)
		if err != nil {
			log.Printf("Error scanning device fraud row: %v", err)
			continue
		}
		results = append(results, df)
	}

	h.writeTRPCResponse(w, results)
}

// getGeoFraudAnalysis returns geographic fraud patterns
func (h *Handler) getGeoFraudAnalysis(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	if userID == "" {
		h.writeErrorResponse(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		h.writeErrorResponse(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Days int `json:"days"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		req.Days = 30
	}

	startDate := time.Now().AddDate(0, 0, -req.Days)

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	query := `
		SELECT 
			be.country,
			be.region,
			be.city,
			COUNT(*) as total_bids,
			SUM(CASE WHEN p.fraud_risk = true THEN 1 ELSE 0 END) as fraud_bids,
			CASE 
				WHEN COUNT(*) > 0 
				THEN CAST(SUM(CASE WHEN p.fraud_risk = true THEN 1 ELSE 0 END) AS FLOAT) / COUNT(*)
				ELSE 0 
			END as fraud_rate
		FROM bid_events be
		LEFT JOIN predictions p ON be.id = p.bid_event_id
		WHERE be.user_id = $1 AND be.timestamp >= $2
		GROUP BY be.country, be.region, be.city
		HAVING SUM(CASE WHEN p.fraud_risk = true THEN 1 ELSE 0 END) > 0
		ORDER BY fraud_bids DESC
		LIMIT 30
	`

	rows, err := h.bidStore.DB().QueryContext(ctx, query, userUUID, startDate)
	if err != nil {
		log.Printf("Failed to get geo fraud analysis: %v", err)
		h.writeErrorResponse(w, "Failed to retrieve geo fraud analysis", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type geoFraud struct {
		Country   string  `json:"country"`
		Region    string  `json:"region"`
		City      string  `json:"city"`
		TotalBids int     `json:"total_bids"`
		FraudBids int     `json:"fraud_bids"`
		FraudRate float64 `json:"fraud_rate"`
	}

	var results []geoFraud
	for rows.Next() {
		var gf geoFraud
		err := rows.Scan(&gf.Country, &gf.Region, &gf.City, &gf.TotalBids, &gf.FraudBids, &gf.FraudRate)
		if err != nil {
			log.Printf("Error scanning geo fraud row: %v", err)
			continue
		}
		results = append(results, gf)
	}

	h.writeTRPCResponse(w, results)
}

// createFraudAlert creates a new fraud alert (for manual reporting or automated detection)
func (h *Handler) createFraudAlert(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	if userID == "" {
		h.writeErrorResponse(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		h.writeErrorResponse(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req struct {
		CampaignID      string   `json:"campaign_id"`
		AlertType       string   `json:"alert_type"`
		Severity        int      `json:"severity"`
		Description     string   `json:"description"`
		AffectedUserIDs []string `json:"affected_user_ids"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.CampaignID == "" || req.AlertType == "" || req.Description == "" {
		h.writeErrorResponse(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	campaignID, err := uuid.Parse(req.CampaignID)
	if err != nil {
		h.writeErrorResponse(w, "Invalid campaign ID format", http.StatusBadRequest)
		return
	}

	// Verify campaign ownership
	var campaignUserID uuid.UUID
	err = h.bidStore.DB().QueryRow(`SELECT user_id FROM campaigns WHERE id = $1`, campaignID).Scan(&campaignUserID)
	if err != nil {
		h.writeErrorResponse(w, "Campaign not found", http.StatusNotFound)
		return
	}

	if campaignUserID != userUUID {
		h.writeErrorResponse(w, "Unauthorized", http.StatusForbidden)
		return
	}

	// Create the alert
	alert := &models.FraudAlert{
		ID:              uuid.New(),
		CampaignID:      campaignID,
		AlertType:       req.AlertType,
		Severity:        req.Severity,
		Description:     req.Description,
		AffectedUserIDs: req.AffectedUserIDs,
		DetectedAt:      time.Now(),
		Status:          "active",
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	fraudDetector := fraud.NewFraudDetector(h.bidStore.DB())
	err = fraudDetector.CreateFraudAlert(ctx, alert)
	if err != nil {
		log.Printf("Failed to create fraud alert: %v", err)
		h.writeErrorResponse(w, "Failed to create alert", http.StatusInternalServerError)
		return
	}

	h.writeTRPCResponse(w, map[string]interface{}{
		"success":  true,
		"alert_id": alert.ID.String(),
		"message":  "Fraud alert created successfully",
	})
}
