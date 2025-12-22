package trpc

import (
	"context"
	"fmt"
	"time"

	"github.com/baskint/bidding-analysis/internal/models"
	"github.com/google/uuid"
)

// Request types for bid operations
type BidSubmitRequest struct {
	CampaignID      string    `json:"campaign_id"`
	UserID          uuid.UUID `json:"user_id"`
	BidPrice        float64   `json:"bid_price"`
	FloorPrice      float64   `json:"floor_price"`
	DeviceType      string    `json:"device_type"`
	OS              string    `json:"os"`
	Browser         string    `json:"browser"`
	Country         string    `json:"country"`
	Region          string    `json:"region"`
	City            string    `json:"city"`
	Keywords        []string  `json:"keywords"`
	SegmentID       string    `json:"segment_id"`
	SegmentCategory string    `json:"segment_category"`
}

type BidPredictionRequest struct {
	CampaignID      string   `json:"campaign_id"`
	UserSegment     string   `json:"user_segment"`
	DeviceType      string   `json:"device_type"`
	Country         string   `json:"country"`
	FloorPrice      float64  `json:"floor_price"`
	Keywords        []string `json:"keywords"`
	EngagementScore float64  `json:"engagement_score"`
	ConversionProb  float64  `json:"conversion_probability"`
}

// Response types
type BidSubmitResponse struct {
	BidEventID     string  `json:"bid_event_id"`
	Status         string  `json:"status"`
	ProcessedAt    string  `json:"processed_at"`
	WinProbability float64 `json:"win_probability"`
	Message        string  `json:"message"`
}

type BidPredictionResponse struct {
	PredictedBid float64 `json:"predicted_bid"`
	Confidence   float64 `json:"confidence"`
	Strategy     string  `json:"strategy"`
	FraudRisk    bool    `json:"fraud_risk"`
	Reasoning    string  `json:"reasoning"`
}

type BidStreamData struct {
	BidEventID      string    `json:"bid_event_id"`
	CampaignID      string    `json:"campaign_id"`
	BidPrice        float64   `json:"bid_price"`
	WinPrice        *float64  `json:"win_price,omitempty"`
	Won             bool      `json:"won"`
	Converted       bool      `json:"converted"`
	Timestamp       time.Time `json:"timestamp"`
	SegmentCategory string    `json:"segment_category"`
	DeviceType      string    `json:"device_type"`
	Country         string    `json:"country"`
}

// ============================================================================
// BID HANDLERS
// ============================================================================

// handleSubmitBid processes a new bid submission
func (h *Handler) handleSubmitBid(ctx context.Context, userID uuid.UUID, req interface{}) (interface{}, error) {
	params := req.(*BidSubmitRequest)

	// Validate required fields
	if params.CampaignID == "" {
		return nil, fmt.Errorf("campaign ID is required")
	}
	if params.BidPrice <= 0 {
		return nil, fmt.Errorf("bid price must be positive")
	}
	if params.FloorPrice <= 0 {
		return nil, fmt.Errorf("floor price must be positive")
	}

	// Parse campaign ID
	campaignID, err := uuid.Parse(params.CampaignID)
	if err != nil {
		return nil, fmt.Errorf("invalid campaign ID format")
	}

	// Create bid event
	bidEvent := &models.BidEvent{
		CampaignID:      campaignID,
		UserID:          params.UserID,
		BidPrice:        params.BidPrice,
		FloorPrice:      params.FloorPrice,
		Won:             false, // Will be determined by auction logic
		Converted:       false,
		SegmentID:       params.SegmentID,
		SegmentCategory: params.SegmentCategory,
		Country:         params.Country,
		Region:          params.Region,
		City:            params.City,
		DeviceType:      params.DeviceType,
		OS:              params.OS,
		Browser:         params.Browser,
		IsMobile:        params.DeviceType == "mobile",
		Timestamp:       time.Now(),
	}

	// Convert keywords
	if len(params.Keywords) > 0 {
		bidEvent.Keywords = params.Keywords
	}

	// Simulate auction logic
	winProbability := calculateWinProbability(params.BidPrice, params.FloorPrice)
	bidEvent.Won = winProbability > 0.5

	if bidEvent.Won {
		// Calculate win price (typically 85-95% of bid price)
		winPrice := params.BidPrice * (0.85 + 0.1*(winProbability-0.5)*2)
		bidEvent.WinPrice = &winPrice
	}

	// Store the bid event
	if err := h.bidStore.StoreBidEvent(bidEvent); err != nil {
		return nil, fmt.Errorf("failed to process bid: %w", err)
	}

	// Prepare response
	response := BidSubmitResponse{
		BidEventID:     bidEvent.ID.String(),
		Status:         getAuctionStatus(bidEvent.Won),
		ProcessedAt:    bidEvent.Timestamp.Format(time.RFC3339),
		WinProbability: winProbability,
		Message:        generateBidMessage(bidEvent.Won, params.BidPrice, bidEvent.WinPrice),
	}

	return response, nil
}

// handlePredictBid uses AI to predict optimal bid price
func (h *Handler) handlePredictBid(ctx context.Context, userID uuid.UUID, req interface{}) (interface{}, error) {
	params := req.(*BidPredictionRequest)

	// Validate required fields
	if params.CampaignID == "" || params.FloorPrice <= 0 {
		return nil, fmt.Errorf("campaign ID and floor price are required")
	}

	// Create ML prediction request
	mlRequest := &models.BidRequest{
		CampaignID: uuid.MustParse(params.CampaignID),
		UserSegment: models.UserSegment{
			SegmentID:             params.UserSegment,
			Category:              params.UserSegment,
			EngagementScore:       params.EngagementScore,
			ConversionProbability: params.ConversionProb,
		},
		GeoLocation: models.GeoLocation{
			Country: params.Country,
		},
		DeviceInfo: models.DeviceInfo{
			DeviceType: params.DeviceType,
			IsMobile:   params.DeviceType == "mobile",
		},
		FloorPrice: params.FloorPrice,
		Keywords:   params.Keywords,
		Timestamp:  time.Now(),
	}

	// Get AI prediction
	prediction, err := h.predictor.PredictBid(ctx, mlRequest)
	if err != nil {
		// Fallback to rule-based prediction
		prediction = h.getFallbackPrediction(params)
	}

	response := BidPredictionResponse{
		PredictedBid: prediction.BidPrice,
		Confidence:   prediction.Confidence,
		Strategy:     prediction.Strategy,
		FraudRisk:    prediction.FraudRisk,
		Reasoning:    generatePredictionReasoning(prediction, params),
	}

	return response, nil
}

// handleGetBidStream returns recent bid activity
func (h *Handler) handleGetBidStream(ctx context.Context, userID uuid.UUID, req interface{}) (interface{}, error) {
	// Get recent bids (last 50)
	bids, err := h.bidStore.GetRecentBids(ctx, 50)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve bid stream: %w", err)
	}

	// Convert to stream data
	streamData := make([]BidStreamData, len(bids))
	for i, bid := range bids {
		streamData[i] = BidStreamData{
			BidEventID:      bid.ID.String(),
			CampaignID:      bid.CampaignID.String(),
			BidPrice:        bid.BidPrice,
			WinPrice:        bid.WinPrice,
			Won:             bid.Won,
			Converted:       bid.Converted,
			Timestamp:       bid.Timestamp,
			SegmentCategory: bid.SegmentCategory,
			DeviceType:      bid.DeviceType,
			Country:         bid.Country,
		}
	}

	return streamData, nil
}

// processBid is a legacy endpoint - redirects to handleSubmitBid
func (h *Handler) processBid(ctx context.Context, userID uuid.UUID, req interface{}) (interface{}, error) {
	return h.handleSubmitBid(ctx, userID, req)
}

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================

func calculateWinProbability(bidPrice, floorPrice float64) float64 {
	ratio := bidPrice / floorPrice

	if ratio <= 1.0 {
		return 0.1 // Very low chance if bid <= floor
	} else if ratio >= 3.0 {
		return 0.95 // Very high chance if bid >= 3x floor
	}

	// Linear interpolation between 0.1 and 0.95
	return 0.1 + (0.85 * (ratio - 1.0) / 2.0)
}

func getAuctionStatus(won bool) string {
	if won {
		return "won"
	}
	return "lost"
}

func generateBidMessage(won bool, bidPrice float64, winPrice *float64) string {
	if won && winPrice != nil {
		return fmt.Sprintf("Bid won! Paid $%.4f for $%.4f bid", *winPrice, bidPrice)
	} else if won {
		return fmt.Sprintf("Bid won at $%.4f", bidPrice)
	}
	return fmt.Sprintf("Bid lost at $%.4f", bidPrice)
}

func generatePredictionReasoning(prediction *models.BidResponse, req *BidPredictionRequest) string {
	confidence := "medium"
	if prediction.Confidence > 0.8 {
		confidence = "high"
	} else if prediction.Confidence < 0.4 {
		confidence = "low"
	}

	return fmt.Sprintf("AI recommends $%.4f based on %s confidence analysis of user segment '%s' and device type '%s'",
		prediction.BidPrice, confidence, req.UserSegment, req.DeviceType)
}

func (h *Handler) getFallbackPrediction(req *BidPredictionRequest) *models.BidResponse {
	// Simple rule-based fallback
	baseBid := req.FloorPrice * 1.5

	// Adjust based on engagement score
	if req.EngagementScore > 0.7 {
		baseBid *= 1.2
	} else if req.EngagementScore < 0.3 {
		baseBid *= 0.8
	}

	// Adjust based on device type
	if req.DeviceType == "mobile" {
		baseBid *= 1.1 // Mobile traffic premium
	}

	return &models.BidResponse{
		BidPrice:   baseBid,
		Confidence: 0.6,
		Strategy:   "rule_based_fallback",
		FraudRisk:  false,
	}
}
