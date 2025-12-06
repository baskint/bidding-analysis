#!/usr/bin/env python3
"""
Quick script to convert your existing JSON model to binary format for Go/leaves
Run this after training to create the binary version
"""

import xgboost as xgb
import json
from pathlib import Path


def convert_json_to_binary():
    """Convert existing JSON model to binary format"""

    models_dir = Path("models")

    # Load the JSON model
    json_model = models_dir / "bid_optimizer_latest.json"

    if not json_model.exists():
        print(f"❌ Model not found: {json_model}")
        print("Please run train_model.py first")
        return

    print(f"📖 Loading JSON model from {json_model}")

    # Load the model
    model = xgb.Booster()
    model.load_model(json_model)

    # Save in binary format
    binary_model = models_dir / "bid_optimizer_latest.bin"
    print(f"💾 Saving binary model to {binary_model}")
    model.save_model(binary_model)

    # Also save as .ubj (universal binary JSON - another format leaves supports)
    ubj_model = models_dir / "bid_optimizer_latest.ubj"
    print(f"💾 Saving UBJ model to {ubj_model}")
    model.save_model(ubj_model)

    print()
    print("✅ Conversion complete!")
    print()
    print("Binary formats created:")
    print(f"  - {binary_model} (XGBoost binary)")
    print(f"  - {ubj_model} (Universal Binary JSON)")
    print()
    print("Now update your Go code to use:")
    print(f'  modelPath: "models/bid_optimizer_latest.bin"')
    print()


if __name__ == "__main__":
    convert_json_to_binary()
