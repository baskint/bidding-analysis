# AI/ML Bid Optimization: Analysis & Recommendations

## Current System Analysis

### Architecture Overview
Your platform is a Go-based bid optimization system with:
- **Backend**: Go (gRPC + tRPC APIs)
- **Frontend**: Next.js/React
- **Database**: PostgreSQL
- **Current ML**: OpenAI API calls for predictions

### Current ML Implementation

**Strengths:**
- Clean interface abstraction (`AIClient`)
- Fallback to rule-based predictions
- Historical data integration
- Campaign performance tracking

**Limitations:**
1. **No actual model training** - relies on OpenAI's general-purpose LLM
2. **No fine-tuning** - each prediction is zero-shot
3. **Expensive** - API calls for every bid prediction
4. **Latency** - Network round-trips for real-time bidding
5. **No model ownership** - dependent on external service
6. **Limited customization** - can't optimize for your specific data patterns

---

## Machine Learning Options in Go

### 1. **GoLearn** (Traditional ML)
**Repository**: `github.com/sjwhitworth/golearn`

**Best For**: Classical ML algorithms (Random Forests, Decision Trees, Linear Regression)

```go
import (
    "github.com/sjwhitworth/golearn/base"
    "github.com/sjwhitworth/golearn/ensemble"
    "github.com/sjwhitworth/golearn/evaluation"
)

// Example: Random Forest for bid prediction
func trainBidModel(trainingData base.FixedDataGrid) *ensemble.RandomForest {
    rf := ensemble.NewRandomForest(100, 3) // 100 trees, 3 features per split
    rf.Fit(trainingData)
    return rf
}
```

**Pros:**
- Pure Go implementation
- Fast inference
- Good for tabular data (perfect for bid features)
- No external dependencies

**Cons:**
- Limited algorithm selection
- Less active development
- No deep learning

### 2. **Gorgonia** (Deep Learning)
**Repository**: `gorgonia.org/gorgonia`

**Best For**: Neural networks, gradient-based optimization

```go
import (
    "gorgonia.org/gorgonia"
    "gorgonia.org/tensor"
)

// Example: Neural network for bid prediction
type BidNN struct {
    g      *gorgonia.ExprGraph
    w1, w2 *gorgonia.Node // weights
    b1, b2 *gorgonia.Node // biases
    pred   *gorgonia.Node // prediction output
}

func (m *BidNN) Forward(input *gorgonia.Node) *gorgonia.Node {
    h := gorgonia.Must(gorgonia.Add(
        gorgonia.Must(gorgonia.Mul(input, m.w1)),
        m.b1,
    ))
    h = gorgonia.Must(gorgonia.Rectify(h))
    
    pred := gorgonia.Must(gorgonia.Add(
        gorgonia.Must(gorgonia.Mul(h, m.w2)),
        m.b2,
    ))
    return pred
}
```

**Pros:**
- Full neural network capabilities
- Automatic differentiation
- GPU support
- Pure Go

**Cons:**
- Steeper learning curve
- More complex to implement
- Less documentation than Python alternatives

### 3. **TensorFlow Go Bindings**
**Repository**: `github.com/tensorflow/tensorflow/go`

**Best For**: Using pre-trained TensorFlow models, production inference

```go
import (
    tf "github.com/tensorflow/tensorflow/tensorflow/go"
)

// Example: Load and use SavedModel
func loadModel(modelPath string) (*tf.SavedModel, error) {
    model, err := tf.LoadSavedModel(
        modelPath,
        []string{"serve"},
        nil,
    )
    return model, err
}

func predict(model *tf.SavedModel, features map[string]*tf.Tensor) (*tf.Tensor, error) {
    results, err := model.Session.Run(
        features,
        []tf.Output{
            model.Graph.Operation("StatefulPartitionedCall").Output(0),
        },
        nil,
    )
    return results[0], err
}
```

**Pros:**
- Industry standard
- Well-tested inference
- Can use Python-trained models
- Good performance

**Cons:**
- C++ dependency (CGO)
- Training must be done in Python
- Binary size increases
- Complex setup

### 4. **ONNX Runtime Go**
**Repository**: `github.com/yalue/onnxruntime_go`

**Best For**: Model portability, production inference

```go
import "github.com/yalue/onnxruntime_go"

// Example: Run ONNX model
func runONNXModel(modelPath string, inputData []float32) ([]float32, error) {
    // Initialize ONNX Runtime
    onnxruntime_go.InitializeEnvironment()
    defer onnxruntime_go.DestroyEnvironment()
    
    // Load model
    session, err := onnxruntime_go.NewAdvancedSession(
        modelPath,
        []string{"input"},
        []string{"output"},
        inputData,
        []int64{1, len(inputData)},
        nil,
    )
    defer session.Destroy()
    
    // Run inference
    err = session.Run()
    output := session.GetOutputTensor(0)
    
    return output.GetData(), err
}
```

**Pros:**
- Framework-agnostic (PyTorch, TensorFlow, scikit-learn → ONNX)
- Optimized inference
- Cross-platform
- Model portability

**Cons:**
- Still requires Python for training
- CGO dependency
- Limited Go bindings documentation

### 5. **Hybrid Approach: Python Training + Go Inference** ⭐ **RECOMMENDED**

**Architecture:**
```
┌─────────────────────────────────────────┐
│ Training Pipeline (Python)              │
│ - scikit-learn / XGBoost / LightGBM    │
│ - PyTorch / TensorFlow                  │
│ - Feature engineering                   │
│ - Hyperparameter tuning                 │
│ - Model validation                      │
└───────────┬─────────────────────────────┘
            │
            ↓ Export (ONNX/SavedModel/Pickle)
            │
┌───────────┴─────────────────────────────┐
│ Production Inference (Go)               │
│ - Fast, low-latency predictions         │
│ - Model serving via gRPC                │
│ - Feature computation                   │
│ - A/B testing                           │
└─────────────────────────────────────────┘
```

---

## Recommended ML Architecture

### Phase 1: Quick Win - Traditional ML (2-4 weeks)

**1. Use XGBoost/LightGBM (Python Training)**

```python
# training/train_bid_model.py
import xgboost as xgb
import pandas as pd
import numpy as np
from sklearn.model_selection import train_test_split
import onnxmltools
from skl2onnx import convert_sklearn
from skl2onnx.common.data_types import FloatTensorType

class BidOptimizer:
    def __init__(self):
        self.model = None
        self.feature_columns = [
            'floor_price',
            'engagement_score',
            'conversion_probability',
            'historical_win_rate',
            'historical_avg_bid',
            'device_type_encoded',
            'segment_category_encoded',
            'hour_of_day',
            'day_of_week',
            'country_encoded'
        ]
    
    def prepare_features(self, df):
        """Feature engineering"""
        # Time features
        df['hour_of_day'] = pd.to_datetime(df['timestamp']).dt.hour
        df['day_of_week'] = pd.to_datetime(df['timestamp']).dt.dayofweek
        
        # Categorical encoding
        df['device_type_encoded'] = pd.Categorical(df['device_type']).codes
        df['segment_category_encoded'] = pd.Categorical(df['segment_category']).codes
        df['country_encoded'] = pd.Categorical(df['country']).codes
        
        # Historical features (from aggregate tables)
        df = df.merge(self.get_historical_features(), on='campaign_id', how='left')
        
        return df[self.feature_columns]
    
    def train(self, training_data_path):
        """Train the model"""
        df = pd.read_csv(training_data_path)
        
        X = self.prepare_features(df)
        y = df['optimal_bid']  # Target: actual win price or computed optimal bid
        
        X_train, X_val, y_train, y_val = train_test_split(
            X, y, test_size=0.2, random_state=42
        )
        
        # XGBoost model
        self.model = xgb.XGBRegressor(
            n_estimators=200,
            max_depth=6,
            learning_rate=0.1,
            subsample=0.8,
            colsample_bytree=0.8,
            objective='reg:squarederror',
            n_jobs=-1
        )
        
        self.model.fit(
            X_train, y_train,
            eval_set=[(X_val, y_val)],
            early_stopping_rounds=10,
            verbose=True
        )
        
        # Evaluate
        train_score = self.model.score(X_train, y_train)
        val_score = self.model.score(X_val, y_val)
        
        print(f"Train R²: {train_score:.4f}")
        print(f"Val R²: {val_score:.4f}")
        
        return {
            'train_r2': train_score,
            'val_r2': val_score,
            'feature_importance': dict(zip(
                self.feature_columns,
                self.model.feature_importances_
            ))
        }
    
    def export_to_onnx(self, output_path):
        """Export model to ONNX for Go inference"""
        initial_type = [('float_input', FloatTensorType([None, len(self.feature_columns)]))]
        
        onnx_model = convert_sklearn(
            self.model,
            initial_types=initial_type,
            target_opset=12
        )
        
        with open(output_path, "wb") as f:
            f.write(onnx_model.SerializeToString())
        
        print(f"Model exported to {output_path}")

# Usage
optimizer = BidOptimizer()
metrics = optimizer.train('data/training_bids.csv')
optimizer.export_to_onnx('models/bid_optimizer_v1.onnx')
```

**2. Go Inference Service**

```go
// internal/ml/onnx_predictor.go
package ml

import (
    "context"
    "fmt"
    "github.com/yalue/onnxruntime_go"
    "github.com/baskint/bidding-analysis/internal/models"
)

type ONNXPredictor struct {
    session      *onnxruntime_go.AdvancedSession
    featureCount int
}

func NewONNXPredictor(modelPath string) (*ONNXPredictor, error) {
    err := onnxruntime_go.InitializeEnvironment()
    if err != nil {
        return nil, fmt.Errorf("failed to initialize ONNX: %w", err)
    }
    
    // Feature count must match training
    featureCount := 10
    
    // Create dummy input for session initialization
    dummyInput := make([]float32, featureCount)
    
    session, err := onnxruntime_go.NewAdvancedSession(
        modelPath,
        []string{"float_input"},
        []string{"output"},
        dummyInput,
        []int64{1, int64(featureCount)},
        nil,
    )
    if err != nil {
        return nil, fmt.Errorf("failed to load model: %w", err)
    }
    
    return &ONNXPredictor{
        session:      session,
        featureCount: featureCount,
    }, nil
}

func (p *ONNXPredictor) PredictBidPrice(
    ctx context.Context,
    req *models.BidRequest,
    historical []*models.BidEvent,
) (*models.BidResponse, error) {
    // Prepare features
    features := p.extractFeatures(req, historical)
    
    // Set input
    inputTensor := p.session.GetInputTensor(0)
    copy(inputTensor.GetData(), features)
    
    // Run inference
    err := p.session.Run()
    if err != nil {
        return nil, fmt.Errorf("inference failed: %w", err)
    }
    
    // Get prediction
    outputTensor := p.session.GetOutputTensor(0)
    prediction := outputTensor.GetData()[0]
    
    return &models.BidResponse{
        BidPrice:   float64(prediction),
        Confidence: 0.85, // Can be computed from model uncertainty
        Strategy:   "ml_onnx_v1",
        FraudRisk:  false,
    }, nil
}

func (p *ONNXPredictor) extractFeatures(
    req *models.BidRequest,
    historical []*models.BidEvent,
) []float32 {
    features := make([]float32, p.featureCount)
    
    // Feature 0: floor_price
    features[0] = float32(req.FloorPrice)
    
    // Feature 1: engagement_score
    features[1] = float32(req.UserSegment.EngagementScore)
    
    // Feature 2: conversion_probability
    features[2] = float32(req.UserSegment.ConversionProbability)
    
    // Feature 3-4: historical_win_rate, historical_avg_bid
    if len(historical) > 0 {
        var wonCount, totalBid float32
        for _, bid := range historical {
            if bid.Won {
                wonCount++
            }
            totalBid += float32(bid.BidPrice)
        }
        features[3] = wonCount / float32(len(historical))
        features[4] = totalBid / float32(len(historical))
    }
    
    // Feature 5: device_type_encoded
    features[5] = p.encodeDeviceType(req.DeviceInfo.DeviceType)
    
    // Feature 6: segment_category_encoded
    features[6] = p.encodeSegmentCategory(req.UserSegment.Category)
    
    // Feature 7-8: hour_of_day, day_of_week
    features[7] = float32(req.Timestamp.Hour())
    features[8] = float32(req.Timestamp.Weekday())
    
    // Feature 9: country_encoded
    features[9] = p.encodeCountry(req.GeoLocation.Country)
    
    return features
}

func (p *ONNXPredictor) encodeDeviceType(deviceType string) float32 {
    // Must match Python encoding
    mapping := map[string]float32{
        "desktop": 0,
        "mobile":  1,
        "tablet":  2,
    }
    if val, ok := mapping[deviceType]; ok {
        return val
    }
    return -1
}

func (p *ONNXPredictor) encodeSegmentCategory(category string) float32 {
    mapping := map[string]float32{
        "premium":   0,
        "standard":  1,
        "value":     2,
        "new_user":  3,
    }
    if val, ok := mapping[category]; ok {
        return val
    }
    return -1
}

func (p *ONNXPredictor) encodeCountry(country string) float32 {
    // Top countries encoded, others as -1
    mapping := map[string]float32{
        "US": 0, "GB": 1, "CA": 2, "AU": 3, "DE": 4,
        "FR": 5, "JP": 6, "CN": 7, "IN": 8,
    }
    if val, ok := mapping[country]; ok {
        return val
    }
    return -1
}

func (p *ONNXPredictor) Close() error {
    if p.session != nil {
        p.session.Destroy()
    }
    onnxruntime_go.DestroyEnvironment()
    return nil
}
```

**3. Integration into Existing System**

```go
// internal/ml/predictor.go - Update NewPredictor
func NewPredictor(apiKey string, modelPath string, bidStore *store.BidStore) *Predictor {
    var client AIClient
    
    // Try ONNX model first
    if modelPath != "" {
        onnxPredictor, err := NewONNXPredictor(modelPath)
        if err == nil {
            client = onnxPredictor
            log.Printf("Using ONNX model: %s", modelPath)
        } else {
            log.Printf("Failed to load ONNX model: %v, falling back", err)
        }
    }
    
    // Fallback to OpenAI
    if client == nil {
        if apiKey == "" || apiKey == "your_openai_key" {
            client = NewMockOpenAIClient()
        } else {
            client = NewOpenAIClient(apiKey)
        }
    }
    
    return &Predictor{
        openaiClient: client,
        bidStore:     bidStore,
        modelVersion: "v1.0.0",
    }
}
```

---

## Training Data Pipeline

### 1. Data Collection Schema

```sql
-- Create training data table
CREATE TABLE bid_training_data (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    campaign_id UUID NOT NULL,
    
    -- Input features
    floor_price DECIMAL(10,4) NOT NULL,
    engagement_score DECIMAL(5,4),
    conversion_probability DECIMAL(5,4),
    device_type VARCHAR(50),
    segment_category VARCHAR(100),
    country VARCHAR(3),
    region VARCHAR(100),
    hour_of_day INTEGER,
    day_of_week INTEGER,
    
    -- Historical context features
    historical_win_rate DECIMAL(5,4),
    historical_avg_bid DECIMAL(10,4),
    historical_avg_win_price DECIMAL(10,4),
    campaign_spend_last_7d DECIMAL(12,2),
    campaign_conversions_last_7d INTEGER,
    
    -- Target variable
    optimal_bid DECIMAL(10,4) NOT NULL,
    actual_outcome VARCHAR(20), -- 'won', 'lost', 'converted'
    
    -- Metadata
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    training_set VARCHAR(50), -- 'train', 'validation', 'test'
    
    FOREIGN KEY (campaign_id) REFERENCES campaigns(id)
);

CREATE INDEX idx_training_data_campaign ON bid_training_data(campaign_id);
CREATE INDEX idx_training_data_created ON bid_training_data(created_at DESC);
CREATE INDEX idx_training_data_training_set ON bid_training_data(training_set);
```

### 2. Data Generation Service

```go
// cmd/training-data-generator/main.go
package main

import (
    "context"
    "log"
    "time"
    
    "github.com/baskint/bidding-analysis/internal/store"
    "github.com/google/uuid"
)

type TrainingDataGenerator struct {
    bidStore      *store.BidStore
    campaignStore *store.CampaignStore
    db            *sqlx.DB
}

func (g *TrainingDataGenerator) GenerateTrainingData(
    ctx context.Context,
    startDate, endDate time.Time,
) error {
    log.Printf("Generating training data from %v to %v", startDate, endDate)
    
    // Get all campaigns
    campaigns, err := g.campaignStore.ListAll(ctx)
    if err != nil {
        return err
    }
    
    for _, campaign := range campaigns {
        log.Printf("Processing campaign: %s", campaign.Name)
        
        // Get bid events
        bidEvents, err := g.bidStore.GetCampaignBids(
            ctx,
            campaign.ID,
            startDate,
        )
        if err != nil {
            continue
        }
        
        // Transform to training examples
        for _, bid := range bidEvents {
            trainingExample := g.createTrainingExample(ctx, bid, campaign.ID)
            
            // Insert into training table
            _, err = g.db.ExecContext(ctx, `
                INSERT INTO bid_training_data (
                    campaign_id,
                    floor_price,
                    engagement_score,
                    conversion_probability,
                    device_type,
                    segment_category,
                    country,
                    hour_of_day,
                    day_of_week,
                    historical_win_rate,
                    historical_avg_bid,
                    optimal_bid,
                    actual_outcome,
                    training_set
                ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
            `,
                trainingExample.CampaignID,
                trainingExample.FloorPrice,
                trainingExample.EngagementScore,
                trainingExample.ConversionProbability,
                trainingExample.DeviceType,
                trainingExample.SegmentCategory,
                trainingExample.Country,
                trainingExample.HourOfDay,
                trainingExample.DayOfWeek,
                trainingExample.HistoricalWinRate,
                trainingExample.HistoricalAvgBid,
                trainingExample.OptimalBid,
                trainingExample.ActualOutcome,
                g.assignTrainingSet(), // Random split
            )
        }
    }
    
    return nil
}

func (g *TrainingDataGenerator) createTrainingExample(
    ctx context.Context,
    bid *models.BidEvent,
    campaignID uuid.UUID,
) *TrainingExample {
    // Compute optimal bid (target variable)
    optimalBid := g.computeOptimalBid(bid)
    
    // Get historical features
    historical := g.getHistoricalFeatures(ctx, campaignID, bid.Timestamp)
    
    return &TrainingExample{
        CampaignID:            campaignID,
        FloorPrice:            bid.FloorPrice,
        EngagementScore:       bid.EngagementScore,
        ConversionProbability: bid.ConversionProbability,
        DeviceType:            bid.DeviceType,
        SegmentCategory:       bid.SegmentCategory,
        Country:               bid.Country,
        HourOfDay:             bid.Timestamp.Hour(),
        DayOfWeek:             int(bid.Timestamp.Weekday()),
        HistoricalWinRate:     historical.WinRate,
        HistoricalAvgBid:      historical.AvgBid,
        OptimalBid:            optimalBid,
        ActualOutcome:         g.determineOutcome(bid),
    }
}

func (g *TrainingDataGenerator) computeOptimalBid(bid *models.BidEvent) float64 {
    // Strategy 1: If won, use actual win price
    if bid.Won && bid.WinPrice != nil {
        return *bid.WinPrice
    }
    
    // Strategy 2: If lost, estimate what would have won
    // (could use percentile of winning bids in similar situations)
    // For now, use bid price + 10%
    return bid.BidPrice * 1.1
}

func (g *TrainingDataGenerator) assignTrainingSet() string {
    r := rand.Float64()
    switch {
    case r < 0.7:
        return "train"
    case r < 0.85:
        return "validation"
    default:
        return "test"
    }
}

type TrainingExample struct {
    CampaignID            uuid.UUID
    FloorPrice            float64
    EngagementScore       *float64
    ConversionProbability *float64
    DeviceType            string
    SegmentCategory       string
    Country               string
    HourOfDay             int
    DayOfWeek             int
    HistoricalWinRate     float64
    HistoricalAvgBid      float64
    OptimalBid            float64
    ActualOutcome         string
}
```

### 3. Training Pipeline (Python)

```python
# training/pipeline.py
import psycopg2
import pandas as pd
from sqlalchemy import create_engine
import yaml

class TrainingPipeline:
    def __init__(self, config_path='config.yaml'):
        with open(config_path) as f:
            self.config = yaml.safe_load(f)
        
        self.engine = create_engine(self.config['database_url'])
    
    def load_training_data(self):
        """Load data from PostgreSQL"""
        query = """
        SELECT 
            floor_price,
            engagement_score,
            conversion_probability,
            device_type,
            segment_category,
            country,
            hour_of_day,
            day_of_week,
            historical_win_rate,
            historical_avg_bid,
            optimal_bid,
            training_set
        FROM bid_training_data
        WHERE created_at >= NOW() - INTERVAL '30 days'
        """
        
        df = pd.read_sql(query, self.engine)
        return df
    
    def run_training(self):
        """Full training pipeline"""
        print("Loading data...")
        df = self.load_training_data()
        
        print(f"Total samples: {len(df)}")
        print(f"Train: {len(df[df.training_set=='train'])}")
        print(f"Val: {len(df[df.training_set=='validation'])}")
        print(f"Test: {len(df[df.training_set=='test'])}")
        
        # Train model
        optimizer = BidOptimizer()
        metrics = optimizer.train(df)
        
        # Export model
        model_path = f"models/bid_optimizer_{datetime.now().strftime('%Y%m%d_%H%M%S')}.onnx"
        optimizer.export_to_onnx(model_path)
        
        # Save metrics
        self.save_metrics(metrics, model_path)
        
        return model_path
    
    def save_metrics(self, metrics, model_path):
        """Save training metrics to database"""
        query = """
        INSERT INTO ml_model_metrics (
            model_path,
            train_r2,
            val_r2,
            feature_importance,
            created_at
        ) VALUES (%s, %s, %s, %s, NOW())
        """
        
        with self.engine.connect() as conn:
            conn.execute(query, (
                model_path,
                metrics['train_r2'],
                metrics['val_r2'],
                json.dumps(metrics['feature_importance'])
            ))

# Run pipeline
if __name__ == '__main__':
    pipeline = TrainingPipeline()
    model_path = pipeline.run_training()
    print(f"Model saved to: {model_path}")
```

---

## Model Training & Fine-Tuning Strategy

### 1. Initial Training

```bash
# Step 1: Generate training data (Go)
go run cmd/training-data-generator/main.go --start-date=2024-01-01 --end-date=2024-12-01

# Step 2: Train initial model (Python)
python training/pipeline.py

# Step 3: Evaluate model
python training/evaluate.py --model=models/bid_optimizer_20241201.onnx

# Step 4: Deploy to production (Go)
cp models/bid_optimizer_20241201.onnx /app/models/current.onnx
systemctl restart bidding-service
```

### 2. Continuous Learning (Daily Retraining)

```python
# training/continuous_learning.py
import schedule
import time

def retrain_models():
    """Daily retraining job"""
    pipeline = TrainingPipeline()
    
    # Load last 30 days of data
    df = pipeline.load_training_data()
    
    # Train new model
    new_model_path = pipeline.run_training()
    
    # A/B test new model
    deploy_for_testing(new_model_path, traffic_percentage=10)

# Schedule daily at 2 AM
schedule.every().day.at("02:00").do(retrain_models)

while True:
    schedule.run_pending()
    time.sleep(60)
```

### 3. Online Learning (Advanced)

```go
// internal/ml/online_learner.go
package ml

import (
    "context"
    "time"
)

type OnlineLearner struct {
    predictor    *ONNXPredictor
    feedbackChan chan *BidFeedback
    batchSize    int
    updateFreq   time.Duration
}

type BidFeedback struct {
    PredictionID string
    ActualBid    float64
    Won          bool
    WinPrice     *float64
    Converted    bool
}

func (l *OnlineLearner) CollectFeedback(feedback *BidFeedback) {
    l.feedbackChan <- feedback
}

func (l *OnlineLearner) StartLearning(ctx context.Context) {
    batch := make([]*BidFeedback, 0, l.batchSize)
    ticker := time.NewTicker(l.updateFreq)
    
    for {
        select {
        case <-ctx.Done():
            return
        case feedback := <-l.feedbackChan:
            batch = append(batch, feedback)
            
            if len(batch) >= l.batchSize {
                go l.updateModel(batch)
                batch = make([]*BidFeedback, 0, l.batchSize)
            }
        case <-ticker.C:
            if len(batch) > 0 {
                go l.updateModel(batch)
                batch = make([]*BidFeedback, 0, l.batchSize)
            }
        }
    }
}

func (l *OnlineLearner) updateModel(batch []*BidFeedback) {
    // Save feedback to database for next training cycle
    // Trigger model retrain if error threshold exceeded
}
```

### 4. Model Versioning

```go
// internal/ml/model_manager.go
package ml

type ModelManager struct {
    models map[string]*ONNXPredictor
    current string
    db *store.MLModelStore
}

func (m *ModelManager) LoadModel(version string) error {
    modelPath := fmt.Sprintf("models/bid_optimizer_%s.onnx", version)
    
    predictor, err := NewONNXPredictor(modelPath)
    if err != nil {
        return err
    }
    
    m.models[version] = predictor
    m.current = version
    
    // Log to database
    m.db.LogModelDeployment(version, modelPath)
    
    return nil
}

func (m *ModelManager) GetPredictor() *ONNXPredictor {
    return m.models[m.current]
}

func (m *ModelManager) Rollback(previousVersion string) error {
    if _, exists := m.models[previousVersion]; !exists {
        return fmt.Errorf("version not loaded: %s", previousVersion)
    }
    
    m.current = previousVersion
    return nil
}
```

---

## Performance Optimization

### 1. Batch Predictions

```go
// For high-throughput scenarios
func (p *ONNXPredictor) PredictBatch(
    requests []*models.BidRequest,
) ([]*models.BidResponse, error) {
    batchSize := len(requests)
    features := make([]float32, batchSize * p.featureCount)
    
    // Extract features for all requests
    for i, req := range requests {
        reqFeatures := p.extractFeatures(req, nil)
        copy(features[i*p.featureCount:], reqFeatures)
    }
    
    // Single inference for entire batch
    inputTensor := p.session.GetInputTensor(0)
    copy(inputTensor.GetData(), features)
    
    err := p.session.Run()
    if err != nil {
        return nil, err
    }
    
    // Parse results
    outputTensor := p.session.GetOutputTensor(0)
    predictions := outputTensor.GetData()
    
    responses := make([]*models.BidResponse, batchSize)
    for i := 0; i < batchSize; i++ {
        responses[i] = &models.BidResponse{
            BidPrice:   float64(predictions[i]),
            Confidence: 0.85,
            Strategy:   "ml_onnx_batch",
        }
    }
    
    return responses, nil
}
```

### 2. Feature Caching

```go
type FeatureCache struct {
    cache *bigcache.BigCache
    ttl   time.Duration
}

func (c *FeatureCache) GetHistoricalFeatures(
    campaignID uuid.UUID,
) (*HistoricalFeatures, bool) {
    key := fmt.Sprintf("hist:%s", campaignID)
    data, err := c.cache.Get(key)
    if err != nil {
        return nil, false
    }
    
    var features HistoricalFeatures
    json.Unmarshal(data, &features)
    return &features, true
}
```

---

## Monitoring & Evaluation

### 1. Model Performance Metrics

```go
// internal/ml/metrics.go
type ModelMetrics struct {
    RMSE           float64
    MAE            float64
    R2Score        float64
    PredictionTime time.Duration
    ThroughputRPS  float64
}

func (p *ONNXPredictor) EvaluateOnTestSet(
    testData []*TestExample,
) *ModelMetrics {
    var predictions []float64
    var actuals []float64
    
    start := time.Now()
    
    for _, example := range testData {
        pred, _ := p.PredictBidPrice(context.Background(), example.Request, nil)
        predictions = append(predictions, pred.BidPrice)
        actuals = append(actuals, example.ActualOptimalBid)
    }
    
    elapsed := time.Since(start)
    
    return &ModelMetrics{
        RMSE:           calculateRMSE(predictions, actuals),
        MAE:            calculateMAE(predictions, actuals),
        R2Score:        calculateR2(predictions, actuals),
        PredictionTime: elapsed / time.Duration(len(testData)),
        ThroughputRPS:  float64(len(testData)) / elapsed.Seconds(),
    }
}
```

### 2. Business Metrics Dashboard

```sql
-- Track model impact
CREATE TABLE model_performance_log (
    id UUID PRIMARY KEY,
    model_version VARCHAR(50),
    date DATE,
    total_bids INTEGER,
    win_rate DECIMAL(5,4),
    avg_bid_price DECIMAL(10,4),
    avg_win_price DECIMAL(10,4),
    total_spend DECIMAL(12,2),
    total_conversions INTEGER,
    roi DECIMAL(8,4),
    created_at TIMESTAMP DEFAULT NOW()
);

-- Daily aggregation query
INSERT INTO model_performance_log
SELECT
    gen_random_uuid(),
    'v2.0.0',
    CURRENT_DATE,
    COUNT(*),
    AVG(CASE WHEN won THEN 1 ELSE 0 END),
    AVG(bid_price),
    AVG(win_price),
    SUM(CASE WHEN won THEN win_price ELSE 0 END),
    SUM(CASE WHEN converted THEN 1 ELSE 0 END),
    SUM(CASE WHEN converted THEN 1 ELSE 0 END)::DECIMAL / 
        NULLIF(SUM(CASE WHEN won THEN win_price ELSE 0 END), 0)
FROM bid_events
WHERE timestamp >= CURRENT_DATE
  AND timestamp < CURRENT_DATE + INTERVAL '1 day';
```

---

## Implementation Roadmap

### Phase 1: Foundation (Weeks 1-2)
- [ ] Set up Python training environment
- [ ] Implement training data generation in Go
- [ ] Create initial XGBoost model
- [ ] Export to ONNX
- [ ] Integrate ONNX predictor into Go service

### Phase 2: Production (Weeks 3-4)
- [ ] Deploy model to staging
- [ ] A/B test against existing OpenAI predictor
- [ ] Monitor metrics
- [ ] Optimize inference performance
- [ ] Deploy to production

### Phase 3: Automation (Weeks 5-6)
- [ ] Automated daily retraining pipeline
- [ ] Model versioning system
- [ ] Rollback mechanisms
- [ ] Alert system for model degradation

### Phase 4: Advanced (Weeks 7-8)
- [ ] Feature store implementation
- [ ] Online learning experiments
- [ ] Multi-armed bandit for model selection
- [ ] Advanced feature engineering

---

## Cost-Benefit Analysis

### Current System (OpenAI API)
- **Cost**: $0.01 per prediction (GPT-4)
- **Latency**: 500-1000ms
- **Daily predictions**: 100,000
- **Monthly cost**: ~$30,000

### Proposed System (ONNX + Go)
- **Initial setup**: 2-3 weeks dev time
- **Cost**: $0.0001 per prediction (compute only)
- **Latency**: 5-10ms
- **Daily predictions**: 100,000
- **Monthly cost**: ~$300

**Savings**: ~$29,700/month (~$356,400/year)
**ROI**: 2-3 weeks payback period

---

## Conclusion

**Recommended Approach:**
1. Start with XGBoost/LightGBM trained in Python
2. Export to ONNX for Go inference
3. Build automated retraining pipeline
4. Gradually transition from OpenAI to custom models
5. Keep OpenAI as fallback for edge cases

This hybrid approach gives you:
- ✅ Fast inference in Go (5-10ms vs 500ms+)
- ✅ Low cost (<1% of current OpenAI costs)
- ✅ Full control over model and data
- ✅ Easy experimentation in Python
- ✅ Production reliability in Go
- ✅ Proven ML frameworks (XGBoost/ONNX)

Let me know if you'd like me to elaborate on any section or provide more code examples!
