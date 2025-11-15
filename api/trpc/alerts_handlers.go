package trpc

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// Alert types and structures
type AlertType string

const (
	AlertTypeFraud       AlertType = "fraud"
	AlertTypeBudget      AlertType = "budget"
	AlertTypePerformance AlertType = "performance"
	AlertTypeModel       AlertType = "model"
	AlertTypeSystem      AlertType = "system"
	AlertTypeCampaign    AlertType = "campaign"
)

type AlertSeverity string

const (
	SeverityLow      AlertSeverity = "low"
	SeverityMedium   AlertSeverity = "medium"
	SeverityHigh     AlertSeverity = "high"
	SeverityCritical AlertSeverity = "critical"
)

type AlertStatus string

const (
	StatusUnread       AlertStatus = "unread"
	StatusRead         AlertStatus = "read"
	StatusAcknowledged AlertStatus = "acknowledged"
	StatusResolved     AlertStatus = "resolved"
	StatusDismissed    AlertStatus = "dismissed"
)

type Alert struct {
	ID             uuid.UUID              `json:"id"`
	Type           AlertType              `json:"type"`
	Severity       AlertSeverity          `json:"severity"`
	Status         AlertStatus            `json:"status"`
	Title          string                 `json:"title"`
	Message        string                 `json:"message"`
	CampaignID     *uuid.UUID             `json:"campaign_id,omitempty"`
	CampaignName   string                 `json:"campaign_name,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
	AcknowledgedAt *time.Time             `json:"acknowledged_at,omitempty"`
	ResolvedAt     *time.Time             `json:"resolved_at,omitempty"`
	Notes          string                 `json:"notes,omitempty"`
}

type AlertOverview struct {
	TotalAlerts      int                      `json:"total_alerts"`
	UnreadAlerts     int                      `json:"unread_alerts"`
	CriticalAlerts   int                      `json:"critical_alerts"`
	AlertsByType     map[AlertType]int        `json:"alerts_by_type"`
	AlertsBySeverity map[AlertSeverity]int    `json:"alerts_by_severity"`
	RecentTrend      []map[string]interface{} `json:"recent_trend"`
}

// getAlerts retrieves all alerts with optional filtering
func (h *Handler) getAlerts(w http.ResponseWriter, r *http.Request) {
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
		Type       string `json:"type"`
		Severity   string `json:"severity"`
		Status     string `json:"status"`
		CampaignID string `json:"campaign_id"`
		StartDate  string `json:"start_date"`
		EndDate    string `json:"end_date"`
		Limit      int    `json:"limit"`
		Offset     int    `json:"offset"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode request, using defaults: %v", err)
	}

	if req.Limit == 0 {
		req.Limit = 100
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	// Build query with filters
	query := `
		SELECT 
			a.id, a.type, a.severity, a.status, a.title, a.message,
			a.campaign_id, COALESCE(c.name, '') as campaign_name,
			COALESCE(a.metadata, '{}') as metadata, 
			a.created_at, a.updated_at,
			a.acknowledged_at, a.resolved_at, 
			COALESCE(a.notes, '') as notes
		FROM alerts a
		LEFT JOIN campaigns c ON a.campaign_id = c.id
		WHERE a.user_id = $1
	`
	args := []interface{}{userUUID}
	argCount := 1

	if req.Type != "" {
		argCount++
		query += fmt.Sprintf(" AND a.type = $%d", argCount)
		args = append(args, req.Type)
	}

	if req.Severity != "" {
		argCount++
		query += fmt.Sprintf(" AND a.severity = $%d", argCount)
		args = append(args, req.Severity)
	}

	if req.Status != "" {
		argCount++
		query += fmt.Sprintf(" AND a.status = $%d", argCount)
		args = append(args, req.Status)
	}

	if req.CampaignID != "" {
		campaignUUID, err := uuid.Parse(req.CampaignID)
		if err == nil {
			argCount++
			query += fmt.Sprintf(" AND a.campaign_id = $%d", argCount)
			args = append(args, campaignUUID)
		}
	}

	query += fmt.Sprintf(" ORDER BY a.created_at DESC LIMIT $%d", argCount+1)
	args = append(args, req.Limit)

	if req.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argCount+2)
		args = append(args, req.Offset)
	}

	rows, err := h.bidStore.DB().QueryContext(ctx, query, args...)
	if err != nil {
		log.Printf("Failed to query alerts: %v", err)
		h.writeErrorResponse(w, "Failed to retrieve alerts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var alerts []Alert
	for rows.Next() {
		var alert Alert
		var metadataJSON []byte
		var campaignID sql.NullString
		var acknowledgedAt sql.NullTime
		var resolvedAt sql.NullTime

		err := rows.Scan(
			&alert.ID, &alert.Type, &alert.Severity, &alert.Status,
			&alert.Title, &alert.Message, &campaignID, &alert.CampaignName,
			&metadataJSON, &alert.CreatedAt, &alert.UpdatedAt,
			&acknowledgedAt, &resolvedAt, &alert.Notes,
		)
		if err != nil {
			log.Printf("Error scanning alert row: %v", err)
			continue
		}

		// Handle nullable campaign_id
		if campaignID.Valid {
			cid, err := uuid.Parse(campaignID.String)
			if err == nil {
				alert.CampaignID = &cid
			}
		}

		// Handle nullable timestamps
		if acknowledgedAt.Valid {
			alert.AcknowledgedAt = &acknowledgedAt.Time
		}

		if resolvedAt.Valid {
			alert.ResolvedAt = &resolvedAt.Time
		}

		// Parse metadata JSON
		if len(metadataJSON) > 0 && string(metadataJSON) != "{}" {
			json.Unmarshal(metadataJSON, &alert.Metadata)
		}

		alerts = append(alerts, alert)
	}

	if alerts == nil {
		alerts = []Alert{}
	}

	h.writeTRPCResponse(w, alerts)
}

// getAlertOverview retrieves alert statistics
func (h *Handler) getAlertOverview(w http.ResponseWriter, r *http.Request) {
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

	startDate := time.Now().AddDate(0, 0, -req.Days)

	// Get total and unread counts
	var totalAlerts, unreadAlerts, criticalAlerts int
	query := `
		SELECT 
			COUNT(*) as total,
			SUM(CASE WHEN status = 'unread' THEN 1 ELSE 0 END) as unread,
			SUM(CASE WHEN severity = 'critical' THEN 1 ELSE 0 END) as critical
		FROM alerts
		WHERE user_id = $1 AND created_at >= $2
	`
	err = h.bidStore.DB().QueryRowContext(ctx, query, userUUID, startDate).Scan(
		&totalAlerts, &unreadAlerts, &criticalAlerts,
	)
	if err != nil {
		log.Printf("Failed to get alert counts: %v", err)
		h.writeErrorResponse(w, "Failed to retrieve alert overview", http.StatusInternalServerError)
		return
	}

	// Get counts by type
	alertsByType := make(map[AlertType]int)
	rows, err := h.bidStore.DB().QueryContext(ctx, `
		SELECT type, COUNT(*) as count
		FROM alerts
		WHERE user_id = $1 AND created_at >= $2
		GROUP BY type
	`, userUUID, startDate)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var alertType AlertType
			var count int
			if err := rows.Scan(&alertType, &count); err == nil {
				alertsByType[alertType] = count
			}
		}
	}

	// Get counts by severity
	alertsBySeverity := make(map[AlertSeverity]int)
	rows2, err := h.bidStore.DB().QueryContext(ctx, `
		SELECT severity, COUNT(*) as count
		FROM alerts
		WHERE user_id = $1 AND created_at >= $2
		GROUP BY severity
	`, userUUID, startDate)
	if err == nil {
		defer rows2.Close()
		for rows2.Next() {
			var severity AlertSeverity
			var count int
			if err := rows2.Scan(&severity, &count); err == nil {
				alertsBySeverity[severity] = count
			}
		}
	}

	// Get recent trend (daily counts)
	rows3, err := h.bidStore.DB().QueryContext(ctx, `
		SELECT 
			DATE(created_at) as date,
			COUNT(*) as count
		FROM alerts
		WHERE user_id = $1 AND created_at >= $2
		GROUP BY DATE(created_at)
		ORDER BY date DESC
		LIMIT 30
	`, userUUID, startDate)

	var recentTrend []map[string]interface{}
	if err == nil {
		defer rows3.Close()
		for rows3.Next() {
			var date time.Time
			var count int
			if err := rows3.Scan(&date, &count); err == nil {
				recentTrend = append(recentTrend, map[string]interface{}{
					"date":  date.Format("2006-01-02"),
					"count": count,
				})
			}
		}
	}

	overview := AlertOverview{
		TotalAlerts:      totalAlerts,
		UnreadAlerts:     unreadAlerts,
		CriticalAlerts:   criticalAlerts,
		AlertsByType:     alertsByType,
		AlertsBySeverity: alertsBySeverity,
		RecentTrend:      recentTrend,
	}

	h.writeTRPCResponse(w, overview)
}

// updateAlertStatus updates an alert's status
func (h *Handler) updateAlertStatus(w http.ResponseWriter, r *http.Request) {
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
		AlertID string      `json:"alert_id"`
		Status  AlertStatus `json:"status"`
		Notes   string      `json:"notes"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	alertID, err := uuid.Parse(req.AlertID)
	if err != nil {
		h.writeErrorResponse(w, "Invalid alert ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// Verify ownership
	var ownerID uuid.UUID
	err = h.bidStore.DB().QueryRowContext(ctx, `SELECT user_id FROM alerts WHERE id = $1`, alertID).Scan(&ownerID)
	if err != nil {
		if err == sql.ErrNoRows {
			h.writeErrorResponse(w, "Alert not found", http.StatusNotFound)
		} else {
			h.writeErrorResponse(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	if ownerID != userUUID {
		h.writeErrorResponse(w, "Unauthorized", http.StatusForbidden)
		return
	}

	// Update status
	query := `
		UPDATE alerts
		SET status = $1, updated_at = $2, notes = $3
	`
	args := []interface{}{req.Status, time.Now(), req.Notes}

	if req.Status == StatusAcknowledged {
		query += ", acknowledged_at = $4"
		args = append(args, time.Now())
	} else if req.Status == StatusResolved {
		query += ", resolved_at = $4"
		args = append(args, time.Now())
	}

	query += fmt.Sprintf(" WHERE id = $%d", len(args)+1)
	args = append(args, alertID)

	_, err = h.bidStore.DB().ExecContext(ctx, query, args...)
	if err != nil {
		log.Printf("Failed to update alert status: %v", err)
		h.writeErrorResponse(w, "Failed to update alert", http.StatusInternalServerError)
		return
	}

	h.writeTRPCResponse(w, map[string]interface{}{
		"success": true,
		"message": "Alert status updated successfully",
	})
}

// bulkUpdateAlerts updates multiple alerts' statuses
func (h *Handler) bulkUpdateAlerts(w http.ResponseWriter, r *http.Request) {
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
		AlertIDs []string    `json:"alert_ids"`
		Status   AlertStatus `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(req.AlertIDs) == 0 {
		h.writeErrorResponse(w, "No alert IDs provided", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	// Convert alert IDs to UUIDs
	var alertUUIDs []uuid.UUID
	for _, idStr := range req.AlertIDs {
		id, err := uuid.Parse(idStr)
		if err == nil {
			alertUUIDs = append(alertUUIDs, id)
		}
	}

	// Update all alerts owned by user
	query := `
		UPDATE alerts
		SET status = $1, updated_at = $2
		WHERE user_id = $3 AND id = ANY($4)
	`

	result, err := h.bidStore.DB().ExecContext(ctx, query, req.Status, time.Now(), userUUID, alertUUIDs)
	if err != nil {
		log.Printf("Failed to bulk update alerts: %v", err)
		h.writeErrorResponse(w, "Failed to update alerts", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()

	h.writeTRPCResponse(w, map[string]interface{}{
		"success":       true,
		"message":       "Alerts updated successfully",
		"updated_count": rowsAffected,
	})
}
