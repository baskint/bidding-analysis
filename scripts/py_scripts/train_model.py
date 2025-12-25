#!/usr/bin/env python3
"""
Complete Training Pipeline for Bid Optimization
Uses XGBoost with ONNX export for Go inference
"""

import os
import sys
import json
import argparse
import uuid
import shutil
from datetime import datetime, timedelta
from typing import Dict, List, Tuple

import numpy as np
import pandas as pd
import psycopg2
import xgboost as xgb
from sklearn.model_selection import train_test_split, cross_val_score
from sklearn.metrics import mean_squared_error, mean_absolute_error, r2_score
import onnxmltools
from onnxconverter_common import FloatTensorType
import yaml


class BidOptimizer:
  """
  Bid Optimization Model Trainer

  Features:
  - Automatic feature engineering
  - Hyperparameter tuning
  - Model validation
  - ONNX export for Go inference
  """

  FEATURE_COLUMNS = [
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
  ]

  def __init__(self, config_path: str = "config.yaml"):
    """Initialize with configuration"""
    with open(config_path) as f:
      self.config = yaml.safe_load(f)

    self.model = None
    self.feature_encoders = {}
    self.scaler = None

  def connect_db(self):
    """Connect to PostgreSQL database"""
    return psycopg2.connect(
        host=self.config["database"]["host"],
        port=self.config["database"]["port"],
        database=self.config["database"]["name"],
        user=self.config["database"]["user"],
        password=self.config["database"]["password"],
    )

  def generate_synthetic_data(self, n_samples: int = 5000) -> pd.DataFrame:
    """
    Generate synthetic training data (no database required!)

    Args:
        n_samples: Number of samples to generate

    Returns:
        DataFrame with synthetic training data
    """
    np.random.seed(42)

    device_types = ["mobile", "desktop", "tablet"]
    categories = ["premium", "standard", "value", "new_user"]
    countries = ["US", "GB", "CA", "AU", "DE", "FR", "JP"]

    data = {
        "campaign_id": [str(uuid.uuid4()) for _ in range(n_samples)],
        "floor_price": np.random.uniform(1.0, 5.0, n_samples),
        # Skewed distribution
        "engagement_score": np.random.beta(2, 5, n_samples),
        "conversion_probability": np.random.beta(
            2, 8, n_samples
        ),  # Lower conversion rates
        "device_type": np.random.choice(device_types, n_samples),
        "segment_category": np.random.choice(categories, n_samples),
        "country": np.random.choice(
            countries, n_samples, p=[0.4, 0.2, 0.15, 0.1, 0.05, 0.05, 0.05]
        ),
        "created_at": pd.date_range(
            end=datetime.now(), periods=n_samples, freq="5min"
        ),
    }

    df = pd.DataFrame(data)

    # Add temporal features
    df["hour_of_day"] = pd.to_datetime(df["created_at"]).dt.hour
    df["day_of_week"] = pd.to_datetime(df["created_at"]).dt.dayofweek

    # Historical features (simulated)
    df["historical_win_rate"] = np.random.beta(4, 6, n_samples)  # Around 0.4
    df["historical_avg_bid"] = df["floor_price"] * np.random.uniform(
        1.2, 2.0, n_samples
    )
    df["historical_avg_win_price"] = df["historical_avg_bid"] * np.random.uniform(
        0.8, 1.2, n_samples
    )
    df["campaign_spend_last_7d"] = np.random.uniform(100, 10000, n_samples)
    df["campaign_conversions_last_7d"] = np.random.poisson(5, n_samples)

    # Generate realistic optimal_bid (target variable)
    # Formula considers multiple factors
    df["optimal_bid"] = (
        df["floor_price"] * 1.3  # Base markup
        + df["engagement_score"] * 0.8  # User engagement value
        + df["conversion_probability"] * 3.0  # Conversion value
        + (df["historical_win_rate"] - 0.4) * 1.5  # Historical performance
        + np.random.normal(0, 0.2, n_samples)  # Some noise
    )

    # Ensure optimal_bid is reasonable
    df["optimal_bid"] = df["optimal_bid"].clip(
        lower=df["floor_price"] * 1.01,  # At least 1% above floor
        upper=df["floor_price"] * 5.0,  # Max 5x floor price
    )

    # Training set split
    train_ratio = 0.7
    val_ratio = 0.15
    df["training_set"] = np.random.choice(
        ["train", "validation", "test"],
        n_samples,
        p=[train_ratio, val_ratio, 1 - train_ratio - val_ratio],
    )

    print(f"✅ Generated {n_samples} synthetic samples")
    print(
        f"   Floor price range: ${df['floor_price'].min():.2f} - ${df['floor_price'].max():.2f}"
    )
    print(
        f"   Optimal bid range: ${df['optimal_bid'].min():.2f} - ${df['optimal_bid'].max():.2f}"
    )
    print(f"   Avg conversion prob: {df['conversion_probability'].mean():.2%}")

    return df

  def load_training_data(self, days: int = 30) -> pd.DataFrame:
    """
    Load training data from PostgreSQL

    Args:
        days: Number of days of historical data to load

    Returns:
        DataFrame with training examples
    """
    try:
      query = f"""
            SELECT 
                btd.*,
                c.budget as campaign_budget,
                c.status as campaign_status
            FROM bid_training_data btd
            JOIN campaigns c ON btd.campaign_id = c.id
            WHERE btd.created_at >= NOW() - INTERVAL '{days} days'
              AND btd.optimal_bid > 0
              AND btd.optimal_bid < 100  -- Remove outliers
            ORDER BY btd.created_at DESC
            """

      conn = self.connect_db()
      df = pd.read_sql(query, conn)
      conn.close()

      print(f"Loaded {len(df)} training examples from database")
      return df
    except Exception as e:
      print(f"⚠️  Could not load data from database: {e}")
      print(f"   Use --synthetic flag to generate synthetic data")
      raise

  def engineer_features(self, df: pd.DataFrame) -> pd.DataFrame:
    """
    Feature engineering pipeline

    Args:
        df: Raw dataframe

    Returns:
        DataFrame with engineered features
    """
    print("Engineering features...")

    # Time-based features
    df["hour_of_day"] = pd.to_datetime(df["created_at"]).dt.hour
    df["day_of_week"] = pd.to_datetime(df["created_at"]).dt.dayofweek
    df["is_weekend"] = df["day_of_week"].isin([5, 6]).astype(int)
    df["is_business_hours"] = df["hour_of_day"].between(9, 17).astype(int)

    # Categorical encoding with frequency
    for col in ["device_type", "segment_category", "country"]:
      encoded_col = f"{col}_encoded"

      # Use frequency encoding
      freq_encoding = df[col].value_counts().to_dict()
      df[encoded_col] = df[col].map(freq_encoding).fillna(0)

      # Save encoder for inference
      self.feature_encoders[col] = freq_encoding

    # Interaction features
    df["engagement_x_conversion"] = (
        df["engagement_score"] * df["conversion_probability"]
    )
    df["floor_price_ratio"] = df["floor_price"] / \
        (df["historical_avg_bid"] + 0.01)

    # Handle missing values
    df["engagement_score"] = df["engagement_score"].fillna(0.5)
    df["conversion_probability"] = df["conversion_probability"].fillna(0.05)
    df["historical_win_rate"] = df["historical_win_rate"].fillna(0.3)
    df["historical_avg_bid"] = df["historical_avg_bid"].fillna(
        df["floor_price"])
    df["historical_avg_win_price"] = df["historical_avg_win_price"].fillna(
        df["floor_price"]
    )

    return df

  def prepare_dataset(
      self, df: pd.DataFrame, test_size: float = 0.2
  ) -> Tuple[pd.DataFrame, pd.Series, pd.DataFrame, pd.Series]:
    """
    Prepare train/val split

    Args:
        df: Engineered dataframe
        test_size: Validation set proportion

    Returns:
        X_train, y_train, X_val, y_val
    """
    # Select features
    X = df[self.FEATURE_COLUMNS].copy()
    y = df["optimal_bid"].copy()

    # Remove any rows with NaN
    mask = ~(X.isna().any(axis=1) | y.isna())
    X = X[mask]
    y = y[mask]

    print(f"Dataset size after cleaning: {len(X)}")
    print(f"Feature columns: {self.FEATURE_COLUMNS}")

    # Split
    X_train, X_val, y_train, y_val = train_test_split(
        X, y, test_size=test_size, random_state=42, shuffle=True
    )

    print(f"Train set: {len(X_train)} samples")
    print(f"Validation set: {len(X_val)} samples")

    return X_train, y_train, X_val, y_val

  def train(
      self,
      X_train: pd.DataFrame,
      y_train: pd.Series,
      X_val: pd.DataFrame,
      y_val: pd.Series,
      params: Dict = None,
  ) -> Dict:
    """
    Train XGBoost model

    Args:
        X_train, y_train: Training data
        X_val, y_val: Validation data
        params: Optional hyperparameters

    Returns:
        Dictionary of training metrics
    """
    print("Training XGBoost model...")

    # Default hyperparameters
    if params is None:
      params = {
          "n_estimators": 300,
          "max_depth": 8,
          "learning_rate": 0.05,
          "subsample": 0.8,
          "colsample_bytree": 0.8,
          "min_child_weight": 3,
          "gamma": 0.1,
          "reg_alpha": 0.1,
          "reg_lambda": 1.0,
          "objective": "reg:squarederror",
          "n_jobs": -1,
          "random_state": 42,
      }

    # Extract early_stopping_rounds if present (not a model parameter)
    early_stopping_rounds = params.pop("early_stopping_rounds", 20)

    # Remove any other non-model parameters
    params.pop("eval_metric", None)

    self.model = xgb.XGBRegressor(**params)

    # Train with early stopping
    # Try the simpler approach that works with all versions
    try:
      self.model.fit(
          X_train,
          y_train,
          eval_set=[(X_val, y_val)],
          verbose=False,  # Reduce output
      )
    except Exception as e:
      print(f"Training with basic configuration")
      self.model.fit(X_train, y_train)

    # Evaluate
    train_pred = self.model.predict(X_train)
    val_pred = self.model.predict(X_val)

    metrics = {
        "train_rmse": np.sqrt(mean_squared_error(y_train, train_pred)),
        "val_rmse": np.sqrt(mean_squared_error(y_val, val_pred)),
        "train_mae": mean_absolute_error(y_train, train_pred),
        "val_mae": mean_absolute_error(y_val, val_pred),
        "train_r2": r2_score(y_train, train_pred),
        "val_r2": r2_score(y_val, val_pred),
        "feature_importance": dict(
            zip(self.FEATURE_COLUMNS, self.model.feature_importances_.tolist())
        ),
    }

    print("\n=== Training Results ===")
    print(f"Train RMSE: {metrics['train_rmse']:.4f}")
    print(f"Val RMSE: {metrics['val_rmse']:.4f}")
    print(f"Train R²: {metrics['train_r2']:.4f}")
    print(f"Val R²: {metrics['val_r2']:.4f}")

    print("\n=== Top 5 Important Features ===")
    sorted_features = sorted(
        metrics["feature_importance"].items(), key=lambda x: x[1], reverse=True
    )
    for feature, importance in sorted_features[:5]:
      print(f"{feature}: {importance:.4f}")

    return metrics

  def export_to_onnx(self, output_path: str):
    """
    Export trained model to ONNX format for Go inference
    Fallback to XGBoost native format if ONNX fails

    Args:
        output_path: Path to save ONNX model
    """
    print(f"\nExporting model: {output_path}")

    # Create output directory
    os.makedirs(
        os.path.dirname(output_path) if os.path.dirname(
            output_path) else "models",
        exist_ok=True,
    )

    # Try ONNX export
    try:
      from onnxmltools.convert import convert_xgboost as convert_xgb_new

      # Create ONNX model
      initial_types = [
          ("float_input", FloatTensorType([None, len(self.FEATURE_COLUMNS)]))
      ]

      onnx_model = convert_xgb_new(
          self.model, initial_types=initial_types, target_opset=12
      )

      # Save ONNX model
      with open(output_path, "wb") as f:
        f.write(onnx_model.SerializeToString())

      print(f"✅ Model exported to ONNX: {output_path}")

    except (TypeError, AttributeError, ValueError) as e:
      # ONNX export failed, use XGBoost native format
      print(f"⚠️  ONNX export failed: {e}")
      print(f"   Saving as XGBoost native format instead...")

      # Change extension to .json
      xgb_path = output_path.replace(".onnx", ".json")
      self.model.save_model(xgb_path)

      print(f"✅ Model saved to XGBoost format: {xgb_path}")
      print(f"   Note: For Go integration, you'll need to:")
      print(f"   1. Fix ONNX package versions, OR")
      print(f"   2. Use xgboost-go library")

      # Update output_path for encoder saving
      output_path = xgb_path

    # Save feature encoders for Go
    encoder_path = output_path.replace(".onnx", "_encoders.json").replace(
        ".json", "_encoders.json"
    )
    with open(encoder_path, "w") as f:
      json.dump(self.feature_encoders, f, indent=2)

    print(f"✅ Feature encoders saved to: {encoder_path}")

  def save_metadata(self, metrics: Dict, output_path: str):
    """Save model metadata to database (if available)"""
    # Check if database is enabled
    if not self.config.get("database", {}).get("enabled", True):
      print("Database disabled, skipping metadata save")
      return

    try:
      conn = self.connect_db()
      cursor = conn.cursor()

      query = """
            INSERT INTO ml_model_metadata (
                model_path,
                model_type,
                train_rmse,
                val_rmse,
                train_r2,
                val_r2,
                feature_importance,
                created_at
            ) VALUES (%s, %s, %s, %s, %s, %s, %s, NOW())
            """

      # Convert numpy types to Python native types
      cursor.execute(
          query,
          (
              output_path,
              "xgboost_bid_optimizer",
              # Convert np.float64 to Python float
              float(metrics["train_rmse"]),
              float(metrics["val_rmse"]),
              float(metrics["train_r2"]),
              float(metrics["val_r2"]),
              json.dumps(metrics["feature_importance"]),
          ),
      )

      conn.commit()
      conn.close()
      print(f"✅ Metadata saved to database")
    except Exception as e:
      print(f"⚠️  Could not save metadata to database: {e}")
      print(f"   Continuing without database logging...")


def main():
  parser = argparse.ArgumentParser(description="Train bid optimization model")
  parser.add_argument("--config", default="config.yaml",
                      help="Config file path")
  parser.add_argument("--days", type=int, default=30,
                      help="Days of training data")
  parser.add_argument(
      "--output", default="models/bid_optimizer.onnx", help="Output model path"
  )
  parser.add_argument(
      "--synthetic",
      action="store_true",
      help="Use synthetic data (no database needed)",
  )
  parser.add_argument(
      "--samples", type=int, default=5000, help="Number of synthetic samples"
  )

  args = parser.parse_args()

  # Initialize trainer
  optimizer = BidOptimizer(args.config)

  # Load and prepare data
  print("Loading training data...")
  if args.synthetic:
    print(f"📊 Generating {args.samples} synthetic samples...")
    df = optimizer.generate_synthetic_data(args.samples)
  else:
    df = optimizer.load_training_data(days=args.days)

  df = optimizer.engineer_features(df)
  X_train, y_train, X_val, y_val = optimizer.prepare_dataset(df)

  # Train model
  metrics = optimizer.train(X_train, y_train, X_val, y_val)

  # Export model
  timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
  model_path = args.output.replace(".onnx", f"_{timestamp}.onnx")

  # Save model as JSON (ONNX has compatibility issues)
  os.makedirs(os.path.dirname(model_path) if os.path.dirname(
      model_path) else "models", exist_ok=True)
  model_json = model_path.replace(".onnx", ".json")
  optimizer.model.save_model(model_json)
  print(f"✅ Model saved: {model_json}")

  # Save encoders
  encoder_path = model_json.replace(".json", "_encoders.json")
  with open(encoder_path, 'w') as f:
    json.dump(optimizer.feature_encoders, f, indent=2)
  print(f"✅ Encoders saved: {encoder_path}")

  actual_model_path = model_json

  # Determine actual saved model path (might be .json if ONNX failed)
  actual_model_path = model_path
  if not os.path.exists(model_path):
    # ONNX failed, model was saved as JSON
    actual_model_path = model_path.replace(".onnx", ".json")

  optimizer.save_metadata(metrics, actual_model_path)

  print(f"\n✅ Training complete! Model saved to: {actual_model_path}")

  # Create copy to "latest" (Windows-compatible)
  # Determine the correct extension for latest
  if actual_model_path.endswith(".json"):
    latest_path = args.output.replace(".onnx", "_latest.json")
  else:
    latest_path = args.output.replace(".onnx", "_latest.onnx")

  # Remove old latest if it exists
  if os.path.exists(latest_path):
    os.remove(latest_path)

  # Copy the actual model file
  shutil.copy2(actual_model_path, latest_path)
  print(f"✅ Latest model: {latest_path}")

  # Also copy encoders to latest
  encoder_src = actual_model_path.replace(".onnx", "_encoders.json").replace(
      ".json", "_encoders.json"
  )
  encoder_dst = latest_path.replace(".onnx", "_encoders.json").replace(
      ".json", "_encoders.json"
  )

  if os.path.exists(encoder_src):
    if os.path.exists(encoder_dst):
      os.remove(encoder_dst)
    shutil.copy2(encoder_src, encoder_dst)
    print(f"✅ Latest encoders: {encoder_dst}")

  print("\n" + "=" * 60)
  print("🎉 Training Pipeline Complete!")
  print("=" * 60)
  print(f"\nYour trained model is ready:")
  print(f"  Model: {actual_model_path}")
  print(f"  Latest: {latest_path}")
  print(f"\nModel Performance:")
  print(f"  Train R²: {metrics['train_r2']:.4f}")
  print(f"  Val R²: {metrics['val_r2']:.4f}")
  print(f"  Val RMSE: ${metrics['val_rmse']:.4f}")
  print(f"\nNext steps:")
  print(f"  1. Test predictions: python test_predictions.py")
  print(f"  2. Integrate with your Go service")
  print(f"  3. Deploy to production")
  print("=" * 60)


if __name__ == "__main__":
  main()
