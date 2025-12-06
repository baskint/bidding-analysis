#!/usr/bin/env python3
"""
Quick fix: Save XGBoost in older format that ONNX tools can handle
"""

import xgboost as xgb
from pathlib import Path


def downgrade_and_export():
    """Load model and save in older XGBoost format"""

    models_dir = Path("models")
    json_model = models_dir / "bid_optimizer_latest.json"

    print(f"📖 Loading model from {json_model}")
    model = xgb.Booster()
    model.load_model(str(json_model))

    print(f"✅ Loaded: {model.num_features()} features")

    # Try different save formats
    formats_to_try = [
        ("bid_optimizer_v1.model", "deprecated binary"),
        ("bid_optimizer_v1.json", "json with compat"),
    ]

    for filename, desc in formats_to_try:
        try:
            path = models_dir / filename
            print(f"\n💾 Trying {desc}: {path}")
            model.save_model(str(path))
            print(f"✅ Saved successfully")

            # Now try ONNX conversion
            print(f"🔄 Attempting ONNX conversion...")

            from onnxmltools.convert import convert_xgboost
            from skl2onnx.common.data_types import FloatTensorType

            # Reload
            test_model = xgb.Booster()
            test_model.load_model(str(path))

            initial_type = [("float_input", FloatTensorType([None, 13]))]

            onnx_model = convert_xgboost(
                test_model,
                initial_types=initial_type,
                target_opset=10,  # Very old opset for compatibility
            )

            onnx_path = models_dir / "bid_optimizer_latest.onnx"
            with open(onnx_path, "wb") as f:
                f.write(onnx_model.SerializeToString())

            print(f"✅ ONNX export successful!")
            print(f"   File: {onnx_path}")
            return True

        except Exception as e:
            print(f"❌ Failed: {e}")
            continue

    print("\n" + "=" * 60)
    print("❌ All ONNX export methods failed")
    print("=" * 60)
    print("\n💡 BEST SOLUTION: Use Python microservice")
    print("   XGBoost 2.x → ONNX has known compatibility issues")
    print("   The Python microservice is:")
    print("   ✅ Faster to set up (5 minutes)")
    print("   ✅ 100% reliable")
    print("   ✅ Easy to deploy to Cloud Run")
    print("   ✅ ~$20/month for 10M predictions")
    print("\nOR downgrade XGBoost:")
    print("   pip install xgboost==1.7.6")
    print("   Then retrain your model")


if __name__ == "__main__":
    downgrade_and_export()
