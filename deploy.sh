#!/bin/bash

# Google Cloud Run Deployment Script
# Deploys both ML service and Go API service

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
REGION="us-central1"
GO_SERVICE_NAME="bidding-analysis"
ML_SERVICE_NAME="ml-predictor"
GO_IMAGE_NAME="gcr.io/${PROJECT_ID}/${GO_SERVICE_NAME}:latest"

# ============================================
# STEP 1: Deploy ML Service
# ============================================
echo ""
echo "🤖 Step 1: Deploying ML Prediction Service..."
cd ml-service

gcloud run deploy ${ML_SERVICE_NAME} \
  --source . \
  --region=${REGION} \
  --platform=managed \
  --memory=512Mi \
  --cpu=1 \
  --max-instances=10 \
  --timeout=60 \
  --allow-unauthenticated \
  --project=${PROJECT_ID}

if [ $? -ne 0 ]; then
    echo "❌ ML service deployment failed!"
    exit 1
fi

# Get ML service URL
export ML_SERVICE_URL=$(gcloud run services describe ${ML_SERVICE_NAME} \
  --region=${REGION} \
  --format='value(status.url)' \
  --project=${PROJECT_ID})

echo "✅ ML service deployed at: $ML_SERVICE_URL"

# Return to project root
cd ..

# ============================================
# STEP 2: Build and Deploy Go Service
# ============================================
echo ""
echo "📦 Step 2: Building Go service Docker image..."
docker build -t ${GO_IMAGE_NAME} .

if [ $? -ne 0 ]; then
    echo "❌ Docker build failed!"
    exit 1
fi

echo "📤 Step 3: Pushing image to Google Container Registry..."
docker push ${GO_IMAGE_NAME}

if [ $? -ne 0 ]; then
    echo "❌ Docker push failed!"
    exit 1
fi