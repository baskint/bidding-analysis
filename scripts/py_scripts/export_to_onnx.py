#!/usr/bin/env python3
"""
Export XGBoost model to ONNX format for Go
Requires: pip install onnxmltools skl2onnx onnxruntime
"""

import xgboost as xgb
import numpy as np
from pathlib import Path

# Try to import ONNX conversion libraries
try:
  import onnxmltools
  from onnxmltools.convert import convert_xgboost
  from skl2onnx.common.data_types import FloatTensorType
except ImportError:
  print("❌ Missing dependencies. Please install:")
  print("   pip install onnxmltools skl2onnx onnxruntime")
  exit(1)


def export_to_onnx():
  """Export XGBoost model to ONNX format"""

  models_dir = Path("models")
  json_model = models_dir / "bid_optimizer_latest.json"

  if not json_model.exists():
    print(f"❌ Model not found: {json_model}")
    return

  print(f"📖 Loading XGBoost model from {json_model}")
  model = xgb.Booster()
  model.load_model(str(json_model))

  print(
      f"✅ Model loaded: {model.num_features()} features, {model.num_boosted_rounds()} rounds"
  )

  # Define input type (13 features, all float32)
  initial_type = [("float_input", FloatTensorType([None, 13]))]

  print("🔄 Converting to ONNX format...")
  try:
    onnx_model = convert_xgboost(
        model, initial_types=initial_type, target_opset=12)

    # Save ONNX model
    onnx_path = models_dir / "bid_optimizer_latest.onnx"

    with open(onnx_path, "wb") as f:
      f.write(onnx_model.SerializeToString())

    size_mb = onnx_path.stat().st_size / (1024 * 1024)
    print(f"✅ ONNX model saved to {onnx_path} ({size_mb:.2f} MB)")

    # Test the ONNX model
    print("\n🧪 Testing ONNX model...")
    import onnxruntime as rt

    sess = rt.InferenceSession(str(onnx_path))
    input_name = sess.get_inputs()[0].name

    # Test prediction with dummy data
    test_input = np.array(
        [
            [
                2.5,  # floor_price
                0.75,  # engagement_score
                0.20,  # conversion_probability
                0.50,  # historical_win_rate
                2.80,  # historical_avg_bid
                3.00,  # historical_avg_win_price
                1.0,  # device_type (encoded)
                2.0,  # segment_category (encoded)
                14,  # hour_of_day
                2,  # day_of_week
                1.0,  # country (encoded)
                250.0,  # campaign_spend_last_7d
                8.0,  # campaign_conversions_last_7d
            ]
        ],
        dtype=np.float32,
    )

    pred = sess.run(None, {input_name: test_input})
    print(f"✅ Test prediction: ${pred[0][0]:.2f}")

    print("\n" + "=" * 60)
    print("🎉 Success! ONNX model is ready for Go!")
    print("=" * 60)
    print("\nNext steps:")
    print("1. Install Go ONNX runtime:")
    print("   go get github.com/yalue/onnxruntime_go")
    print("\n2. Update your Go code to use ONNX predictor")
    print(f"\n3. Model file: {onnx_path}")

  except Exception as e:
    print(f"❌ ONNX conversion failed: {e}")
    print("\nTroubleshooting:")
    print("- Make sure all packages are up to date:")
    print("  pip install --upgrade onnxmltools skl2onnx onnxruntime xgboost")


if __name__ == "__main__":
  export_to_onnx()
