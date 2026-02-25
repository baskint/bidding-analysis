# Scripts

This directory contains tooling for data generation, ML training, and bid simulation.

---

## Directory Structure

```
scripts/
├── py_scripts/              # Python ML training pipeline
│   ├── train_model.py       # Full XGBoost training pipeline (DB or synthetic)
│   ├── test_training.py     # Verify your Python/ML setup works
│   ├── test_predictions.py  # Test model predictions after training
│   ├── export_to_onnx.py    # Export trained model to ONNX format
│   ├── fix_onnx_export.py   # Fixes for ONNX compatibility issues
│   ├── convert_model_to_binary.py  # Convert model to binary format
│   ├── config.yaml          # Database + training configuration
│   └── requirements.txt     # Python dependencies for this pipeline
│
├── bid-simulator.ts         # Continuously submits live bid events to the API
├── generate-test-data.ts    # Bulk-generates historical bid events (batched)
├── onnx_predictor.go        # Go inference code using ONNX Runtime
├── Makefile.ml              # Make targets for the full ML pipeline
├── setup.ps1                # Windows setup script (Python + packages)
├── test-model.ps1           # Windows: test model predictions
├── train.ps1                # Windows: run training pipeline
├── test_data.sql            # SQL schema + seed data for training table
├── requirements.txt         # Top-level Python dependencies
└── package.json             # Node.js dependencies for TypeScript scripts
```

---

## Python ML Training Pipeline (`py_scripts/`)

A more advanced training pipeline than `ml-service/train_13_features.py`. Supports loading data from PostgreSQL or generating synthetic data.

### Setup

```bash
cd scripts
python3 -m venv venv
source venv/bin/activate        # Linux/Mac
# .\venv\Scripts\Activate.ps1   # Windows (or run setup.ps1)
pip install -r py_scripts/requirements.txt
```

### Train with synthetic data (no database needed)
```bash
cd py_scripts
python train_model.py --synthetic --samples 5000
```

### Train with database data
```bash
# Edit config.yaml with your PostgreSQL credentials first
python train_model.py --days 30 --output models/bid_optimizer.onnx
```

### Verify setup
```bash
python test_training.py
```

### Test predictions after training
```bash
python test_predictions.py
```

### Key differences from `ml-service/train_13_features.py`

| Feature | `train_13_features.py` | `py_scripts/train_model.py` |
|---|---|---|
| Data source | CSV file or synthetic | PostgreSQL or synthetic |
| Feature engineering | None | Temporal + interaction features |
| Hyperparameters | Basic | Tuned (300 trees, early stopping) |
| Output format | XGBoost JSON | XGBoost JSON + ONNX attempt |
| Encoder type | Fixed integer mapping | Frequency encoding |
| Use case | Quick local training | Production pipeline |

---

## TypeScript Data Scripts

Requires Node.js. Install dependencies first:

```bash
cd scripts
npm install
```

Both scripts require a running backend API and a valid JWT token:
```bash
export API_URL="http://localhost:8080"
export AUTH_TOKEN="your-jwt-token"
```

### `generate-test-data.ts` — Bulk historical data generator

Generates and submits 1000 bid events in batches of 50. Uses weighted distributions to produce realistic data (60% mobile, 40% US traffic, etc.).

```bash
npx ts-node generate-test-data.ts
```

### `bid-simulator.ts` — Live bid simulator

Continuously submits one bid every 2–5 seconds until stopped. Useful for testing real-time dashboards.

```bash
npx ts-node bid-simulator.ts
# Press Ctrl+C to stop
```

---

## Makefile (`Makefile.ml`)

Orchestrates the full ML pipeline. Run from the `scripts/` directory:

```bash
make help           # Show all available commands

make setup          # Create venv and install dependencies
make train          # Train model from database (last 30 days)
make evaluate       # Evaluate model on test set
make export         # Export model to ONNX
make test-go        # Run Go ONNX inference tests
make monitor        # Check model performance (last 24h)
make clean          # Remove generated model/log files

# Advanced
make train-tuned        # Train with hyperparameter tuning
make deploy-gradual     # Gradual rollout (10% → 30% → 50% → 100%)
make check-drift        # Detect model drift vs baseline
make full-pipeline      # generate-data → train → evaluate → test-go
```

---

## Windows Setup (`setup.ps1`)

Automates the full Windows setup: Python venv, all ML packages, and directory structure.

```powershell
.\setup.ps1
# Skip steps if needed:
.\setup.ps1 -SkipVenv
.\setup.ps1 -SkipPackages
```

---

## Go ONNX Predictor (`onnx_predictor.go`)

Reference implementation for running the trained model inside a Go service using the ONNX Runtime. Handles feature extraction, categorical encoding, and batch predictions.

See `ml-recommendations.md` for the full architectural context and rationale for the Python-train / Go-infer hybrid approach.

---

## Database Schema (`test_data.sql`)

Contains the `bid_training_data` table schema and sample seed data used by `py_scripts/train_model.py` when loading from PostgreSQL.

```bash
psql -U postgres -d bidding_analysis -f test_data.sql
```
