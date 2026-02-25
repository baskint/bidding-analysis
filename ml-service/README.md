# ML Service — Bid Optimizer

A Flask REST API that serves an XGBoost bid price prediction model.

---

## Setup

### Linux / Mac
```bash
cd ml-service
python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt
```

### Windows
```bash
cd ml-service
py -3.12 -m venv venv
venv\Scripts\activate
pip install -r requirements.txt
```

---

## Training the Model

The training script generates the model files that `ml_service.py` needs.
It supports two modes: **synthetic data** (default) or a **real CSV file**.

### Synthetic data (quick start)
```bash
python train_13_features.py
```

### Real CSV data
```bash
python train_13_features.py --csv sample_training_data.csv
```

The CSV must contain these columns (see `sample_training_data.csv` for an example):

| Column | Type | Notes |
|---|---|---|
| `floor_price` | float | Minimum bid price |
| `engagement_score` | float | 0.0 – 1.0 |
| `conversion_probability` | float | 0.0 – 1.0 |
| `historical_win_rate` | float | 0.0 – 1.0 |
| `historical_avg_bid` | float | |
| `historical_avg_win_price` | float | |
| `device_type` | string | `desktop`, `mobile`, `tablet` |
| `segment_category` | string | `premium`, `standard`, `value`, `new_user` |
| `hour_of_day` | int | 0 – 23 |
| `day_of_week` | int | 0 (Mon) – 6 (Sun) |
| `country` | string | `US`, `GB`, `CA`, `AU`, `DE`, `FR`, `JP`, `CN`, `IN` |
| `campaign_spend_last_7d` | float | |
| `campaign_conversions_last_7d` | float | |
| `optimal_bid` | float | **Target variable** |

Training outputs three files to `../models/`:
- `bid_optimizer_latest.json` — XGBoost model
- `bid_optimizer_latest_encoders.json` — categorical encoders
- `model_info.json` — feature list, metrics, and data source used

---

## Running the Service

```bash
python ml_service.py
```

The service starts on **port 5001**.

---

## API Endpoints

### `GET /health`
```bash
curl http://localhost:5001/health
```
```json
{"model_loaded": true, "status": "healthy"}
```

### `POST /predict`
```bash
curl -X POST http://localhost:5001/predict \
  -H "Content-Type: application/json" \
  -d '{
    "features": {
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
  }'
```
```json
{"model_version": "bid_optimizer_latest", "predicted_bid": 3.42}
```

---

## Running Tests

```bash
python test_service.py
```

The service must be running before executing the tests.

---

## Docker

```bash
docker build -t ml-service .
docker run -p 5001:5001 ml-service
```
