# 🪟 Windows Quick Start Guide

**Goal:** Get your first ML model trained and running in 20 minutes!

## 🚀 Phase 1: Setup (5 minutes)

### Step 1: Open PowerShell
- Press `Win + X`
- Select "Windows PowerShell" or "Terminal"
- Navigate to your project directory:
```powershell
cd C:\path\to\your\bidding-analysis
```

### Step 2: Run Setup Script
```powershell
# If you get an execution policy error:
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser

# Run the setup
.\setup.ps1
```

**What this does:**
- Creates Python virtual environment
- Installs all ML packages (XGBoost, ONNX, etc.)
- Creates necessary directories
- Takes 3-5 minutes

### Step 3: Activate Environment
```powershell
.\venv\Scripts\Activate.ps1

# You should see (venv) in your prompt
```

## 🧪 Phase 2: Test Everything Works (5 minutes)

### Run the Test Script
```powershell
python test_training.py
```

**Expected output:**
```
[1/6] Testing basic imports...
✅ numpy: 1.26.3
✅ pandas: 2.2.0

[2/6] Testing ML packages...
✅ xgboost: 2.0.3
✅ scikit-learn: imported successfully

[3/6] Creating synthetic training data...
✅ Created 1000 training samples

[4/6] Training XGBoost model...
✅ Model training complete
   Test R²: 0.9234

[5/6] Testing ONNX export...
✅ Model exported to: models/test_model.onnx

[6/6] Testing ONNX inference...
✅ ONNX model loaded
✅ Predictions match!

🎉 SUCCESS! Your ML environment is fully functional!
```

**If you see ✅ on all 6 steps - you're ready to go!**

## 🎯 Phase 3: Train Your First Real Model (10 minutes)

### Option A: Test with Synthetic Data (No Database Needed)

If you don't have real data yet, train on synthetic data:

```powershell
# This uses the training script with mock data
python train_model.py --config config.yaml --days 30 --output models\my_first_model.onnx
```

### Option B: Use Real Data (Database Required)

**Step 1:** Configure database connection

Edit `config.yaml`:
```yaml
database:
  host: localhost
  port: 5432
  name: bidding_analysis
  user: postgres
  password: YOUR_PASSWORD  # Change this!
```

**Step 2:** Generate training data

```powershell
# Run your Go data generator
go run .\cmd\training-data-generator\main.go --days=30

# Or manually insert test data in PostgreSQL (see WINDOWS_SETUP_GUIDE.md)
```

**Step 3:** Train model

```powershell
.\train.ps1
```

## 📊 Phase 4: Test Your Model (2 minutes)

```powershell
# Test the model you just trained
.\test-model.ps1 -Model "models\bid_optimizer_latest.onnx"
```

**Expected output:**
```
[1/3] Testing ONNX model loading...
✅ Model loaded successfully

[2/3] Testing predictions...
✅ Prediction successful
   Predicted optimal bid: $3.2450

[3/3] Running performance benchmark...
✅ Performance benchmark complete
   Average inference time: 2.34 ms
   Throughput: 427 predictions/second
```

## ✅ Success Checklist

You're ready when:
- [x] Setup script ran without errors
- [x] Test script shows all ✅
- [x] Model training completed
- [x] Model test passed
- [x] You have a `.onnx` file in `models/` directory

## 🎓 Next Steps

### 1. Integrate with Your Go Service

Copy these files to your Go project:
```powershell
# Copy the ONNX predictor
Copy-Item onnx_predictor.go -Destination internal\mlonnx\

# Copy your trained model
Copy-Item models\bid_optimizer_latest.onnx -Destination models\
Copy-Item models\bid_optimizer_latest_encoders.json -Destination models\
```

### 2. Update Your Go Code

In `internal/ml/predictor.go`:
```go
import "github.com/baskint/bidding-analysis/internal/mlonnx"

func NewPredictor(config *Config, bidStore *store.BidStore) *Predictor {
    var client AIClient
    
    // Try ONNX model first
    onnxPredictor, err := mlonnx.NewONNXPredictor(
        "models/bid_optimizer_latest.onnx",
        "models/bid_optimizer_latest_encoders.json",
    )
    if err == nil {
        client = onnxPredictor
        log.Println("✅ Using ONNX ML model")
    } else {
        // Fallback to OpenAI
        client = NewOpenAIClient(config.OpenAIKey)
        log.Printf("⚠️  Using OpenAI fallback: %v", err)
    }
    
    return &Predictor{
        openaiClient: client,
        bidStore:     bidStore,
    }
}
```

### 3. Test in Go

```powershell
# Test your Go service
go test .\internal\ml\... -v

# Or run the service
go run .\cmd\server\main.go
```

### 4. Deploy Gradually

Start with A/B testing:
- Route 10% of traffic to new model
- Monitor metrics for 24 hours
- If good, increase to 50%
- Then 100%

## 🔧 Troubleshooting

### "Execution Policy" Error
```powershell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

### "Python not found"
Install Python 3.10+ from: https://www.python.org/downloads/

### "pip install failed"
```powershell
# Activate venv first
.\venv\Scripts\Activate.ps1

# Then try again
pip install xgboost
```

### "No module named 'psycopg2'"
```powershell
pip install psycopg2-binary
```

### "ONNX Runtime error"
```powershell
pip install onnxruntime
```

### "Go build error with ONNX"
Download ONNX Runtime DLL:
```powershell
# Run this in PowerShell
$url = "https://github.com/microsoft/onnxruntime/releases/download/v1.16.3/onnxruntime-win-x64-1.16.3.zip"
Invoke-WebRequest -Uri $url -OutFile onnxruntime.zip
Expand-Archive onnxruntime.zip
Copy-Item "onnxruntime\onnxruntime-win-x64-1.16.3\lib\onnxruntime.dll" -Destination .
```

## 💡 Pro Tips

### Daily Workflow
```powershell
# 1. Activate environment
.\venv\Scripts\Activate.ps1

# 2. Generate new training data
go run .\cmd\training-data-generator\main.go --days=7

# 3. Train new model
.\train.ps1

# 4. Test it
.\test-model.ps1

# 5. Deploy if good
Copy-Item models\bid_optimizer_latest.onnx -Destination production\
```

### Check Model Performance
```powershell
# See all trained models
Get-ChildItem models\*.onnx | Select-Object Name, Length, LastWriteTime

# Check logs
Get-Content logs\training_*.log | Select-Object -Last 20
```

### Database Quick Check
```powershell
# Connect to database
psql -U postgres -d bidding_analysis

# Check training data
SELECT COUNT(*) FROM bid_training_data;
SELECT training_set, COUNT(*) FROM bid_training_data GROUP BY training_set;
```

## 📚 What You've Learned

✅ Set up Python ML environment
✅ Train XGBoost models
✅ Export to ONNX format
✅ Test model performance
✅ Integrate with Go service

## 🎯 Performance Targets

Your model should achieve:
- **Inference time:** < 10ms per prediction
- **Throughput:** > 100 predictions/second
- **Accuracy (R²):** > 0.80 on validation set
- **Cost:** $300/month (vs $30K with OpenAI)

## 📖 Reference

**Full documentation:**
- [WINDOWS_SETUP_GUIDE.md](WINDOWS_SETUP_GUIDE.md) - Detailed setup
- [ml-recommendations.md](ml-recommendations.md) - Architecture details
- [ML_IMPLEMENTATION_GUIDE.md](ML_IMPLEMENTATION_GUIDE.md) - Complete guide

**Scripts:**
- `setup.ps1` - Initial setup
- `train.ps1` - Train model
- `test-model.ps1` - Test model
- `test_training.py` - Verify setup

**Need help?** Check the troubleshooting section or the full guides!

---

**🎉 Congratulations!** You've successfully set up ML training on Windows!
