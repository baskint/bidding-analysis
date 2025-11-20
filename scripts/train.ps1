# train.ps1 - Train ML model
param(
    [int]$Days = 30,
    [string]$Config = "config.yaml",
    [string]$Output = "models\bid_optimizer.onnx"
)

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  Training ML Model" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# Check if venv exists
if (-not (Test-Path "venv\Scripts\Activate.ps1")) {
    Write-Host "❌ Virtual environment not found" -ForegroundColor Red
    Write-Host "Run: .\setup.ps1" -ForegroundColor Yellow
    exit 1
}

# Activate venv
Write-Host "Activating virtual environment..." -ForegroundColor Yellow
& ".\venv\Scripts\Activate.ps1"

# Check if config exists
if (-not (Test-Path $Config)) {
    Write-Host "❌ Config file not found: $Config" -ForegroundColor Red
    Write-Host "Create config.yaml with your database settings" -ForegroundColor Yellow
    exit 1
}

# Create models directory
New-Item -ItemType Directory -Force -Path "models" | Out-Null

# Create logs directory
New-Item -ItemType Directory -Force -Path "logs" | Out-Null

# Get timestamp for log file
$timestamp = Get-Date -Format "yyyyMMdd_HHmmss"
$logFile = "logs\training_$timestamp.log"

# Run training
Write-Host ""
Write-Host "Training configuration:" -ForegroundColor Cyan
Write-Host "  Days of data: $Days"
Write-Host "  Config: $Config"
Write-Host "  Output: $Output"
Write-Host "  Log: $logFile"
Write-Host ""

Write-Host "Starting training..." -ForegroundColor Yellow
Write-Host "(This may take 2-10 minutes depending on data size)" -ForegroundColor Gray
Write-Host ""

python train_model.py --config $Config --days $Days --output $Output 2>&1 | Tee-Object -FilePath $logFile

if ($LASTEXITCODE -eq 0) {
    Write-Host ""
    Write-Host "========================================" -ForegroundColor Green
    Write-Host "  ✅ Training Complete!" -ForegroundColor Green
    Write-Host "========================================" -ForegroundColor Green
    Write-Host ""
    Write-Host "Model saved to: $Output" -ForegroundColor Cyan
    Write-Host "Log saved to: $logFile" -ForegroundColor Gray
    Write-Host ""
    Write-Host "Next steps:" -ForegroundColor Yellow
    Write-Host "  1. Review the training log"
    Write-Host "  2. Test the model: .\test-model.ps1"
    Write-Host "  3. If satisfied, deploy: Copy-Item '$Output' -Destination 'production\'"
    Write-Host ""
} else {
    Write-Host ""
    Write-Host "========================================" -ForegroundColor Red
    Write-Host "  ❌ Training Failed" -ForegroundColor Red
    Write-Host "========================================" -ForegroundColor Red
    Write-Host ""
    Write-Host "Check the log file for details: $logFile" -ForegroundColor Yellow
    Write-Host ""
    exit 1
}
