package ml

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/baskint/bidding-analysis/internal/mlpredictor"
	"github.com/baskint/bidding-analysis/internal/models"
	"github.com/baskint/bidding-analysis/internal/store"
	"github.com/google/uuid"
)

// MLPredictor wraps the ML predictor to work with existing code
type MLPredictor struct {
	predictor mlpredictor.Predictor
	bidStore  *store.BidStore
}

// NewMLPredictor creates a predictor using the ML model (ONNX)
func NewMLPredictor(modelPath, encodersPath string, bidStore *store.BidStore) (*Predictor, error) {
	// Try ONNX predictor (pure Go, no Python needed!)
	log.Println("🤖 Loading ONNX model...")

	mlPred, err := mlpredictor.NewBidPredictorONNX(modelPath, encodersPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load ONNX model: %w", err)
	}

	log.Println("✅ ONNX model loaded successfully!")

	// Wrap in adapter
	adapter := &MLPredictor{
		predictor: mlPred,
		bidStore:  bidStore,
	}

	return &Predictor{
		openaiClient: adapter,
		bidStore:     bidStore,
		modelVersion: "ml-onnx-v1.0.0",
	}, nil
}

// PredictBidPrice implements AIClient interface for ML model
func (m *MLPredictor) PredictBidPrice(ctx context.Context, req *models.BidRequest, historicalData []*models.BidEvent) (*models.BidResponse, error) {
	// Convert BidRequest to mlpredictor.BidFeatures
	features := m.extractFeatures(req, historicalData)

	// Get ML prediction
	bidPrice, err := m.predictor.Predict(features)
	if err != nil {
		return nil, fmt.Errorf("ML prediction failed: %w", err)
	}

	// Build response
	response := &models.BidResponse{
		BidPrice:     bidPrice,
		Confidence:   0.90,
		Strategy:     "ml_optimized",
		FraudRisk:    false,
		PredictionID: uuid.New().String(),
	}

	return response, nil
}

// AnalyzeAudienceSegment implements AIClient interface
func (m *MLPredictor) AnalyzeAudienceSegment(ctx context.Context, bidEvents []*models.BidEvent) (*AudienceAnalysis, error) {
	if len(bidEvents) == 0 {
		return &AudienceAnalysis{
			Segments: []string{},
			Insights: []string{"Insufficient data for analysis"},
		}, nil
	}

	var conversions int
	deviceTypes := make(map[string]int)
	segments := make(map[string]bool)

	for _, event := range bidEvents {
		if event.Converted {
			conversions++
		}
		deviceTypes[event.DeviceType]++
		if event.SegmentCategory != "" {
			segments[event.SegmentCategory] = true
		}
	}

	conversionRate := float64(conversions) / float64(len(bidEvents))
	insights := []string{
		fmt.Sprintf("Analyzed %d bid events", len(bidEvents)),
		fmt.Sprintf("Conversion rate: %.2f%%", conversionRate*100),
	}

	var topDevice string
	var maxCount int
	for device, count := range deviceTypes {
		if count > maxCount {
			topDevice = device
			maxCount = count
		}
	}
	if topDevice != "" {
		insights = append(insights, fmt.Sprintf("Top performing device: %s", topDevice))
	}

	segmentList := make([]string, 0, len(segments))
	for segment := range segments {
		segmentList = append(segmentList, segment)
	}

	return &AudienceAnalysis{
		Segments: segmentList,
		Insights: insights,
	}, nil
}

// DetectFraud implements AIClient interface
func (m *MLPredictor) DetectFraud(ctx context.Context, bidEvents []*models.BidEvent) (*FraudAnalysis, error) {
	if len(bidEvents) < 10 {
		return &FraudAnalysis{
			FraudDetected: false,
			Confidence:    0.5,
			Patterns:      []string{"Insufficient data for fraud analysis"},
			Severity:      0,
		}, nil
	}

	userActivity := make(map[uuid.UUID]int)
	userConversions := make(map[uuid.UUID]int)

	for _, event := range bidEvents {
		userActivity[event.UserID]++
		if event.Converted {
			userConversions[event.UserID]++
		}
	}

	var suspiciousCount int
	var patterns []string

	for userID, count := range userActivity {
		if count > 50 {
			suspiciousCount++
			patterns = append(patterns, fmt.Sprintf("User %s: %d bids (abnormally high)", userID, count))
		}

		conversions := userConversions[userID]
		if conversions > 0 {
			convRate := float64(conversions) / float64(count)
			if convRate > 0.8 && count > 10 {
				patterns = append(patterns, fmt.Sprintf("User %s: %.0f%% conversion rate (suspicious)", userID, convRate*100))
			}
		}
	}

	fraudDetected := suspiciousCount > 0
	severity := suspiciousCount
	if severity > 10 {
		severity = 10
	}

	if !fraudDetected {
		patterns = []string{"No obvious fraud patterns detected"}
	}

	return &FraudAnalysis{
		FraudDetected: fraudDetected,
		Confidence:    0.70,
		Patterns:      patterns,
		Severity:      severity,
	}, nil
}

// extractFeatures converts your BidRequest to mlpredictor.BidFeatures
func (m *MLPredictor) extractFeatures(req *models.BidRequest, historicalData []*models.BidEvent) mlpredictor.BidFeatures {
	now := time.Now()
	stats := m.calculateHistoricalStats(historicalData)

	floorPrice := req.FloorPrice
	engagementScore := 0.5
	if req.UserSegment.EngagementScore > 0 {
		engagementScore = req.UserSegment.EngagementScore
	}

	conversionProb := 0.1
	if req.UserSegment.ConversionProbability > 0 {
		conversionProb = req.UserSegment.ConversionProbability
	}

	deviceType := "unknown"
	if req.DeviceInfo.DeviceType != "" {
		deviceType = req.DeviceInfo.DeviceType
	}

	segmentCategory := "standard"
	if req.UserSegment.Category != "" {
		segmentCategory = req.UserSegment.Category
	}

	country := "US"
	if req.GeoLocation.Country != "" {
		country = req.GeoLocation.Country
	}

	return mlpredictor.BidFeatures{
		FloorPrice:                floorPrice,
		EngagementScore:           engagementScore,
		ConversionProbability:     conversionProb,
		HistoricalWinRate:         stats.WinRate,
		HistoricalAvgBid:          stats.AvgBid,
		HistoricalAvgWinPrice:     stats.AvgWinPrice,
		DeviceType:                deviceType,
		SegmentCategory:           segmentCategory,
		Country:                   country,
		HourOfDay:                 now.Hour(),
		DayOfWeek:                 int(now.Weekday()),
		CampaignSpendLast7d:       stats.SpendLast7d,
		CampaignConversionsLast7d: stats.ConversionsLast7d,
	}
}

type HistoricalStats struct {
	WinRate           float64
	AvgBid            float64
	AvgWinPrice       float64
	SpendLast7d       float64
	ConversionsLast7d float64
}

func (m *MLPredictor) calculateHistoricalStats(bidEvents []*models.BidEvent) HistoricalStats {
	if len(bidEvents) == 0 {
		return HistoricalStats{
			WinRate:           0.4,
			AvgBid:            2.5,
			AvgWinPrice:       2.7,
			SpendLast7d:       100.0,
			ConversionsLast7d: 3.0,
		}
	}

	var totalBids, totalWins, totalBidAmount, totalWinAmount, totalSpend, totalConversions float64

	for _, event := range bidEvents {
		totalBids++
		totalBidAmount += event.BidPrice

		if event.Won {
			totalWins++
			if event.WinPrice != nil {
				totalWinAmount += *event.WinPrice
				totalSpend += *event.WinPrice
			}
		}

		if event.Converted {
			totalConversions++
		}
	}

	stats := HistoricalStats{
		SpendLast7d:       totalSpend,
		ConversionsLast7d: totalConversions,
	}

	if totalBids > 0 {
		stats.WinRate = totalWins / totalBids
		stats.AvgBid = totalBidAmount / totalBids
	}

	if totalWins > 0 {
		stats.AvgWinPrice = totalWinAmount / totalWins
	}

	return stats
}

func (m *MLPredictor) Close() error {
	if m.predictor != nil {
		return m.predictor.Close()
	}
	return nil
}
