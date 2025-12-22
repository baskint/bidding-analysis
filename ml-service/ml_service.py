# ml-service/ml_service.py
from flask import Flask, request, jsonify
import xgboost as xgb
import numpy as np
import json
import logging
import os

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

app = Flask(__name__)

# Global variables for model
model = None
encoders = None
model_loaded = False
feature_names = [
    'floor_price', 'engagement_score', 'conversion_probability',
    'historical_win_rate', 'historical_avg_bid', 'historical_avg_win_price',
    'device_type_encoded', 'segment_category_encoded', 'hour_of_day',
    'day_of_week', 'country_encoded', 'campaign_spend_last_7d',
    'campaign_conversions_last_7d'
]


def load_model():
  """Load the XGBoost model and encoders"""
  global model, encoders, model_loaded

  try:
    # In production (Cloud Run), models are in ./models/
    # In local dev, models are in ../models/
    if os.path.exists('models/bid_optimizer_latest.json'):
      model_path = 'models/bid_optimizer_latest.json'
      encoders_path = 'models/bid_optimizer_latest_encoders.json'
    else:
      model_path = '../models/bid_optimizer_latest.json'
      encoders_path = '../models/bid_optimizer_latest_encoders.json'

    logger.info(f"Loading model from {model_path}")
    model = xgb.Booster()
    model.load_model(model_path)

    logger.info(f"Loading encoders from {encoders_path}")
    with open(encoders_path, 'r') as f:
      encoders = json.load(f)

    logger.info("✅ Model loaded successfully!")
    logger.info(f"   Features: {model.num_features()}")
    logger.info(f"   Trees: {model.num_boosted_rounds()}")

    model_loaded = True
    return True

  except FileNotFoundError as e:
    logger.error(f"Model files not found: {e}")
    model_loaded = False
    return False
  except Exception as e:
    logger.error(f"Failed to load model: {e}")
    model_loaded = False
    return False


@app.route('/health', methods=['GET'])
def health():
  """Health check endpoint"""
  return jsonify({
      'status': 'healthy' if model_loaded else 'unhealthy',
      'model_loaded': model_loaded
  })


@app.route('/predict', methods=['POST'])
def predict():
  """Prediction endpoint"""
  if not model_loaded:
    return jsonify({'error': 'Model not loaded'}), 503

  try:
    data = request.json
    features = data.get('features', {})

    # Encode categorical features
    device_type_encoded = encoders['device_type'].get(
        features.get('device_type', 'unknown'), 0)
    segment_encoded = encoders['segment_category'].get(
        features.get('segment_category', 'unknown'), 0)
    country_encoded = encoders['country'].get(
        features.get('country', 'unknown'), 0)

    # Create feature array in correct order
    feature_array = np.array([[
        features.get('floor_price', 0.0),
        features.get('engagement_score', 0.0),
        features.get('conversion_probability', 0.0),
        features.get('historical_win_rate', 0.0),
        features.get('historical_avg_bid', 0.0),
        features.get('historical_avg_win_price', 0.0),
        device_type_encoded,
        segment_encoded,
        features.get('hour_of_day', 0),
        features.get('day_of_week', 0),
        country_encoded,
        features.get('campaign_spend_last_7d', 0.0),
        features.get('campaign_conversions_last_7d', 0.0)
    ]], dtype=np.float32)

    # Make prediction with feature names
    dmatrix = xgb.DMatrix(feature_array, feature_names=feature_names)
    prediction = float(model.predict(dmatrix)[0])

    return jsonify({
        'predicted_bid': prediction,
        'model_version': 'bid_optimizer_latest'
    })

  except Exception as e:
    logger.error(f"Prediction error: {e}")
    return jsonify({'error': str(e)}), 500


if __name__ == '__main__':
  logger.info("🚀 Starting on port 5001")
  load_model()
  app.run(host='0.0.0.0', port=5001, debug=False)
else:
  # When running with gunicorn
  load_model()
