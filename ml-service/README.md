# INSTALL INSTRUCTIONS

## Commands - Windows
```
python -m venv venv
.\venv\Scripts\activate
pip install -r requirements.txt
```

### Commandws - Linux
```
# Navigate to ml-service directory
cd ml-service

# Create virtual environment
python3 -m venv venv

# Activate (Linux/Mac uses 'source' not backslash)
source venv/bin/activate

# Install dependencies
pip install -r requirements.txt
```

### install files
```
# Verify the model files
ls -lh models/bid_optimizer_latest.json
ls -lh models/bid_optimizer_latest_encoders.json

# Copy to where ML service expects them
cp models/bid_optimizer_latest.json ../../models/
cp models/bid_optimizer_latest_encoders.json ../../models/

# Verify
ls -lh ../../models/bid_optimizer_latest.*
```
# Test 2: Make a prediction
curl -X POST http://localhost:5000/predict \
  -H "Content-Type: application/json" \
  -d '{
    "floor_price": 2.5,
    "engagement_score": 0.75,
    "conversion_probability": 0.2,
    "device_type": "desktop",
    "segment_category": "premium",
    "country": "US",
    "hour_of_day": 14,
    "day_of_week": 2
  }'

  autopep8
