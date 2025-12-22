package mlpredictor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

// BidFeatures represents input features for prediction
type BidFeatures struct {
	FloorPrice                float64
	EngagementScore           float64
	ConversionProbability     float64
	HistoricalWinRate         float64
	HistoricalAvgBid          float64
	HistoricalAvgWinPrice     float64
	DeviceType                string
	SegmentCategory           string
	Country                   string
	HourOfDay                 int
	DayOfWeek                 int
	CampaignSpendLast7d       float64
	CampaignConversionsLast7d float64
}

// Predictor interface
type Predictor interface {
	Predict(features BidFeatures) (float64, error)
	PredictBatch(batch []BidFeatures) ([]float64, error)
	GetModelInfo() map[string]interface{}
	ReloadModel() error
	Close() error
}

// BidPredictorHTTP calls Python ML service via HTTP
type BidPredictorHTTP struct {
	serviceURL string
	client     *http.Client
	mu         sync.RWMutex
}

// NewBidPredictorHTTP creates a predictor that calls Python service
func NewBidPredictorHTTP(serviceURL string) (Predictor, error) {
	if serviceURL == "" {
		serviceURL = os.Getenv("ML_SERVICE_URL")
		if serviceURL == "" {
			serviceURL = "http://localhost:5001"
		}
	}

	// increase timeout to 30s to allow for complex model loading
	p := &BidPredictorHTTP{
		serviceURL: serviceURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	// Check if service is available
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := p.healthCheck(ctx); err != nil {
		return nil, fmt.Errorf("ML service not available at %s: %w", serviceURL, err)
	}

	return p, nil
}

// Predict returns the optimal bid
func (p *BidPredictorHTTP) Predict(features BidFeatures) (float64, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	// Wrap features in "features" key as expected by Python service
	reqData := map[string]interface{}{
		"features": map[string]interface{}{
			"floor_price":                  features.FloorPrice,
			"engagement_score":             features.EngagementScore,
			"conversion_probability":       features.ConversionProbability,
			"historical_win_rate":          features.HistoricalWinRate,
			"historical_avg_bid":           features.HistoricalAvgBid,
			"historical_avg_win_price":     features.HistoricalAvgWinPrice,
			"device_type":                  features.DeviceType,
			"segment_category":             features.SegmentCategory,
			"hour_of_day":                  features.HourOfDay,
			"day_of_week":                  features.DayOfWeek,
			"country":                      features.Country,
			"campaign_spend_last_7d":       features.CampaignSpendLast7d,
			"campaign_conversions_last_7d": features.CampaignConversionsLast7d,
		},
	}

	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal request: %w", err)
	}

	// increase this to 30s to allow for complex model inference
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", p.serviceURL+"/predict", bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to call ML service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("ML service returned status %d", resp.StatusCode)
	}

	// Python service returns {"predicted_bid": X, "model_version": Y}
	var result struct {
		PredictedBid float64 `json:"predicted_bid"`
		ModelVersion string  `json:"model_version"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.PredictedBid, nil
}

// PredictBatch makes predictions for multiple requests
func (p *BidPredictorHTTP) PredictBatch(batch []BidFeatures) ([]float64, error) {
	predictions := make([]float64, len(batch))
	for i, features := range batch {
		pred, err := p.Predict(features)
		if err != nil {
			return nil, fmt.Errorf("batch prediction failed at index %d: %w", i, err)
		}
		predictions[i] = pred
	}
	return predictions, nil
}

// healthCheck verifies the ML service is running
func (p *BidPredictorHTTP) healthCheck(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", p.serviceURL+"/health", nil)
	if err != nil {
		return err
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check failed with status %d", resp.StatusCode)
	}

	return nil
}

// GetModelInfo returns information about the ML service
func (p *BidPredictorHTTP) GetModelInfo() map[string]interface{} {
	return map[string]interface{}{
		"service_url":  p.serviceURL,
		"model_type":   "http_python_service",
		"model_loaded": true,
	}
}

// ReloadModel is a no-op for HTTP predictor
func (p *BidPredictorHTTP) ReloadModel() error {
	return nil
}

// Close cleans up resources
func (p *BidPredictorHTTP) Close() error {
	p.client.CloseIdleConnections()
	return nil
}
