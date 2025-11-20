# 🪟 Windows 11 Setup Guide - ML Bid Optimization

**Designed for:** Windows 11, Python 3.14, Visual Studio 2022

Let's get you up and running step by step!

## 📋 Prerequisites Check

### What You Already Have ✅
- Windows 11
- Python 3.14 ✅
- Visual Studio 2022 ✅

### What We Need to Install

Open PowerShell as Administrator and let's check what's missing:

```powershell
# Check Python
python --version
# Should show: Python 3.14.0

# Check Go (we need this)
go version
# If error: we'll install it

# Check PostgreSQL (we need this)
psql --version
# If error: we'll install it

# Check Git
git --version
# If error: we'll install it
```

## 🚀 Part 1: Install Missing Tools (15 minutes)

### Install Go (if not installed)

1. Download Go for Windows:
   - Visit: https://go.dev/dl/
   - Download: `go1.21.x.windows-amd64.msi` (or latest)
   - Run installer, use default settings

2. Verify installation:
```powershell
go version
# Should show: go version go1.21.x windows/amd64
```

### Install PostgreSQL (if not installed)

1. Download PostgreSQL for Windows:
   - Visit: https://www.postgresql.org/download/windows/
   - Download: PostgreSQL 15 or 16 installer
   - Run installer:
     - Password: Choose a strong password (save it!)
     - Port: 5432 (default)
     - Locale: Default

2. Verify installation:
```powershell
psql --version
# Should show: psql (PostgreSQL) 15.x
```

3. Create database:
```powershell
# Open SQL Shell (psql) from Start Menu
# Or use PowerShell:
psql -U postgres

# In psql:
CREATE DATABASE bidding_analysis;
\q
```

### Install Git (if not installed)

1. Download Git for Windows:
   - Visit: https://git-scm.com/download/win
   - Download and run installer
   - Use default settings

## 🐍 Part 2: Set Up Python Environment (10 minutes)

### Create Virtual Environment

Open PowerShell in your project directory:

```powershell
# Navigate to your project
cd C:\path\to\bidding-analysis

# Create virtual environment
python -m venv venv

# Activate it
.\venv\Scripts\Activate.ps1

# If you get an error about execution policy:
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
# Then try activating again

# You should see (venv) in your prompt
```

### Install Python Packages

```powershell
# Make sure venv is activated (you see (venv) in prompt)

# Upgrade pip first
python -m pip install --upgrade pip

# Install ML packages
pip install xgboost==2.0.3
pip install lightgbm==4.1.0
pip install scikit-learn==1.4.0
pip install numpy==1.26.3
pip install pandas==2.2.0

# Install ONNX
pip install onnxmltools==1.12.0
pip install onnxruntime==1.16.3
pip install skl2onnx==1.16.0

# Install database driver
pip install psycopg2-binary==2.9.9

# Install utilities
pip install pyyaml==6.0.1
pip install jupyter

# This should take 2-5 minutes
```

### Verify Installation

```powershell
python -c "import xgboost; print('XGBoost:', xgboost.__version__)"
python -c "import onnxruntime; print('ONNX Runtime:', onnxruntime.__version__)"
python -c "import pandas; print('Pandas:', pandas.__version__)"

# All should print versions without errors
```

## 🗄️ Part 3: Set Up Database Schema (5 minutes)

### Configure Database Connection

Create a file called `config.yaml` in your project root:

```yaml
database:
  host: localhost
  port: 5432
  name: bidding_analysis
  user: postgres
  password: YOUR_PASSWORD_HERE  # Change this!

training:
  days_of_history: 30
  test_size: 0.2
```

### Run Existing Migrations

If you have migrations set up:

```powershell
# Check if migrate tool is installed
migrate --version

# If not, install it:
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run migrations
migrate -database "postgres://postgres:YOUR_PASSWORD@localhost:5432/bidding_analysis?sslmode=disable" -path migrations up
```

Or manually create the training table:

```sql
-- Copy this SQL and run in pgAdmin or psql

CREATE TABLE IF NOT EXISTS bid_training_data (
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
    campaign_spend_last_7d DECIMAL(12,2),
    campaign_conversions_last_7d INTEGER,
    
    -- Target variable
    optimal_bid DECIMAL(10,4) NOT NULL,
    actual_outcome VARCHAR(20),
    
    -- Metadata
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    training_set VARCHAR(50) CHECK (training_set IN ('train', 'validation', 'test'))
);

CREATE INDEX idx_training_data_campaign ON bid_training_data(campaign_id);
CREATE INDEX idx_training_data_created ON bid_training_data(created_at DESC);
CREATE INDEX idx_training_data_training_set ON bid_training_data(training_set);

-- Model metadata table
CREATE TABLE IF NOT EXISTS ml_model_metadata (
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

## 🧪 Part 4: First Test - Train a Simple Model (15 minutes)

Let's test the training pipeline with a simple script first.

### Create Test Script

Create `test_training.py`:

```python
import pandas as pd
import numpy as np
from sklearn.model_selection import train_test_split
import xgboost as xgb
from sklearn.metrics import mean_squared_error, r2_score

print("🚀 Testing ML Training Pipeline")
print("-" * 50)

# 1. Create synthetic data (we'll use real data later)
print("\n📊 Creating synthetic training data...")
np.random.seed(42)
n_samples = 1000

data = {
    'floor_price': np.random.uniform(0.5, 5.0, n_samples),
    'engagement_score': np.random.uniform(0, 1, n_samples),
    'conversion_probability': np.random.uniform(0, 0.3, n_samples),
    'historical_win_rate': np.random.uniform(0.2, 0.8, n_samples),
    'device_type': np.random.choice(['mobile', 'desktop', 'tablet'], n_samples),
    'hour_of_day': np.random.randint(0, 24, n_samples),
}

# Create target: optimal_bid (simple formula for testing)
data['optimal_bid'] = (
    data['floor_price'] * 1.2 +
    data['engagement_score'] * 0.5 +
    data['conversion_probability'] * 2.0 +
    np.random.normal(0, 0.1, n_samples)
)

df = pd.DataFrame(data)
print(f"✅ Created {len(df)} samples")

# 2. Feature engineering
print("\n🔧 Engineering features...")
df['device_type_encoded'] = pd.Categorical(df['device_type']).codes
X = df[['floor_price', 'engagement_score', 'conversion_probability', 
        'historical_win_rate', 'device_type_encoded', 'hour_of_day']]
y = df['optimal_bid']
print(f"✅ Features: {list(X.columns)}")

# 3. Train/test split
print("\n✂️ Splitting data...")
X_train, X_test, y_train, y_test = train_test_split(
    X, y, test_size=0.2, random_state=42
)
print(f"✅ Train: {len(X_train)}, Test: {len(X_test)}")

# 4. Train XGBoost
print("\n🎓 Training XGBoost model...")
model = xgb.XGBRegressor(
    n_estimators=100,
    max_depth=5,
    learning_rate=0.1,
    random_state=42
)
model.fit(X_train, y_train)
print("✅ Model trained!")

# 5. Evaluate
print("\n📈 Evaluating model...")
train_pred = model.predict(X_train)
test_pred = model.predict(X_test)

train_rmse = np.sqrt(mean_squared_error(y_train, train_pred))
test_rmse = np.sqrt(mean_squared_error(y_test, test_pred))
train_r2 = r2_score(y_train, train_pred)
test_r2 = r2_score(y_test, test_pred)

print(f"✅ Train RMSE: {train_rmse:.4f}")
print(f"✅ Test RMSE: {test_rmse:.4f}")
print(f"✅ Train R²: {train_r2:.4f}")
print(f"✅ Test R²: {test_r2:.4f}")

# 6. Test ONNX export
print("\n📦 Testing ONNX export...")
try:
    import onnxmltools
    from onnxmltools.convert.xgboost.operator_converters.XGBoost import convert_xgboost
    from onnxconverter_common import FloatTensorType
    
    initial_type = [('float_input', FloatTensorType([None, len(X.columns)]))]
    onnx_model = convert_xgboost(model, initial_types=initial_type, target_opset=12)
    
    with open('test_model.onnx', 'wb') as f:
        f.write(onnx_model.SerializeToString())
    
    print("✅ Model exported to test_model.onnx")
    
    # Test loading
    import onnxruntime as ort
    sess = ort.InferenceSession('test_model.onnx')
    print("✅ ONNX model loads successfully!")
    
    # Test prediction
    import numpy as np
    test_input = X_test.iloc[0:1].values.astype(np.float32)
    result = sess.run(None, {'float_input': test_input})
    print(f"✅ ONNX prediction: {result[0][0][0]:.4f}")
    print(f"✅ XGBoost prediction: {test_pred[0]:.4f}")
    
except Exception as e:
    print(f"❌ ONNX export failed: {e}")

print("\n" + "="*50)
print("🎉 All tests passed! Your ML environment is ready!")
print("="*50)
```

### Run Test

```powershell
# Make sure venv is activated
python test_training.py
```

**Expected output:**
```
🚀 Testing ML Training Pipeline
--------------------------------------------------
📊 Creating synthetic training data...
✅ Created 1000 samples
🔧 Engineering features...
✅ Features: ['floor_price', 'engagement_score', ...]
✂️ Splitting data...
✅ Train: 800, Test: 200
🎓 Training XGBoost model...
✅ Model trained!
📈 Evaluating model...
✅ Train RMSE: 0.1234
✅ Test RMSE: 0.1456
✅ Train R²: 0.8765
✅ Test R²: 0.8543
📦 Testing ONNX export...
✅ Model exported to test_model.onnx
✅ ONNX model loads successfully!
✅ ONNX prediction: 2.3456
✅ XGBoost prediction: 2.3457
==================================================
🎉 All tests passed! Your ML environment is ready!
==================================================
```

## 🎯 Part 5: Install Go ONNX Runtime (10 minutes)

Now let's set up Go to use ONNX models.

### Install ONNX Runtime for Go

```powershell
# In your project directory
go get github.com/yalue/onnxruntime_go@latest
```

### Download ONNX Runtime DLL

Windows needs the ONNX Runtime DLL:

1. Visit: https://github.com/microsoft/onnxruntime/releases
2. Download: `onnxruntime-win-x64-{version}.zip` (latest version)
3. Extract the ZIP file
4. Copy `onnxruntime.dll` to your project directory or system PATH

**Or use PowerShell:**

```powershell
# Download ONNX Runtime (version 1.16.3)
$url = "https://github.com/microsoft/onnxruntime/releases/download/v1.16.3/onnxruntime-win-x64-1.16.3.zip"
$output = "onnxruntime.zip"

Invoke-WebRequest -Uri $url -OutFile $output
Expand-Archive -Path $output -DestinationPath "onnxruntime"

# Copy DLL to project directory
Copy-Item "onnxruntime\onnxruntime-win-x64-1.16.3\lib\onnxruntime.dll" -Destination "."

# Clean up
Remove-Item $output
Remove-Item -Recurse "onnxruntime"

Write-Host "✅ ONNX Runtime DLL installed"
```

### Test Go ONNX Runtime

Create `test_onnx.go`:

```go
package main

import (
	"fmt"
	"log"

	"github.com/yalue/onnxruntime_go"
)

func main() {
	fmt.Println("🧪 Testing ONNX Runtime in Go")
	fmt.Println(strings.Repeat("-", 50))

	// Initialize ONNX Runtime
	err := onnxruntime_go.InitializeEnvironment()
	if err != nil {
		log.Fatalf("❌ Failed to initialize ONNX: %v", err)
	}
	defer onnxruntime_go.DestroyEnvironment()
	fmt.Println("✅ ONNX Runtime initialized")

	// Load model (make sure test_model.onnx exists)
	modelPath := "test_model.onnx"
	
	// Create dummy input (6 features from our test)
	inputData := []float32{1.5, 0.7, 0.15, 0.45, 1.0, 12.0}
	
	session, err := onnxruntime_go.NewAdvancedSession(
		modelPath,
		[]string{"float_input"},
		[]string{"output"},
		inputData,
		[]int64{1, 6},
		nil,
	)
	if err != nil {
		log.Fatalf("❌ Failed to load model: %v", err)
	}
	defer session.Destroy()
	fmt.Println("✅ Model loaded successfully")

	// Run inference
	err = session.Run()
	if err != nil {
		log.Fatalf("❌ Inference failed: %v", err)
	}
	fmt.Println("✅ Inference completed")

	// Get result
	outputTensor := session.GetOutputTensor(0)
	result := outputTensor.GetData()[0]
	fmt.Printf("✅ Prediction: %.4f\n", result)

	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("🎉 Go ONNX Runtime is working!")
}
```

Run it:

```powershell
go run test_onnx.go
```

## 📚 Part 6: PowerShell Scripts Instead of Makefile

Since Windows doesn't have `make`, here are PowerShell equivalents.

### Create `train.ps1`

```powershell
# train.ps1 - Train ML model
Write-Host "🎓 Training ML Model" -ForegroundColor Cyan
Write-Host ("-" * 50)

# Activate venv
.\venv\Scripts\Activate.ps1

# Create models directory
New-Item -ItemType Directory -Force -Path "models" | Out-Null

# Run training
python train_model.py --config config.yaml --days 30 --output "models\bid_optimizer.onnx"

Write-Host "✅ Training complete!" -ForegroundColor Green
```

### Create `test-go.ps1`

```powershell
# test-go.ps1 - Test Go inference
Write-Host "🧪 Testing Go Inference" -ForegroundColor Cyan
Write-Host ("-" * 50)

# Run Go tests
go test -v .\internal\mlonnx\...

# Run benchmark if exists
if (Test-Path "cmd\benchmark\main.go") {
    go run cmd\benchmark\main.go --model "models\bid_optimizer_latest.onnx"
}

Write-Host "✅ Go tests complete!" -ForegroundColor Green
```

### Create `setup.ps1`

```powershell
# setup.ps1 - Complete setup script
Write-Host "🚀 Setting up ML Pipeline" -ForegroundColor Cyan
Write-Host ("=" * 50)

# Create virtual environment
Write-Host "`n📦 Creating Python virtual environment..."
python -m venv venv

# Activate it
Write-Host "📦 Activating virtual environment..."
.\venv\Scripts\Activate.ps1

# Upgrade pip
Write-Host "`n📦 Upgrading pip..."
python -m pip install --upgrade pip

# Install packages
Write-Host "`n📦 Installing Python packages..."
pip install -r requirements.txt

Write-Host "`n✅ Setup complete!" -ForegroundColor Green
Write-Host "To activate environment later, run: .\venv\Scripts\Activate.ps1"
```

### Make Scripts Executable

```powershell
# Allow scripts to run
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

## 🎯 Part 7: Your First Real Training Run

Now let's train on your actual bid data!

### Step 1: Generate Training Data

You'll need to create the Go training data generator. For now, you can manually insert some test data:

```sql
-- Run this in pgAdmin or psql to create sample training data
INSERT INTO bid_training_data (
    campaign_id, floor_price, engagement_score, conversion_probability,
    device_type, segment_category, country, hour_of_day, day_of_week,
    historical_win_rate, historical_avg_bid, optimal_bid, actual_outcome, training_set
)
SELECT
    gen_random_uuid(),
    random() * 4 + 1,  -- floor_price between 1 and 5
    random(),  -- engagement_score
    random() * 0.3,  -- conversion_probability
    (ARRAY['mobile', 'desktop', 'tablet'])[floor(random() * 3 + 1)],
    (ARRAY['premium', 'standard', 'value'])[floor(random() * 3 + 1)],
    (ARRAY['US', 'GB', 'CA', 'AU'])[floor(random() * 4 + 1)],
    floor(random() * 24)::int,  -- hour_of_day
    floor(random() * 7)::int,  -- day_of_week
    random() * 0.6 + 0.2,  -- historical_win_rate
    random() * 4 + 1,  -- historical_avg_bid
    random() * 5 + 1.5,  -- optimal_bid (target)
    'won',
    CASE WHEN random() < 0.7 THEN 'train' 
         WHEN random() < 0.85 THEN 'validation' 
         ELSE 'test' END
FROM generate_series(1, 5000);  -- Generate 5000 samples
```

### Step 2: Train Model

```powershell
# Run the training script
.\train.ps1
```

### Step 3: Test in Go

```powershell
# Test Go integration
.\test-go.ps1
```

## 🔍 Troubleshooting

### Python Issues

**Error: "venv\Scripts\Activate.ps1 is not digitally signed"**
```powershell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

**Error: "No module named 'xgboost'"**
```powershell
# Make sure venv is activated (you see (venv) in prompt)
pip install xgboost
```

### Go Issues

**Error: "onnxruntime.dll not found"**
```powershell
# Make sure onnxruntime.dll is in your project directory or PATH
# Download again using the PowerShell script above
```

**Error: "CGO_ENABLED"**
```powershell
# Set environment variable
$env:CGO_ENABLED="1"
```

### Database Issues

**Error: "Connection refused"**
```powershell
# Check if PostgreSQL is running
Get-Service postgresql*

# Start if stopped
Start-Service postgresql-x64-15  # Adjust version number
```

## 📝 Next Steps

Once everything is working:

1. **Review your data**: Look at the `bid_training_data` table
2. **Adjust features**: Edit `train_model.py` to match your needs
3. **Train regularly**: Set up Windows Task Scheduler for daily retraining
4. **Monitor performance**: Track model metrics over time
5. **Deploy gradually**: A/B test new models

## 🎓 Quick Reference

### Daily Workflow

```powershell
# 1. Activate environment
.\venv\Scripts\Activate.ps1

# 2. Generate new data (once you have the Go tool)
go run .\cmd\training-data-generator\main.go --days 7

# 3. Train model
python train_model.py --config config.yaml

# 4. Test model
python test_model.py

# 5. Deploy (if tests pass)
Copy-Item "models\bid_optimizer_latest.onnx" -Destination "production\"
```

### Useful Commands

```powershell
# Check what's running
Get-Process | Where-Object {$_.Name -like "*python*"}
Get-Process | Where-Object {$_.Name -like "*postgres*"}

# Database connection test
psql -U postgres -d bidding_analysis -c "SELECT COUNT(*) FROM bid_training_data;"

# Check model file
Get-ChildItem models\*.onnx | Select-Object Name, Length, LastWriteTime
```

## 🎉 You're Ready!

You now have:
- ✅ Python 3.14 environment with ML packages
- ✅ XGBoost and ONNX working
- ✅ Go with ONNX Runtime
- ✅ PostgreSQL database
- ✅ Training and testing scripts
- ✅ Windows-friendly workflow

Start with the test scripts to make sure everything works, then move to real data!

---

**Need help?** Check the main guides or ask specific questions about any step.
