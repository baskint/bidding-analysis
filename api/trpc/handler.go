package trpc

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/baskint/bidding-analysis/internal/ml"
	"github.com/baskint/bidding-analysis/internal/store"
)

// Handler contains the dependencies for tRPC handlers
type Handler struct {
	bidStore      *store.BidStore
	campaignStore *store.CampaignStore
	userStore     *store.UserStore
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
	// You'll need to get the database connection from one of the existing stores
	db := bidStore.DB() // You'll need to add this method to BidStore
	userStore := store.NewUserStore(db)

	return &Handler{
		bidStore:      bidStore,
		campaignStore: campaignStore,
		userStore:     userStore,
		predictor:     predictor,
		jwtSecret:     jwtSecret,
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

	// Auth procedures (public - no authentication required)
	api.HandleFunc("/auth.login", h.login).Methods("POST")
	api.HandleFunc("/auth.register", h.register).Methods("POST")

	// Protected routes (require authentication)
	protected := api.PathPrefix("").Subrouter()
	protected.Use(h.authMiddleware)

	// Auth procedures (protected)
	protected.HandleFunc("/auth.me", h.getMe).Methods("GET", "POST")

	// Bidding procedures (now protected)
	protected.HandleFunc("/bidding.processBid", h.processBid).Methods("POST")

	// Campaign procedures (now protected)
	protected.HandleFunc("/campaign.getStats", h.getCampaignStats).Methods("GET", "POST")
	protected.HandleFunc("/campaign.getBidHistory", h.getBidHistory).Methods("GET", "POST")

	// Analytics procedures (now protected)
	protected.HandleFunc("/analytics.getFraudAlerts", h.getFraudAlerts).Methods("GET", "POST")
	protected.HandleFunc("/analytics.getModelAccuracy", h.getModelAccuracy).Methods("GET", "POST")

	return router
}
