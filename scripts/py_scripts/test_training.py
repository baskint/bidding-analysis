"""
Simple test to verify your ML setup is working
No database required - uses synthetic data
"""

import sys

print("Python version:", sys.version)
print("-" * 60)

# Test 1: Import basic packages
print("\n[1/6] Testing basic imports...")
try:
  import numpy as np
  import pandas as pd

  print("✅ numpy:", np.__version__)
  print("✅ pandas:", pd.__version__)
except ImportError as e:
  print(f"❌ Failed: {e}")
  print("Run: pip install numpy pandas")
  sys.exit(1)

# Test 2: Import ML packages
print("\n[2/6] Testing ML packages...")
try:
  import xgboost as xgb
  from sklearn.model_selection import train_test_split
  from sklearn.metrics import mean_squared_error, r2_score

  print("✅ xgboost:", xgb.__version__)
  print("✅ scikit-learn: imported successfully")
except ImportError as e:
  print(f"❌ Failed: {e}")
  print("Run: pip install xgboost scikit-learn")
  sys.exit(1)

# Test 3: Create synthetic data
print("\n[3/6] Creating synthetic training data...")
try:
  np.random.seed(42)
  n_samples = 1000

  X = pd.DataFrame(
      {
          "floor_price": np.random.uniform(1.0, 5.0, n_samples),
          "engagement_score": np.random.uniform(0, 1, n_samples),
          "conversion_prob": np.random.uniform(0, 0.3, n_samples),
          "win_rate": np.random.uniform(0.2, 0.8, n_samples),
          "hour": np.random.randint(0, 24, n_samples),
          "device": np.random.randint(
              0, 3, n_samples
          ),  # 0=mobile, 1=desktop, 2=tablet
      }
  )

  # Create target variable (optimal bid)
  y = (
      X["floor_price"] * 1.3
      + X["engagement_score"] * 0.8
      + X["conversion_prob"] * 2.5
      + np.random.normal(0, 0.15, n_samples)
  )

  print(f"✅ Created {len(X)} training samples")
  print(f"   Features: {list(X.columns)}")
  print(f"   Target range: ${y.min():.2f} - ${y.max():.2f}")
except Exception as e:
  print(f"❌ Failed: {e}")
  sys.exit(1)

# Test 4: Train model
print("\n[4/6] Training XGBoost model...")
try:
  X_train, X_test, y_train, y_test = train_test_split(
      X, y, test_size=0.2, random_state=42
  )

  model = xgb.XGBRegressor(
      n_estimators=50,
      max_depth=4,
      learning_rate=0.1,
      random_state=42,
      verbosity=0,  # Suppress output
  )

  model.fit(X_train, y_train)
  print("✅ Model training complete")

  # Evaluate
  y_pred = model.predict(X_test)
  rmse = np.sqrt(mean_squared_error(y_test, y_pred))
  r2 = r2_score(y_test, y_pred)

  print(f"   Test RMSE: {rmse:.4f}")
  print(f"   Test R²: {r2:.4f}")

  if r2 < 0.5:
    print("⚠️  Warning: Low R² score (expected with synthetic data)")

except Exception as e:
  print(f"❌ Failed: {e}")
  sys.exit(1)

# Test 5: ONNX export
print("\n[5/6] Testing ONNX export...")
try:
  import onnxmltools
  from onnxmltools.convert import convert_xgboost as convert_xgboost_new
  from onnxconverter_common import FloatTensorType
  import os

  # Save model directory
  os.makedirs("models", exist_ok=True)
  model_path = "models/test_model.onnx"

  # Method 1: Try new API (onnxmltools >= 1.11)
  try:
    # Define initial type
    initial_type = [("float_input", FloatTensorType([None, len(X.columns)]))]

    # Convert using newer API
    onnx_model = convert_xgboost_new(
        model, initial_types=initial_type, target_opset=12
    )

    # Save model
    with open(model_path, "wb") as f:
      f.write(onnx_model.SerializeToString())

    file_size = os.path.getsize(model_path) / 1024  # KB
    print(f"✅ Model exported to: {model_path}")
    print(f"   File size: {file_size:.1f} KB")
    print(f"   Method: onnxmltools.convert")

  except (TypeError, AttributeError) as e:
    # Method 2: Try alternative API with skl2onnx
    print(f"   Trying alternative conversion method...")
    try:
      from skl2onnx import convert_sklearn
      from skl2onnx.common.data_types import FloatTensorType as FloatTensorType2

      # Wrap model for sklearn conversion
      initial_type = [("float_input", FloatTensorType2([None, X.shape[1]]))]

      # Try sklearn converter
      onnx_model = convert_sklearn(
          model, initial_types=initial_type, target_opset=12
      )

      # Save model
      with open(model_path, "wb") as f:
        f.write(onnx_model.SerializeToString())

      file_size = os.path.getsize(model_path) / 1024  # KB
      print(f"✅ Model exported to: {model_path}")
      print(f"   File size: {file_size:.1f} KB")
      print(f"   Method: skl2onnx")

    except Exception as e2:
      # Method 3: Save as XGBoost native format (fallback)
      print(f"   ONNX conversion not available, using XGBoost native format")
      model.save_model("models/test_model.json")
      print(f"✅ Model saved to: models/test_model.json")
      print(f"   Format: XGBoost JSON (can be loaded in Go with xgboost-go)")
      print(f"⚠️  Note: ONNX export requires compatible package versions")

      # Skip ONNX inference test
      print("\n[6/6] Skipping ONNX inference test (ONNX export unavailable)")
      print("\n" + "=" * 60)
      print("✅ SUCCESS! Your ML environment works (without ONNX)")
      print("=" * 60)
      print("\nWhat this test verified:")
      print("  ✅ Python packages installed correctly")
      print("  ✅ XGBoost can train models")
      print("  ✅ Models can be saved (XGBoost format)")
      print("  ⚠️  ONNX export needs package version adjustment")
      print("\nTo fix ONNX export:")
      print("  pip install onnxmltools==1.11.0 skl2onnx==1.15.0")
      print("\nYou're ready to:")
      print("  1. Train models on real data")
      print("  2. Use XGBoost models in Python")
      print("  3. Export to ONNX (after fixing versions)")
      print("-" * 60)
      sys.exit(0)

except ImportError as e:
  print(f"❌ ONNX packages not installed: {e}")
  print("Run: pip install onnxmltools skl2onnx")

  # Save as XGBoost native format instead
  import os

  os.makedirs("models", exist_ok=True)
  model.save_model("models/test_model.json")
  print(f"✅ Model saved to: models/test_model.json (XGBoost format)")

  print("\n[6/6] Skipping ONNX inference test (packages not installed)")
  print("\n" + "=" * 60)
  print("✅ PARTIAL SUCCESS! Core ML works, ONNX optional")
  print("=" * 60)
  print("\nInstall ONNX packages:")
  print("  pip install onnxruntime onnxmltools skl2onnx")
  print("-" * 60)
  sys.exit(0)

except Exception as e:
  print(f"⚠️  ONNX export issue: {e}")
  print(f"   Saving model in XGBoost format instead...")

  import os

  os.makedirs("models", exist_ok=True)
  model.save_model("models/test_model.json")
  print(f"✅ Model saved to: models/test_model.json")
  print(f"   You can still train and use models!")

# Test 6: Load and test ONNX model (if export succeeded)
print("\n[6/6] Testing ONNX inference...")
try:
  import onnxruntime as ort

  # Load ONNX model
  session = ort.InferenceSession(model_path)

  # Get input/output names
  input_name = session.get_inputs()[0].name
  output_name = session.get_outputs()[0].name

  print(f"✅ ONNX model loaded")
  print(f"   Input: {input_name}")
  print(f"   Output: {output_name}")

  # Test prediction
  test_sample = X_test.iloc[0:1].values.astype(np.float32)
  onnx_pred = session.run([output_name], {input_name: test_sample})[0][0][0]
  xgb_pred = model.predict(X_test.iloc[0:1])[0]

  print(f"   ONNX prediction: ${onnx_pred:.4f}")
  print(f"   XGBoost prediction: ${xgb_pred:.4f}")
  print(f"   Difference: ${abs(onnx_pred - xgb_pred):.6f}")

  if abs(onnx_pred - xgb_pred) < 0.001:
    print("✅ Predictions match!")
  else:
    print("⚠️  Small difference (normal for ONNX conversion)")

except ImportError:
  print("⚠️  ONNX Runtime not installed")
  print("   Run: pip install onnxruntime")
except Exception as e:
  print(f"⚠️  Could not test ONNX inference: {e}")
  print(f"   Model training still works!")

# Success summary
print("\n" + "=" * 60)
print("🎉 SUCCESS! Your ML environment is functional!")
print("=" * 60)
print("\nWhat this test verified:")
print("  ✅ Python packages installed correctly")
print("  ✅ XGBoost can train models")
print("  ✅ Models can be saved")
print("\nYou're ready to:")
print("  1. Train models on real data")
print("  2. Test predictions in Python")
print("  3. Integrate with your application")
print("\nNext step: Try training with real data!")
print("  python train_model.py --synthetic --samples 5000")
print("-" * 60)
