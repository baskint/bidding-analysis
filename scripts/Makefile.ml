# Makefile for ML Pipeline
# Manages training, deployment, and monitoring of ML models

.PHONY: help setup train evaluate export deploy monitor clean

# Variables
PYTHON := python3
PIP := pip3
CONFIG := config.yaml
MODELS_DIR := models
DATA_DIR := data
LOGS_DIR := logs
VENV := venv

# Default target
help:
	@echo "ML Pipeline Commands:"
	@echo "  make setup        - Set up Python environment"
	@echo "  make train        - Train new model"
	@echo "  make evaluate     - Evaluate model performance"
	@echo "  make export       - Export model to ONNX"
	@echo "  make deploy       - Deploy model to production"
	@echo "  make test-go      - Test Go inference"
	@echo "  make monitor      - Monitor model performance"
	@echo "  make clean        - Clean generated files"
	@echo ""
	@echo "Data Commands:"
	@echo "  make generate-data    - Generate training data from DB"
	@echo "  make validate-data    - Validate training data quality"
	@echo ""
	@echo "Development:"
	@echo "  make jupyter      - Start Jupyter notebook"
	@echo "  make lint         - Run Python linters"

# Setup Python environment
setup:
	@echo "Setting up Python environment..."
	$(PYTHON) -m venv $(VENV)
	. $(VENV)/bin/activate && $(PIP) install --upgrade pip
	. $(VENV)/bin/activate && $(PIP) install -r requirements.txt
	@echo "✅ Setup complete! Activate with: source venv/bin/activate"

# Generate training data from PostgreSQL
generate-data:
	@echo "Generating training data..."
	go run cmd/training-data-generator/main.go --days=30
	@echo "✅ Training data generated"

# Validate data quality
validate-data:
	@echo "Validating training data..."
	. $(VENV)/bin/activate && $(PYTHON) scripts/validate_data.py
	@echo "✅ Data validation complete"

# Train model
train:
	@echo "Training model..."
	mkdir -p $(MODELS_DIR) $(LOGS_DIR)
	. $(VENV)/bin/activate && $(PYTHON) train_model.py \
		--config $(CONFIG) \
		--days 30 \
		--output $(MODELS_DIR)/bid_optimizer.onnx \
		2>&1 | tee $(LOGS_DIR)/training_$$(date +%Y%m%d_%H%M%S).log
	@echo "✅ Training complete"

# Train with hyperparameter tuning
train-tuned:
	@echo "Training with hyperparameter tuning..."
	mkdir -p $(MODELS_DIR) $(LOGS_DIR)
	. $(VENV)/bin/activate && $(PYTHON) scripts/hyperparameter_tuning.py \
		--config $(CONFIG) \
		--trials 50 \
		2>&1 | tee $(LOGS_DIR)/tuning_$$(date +%Y%m%d_%H%M%S).log
	@echo "✅ Hyperparameter tuning complete"

# Evaluate model on test set
evaluate:
	@echo "Evaluating model..."
	. $(VENV)/bin/activate && $(PYTHON) scripts/evaluate_model.py \
		--model $(MODELS_DIR)/bid_optimizer_latest.onnx \
		--test-data $(DATA_DIR)/test_set.csv
	@echo "✅ Evaluation complete"

# Compare models (A/B testing preparation)
compare-models:
	@echo "Comparing models..."
	. $(VENV)/bin/activate && $(PYTHON) scripts/compare_models.py \
		--model1 $(MODELS_DIR)/bid_optimizer_v1.onnx \
		--model2 $(MODELS_DIR)/bid_optimizer_latest.onnx
	@echo "✅ Model comparison complete"

# Test Go inference
test-go:
	@echo "Testing Go inference..."
	go test -v ./internal/mlonnx/...
	@echo "Testing inference speed..."
	go run cmd/benchmark/main.go --model $(MODELS_DIR)/bid_optimizer_latest.onnx
	@echo "✅ Go tests complete"

# Build Go service with ONNX support
build-go:
	@echo "Building Go service..."
	go build -o bin/bidding-server cmd/server/main.go
	@echo "✅ Build complete"

# Deploy model to production
deploy:
	@echo "Deploying model..."
	@echo "1. Copying model to production..."
	cp $(MODELS_DIR)/bid_optimizer_latest.onnx /app/models/production.onnx
	cp $(MODELS_DIR)/bid_optimizer_latest_encoders.json /app/models/production_encoders.json
	@echo "2. Updating model version in database..."
	. $(VENV)/bin/activate && $(PYTHON) scripts/update_model_version.py \
		--model-path /app/models/production.onnx
	@echo "3. Restarting service..."
	@# In production, use: systemctl restart bidding-service
	@echo "✅ Deployment complete"

# Deploy with gradual rollout (A/B testing)
deploy-gradual:
	@echo "Starting gradual rollout..."
	. $(VENV)/bin/activate && $(PYTHON) scripts/gradual_rollout.py \
		--model $(MODELS_DIR)/bid_optimizer_latest.onnx \
		--initial-traffic 10 \
		--increment 20 \
		--interval 24
	@echo "✅ Gradual rollout initiated"

# Monitor model performance
monitor:
	@echo "Monitoring model performance..."
	. $(VENV)/bin/activate && $(PYTHON) scripts/monitor_model.py \
		--lookback-hours 24
	@echo "✅ Monitoring complete"

# Real-time monitoring dashboard
monitor-live:
	@echo "Starting live monitoring dashboard..."
	. $(VENV)/bin/activate && streamlit run scripts/monitoring_dashboard.py

# Start MLflow tracking server
mlflow-server:
	@echo "Starting MLflow tracking server..."
	. $(VENV)/bin/activate && mlflow ui --host 0.0.0.0 --port 5000

# Start Jupyter notebook
jupyter:
	@echo "Starting Jupyter notebook..."
	. $(VENV)/bin/activate && jupyter notebook

# Continuous training (scheduled job)
retrain:
	@echo "Running continuous training..."
	. $(VENV)/bin/activate && $(PYTHON) scripts/continuous_training.py
	@echo "✅ Continuous training cycle complete"

# Run all linters
lint:
	@echo "Running linters..."
	. $(VENV)/bin/activate && black train_model.py scripts/*.py
	. $(VENV)/bin/activate && flake8 train_model.py scripts/
	. $(VENV)/bin/activate && mypy train_model.py
	@echo "✅ Linting complete"

# Run tests
test-python:
	@echo "Running Python tests..."
	. $(VENV)/bin/activate && pytest tests/ -v
	@echo "✅ Tests complete"

# Generate training data report
data-report:
	@echo "Generating data report..."
	. $(VENV)/bin/activate && $(PYTHON) scripts/generate_data_report.py \
		--output $(LOGS_DIR)/data_report_$$(date +%Y%m%d).html
	@echo "✅ Report generated: $(LOGS_DIR)/data_report_$$(date +%Y%m%d).html"

# Clean generated files
clean:
	@echo "Cleaning generated files..."
	rm -rf $(MODELS_DIR)/*.onnx
	rm -rf $(MODELS_DIR)/*_encoders.json
	rm -rf $(LOGS_DIR)/*.log
	rm -rf __pycache__ **/__pycache__
	rm -rf .pytest_cache
	@echo "✅ Cleanup complete"

# Clean everything including venv
clean-all: clean
	@echo "Removing virtual environment..."
	rm -rf $(VENV)
	@echo "✅ Full cleanup complete"

# Full pipeline: data generation -> training -> evaluation -> deployment
full-pipeline:
	@echo "Running full ML pipeline..."
	$(MAKE) generate-data
	$(MAKE) validate-data
	$(MAKE) train
	$(MAKE) evaluate
	$(MAKE) test-go
	@echo "Pipeline complete! Review results before deploying."
	@echo "To deploy: make deploy"

# Docker commands
docker-build:
	@echo "Building Docker image..."
	docker build -t bidding-ml-trainer -f Dockerfile.ml .
	@echo "✅ Docker image built"

docker-train:
	@echo "Training in Docker..."
	docker run --rm \
		-v $(PWD)/models:/app/models \
		-v $(PWD)/config.yaml:/app/config.yaml \
		bidding-ml-trainer python train_model.py
	@echo "✅ Docker training complete"

# Database operations
db-backup:
	@echo "Backing up training data..."
	pg_dump -h localhost -U postgres bidding_analysis \
		-t bid_training_data \
		-t ml_model_metadata \
		> $(DATA_DIR)/backup_$$(date +%Y%m%d).sql
	@echo "✅ Database backup complete"

# Feature store sync
sync-features:
	@echo "Syncing feature store..."
	. $(VENV)/bin/activate && $(PYTHON) scripts/sync_feature_store.py
	@echo "✅ Feature store synced"

# Check model drift
check-drift:
	@echo "Checking for model drift..."
	. $(VENV)/bin/activate && $(PYTHON) scripts/check_model_drift.py \
		--baseline-date 2024-01-01 \
		--current-date $$(date +%Y-%m-%d)
	@echo "✅ Drift check complete"

# Performance benchmarks
benchmark:
	@echo "Running performance benchmarks..."
	. $(VENV)/bin/activate && $(PYTHON) scripts/benchmark_inference.py
	go run cmd/benchmark/main.go --iterations 10000
	@echo "✅ Benchmarks complete"

.DEFAULT_GOAL := help
