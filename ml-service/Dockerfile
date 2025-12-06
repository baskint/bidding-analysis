FROM python:3.11-slim

WORKDIR /app

# Copy requirements and install
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy service code
COPY ml_service.py .

# Copy models from parent directory
COPY ../models/bid_optimizer_latest.json models/
COPY ../models/bid_optimizer_latest_encoders.json models/

# Cloud Run uses PORT environment variable
ENV PORT=8080

# Use gunicorn for production
CMD exec gunicorn --bind :$PORT --workers 1 --threads 8 --timeout 60 ml_service:app