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

// MLPredictor wraps the ML predictor
type MLPredictor struct {
	predictor mlpredictor.Predictor
	bidStore  *store.BidStore
}

// NewMLPredictor creates a predictor using the ML service
func NewMLPredictor(mlServiceURL, encodersPath string, bidStore *store.BidStore) (*Predictor, error) {
	log.Println("🤖 Connecting to ML service...")
	
	mlPred, err := mlpredictor.NewBidPredictorHTTP(mlServiceURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ML service: %w", err)
	}

	log.Println("✅ ML service connected!")

	adapter := &MLPredictor{
		predictor: mlPred,
		bidStore:  bidStore,
	}

	return &Predictor{
		openaiClient: adapter,
		bidStore:     bidStore,
		modelVersion: "ml-http-v1.0.0",
	}, nil
}

// PredictBidPrice implements AIClient interface
func (m *MLPredictor) PredictBidPrice(ctx context.Context, req *models.BidRequest, historicalData []*models.BidEvent) (*models.BidResponse, error) {
	features := m.extractFeatures(req, historicalData)

	bidPrice, err := m.predictor.Predict(features)
	if err != nil {
		return nil, fmt.Errorf("ML prediction failed: %w", err)
	}

	return &models.BidResponse{
		BidPrice:     bidPrice,
		Confidence:   0.90,
		Strategy:     "ml_optimized",
		FraudRisk:    false,
		PredictionID: uuid.New().String(),
	}, nil
}

// AnalyzeAudienceSegment implements AIClient interface
func (m *MLPredictor) AnalyzeAudienceSegment(ctx context.Context, bidEvents []*models.BidEvent) (*AudienceAnalysis, error) {
	if len(bidEvents) == 0 {
		return &AudienceAnalysis{
			Segments: []string{},
			Insights: []string{"Insufficient data"},
		}, nil
	}

	var conversions int
	for _, event := range bidEvents {
		if event.Converted {
			conversions++
		}
	}

	return &AudienceAnalysis{
		Segments: []string{"standard"},
		Insights: []string{fmt.Sprintf("Analyzed %d events", len(bidEvents))},
	}, nil
}

// DetectFraud implements AIClient interface
func (m *MLPredictor) DetectFraud(ctx context.Context, bidEvents []*models.BidEvent) (*FraudAnalysis, error) {
	return &FraudAnalysis{
		FraudDetected: false,
		Confidence:    0.70,
		Patterns:      []string{"No fraud detected"},
		Severity:      0,
	}, nil
}

func (m *MLPredictor) extractFeatures(req *models.BidRequest, historicalData []*models.BidEvent) mlpredictor.BidFeatures {
	now := time.Now()
	stats := m.calculateHistoricalStats(historicalData)

	return mlpredictor.BidFeatures{
		FloorPrice:                req.FloorPrice,
		EngagementScore:           req.UserSegment.EngagementScore,
		ConversionProbability:     req.UserSegment.ConversionProbability,
		HistoricalWinRate:         stats.WinRate,
		HistoricalAvgBid:          stats.AvgBid,
		HistoricalAvgWinPrice:     stats.AvgWinPrice,
		DeviceType:                req.DeviceInfo.DeviceType,
		SegmentCategory:           req.UserSegment.Category,
		Country:                   req.GeoLocation.Country,
		HourOfDay:                 now.Hour(),
		DayOfWeek:                 int(now.Weekday()),
		CampaignSpendLast7d:       stats.SpendLast7d,
		CampaignConversionsLast7d: stats.ConversionsLast7d,
	}
}

type HistoricalStats struct {
	WinRate, AvgBid, AvgWinPrice, SpendLast7d, ConversionsLast7d float64
}

func (m *MLPredictor) calculateHistoricalStats(bidEvents []*models.BidEvent) HistoricalStats {
	if len(bidEvents) == 0 {
		return HistoricalStats{0.4, 2.5, 2.7, 100.0, 3.0}
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

	stats := HistoricalStats{SpendLast7d: totalSpend, ConversionsLast7d: totalConversions}
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
