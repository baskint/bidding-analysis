package mlpredictor

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/owulveryck/onnx-go"
	"github.com/owulveryck/onnx-go/backend/x/gorgonnx"
	"gorgonia.org/tensor"
)

// BidPredictorONNXSimple uses pure Go ONNX library (no C dependencies)
type BidPredictorONNXSimple struct {
	backend      *gorgonnx.Graph
	encoders     map[string]map[string]float64
	mu           sync.RWMutex
	modelPath    string
	encodersPath string
}

// NewBidPredictorONNXSimple creates predictor using pure Go ONNX
func NewBidPredictorONNXSimple(modelPath, encodersPath string) (*BidPredictorONNXSimple, error) {
	p := &BidPredictorONNXSimple{
		modelPath:    modelPath,
		encodersPath: encodersPath,
	}

	if err := p.LoadModel(); err != nil {
		return nil, fmt.Errorf("failed to load model: %w", err)
	}

	return p, nil
}

// LoadModel loads the ONNX model and feature encoders
func (p *BidPredictorONNXSimple) LoadModel() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Load ONNX model
	backend := gorgonnx.NewGraph()
	model := onnx.NewModel(backend)

	modelFile, err := os.Open(p.modelPath)
	if err != nil {
		return fmt.Errorf("failed to open model file: %w", err)
	}
	defer modelFile.Close()

	err = model.UnmarshalBinary(modelFile)
	if err != nil {
		return fmt.Errorf("failed to load ONNX model: %w", err)
	}

	p.backend = backend

	// Load feature encoders
	encodersData, err := os.ReadFile(p.encodersPath)
	if err != nil {
		return fmt.Errorf("failed to read encoders: %w", err)
	}

	if err := json.Unmarshal(encodersData, &p.encoders); err != nil {
		return fmt.Errorf("failed to parse encoders: %w", err)
	}

	return nil
}

// Predict returns the optimal bid for given features
func (p *BidPredictorONNXSimple) Predict(features BidFeatures) (float64, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.backend == nil {
		return 0, fmt.Errorf("model not loaded")
	}

	// Convert features to float32 array
	inputData := []float32{
		float32(features.FloorPrice),
		float32(features.EngagementScore),
		float32(features.ConversionProbability),
		float32(features.HistoricalWinRate),
		float32(features.HistoricalAvgBid),
		float32(features.HistoricalAvgWinPrice),
		float32(p.encodeFeature("device_type", features.DeviceType)),
		float32(p.encodeFeature("segment_category", features.SegmentCategory)),
		float32(features.HourOfDay),
		float32(features.DayOfWeek),
		float32(p.encodeFeature("country", features.Country)),
		float32(features.CampaignSpendLast7d),
		float32(features.CampaignConversionsLast7d),
	}

	// Create input tensor
	inputTensor := tensor.New(
		tensor.WithShape(1, 13),
		tensor.WithBacking(inputData),
	)

	// Set input
	err := p.backend.SetInput(0, inputTensor)
	if err != nil {
		return 0, fmt.Errorf("failed to set input: %w", err)
	}

	// Run inference
	err = p.backend.Run()
	if err != nil {
		return 0, fmt.Errorf("inference failed: %w", err)
	}

	// Get output
	output, err := p.backend.GetOutputTensors()
	if err != nil {
		return 0, fmt.Errorf("failed to get output: %w", err)
	}

	if len(output) == 0 {
		return 0, fmt.Errorf("no output from model")
	}

	// Extract prediction
	outputData := output[0].Data().([]float32)
	if len(outputData) == 0 {
		return 0, fmt.Errorf("empty output")
	}

	prediction := float64(outputData[0])

	// Ensure above floor price
	if prediction < features.FloorPrice {
		prediction = features.FloorPrice * 1.01
	}

	return prediction, nil
}

// PredictBatch makes predictions for multiple bid requests
func (p *BidPredictorONNXSimple) PredictBatch(batch []BidFeatures) ([]float64, error) {
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

// encodeFeature encodes a categorical feature
func (p *BidPredictorONNXSimple) encodeFeature(featureName, value string) float64 {
	if encoder, ok := p.encoders[featureName]; ok {
		if encoded, ok := encoder[value]; ok {
			return encoded
		}
	}
	return 0.0
}

// GetModelInfo returns model information
func (p *BidPredictorONNXSimple) GetModelInfo() map[string]interface{} {
	return map[string]interface{}{
		"model_path":    p.modelPath,
		"encoders_path": p.encodersPath,
		"model_loaded":  p.backend != nil,
		"model_type":    "onnx_pure_go",
	}
}

// ReloadModel reloads the model
func (p *BidPredictorONNXSimple) ReloadModel() error {
	return p.LoadModel()
}

// Close cleans up resources
func (p *BidPredictorONNXSimple) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.backend = nil
	return nil
}
