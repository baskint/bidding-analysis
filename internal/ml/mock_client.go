package ml

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"

	"github.com/baskint/bidding-analysis/internal/models"
)

// MockOpenAIClient simulates OpenAI responses for development
type MockOpenAIClient struct {
	rng *rand.Rand
}

// NewMockOpenAIClient creates a new mock client
func NewMockOpenAIClient() *MockOpenAIClient {
	return &MockOpenAIClient{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// PredictBidPrice simulates bid price prediction
func (m *MockOpenAIClient) PredictBidPrice(ctx context.Context, req *models.BidRequest, historicalData []*models.BidEvent) (*models.BidResponse, error) {
	// Simulate processing delay
	time.Sleep(time.Millisecond * time.Duration(50+m.rng.Intn(100)))

	// Base bid calculation
	baseBid := req.FloorPrice * (1.1 + m.rng.Float64()*0.8) // 1.1x to 1.9x floor price

	// Adjust based on user segment
	if req.UserSegment.ConversionProbability > 0.1 {
		multiplier := 1 + req.UserSegment.ConversionProbability*0.5
		baseBid *= multiplier
	}

	// Adjust based on engagement score
	if req.UserSegment.EngagementScore > 0.5 {
		baseBid *= (1 + req.UserSegment.EngagementScore*0.2)
	}

	// Adjust based on device type
	switch req.DeviceInfo.DeviceType {
	case "mobile":
		baseBid *= 0.95 // Slightly lower for mobile
	case "desktop":
		baseBid *= 1.05 // Slightly higher for desktop
	}

	// Adjust based on historical performance
	if len(historicalData) > 0 {
		var winRate, conversionRate float64
		var wonBids, conversions int

		for _, bid := range historicalData {
			if bid.Won {
				wonBids++
			}
			if bid.Converted {
				conversions++
			}
		}

		winRate = float64(wonBids) / float64(len(historicalData))
		if wonBids > 0 {
			conversionRate = float64(conversions) / float64(wonBids)
		}

		// Adjust bid based on performance
		if winRate > 0.6 {
			baseBid *= 1.1 // Increase if winning too often
		} else if winRate < 0.2 {
			baseBid *= 0.9 // Decrease if winning too rarely
		}

		if conversionRate > 0.1 {
			baseBid *= 1.15 // Increase for high-converting campaigns
		}
	}

	// Calculate confidence based on data quality
	confidence := 0.6 + m.rng.Float64()*0.3 // Base confidence 0.6-0.9

	if len(historicalData) > 50 {
		confidence += 0.1 // More confidence with more data
	}
	if len(historicalData) > 100 {
		confidence += 0.05
	}

	// Determine strategy
	strategy := m.determineStrategy(req, historicalData, baseBid)

	// Simple fraud risk assessment
	fraudRisk := m.assessFraudRisk(req, historicalData)

	return &models.BidResponse{
		BidPrice:     baseBid,
		Confidence:   confidence,
		Strategy:     strategy,
		FraudRisk:    fraudRisk,
		PredictionID: uuid.New().String(),
	}, nil
}

// AnalyzeAudienceSegment simulates audience analysis
func (m *MockOpenAIClient) AnalyzeAudienceSegment(ctx context.Context, bidEvents []*models.BidEvent) (*AudienceAnalysis, error) {
	// Simulate processing delay
	time.Sleep(time.Millisecond * time.Duration(100+m.rng.Intn(200)))

	deviceTypes := make(map[string]int)
	geoRegions := make(map[string]int)

	for _, bid := range bidEvents {
		deviceTypes[bid.DeviceType]++
		region := fmt.Sprintf("%s, %s", bid.Region, bid.Country)
		geoRegions[region]++
	}

	var segments []string
	var insights []string

	// Generate mock segments
	if deviceTypes["mobile"] > len(bidEvents)/2 {
		segments = append(segments, "mobile_heavy_users")
		insights = append(insights, "Campaign shows strong mobile user engagement")
	}

	if deviceTypes["desktop"] > len(bidEvents)/3 {
		segments = append(segments, "desktop_professionals")
		insights = append(insights, "Desktop users show higher engagement during business hours")
	}

	// Generate geographic insights
	if len(geoRegions) > 10 {
		insights = append(insights, "Campaign has broad geographic reach across multiple regions")
	} else {
		insights = append(insights, "Campaign is geographically concentrated - consider expansion")
	}

	// Generate time-based insights
	insights = append(insights, "Peak engagement occurs during 2-4 PM and 7-9 PM time slots")

	return &AudienceAnalysis{
		Segments: segments,
		Insights: insights,
	}, nil
}

// DetectFraud simulates fraud detection
func (m *MockOpenAIClient) DetectFraud(ctx context.Context, bidEvents []*models.BidEvent) (*FraudAnalysis, error) {
	// Simulate processing delay
	time.Sleep(time.Millisecond * time.Duration(150+m.rng.Intn(100)))

	userActivity := make(map[string]int)
	userConversions := make(map[string]int)

	for _, bid := range bidEvents {
		userActivity[bid.UserID]++
		if bid.Converted {
			userConversions[bid.UserID]++
		}
	}

	var patterns []string
	fraudDetected := false
	confidence := 0.8
	severity := 1

	// Check for suspicious patterns
	for userID, activity := range userActivity {
		conversions := userConversions[userID]
		conversionRate := float64(conversions) / float64(activity)

		// Detect unusual activity
		if activity > 100 {
			patterns = append(patterns, fmt.Sprintf("User %s shows excessive activity (%d events)", userID[:8], activity))
			fraudDetected = true
			severity = max(severity, 6)
		}

		if conversionRate > 0.8 && activity > 10 {
			patterns = append(patterns, fmt.Sprintf("User %s has suspiciously high conversion rate (%.1f%%)", userID[:8], conversionRate*100))
			fraudDetected = true
			severity = max(severity, 7)
		}
	}

	// Random chance of detecting other fraud types
	if m.rng.Float64() < 0.1 { // 10% chance
		patterns = append(patterns, "Detected potential click farming based on timing patterns")
		fraudDetected = true
		severity = max(severity, 5)
	}

	if m.rng.Float64() < 0.05 { // 5% chance
		patterns = append(patterns, "Geographic anomaly detected - unusual traffic from specific regions")
		fraudDetected = true
		severity = max(severity, 4)
	}

	// Adjust confidence based on data quality
	if len(bidEvents) < 20 {
		confidence *= 0.7 // Lower confidence with less data
	}

	return &FraudAnalysis{
		FraudDetected: fraudDetected,
		Confidence:    confidence,
		Patterns:      patterns,
		Severity:      severity,
	}, nil
}

// determineStrategy selects an appropriate bidding strategy
func (m *MockOpenAIClient) determineStrategy(req *models.BidRequest, historical []*models.BidEvent, bidPrice float64) string {
	strategies := []string{
		"aggressive_targeting",
		"conservative_bidding",
		"performance_optimized",
		"brand_awareness",
		"conversion_focused",
	}

	// Strategy selection based on context
	if req.UserSegment.ConversionProbability > 0.2 {
		return "conversion_focused"
	}

	if len(historical) > 0 {
		var avgBid float64
		for _, bid := range historical {
			avgBid += bid.BidPrice
		}
		avgBid /= float64(len(historical))

		if bidPrice > avgBid*1.2 {
			return "aggressive_targeting"
		} else if bidPrice < avgBid*0.8 {
			return "conservative_bidding"
		}
		return "performance_optimized"
	}

	// Random strategy for new campaigns
	return strategies[m.rng.Intn(len(strategies))]
}

// assessFraudRisk performs basic fraud risk assessment
func (m *MockOpenAIClient) assessFraudRisk(req *models.BidRequest, historical []*models.BidEvent) bool {
	// Low chance of fraud for most requests
	if m.rng.Float64() < 0.05 { // 5% chance
		return true
	}

	// Check for suspicious user patterns
	if len(historical) > 0 {
		userCounts := make(map[string]int)
		for _, bid := range historical {
			userCounts[bid.UserID]++
		}

		// High fraud risk if any user has excessive activity
		for _, count := range userCounts {
			if count > 50 {
				return true
			}
		}
	}

	return false
}
