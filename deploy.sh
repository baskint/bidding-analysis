#!/bin/bash

# Google Cloud Run Deployment Script
# Run this from your project root directory

set -e  # Exit on any error

echo "🚀 Starting deployment to Google Cloud Run..."

# Load environment variables from .env.neon
if [ -f .env.neon ]; then
    echo "📋 Loading environment variables from .env.neon..."
    source .env.neon
    echo "✅ Environment variables loaded"
else
    echo "❌ .env.neon file not found!"
    exit 1
fi

# Configuration
PROJECT_ID="bidding-analysis"
SERVICE_NAME="bidding-analysis"
REGION="us-central1"
IMAGE_NAME="gcr.io/${PROJECT_ID}/${SERVICE_NAME}:latest"

echo "📦 Step 1: Building Docker image..."
docker build -t ${IMAGE_NAME} .

if [ $? -ne 0 ]; then
    echo "❌ Docker build failed!"
    exit 1
fi

echo "📤 Step 2: Pushing image to Google Container Registry..."
docker push ${IMAGE_NAME}

if [ $? -ne 0 ]; then
    echo "❌ Docker push failed!"
    exit 1
fi

cat > /tmp/env.yaml << EOF
DATABASE_URL: "${DATABASE_URL}"
DB_HOST: "${DB_HOST}"
DB_USER: "${DB_USER}"
DB_PASSWORD: "${DB_PASSWORD}"
DB_NAME: "${DB_NAME}"
DB_PORT: "${DB_PORT}"
DB_SSL_MODE: "${DB_SSL_MODE}"
DB_STATEMENT_CACHE_MODE: "${DB_STATEMENT_CACHE_MODE}"
OPENAI_API_KEY: "${OPENAI_API_KEY}"
ALLOWED_ORIGINS: "${ALLOWED_ORIGINS}"
EOF

echo "☁️  Step 3: Deploying to Cloud Run..."
gcloud run deploy ${SERVICE_NAME} \
  --image ${IMAGE_NAME} \
  --platform managed \
  --region ${REGION} \
  --allow-unauthenticated \
  --port 8080 \
  --timeout 300 \
  --memory 1Gi \
  --cpu 1 \
  --concurrency 80 \
  --max-instances 10 \
  --project ${PROJECT_ID} \
 --env-vars-file /tmp/env.yaml

if [ $? -eq 0 ]; then
    echo "✅ Deployment successful!"
    echo ""
    echo "🌐 Your API is live at:"
    echo "https://${SERVICE_NAME}-539382269313.${REGION}.run.app"
    echo ""
    echo "🧪 Test endpoints:"
    echo "curl https://${SERVICE_NAME}-539382269313.${REGION}.run.app/health"
    echo "curl https://${SERVICE_NAME}-539382269313.${REGION}.run.app/trpc/debug"
    echo ""
    echo "📊 Monitor logs:"
    echo "gcloud logging read \"resource.type=cloud_run_revision AND resource.labels.service_name=${SERVICE_NAME}\" --project ${PROJECT_ID} --limit 20"
else
    echo "❌ Deployment failed!"
    exit 1
fi
