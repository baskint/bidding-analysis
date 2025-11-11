package trpc

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// PerformanceMetrics represents overall performance metrics
type PerformanceMetrics struct {
	TotalBids       int64   `json:"total_bids"`
	WonBids         int64   `json:"won_bids"`
	Conversions     int64   `json:"conversions"`
	TotalSpend      float64 `json:"total_spend"`
	Revenue         float64 `json:"revenue"`
	WinRate         float64 `json:"win_rate"`
	ConversionRate  float64 `json:"conversion_rate"`
	AverageBid      float64 `json:"average_bid"`
	CPA             float64 `json:"cpa"`  // Cost Per Acquisition
	ROAS            float64 `json:"roas"` // Return On Ad Spend
	FraudDetections int64   `json:"fraud_detections"`
}

// KeywordAnalysis represents keyword performance data
type KeywordAnalysis struct {
	Keyword        string  `json:"keyword"`
	TotalBids      int64   `json:"total_bids"`
	WonBids        int64   `json:"won_bids"`
	Conversions    int64   `json:"conversions"`
	Spend          float64 `json:"spend"`
	Revenue        float64 `json:"revenue"`
	WinRate        float64 `json:"win_rate"`
	ConversionRate float64 `json:"conversion_rate"`
	CPA            float64 `json:"cpa"`
	ROAS           float64 `json:"roas"`
}

// DeviceBreakdown represents device-specific performance
type DeviceBreakdown struct {
	DeviceType     string  `json:"device_type"`
	TotalBids      int64   `json:"total_bids"`
	WonBids        int64   `json:"won_bids"`
	Conversions    int64   `json:"conversions"`
	Spend          float64 `json:"spend"`
	WinRate        float64 `json:"win_rate"`
	ConversionRate float64 `json:"conversion_rate"`
	AverageBid     float64 `json:"average_bid"`
}

// GeoBreakdown represents geographic performance
type GeoBreakdown struct {
	Country        string  `json:"country"`
	Region         string  `json:"region"`
	TotalBids      int64   `json:"total_bids"`
	WonBids        int64   `json:"won_bids"`
	Conversions    int64   `json:"conversions"`
	Spend          float64 `json:"spend"`
	WinRate        float64 `json:"win_rate"`
	ConversionRate float64 `json:"conversion_rate"`
}

// HourlyPerformance represents time-based performance
type HourlyPerformance struct {
	Hour           int     `json:"hour"`
	TotalBids      int64   `json:"total_bids"`
	WonBids        int64   `json:"won_bids"`
	Conversions    int64   `json:"conversions"`
	Spend          float64 `json:"spend"`
	WinRate        float64 `json:"win_rate"`
	ConversionRate float64 `json:"conversion_rate"`
	AverageBid     float64 `json:"average_bid"`
}

// DailyTrend represents daily performance trends
type DailyTrend struct {
	Date           string  `json:"date"`
	TotalBids      int64   `json:"total_bids"`
	WonBids        int64   `json:"won_bids"`
	Conversions    int64   `json:"conversions"`
	Spend          float64 `json:"spend"`
	Revenue        float64 `json:"revenue"`
	WinRate        float64 `json:"win_rate"`
	ConversionRate float64 `json:"conversion_rate"`
	CPA            float64 `json:"cpa"`
}

// CompetitiveAnalysis represents competitive insights
type CompetitiveAnalysis struct {
	SegmentCategory      string  `json:"segment_category"`
	OurWinRate           float64 `json:"our_win_rate"`
	MarketAverageBid     float64 `json:"market_average_bid"`
	OurAverageBid        float64 `json:"our_average_bid"`
	AverageFloorPrice    float64 `json:"average_floor_price"`
	CompetitionIntensity float64 `json:"competition_intensity"` // Bid/Floor ratio
	TotalOpportunities   int64   `json:"total_opportunities"`
}

// getPerformanceOverview returns overall performance metrics
func (h *Handler) getPerformanceOverview(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r.Context())
	if userID == "" {
		h.writeErrorResponse(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	// Parse date range from query or use defaults
	var req struct {
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	startDate, endDate := parseDateRange(req.StartDate, req.EndDate)

	// Query to get overall performance metrics
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
	err := h.bidStore.DB().QueryRow(query, userID, startDate, endDate).Scan(
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
		log.Printf("Failed to get performance overview: %v", err)
		h.writeErrorResponse(w, "Failed to retrieve performance metrics", http.StatusInternalServerError)
		return
	}

	// Calculate derived metrics
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
	userID := getUserIDFromContext(r.Context())
	if userID == "" {
		log.Println("DEBUG: UserID is empty, unauthorized request.")
		h.writeErrorResponse(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	// We do NOT parse to UUID here because bid_events.user_id is VARCHAR,
	// and passing a Go UUID type to a VARCHAR column caused the 'syntax error at or near ".."'
	log.Printf("DEBUG: Processing getKeywordAnalysis for UserID: %s", userID)

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

	log.Printf("DEBUG: Params - StartDate: %v, EndDate: %v, Limit: %d", startDate, endDate, limit)

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

	// Pass the string userID directly. $1=userID, $2=startDate, $3=endDate, $4=limit
	rows, err := h.bidStore.DB().Query(query, userID, startDate, endDate, limit)
	if err != nil {
		log.Printf("DEBUG ERROR: Failed to execute keyword query: %v", err)
		h.writeErrorResponse(w, "Failed to retrieve keyword analysis", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var keywords []KeywordAnalysis
	rowsScanned := 0
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
			log.Printf("DEBUG ERROR: Failed to scan row for keyword analysis: %v", err)
			continue
		}
		keywords = append(keywords, kw)
		rowsScanned++
	}

	log.Printf("DEBUG: Query completed. Successfully scanned %d rows.", rowsScanned)

	if err := rows.Err(); err != nil {
		log.Printf("DEBUG ERROR: Error after row iteration: %v", err)
	}

	h.writeTRPCResponse(w, keywords)
}

// getDeviceBreakdown returns device-specific performance
func (h *Handler) getDeviceBreakdown(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r.Context())
	if userID == "" {
		log.Println("DEBUG: DeviceBreakdown: UserID is empty, unauthorized request.")
		h.writeErrorResponse(w, "User not found in context", http.StatusUnauthorized)
		return
	}
	log.Printf("DEBUG: DeviceBreakdown: Processing for UserID: %s", userID)

	var req struct {
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	startDate, endDate := parseDateRange(req.StartDate, req.EndDate)
	log.Printf("DEBUG: DeviceBreakdown: Params - StartDate: %v, EndDate: %v", startDate, endDate)

	query := `
        SELECT 
            device_type,
            COUNT(*) as total_bids,
            SUM(CASE WHEN won THEN 1 ELSE 0 END) as won_bids,
            SUM(CASE WHEN converted THEN 1 ELSE 0 END) as conversions,
            SUM(CASE WHEN won THEN bid_price ELSE 0 END) as spend,
            -- ... (rest of calculated metrics)
            CASE WHEN COUNT(*) > 0 THEN CAST(SUM(CASE WHEN won THEN 1 ELSE 0 END) AS FLOAT) / COUNT(*) ELSE 0 END as win_rate,
            CASE WHEN SUM(CASE WHEN won THEN 1 ELSE 0 END) > 0 THEN CAST(SUM(CASE WHEN converted THEN 1 ELSE 0 END) AS FLOAT) / SUM(CASE WHEN won THEN 1 ELSE 0 END) ELSE 0 END as conversion_rate,
            CASE WHEN COUNT(*) > 0 THEN SUM(bid_price) / COUNT(*) ELSE 0 END as average_bid
        FROM bid_events
        WHERE user_id = $1 
            AND timestamp BETWEEN $2 AND $3
        GROUP BY device_type
        ORDER BY total_bids DESC
    `

	rows, err := h.bidStore.DB().Query(query, userID, startDate, endDate)
	if err != nil {
		log.Printf("DEBUG ERROR: Failed to get device breakdown: %v", err)
		h.writeErrorResponse(w, "Failed to retrieve device breakdown", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var devices []DeviceBreakdown
	rowsScanned := 0
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
			log.Printf("DEBUG ERROR: DeviceBreakdown: Failed to scan row: %v", err)
			continue
		}
		devices = append(devices, device)
		rowsScanned++
	}

	log.Printf("DEBUG: DeviceBreakdown: Query completed. Successfully scanned %d rows.", rowsScanned)

	h.writeTRPCResponse(w, devices)
}

// getGeoBreakdown returns geographic performance
// getGeoBreakdown returns geographic performance
func (h *Handler) getGeoBreakdown(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r.Context())
	if userID == "" {
		log.Println("DEBUG: GeoBreakdown: UserID is empty, unauthorized request.")
		h.writeErrorResponse(w, "User not found in context", http.StatusUnauthorized)
		return
	}
	log.Printf("DEBUG: GeoBreakdown: Processing for UserID: %s", userID)

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
	log.Printf("DEBUG: GeoBreakdown: Params - StartDate: %v, EndDate: %v, Limit: %d", startDate, endDate, limit)

	query := `
		SELECT 
			country,
			region,
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

	// Pass the string userID directly for $1.
	rows, err := h.bidStore.DB().Query(query, userID, startDate, endDate, limit)
	if err != nil {
		log.Printf("DEBUG ERROR: Failed to get geo breakdown: %v", err)
		h.writeErrorResponse(w, "Failed to retrieve geographic breakdown", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var geos []GeoBreakdown
	rowsScanned := 0
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
			log.Printf("DEBUG ERROR: GeoBreakdown: Failed to scan row: %v", err)
			continue
		}
		geos = append(geos, geo)
		rowsScanned++
	}

	log.Printf("DEBUG: GeoBreakdown: Query completed. Successfully scanned %d rows.", rowsScanned)

	h.writeTRPCResponse(w, geos)
}

// getHourlyPerformance returns time-based performance by hour
func (h *Handler) getHourlyPerformance(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r.Context())
	if userID == "" {
		h.writeErrorResponse(w, "User not found in context", http.StatusUnauthorized)
		return
	}

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

	rows, err := h.bidStore.DB().Query(query, userID, startDate, endDate)
	if err != nil {
		log.Printf("Failed to get hourly performance: %v", err)
		h.writeErrorResponse(w, "Failed to retrieve hourly performance", http.StatusInternalServerError)
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
			continue
		}
		h.Hour = int(hourFloat)
		hourly = append(hourly, h)
	}

	h.writeTRPCResponse(w, hourly)
}

// getDailyTrends returns daily performance trends
func (h *Handler) getDailyTrends(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r.Context())
	if userID == "" {
		h.writeErrorResponse(w, "User not found in context", http.StatusUnauthorized)
		return
	}

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

	rows, err := h.bidStore.DB().Query(query, userID, startDate, endDate)
	if err != nil {
		log.Printf("Failed to get daily trends: %v", err)
		h.writeErrorResponse(w, "Failed to retrieve daily trends", http.StatusInternalServerError)
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
			continue
		}
		trend.Date = date.Format("2006-01-02")
		trends = append(trends, trend)
	}

	h.writeTRPCResponse(w, trends)
}

// getCompetitiveAnalysis returns competitive insights
func (h *Handler) getCompetitiveAnalysis(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r.Context())
	if userID == "" {
		log.Println("DEBUG: CompetitiveAnalysis: UserID is empty, unauthorized request.")
		h.writeErrorResponse(w, "User not found in context", http.StatusUnauthorized)
		return
	}
	log.Printf("DEBUG: CompetitiveAnalysis: Processing for UserID: %s", userID)

	var req struct {
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	startDate, endDate := parseDateRange(req.StartDate, req.EndDate)
	log.Printf("DEBUG: CompetitiveAnalysis: Params - StartDate: %v, EndDate: %v", startDate, endDate)

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
        -- NOTE: Since this is competitive analysis, the WHERE clause filters by time, 
        -- but does NOT filter by user_id, allowing all bids (market data) to be included.
		GROUP BY segment_category
		ORDER BY total_opportunities DESC
		LIMIT 10
	`

	// $1 = userID (string), $2 = startDate, $3 = endDate
	rows, err := h.bidStore.DB().Query(query, userID, startDate, endDate)
	if err != nil {
		log.Printf("DEBUG ERROR: Failed to get competitive analysis: %v", err)
		h.writeErrorResponse(w, "Failed to retrieve competitive analysis", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var competitive []CompetitiveAnalysis
	rowsScanned := 0
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
			log.Printf("DEBUG ERROR: CompetitiveAnalysis: Failed to scan row: %v", err)
			continue
		}
		competitive = append(competitive, comp)
		rowsScanned++
	}

	log.Printf("DEBUG: CompetitiveAnalysis: Query completed. Successfully scanned %d rows.", rowsScanned)

	h.writeTRPCResponse(w, competitive)
}

// getCampaignComparison returns comparison metrics across campaigns
func (h *Handler) getCampaignComparison(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r.Context())
	if userID == "" {
		h.writeErrorResponse(w, "User not found in context", http.StatusUnauthorized)
		return
	}

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
		campaignsFilter = " AND campaign_id = ANY($4)"
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
	var err error

	userUUID, _ := uuid.Parse(userID)
	if len(req.CampaignIDs) > 0 {
		rows, err = h.bidStore.DB().Query(query, userUUID, startDate, endDate, req.CampaignIDs)
	} else {
		rows, err = h.bidStore.DB().Query(query, userUUID, startDate, endDate)
	}

	if err != nil {
		log.Printf("Failed to get campaign comparison: %v", err)
		h.writeErrorResponse(w, "Failed to retrieve campaign comparison", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type CampaignComparison struct {
		CampaignID     string  `json:"campaign_id"`
		CampaignName   string  `json:"campaign_name"`
		TotalBids      int64   `json:"total_bids"`
		WonBids        int64   `json:"won_bids"`
		Conversions    int64   `json:"conversions"`
		Spend          float64 `json:"spend"`
		WinRate        float64 `json:"win_rate"`
		ConversionRate float64 `json:"conversion_rate"`
	}

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

	if startDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = parsed
		}
	}

	if endDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", endDateStr); err == nil {
			endDate = parsed
		}
	}

	// Ensure endDate includes full day
	endDate = endDate.Add(24 * time.Hour).Add(-1 * time.Second)

	return startDate, endDate
}
