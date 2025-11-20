// Package mlonnx provides ONNX-based ML inference for bid optimization
package mlonnx

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/baskint/bidding-analysis/internal/models"
	"github.com/google/uuid"
	"github.com/yalue/onnxruntime_go"
)

// ONNXPredictor handles ML predictions using ONNX models
type ONNXPredictor struct {
	session         *onnxruntime_go.AdvancedSession
	featureCount    int
	featureEncoders map[string]map[string]float32
	mu              sync.RWMutex
	modelVersion    string
	loadedAt        time.Time
}

// FeatureEncoders stores categorical encoding mappings
type FeatureEncoders struct {
	DeviceType      map[string]float32 `json:"device_type"`
	SegmentCategory map[string]float32 `json:"segment_category"`
	Country         map[string]float32 `json:"country"`
}

// NewONNXPredictor creates a new ONNX-based predictor
func NewONNXPredictor(modelPath, encodersPath string) (*ONNXPredictor, error) {
	// Initialize ONNX Runtime
	err := onnxruntime_go.InitializeEnvironment()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize ONNX: %w", err)
	}

	// Load feature encoders
	encoders, err := loadEncoders(encodersPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load encoders: %w", err)
	}

	// Feature count must match training
	featureCount := 13

	// Create dummy input for session initialization
	dummyInput := make([]float32, featureCount)

	// Create ONNX session
	session, err := onnxruntime_go.NewAdvancedSession(
		modelPath,
		[]string{"float_input"},
		[]string{"output"},
		dummyInput,
		[]int64{1, int64(featureCount)},
		nil,
	)
	if err != nil {
		onnxruntime_go.DestroyEnvironment()
		return nil, fmt.Errorf("failed to load ONNX model: %w", err)
	}

	return &ONNXPredictor{
		session:         session,
		featureCount:    featureCount,
		featureEncoders: encoders,
		modelVersion:    extractVersionFromPath(modelPath),
		loadedAt:        time.Now(),
	}, nil
}

// PredictBidPrice predicts optimal bid price for a request
func (p *ONNXPredictor) PredictBidPrice(
	ctx context.Context,
	req *models.BidRequest,
	historicalData []*models.BidEvent,
) (*models.BidResponse, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	// Extract features
	features, err := p.extractFeatures(req, historicalData)
	if err != nil {
		return nil, fmt.Errorf("feature extraction failed: %w", err)
	}

	// Set input tensor
	inputTensor := p.session.GetInputTensor(0)
	inputData := inputTensor.GetData()
	copy(inputData, features)

	// Run inference
	startTime := time.Now()
	err = p.session.Run()
	inferenceTime := time.Since(startTime)

	if err != nil {
		return nil, fmt.Errorf("inference failed: %w", err)
	}

	// Get prediction
	outputTensor := p.session.GetOutputTensor(0)
	prediction := outputTensor.GetData()[0]

	// Build response
	response := &models.BidResponse{
		BidPrice:     float64(prediction),
		Confidence:   p.calculateConfidence(features, historicalData),
		Strategy:     fmt.Sprintf("ml_onnx_%s", p.modelVersion),
		FraudRisk:    false, // Can be enhanced with separate fraud model
		PredictionID: uuid.New().String(),
	}

	// Validate prediction
	response = p.validatePrediction(response, req)

	// Log inference time for monitoring
	if inferenceTime > 10*time.Millisecond {
		fmt.Printf("Warning: Slow inference time: %v\n", inferenceTime)
	}

	return response, nil
}

// PredictBatch performs batch predictions for multiple requests
func (p *ONNXPredictor) PredictBatch(
	ctx context.Context,
	requests []*models.BidRequest,
	historicalDataMap map[uuid.UUID][]*models.BidEvent,
) ([]*models.BidResponse, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	batchSize := len(requests)
	if batchSize == 0 {
		return nil, fmt.Errorf("empty batch")
	}

	// Extract features for all requests
	allFeatures := make([]float32, batchSize*p.featureCount)
	for i, req := range requests {
		historical := historicalDataMap[req.CampaignID]
		features, err := p.extractFeatures(req, historical)
		if err != nil {
			return nil, fmt.Errorf("batch feature extraction failed at index %d: %w", i, err)
		}
		copy(allFeatures[i*p.featureCount:], features)
	}

	// Set input tensor
	inputTensor := p.session.GetInputTensor(0)
	copy(inputTensor.GetData(), allFeatures)

	// Run batch inference
	startTime := time.Now()
	err := p.session.Run()
	if err != nil {
		return nil, fmt.Errorf("batch inference failed: %w", err)
	}
	inferenceTime := time.Since(startTime)

	// Get predictions
	outputTensor := p.session.GetOutputTensor(0)
	predictions := outputTensor.GetData()

	// Build responses
	responses := make([]*models.BidResponse, batchSize)
	for i := 0; i < batchSize; i++ {
		historical := historicalDataMap[requests[i].CampaignID]
		responses[i] = &models.BidResponse{
			BidPrice:     float64(predictions[i]),
			Confidence:   p.calculateConfidence(nil, historical),
			Strategy:     fmt.Sprintf("ml_onnx_batch_%s", p.modelVersion),
			FraudRisk:    false,
			PredictionID: uuid.New().String(),
		}
		responses[i] = p.validatePrediction(responses[i], requests[i])
	}

	fmt.Printf("Batch inference: %d requests in %v (%.2f ms/request)\n",
		batchSize, inferenceTime, float64(inferenceTime.Milliseconds())/float64(batchSize))

	return responses, nil
}

// extractFeatures converts bid request to feature vector
func (p *ONNXPredictor) extractFeatures(
	req *models.BidRequest,
	historical []*models.BidEvent,
) ([]float32, error) {
	features := make([]float32, p.featureCount)

	// Feature 0: floor_price
	features[0] = float32(req.FloorPrice)

	// Feature 1: engagement_score
	features[1] = float32(req.UserSegment.EngagementScore)
	if features[1] == 0 {
		features[1] = 0.5 // Default
	}

	// Feature 2: conversion_probability
	features[2] = float32(req.UserSegment.ConversionProbability)
	if features[2] == 0 {
		features[2] = 0.05 // Default
	}

	// Features 3-5: Historical metrics
	if len(historical) > 0 {
		var wonCount, totalBid, totalWinPrice float32
		var winCount int
		for _, bid := range historical {
			if bid.Won {
				wonCount++
				if bid.WinPrice != nil {
					totalWinPrice += float32(*bid.WinPrice)
					winCount++
				}
			}
			totalBid += float32(bid.BidPrice)
		}
		features[3] = wonCount / float32(len(historical))                   // win_rate
		features[4] = totalBid / float32(len(historical))                   // avg_bid
		if winCount > 0 {
			features[5] = totalWinPrice / float32(winCount) // avg_win_price
		} else {
			features[5] = features[0] // Default to floor price
		}
	} else {
		features[3] = 0.3                // Default win rate
		features[4] = float32(req.FloorPrice) // Default avg bid
		features[5] = float32(req.FloorPrice) // Default avg win price
	}

	// Feature 6: device_type_encoded
	features[6] = p.encodeCategory("device_type", req.DeviceInfo.DeviceType)

	// Feature 7: segment_category_encoded
	features[7] = p.encodeCategory("segment_category", req.UserSegment.Category)

	// Feature 8: hour_of_day
	features[8] = float32(req.Timestamp.Hour())

	// Feature 9: day_of_week
	features[9] = float32(req.Timestamp.Weekday())

	// Feature 10: country_encoded
	features[10] = p.encodeCategory("country", req.GeoLocation.Country)

	// Features 11-12: Campaign metrics (would need to be queried)
	// For now, using defaults
	features[11] = 0 // campaign_spend_last_7d (needs DB query)
	features[12] = 0 // campaign_conversions_last_7d (needs DB query)

	return features, nil
}

// encodeCategory encodes categorical variable using learned encoders
func (p *ONNXPredictor) encodeCategory(category, value string) float32 {
	if encoder, exists := p.featureEncoders[category]; exists {
		if encoded, found := encoder[value]; found {
			return encoded
		}
	}
	return 0 // Unknown category
}

// calculateConfidence computes prediction confidence
func (p *ONNXPredictor) calculateConfidence(
	features []float32,
	historical []*models.BidEvent,
) float64 {
	// Base confidence
	confidence := 0.85

	// Reduce confidence if little historical data
	if len(historical) < 10 {
		confidence -= 0.1
	}
	if len(historical) < 5 {
		confidence -= 0.1
	}

	// Could add more sophisticated confidence estimation:
	// - Model uncertainty from ensemble predictions
	// - Feature-based uncertainty
	// - Historical prediction accuracy

	return confidence
}

// validatePrediction ensures prediction is within reasonable bounds
func (p *ONNXPredictor) validatePrediction(
	prediction *models.BidResponse,
	req *models.BidRequest,
) *models.BidResponse {
	// Ensure bid is at least floor price
	if prediction.BidPrice < req.FloorPrice {
		prediction.BidPrice = req.FloorPrice * 1.05
		prediction.Confidence *= 0.9
	}

	// Cap maximum bid at reasonable level
	maxBid := req.FloorPrice * 5.0
	if prediction.BidPrice > maxBid {
		prediction.BidPrice = maxBid
		prediction.Confidence *= 0.8
	}

	// Ensure confidence is between 0 and 1
	if prediction.Confidence > 1.0 {
		prediction.Confidence = 1.0
	}
	if prediction.Confidence < 0.0 {
		prediction.Confidence = 0.0
	}

	return prediction
}

// Close releases ONNX resources
func (p *ONNXPredictor) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.session != nil {
		p.session.Destroy()
		p.session = nil
	}
	onnxruntime_go.DestroyEnvironment()
	return nil
}

// GetModelInfo returns model metadata
func (p *ONNXPredictor) GetModelInfo() map[string]interface{} {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return map[string]interface{}{
		"model_version": p.modelVersion,
		"loaded_at":     p.loadedAt,
		"feature_count": p.featureCount,
		"uptime":        time.Since(p.loadedAt).String(),
	}
}

// Helper functions

func loadEncoders(path string) (map[string]map[string]float32, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var encoders FeatureEncoders
	if err := json.Unmarshal(data, &encoders); err != nil {
		return nil, err
	}

	return map[string]map[string]float32{
		"device_type":      encoders.DeviceType,
		"segment_category": encoders.SegmentCategory,
		"country":          encoders.Country,
	}, nil
}

func extractVersionFromPath(path string) string {
	// Extract version from path like "models/bid_optimizer_20241201_120000.onnx"
	// Returns "20241201_120000"
	base := filepath.Base(path)
	parts := strings.Split(base, "_")
	if len(parts) >= 3 {
		version := strings.TrimSuffix(parts[len(parts)-2]+"_"+parts[len(parts)-1], ".onnx")
		return version
	}
	return "unknown"
}

// AnalyzeAudienceSegment implements the AIClient interface
func (p *ONNXPredictor) AnalyzeAudienceSegment(
	ctx context.Context,
	bidEvents []*models.BidEvent,
) (*AudienceAnalysis, error) {
	// This would require a separate ONNX model for audience analysis
	// For now, return a placeholder
	return &AudienceAnalysis{
		Segments: []string{"premium", "standard"},
		Insights: []string{"Audience analysis requires separate model"},
	}, nil
}

// DetectFraud implements the AIClient interface
func (p *ONNXPredictor) DetectFraud(
	ctx context.Context,
	bidEvents []*models.BidEvent,
) (*FraudAnalysis, error) {
	// This would require a separate ONNX model for fraud detection
	// For now, return a placeholder
	return &FraudAnalysis{
		FraudDetected: false,
		Confidence:    0.5,
		Patterns:      []string{},
		Severity:      0,
	}, nil
}

// Supporting types (matching predictor.go)
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
