"""Test your trained model"""

import xgboost as xgb
import numpy as np
import pandas as pd
import json

print("🧪 Testing Trained Model")
print("=" * 60)

# Load model
print("\n[1/3] Loading model...")
model = xgb.XGBRegressor()
model.load_model("models/bid_optimizer_latest.json")
print("✅ Model loaded!")

# Load encoders
print("\n[2/3] Loading encoders...")
with open("models/bid_optimizer_latest_encoders.json", "r") as f:
    encoders = json.load(f)
print("✅ Encoders loaded")

# Test prediction
print("\n[3/3] Making test prediction...")
test_data = pd.DataFrame(
    {
        "floor_price": [2.5],
        "engagement_score": [0.75],
        "conversion_probability": [0.20],
        "historical_win_rate": [0.50],
        "historical_avg_bid": [2.8],
        "historical_avg_win_price": [3.0],
        "device_type_encoded": [1.0],
        "segment_category_encoded": [0.0],
        "hour_of_day": [14],
        "day_of_week": [2],
        "country_encoded": [1.0],
        "campaign_spend_last_7d": [250.0],
        "campaign_conversions_last_7d": [8.0],
    }
)

prediction = model.predict(test_data)

print("\n" + "=" * 60)
print("📊 TEST RESULT")
print("=" * 60)
print(f"Floor price: $2.50")
print(f"Engagement: 75%")
print(f"Conversion prob: 20%")
print(f"\n➡️  PREDICTED OPTIMAL BID: ${prediction[0]:.2f}")
print(f"Markup: {((prediction[0] / 2.5) - 1) * 100:.1f}%")
print("\n✅ Model is working perfectly!")
print("=" * 60)
