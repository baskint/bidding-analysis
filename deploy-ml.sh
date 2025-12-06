#!/bin/bash
set -e

echo "🚀 Deploying ML Service to Cloud Run"

PROJECT_ID="bidding-analysis"
REGION="us-central1"

cd ml-service

# Remove symlinks first
echo "Copying model files..."
rm -f bid_optimizer_latest.json bid_optimizer_latest_encoders.json

# Now copy the actual files
cp ../models/bid_optimizer_latest.json bid_optimizer_latest.json
cp ../models/bid_optimizer_latest_encoders.json bid_optimizer_latest_encoders.json

# Verify they're real files
ls -lh bid_optimizer_latest*.json

# Deploy ML service
echo "Deploying ml-predictor service..."
gcloud run deploy ml-predictor \
  --source . \
  --project="${PROJECT_ID}" \
  --region="${REGION}" \
  --platform=managed \
  --memory=2Gi \
  --cpu=2 \
  --min-instances=0 \
  --max-instances=10 \
  --allow-unauthenticated \
  --timeout=120

# Get ML service URL
ML_SERVICE_URL=$(gcloud run services describe ml-predictor \
  --project="${PROJECT_ID}" \
  --region="${REGION}" \
  --format='value(status.url)')

echo "✅ ML Service deployed at: ${ML_SERVICE_URL}"

# Clean up - restore symlinks for local dev
rm bid_optimizer_latest.json bid_optimizer_latest_encoders.json
ln -sf ../models/bid_optimizer_latest.json bid_optimizer_latest.json
ln -sf ../models/bid_optimizer_latest_encoders.json bid_optimizer_latest_encoders.json

cd ..
