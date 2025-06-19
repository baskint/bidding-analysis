package ml

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/baskint/bidding-analysis/internal/models"
	"github.com/baskint/bidding-analysis/internal/store"
)

// AIClient interface for AI prediction services
type AIClient interface {
	PredictBidPrice(ctx context.Context, req *models.BidRequest, historicalData []*models.BidEvent) (*models.BidResponse, error)
	AnalyzeAudienceSegment(ctx context.Context, bidEvents []*models.BidEvent) (*AudienceAnalysis, error)
	DetectFraud(ctx context.Context, bidEvents []*models.BidEvent) (*FraudAnalysis, error)
}

// Predictor handles ML predictions for bidding
type Predictor struct {
	openaiClient AIClient
	bidStore     *store.BidStore
	modelVersion string
}

// NewPredictor creates a new ML predictor
func NewPredictor(apiKey string, bidStore *store.BidStore) *Predictor {
	var openaiClient AIClient

	if apiKey == "" || apiKey == "your_openai_key" {
		// Use mock client if no API key provided
		openaiClient = NewMockOpenAIClient()
	} else {
		// Use real OpenAI client
		openaiClient = NewOpenAIClient(apiKey)
	}

	return &Predictor{
		openaiClient: openaiClient,
		bidStore:     bidStore,
		modelVersion: "v1.0.0",
	}
}

// PredictOptimalBid predicts the optimal bid price for a request
func (p *Predictor) PredictOptimalBid(ctx context.Context, req *models.BidRequest) (*models.BidResponse, error) {
	// Get historical data for this campaign
	historicalData, err := p.bidStore.GetRecentBids(req.CampaignID.String(), 100)
	if err != nil {
		return nil, fmt.Errorf("failed to get historical data: %w", err)
	}

	// Use OpenAI for prediction
	prediction, err := p.openaiClient.PredictBidPrice(ctx, req, historicalData)
	if err != nil {
		// Fallback to rule-based prediction if OpenAI fails
		return p.fallbackPrediction(req, historicalData), nil
	}

	// Add prediction ID for tracking
	prediction.PredictionID = uuid.New().String()

	// Validate and adjust prediction
	prediction = p.validatePrediction(prediction, req)

	return prediction, nil
}

// AnalyzeCampaignPerformance provides campaign performance insights
func (p *Predictor) AnalyzeCampaignPerformance(ctx context.Context, campaignID uuid.UUID, days int) (*CampaignAnalysis, error) {
	startTime := time.Now().AddDate(0, 0, -days)
	endTime := time.Now()

	// Get recent bid data
	bidData, err := p.bidStore.GetBidHistory(campaignID.String(), startTime, endTime, 1000, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get bid history: %w", err)
	}

	if len(bidData) == 0 {
		return &CampaignAnalysis{
			CampaignID: campaignID,
			Message:    "No bid data available for analysis",
		}, nil
	}

	// Calculate basic metrics
	analysis := &CampaignAnalysis{
		CampaignID: campaignID,
		Period:     fmt.Sprintf("%d days", days),
		TotalBids:  len(bidData),
	}

	var wonBids, conversions int
	var totalSpend, totalWinPrice float64
	deviceTypes := make(map[string]int)
	geoRegions := make(map[string]int)

	for _, bid := range bidData {
		totalSpend += bid.BidPrice

		if bid.Won {
			wonBids++
			if bid.WinPrice != nil {
				totalWinPrice += *bid.WinPrice
			}
		}

		if bid.Converted {
			conversions++
		}

		// Track device types
		deviceTypes[bid.DeviceType]++

		// Track geo regions
		region := fmt.Sprintf("%s, %s", bid.Region, bid.Country)
		geoRegions[region]++
	}

	// Calculate rates
	analysis.WinRate = float64(wonBids) / float64(len(bidData))
	if wonBids > 0 {
		analysis.ConversionRate = float64(conversions) / float64(wonBids)
		analysis.AvgWinPrice = totalWinPrice / float64(wonBids)
	}
	analysis.AvgBidPrice = totalSpend / float64(len(bidData))
	analysis.TotalSpend = totalSpend

	// Generate insights using OpenAI
	insights, err := p.generateInsights(ctx, bidData)
	if err == nil {
		analysis.Insights = insights
	}

	// Generate recommendations
	analysis.Recommendations = p.generateRecommendations(analysis)

	return analysis, nil
}

// DetectFraudPatterns analyzes bid data for fraud patterns
func (p *Predictor) DetectFraudPatterns(ctx context.Context, campaignID uuid.UUID) (*FraudDetectionResult, error) {
	// Get recent data for analysis
	startTime := time.Now().AddDate(0, 0, -7) // Last 7 days
	endTime := time.Now()

	bidData, err := p.bidStore.GetBidHistory(campaignID.String(), startTime, endTime, 500, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get bid data: %w", err)
	}

	if len(bidData) < 10 {
		return &FraudDetectionResult{
			CampaignID:    campaignID,
			FraudDetected: false,
			Message:       "Insufficient data for fraud analysis",
		}, nil
	}

	// Rule-based fraud detection
	ruleBasedResult := p.ruleBasedFraudDetection(bidData)

	// AI-powered fraud detection
	aiResult, err := p.openaiClient.DetectFraud(ctx, bidData)
	if err != nil {
		// Use only rule-based if AI fails
		return ruleBasedResult, nil
	}

	// Combine results
	return p.combineFraudResults(ruleBasedResult, aiResult), nil
}

// fallbackPrediction provides rule-based prediction when AI is unavailable
func (p *Predictor) fallbackPrediction(req *models.BidRequest, historical []*models.BidEvent) *models.BidResponse {
	baseBid := req.FloorPrice * 1.2 // Start 20% above floor

	// Adjust based on user segment
	if req.UserSegment.ConversionProbability > 0.1 {
		baseBid *= (1 + req.UserSegment.ConversionProbability)
	}

	// Adjust based on engagement score
	if req.UserSegment.EngagementScore > 0.7 {
		baseBid *= 1.15
	}

	// Adjust based on historical performance
	if len(historical) > 0 {
		var avgWinPrice float64
		var winCount int

		for _, bid := range historical {
			if bid.Won && bid.WinPrice != nil {
				avgWinPrice += *bid.WinPrice
				winCount++
			}
		}

		if winCount > 0 {
			avgWinPrice /= float64(winCount)
			// Use historical average as reference
			baseBid = (baseBid + avgWinPrice) / 2
		}
	}

	return &models.BidResponse{
		BidPrice:     baseBid,
		Confidence:   0.6, // Medium confidence for rule-based
		Strategy:     "rule_based_fallback",
		FraudRisk:    false,
		PredictionID: uuid.New().String(),
	}
}

// validatePrediction ensures prediction is within reasonable bounds
func (p *Predictor) validatePrediction(prediction *models.BidResponse, req *models.BidRequest) *models.BidResponse {
	// Ensure bid is at least floor price
	if prediction.BidPrice < req.FloorPrice {
		prediction.BidPrice = req.FloorPrice * 1.05
	}

	// Cap maximum bid at reasonable level
	maxBid := req.FloorPrice * 5.0
	if prediction.BidPrice > maxBid {
		prediction.BidPrice = maxBid
		prediction.Confidence *= 0.8 // Reduce confidence for capped bids
	}

	// Ensure confidence is between 0 and 1
	if prediction.Confidence > 1.0 {
		prediction.Confidence = 1.0
	}
	if prediction.Confidence < 0.0 {
		prediction.Confidence = 0.0
	}

	return prediction
}

// generateInsights creates performance insights
func (p *Predictor) generateInsights(ctx context.Context, bidData []*models.BidEvent) ([]string, error) {
	// Use OpenAI to generate insights
	analysis, err := p.openaiClient.AnalyzeAudienceSegment(ctx, bidData)
	if err != nil {
		// Fallback to rule-based insights
		return p.ruleBasedInsights(bidData), nil
	}

	return analysis.Insights, nil
}

// generateRecommendations creates actionable recommendations
func (p *Predictor) generateRecommendations(analysis *CampaignAnalysis) []string {
	var recommendations []string

	if analysis.WinRate < 0.3 {
		recommendations = append(recommendations, "Consider increasing bid prices to improve win rate")
	}

	if analysis.ConversionRate < 0.05 {
		recommendations = append(recommendations, "Review audience targeting - conversion rate is below industry average")
	}

	if analysis.AvgBidPrice > analysis.AvgWinPrice*1.5 {
		recommendations = append(recommendations, "Optimize bidding strategy - you may be overbidding")
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations, "Campaign performance looks healthy - continue current strategy")
	}

	return recommendations
}

// ruleBasedFraudDetection implements basic fraud detection rules
func (p *Predictor) ruleBasedFraudDetection(bidData []*models.BidEvent) *FraudDetectionResult {
	userActivity := make(map[string]int)
	userConversions := make(map[string]int)

	for _, bid := range bidData {
		userActivity[bid.UserID]++
		if bid.Converted {
			userConversions[bid.UserID]++
		}
	}

	// Check for suspicious patterns
	var suspiciousUsers []string
	for userID, activity := range userActivity {
		conversions := userConversions[userID]
		conversionRate := float64(conversions) / float64(activity)

		// Flag users with unusually high activity or conversion rates
		if activity > 50 || conversionRate > 0.5 {
			suspiciousUsers = append(suspiciousUsers, userID)
		}
	}

	fraudDetected := len(suspiciousUsers) > 0
	severity := len(suspiciousUsers) * 2

	if severity > 10 {
		severity = 10
	}

	return &FraudDetectionResult{
		FraudDetected:   fraudDetected,
		Confidence:      0.7,
		SuspiciousUsers: suspiciousUsers,
		Severity:        severity,
		DetectionMethod: "rule_based",
	}
}

// ruleBasedInsights generates basic insights from bid data
func (p *Predictor) ruleBasedInsights(bidData []*models.BidEvent) []string {
	var insights []string

	// Analyze device performance
	devicePerf := make(map[string]float64)
	deviceCount := make(map[string]int)

	for _, bid := range bidData {
		deviceCount[bid.DeviceType]++
		if bid.Converted {
			devicePerf[bid.DeviceType]++
		}
	}

	for device, conversions := range devicePerf {
		rate := conversions / float64(deviceCount[device])
		if rate > 0.1 {
			insights = append(insights, fmt.Sprintf("%s devices show strong performance (%.1f%% conversion rate)", device, rate*100))
		}
	}

	return insights
}

func (p *Predictor) combineFraudResults(ruleResult *FraudDetectionResult, aiResult *FraudAnalysis) *FraudDetectionResult {
	// Combine confidence scores
	combinedConfidence := (ruleResult.Confidence + aiResult.Confidence) / 2

	return &FraudDetectionResult{
		FraudDetected:   ruleResult.FraudDetected || aiResult.FraudDetected,
		Confidence:      combinedConfidence,
		SuspiciousUsers: ruleResult.SuspiciousUsers,
		Severity:        max(ruleResult.Severity, aiResult.Severity),
		DetectionMethod: "combined",
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Supporting types
type CampaignAnalysis struct {
	CampaignID      uuid.UUID `json:"campaign_id"`
	Period          string    `json:"period"`
	TotalBids       int       `json:"total_bids"`
	WinRate         float64   `json:"win_rate"`
	ConversionRate  float64   `json:"conversion_rate"`
	AvgBidPrice     float64   `json:"avg_bid_price"`
	AvgWinPrice     float64   `json:"avg_win_price"`
	TotalSpend      float64   `json:"total_spend"`
	Insights        []string  `json:"insights"`
	Recommendations []string  `json:"recommendations"`
	Message         string    `json:"message,omitempty"`
}

type FraudDetectionResult struct {
	CampaignID      uuid.UUID `json:"campaign_id"`
	FraudDetected   bool      `json:"fraud_detected"`
	Confidence      float64   `json:"confidence"`
	SuspiciousUsers []string  `json:"suspicious_users"`
	Severity        int       `json:"severity"`
	DetectionMethod string    `json:"detection_method"`
	Message         string    `json:"message,omitempty"`
}
