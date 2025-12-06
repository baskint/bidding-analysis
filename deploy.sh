#!/bin/bash
set -e

echo "🚀 Deploying ML Service and Go API to Cloud Run"

# Load environment variables
if [ -f .env.neon ]; then
    export $(grep -v '^#' .env.neon | xargs)
fi

PROJECT_ID="bidding-analysis"
REGION="us-central1"

# ============================================
# STAGE 1: Deploy ML Service
# ============================================
echo "📦 Stage 1: Deploying ML Service..."

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
  --memory=512Mi \
  --cpu=1 \
  --min-instances=0 \
  --max-instances=10 \
  --allow-unauthenticated \
  --timeout=60

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

# ============================================
# STAGE 2: Deploy Go Service
# ============================================
echo "📦 Stage 2: Deploying Go Service..."

# Build and push Docker image
echo "Building Go service Docker image..."
docker build -t gcr.io/${PROJECT_ID}/bidding-analysis:latest .
docker push gcr.io/${PROJECT_ID}/bidding-analysis:latest

# Create temporary env file (WITHOUT PORT - Cloud Run sets that automatically)
cat > /tmp/env.yaml << ENVEOF
DATABASE_URL: "${DATABASE_URL}"
ML_SERVICE_URL: "${ML_SERVICE_URL}"
OPENAI_API_KEY: "${OPENAI_API_KEY}"
ENVIRONMENT: "production"
ENVEOF

# Deploy Go service
echo "Deploying bidding-analysis service..."
gcloud run deploy bidding-analysis \
  --image=gcr.io/${PROJECT_ID}/bidding-analysis:latest \
  --project="${PROJECT_ID}" \
  --region="${REGION}" \
  --platform=managed \
  --env-vars-file=/tmp/env.yaml \
  --memory=512Mi \
  --cpu=1 \
  --min-instances=0 \
  --max-instances=10 \
  --allow-unauthenticated \
  --timeout=60

# Clean up
rm /tmp/env.yaml

echo "✅ Deployment complete!"
echo ""
echo "Services:"
echo "  ML Service: ${ML_SERVICE_URL}"
echo "  Go API: $(gcloud run services describe bidding-analysis --project=${PROJECT_ID} --region=${REGION} --format='value(status.url)')"