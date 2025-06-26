#!/bin/bash

# Frontend Deployment Script for Firebase Hosting
# This script builds and deploys the Next.js app to Firebase Hosting

set -e  # Exit on any error

echo "🚀 Starting Frontend Deployment to Firebase Hosting..."
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
PRODUCTION_URL="https://bidding-analysis-539382269313.us-central1.run.app"

echo "📋 Configuration:"
echo "  Production API URL: ${PRODUCTION_URL}"
echo ""

# Step 1: Verify environment files
echo "🔍 Step 1: Checking environment files..."
if [ ! -f .env.production ]; then
    echo -e "${RED}❌ .env.production file not found!${NC}"
    exit 1
fi

echo -e "${GREEN}✅ .env.production found${NC}"

# Check if .env.local exists and warn about potential conflicts
if [ -f .env.local ]; then
    if grep -q "NEXT_PUBLIC_API_URL" .env.local; then
        echo -e "${YELLOW}⚠️  Warning: .env.local contains NEXT_PUBLIC_API_URL which may override production settings${NC}"
        echo "Consider removing NEXT_PUBLIC_API_URL from .env.local for production builds"
        echo ""
    fi
fi

# Verify .env.production content
echo "📄 Checking .env.production content:"
if grep -q "NEXT_PUBLIC_API_URL.*bidding-analysis" .env.production; then
    echo -e "${GREEN}✅ Production API URL found in .env.production${NC}"
else
    echo -e "${RED}❌ Production API URL not found in .env.production${NC}"
    echo "Expected: NEXT_PUBLIC_API_URL=${PRODUCTION_URL}"
    exit 1
fi
echo ""

# Step 2: Clean previous builds
echo "🧹 Step 2: Cleaning previous builds..."
rm -rf .next out node_modules/.cache
echo -e "${GREEN}✅ Cleaned build artifacts${NC}"
echo ""

# Step 3: Build and export with production environment
echo "📦 Step 3: Building with production environment (static export enabled in next.config.js)..."
echo "Setting NEXT_PUBLIC_API_URL=${PRODUCTION_URL}"

# Build with explicit environment variable (exports automatically due to next.config.js)
NEXT_PUBLIC_API_URL="${PRODUCTION_URL}" yarn build

if [ $? -ne 0 ]; then
    echo -e "${RED}❌ Build failed!${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Build completed (static export generated)${NC}"
echo ""

# Step 4: Verify the production URL is in the built files
echo "🔍 Step 4: Verifying production URL in built files..."
if grep -r "bidding-analysis-539382269313" out/ > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Production URL found in built files${NC}"
else
    echo -e "${RED}❌ Production URL NOT found in built files${NC}"
    echo "The build may still be using localhost. Check your environment setup."
    
    # Check if localhost is still present
    if grep -r "localhost:8080" out/ > /dev/null 2>&1; then
        echo -e "${RED}❌ Found localhost:8080 in built files - environment variable not applied${NC}"
    fi
    exit 1
fi

# Check for any remaining localhost references
if grep -r "localhost:8080" out/ > /dev/null 2>&1; then
    echo -e "${YELLOW}⚠️  Warning: Found localhost:8080 references in built files${NC}"
fi
echo ""

# Step 5: Deploy to Firebase
echo "☁️  Step 5: Deploying to Firebase Hosting..."
yarn dlx firebase-tools deploy

if [ $? -eq 0 ]; then
    echo ""
    echo -e "${GREEN}✅ Deployment successful!${NC}"
    echo ""
    echo "🌐 Your frontend is live at:"
    echo "   https://bidding-analysis.web.app"
    echo "   https://bidding-analysis.firebaseapp.com"
    echo ""
    echo "🧪 Test the deployment:"
    echo "   1. Open the app in your browser"
    echo "   2. Check browser console for API URL (should show production URL)"
    echo "   3. Test authentication flow"
    echo ""
    echo "📊 Monitor deployment:"
    echo "   Firebase Console: https://console.firebase.google.com/project/bidding-analysis"
    echo ""
else
    echo -e "${RED}❌ Deployment failed!${NC}"
    exit 1
fi

echo "🎉 Frontend deployment completed successfully!" 
