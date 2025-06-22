package trpc

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/baskint/bidding-analysis/internal/models"
)

// healthCheck handles health check requests
func (h *Handler) healthCheck(w http.ResponseWriter, r *http.Request) {
	result := map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now(),
	}
	writeSuccess(w, result)
}

// debugEndpoint handles debug requests
func (h *Handler) debugEndpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("DEBUG: /debug endpoint hit!\n")
	writeSuccess(w, map[string]string{"message": "debug endpoint working"})
}

// processBid handles bid processing requests
func (h *Handler) processBid(w http.ResponseWriter, r *http.Request) {
	var input ProcessBidInput

	if err := parseInput(r, &input); err != nil {
		writeError(w, 400, "Invalid input", err)
		return
	}

	// Parse campaign ID
	campaignID, err := uuid.Parse(input.CampaignID)
	if err != nil {
		writeError(w, 400, "Invalid campaign ID", err)
		return
	}

	// Build bid request
	bidRequest := &models.BidRequest{
		CampaignID: campaignID,
		UserID:     input.UserID,
		UserSegment: models.UserSegment{
			SegmentID:             input.SegmentID,
			Category:              input.SegmentCategory,
			EngagementScore:       input.EngagementScore,
			ConversionProbability: input.ConversionProbability,
		},
		GeoLocation: models.GeoLocation{
			Country: input.Country,
			Region:  input.Region,
			City:    input.City,
		},
		DeviceInfo: models.DeviceInfo{
			DeviceType: input.DeviceType,
			OS:         input.OS,
			Browser:    input.Browser,
			IsMobile:   input.DeviceType == "mobile",
		},
		FloorPrice: input.FloorPrice,
		Keywords:   input.Keywords,
		Timestamp:  time.Now(),
	}

	// Get prediction
	prediction, err := h.predictor.PredictOptimalBid(r.Context(), bidRequest)
	if err != nil {
		writeError(w, 500, "Failed to get bid prediction", err)
		return
	}

	writeSuccess(w, prediction)
}

// getCampaignStats handles campaign statistics requests
func (h *Handler) getCampaignStats(w http.ResponseWriter, r *http.Request) {
	campaignID := r.URL.Query().Get("campaignId")
	startTime := r.URL.Query().Get("startTime")
	endTime := r.URL.Query().Get("endTime")

	if campaignID == "" {
		writeError(w, 400, "Missing campaignId parameter", nil)
		return
	}

	// Parse campaign ID
	id, err := uuid.Parse(campaignID)
	if err != nil {
		writeError(w, 400, "Invalid campaign ID", err)
		return
	}

	// Parse time range
	var start, end time.Time
	if startTime != "" {
		start, err = time.Parse("2006-01-02", startTime)
		if err != nil {
			writeError(w, 400, "Invalid startTime format (use YYYY-MM-DD)", err)
			return
		}
	} else {
		start = time.Now().AddDate(0, 0, -30) // Default to last 30 days
	}

	if endTime != "" {
		end, err = time.Parse("2006-01-02", endTime)
		if err != nil {
			writeError(w, 400, "Invalid endTime format (use YYYY-MM-DD)", err)
			return
		}
	} else {
		end = time.Now()
	}

	// Get campaign statistics using the ML predictor
	stats, err := h.predictor.AnalyzeCampaignPerformance(r.Context(), id, int(end.Sub(start).Hours()/24))
	if err != nil {
		writeError(w, 500, "Failed to get campaign statistics", err)
		return
	}

	writeSuccess(w, stats)
}

// getBidHistory handles bid history requests
func (h *Handler) getBidHistory(w http.ResponseWriter, r *http.Request) {
	var input BidHistoryInput

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
	var input FraudAlertsInput

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
	var input ModelAccuracyInput

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
