# Use an official Python runtime as a parent image
FROM python:3.8-slim

# Set the working directory in the container
WORKDIR /app/ml-microservices

# Install any needed packages specified in requirements.txt
COPY requirements.txt /app/
RUN pip install --no-cache-dir -r requirements.txt
RUN pip install transformers flask elevenlabs deepgram

# Copy the current directory contents into the container at /app
COPY ml-microservices/processing.py /app/ml-microservices/

# Make port 80 available to the world outside this container
EXPOSE 80

# Run app.py when the container launches
CMD ["python", "processing.py"]
