package trpc

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"

	"github.com/baskint/bidding-analysis/internal/ml"
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

	// Dashboard metrics - Using WithAuth wrapper
	protected.HandleFunc("/dashboard.metrics", h.WithAuthNoBody(h.getDashboardMetrics)).Methods("GET")
	protected.HandleFunc("/campaign.stats", h.WithAuthNoBody(h.getCampaignStats)).Methods("GET")
	protected.HandleFunc("/bid.history", h.WithAuthNoBody(h.getBidHistory)).Methods("GET")
	protected.HandleFunc("/fraud.alerts", h.WithAuthNoBody(h.getFraudAlerts)).Methods("GET")
	protected.HandleFunc("/model.accuracy", h.WithAuthNoBody(h.getModelAccuracy)).Methods("GET")

	// Add these routes (they can call the same handlers)
	protected.HandleFunc("/analytics.getDashboardMetrics", h.WithAuthNoBody(h.getDashboardMetrics)).Methods("GET", "POST")
	protected.HandleFunc("/campaign.getStats", h.WithAuthNoBody(h.getCampaignStats)).Methods("GET", "POST")
	protected.HandleFunc("/analytics.getFraudAlerts", h.WithAuthNoBody(h.getFraudAlerts)).Methods("GET", "POST")

	// Bidding procedures - Using WithAuth wrapper
	protected.HandleFunc("/bidding.submit", h.WithAuth(h.handleSubmitBid, &BidSubmitRequest{})).Methods("POST")
	protected.HandleFunc("/bidding.predict", h.WithAuth(h.handlePredictBid, &BidPredictionRequest{})).Methods("POST")
	protected.HandleFunc("/bidding.stream", h.WithAuthNoBody(h.handleGetBidStream)).Methods("GET", "POST")
	protected.HandleFunc("/bid.process", h.WithAuth(h.processBid, &BidSubmitRequest{})).Methods("POST")

	// Campaign procedures - Using WithAuth wrapper
	protected.HandleFunc("/campaign.list", h.WithAuthNoBody(h.listCampaigns)).Methods("GET", "POST")
	protected.HandleFunc("/campaign.create", h.WithAuth(h.createCampaign, &CreateCampaignRequest{})).Methods("POST")
	protected.HandleFunc("/campaign.get", h.WithAuthQuery(h.getCampaign)).Methods("GET", "POST")
	protected.HandleFunc("/campaign.update", h.WithAuth(h.updateCampaign, &UpdateCampaignRequest{})).Methods("POST")
	protected.HandleFunc("/campaign.delete", h.WithAuth(h.deleteCampaign, &DeleteCampaignRequest{})).Methods("POST")
	protected.HandleFunc("/campaign.pause", h.WithAuth(h.pauseCampaign, &PauseCampaignRequest{})).Methods("POST")
	protected.HandleFunc("/campaign.activate", h.WithAuth(h.activateCampaign, &ActivateCampaignRequest{})).Methods("POST")
	protected.HandleFunc("/campaign.listWithMetrics", h.WithAuthNoBody(h.listCampaignsEnhanced)).Methods("GET", "POST")
	protected.HandleFunc("/campaign.getDailyMetrics", h.WithAuthQuery(h.getDailyMetrics)).Methods("GET", "POST")

	// Analytics page procedures - Using WithAuth wrapper
	protected.HandleFunc("/analytics.getPerformanceOverview", h.WithAuth(h.getPerformanceOverview, &DateRangeRequest{})).Methods("GET", "POST")
	protected.HandleFunc("/analytics.getKeywordAnalysis", h.WithAuth(h.getKeywordAnalysis, &KeywordAnalysisRequest{})).Methods("GET", "POST")
	protected.HandleFunc("/analytics.getDeviceBreakdown", h.WithAuth(h.getDeviceBreakdown, &DateRangeRequest{})).Methods("GET", "POST")
	protected.HandleFunc("/analytics.getGeoBreakdown", h.WithAuth(h.getGeoBreakdown, &DateRangeRequest{})).Methods("GET", "POST")
	protected.HandleFunc("/analytics.getHourlyPerformance", h.WithAuth(h.getHourlyPerformance, &DateRangeRequest{})).Methods("GET", "POST")
	protected.HandleFunc("/analytics.getDailyTrends", h.WithAuth(h.getDailyTrends, &DateRangeRequest{})).Methods("GET", "POST")
	protected.HandleFunc("/analytics.getCompetitiveAnalysis", h.WithAuth(h.getCompetitiveAnalysis, &CompetitiveAnalysisRequest{})).Methods("GET", "POST")
	// protected.HandleFunc("/analytics.getCampaignComparison", h.getCampaignComparison).Methods("GET", "POST") // TODO: Refactor later

	// ML Model procedures - Using WithAuthQuery for GET endpoints
	protected.HandleFunc("/mlModel.list", h.WithAuthQuery(h.listMLModels)).Methods("GET", "POST")
	protected.HandleFunc("/mlModel.get", h.WithAuthQuery(h.getMLModel)).Methods("GET", "POST")
	protected.HandleFunc("/mlModel.create", h.WithAuth(h.createMLModel, &CreateMLModelRequest{})).Methods("POST")
	protected.HandleFunc("/mlModel.update", h.WithAuth(h.updateMLModel, &UpdateMLModelRequest{})).Methods("POST")
	protected.HandleFunc("/mlModel.delete", h.WithAuth(h.deleteMLModel, &DeleteMLModelRequest{})).Methods("POST")
	protected.HandleFunc("/mlModel.setDefault", h.WithAuth(h.setDefaultMLModel, &SetDefaultMLModelRequest{})).Methods("POST")
	protected.HandleFunc("/mlModel.getDefault", h.WithAuthQuery(h.getDefaultMLModel)).Methods("GET", "POST")

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

	// Settings and Integrations - Using WithAuth wrapper
	protected.HandleFunc("/settings.get", h.WithAuthNoBody(h.GetUserSettings)).Methods("GET", "POST")
	protected.HandleFunc("/settings.update", h.WithAuth(h.UpdateUserSettings, &UpdateUserSettingsRequest{})).Methods("POST")
	protected.HandleFunc("/settings.regenerateAPIKey", h.WithAuthNoBody(h.RegenerateAPIKey)).Methods("POST")
	protected.HandleFunc("/integrations.list", h.WithAuthNoBody(h.ListIntegrations)).Methods("GET", "POST")
	protected.HandleFunc("/integrations.get", h.WithAuth(h.GetIntegration, &GetIntegrationRequest{})).Methods("GET", "POST")
	protected.HandleFunc("/integrations.create", h.WithAuth(h.CreateIntegration, &CreateIntegrationRequest{})).Methods("POST")
	protected.HandleFunc("/integrations.update", h.WithAuth(h.UpdateIntegration, &UpdateIntegrationRequest{})).Methods("POST")
	protected.HandleFunc("/integrations.delete", h.WithAuth(h.DeleteIntegration, &DeleteIntegrationRequest{})).Methods("POST")
	protected.HandleFunc("/integrations.test", h.WithAuth(h.TestIntegration, &TestIntegrationRequest{})).Methods("POST")
	protected.HandleFunc("/billing.get", h.WithAuthNoBody(h.GetBillingInfo)).Methods("GET", "POST")

	// Billing procedures
	protected.HandleFunc("/billing.get", h.WithAuthNoBody(h.GetBillingInfo)).Methods("GET", "POST")

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
