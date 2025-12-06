# Dockerfile for ML Prediction Service
FROM python:3.11-slim

WORKDIR /app

# Install dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy service code
COPY ml_service.py .

# Copy models (these will be in the parent directory during build)
COPY ../models/bid_optimizer_latest.json models/
COPY ../models/bid_optimizer_latest_encoders.json models/

# Expose port
EXPOSE 8080

# Use gunicorn for production
CMD exec gunicorn --bind :$PORT --workers 1 --threads 8 --timeout 0 ml_service:app
