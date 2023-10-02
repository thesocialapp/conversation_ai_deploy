# Use an official Python runtime as a parent image
FROM python:3.8-slim

# Set the working directory in the container to /app
WORKDIR /app

# Add the current directory contents into the container at /app
ADD . /app

# Install any needed packages specified in requirements.txt
RUN pip install --no-cache-dir -r requirements.txt

# Copy the Go.Sum file into the Docker image
COPY go.sum /app

# Run go mod download to download the required dependencies
RUN go mod download

# Build the Go application
RUN go build -o main .

# Make port 80 available to the world outside this container
EXPOSE 80

# Define environment variable
ENV NAME World

# Set the entrypoint to run the built Go application
CMD ["./main"]
