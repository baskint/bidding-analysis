# ML Implementation Guide for Bid Optimization

This guide helps you implement a production-ready ML pipeline for your bid optimization platform.

## 🎯 Quick Start (30 minutes)

```bash
# 1. Set up Python environment
make -f Makefile.ml setup

# 2. Activate environment
source venv/bin/activate

# 3. Configure database connection
cp config.yaml.example config.yaml
# Edit config.yaml with your database credentials

# 4. Generate training data
make -f Makefile.ml generate-data

# 5. Train your first model
make -f Makefile.ml train

# 6. Test Go inference
make -f Makefile.ml test-go
```

## 📋 Prerequisites

### System Requirements
- Go 1.21+
- Python 3.10+
- PostgreSQL 14+
- 8GB RAM minimum (16GB recommended)
- 10GB disk space for models and data

### Go Dependencies
Add to your `go.mod`:
```go
require (
    github.com/yalue/onnxruntime_go v1.5.0
)
```

### System Libraries
For ONNX Runtime:
```bash
# Ubuntu/Debian
sudo apt-get install -y libgomp1

# macOS
brew install libomp

# Windows
# Download ONNX Runtime from https://github.com/microsoft/onnxruntime/releases
```

## 🏗️ Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                    Data Collection Layer                     │
│  • PostgreSQL (bid_events, campaigns)                        │
│  • Training data generator (Go)                              │
└───────────────────────┬─────────────────────────────────────┘
                        │
                        ↓
┌─────────────────────────────────────────────────────────────┐
│                   Training Pipeline (Python)                 │
│  • Feature engineering                                       │
│  • XGBoost/LightGBM training                                 │
│  • Model validation                                          │
│  • ONNX export                                               │
└───────────────────────┬─────────────────────────────────────┘
                        │
                        ↓
┌─────────────────────────────────────────────────────────────┐
│                 Production Inference (Go)                    │
│  • ONNX Runtime                                              │
│  • Fast predictions (<10ms)                                  │
│  • Batch processing                                          │
│  • Model versioning                                          │
└─────────────────────────────────────────────────────────────┘
```

## 📁 Project Structure

```
bidding-analysis/
├── cmd/
│   ├── server/                 # Main API server
│   ├── training-data-generator/# Data generation tool
│   └── benchmark/              # Inference benchmarks
├── internal/
│   ├── ml/                     # Current ML code
│   └── mlonnx/                 # New ONNX predictor
├── models/                     # Trained ONNX models
│   ├── bid_optimizer_latest.onnx
│   └── bid_optimizer_latest_encoders.json
├── training/                   # Python training code
│   ├── train_model.py         # Main training script
│   ├── config.yaml            # Configuration
│   ├── requirements.txt       # Python dependencies
│   └── scripts/               # Helper scripts
│       ├── evaluate_model.py
│       ├── monitor_model.py
│       └── continuous_training.py
└── migrations/
    └── 0004_training_data.sql # Training data schema
```

## 🔧 Step-by-Step Implementation

### Step 1: Database Setup

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
    hour_of_day INTEGER,
    day_of_week INTEGER,
    
    -- Historical features
    historical_win_rate DECIMAL(5,4),
    historical_avg_bid DECIMAL(10,4),
    historical_avg_win_price DECIMAL(10,4),
    
    -- Target
    optimal_bid DECIMAL(10,4) NOT NULL,
    actual_outcome VARCHAR(20),
    
    -- Metadata
    training_set VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    FOREIGN KEY (campaign_id) REFERENCES campaigns(id)
);

CREATE INDEX idx_training_data_campaign 
ON bid_training_data(campaign_id);

-- Metrics tracking
CREATE TABLE ml_model_metadata (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    model_path VARCHAR(500) NOT NULL,
    model_type VARCHAR(100),
    train_rmse DECIMAL(10,6),
    val_rmse DECIMAL(10,6),
    train_r2 DECIMAL(10,6),
    val_r2 DECIMAL(10,6),
    feature_importance JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

### Step 2: Generate Training Data

Create `cmd/training-data-generator/main.go`:

```go
package main

import (
    "context"
    "flag"
    "log"
    "time"
    
    "github.com/baskint/bidding-analysis/internal/store"
)

func main() {
    days := flag.Int("days", 30, "Days of historical data")
    flag.Parse()
    
    // Initialize stores
    db := store.NewPostgresStore(os.Getenv("DATABASE_URL"))
    defer db.Close()
    
    generator := NewTrainingDataGenerator(db)
    
    startDate := time.Now().AddDate(0, 0, -*days)
    endDate := time.Now()
    
    log.Printf("Generating training data from %v to %v\n", startDate, endDate)
    
    err := generator.Generate(context.Background(), startDate, endDate)
    if err != nil {
        log.Fatalf("Failed to generate data: %v", err)
    }
    
    log.Println("✅ Training data generation complete")
}
```

Run:
```bash
go run cmd/training-data-generator/main.go --days=30
```

### Step 3: Train Model

```bash
# Install dependencies
source venv/bin/activate
pip install -r requirements.txt

# Configure
cp config.yaml.example config.yaml
# Edit database credentials in config.yaml

# Train
python training/train_model.py \
    --config config.yaml \
    --days 30 \
    --output models/bid_optimizer.onnx
```

Output:
```
Loading training data...
Loaded 50000 training examples
Engineering features...
Dataset size after cleaning: 49823
Train set: 34876 samples
Validation set: 14947 samples
Training XGBoost model...
...
=== Training Results ===
Train RMSE: 0.1234
Val RMSE: 0.1456
Train R²: 0.8765
Val R²: 0.8543

✅ Training complete! Model saved to: models/bid_optimizer_20241201_120000.onnx
```

### Step 4: Integrate ONNX Predictor in Go

Update `internal/ml/predictor.go`:

```go
package ml

import (
    "github.com/baskint/bidding-analysis/internal/mlonnx"
)

func NewPredictor(config *Config, bidStore *store.BidStore) *Predictor {
    var client AIClient
    
    // Try to load ONNX model
    if config.ModelPath != "" {
        onnxPredictor, err := mlonnx.NewONNXPredictor(
            config.ModelPath,
            config.EncodersPath,
        )
        if err == nil {
            client = onnxPredictor
            log.Printf("✅ Using ONNX model: %s", config.ModelPath)
        } else {
            log.Printf("⚠️  Failed to load ONNX: %v, using fallback", err)
        }
    }
    
    // Fallback to OpenAI or mock
    if client == nil {
        if config.OpenAIKey != "" {
            client = NewOpenAIClient(config.OpenAIKey)
        } else {
            client = NewMockOpenAIClient()
        }
    }
    
    return &Predictor{
        openaiClient: client,
        bidStore:     bidStore,
    }
}
```

Update `.env`:
```bash
# ML Model Configuration
MODEL_PATH=models/bid_optimizer_latest.onnx
ENCODERS_PATH=models/bid_optimizer_latest_encoders.json
```

### Step 5: Test & Deploy

```bash
# Test Go inference
make -f Makefile.ml test-go

# Benchmark performance
go run cmd/benchmark/main.go \
    --model models/bid_optimizer_latest.onnx \
    --iterations 10000

# Expected output:
# Average inference time: 5.2ms
# Throughput: 192 predictions/second
# Memory usage: 45MB

# Deploy to production
make -f Makefile.ml deploy
```

## 📊 Model Performance Tracking

### Monitor Metrics

```bash
# View recent model performance
make -f Makefile.ml monitor

# Start live dashboard
make -f Makefile.ml monitor-live
```

### Key Metrics to Track

1. **Prediction Accuracy**
   - RMSE (Root Mean Squared Error)
   - MAE (Mean Absolute Error)
   - R² Score

2. **Business Metrics**
   - Win rate improvement
   - Cost per conversion
   - ROI vs previous model

3. **System Metrics**
   - Inference latency (target: <10ms)
   - Throughput (predictions/sec)
   - Memory usage

## 🔄 Continuous Learning

### Daily Retraining

Set up cron job:
```bash
# Edit crontab
crontab -e

# Add daily retraining at 2 AM
0 2 * * * cd /path/to/bidding-analysis && make -f Makefile.ml retrain
```

Or use the Python scheduler:
```python
# training/scripts/continuous_training.py
import schedule
import time

def retrain_pipeline():
    # Generate new data
    os.system("make -f Makefile.ml generate-data")
    
    # Train new model
    os.system("make -f Makefile.ml train")
    
    # Evaluate
    metrics = evaluate_model()
    
    # Deploy if better
    if metrics['val_r2'] > current_model_r2:
        os.system("make -f Makefile.ml deploy-gradual")

schedule.every().day.at("02:00").do(retrain_pipeline)

while True:
    schedule.run_pending()
    time.sleep(60)
```

## 🔬 A/B Testing

### Gradual Rollout

```bash
# Deploy to 10% of traffic
make -f Makefile.ml deploy-gradual
```

This will:
1. Deploy new model alongside existing one
2. Route 10% of traffic to new model
3. Monitor metrics for 24 hours
4. Automatically increase to 30% if metrics improve
5. Continue until 100% or rollback if metrics degrade

### Manual A/B Test

Update `internal/ml/model_manager.go`:

```go
type ModelManager struct {
    models map[string]AIClient
    trafficSplit map[string]float64
}

func (m *ModelManager) GetPredictor() AIClient {
    // Route traffic based on split
    rand := rand.Float64()
    
    var cumulative float64
    for modelVersion, percentage := range m.trafficSplit {
        cumulative += percentage
        if rand < cumulative {
            return m.models[modelVersion]
        }
    }
    
    return m.models["default"]
}
```

## 🐛 Troubleshooting

### Model Not Loading

```bash
# Check ONNX Runtime installation
go test github.com/yalue/onnxruntime_go

# Verify model file
ls -lh models/bid_optimizer_latest.onnx

# Test Python → ONNX export
python -c "
import onnxruntime as ort
sess = ort.InferenceSession('models/bid_optimizer_latest.onnx')
print('✅ Model loads correctly')
"
```

### Poor Predictions

1. **Check training data quality**
   ```bash
   make -f Makefile.ml validate-data
   ```

2. **Verify feature encoding consistency**
   - Ensure Go feature extraction matches Python training

3. **Compare predictions**
   ```python
   # Compare Python vs Go predictions
   python scripts/compare_predictions.py
   ```

### Memory Issues

```go
// Use batch prediction for high throughput
responses, err := predictor.PredictBatch(ctx, requests, historicalMap)

// Or implement model pooling
type ModelPool struct {
    models chan *ONNXPredictor
}
```

## 📈 Performance Optimization

### 1. Feature Caching

```go
type FeatureCache struct {
    cache *bigcache.BigCache
}

func (c *FeatureCache) GetOrCompute(
    key string,
    computeFn func() []float32,
) []float32 {
    if cached, err := c.cache.Get(key); err == nil {
        return deserialize(cached)
    }
    
    features := computeFn()
    c.cache.Set(key, serialize(features))
    return features
}
```

### 2. Batch Processing

Process multiple bids in a single inference call:
```go
// Instead of:
for _, req := range requests {
    pred, _ := predictor.PredictBidPrice(ctx, req, nil)
}

// Use:
preds, _ := predictor.PredictBatch(ctx, requests, historicalMap)
```

### 3. Model Warm-up

```go
func warmUpModel(predictor *ONNXPredictor) {
    dummyRequests := generateDummyRequests(100)
    predictor.PredictBatch(context.Background(), dummyRequests, nil)
}
```

## 🎓 Advanced Topics

### Online Learning

Implement feedback loop:
```go
type FeedbackCollector struct {
    buffer []*BidFeedback
}

func (c *FeedbackCollector) CollectOutcome(
    predictionID string,
    won bool,
    winPrice *float64,
    converted bool,
) {
    feedback := &BidFeedback{
        PredictionID: predictionID,
        Won: won,
        WinPrice: winPrice,
        Converted: converted,
        Timestamp: time.Now(),
    }
    
    c.buffer = append(c.buffer, feedback)
    
    // Trigger retraining when buffer full
    if len(c.buffer) >= 10000 {
        go c.triggerRetraining()
    }
}
```

### Multi-Model Ensemble

```go
type EnsemblePredictor struct {
    models []AIClient
    weights []float64
}

func (e *EnsemblePredictor) PredictBidPrice(
    ctx context.Context,
    req *models.BidRequest,
    historical []*models.BidEvent,
) (*models.BidResponse, error) {
    predictions := make([]*models.BidResponse, len(e.models))
    
    // Get predictions from all models
    for i, model := range e.models {
        pred, err := model.PredictBidPrice(ctx, req, historical)
        if err != nil {
            continue
        }
        predictions[i] = pred
    }
    
    // Weighted average
    return e.combine(predictions), nil
}
```

### Feature Store Integration

Use [Feast](https://feast.dev/) for feature management:
```go
import "github.com/feast-dev/feast/sdk/go"

type FeatureStore struct {
    client *feast.Client
}

func (s *FeatureStore) GetFeatures(
    entityID string,
) (map[string]float32, error) {
    features, err := s.client.GetOnlineFeatures(
        context.Background(),
        []string{"campaign_id"},
        []feast.Row{{EntityKey: entityID}},
        []string{
            "campaign_stats:win_rate_7d",
            "campaign_stats:avg_bid_7d",
            "user_segment:engagement_score",
        },
    )
    return features, err
}
```

## 📚 Additional Resources

- [XGBoost Documentation](https://xgboost.readthedocs.io/)
- [ONNX Runtime Go](https://github.com/yalue/onnxruntime_go)
- [MLflow for Experiment Tracking](https://mlflow.org/)
- [Feast Feature Store](https://feast.dev/)

## 🤝 Contributing

When adding new features:
1. Update `FEATURE_COLUMNS` in both Python and Go
2. Add tests for feature extraction
3. Document feature meaning
4. Update model version

## 📝 Checklist

- [ ] Database schema created
- [ ] Training data generator implemented
- [ ] Python environment set up
- [ ] First model trained
- [ ] ONNX export working
- [ ] Go integration complete
- [ ] Tests passing
- [ ] Monitoring dashboard set up
- [ ] Deployment pipeline configured
- [ ] Documentation updated

## 🎯 Next Steps

1. **Week 1**: Set up infrastructure and train first model
2. **Week 2**: Integrate with Go service and A/B test
3. **Week 3**: Implement automated retraining
4. **Week 4**: Optimize and scale

Need help? Check the main ML recommendations document or create an issue!
