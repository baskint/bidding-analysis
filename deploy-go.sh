#!/bin/bash
set -e

echo "🚀 Deploying Go Service to Cloud Run"

# Load environment variables
if [ -f .env.neon ]; then
    export $(grep -v '^#' .env.neon | xargs)
fi

PROJECT_ID="bidding-analysis"
REGION="us-central1"

# Get ML service URL
ML_SERVICE_URL=$(gcloud run services describe ml-predictor \
  --project="${PROJECT_ID}" \
  --region="${REGION}" \
  --format='value(status.url)' 2>/dev/null || echo "")

if [ -z "$ML_SERVICE_URL" ]; then
    echo "⚠️  Warning: ML service not found. Deploy it first with ./deploy-ml.sh"
    echo "Continuing anyway..."
fi

# Build and push Docker image
echo "Building Go service Docker image..."
docker build -t gcr.io/${PROJECT_ID}/bidding-analysis:latest .
docker push gcr.io/${PROJECT_ID}/bidding-analysis:latest

# Create temporary env file with all required variables
cat > /tmp/env.yaml << ENVEOF
ML_SERVICE_URL: "${ML_SERVICE_URL}"
OPENAI_API_KEY: "${OPENAI_API_KEY}"
ENVIRONMENT: "production"
DB_HOST: "${DB_HOST}"
DB_PORT: "${DB_PORT}"
DB_NAME: "${DB_NAME}"
DB_USER: "${DB_USER}"
DB_PASSWORD: "${DB_PASSWORD}"
DB_SSL_MODE: "${DB_SSL_MODE}"
JWT_SECRET: "${JWT_SECRET}"
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

GO_API_URL=$(gcloud run services describe bidding-analysis \
  --project="${PROJECT_ID}" \
  --region="${REGION}" \
  --format='value(status.url)')

echo "✅ Deployment complete!"
echo ""
echo "Services:"
echo "  ML Service: ${ML_SERVICE_URL}"
echo "  Go API: ${GO_API_URL}"
