package trpc

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/baskint/bidding-analysis/internal/ml"
	"github.com/baskint/bidding-analysis/internal/store"
)

// Handler contains the dependencies for tRPC handlers
type Handler struct {
	bidStore      *store.BidStore
	campaignStore *store.CampaignStore
	predictor     *ml.Predictor
}

// NewHandler creates a new tRPC Handler instance
func NewHandler(bidStore *store.BidStore, campaignStore *store.CampaignStore, predictor *ml.Predictor) *Handler {
	return &Handler{
		bidStore:      bidStore,
		campaignStore: campaignStore,
		predictor:     predictor,
	}
}

// SetupRoutes configures all tRPC routes
func (h *Handler) SetupRoutes() http.Handler {
	router := mux.NewRouter()

	// Apply middleware
	router.Use(corsMiddleware)
	router.Use(loggingMiddleware)

	// Health check
	router.HandleFunc("/health", h.healthCheck).Methods("GET")

	// tRPC routes
	api := router.PathPrefix("/trpc").Subrouter()

	// Debug endpoint
	api.HandleFunc("/debug", h.debugEndpoint).Methods("GET")

	// Bidding procedures
	api.HandleFunc("/bidding.processBid", h.processBid).Methods("POST")

	// Campaign procedures
	api.HandleFunc("/campaign.getStats", h.getCampaignStats).Methods("GET", "POST")
	api.HandleFunc("/campaign.getBidHistory", h.getBidHistory).Methods("GET", "POST")

	// Analytics procedures
	api.HandleFunc("/analytics.getFraudAlerts", h.getFraudAlerts).Methods("GET", "POST")
	api.HandleFunc("/analytics.getModelAccuracy", h.getModelAccuracy).Methods("GET", "POST")

	return router
}
