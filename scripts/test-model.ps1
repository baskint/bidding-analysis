# test-model.ps1 - Test trained model
param(
    [string]$Model = "models\bid_optimizer_latest.onnx"
)

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  Testing ML Model" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# Activate venv
if (Test-Path "venv\Scripts\Activate.ps1") {
    & ".\venv\Scripts\Activate.ps1"
}

# Check if model exists
if (-not (Test-Path $Model)) {
    Write-Host "❌ Model not found: $Model" -ForegroundColor Red
    Write-Host ""
    Write-Host "Available models:" -ForegroundColor Yellow
    Get-ChildItem -Path "models" -Filter "*.onnx" | ForEach-Object {
        Write-Host "  - $($_.Name)"
    }
    Write-Host ""
    Write-Host "Train a model first: .\train.ps1" -ForegroundColor Yellow
    exit 1
}

# Get model info
$modelInfo = Get-Item $Model
$modelSize = [math]::Round($modelInfo.Length / 1KB, 2)

Write-Host "Model Information:" -ForegroundColor Cyan
Write-Host "  Path: $Model"
Write-Host "  Size: $modelSize KB"
Write-Host "  Modified: $($modelInfo.LastWriteTime)"
Write-Host ""

# Test 1: Load model with ONNX Runtime
Write-Host "[1/3] Testing ONNX model loading..." -ForegroundColor Yellow

$testScript = @"
import onnxruntime as ort
import sys

try:
    session = ort.InferenceSession('$Model')
    print('✅ Model loaded successfully')
    
    # Get model info
    inputs = session.get_inputs()
    outputs = session.get_outputs()
    
    print(f'   Inputs: {len(inputs)}')
    for inp in inputs:
        print(f'     - {inp.name}: {inp.shape}')
    
    print(f'   Outputs: {len(outputs)}')
    for out in outputs:
        print(f'     - {out.name}: {out.shape}')
    
    sys.exit(0)
except Exception as e:
    print(f'❌ Failed to load model: {e}')
    sys.exit(1)
"@

$testScript | python

if ($LASTEXITCODE -ne 0) {
    Write-Host ""
    Write-Host "Model loading failed. Check if ONNX Runtime is installed:" -ForegroundColor Red
    Write-Host "  pip install onnxruntime" -ForegroundColor Yellow
    exit 1
}

# Test 2: Test prediction with sample data
Write-Host ""
Write-Host "[2/3] Testing predictions with sample data..." -ForegroundColor Yellow

$predTest = @"
import onnxruntime as ort
import numpy as np

try:
    session = ort.InferenceSession('$Model')
    input_name = session.get_inputs()[0].name
    
    # Create sample input (adjust feature count to match your model)
    # This is a sample bid request
    sample_input = np.array([[
        2.5,    # floor_price
        0.75,   # engagement_score
        0.15,   # conversion_probability
        0.45,   # historical_win_rate
        2.8,    # historical_avg_bid
        3.0,    # historical_avg_win_price
        1.0,    # device_type_encoded (mobile)
        0.0,    # segment_category_encoded (premium)
        14.0,   # hour_of_day (2 PM)
        2.0,    # day_of_week (Tuesday)
        1.0,    # country_encoded (US)
        150.0,  # campaign_spend_last_7d
        5.0     # campaign_conversions_last_7d
    ]], dtype=np.float32)
    
    # Run prediction
    result = session.run(None, {input_name: sample_input})
    predicted_bid = result[0][0][0]
    
    print(f'✅ Prediction successful')
    print(f'   Sample input: Floor price \$2.50, High engagement')
    print(f'   Predicted optimal bid: \${predicted_bid:.4f}')
    
    # Sanity check
    if predicted_bid < 0.1 or predicted_bid > 100:
        print(f'⚠️  Warning: Unusual prediction value')
    else:
        print(f'✅ Prediction looks reasonable')
    
except Exception as e:
    print(f'❌ Prediction failed: {e}')
    print(f'   Note: Feature count may not match your model')
    import sys
    sys.exit(1)
"@

$predTest | python

if ($LASTEXITCODE -ne 0) {
    Write-Host ""
    Write-Host "⚠️  Prediction test failed (may need to adjust feature count)" -ForegroundColor Yellow
}

# Test 3: Performance benchmark
Write-Host ""
Write-Host "[3/3] Running performance benchmark..." -ForegroundColor Yellow

$benchScript = @"
import onnxruntime as ort
import numpy as np
import time

try:
    session = ort.InferenceSession('$Model')
    input_name = session.get_inputs()[0].name
    
    # Create random input
    sample_input = np.random.randn(1, 13).astype(np.float32)
    
    # Warmup
    for _ in range(10):
        session.run(None, {input_name: sample_input})
    
    # Benchmark
    n_iterations = 1000
    start = time.time()
    
    for _ in range(n_iterations):
        session.run(None, {input_name: sample_input})
    
    elapsed = time.time() - start
    avg_time = (elapsed / n_iterations) * 1000  # ms
    throughput = n_iterations / elapsed
    
    print(f'✅ Performance benchmark complete')
    print(f'   Iterations: {n_iterations}')
    print(f'   Average inference time: {avg_time:.2f} ms')
    print(f'   Throughput: {throughput:.0f} predictions/second')
    
    if avg_time < 10:
        print(f'✅ Excellent performance (< 10ms)')
    elif avg_time < 50:
        print(f'✅ Good performance (< 50ms)')
    else:
        print(f'⚠️  Slower than expected (may need optimization)')
    
except Exception as e:
    print(f'❌ Benchmark failed: {e}')
"@

$benchScript | python

# Summary
Write-Host ""
Write-Host "========================================" -ForegroundColor Green
Write-Host "  ✅ Model Testing Complete!" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Green
Write-Host ""
Write-Host "Your model is ready to use!" -ForegroundColor Cyan
Write-Host ""
Write-Host "Next steps:" -ForegroundColor Yellow
Write-Host "  1. Integrate with Go service (see onnx_predictor.go)"
Write-Host "  2. Test in staging environment"
Write-Host "  3. Deploy to production"
Write-Host ""
Write-Host "To use in Go:" -ForegroundColor Cyan
Write-Host '  predictor, err := mlonnx.NewONNXPredictor("$Model", "models\bid_optimizer_latest_encoders.json")'
Write-Host ""
