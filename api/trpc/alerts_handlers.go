package trpc

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
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
	CampaignID     string                 `json:"campaign_id,omitempty"`
	CampaignName   string                 `json:"campaign_name,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
	AcknowledgedAt *time.Time             `json:"acknowledged_at,omitempty"`
	ResolvedAt     *time.Time             `json:"resolved_at,omitempty"`
	Notes          string                 `json:"notes,omitempty"`
}

// Request types
type GetAlertsRequest struct {
	Type       string `json:"type"`
	Severity   string `json:"severity"`
	Status     string `json:"status"`
	CampaignID string `json:"campaign_id"`
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
	Limit      int    `json:"limit"`
	Offset     int    `json:"offset"`
}

type AlertOverviewRequest struct {
	Days int `json:"days"`
}

type UpdateAlertStatusRequest struct {
	AlertID string `json:"alert_id"`
	Status  string `json:"status"`
	Notes   string `json:"notes"`
}

type BulkUpdateAlertsRequest struct {
	AlertIDs []string `json:"alert_ids"`
	Status   string   `json:"status"`
}

// ============================================================================
// REFACTORED HANDLERS
// ============================================================================

// getAlerts returns filtered alerts
func (h *Handler) getAlerts(ctx context.Context, userID uuid.UUID, req interface{}) (interface{}, error) {
	params := req.(*GetAlertsRequest)

	// Set default limit
	if params.Limit == 0 {
		params.Limit = 100
	}

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
	args := []interface{}{userID}
	argCount := 1

	// Add filters
	if params.Type != "" {
		argCount++
		query += fmt.Sprintf(" AND a.type = $%d", argCount)
		args = append(args, params.Type)
	}

	if params.Severity != "" {
		argCount++
		query += fmt.Sprintf(" AND a.severity = $%d", argCount)
		args = append(args, params.Severity)
	}

	if params.Status != "" {
		argCount++
		query += fmt.Sprintf(" AND a.status = $%d", argCount)
		args = append(args, params.Status)
	}

	if params.CampaignID != "" {
		campaignUUID, err := uuid.Parse(params.CampaignID)
		if err == nil {
			argCount++
			query += fmt.Sprintf(" AND a.campaign_id = $%d", argCount)
			args = append(args, campaignUUID)
		}
	}

	if params.StartDate != "" {
		startTime, err := time.Parse("2006-01-02", params.StartDate)
		if err == nil {
			argCount++
			query += fmt.Sprintf(" AND a.created_at >= $%d", argCount)
			args = append(args, startTime)
		}
	}

	if params.EndDate != "" {
		endTime, err := time.Parse("2006-01-02", params.EndDate)
		if err == nil {
			endTime = endTime.Add(24 * time.Hour)
			argCount++
			query += fmt.Sprintf(" AND a.created_at < $%d", argCount)
			args = append(args, endTime)
		}
	}

	query += " ORDER BY a.created_at DESC"

	argCount++
	query += fmt.Sprintf(" LIMIT $%d", argCount)
	args = append(args, params.Limit)

	if params.Offset > 0 {
		argCount++
		query += fmt.Sprintf(" OFFSET $%d", argCount)
		args = append(args, params.Offset)
	}

	// Execute query
	rows, err := h.bidStore.DB().QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query alerts: %w", err)
	}
	defer rows.Close()

	var alerts []Alert
	for rows.Next() {
		var alert Alert
		var campaignID sql.NullString
		var campaignName sql.NullString
		var acknowledgedAt sql.NullTime
		var resolvedAt sql.NullTime
		var metadataJSON string

		err := rows.Scan(
			&alert.ID,
			&alert.Type,
			&alert.Severity,
			&alert.Status,
			&alert.Title,
			&alert.Message,
			&campaignID,
			&campaignName,
			&metadataJSON,
			&alert.CreatedAt,
			&alert.UpdatedAt,
			&acknowledgedAt,
			&resolvedAt,
			&alert.Notes,
		)
		if err != nil {
			log.Printf("Error scanning alert row: %v", err)
			continue
		}

		if campaignID.Valid {
			alert.CampaignID = campaignID.String
		}
		if campaignName.Valid {
			alert.CampaignName = campaignName.String
		}
		if acknowledgedAt.Valid {
			alert.AcknowledgedAt = &acknowledgedAt.Time
		}
		if resolvedAt.Valid {
			alert.ResolvedAt = &resolvedAt.Time
		}

		// Parse metadata JSON if present
		if metadataJSON != "" && metadataJSON != "{}" {
			var metadata map[string]interface{}
			if err := json.Unmarshal([]byte(metadataJSON), &metadata); err == nil {
				alert.Metadata = metadata
			}
		}

		alerts = append(alerts, alert)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating alert rows: %w", err)
	}

	// Ensure we always return an array, even if empty
	if alerts == nil {
		alerts = []Alert{}
	}

	return alerts, nil
}

// getAlertOverview returns alert statistics
func (h *Handler) getAlertOverview(ctx context.Context, userID uuid.UUID, req interface{}) (interface{}, error) {
	params := req.(*AlertOverviewRequest)

	days := params.Days
	if days <= 0 || days > 365 {
		days = 30
	}

	startDate := time.Now().AddDate(0, 0, -days)

	// Get total counts
	query := `
		SELECT
			COALESCE(COUNT(*), 0) as total_alerts,
			COALESCE(SUM(CASE WHEN status = 'unread' THEN 1 ELSE 0 END), 0) as unread_alerts,
			COALESCE(SUM(CASE WHEN severity = 'critical' THEN 1 ELSE 0 END), 0) as critical_alerts
		FROM alerts
		WHERE user_id = $1 AND created_at >= $2
	`

	var overview struct {
		TotalAlerts    int `json:"total_alerts"`
		UnreadAlerts   int `json:"unread_alerts"`
		CriticalAlerts int `json:"critical_alerts"`
	}

	err := h.bidStore.DB().QueryRowContext(ctx, query, userID, startDate).Scan(
		&overview.TotalAlerts,
		&overview.UnreadAlerts,
		&overview.CriticalAlerts,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get alert overview: %w", err)
	}

	// Get alerts by type
	alertsByType := make(map[string]int)
	typeQuery := `
		SELECT type, COUNT(*) as count
		FROM alerts
		WHERE user_id = $1 AND created_at >= $2
		GROUP BY type
	`
	rows, err := h.bidStore.DB().QueryContext(ctx, typeQuery, userID, startDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get alerts by type: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var alertType string
		var count int
		if err := rows.Scan(&alertType, &count); err != nil {
			log.Printf("Error scanning alert type: %v", err)
			continue
		}
		alertsByType[alertType] = count
	}

	// Get alerts by severity
	alertsBySeverity := make(map[string]int)
	severityQuery := `
		SELECT severity, COUNT(*) as count
		FROM alerts
		WHERE user_id = $1 AND created_at >= $2
		GROUP BY severity
	`
	rows, err = h.bidStore.DB().QueryContext(ctx, severityQuery, userID, startDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get alerts by severity: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var severity string
		var count int
		if err := rows.Scan(&severity, &count); err != nil {
			log.Printf("Error scanning severity: %v", err)
			continue
		}
		alertsBySeverity[severity] = count
	}

	// Get recent trend (last 30 days)
	trendQuery := `
		SELECT DATE(created_at) as date, COUNT(*) as count
		FROM alerts
		WHERE user_id = $1 AND created_at >= $2
		GROUP BY DATE(created_at)
		ORDER BY date
	`
	rows, err = h.bidStore.DB().QueryContext(ctx, trendQuery, userID, startDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get alert trend: %w", err)
	}
	defer rows.Close()

	var recentTrend []struct {
		Date  string `json:"date"`
		Count int    `json:"count"`
	}
	for rows.Next() {
		var date time.Time
		var count int
		if err := rows.Scan(&date, &count); err != nil {
			log.Printf("Error scanning trend: %v", err)
			continue
		}
		recentTrend = append(recentTrend, struct {
			Date  string `json:"date"`
			Count int    `json:"count"`
		}{
			Date:  date.Format("2006-01-02"),
			Count: count,
		})
	}

	return map[string]interface{}{
		"total_alerts":       overview.TotalAlerts,
		"unread_alerts":      overview.UnreadAlerts,
		"critical_alerts":    overview.CriticalAlerts,
		"alerts_by_type":     alertsByType,
		"alerts_by_severity": alertsBySeverity,
		"recent_trend":       recentTrend,
	}, nil
}

// updateAlertStatus updates a single alert's status
func (h *Handler) updateAlertStatus(ctx context.Context, userID uuid.UUID, req interface{}) (interface{}, error) {
	params := req.(*UpdateAlertStatusRequest)

	if params.AlertID == "" {
		return nil, fmt.Errorf("alert_id is required")
	}

	alertUUID, err := uuid.Parse(params.AlertID)
	if err != nil {
		return nil, fmt.Errorf("invalid alert_id format")
	}

	validStatuses := map[string]bool{
		"unread":       true,
		"read":         true,
		"acknowledged": true,
		"resolved":     true,
		"dismissed":    true,
	}

	if !validStatuses[params.Status] {
		return nil, fmt.Errorf("invalid status: must be one of unread, read, acknowledged, resolved, dismissed")
	}

	// Update query
	updateFields := []string{"status = $3", "updated_at = NOW()"}
	args := []interface{}{alertUUID, userID, params.Status}
	argCount := 3

	if params.Status == "acknowledged" {
		updateFields = append(updateFields, "acknowledged_at = NOW()")
	}

	if params.Status == "resolved" {
		updateFields = append(updateFields, "resolved_at = NOW()")
	}

	if params.Notes != "" {
		argCount++
		updateFields = append(updateFields, fmt.Sprintf("notes = $%d", argCount))
		args = append(args, params.Notes)
	}

	query := fmt.Sprintf(`
		UPDATE alerts
		SET %s
		WHERE id = $1 AND user_id = $2
	`, joinStrings(updateFields, ", "))

	result, err := h.bidStore.DB().ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update alert status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("alert not found or unauthorized")
	}

	return map[string]interface{}{
		"success": true,
		"message": "Alert status updated successfully",
	}, nil
}

// bulkUpdateAlerts updates multiple alerts at once
func (h *Handler) bulkUpdateAlerts(ctx context.Context, userID uuid.UUID, req interface{}) (interface{}, error) {
	params := req.(*BulkUpdateAlertsRequest)

	if len(params.AlertIDs) == 0 {
		return nil, fmt.Errorf("alert_ids is required")
	}

	if len(params.AlertIDs) > 100 {
		return nil, fmt.Errorf("cannot update more than 100 alerts at once")
	}

	validStatuses := map[string]bool{
		"unread":       true,
		"read":         true,
		"acknowledged": true,
		"resolved":     true,
		"dismissed":    true,
	}

	if !validStatuses[params.Status] {
		return nil, fmt.Errorf("invalid status")
	}

	// Parse alert IDs
	var alertUUIDs []uuid.UUID
	for _, idStr := range params.AlertIDs {
		alertUUID, err := uuid.Parse(idStr)
		if err != nil {
			log.Printf("Invalid alert ID: %s", idStr)
			continue
		}
		alertUUIDs = append(alertUUIDs, alertUUID)
	}

	if len(alertUUIDs) == 0 {
		return nil, fmt.Errorf("no valid alert IDs provided")
	}

	// Build update query
	updateFields := "status = $2, updated_at = NOW()"
	if params.Status == "acknowledged" {
		updateFields += ", acknowledged_at = NOW()"
	}
	if params.Status == "resolved" {
		updateFields += ", resolved_at = NOW()"
	}

	query := fmt.Sprintf(`
		UPDATE alerts
		SET %s
		WHERE user_id = $1 AND id = ANY($3)
	`, updateFields)

	result, err := h.bidStore.DB().ExecContext(ctx, query, userID, params.Status, alertUUIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to bulk update alerts: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return map[string]interface{}{
		"success":       true,
		"message":       "Alerts updated successfully",
		"updated_count": rowsAffected,
	}, nil
}

// Helper function
func joinStrings(strings []string, separator string) string {
	result := ""
	for i, s := range strings {
		if i > 0 {
			result += separator
		}
		result += s
	}
	return result
}
