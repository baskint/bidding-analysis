# setup.ps1 - Complete ML Pipeline Setup for Windows
# Run this script to set up everything you need

param(
    [switch]$SkipVenv,
    [switch]$SkipPackages
)

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  ML Bid Optimization - Windows Setup  " -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# Check Python
Write-Host "[1/6] Checking Python installation..." -ForegroundColor Yellow
try {
    $pythonVersion = python --version 2>&1
    Write-Host "✅ $pythonVersion" -ForegroundColor Green
} catch {
    Write-Host "❌ Python not found. Please install Python 3.10 or higher" -ForegroundColor Red
    exit 1
}

# Check Go
Write-Host ""
Write-Host "[2/6] Checking Go installation..." -ForegroundColor Yellow
try {
    $goVersion = go version
    Write-Host "✅ $goVersion" -ForegroundColor Green
} catch {
    Write-Host "⚠️  Go not found. Install from: https://go.dev/dl/" -ForegroundColor Yellow
    Write-Host "   (You can continue without Go for now)" -ForegroundColor Gray
}

# Check PostgreSQL
Write-Host ""
Write-Host "[3/6] Checking PostgreSQL..." -ForegroundColor Yellow
try {
    $pgVersion = psql --version 2>&1
    Write-Host "✅ $pgVersion" -ForegroundColor Green
} catch {
    Write-Host "⚠️  PostgreSQL not found. Install from: https://www.postgresql.org/download/windows/" -ForegroundColor Yellow
    Write-Host "   (You can continue without PostgreSQL for testing)" -ForegroundColor Gray
}

if (-not $SkipVenv) {
    # Create virtual environment
    Write-Host ""
    Write-Host "[4/6] Creating Python virtual environment..." -ForegroundColor Yellow
    if (Test-Path "venv") {
        Write-Host "✅ Virtual environment already exists" -ForegroundColor Green
    } else {
        python -m venv venv
        Write-Host "✅ Virtual environment created" -ForegroundColor Green
    }

    # Activate virtual environment
    Write-Host ""
    Write-Host "[5/6] Activating virtual environment..." -ForegroundColor Yellow
    try {
        & ".\venv\Scripts\Activate.ps1"
        Write-Host "✅ Virtual environment activated" -ForegroundColor Green
    } catch {
        Write-Host "⚠️  Execution policy issue. Running fix..." -ForegroundColor Yellow
        Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser -Force
        & ".\venv\Scripts\Activate.ps1"
        Write-Host "✅ Virtual environment activated" -ForegroundColor Green
    }
} else {
    Write-Host ""
    Write-Host "[4/6] Skipping virtual environment creation" -ForegroundColor Gray
    Write-Host "[5/6] Skipping activation" -ForegroundColor Gray
}

if (-not $SkipPackages) {
    # Install Python packages
    Write-Host ""
    Write-Host "[6/6] Installing Python packages..." -ForegroundColor Yellow
    Write-Host "     This may take 2-5 minutes..." -ForegroundColor Gray
    
    # Upgrade pip first
    python -m pip install --upgrade pip --quiet
    
    # Install core ML packages
    Write-Host "     Installing XGBoost..." -ForegroundColor Gray
    pip install xgboost==2.0.3 --quiet
    
    Write-Host "     Installing scikit-learn..." -ForegroundColor Gray
    pip install scikit-learn==1.4.0 --quiet
    
    Write-Host "     Installing numpy and pandas..." -ForegroundColor Gray
    pip install numpy==1.26.3 pandas==2.2.0 --quiet
    
    Write-Host "     Installing ONNX..." -ForegroundColor Gray
    pip install onnxmltools==1.12.0 onnxruntime==1.16.3 skl2onnx==1.16.0 --quiet
    
    Write-Host "     Installing database drivers..." -ForegroundColor Gray
    pip install psycopg2-binary==2.9.9 --quiet
    
    Write-Host "     Installing utilities..." -ForegroundColor Gray
    pip install pyyaml==6.0.1 --quiet
    
    Write-Host "✅ All packages installed successfully!" -ForegroundColor Green
} else {
    Write-Host ""
    Write-Host "[6/6] Skipping package installation" -ForegroundColor Gray
}

# Create directories
Write-Host ""
Write-Host "Creating project directories..." -ForegroundColor Yellow
New-Item -ItemType Directory -Force -Path "models" | Out-Null
New-Item -ItemType Directory -Force -Path "logs" | Out-Null
New-Item -ItemType Directory -Force -Path "data" | Out-Null
Write-Host "✅ Directories created" -ForegroundColor Green

# Summary
Write-Host ""
Write-Host "========================================" -ForegroundColor Green
Write-Host "  ✅ Setup Complete!" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Green
Write-Host ""
Write-Host "Next steps:" -ForegroundColor Cyan
Write-Host "1. Edit config.yaml with your database credentials"
Write-Host "2. Run: python test_training.py (to test your setup)"
Write-Host "3. Run: python train_model.py (to train your first model)"
Write-Host ""
Write-Host "To activate the environment later, run:" -ForegroundColor Yellow
Write-Host "  .\venv\Scripts\Activate.ps1"
Write-Host ""
Write-Host "For help, see: WINDOWS_SETUP_GUIDE.md" -ForegroundColor Gray
Write-Host ""
