package ml

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/sashabaranov/go-openai"

	"github.com/baskint/bidding-analysis/internal/models"
)

// OpenAIClient wraps the OpenAI API client
type OpenAIClient struct {
	client *openai.Client
}

// NewOpenAIClient creates a new OpenAI client
func NewOpenAIClient(apiKey string) *OpenAIClient {
	return &OpenAIClient{
		client: openai.NewClient(apiKey),
	}
}

// PredictBidPrice uses OpenAI to predict optimal bid price
func (c *OpenAIClient) PredictBidPrice(ctx context.Context, req *models.BidRequest, historicalData []*models.BidEvent) (*models.BidResponse, error) {
	prompt := c.buildPrompt(req, historicalData)

	resp, err := c.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are an expert in digital advertising bid optimization. Analyze the data and predict the optimal bid price. Respond with JSON format: {\"bid_price\": number, \"confidence\": number, \"strategy\": \"string\", \"fraud_risk\": boolean}",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			Temperature: 0.3, // Lower temperature for more consistent predictions
		},
	)

	if err != nil {
		return nil, fmt.Errorf("OpenAI API error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	return c.parseResponse(resp.Choices[0].Message.Content)
}

// AnalyzeAudienceSegment uses OpenAI to analyze and segment audiences
func (c *OpenAIClient) AnalyzeAudienceSegment(ctx context.Context, bidEvents []*models.BidEvent) (*AudienceAnalysis, error) {
	prompt := c.buildAudiencePrompt(bidEvents)

	resp, err := c.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are an expert in audience analysis for digital advertising. Analyze user behavior patterns and provide insights. Respond with JSON format.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			Temperature: 0.5,
		},
	)

	if err != nil {
		return nil, fmt.Errorf("OpenAI API error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	return c.parseAudienceResponse(resp.Choices[0].Message.Content)
}

// DetectFraud uses OpenAI to detect potential fraud patterns
func (c *OpenAIClient) DetectFraud(ctx context.Context, bidEvents []*models.BidEvent) (*FraudAnalysis, error) {
	prompt := c.buildFraudPrompt(bidEvents)

	resp, err := c.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are an expert in fraud detection for digital advertising. Analyze bid patterns for suspicious activity. Respond with JSON format: {\"fraud_detected\": boolean, \"confidence\": number, \"patterns\": [\"string\"], \"severity\": number}",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			Temperature: 0.2, // Very low temperature for fraud detection
		},
	)

	if err != nil {
		return nil, fmt.Errorf("OpenAI API error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	return c.parseFraudResponse(resp.Choices[0].Message.Content)
}

// buildPrompt creates a prompt for bid price prediction
func (c *OpenAIClient) buildPrompt(req *models.BidRequest, historical []*models.BidEvent) string {
	var prompt strings.Builder

	prompt.WriteString("Analyze this bid request and predict optimal bid price:\n\n")

	// Current bid request
	prompt.WriteString(fmt.Sprintf("REQUEST:\n"))
	prompt.WriteString(fmt.Sprintf("- Campaign: %s\n", req.CampaignID))
	prompt.WriteString(fmt.Sprintf("- User Segment: %s (%s)\n", req.UserSegment.SegmentID, req.UserSegment.Category))
	prompt.WriteString(fmt.Sprintf("- Engagement Score: %.2f\n", req.UserSegment.EngagementScore))
	prompt.WriteString(fmt.Sprintf("- Conversion Probability: %.2f\n", req.UserSegment.ConversionProbability))
	prompt.WriteString(fmt.Sprintf("- Location: %s, %s, %s\n", req.GeoLocation.City, req.GeoLocation.Region, req.GeoLocation.Country))
	prompt.WriteString(fmt.Sprintf("- Device: %s (%s, %s)\n", req.DeviceInfo.DeviceType, req.DeviceInfo.OS, req.DeviceInfo.Browser))
	prompt.WriteString(fmt.Sprintf("- Floor Price: $%.4f\n", req.FloorPrice))
	prompt.WriteString(fmt.Sprintf("- Keywords: %v\n", req.Keywords))

	// Historical data summary
	if len(historical) > 0 {
		prompt.WriteString(fmt.Sprintf("\nHISTORICAL DATA (%d recent bids):\n", len(historical)))

		var totalBids, wonBids, conversions int
		var totalSpend, totalWinPrice float64

		for _, bid := range historical {
			totalBids++
			if bid.Won {
				wonBids++
				if bid.WinPrice != nil {
					totalWinPrice += *bid.WinPrice
				}
			}
			if bid.Converted {
				conversions++
			}
			totalSpend += bid.BidPrice
		}

		avgBid := totalSpend / float64(totalBids)
		winRate := float64(wonBids) / float64(totalBids)
		conversionRate := float64(conversions) / float64(wonBids)
		if wonBids == 0 {
			conversionRate = 0
		}
		avgWinPrice := totalWinPrice / float64(wonBids)
		if wonBids == 0 {
			avgWinPrice = 0
		}

		prompt.WriteString(fmt.Sprintf("- Win Rate: %.2f%%\n", winRate*100))
		prompt.WriteString(fmt.Sprintf("- Conversion Rate: %.2f%%\n", conversionRate*100))
		prompt.WriteString(fmt.Sprintf("- Average Bid: $%.4f\n", avgBid))
		prompt.WriteString(fmt.Sprintf("- Average Win Price: $%.4f\n", avgWinPrice))
	}

	prompt.WriteString("\nConsider factors like competition, user engagement, conversion probability, and historical performance.")
	prompt.WriteString("\nRecommend a bid price that maximizes ROI while maintaining competitiveness.")

	return prompt.String()
}

// buildAudiencePrompt creates a prompt for audience analysis
func (c *OpenAIClient) buildAudiencePrompt(bidEvents []*models.BidEvent) string {
	// Implementation for audience analysis prompt
	return "Analyze audience patterns from bid data..."
}

// buildFraudPrompt creates a prompt for fraud detection
func (c *OpenAIClient) buildFraudPrompt(bidEvents []*models.BidEvent) string {
	var prompt strings.Builder

	prompt.WriteString("Analyze these bid events for fraud patterns:\n\n")

	// Group by user for pattern analysis
	userStats := make(map[uuid.UUID]*UserStats)

	for _, bid := range bidEvents {
		if _, exists := userStats[bid.UserID]; !exists {
			userStats[bid.UserID] = &UserStats{}
		}
		stats := userStats[bid.UserID]
		stats.TotalBids++
		if bid.Won {
			stats.WonBids++
		}
		if bid.Converted {
			stats.Conversions++
		}
	}

	prompt.WriteString(fmt.Sprintf("SUMMARY:\n"))
	prompt.WriteString(fmt.Sprintf("- Total Events: %d\n", len(bidEvents)))
	prompt.WriteString(fmt.Sprintf("- Unique Users: %d\n", len(userStats)))

	// Look for suspicious patterns
	prompt.WriteString("\nSUSPICIOUS PATTERNS TO CHECK:\n")
	prompt.WriteString("- Users with abnormally high bid/win rates\n")
	prompt.WriteString("- Rapid-fire bidding from same user\n")
	prompt.WriteString("- Conversion rates that are too good to be true\n")
	prompt.WriteString("- Geographic anomalies\n")

	return prompt.String()
}

// parseResponse parses OpenAI response for bid prediction
func (c *OpenAIClient) parseResponse(content string) (*models.BidResponse, error) {
	// Try to extract JSON from the response
	content = strings.TrimSpace(content)

	// Remove markdown code block markers if present
	if strings.HasPrefix(content, "```json") {
		content = strings.TrimPrefix(content, "```json")
		content = strings.TrimSuffix(content, "```")
		content = strings.TrimSpace(content)
	}

	var result struct {
		BidPrice   float64 `json:"bid_price"`
		Confidence float64 `json:"confidence"`
		Strategy   string  `json:"strategy"`
		FraudRisk  bool    `json:"fraud_risk"`
	}

	if err := json.Unmarshal([]byte(content), &result); err != nil {
		// Fallback: try to parse manually
		return c.parseManually(content)
	}

	return &models.BidResponse{
		BidPrice:   result.BidPrice,
		Confidence: result.Confidence,
		Strategy:   result.Strategy,
		FraudRisk:  result.FraudRisk,
	}, nil
}

// parseManually attempts to parse response manually if JSON parsing fails
func (c *OpenAIClient) parseManually(content string) (*models.BidResponse, error) {
	lines := strings.Split(content, "\n")
	response := &models.BidResponse{
		Strategy: "manual_parse",
	}

	for _, line := range lines {
		if strings.Contains(strings.ToLower(line), "bid") && strings.Contains(line, "$") {
			// Extract price
			parts := strings.Split(line, "$")
			if len(parts) > 1 {
				priceStr := strings.Fields(parts[1])[0]
				if price, err := strconv.ParseFloat(priceStr, 64); err == nil {
					response.BidPrice = price
				}
			}
		}
		if strings.Contains(strings.ToLower(line), "confidence") {
			// Extract confidence
			parts := strings.Fields(line)
			for _, part := range parts {
				if strings.HasSuffix(part, "%") {
					confStr := strings.TrimSuffix(part, "%")
					if conf, err := strconv.ParseFloat(confStr, 64); err == nil {
						response.Confidence = conf / 100
					}
				}
			}
		}
	}

	// Set defaults if parsing failed
	if response.BidPrice == 0 {
		response.BidPrice = 1.0 // Default bid
	}
	if response.Confidence == 0 {
		response.Confidence = 0.5 // Default confidence
	}

	return response, nil
}

func (c *OpenAIClient) parseAudienceResponse(content string) (*AudienceAnalysis, error) {
	// Placeholder implementation
	return &AudienceAnalysis{}, nil
}

func (c *OpenAIClient) parseFraudResponse(content string) (*FraudAnalysis, error) {
	// Placeholder implementation
	return &FraudAnalysis{}, nil
}

// Supporting types
type UserStats struct {
	TotalBids   int
	WonBids     int
	Conversions int
}

type AudienceAnalysis struct {
	Segments []string `json:"segments"`
	Insights []string `json:"insights"`
}

type FraudAnalysis struct {
	FraudDetected bool     `json:"fraud_detected"`
	Confidence    float64  `json:"confidence"`
	Patterns      []string `json:"patterns"`
	Severity      int      `json:"severity"`
}
