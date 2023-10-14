# Build stage: Use a Golang-based image to build the binary
FROM golang:1.19-buster AS builder

WORKDIR /app

# Copy your Go source code into the container
COPY . .

# Install Opus and any other dependencies specific to your Alpine-based image
RUN apt-get update && apt-get -y install libopus-dev libopusfile-dev

# Build the Go binary
RUN go build -o main main.go

# Final stage: Use an Alpine Linux-based image for the lightweight runtime image
FROM alpine:3.16

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main .
COPY app.env .


EXPOSE 4400
CMD ["/app/main"]
