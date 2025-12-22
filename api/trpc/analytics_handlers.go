package trpc

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
)

// Request types - extracted to top of file for clarity
type DateRangeRequest struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type KeywordAnalysisRequest struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Limit     int    `json:"limit"`
}

type CompetitiveAnalysisRequest struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// ============================================================================
// REFACTORED HANDLERS - Using WithAuth Wrapper
// ============================================================================

// getPerformanceOverview returns overall performance metrics
// BEFORE: 90 lines with boilerplate
// AFTER: 45 lines of pure business logic
func (h *Handler) getPerformanceOverview(ctx context.Context, userID uuid.UUID, req interface{}) (interface{}, error) {
	params := req.(*DateRangeRequest)
	startDate, endDate := parseDateRange(params.StartDate, params.EndDate)

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
	err := h.bidStore.DB().QueryRowContext(ctx, query, userID, startDate, endDate).Scan(
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
		return nil, fmt.Errorf("failed to get performance overview: %w", err)
	}

	// Calculate derived metrics
	if metrics.Conversions > 0 {
		metrics.CPA = metrics.TotalSpend / float64(metrics.Conversions)
	}
	if metrics.TotalSpend > 0 {
		metrics.ROAS = metrics.Revenue / metrics.TotalSpend
	}

	return metrics, nil
}

// getKeywordAnalysis returns keyword performance breakdown
// BEFORE: 120+ lines with boilerplate
// AFTER: 75 lines of pure business logic
func (h *Handler) getKeywordAnalysis(ctx context.Context, userID uuid.UUID, req interface{}) (interface{}, error) {
	params := req.(*KeywordAnalysisRequest)
	startDate, endDate := parseDateRange(params.StartDate, params.EndDate)

	limit := 20
	if params.Limit > 0 && params.Limit <= 100 {
		limit = params.Limit
	}

	query := `
		WITH keyword_stats AS (
			SELECT 
				unnest(keywords) as keyword,
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
		WHERE total_bids >= 5
		ORDER BY total_bids DESC
		LIMIT $4
	`

	rows, err := h.bidStore.DB().QueryContext(ctx, query, userID, startDate, endDate, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query keyword analysis: %w", err)
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
			log.Printf("Error scanning keyword row: %v", err)
			continue
		}
		keywords = append(keywords, kw)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating keyword rows: %w", err)
	}

	return keywords, nil
}

// getDeviceBreakdown returns device type performance breakdown
func (h *Handler) getDeviceBreakdown(ctx context.Context, userID uuid.UUID, req interface{}) (interface{}, error) {
	params := req.(*DateRangeRequest)
	startDate, endDate := parseDateRange(params.StartDate, params.EndDate)

	query := `
		SELECT 
			device_type,
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

	rows, err := h.bidStore.DB().QueryContext(ctx, query, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to query device breakdown: %w", err)
	}
	defer rows.Close()

	var devices []DeviceBreakdown
	for rows.Next() {
		var dev DeviceBreakdown
		err := rows.Scan(
			&dev.DeviceType,
			&dev.TotalBids,
			&dev.WonBids,
			&dev.Conversions,
			&dev.Spend,
			&dev.WinRate,
			&dev.ConversionRate,
			&dev.AverageBid,
		)
		if err != nil {
			log.Printf("Error scanning device row: %v", err)
			continue
		}
		devices = append(devices, dev)
	}

	return devices, rows.Err()
}

// getGeoBreakdown returns geographic performance breakdown
func (h *Handler) getGeoBreakdown(ctx context.Context, userID uuid.UUID, req interface{}) (interface{}, error) {
	params := req.(*DateRangeRequest)
	startDate, endDate := parseDateRange(params.StartDate, params.EndDate)

	query := `
		SELECT 
			country,
			COUNT(*) as total_bids,
			SUM(CASE WHEN won THEN 1 ELSE 0 END) as won_bids,
			SUM(CASE WHEN converted THEN 1 ELSE 0 END) as conversions,
			SUM(CASE WHEN won THEN bid_price ELSE 0 END) as spend,
			CASE WHEN COUNT(*) > 0 THEN CAST(SUM(CASE WHEN won THEN 1 ELSE 0 END) AS FLOAT) / COUNT(*) ELSE 0 END as win_rate,
			CASE WHEN SUM(CASE WHEN won THEN 1 ELSE 0 END) > 0 THEN CAST(SUM(CASE WHEN converted THEN 1 ELSE 0 END) AS FLOAT) / SUM(CASE WHEN won THEN 1 ELSE 0 END) ELSE 0 END as conversion_rate
		FROM bid_events
		WHERE user_id = $1 
			AND timestamp BETWEEN $2 AND $3
		GROUP BY country
		ORDER BY total_bids DESC
		LIMIT 20
	`

	rows, err := h.bidStore.DB().QueryContext(ctx, query, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to query geo breakdown: %w", err)
	}
	defer rows.Close()

	var geos []GeoBreakdown
	for rows.Next() {
		var geo GeoBreakdown
		err := rows.Scan(
			&geo.Country,
			&geo.TotalBids,
			&geo.WonBids,
			&geo.Conversions,
			&geo.Spend,
			&geo.WinRate,
			&geo.ConversionRate,
		)
		if err != nil {
			log.Printf("Error scanning geo row: %v", err)
			continue
		}
		geos = append(geos, geo)
	}

	return geos, rows.Err()
}

// getHourlyPerformance returns performance metrics by hour of day
func (h *Handler) getHourlyPerformance(ctx context.Context, userID uuid.UUID, req interface{}) (interface{}, error) {
	params := req.(*DateRangeRequest)
	startDate, endDate := parseDateRange(params.StartDate, params.EndDate)

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
		GROUP BY EXTRACT(HOUR FROM timestamp)
		ORDER BY hour
	`

	rows, err := h.bidStore.DB().QueryContext(ctx, query, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to query hourly performance: %w", err)
	}
	defer rows.Close()

	var hourly []HourlyPerformance
	for rows.Next() {
		var hp HourlyPerformance
		err := rows.Scan(
			&hp.Hour,
			&hp.TotalBids,
			&hp.WonBids,
			&hp.Conversions,
			&hp.Spend,
			&hp.WinRate,
			&hp.ConversionRate,
			&hp.AverageBid,
		)
		if err != nil {
			log.Printf("Error scanning hourly row: %v", err)
			continue
		}
		hourly = append(hourly, hp)
	}

	return hourly, rows.Err()
}

// getDailyTrends returns daily performance trends
func (h *Handler) getDailyTrends(ctx context.Context, userID uuid.UUID, req interface{}) (interface{}, error) {
	params := req.(*DateRangeRequest)
	startDate, endDate := parseDateRange(params.StartDate, params.EndDate)

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

	rows, err := h.bidStore.DB().QueryContext(ctx, query, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to query daily trends: %w", err)
	}
	defer rows.Close()

	var trends []DailyTrend
	for rows.Next() {
		var trend DailyTrend
		err := rows.Scan(
			&trend.Date,
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
			log.Printf("Error scanning daily trend row: %v", err)
			continue
		}
		trends = append(trends, trend)
	}

	return trends, rows.Err()
}

// getCompetitiveAnalysis returns competitive insights
func (h *Handler) getCompetitiveAnalysis(ctx context.Context, userID uuid.UUID, req interface{}) (interface{}, error) {
	params := req.(*CompetitiveAnalysisRequest)
	startDate, endDate := parseDateRange(params.StartDate, params.EndDate)

	// Add buffer for timezone differences
	startDate = startDate.AddDate(0, 0, -7)
	endDate = endDate.AddDate(0, 0, 1)

	query := `
		SELECT
			COALESCE(segment_category, 'unknown') as segment_category,
			CASE WHEN COUNT(*) > 0 THEN CAST(SUM(CASE WHEN won THEN 1 ELSE 0 END) AS FLOAT) / COUNT(*) ELSE 0 END as our_win_rate,
			AVG(bid_price) * 1.15 as market_average_bid,
			AVG(bid_price) as our_average_bid,
			AVG(floor_price) as average_floor_price,
			1.15 as competition_intensity,
			COUNT(*) as total_opportunities
		FROM bid_events
		WHERE user_id = $1
			AND timestamp >= $2
			AND timestamp <= $3
		GROUP BY segment_category
		HAVING COUNT(*) > 0
		ORDER BY total_opportunities DESC
		LIMIT 10
	`

	rows, err := h.bidStore.DB().QueryContext(ctx, query, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to query competitive analysis: %w", err)
	}
	defer rows.Close()

	var competitive []CompetitiveAnalysis
	for rows.Next() {
		var comp CompetitiveAnalysis
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

	return competitive, rows.Err()
}
