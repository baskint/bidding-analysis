package trpc

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/baskint/bidding-analysis/internal/ml"
	"github.com/baskint/bidding-analysis/internal/models"
	"github.com/baskint/bidding-analysis/internal/store"
)

// Handler contains the dependencies for tRPC handlers
type Handler struct {
	bidStore      *store.BidStore
	campaignStore *store.CampaignStore
	userStore     *store.UserStore
	mlModelStore  *store.MLModelStore
	settingsStore *store.SettingsStore
	predictor     *ml.Predictor
	jwtSecret     string
}

// NewHandler creates a new tRPC Handler instance
func NewHandler(bidStore *store.BidStore, campaignStore *store.CampaignStore, predictor *ml.Predictor) *Handler {
	// Get JWT secret from environment, with fallback
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-default-development-secret"
		log.Println("Warning: Using default JWT secret. Set JWT_SECRET environment variable for production.")
	}

	// Initialize UserStore using the same database connection as other stores
	db := bidStore.DB() // You'll need to add this method to BidStore
	userStore := store.NewUserStore(db)
	mlModelStore := store.NewMLModelStore(db)
	settingsStore := store.NewSettingsStore(db)

	return &Handler{
		bidStore:      bidStore,
		campaignStore: campaignStore,
		userStore:     userStore,
		mlModelStore:  mlModelStore,
		settingsStore: settingsStore,
		predictor:     predictor,
		jwtSecret:     jwtSecret,
	}
}

// SetupRoutes configures all tRPC routes
func (h *Handler) SetupRoutes() http.Handler {
	router := mux.NewRouter()

	// Apply middleware
	router.Use(loggingMiddleware)

	// Health check
	router.HandleFunc("/health", h.healthCheck).Methods("GET")
	router.HandleFunc("/", h.rootHandler).Methods("GET")

	// tRPC routes
	api := router.PathPrefix("/trpc").Subrouter()

	// Debug endpoint
	api.HandleFunc("/debug", h.debugEndpoint).Methods("GET")

	// Auth procedures (public - no authentication required)
	api.HandleFunc("/auth.login", h.login).Methods("POST")
	api.HandleFunc("/auth.register", h.register).Methods("POST")

	// Protected routes (require authentication)
	protected := api.PathPrefix("").Subrouter()
	protected.Use(h.authMiddleware)

	// Auth procedures (protected)
	protected.HandleFunc("/auth.me", h.getMe).Methods("GET", "POST")

	// Bidding procedures
	protected.HandleFunc("/bidding.submit", h.handleSubmitBid).Methods("POST")
	protected.HandleFunc("/bidding.predict", h.handlePredictBid).Methods("POST")
	protected.HandleFunc("/bidding.stream", h.handleGetBidStream).Methods("GET")
	protected.HandleFunc("/bidding.processBid", h.processBid).Methods("POST") // Legacy endpoint

	// Campaign procedures
	protected.HandleFunc("/campaign.getStats", h.getCampaignStats).Methods("GET", "POST")
	protected.HandleFunc("/campaign.getBidHistory", h.getBidHistory).Methods("GET", "POST")
	protected.HandleFunc("/campaign.list", h.listCampaigns).Methods("GET", "POST")
	protected.HandleFunc("/campaign.create", h.createCampaign).Methods("POST")

	// Analytics procedures
	protected.HandleFunc("/analytics.getFraudAlerts", h.getFraudAlerts).Methods("GET", "POST")
	protected.HandleFunc("/analytics.getModelAccuracy", h.getModelAccuracy).Methods("GET", "POST")
	protected.HandleFunc("/analytics.getDashboardMetrics", h.getDashboardMetrics).Methods("GET", "POST")

	protected.HandleFunc("/campaign.get", h.getCampaign).Methods("GET", "POST")
	protected.HandleFunc("/campaign.update", h.updateCampaign).Methods("POST")
	protected.HandleFunc("/campaign.delete", h.deleteCampaign).Methods("POST")
	protected.HandleFunc("/campaign.pause", h.pauseCampaign).Methods("POST")
	protected.HandleFunc("/campaign.activate", h.activateCampaign).Methods("POST")
	protected.HandleFunc("/campaign.listWithMetrics", h.listCampaignsEnhanced).Methods("GET", "POST")
	protected.HandleFunc("/campaign.getDailyMetrics", h.getDailyMetrics).Methods("GET", "POST")

	// Analysis page procedures
	// Analysis page procedures - Using WithAuth wrapper
	protected.HandleFunc("/analytics.getPerformanceOverview", h.WithAuth(h.getPerformanceOverview, &DateRangeRequest{})).Methods("GET", "POST")
	protected.HandleFunc("/analytics.getKeywordAnalysis", h.WithAuth(h.getKeywordAnalysis, &KeywordAnalysisRequest{})).Methods("GET", "POST")
	protected.HandleFunc("/analytics.getDeviceBreakdown", h.WithAuth(h.getDeviceBreakdown, &DateRangeRequest{})).Methods("GET", "POST")
	protected.HandleFunc("/analytics.getGeoBreakdown", h.WithAuth(h.getGeoBreakdown, &DateRangeRequest{})).Methods("GET", "POST")
	protected.HandleFunc("/analytics.getHourlyPerformance", h.WithAuth(h.getHourlyPerformance, &DateRangeRequest{})).Methods("GET", "POST")
	protected.HandleFunc("/analytics.getDailyTrends", h.WithAuth(h.getDailyTrends, &DateRangeRequest{})).Methods("GET", "POST")
	protected.HandleFunc("/analytics.getCompetitiveAnalysis", h.WithAuth(h.getCompetitiveAnalysis, &CompetitiveAnalysisRequest{})).Methods("GET", "POST")
	// protected.HandleFunc("/analytics.getCampaignComparison", h.getCampaignComparison).Methods("GET", "POST") // TODO: Refactor later

	// ML Model procedures
	protected.HandleFunc("/mlModel.list", h.listMLModels).Methods("GET", "POST")
	protected.HandleFunc("/mlModel.get", h.getMLModel).Methods("GET", "POST")
	protected.HandleFunc("/mlModel.create", h.createMLModel).Methods("POST")
	protected.HandleFunc("/mlModel.update", h.updateMLModel).Methods("POST")
	protected.HandleFunc("/mlModel.delete", h.deleteMLModel).Methods("POST")
	protected.HandleFunc("/mlModel.setDefault", h.setDefaultMLModel).Methods("POST")
	protected.HandleFunc("/mlModel.getDefault", h.getDefaultMLModel).Methods("GET", "POST")

	// Fraud procedures - Using WithAuth wrapper
	protected.HandleFunc("/fraud.getOverview", h.WithAuth(h.getFraudOverview, &FraudOverviewRequest{})).Methods("GET", "POST")
	protected.HandleFunc("/fraud.getAlerts", h.WithAuth(h.getRealFraudAlerts, &FraudAlertsRequest{})).Methods("GET", "POST")
	protected.HandleFunc("/fraud.updateAlert", h.WithAuth(h.updateFraudAlert, &UpdateFraudAlertRequest{})).Methods("POST")
	protected.HandleFunc("/fraud.getTrends", h.WithAuth(h.getFraudTrends, &FraudTrendsRequest{})).Methods("GET", "POST")
	protected.HandleFunc("/fraud.getDeviceAnalysis", h.WithAuth(h.getDeviceFraudAnalysis, &DeviceFraudRequest{})).Methods("GET", "POST")
	protected.HandleFunc("/fraud.getGeoAnalysis", h.WithAuth(h.getGeoFraudAnalysis, &GeoFraudRequest{})).Methods("GET", "POST")
	protected.HandleFunc("/fraud.createAlert", h.WithAuth(h.createFraudAlert, &CreateFraudAlertRequest{})).Methods("POST")

	// Alert procedures - Using WithAuth wrapper
	protected.HandleFunc("/alerts.getAlerts", h.WithAuth(h.getAlerts, &GetAlertsRequest{})).Methods("GET", "POST")
	protected.HandleFunc("/alerts.getOverview", h.WithAuth(h.getAlertOverview, &AlertOverviewRequest{})).Methods("GET", "POST")
	protected.HandleFunc("/alerts.updateStatus", h.WithAuth(h.updateAlertStatus, &UpdateAlertStatusRequest{})).Methods("POST")
	protected.HandleFunc("/alerts.bulkUpdate", h.WithAuth(h.bulkUpdateAlerts, &BulkUpdateAlertsRequest{})).Methods("POST")

	// Settings procedures
	protected.HandleFunc("/settings.get", h.GetUserSettings).Methods("GET", "POST")
	protected.HandleFunc("/settings.update", h.UpdateUserSettings).Methods("POST")
	protected.HandleFunc("/settings.regenerateAPIKey", h.RegenerateAPIKey).Methods("POST")

	// Integration procedures
	protected.HandleFunc("/integrations.list", h.ListIntegrations).Methods("GET", "POST")
	protected.HandleFunc("/integrations.get", h.GetIntegration).Methods("GET", "POST")
	protected.HandleFunc("/integrations.create", h.CreateIntegration).Methods("POST")
	protected.HandleFunc("/integrations.update", h.UpdateIntegration).Methods("POST")
	protected.HandleFunc("/integrations.delete", h.DeleteIntegration).Methods("POST")
	protected.HandleFunc("/integrations.test", h.TestIntegration).Methods("POST")

	// Billing procedures
	protected.HandleFunc("/billing.get", h.GetBillingInfo).Methods("GET", "POST")

	return router
}

// Helper functions
func (h *Handler) writeTRPCResponse(w http.ResponseWriter, data interface{}) {
	response := TRPCResponse{
		Result: &TRPCResult{
			Type: "data",
			Data: data,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) writeErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	response := TRPCResponse{
		Error: &TRPCError{
			Code:    statusCode,
			Message: message,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// Basic handlers
func (h *Handler) rootHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Bidding Analysis API"))
}

func (h *Handler) healthCheck(w http.ResponseWriter, r *http.Request) {
	health := map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now(),
		"service":   "bidding-analysis-api",
	}
	h.writeTRPCResponse(w, health)
}

func (h *Handler) debugEndpoint(w http.ResponseWriter, r *http.Request) {
	debug := map[string]interface{}{
		"status":      "running",
		"version":     "1.0.0",
		"environment": "development",
		"endpoints": map[string][]string{
			"auth":      {"login", "register", "me"},
			"bidding":   {"submit", "predict", "stream", "processBid"},
			"campaign":  {"getStats", "getBidHistory", "list", "create"},
			"analytics": {"getFraudAlerts", "getModelAccuracy", "getDashboardMetrics"},
		},
	}
	h.writeTRPCResponse(w, debug)
}

// Campaign handlers
func (h *Handler) listCampaigns(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	if userID == "" {
		h.writeErrorResponse(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	campaigns, err := h.campaignStore.GetUserCampaigns(r.Context(), userID)
	if err != nil {
		log.Printf("Failed to get campaigns: %v", err)
		h.writeErrorResponse(w, "Failed to retrieve campaigns", http.StatusInternalServerError)
		return
	}

	h.writeTRPCResponse(w, campaigns)
}

func (h *Handler) createCampaign(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	if userID == "" {
		h.writeErrorResponse(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	var req struct {
		Name        string   `json:"name"`
		Budget      *float64 `json:"budget,omitempty"`
		DailyBudget *float64 `json:"daily_budget,omitempty"`
		TargetCPA   *float64 `json:"target_cpa,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	campaign := &models.Campaign{
		Name:        req.Name,
		UserID:      uuid.MustParse(userID),
		Status:      "active",
		Budget:      req.Budget,
		DailyBudget: req.DailyBudget,
		TargetCPA:   req.TargetCPA,
	}

	if err := h.campaignStore.CreateCampaign(campaign); err != nil {
		log.Printf("Failed to create campaign: %v", err)
		h.writeErrorResponse(w, "Failed to create campaign", http.StatusInternalServerError)
		return
	}

	h.writeTRPCResponse(w, campaign)
}

// Analytics handlers
func (h *Handler) getDashboardMetrics(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	if userID == "" {
		h.writeErrorResponse(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	// Mock dashboard metrics - replace with actual database queries
	metrics := map[string]interface{}{
		"total_campaigns": 8,
		"active_bids":     1247,
		"win_rate":        0.348,
		"avg_bid":         2.34,
		"total_spend":     12543.67,
		"conversions":     89,
		"fraud_alerts":    2,
		"model_accuracy":  0.92,
		"last_updated":    time.Now(),
	}

	h.writeTRPCResponse(w, metrics)
}

func (h *Handler) processBid(w http.ResponseWriter, r *http.Request) {
	// Legacy endpoint - redirect to new submit endpoint
	h.handleSubmitBid(w, r)
}

func (h *Handler) getCampaignStats(w http.ResponseWriter, r *http.Request) {
	// Mock implementation
	stats := map[string]interface{}{
		"total_bids":  1500,
		"won_bids":    522,
		"win_rate":    0.348,
		"total_spend": 4567.89,
		"conversions": 45,
		"avg_cpa":     101.51,
	}
	h.writeTRPCResponse(w, stats)
}

func (h *Handler) getBidHistory(w http.ResponseWriter, r *http.Request) {
	bids, err := h.bidStore.GetRecentBids(r.Context(), 20)
	if err != nil {
		h.writeErrorResponse(w, "Failed to get bid history", http.StatusInternalServerError)
		return
	}
	h.writeTRPCResponse(w, bids)
}

func (h *Handler) getFraudAlerts(w http.ResponseWriter, r *http.Request) {
	// Mock fraud alerts
	alerts := []map[string]interface{}{
		{
			"id":          "alert-1",
			"type":        "suspicious_click_velocity",
			"severity":    "medium",
			"campaign_id": "campaign-123",
			"detected_at": time.Now().Add(-2 * time.Hour),
			"status":      "active",
		},
		{
			"id":          "alert-2",
			"type":        "geographic_anomaly",
			"severity":    "high",
			"campaign_id": "campaign-456",
			"detected_at": time.Now().Add(-1 * time.Hour),
			"status":      "investigating",
		},
	}
	h.writeTRPCResponse(w, alerts)
}

func (h *Handler) getModelAccuracy(w http.ResponseWriter, r *http.Request) {
	// Mock model accuracy metrics
	accuracy := map[string]interface{}{
		"current_accuracy":    0.924,
		"last_week_accuracy":  0.918,
		"trend":               "improving",
		"total_predictions":   15420,
		"correct_predictions": 14248,
		"last_updated":        time.Now(),
	}
	h.writeTRPCResponse(w, accuracy)
}

// BidRequest represents an incoming bid request
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

// BidResponse represents the response after processing a bid
type BidSubmitResponse struct {
	BidEventID     string  `json:"bid_event_id"`
	Status         string  `json:"status"`
	ProcessedAt    string  `json:"processed_at"`
	WinProbability float64 `json:"win_probability"`
	Message        string  `json:"message"`
}

// BidPredictionRequest for AI-powered bid predictions
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

// BidPredictionResponse contains AI prediction results
type BidPredictionResponse struct {
	PredictedBid float64 `json:"predicted_bid"`
	Confidence   float64 `json:"confidence"`
	Strategy     string  `json:"strategy"`
	FraudRisk    bool    `json:"fraud_risk"`
	Reasoning    string  `json:"reasoning"`
}

// BidStreamData for real-time bid updates
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

// handleSubmitBid processes a new bid submission
func (h *Handler) handleSubmitBid(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received bid submission request")

	var req BidSubmitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode bid request: %v", err)
		h.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.CampaignID == "" {
		h.writeErrorResponse(w, "Campaign ID is required", http.StatusBadRequest)
		return
	}
	if req.BidPrice <= 0 {
		h.writeErrorResponse(w, "Bid price must be positive", http.StatusBadRequest)
		return
	}
	if req.FloorPrice <= 0 {
		h.writeErrorResponse(w, "Floor price must be positive", http.StatusBadRequest)
		return
	}

	// Parse campaign ID
	campaignID, err := uuid.Parse(req.CampaignID)
	if err != nil {
		h.writeErrorResponse(w, "Invalid campaign ID format", http.StatusBadRequest)
		return
	}

	// Create bid event
	bidEvent := &models.BidEvent{
		CampaignID:      campaignID,
		UserID:          req.UserID,
		BidPrice:        req.BidPrice,
		FloorPrice:      req.FloorPrice,
		Won:             false, // Will be determined by auction logic
		Converted:       false,
		SegmentID:       req.SegmentID,
		SegmentCategory: req.SegmentCategory,
		Country:         req.Country,
		Region:          req.Region,
		City:            req.City,
		DeviceType:      req.DeviceType,
		OS:              req.OS,
		Browser:         req.Browser,
		IsMobile:        req.DeviceType == "mobile",
		Timestamp:       time.Now(),
	}

	// Convert keywords
	if len(req.Keywords) > 0 {
		bidEvent.Keywords = req.Keywords
	}

	// Simulate auction logic (you can enhance this with real auction logic)
	winProbability := calculateWinProbability(req.BidPrice, req.FloorPrice)
	bidEvent.Won = winProbability > 0.5

	if bidEvent.Won {
		// Calculate win price (typically 85-95% of bid price)
		winPrice := req.BidPrice * (0.85 + 0.1*(winProbability-0.5)*2)
		bidEvent.WinPrice = &winPrice
	}

	// Store the bid event
	if err := h.bidStore.StoreBidEvent(bidEvent); err != nil {
		log.Printf("Failed to store bid event: %v", err)
		h.writeErrorResponse(w, "Failed to process bid", http.StatusInternalServerError)
		return
	}

	// Prepare response
	response := BidSubmitResponse{
		BidEventID:     bidEvent.ID.String(),
		Status:         getAuctionStatus(bidEvent.Won),
		ProcessedAt:    bidEvent.Timestamp.Format(time.RFC3339),
		WinProbability: winProbability,
		Message:        generateBidMessage(bidEvent.Won, req.BidPrice, bidEvent.WinPrice),
	}

	log.Printf("Bid processed successfully: %s, Won: %t, Win Probability: %.2f",
		bidEvent.ID.String(), bidEvent.Won, winProbability)

	h.writeTRPCResponse(w, response)
}

// handlePredictBid uses AI to predict optimal bid price
func (h *Handler) handlePredictBid(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received bid prediction request")

	var req BidPredictionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode prediction request: %v", err)
		h.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.CampaignID == "" || req.FloorPrice <= 0 {
		h.writeErrorResponse(w, "Campaign ID and floor price are required", http.StatusBadRequest)
		return
	}

	// Create ML prediction request
	mlRequest := &models.BidRequest{
		CampaignID: uuid.MustParse(req.CampaignID),
		UserSegment: models.UserSegment{
			SegmentID:             req.UserSegment,
			Category:              req.UserSegment,
			EngagementScore:       req.EngagementScore,
			ConversionProbability: req.ConversionProb,
		},
		GeoLocation: models.GeoLocation{
			Country: req.Country,
		},
		DeviceInfo: models.DeviceInfo{
			DeviceType: req.DeviceType,
			IsMobile:   req.DeviceType == "mobile",
		},
		FloorPrice: req.FloorPrice,
		Keywords:   req.Keywords,
		Timestamp:  time.Now(),
	}

	// Get AI prediction
	prediction, err := h.predictor.PredictBid(r.Context(), mlRequest)
	if err != nil {
		log.Printf("ML prediction failed: %v", err)
		// Fallback to rule-based prediction
		prediction = h.getFallbackPrediction(req)
	}

	response := BidPredictionResponse{
		PredictedBid: prediction.BidPrice,
		Confidence:   prediction.Confidence,
		Strategy:     prediction.Strategy,
		FraudRisk:    prediction.FraudRisk,
		Reasoning:    generatePredictionReasoning(prediction, req),
	}

	log.Printf("Bid prediction completed: $%.4f (confidence: %.2f)",
		response.PredictedBid, response.Confidence)

	h.writeTRPCResponse(w, response)
}

// handleGetBidStream returns recent bid activity
func (h *Handler) handleGetBidStream(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received bid stream request")

	// Get recent bids (last 50)
	bids, err := h.bidStore.GetRecentBids(r.Context(), 50)
	if err != nil {
		log.Printf("Failed to get recent bids: %v", err)
		h.writeErrorResponse(w, "Failed to retrieve bid stream", http.StatusInternalServerError)
		return
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

	h.writeTRPCResponse(w, streamData)
}

// Helper functions

func calculateWinProbability(bidPrice, floorPrice float64) float64 {
	// Simple win probability calculation based on bid vs floor price
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

func generatePredictionReasoning(prediction *models.BidResponse, req BidPredictionRequest) string {
	confidence := "medium"
	if prediction.Confidence > 0.8 {
		confidence = "high"
	} else if prediction.Confidence < 0.4 {
		confidence = "low"
	}

	return fmt.Sprintf("AI recommends $%.4f based on %s confidence analysis of user segment '%s' and device type '%s'",
		prediction.BidPrice, confidence, req.UserSegment, req.DeviceType)
}

func (h *Handler) getFallbackPrediction(req BidPredictionRequest) *models.BidResponse {
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

// getCampaign retrieves a single campaign with detailed metrics
func (h *Handler) getCampaign(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	if userID == "" {
		h.writeErrorResponse(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	// Get campaign ID from query parameter or body
	campaignID := r.URL.Query().Get("id")
	if campaignID == "" {
		var req struct {
			ID string `json:"id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err == nil {
			campaignID = req.ID
		}
	}

	if campaignID == "" {
		h.writeErrorResponse(w, "Campaign ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(campaignID)
	if err != nil {
		h.writeErrorResponse(w, "Invalid campaign ID format", http.StatusBadRequest)
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		h.writeErrorResponse(w, "Invalid user ID", http.StatusInternalServerError)
		return
	}

	campaign, err := h.campaignStore.GetCampaignWithMetrics(r.Context(), id, userUUID)
	if err != nil {
		log.Printf("Failed to get campaign: %v", err)
		h.writeErrorResponse(w, "Failed to retrieve campaign", http.StatusInternalServerError)
		return
	}

	h.writeTRPCResponse(w, campaign)
}

// updateCampaign updates an existing campaign
func (h *Handler) updateCampaign(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	if userID == "" {
		h.writeErrorResponse(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	var req struct {
		ID          string   `json:"id"`
		Name        string   `json:"name"`
		Status      string   `json:"status,omitempty"`
		Budget      *float64 `json:"budget,omitempty"`
		DailyBudget *float64 `json:"daily_budget,omitempty"`
		TargetCPA   *float64 `json:"target_cpa,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.ID == "" {
		h.writeErrorResponse(w, "Campaign ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(req.ID)
	if err != nil {
		h.writeErrorResponse(w, "Invalid campaign ID format", http.StatusBadRequest)
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		h.writeErrorResponse(w, "Invalid user ID", http.StatusInternalServerError)
		return
	}

	// Validate inputs
	if req.Name != "" && (len(req.Name) < 3 || len(req.Name) > 255) {
		h.writeErrorResponse(w, "Campaign name must be between 3 and 255 characters", http.StatusBadRequest)
		return
	}

	if req.Status != "" && req.Status != "active" && req.Status != "paused" && req.Status != "archived" {
		h.writeErrorResponse(w, "Invalid status. Must be 'active', 'paused', or 'archived'", http.StatusBadRequest)
		return
	}

	if req.Budget != nil && *req.Budget < 0 {
		h.writeErrorResponse(w, "Budget must be positive", http.StatusBadRequest)
		return
	}

	if req.DailyBudget != nil && *req.DailyBudget < 0 {
		h.writeErrorResponse(w, "Daily budget must be positive", http.StatusBadRequest)
		return
	}

	if req.Budget != nil && req.DailyBudget != nil && *req.DailyBudget > *req.Budget {
		h.writeErrorResponse(w, "Daily budget cannot exceed total budget", http.StatusBadRequest)
		return
	}

	// Get existing campaign first
	existing, err := h.campaignStore.GetCampaign(id)
	if err != nil {
		h.writeErrorResponse(w, "Campaign not found", http.StatusNotFound)
		return
	}

	// Verify ownership
	if existing.UserID != userUUID {
		h.writeErrorResponse(w, "Unauthorized", http.StatusForbidden)
		return
	}

	// Update fields
	if req.Name != "" {
		existing.Name = req.Name
	}
	if req.Status != "" {
		existing.Status = req.Status
	}
	if req.Budget != nil {
		existing.Budget = req.Budget
	}
	if req.DailyBudget != nil {
		existing.DailyBudget = req.DailyBudget
	}
	if req.TargetCPA != nil {
		existing.TargetCPA = req.TargetCPA
	}

	if err := h.campaignStore.UpdateCampaign(existing); err != nil {
		log.Printf("Failed to update campaign: %v", err)
		h.writeErrorResponse(w, "Failed to update campaign", http.StatusInternalServerError)
		return
	}

	h.writeTRPCResponse(w, existing)
}

// deleteCampaign soft deletes a campaign
func (h *Handler) deleteCampaign(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	if userID == "" {
		h.writeErrorResponse(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	var req struct {
		ID string `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.ID == "" {
		h.writeErrorResponse(w, "Campaign ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(req.ID)
	if err != nil {
		h.writeErrorResponse(w, "Invalid campaign ID format", http.StatusBadRequest)
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		h.writeErrorResponse(w, "Invalid user ID", http.StatusInternalServerError)
		return
	}

	if err := h.campaignStore.DeleteCampaign(id, userUUID); err != nil {
		log.Printf("Failed to delete campaign: %v", err)
		h.writeErrorResponse(w, "Failed to delete campaign", http.StatusInternalServerError)
		return
	}

	h.writeTRPCResponse(w, map[string]interface{}{
		"success": true,
		"message": "Campaign archived successfully",
	})
}

// pauseCampaign pauses an active campaign
func (h *Handler) pauseCampaign(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	if userID == "" {
		h.writeErrorResponse(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	var req struct {
		ID string `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.ID == "" {
		h.writeErrorResponse(w, "Campaign ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(req.ID)
	if err != nil {
		h.writeErrorResponse(w, "Invalid campaign ID format", http.StatusBadRequest)
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		h.writeErrorResponse(w, "Invalid user ID", http.StatusInternalServerError)
		return
	}

	if err := h.campaignStore.PauseCampaign(id, userUUID); err != nil {
		log.Printf("Failed to pause campaign: %v", err)
		h.writeErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return updated campaign
	campaign, err := h.campaignStore.GetCampaign(id)
	if err != nil {
		h.writeErrorResponse(w, "Campaign paused but failed to retrieve", http.StatusInternalServerError)
		return
	}

	h.writeTRPCResponse(w, campaign)
}

// activateCampaign activates a paused campaign
func (h *Handler) activateCampaign(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	if userID == "" {
		h.writeErrorResponse(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	var req struct {
		ID string `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.ID == "" {
		h.writeErrorResponse(w, "Campaign ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(req.ID)
	if err != nil {
		h.writeErrorResponse(w, "Invalid campaign ID format", http.StatusBadRequest)
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		h.writeErrorResponse(w, "Invalid user ID", http.StatusInternalServerError)
		return
	}

	if err := h.campaignStore.ActivateCampaign(id, userUUID); err != nil {
		log.Printf("Failed to activate campaign: %v", err)
		h.writeErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return updated campaign
	campaign, err := h.campaignStore.GetCampaign(id)
	if err != nil {
		h.writeErrorResponse(w, "Campaign activated but failed to retrieve", http.StatusInternalServerError)
		return
	}

	h.writeTRPCResponse(w, campaign)
}

// listCampaignsEnhanced lists campaigns with metrics (enhanced version)
func (h *Handler) listCampaignsEnhanced(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	if userID == "" {
		h.writeErrorResponse(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		h.writeErrorResponse(w, "Invalid user ID", http.StatusInternalServerError)
		return
	}

	// Parse pagination parameters (optional)
	var limit = 100 // default
	var offset = 0  // default

	if r.URL.Query().Get("limit") != "" {
		// Parse limit from query param if provided
		// For simplicity, using defaults here
	}

	campaigns, err := h.campaignStore.ListCampaignsWithMetrics(r.Context(), userUUID, limit, offset)
	if err != nil {
		log.Printf("Failed to list campaigns: %v", err)
		h.writeErrorResponse(w, "Failed to retrieve campaigns", http.StatusInternalServerError)
		return
	}

	h.writeTRPCResponse(w, campaigns)
}

// getDailyMetrics retrieves daily metrics for a campaign
func (h *Handler) getDailyMetrics(w http.ResponseWriter, r *http.Request) {
	userID := GetUserIDFromContext(r.Context())
	if userID == "" {
		h.writeErrorResponse(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	campaignID := r.URL.Query().Get("id")
	if campaignID == "" {
		var req struct {
			ID string `json:"id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err == nil {
			campaignID = req.ID
		}
	}

	if campaignID == "" {
		h.writeErrorResponse(w, "Campaign ID is required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(campaignID)
	if err != nil {
		h.writeErrorResponse(w, "Invalid campaign ID format", http.StatusBadRequest)
		return
	}

	// Verify ownership
	userUUID, _ := uuid.Parse(userID)
	campaign, err := h.campaignStore.GetCampaign(id)
	if err != nil || campaign.UserID != userUUID {
		h.writeErrorResponse(w, "Campaign not found or unauthorized", http.StatusNotFound)
		return
	}

	// Get metrics for last 30 days by default
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -30)

	metrics, err := h.campaignStore.GetCampaignDailyMetrics(id, startDate, endDate)
	if err != nil {
		log.Printf("Failed to get daily metrics: %v", err)
		h.writeErrorResponse(w, "Failed to retrieve metrics", http.StatusInternalServerError)
		return
	}

	h.writeTRPCResponse(w, metrics)
}
