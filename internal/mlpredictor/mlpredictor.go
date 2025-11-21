// Package mlpredictor provides ML-based bid prediction using XGBoost models
package mlpredictor

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/dmitryikh/leaves"
)

// BidPredictor handles ML-based bid predictions
type BidPredictor struct {
	model        *leaves.Ensemble
	encoders     map[string]map[string]float64
	mu           sync.RWMutex
	modelPath    string
	encodersPath string
	lastReload   time.Time
}

// BidFeatures represents the input features for bid prediction
type BidFeatures struct {
	// Core features
	FloorPrice            float64
	EngagementScore       float64
	ConversionProbability float64

	// Historical features
	HistoricalWinRate     float64
	HistoricalAvgBid      float64
	HistoricalAvgWinPrice float64

	// Categorical features (will be encoded)
	DeviceType      string
	SegmentCategory string
	Country         string

	// Time features
	HourOfDay int
	DayOfWeek int

	// Campaign features
	CampaignSpendLast7d       float64
	CampaignConversionsLast7d float64
}

// NewBidPredictor creates a new bid predictor
func NewBidPredictor(modelPath, encodersPath string) (*BidPredictor, error) {
	p := &BidPredictor{
		modelPath:    modelPath,
		encodersPath: encodersPath,
	}

	if err := p.LoadModel(); err != nil {
		return nil, fmt.Errorf("failed to load model: %w", err)
	}

	return p, nil
}

// LoadModel loads the XGBoost model and feature encoders
func (p *BidPredictor) LoadModel() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Load XGBoost model (JSON format from Python)
	model, err := leaves.XGEnsembleFromFile(p.modelPath, true)
	if err != nil {
		return fmt.Errorf("failed to load XGBoost model: %w", err)
	}
	p.model = model

	// Load feature encoders
	encodersData, err := os.ReadFile(p.encodersPath)
	if err != nil {
		return fmt.Errorf("failed to read encoders: %w", err)
	}

	if err := json.Unmarshal(encodersData, &p.encoders); err != nil {
		return fmt.Errorf("failed to parse encoders: %w", err)
	}

	p.lastReload = time.Now()

	return nil
}

// Predict returns the optimal bid for given features
func (p *BidPredictor) Predict(features BidFeatures) (float64, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.model == nil {
		return 0, fmt.Errorf("model not loaded")
	}

	// Convert features to float array in correct order
	// Order must match training: see FEATURE_COLUMNS in train_model.py
	featureVector := []float64{
		features.FloorPrice,                                           // 0: floor_price
		features.EngagementScore,                                      // 1: engagement_score
		features.ConversionProbability,                                // 2: conversion_probability
		features.HistoricalWinRate,                                    // 3: historical_win_rate
		features.HistoricalAvgBid,                                     // 4: historical_avg_bid
		features.HistoricalAvgWinPrice,                                // 5: historical_avg_win_price
		p.encodeFeature("device_type", features.DeviceType),           // 6: device_type_encoded
		p.encodeFeature("segment_category", features.SegmentCategory), // 7: segment_category_encoded
		float64(features.HourOfDay),                                   // 8: hour_of_day
		float64(features.DayOfWeek),                                   // 9: day_of_week
		p.encodeFeature("country", features.Country),                  // 10: country_encoded
		features.CampaignSpendLast7d,                                  // 11: campaign_spend_last_7d
		features.CampaignConversionsLast7d,                            // 12: campaign_conversions_last_7d
	}

	// Make prediction
	prediction := p.model.PredictSingle(featureVector, 0)

	// Ensure prediction is at least above floor price
	if prediction < features.FloorPrice {
		prediction = features.FloorPrice * 1.01 // 1% above floor
	}

	return prediction, nil
}

// PredictBatch makes predictions for multiple bid requests efficiently
func (p *BidPredictor) PredictBatch(batch []BidFeatures) ([]float64, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.model == nil {
		return nil, fmt.Errorf("model not loaded")
	}

	predictions := make([]float64, len(batch))

	for i, features := range batch {
		featureVector := []float64{
			features.FloorPrice,
			features.EngagementScore,
			features.ConversionProbability,
			features.HistoricalWinRate,
			features.HistoricalAvgBid,
			features.HistoricalAvgWinPrice,
			p.encodeFeature("device_type", features.DeviceType),
			p.encodeFeature("segment_category", features.SegmentCategory),
			float64(features.HourOfDay),
			float64(features.DayOfWeek),
			p.encodeFeature("country", features.Country),
			features.CampaignSpendLast7d,
			features.CampaignConversionsLast7d,
		}

		prediction := p.model.PredictSingle(featureVector, 0)

		// Ensure above floor
		if prediction < features.FloorPrice {
			prediction = features.FloorPrice * 1.01
		}

		predictions[i] = prediction
	}

	return predictions, nil
}

// encodeFeature encodes a categorical feature using frequency encoding
func (p *BidPredictor) encodeFeature(featureName, value string) float64 {
	if encoder, ok := p.encoders[featureName]; ok {
		if encoded, ok := encoder[value]; ok {
			return encoded
		}
	}
	// Return 0 for unknown categories (will be handled gracefully)
	return 0.0
}

// GetModelInfo returns information about the loaded model
func (p *BidPredictor) GetModelInfo() map[string]interface{} {
	p.mu.RLock()
	defer p.mu.RUnlock()

	info := map[string]interface{}{
		"model_path":    p.modelPath,
		"encoders_path": p.encodersPath,
		"last_reload":   p.lastReload,
		"model_loaded":  p.model != nil,
	}

	if p.model != nil {
		info["n_features"] = p.model.NFeatures()
		info["n_outputs"] = p.model.NOutputGroups()
	}

	return info
}

// ReloadModel reloads the model from disk (useful for hot-swapping)
func (p *BidPredictor) ReloadModel() error {
	return p.LoadModel()
}

// Close cleans up resources
func (p *BidPredictor) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.model = nil
	p.encoders = nil

	return nil
}
