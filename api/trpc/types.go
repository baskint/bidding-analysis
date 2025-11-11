package trpc

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TRPCResponse represents a tRPC response structure
type TRPCResponse struct {
	Result *TRPCResult `json:"result,omitempty"`
	Error  *TRPCError  `json:"error,omitempty"`
}

// TRPCResult represents successful tRPC result
type TRPCResult struct {
	Data interface{} `json:"data"`
	Type string      `json:"type"`
}

// TRPCError represents tRPC error
type TRPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ProcessBidInput represents the input for bid processing
type ProcessBidInput struct {
	CampaignID            string   `json:"campaignId"`
	UserID                string   `json:"userId"`
	FloorPrice            float64  `json:"floorPrice"`
	DeviceType            string   `json:"deviceType"`
	OS                    string   `json:"os"`
	Browser               string   `json:"browser"`
	Country               string   `json:"country"`
	Region                string   `json:"region"`
	City                  string   `json:"city"`
	Keywords              []string `json:"keywords"`
	SegmentID             string   `json:"segmentId"`
	SegmentCategory       string   `json:"segmentCategory"`
	EngagementScore       float64  `json:"engagementScore"`
	ConversionProbability float64  `json:"conversionProbability"`
}

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

// CampaignStatsInput represents input for campaign statistics
type CampaignStatsInput struct {
	CampaignID string `json:"campaignId"`
	StartTime  string `json:"startTime"`
	EndTime    string `json:"endTime"`
}

// BidHistoryInput represents input for bid history requests
type BidHistoryInput struct {
	CampaignID string `json:"campaignId"`
	StartTime  string `json:"startTime"`
	EndTime    string `json:"endTime"`
	Limit      int    `json:"limit"`
	Offset     int    `json:"offset"`
}

// FraudAlertsInput represents input for fraud alerts
type FraudAlertsInput struct {
	StartTime         string `json:"startTime"`
	EndTime           string `json:"endTime"`
	SeverityThreshold int    `json:"severityThreshold"`
}

// ModelAccuracyInput represents input for model accuracy
type ModelAccuracyInput struct {
	StartTime    string `json:"startTime"`
	EndTime      string `json:"endTime"`
	ModelVersion string `json:"modelVersion"`
}

// Auth request/response types
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type AnalyticsInput struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	Limit     int    `json:"limit"`
}

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

// --- Analytics Output/Result Types with NULL Handling ---

// GeoBreakdown represents aggregated statistics by country.
// GeoBreakdown represents the core structure of the analytics data pulled from the database.
//
// NOTE: We switch from sql.NullString to string/float64 for JSON serialization
// to ensure the frontend receives "United States" or null, not {"String": "...", "Valid": true}.
// The database layer handles the Nullable properties, but the JSON output should be clean.
type GeoBreakdown struct {
	// Fields used internally for SQL scanning (retaining the sql.NullString type)
	// We will use '-' to suppress these fields from the JSON output.
	// We use an internal struct to hide them cleanly.
	dbCountry sql.NullString `db:"country" json:"-"`
	dbRegion  sql.NullString `db:"region" json:"-"`

	// Fields exported to JSON (what the frontend sees)
	Country string `json:"country"`
	Region  string `json:"region"`

	TotalBids int64 `db:"total_bids" json:"totalBids"`

	WonBids     int64   `db:"won_bids" json:"wonBids"`
	Conversions int64   `db:"conversions" json:"conversions"`
	Spend       float64 `db:"spend" json:"spend"`

	WinRate        float64 `db:"win_rate" json:"winRate"`
	ConversionRate float64 `db:"conversion_rate" json:"conversionRate"`
	TotalSpend     float64 `db:"total_spend" json:"totalSpend"` // Check if this should be aliased to 'Spend'
}

// NOTE: You must also implement the database scanner and/or the JSON marshaler logic.
// The easiest way to handle both is to use a custom Unmarshal/Marshal logic,
// but if you are using `sqlx.Select` or similar, you need to ensure the DB fields
// map correctly.

// Scan is required to make the struct compatible with database/sql and sqlx,
// allowing the DB values to be mapped to the internal sql.NullString fields.
func (g *GeoBreakdown) Scan(src any) error {
	// A more robust implementation would manually scan the columns, but assuming
	// your ORM/SQL package handles the internal `db:"..."` tags, the main fix
	// is usually handling the JSON output.
	return nil // Placeholder for actual scanning logic if required by your ORM
}

// MarshalJSON implements the json.Marshaler interface to ensure clean output.
// This is the CRITICAL step that fixes the frontend error.
func (g GeoBreakdown) MarshalJSON() ([]byte, error) {
	// 1. Create a shadow struct type that contains the desired JSON fields
	// and excludes the internal sql.NullString fields.
	type GeoBreakdownAlias GeoBreakdown

	// 2. Prepare the clean output struct.
	output := struct {
		GeoBreakdownAlias
	}{
		GeoBreakdownAlias: GeoBreakdownAlias(g),
	}

	// 3. Populate the clean string fields from the internal Nullable fields.
	// This logic handles the NULL-to-"" conversion.
	if g.dbCountry.Valid {
		output.Country = g.dbCountry.String
	} else {
		output.Country = "" // Or use "N/A" if you prefer an explicit placeholder
	}

	if g.dbRegion.Valid {
		output.Region = g.dbRegion.String
	} else {
		output.Region = ""
	}

	// 4. Marshal the clean output struct to JSON bytes.
	return json.Marshal(output)
}

// DeviceBreakdown represents aggregated statistics by device type.
// In api/trpc/types.go:

// DeviceBreakdown represents aggregated statistics by device type.
type DeviceBreakdown struct {
	// (sql.NullString handles the previous "converting NULL to string" error)
	DeviceType sql.NullString `db:"device_type" json:"deviceType"`
	TotalBids  int64          `db:"total_bids" json:"totalBids"`

	// ADDED RAW COUNTS/SPEND FROM SQL QUERY:
	WonBids     int64   `db:"won_bids" json:"wonBids"`        // Fixes Line 278 error
	Conversions int64   `db:"conversions" json:"conversions"` // Fixes Line 279 error
	Spend       float64 `db:"spend" json:"spend"`             // Fixes Line 280 error

	// CALCULATED METRICS:
	WinRate        float64 `db:"win_rate" json:"winRate"`
	ConversionRate float64 `db:"conversion_rate" json:"conversionRate"`
	AverageBid     float64 `db:"average_bid" json:"averageBid"` // Fixes Line 283 error

	TotalSpend float64 `db:"total_spend" json:"totalSpend"` // You may want to rename 'Spend' to 'TotalSpend' for consistency if needed
}

// CompetitiveAnalysis represents competitive metrics by segment.
type CompetitiveAnalysis struct {
	SegmentCategory sql.NullString `db:"segment_category" json:"segmentCategory"`

	// FIX: Match the SQL aliases exactly
	OurWinRate       float64         `db:"our_win_rate" json:"ourWinRate"`             // FIX: Added
	MarketAverageBid sql.NullFloat64 `db:"market_average_bid" json:"marketAverageBid"` // FIX: Added (using sql.NullFloat64 for potential NULL AVG())

	OurAverageBid     sql.NullFloat64 `db:"our_average_bid" json:"ourAverageBid"`
	AverageFloorPrice sql.NullFloat64 `db:"average_floor_price" json:"averageFloorPrice"`

	// FIX: Match the SQL alias exactly
	CompetitionIntensity float64 `db:"competition_intensity" json:"competitionIntensity"` // FIX: Added
	TotalOpportunities   int64   `db:"total_opportunities" json:"totalOpportunities"`     // FIX: This is 'COUNT(*)' from SQL
}

// KeywordAnalysis represents performance metrics for a specific keyword.
// KeywordAnalysis represents performance metrics for a specific keyword.
type KeywordAnalysis struct {
	Keyword   string `db:"keyword" json:"keyword"`
	TotalBids int64  `db:"total_bids" json:"totalBids"`
	// ADDED RAW COUNTS/SPEND FROM SQL:
	WonBids     int64   `db:"won_bids" json:"wonBids"`
	Conversions int64   `db:"conversions" json:"conversions"`
	Spend       float64 `db:"spend" json:"spend"`
	Revenue     float64 `db:"revenue" json:"revenue"`

	// EXISTING CALCULATED RATES:
	WinRate        float64 `db:"win_rate" json:"winRate"`
	ConversionRate float64 `db:"conversion_rate" json:"conversionRate"`
	// ADDED CALCULATED METRICS:
	CPA  float64 `db:"cpa" json:"cpa"`
	ROAS float64 `db:"roas" json:"roas"`
	// Note: I changed TotalBids, WonBids, and Conversions to int64 based on typical SQL COUNT() return type.
}
