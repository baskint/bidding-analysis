#!/usr/bin/env python3
"""
ML Prediction Service
Provides XGBoost-based bid optimization predictions via REST API
"""

from flask import Flask, request, jsonify
import xgboost as xgb
import json
import numpy as np
import os
import logging
from pathlib import Path

# Configure logging
logging.basicConfig(
    level=logging.INFO, format="%(asctime)s - %(name)s - %(levelname)s - %(message)s"
)
logger = logging.getLogger(__name__)

app = Flask(__name__)

# Global model state
model = None
encoders = None
model_info = {}


def load_model():
  """Load the trained XGBoost model and encoders"""
  global model, encoders, model_info

  # Determine models directory
  if os.path.exists("../models"):
    models_dir = Path("../models")
  elif os.path.exists("models"):
    models_dir = Path("models")
  else:
    raise FileNotFoundError("Models directory not found")

  # Load XGBoost model
  model_path = models_dir / "bid_optimizer_latest.json"
  logger.info(f"Loading model from {model_path}")

  if not model_path.exists():
    raise FileNotFoundError(f"Model not found: {model_path}")

  model = xgb.Booster()
  model.load_model(str(model_path))

  # Load encoders
  encoders_path = models_dir / "bid_optimizer_latest_encoders.json"
  logger.info(f"Loading encoders from {encoders_path}")

  if not encoders_path.exists():
    raise FileNotFoundError(f"Encoders not found: {encoders_path}")

  with open(encoders_path, "r") as f:
    encoders = json.load(f)

  # Store model info
  model_info = {
      "num_features": model.num_features(),
      "num_trees": model.num_boosted_rounds(),
      "model_path": str(model_path),
      "encoders_path": str(encoders_path),
  }

  logger.info(f"✅ Model loaded successfully!")
  logger.info(f"   Features: {model_info['num_features']}")
  logger.info(f"   Trees: {model_info['num_trees']}")


@app.route("/", methods=["GET"])
def root():
  """Root endpoint"""
  return jsonify(
      {
          "service": "ML Prediction Service",
          "version": "1.0.0",
          "status": "running",
          "model_loaded": model is not None,
      }
  )


@app.route("/health", methods=["GET"])
def health():
  """Health check"""
  is_healthy = model is not None and encoders is not None
  return jsonify(
      {
          "status": "healthy" if is_healthy else "unhealthy",
          "model_loaded": model is not None,
      }
  ), (200 if is_healthy else 503)


@app.route("/predict", methods=["POST"])
def predict():
  """Predict optimal bid"""
  try:
    if model is None:
      return jsonify({"error": "Model not loaded"}), 503

    data = request.get_json()
    if not data:
      return jsonify({"error": "No JSON data"}), 400

    # Extract features
    features = [
        float(data.get("floor_price", 1.0)),
        float(data.get("engagement_score", 0.5)),
        float(data.get("conversion_probability", 0.1)),
        float(data.get("historical_win_rate", 0.4)),
        float(data.get("historical_avg_bid", 2.5)),
        float(data.get("historical_avg_win_price", 2.7)),
        encode_feature("device_type", data.get("device_type", "unknown")),
        encode_feature("segment_category", data.get(
            "segment_category", "standard")),
        float(data.get("hour_of_day", 12)),
        float(data.get("day_of_week", 1)),
        encode_feature("country", data.get("country", "US")),
        float(data.get("campaign_spend_last_7d", 100.0)),
        float(data.get("campaign_conversions_last_7d", 3.0)),
    ]

    # Predict
    dmatrix = xgb.DMatrix(
        np.array([features]),
        feature_names=[
            "floor_price",
            "engagement_score",
            "conversion_probability",
            "historical_win_rate",
            "historical_avg_bid",
            "historical_avg_win_price",
            "device_type_encoded",
            "segment_category_encoded",
            "hour_of_day",
            "day_of_week",
            "country_encoded",
            "campaign_spend_last_7d",
            "campaign_conversions_last_7d",
        ],
    )
    prediction = float(model.predict(dmatrix)[0])

    # Ensure above floor
    floor_price = float(data.get("floor_price", 1.0))
    if prediction < floor_price:
      prediction = floor_price * 1.01

    return jsonify(
        {
            "bid_price": round(prediction, 4),
            "confidence": 0.90,
            "strategy": "ml_optimized",
        }
    )

  except Exception as e:
    logger.error(f"Error: {e}", exc_info=True)
    return jsonify({"error": str(e)}), 500


def encode_feature(feature_name, value):
  """Encode categorical feature"""
  if encoders and feature_name in encoders:
    return float(encoders[feature_name].get(str(value), 0.0))
  return 0.0


if __name__ == "__main__":
  try:
    load_model()
  except Exception as e:
    logger.error(f"Failed to load model: {e}")

  port = int(os.environ.get("PORT", 5000))
  logger.info(f"🚀 Starting on port {port}")
  app.run(host="0.0.0.0", port=port, debug=False)
