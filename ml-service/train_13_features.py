#!/usr/bin/env python3
"""
Bid Optimizer Model Training - 13 Features Version
Creates an XGBoost model with the exact feature count that ml_service.py expects

Usage:
  python train_13_features.py                             # synthetic data
  python train_13_features.py --csv sample_training_data.csv  # real CSV data
"""

import argparse
import json
import numpy as np
import xgboost as xgb
from sklearn.model_selection import train_test_split
from sklearn.metrics import mean_squared_error, r2_score

# 13 FEATURES to match ml_service.py:
# 0: floor_price
# 1: engagement_score
# 2: conversion_probability
# 3: historical_win_rate
# 4: historical_avg_bid
# 5: historical_avg_win_price
# 6: device_type_encoded
# 7: segment_category_encoded
# 8: hour_of_day
# 9: day_of_week
# 10: country_encoded
# 11: campaign_spend_last_7d
# 12: campaign_conversions_last_7d

FEATURE_COLS = [
    'floor_price', 'engagement_score', 'conversion_probability',
    'historical_win_rate', 'historical_avg_bid', 'historical_avg_win_price',
    'device_type', 'segment_category', 'hour_of_day', 'day_of_week',
    'country', 'campaign_spend_last_7d', 'campaign_conversions_last_7d'
]

encoders = {
    "device_type": {"desktop": 0, "mobile": 1, "tablet": 2},
    "segment_category": {"premium": 0, "standard": 1, "value": 2, "new_user": 3},
    "country": {"US": 0, "GB": 1, "CA": 2, "AU": 3, "DE": 4, "FR": 5, "JP": 6, "CN": 7, "IN": 8}
}


def load_csv(path):
    import pandas as pd
    print(f"Loading data from {path}...")
    df = pd.read_csv(path)

    missing = [c for c in FEATURE_COLS + ['optimal_bid'] if c not in df.columns]
    if missing:
        raise ValueError(f"CSV is missing required columns: {missing}")

    for col, mapping in encoders.items():
        unknown = set(df[col].unique()) - set(mapping.keys())
        if unknown:
            raise ValueError(f"Unknown values in '{col}': {unknown}. Expected: {list(mapping.keys())}")
        df[col] = df[col].map(mapping)

    X = df[FEATURE_COLS].values
    y = df['optimal_bid'].values
    print(f"Loaded {len(df)} rows from CSV.")
    return X, y


def load_synthetic():
    print("No CSV provided — generating synthetic training data...")
    np.random.seed(42)
    n_samples = 1000
    X = np.random.rand(n_samples, 13)
    y = 0.5 + 4.5 * (X[:, 0] * 0.5 + X[:, 1] * 0.3 + X[:, 2] * 0.2) + np.random.normal(0, 0.1, n_samples)
    print(f"Generated {n_samples} synthetic rows.")
    return X, y


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='Train the bid optimizer XGBoost model.')
    parser.add_argument('--csv', type=str, default=None, metavar='FILE',
                        help='Path to a CSV file with real training data (optional)')
    args = parser.parse_args()

    X, y = load_csv(args.csv) if args.csv else load_synthetic()

    X_train, X_val, y_train, y_val = train_test_split(X, y, test_size=0.2, random_state=42)

    model = xgb.XGBRegressor(
        n_estimators=100,
        max_depth=5,
        learning_rate=0.1,
        objective='reg:squarederror',
        random_state=42
    )

    print("Training XGBoost model with 13 features...")
    model.fit(X_train, y_train)

    train_pred = model.predict(X_train)
    val_pred = model.predict(X_val)

    print(f"\n=== Training Results ===")
    print(f"Train RMSE: {np.sqrt(mean_squared_error(y_train, train_pred)):.4f}")
    print(f"Train R²:   {r2_score(y_train, train_pred):.4f}")
    print(f"Val RMSE:   {np.sqrt(mean_squared_error(y_val, val_pred)):.4f}")
    print(f"Val R²:     {r2_score(y_val, val_pred):.4f}")

    import os
    os.makedirs('../models', exist_ok=True)

    model.save_model('../models/bid_optimizer_latest.json')
    print("\n✅ Model saved to ../models/bid_optimizer_latest.json")

    with open('../models/bid_optimizer_latest_encoders.json', 'w') as f:
        json.dump(encoders, f, indent=2)
    print("✅ Encoders saved to ../models/bid_optimizer_latest_encoders.json")

    model_info = {
        'version': '1.0.0',
        'data_source': args.csv if args.csv else 'synthetic',
        'num_features': 13,
        'features': FEATURE_COLS,
        'metrics': {
            'train_r2': float(r2_score(y_train, train_pred)),
            'val_r2': float(r2_score(y_val, val_pred)),
            'train_rmse': float(np.sqrt(mean_squared_error(y_train, train_pred))),
            'val_rmse': float(np.sqrt(mean_squared_error(y_val, val_pred)))
        }
    }

    with open('../models/model_info.json', 'w') as f:
        json.dump(model_info, f, indent=2)
    print("✅ Model info saved to ../models/model_info.json")
    print("\n🚀 All files created! You can now run: python ml_service.py")
