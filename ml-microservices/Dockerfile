# Dockerfile
FROM python:3.8-slim

# Install dependencies
RUN pip install transformers flask

# Copy script
COPY processing.py /app/

# Set work directory
WORKDIR /app

# Command to run the script
CMD ["python", "processing.py"]
