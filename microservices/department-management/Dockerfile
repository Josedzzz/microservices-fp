FROM python:3.11-slim

WORKDIR /app

# Copy requirements and install
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Verify installation (for debugging)
RUN python -c "import uvicorn; print('Uvicorn installed successfully')"

# Copy application code
COPY ./app /app/app

# Expose the port
EXPOSE 8082

# Use the python -m syntax which is more reliable
CMD ["python", "-m", "uvicorn", "app.main:app", "--host", "0.0.0.0", "--port", "8082", "--reload"]
