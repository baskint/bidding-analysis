#!/usr/bin/env python3
"""
Test script for ML Prediction Service
"""

import requests
import json

BASE_URL = "http://localhost:5001"

def test_health():
    """Test health endpoint"""
    print("🧪 Testing /health...")
    response = requests.get(f"{BASE_URL}/health")
    print(f"   Status: {response.status_code}")
    print(f"   Response: {response.json()}")
    assert response.status_code == 200
    assert response.json()["status"] == "healthy"
    print("   ✅ Health check passed")

def test_predict():
    """Test prediction endpoint"""
    print("\n🧪 Testing /predict...")
    
    test_data = {
        "floor_price": 2.5,
        "engagement_score": 0.75,
        "conversion_probability": 0.2,
        "historical_win_rate": 0.5,
        "historical_avg_bid": 2.8,
        "historical_avg_win_price": 3.0,
        "device_type": "desktop",
        "segment_category": "premium",
        "country": "US",
        "hour_of_day": 14,
        "day_of_week": 2,
        "campaign_spend_last_7d": 250.0,
        "campaign_conversions_last_7d": 8.0
    }
    
    response = requests.post(
        f"{BASE_URL}/predict",
        json={"features": test_data},
        headers={"Content-Type": "application/json"}
    )

    print(f"   Status: {response.status_code}")
    result = response.json()
    print(f"   Response: {json.dumps(result, indent=2)}")

    assert response.status_code == 200
    assert "predicted_bid" in result
    assert result["predicted_bid"] > 0

    print(f"   ✅ Prediction passed (bid: ${result['predicted_bid']:.2f})")

def test_multiple_scenarios():
    """Test multiple prediction scenarios"""
    print("\n🧪 Testing multiple scenarios...")
    
    scenarios = [
        {
            "name": "Low engagement",
            "data": {"floor_price": 1.0, "engagement_score": 0.2, "conversion_probability": 0.05}
        },
        {
            "name": "High engagement",
            "data": {"floor_price": 2.0, "engagement_score": 0.9, "conversion_probability": 0.3}
        },
        {
            "name": "Premium segment",
            "data": {"floor_price": 3.0, "segment_category": "premium", "device_type": "mobile"}
        }
    ]
    
    for scenario in scenarios:
        response = requests.post(f"{BASE_URL}/predict", json={"features": scenario["data"]})
        result = response.json()
        print(f"   {scenario['name']}: ${result['predicted_bid']:.2f}")
    
    print("   ✅ Multiple scenarios passed")

if __name__ == "__main__":
    print("🚀 Testing ML Prediction Service")
    print("=" * 50)
    
    try:
        test_health()
        test_predict()
        test_multiple_scenarios()
        
        print("\n" + "=" * 50)
        print("✅ All tests passed!")
        print("=" * 50)
        
    except requests.exceptions.ConnectionError:
        print("\n❌ Error: Could not connect to service")
        print("   Make sure the service is running:")
        print("   python ml_service.py")
        
    except Exception as e:
        print(f"\n❌ Test failed: {e}")
