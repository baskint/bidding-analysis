package trpc

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/baskint/bidding-analysis/internal/store"
)

// Handler contains the dependencies for tRPC handlers
type Handler struct {
	bidStore      *store.BidStore
	campaignStore *store.CampaignStore
}

// NewHandler creates a new tRPC Handler instance
func NewHandler(bidStore *store.BidStore, campaignStore *store.CampaignStore) *Handler {
	return &Handler{
		bidStore:      bidStore,
		campaignStore: campaignStore,
	}
}

// TRPCResponse represents a tRPC response structure
type TRPCResponse struct {
	Result *TRPCResult `json:"result,omitempty"`
	Error  *TRPCError  `json:"error,omitempty"`
}

// TRPCResult represents successful tRPC result
type TRPCResult struct {
	Data interface{} `json:"data"`
	Type string      `json:"type"`
}

// TRPCError represents tRPC error
type TRPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// SetupRoutes configures all tRPC routes
func (h *Handler) SetupRoutes() http.Handler {
	router := mux.NewRouter()

	// CORS middleware
	router.Use(corsMiddleware)

	// Health check
	router.HandleFunc("/health", h.healthCheck).Methods("GET")

	// tRPC routes
	api := router.PathPrefix("/trpc").Subrouter()

	// Campaign procedures
	api.HandleFunc("/campaign.getStats", h.getCampaignStats).Methods("GET", "POST")
	api.HandleFunc("/campaign.getBidHistory", h.getBidHistory).Methods("GET", "POST")

	// Analytics procedures
	api.HandleFunc("/analytics.getFraudAlerts", h.getFraudAlerts).Methods("GET", "POST")
	api.HandleFunc("/analytics.getModelAccuracy", h.getModelAccuracy).Methods("GET", "POST")

	return router
}

// getCampaignStats handles campaign statistics requests
func (h *Handler) getCampaignStats(w http.ResponseWriter, r *http.Request) {
	input := struct {
		CampaignID string `json:"campaignId"`
		StartTime  string `json:"startTime"`
		EndTime    string `json:"endTime"`
	}{}

	if err := parseInput(r, &input); err != nil {
		writeError(w, 400, "Invalid input", err)
		return
	}

	campaignID, err := uuid.Parse(input.CampaignID)
	if err != nil {
		writeError(w, 400, "Invalid campaign ID", err)
		return
	}

	startTime, err := time.Parse("2006-01-02", input.StartTime)
	if err != nil {
		startTime = time.Now().AddDate(0, 0, -7)
	}

	endTime, err := time.Parse("2006-01-02", input.EndTime)
	if err != nil {
		endTime = time.Now()
	}

	stats, err := h.campaignStore.GetCampaignStats(campaignID, startTime, endTime)
	if err != nil {
		writeError(w, 500, "Failed to get campaign stats", err)
		return
	}

	writeSuccess(w, stats)
}

// getBidHistory handles bid history requests
func (h *Handler) getBidHistory(w http.ResponseWriter, r *http.Request) {
	input := struct {
		CampaignID string `json:"campaignId"`
		StartTime  string `json:"startTime"`
		EndTime    string `json:"endTime"`
		Limit      int    `json:"limit"`
		Offset     int    `json:"offset"`
	}{}

	if err := parseInput(r, &input); err != nil {
		writeError(w, 400, "Invalid input", err)
		return
	}

	if input.Limit <= 0 {
		input.Limit = 100
	}
	if input.Offset < 0 {
		input.Offset = 0
	}

	startTime, err := time.Parse("2006-01-02", input.StartTime)
	if err != nil {
		startTime = time.Now().AddDate(0, 0, -1)
	}

	endTime, err := time.Parse("2006-01-02", input.EndTime)
	if err != nil {
		endTime = time.Now()
	}

	bids, err := h.bidStore.GetBidHistory(input.CampaignID, startTime, endTime, input.Limit, input.Offset)
	if err != nil {
		writeError(w, 500, "Failed to get bid history", err)
		return
	}

	result := map[string]interface{}{
		"bids":   bids,
		"limit":  input.Limit,
		"offset": input.Offset,
	}

	writeSuccess(w, result)
}

// getFraudAlerts handles fraud alerts requests
func (h *Handler) getFraudAlerts(w http.ResponseWriter, r *http.Request) {
	input := struct {
		StartTime         string `json:"startTime"`
		EndTime           string `json:"endTime"`
		SeverityThreshold int    `json:"severityThreshold"`
	}{}

	if err := parseInput(r, &input); err != nil {
		writeError(w, 400, "Invalid input", err)
		return
	}

	startTime, err := time.Parse("2006-01-02", input.StartTime)
	if err != nil {
		startTime = time.Now().AddDate(0, 0, -7)
	}

	endTime, err := time.Parse("2006-01-02", input.EndTime)
	if err != nil {
		endTime = time.Now()
	}

	if input.SeverityThreshold < 1 || input.SeverityThreshold > 10 {
		input.SeverityThreshold = 5
	}

	alerts, err := h.campaignStore.GetFraudAlerts(startTime, endTime, input.SeverityThreshold)
	if err != nil {
		writeError(w, 500, "Failed to get fraud alerts", err)
		return
	}

	writeSuccess(w, alerts)
}

// getModelAccuracy handles model accuracy requests
func (h *Handler) getModelAccuracy(w http.ResponseWriter, r *http.Request) {
	input := struct {
		StartTime    string `json:"startTime"`
		EndTime      string `json:"endTime"`
		ModelVersion string `json:"modelVersion"`
	}{}

	if err := parseInput(r, &input); err != nil {
		writeError(w, 400, "Invalid input", err)
		return
	}

	startTime, err := time.Parse("2006-01-02", input.StartTime)
	if err != nil {
		startTime = time.Now().AddDate(0, 0, -7)
	}

	endTime, err := time.Parse("2006-01-02", input.EndTime)
	if err != nil {
		endTime = time.Now()
	}

	metrics, err := h.campaignStore.GetModelAccuracy(startTime, endTime, input.ModelVersion)
	if err != nil {
		writeError(w, 500, "Failed to get model accuracy", err)
		return
	}

	writeSuccess(w, metrics)
}

// healthCheck handles health check requests
func (h *Handler) healthCheck(w http.ResponseWriter, r *http.Request) {
	result := map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now(),
	}
	writeSuccess(w, result)
}

// Helper functions

func parseInput(r *http.Request, input interface{}) error {
	if r.Method == "GET" {
		// Parse query parameters for GET requests
		return parseQueryParams(r, input)
	}

	// Parse JSON body for POST requests
	return json.NewDecoder(r.Body).Decode(input)
}

func parseQueryParams(r *http.Request, input interface{}) error {
	// This is a simplified implementation
	// In a real tRPC setup, you'd parse the query parameters properly
	query := r.URL.Query()

	// Convert to JSON and back for simplicity
	jsonData := make(map[string]interface{})
	for key, values := range query {
		if len(values) > 0 {
			// Try to parse as int, fallback to string
			if intVal, err := strconv.Atoi(values[0]); err == nil {
				jsonData[key] = intVal
			} else {
				jsonData[key] = values[0]
			}
		}
	}

	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		return err
	}

	return json.Unmarshal(jsonBytes, input)
}

func writeSuccess(w http.ResponseWriter, data interface{}) {
	response := TRPCResponse{
		Result: &TRPCResult{
			Data: data,
			Type: "data",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func writeError(w http.ResponseWriter, code int, message string, err error) {
	response := TRPCResponse{
		Error: &TRPCError{
			Code:    code,
			Message: message,
			Data:    err.Error(),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
