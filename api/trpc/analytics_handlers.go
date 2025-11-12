package trpc

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// NOTE: All struct definitions (PerformanceMetrics, KeywordAnalysis, etc.)
// have been removed from this file to resolve the "redeclared" compilation
// error. They are now assumed to be defined in a shared file like types.go
// within the 'trpc' package, as suggested by your error logs.

// (Handler, GetUserIDFromContext, writeErrorResponse, writeTRPCResponse are assumed to be defined elsewhere)

// --- Handlers with Context Safety Applied ---

// getPerformanceOverview returns overall performance metrics
func (h *Handler) getPerformanceOverview(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	if userID == "" {
		h.writeErrorResponse(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	// Parse userID to UUID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		log.Printf("Invalid user ID format: %v", err)
		h.writeErrorResponse(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Setup context with a 60 second database timeout
	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()

	// FIX: Decode as plain object, NOT array
	var req struct {
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode request body: %v", err)
		h.writeErrorResponse(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	startDate, endDate := parseDateRange(req.StartDate, req.EndDate)

	query := `
		SELECT
			COALESCE(COUNT(*), 0) as total_bids,
			COALESCE(SUM(CASE WHEN won THEN 1 ELSE 0 END), 0) as won_bids,
			COALESCE(SUM(CASE WHEN converted THEN 1 ELSE 0 END), 0) as conversions,
			COALESCE(SUM(CASE WHEN won THEN bid_price ELSE 0 END), 0) as total_spend,
			COALESCE(SUM(CASE WHEN converted THEN bid_price * 2.5 ELSE 0 END), 0) as revenue,
			CASE 
				WHEN COUNT(*) > 0 THEN CAST(SUM(CASE WHEN won THEN 1 ELSE 0 END) AS FLOAT) / COUNT(*)
				ELSE 0 
			END as win_rate,
			CASE 
				WHEN SUM(CASE WHEN won THEN 1 ELSE 0 END) > 0 
				THEN CAST(SUM(CASE WHEN converted THEN 1 ELSE 0 END) AS FLOAT) / SUM(CASE WHEN won THEN 1 ELSE 0 END)
				ELSE 0 
			END as conversion_rate,
			CASE 
				WHEN COUNT(*) > 0 THEN SUM(bid_price) / COUNT(*)
				ELSE 0 
			END as average_bid
		FROM bid_events
		WHERE user_id = $1 
			AND timestamp BETWEEN $2 AND $3
	`

	var metrics PerformanceMetrics

	err = h.bidStore.DB().QueryRowContext(ctx, query, userUUID, startDate, endDate).Scan(
		&metrics.TotalBids,
		&metrics.WonBids,
		&metrics.Conversions,
		&metrics.TotalSpend,
		&metrics.Revenue,
		&metrics.WinRate,
		&metrics.ConversionRate,
		&metrics.AverageBid,
	)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Printf("Failed to get performance overview: Query timed out after 60s")
			h.writeErrorResponse(w, "Query timed out. Try a smaller date range.", http.StatusGatewayTimeout) // 504
		} else {
			log.Printf("Failed to get performance overview: %v", err)
			h.writeErrorResponse(w, "Failed to retrieve performance metrics", http.StatusInternalServerError) // 500
		}
		return
	}

	if metrics.Conversions > 0 {
		metrics.CPA = metrics.TotalSpend / float64(metrics.Conversions)
	}
	if metrics.TotalSpend > 0 {
		metrics.ROAS = metrics.Revenue / metrics.TotalSpend
	}

	h.writeTRPCResponse(w, metrics)
}

// getKeywordAnalysis returns keyword performance breakdown
func (h *Handler) getKeywordAnalysis(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	if userID == "" {
		h.writeErrorResponse(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	// Parse userID to UUID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		log.Printf("Invalid user ID format: %v", err)
		h.writeErrorResponse(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()

	var req struct {
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
		Limit     int    `json:"limit"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	startDate, endDate := parseDateRange(req.StartDate, req.EndDate)
	limit := 20
	if req.Limit > 0 && req.Limit <= 100 {
		limit = req.Limit
	}

	query := `
		WITH keyword_stats AS (
			SELECT 
				UNNEST(keywords) as keyword,
				COUNT(*) as total_bids,
				SUM(CASE WHEN won THEN 1 ELSE 0 END) as won_bids,
				SUM(CASE WHEN converted THEN 1 ELSE 0 END) as conversions,
				SUM(CASE WHEN won THEN bid_price ELSE 0 END) as spend,
				SUM(CASE WHEN converted THEN bid_price * 2.5 ELSE 0 END) as revenue
			FROM bid_events
			WHERE user_id = $1 
				AND timestamp BETWEEN $2 AND $3
				AND keywords IS NOT NULL
			GROUP BY keyword
		)
		SELECT 
			keyword,
			total_bids,
			won_bids,
			conversions,
			spend,
			revenue,
			CASE WHEN total_bids > 0 THEN CAST(won_bids AS FLOAT) / total_bids ELSE 0 END as win_rate,
			CASE WHEN won_bids > 0 THEN CAST(conversions AS FLOAT) / won_bids ELSE 0 END as conversion_rate,
			CASE WHEN conversions > 0 THEN spend / conversions ELSE 0 END as cpa,
			CASE WHEN spend > 0 THEN revenue / spend ELSE 0 END as roas
		FROM keyword_stats
		ORDER BY spend DESC
		LIMIT $4
	`

	// Using QueryContext with userUUID
	rows, err := h.bidStore.DB().QueryContext(ctx, query, userUUID, startDate, endDate, limit)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Printf("Failed to get keyword analysis: Query timed out after 60s")
			h.writeErrorResponse(w, "Query timed out. Try a smaller date range.", http.StatusGatewayTimeout)
		} else {
			log.Printf("Failed to execute keyword query: %v", err)
			h.writeErrorResponse(w, "Failed to retrieve keyword analysis", http.StatusInternalServerError)
		}
		return
	}
	defer rows.Close()

	var keywords []KeywordAnalysis
	for rows.Next() {
		var kw KeywordAnalysis
		err := rows.Scan(
			&kw.Keyword,
			&kw.TotalBids,
			&kw.WonBids,
			&kw.Conversions,
			&kw.Spend,
			&kw.Revenue,
			&kw.WinRate,
			&kw.ConversionRate,
			&kw.CPA,
			&kw.ROAS,
		)
		if err != nil {
			log.Printf("Failed to scan row for keyword analysis: %v", err)
			continue
		}
		keywords = append(keywords, kw)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error after row iteration for keyword analysis: %v", err)
	}

	h.writeTRPCResponse(w, keywords)
}

// getDeviceBreakdown returns device-specific performance
func (h *Handler) getDeviceBreakdown(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	if userID == "" {
		h.writeErrorResponse(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	// Parse userID to UUID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		log.Printf("Invalid user ID format: %v", err)
		h.writeErrorResponse(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()

	var req struct {
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	startDate, endDate := parseDateRange(req.StartDate, req.EndDate)

	query := `
		SELECT 
			COALESCE(device_type, 'Unknown') AS device_type,
			COUNT(*) as total_bids,
			SUM(CASE WHEN won THEN 1 ELSE 0 END) as won_bids,
			SUM(CASE WHEN converted THEN 1 ELSE 0 END) as conversions,
			SUM(CASE WHEN won THEN bid_price ELSE 0 END) as spend,
			CASE WHEN COUNT(*) > 0 THEN CAST(SUM(CASE WHEN won THEN 1 ELSE 0 END) AS FLOAT) / COUNT(*) ELSE 0 END as win_rate,
			CASE WHEN SUM(CASE WHEN won THEN 1 ELSE 0 END) > 0 THEN CAST(SUM(CASE WHEN converted THEN 1 ELSE 0 END) AS FLOAT) / SUM(CASE WHEN won THEN 1 ELSE 0 END) ELSE 0 END as conversion_rate,
			CASE WHEN COUNT(*) > 0 THEN SUM(bid_price) / COUNT(*) ELSE 0 END as average_bid
		FROM bid_events
		WHERE user_id = $1 
			AND timestamp BETWEEN $2 AND $3
		GROUP BY device_type
		ORDER BY total_bids DESC
	`

	// Using QueryContext with userUUID
	rows, err := h.bidStore.DB().QueryContext(ctx, query, userUUID, startDate, endDate)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			h.writeErrorResponse(w, "Query timed out. Try a smaller date range.", http.StatusGatewayTimeout)
		} else {
			log.Printf("Failed to get device breakdown: %v", err)
			h.writeErrorResponse(w, "Failed to retrieve device breakdown", http.StatusInternalServerError)
		}
		return
	}
	defer rows.Close()

	var devices []DeviceBreakdown
	for rows.Next() {
		var device DeviceBreakdown
		err := rows.Scan(
			&device.DeviceType,
			&device.TotalBids,
			&device.WonBids,
			&device.Conversions,
			&device.Spend,
			&device.WinRate,
			&device.ConversionRate,
			&device.AverageBid,
		)
		if err != nil {
			log.Printf("DeviceBreakdown: Failed to scan row: %v", err)
			continue
		}
		devices = append(devices, device)
	}

	h.writeTRPCResponse(w, devices)
}

// getGeoBreakdown returns geographic performance
func (h *Handler) getGeoBreakdown(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	if userID == "" {
		h.writeErrorResponse(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	// Parse userID to UUID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		log.Printf("Invalid user ID format: %v", err)
		h.writeErrorResponse(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()

	var req struct {
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
		Limit     int    `json:"limit"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	startDate, endDate := parseDateRange(req.StartDate, req.EndDate)
	limit := 20
	if req.Limit > 0 && req.Limit <= 100 {
		limit = req.Limit
	}

	query := `
		SELECT 
			COALESCE(country, 'Unknown') AS country,
			COALESCE(region, 'Unknown') AS region,
			COUNT(*) as total_bids,
			SUM(CASE WHEN won THEN 1 ELSE 0 END) as won_bids,
			SUM(CASE WHEN converted THEN 1 ELSE 0 END) as conversions,
			SUM(CASE WHEN won THEN bid_price ELSE 0 END) as spend,
			CASE WHEN COUNT(*) > 0 THEN CAST(SUM(CASE WHEN won THEN 1 ELSE 0 END) AS FLOAT) / COUNT(*) ELSE 0 END as win_rate,
			CASE WHEN SUM(CASE WHEN won THEN 1 ELSE 0 END) > 0 THEN CAST(SUM(CASE WHEN converted THEN 1 ELSE 0 END) AS FLOAT) / SUM(CASE WHEN won THEN 1 ELSE 0 END) ELSE 0 END as conversion_rate
		FROM bid_events
		WHERE user_id = $1 
			AND timestamp BETWEEN $2 AND $3
		GROUP BY country, region
		ORDER BY total_bids DESC
		LIMIT $4
	`

	// Using QueryContext with userUUID
	rows, err := h.bidStore.DB().QueryContext(ctx, query, userUUID, startDate, endDate, limit)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			h.writeErrorResponse(w, "Query timed out. Try a smaller date range.", http.StatusGatewayTimeout)
		} else {
			log.Printf("Failed to get geo breakdown: %v", err)
			h.writeErrorResponse(w, "Failed to retrieve geographic breakdown", http.StatusInternalServerError)
		}
		return
	}
	defer rows.Close()

	var geos []GeoBreakdown
	for rows.Next() {
		var geo GeoBreakdown
		err := rows.Scan(
			&geo.Country,
			&geo.Region,
			&geo.TotalBids,
			&geo.WonBids,
			&geo.Conversions,
			&geo.Spend,
			&geo.WinRate,
			&geo.ConversionRate,
		)
		if err != nil {
			log.Printf("GeoBreakdown: Failed to scan row: %v", err)
			continue
		}
		geos = append(geos, geo)
	}

	h.writeTRPCResponse(w, geos)
}

// getHourlyPerformance returns time-based performance by hour
func (h *Handler) getHourlyPerformance(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	if userID == "" {
		h.writeErrorResponse(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	// Parse userID to UUID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		log.Printf("Invalid user ID format: %v", err)
		h.writeErrorResponse(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()

	var req struct {
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	startDate, endDate := parseDateRange(req.StartDate, req.EndDate)

	query := `
		SELECT 
			EXTRACT(HOUR FROM timestamp) as hour,
			COUNT(*) as total_bids,
			SUM(CASE WHEN won THEN 1 ELSE 0 END) as won_bids,
			SUM(CASE WHEN converted THEN 1 ELSE 0 END) as conversions,
			SUM(CASE WHEN won THEN bid_price ELSE 0 END) as spend,
			CASE WHEN COUNT(*) > 0 THEN CAST(SUM(CASE WHEN won THEN 1 ELSE 0 END) AS FLOAT) / COUNT(*) ELSE 0 END as win_rate,
			CASE WHEN SUM(CASE WHEN won THEN 1 ELSE 0 END) > 0 THEN CAST(SUM(CASE WHEN converted THEN 1 ELSE 0 END) AS FLOAT) / SUM(CASE WHEN won THEN 1 ELSE 0 END) ELSE 0 END as conversion_rate,
			CASE WHEN COUNT(*) > 0 THEN SUM(bid_price) / COUNT(*) ELSE 0 END as average_bid
		FROM bid_events
		WHERE user_id = $1 
			AND timestamp BETWEEN $2 AND $3
		GROUP BY hour
		ORDER BY hour
	`

	// Using QueryContext with userUUID
	rows, err := h.bidStore.DB().QueryContext(ctx, query, userUUID, startDate, endDate)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			h.writeErrorResponse(w, "Query timed out. Try a smaller date range.", http.StatusGatewayTimeout)
		} else {
			log.Printf("Failed to get hourly performance: %v", err)
			h.writeErrorResponse(w, "Failed to retrieve hourly performance", http.StatusInternalServerError)
		}
		return
	}
	defer rows.Close()

	var hourly []HourlyPerformance
	for rows.Next() {
		var h HourlyPerformance
		var hourFloat float64
		err := rows.Scan(
			&hourFloat,
			&h.TotalBids,
			&h.WonBids,
			&h.Conversions,
			&h.Spend,
			&h.WinRate,
			&h.ConversionRate,
			&h.AverageBid,
		)
		if err != nil {
			log.Printf("HourlyPerformance: Failed to scan row: %v", err)
			continue
		}
		h.Hour = int(hourFloat)
		hourly = append(hourly, h)
	}

	h.writeTRPCResponse(w, hourly)
}

// getDailyTrends returns daily performance trends
func (h *Handler) getDailyTrends(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	if userID == "" {
		h.writeErrorResponse(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	// Parse userID to UUID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		log.Printf("Invalid user ID format: %v", err)
		h.writeErrorResponse(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()

	var req struct {
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	startDate, endDate := parseDateRange(req.StartDate, req.EndDate)

	query := `
		SELECT 
			DATE(timestamp) as date,
			COUNT(*) as total_bids,
			SUM(CASE WHEN won THEN 1 ELSE 0 END) as won_bids,
			SUM(CASE WHEN converted THEN 1 ELSE 0 END) as conversions,
			SUM(CASE WHEN won THEN bid_price ELSE 0 END) as spend,
			SUM(CASE WHEN converted THEN bid_price * 2.5 ELSE 0 END) as revenue,
			CASE WHEN COUNT(*) > 0 THEN CAST(SUM(CASE WHEN won THEN 1 ELSE 0 END) AS FLOAT) / COUNT(*) ELSE 0 END as win_rate,
			CASE WHEN SUM(CASE WHEN won THEN 1 ELSE 0 END) > 0 THEN CAST(SUM(CASE WHEN converted THEN 1 ELSE 0 END) AS FLOAT) / SUM(CASE WHEN won THEN 1 ELSE 0 END) ELSE 0 END as conversion_rate,
			CASE WHEN SUM(CASE WHEN converted THEN 1 ELSE 0 END) > 0 THEN SUM(CASE WHEN won THEN bid_price ELSE 0 END) / SUM(CASE WHEN converted THEN 1 ELSE 0 END) ELSE 0 END as cpa
		FROM bid_events
		WHERE user_id = $1 
			AND timestamp BETWEEN $2 AND $3
		GROUP BY DATE(timestamp)
		ORDER BY date
	`

	// Using QueryContext with userUUID
	rows, err := h.bidStore.DB().QueryContext(ctx, query, userUUID, startDate, endDate)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			h.writeErrorResponse(w, "Query timed out. Try a smaller date range.", http.StatusGatewayTimeout)
		} else {
			log.Printf("Failed to get daily trends: %v", err)
			h.writeErrorResponse(w, "Failed to retrieve daily trends", http.StatusInternalServerError)
		}
		return
	}
	defer rows.Close()

	var trends []DailyTrend
	for rows.Next() {
		var trend DailyTrend
		var date time.Time
		err := rows.Scan(
			&date,
			&trend.TotalBids,
			&trend.WonBids,
			&trend.Conversions,
			&trend.Spend,
			&trend.Revenue,
			&trend.WinRate,
			&trend.ConversionRate,
			&trend.CPA,
		)
		if err != nil {
			log.Printf("DailyTrends: Failed to scan row: %v", err)
			continue
		}
		trend.Date = date.Format("2006-01-02")
		trends = append(trends, trend)
	}

	h.writeTRPCResponse(w, trends)
}

// getCompetitiveAnalysis returns competitive insights
func (h *Handler) getCompetitiveAnalysis(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	if userID == "" {
		h.writeErrorResponse(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	// Parse userID to UUID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		log.Printf("Invalid user ID format: %v", err)
		h.writeErrorResponse(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()

	var req struct {
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	startDate, endDate := parseDateRange(req.StartDate, req.EndDate)

	query := `
		SELECT 
			segment_category,
			CASE WHEN COUNT(*) > 0 THEN CAST(SUM(CASE WHEN won THEN 1 ELSE 0 END) AS FLOAT) / COUNT(*) ELSE 0 END as our_win_rate,
			AVG(bid_price) as market_average_bid,
			AVG(CASE WHEN user_id = $1 THEN bid_price ELSE NULL END) as our_average_bid,
			AVG(floor_price) as average_floor_price,
			CASE WHEN AVG(floor_price) > 0 THEN AVG(bid_price) / AVG(floor_price) ELSE 0 END as competition_intensity,
			COUNT(*) as total_opportunities
		FROM bid_events
		WHERE timestamp BETWEEN $2 AND $3
		GROUP BY segment_category
		ORDER BY total_opportunities DESC
		LIMIT 10
	`
	log.Printf("[DEBUG] Executing competitive analysis query:\n%s\nArgs: userUUID=%s, startDate=%v, endDate=%v",
		query, userUUID, startDate, endDate)
	// Using QueryContext with userUUID
	rows, err := h.bidStore.DB().QueryContext(ctx, query, userUUID, startDate, endDate)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			h.writeErrorResponse(w, "Query timed out. Try a smaller date range.", http.StatusGatewayTimeout)
		} else {
			log.Printf("Failed to get competitive analysis: %v", err)
			h.writeErrorResponse(w, "Failed to retrieve competitive analysis", http.StatusInternalServerError)
		}
		return
	}
	defer rows.Close()

	var competitive []CompetitiveAnalysis
	for rows.Next() {
		var comp CompetitiveAnalysis
		// Note: Using sql.NullFloat64 for OurAverageBid to correctly handle NULL if no bids were found for this user/segment
		err := rows.Scan(
			&comp.SegmentCategory,
			&comp.OurWinRate,
			&comp.MarketAverageBid,
			&comp.OurAverageBid,
			&comp.AverageFloorPrice,
			&comp.CompetitionIntensity,
			&comp.TotalOpportunities,
		)
		if err != nil {
			log.Printf("CompetitiveAnalysis: Failed to scan row: %v", err)
			continue
		}
		competitive = append(competitive, comp)
	}

	h.writeTRPCResponse(w, competitive)
}

// getCampaignComparison returns comparison metrics across campaigns
func (h *Handler) getCampaignComparison(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	if userID == "" {
		h.writeErrorResponse(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	// Parse userID to UUID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		log.Printf("Invalid user ID format: %v", err)
		h.writeErrorResponse(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 60*time.Second)
	defer cancel()

	var req struct {
		StartDate   string   `json:"start_date"`
		EndDate     string   `json:"end_date"`
		CampaignIDs []string `json:"campaign_ids"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	startDate, endDate := parseDateRange(req.StartDate, req.EndDate)

	// If specific campaigns requested, use them; otherwise get all user campaigns
	campaignsFilter := ""
	if len(req.CampaignIDs) > 0 {
		campaignsFilter = " AND c.id = ANY($4)"
	}

	query := `
		SELECT 
			c.id,
			c.name,
			COUNT(be.id) as total_bids,
			SUM(CASE WHEN be.won THEN 1 ELSE 0 END) as won_bids,
			SUM(CASE WHEN be.converted THEN 1 ELSE 0 END) as conversions,
			SUM(CASE WHEN be.won THEN be.bid_price ELSE 0 END) as spend,
			CASE WHEN COUNT(be.id) > 0 THEN CAST(SUM(CASE WHEN be.won THEN 1 ELSE 0 END) AS FLOAT) / COUNT(be.id) ELSE 0 END as win_rate,
			CASE WHEN SUM(CASE WHEN be.won THEN 1 ELSE 0 END) > 0 THEN CAST(SUM(CASE WHEN be.converted THEN 1 ELSE 0 END) AS FLOAT) / SUM(CASE WHEN be.won THEN 1 ELSE 0 END) ELSE 0 END as conversion_rate
		FROM campaigns c
		LEFT JOIN bid_events be ON c.id = be.campaign_id 
			AND be.timestamp BETWEEN $2 AND $3
		WHERE c.user_id = $1` + campaignsFilter + `
		GROUP BY c.id, c.name
		ORDER BY spend DESC
	`

	var rows *sql.Rows

	// Execute query based on whether campaign IDs are provided
	if len(req.CampaignIDs) > 0 {
		// Using QueryContext with userUUID
		rows, err = h.bidStore.DB().QueryContext(ctx, query, userUUID, startDate, endDate, req.CampaignIDs)
	} else {
		// Using QueryContext with userUUID
		rows, err = h.bidStore.DB().QueryContext(ctx, query, userUUID, startDate, endDate)
	}

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			h.writeErrorResponse(w, "Query timed out. Try a smaller date range.", http.StatusGatewayTimeout)
		} else {
			log.Printf("Failed to get campaign comparison: %v", err)
			h.writeErrorResponse(w, "Failed to retrieve campaign comparison", http.StatusInternalServerError)
		}
		return
	}
	defer rows.Close()

	var campaigns []CampaignComparison
	for rows.Next() {
		var camp CampaignComparison
		var id uuid.UUID
		err := rows.Scan(
			&id,
			&camp.CampaignName,
			&camp.TotalBids,
			&camp.WonBids,
			&camp.Conversions,
			&camp.Spend,
			&camp.WinRate,
			&camp.ConversionRate,
		)
		if err != nil {
			log.Printf("CampaignComparison: Failed to scan row: %v", err)
			continue
		}
		camp.CampaignID = id.String()
		campaigns = append(campaigns, camp)
	}

	h.writeTRPCResponse(w, campaigns)
}

// Helper function to parse date range with defaults
func parseDateRange(startDateStr, endDateStr string) (time.Time, time.Time) {
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -30) // Default to last 30 days

	// Handle Z or lack of explicit time by using ParseInLocation and setting time to 00:00:00
	location, _ := time.LoadLocation("UTC")

	if startDateStr != "" {
		if parsed, err := time.ParseInLocation("2006-01-02", startDateStr, location); err == nil {
			startDate = parsed
		} else if parsed, err := time.Parse(time.RFC3339, startDateStr); err == nil {
			startDate = parsed.In(location)
		}
	}

	if endDateStr != "" {
		if parsed, err := time.ParseInLocation("2006-01-02", endDateStr, location); err == nil {
			endDate = parsed
		} else if parsed, err := time.Parse(time.RFC3339, endDateStr); err == nil {
			endDate = parsed.In(location)
		}
	}

	// Ensure endDate includes full day (up to 23:59:59.999...)
	// Only add if the date parsing was ambiguous (e.g., "YYYY-MM-DD")
	// If RFC3339 was used (like in the curl), it already has time info, but we apply the shift here defensively.
	endDate = endDate.Add(24 * time.Hour).Add(-1 * time.Second)

	return startDate, endDate
}
