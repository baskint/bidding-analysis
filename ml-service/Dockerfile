FROM python:3.11-slim

WORKDIR /app

# Copy requirements and install
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy service code
COPY ml_service.py .

# Copy model files (now in same directory)
COPY bid_optimizer_latest.json models/bid_optimizer_latest.json
COPY bid_optimizer_latest_encoders.json models/bid_optimizer_latest_encoders.json

# Cloud Run sets PORT env variable
ENV PORT=8080

# Use gunicorn for production
CMD exec gunicorn --bind :$PORT --workers 1 --threads 8 --timeout 60 ml_service:app