# We initialize Golang first and use the alpine version
# We set the working directory to /app and copy all files inside
# We build the binary and name it main
FROM golang:1.19-alpine3.16 AS builder

WORKDIR /app
# Copy all files from the current directory to the working directory
# including the go.mod and go.sum files
COPY . .
# We build the binary from the backend folder
RUN go build -o main main.go

FROM alpine:3.16
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .

# # Path: Dockerfile
# # We initialize a new container from scratch
# # We copy the binary from the builder container to the new container
# # We set the entrypoint to the binary
# FROM alpine:3.16

# # We set up python and build it
# FROM python:3.9-alpine3.14 as pybuilder

# # We set the working directory to /internal/py and copy all files inside
# WORKDIR /internal/py
# COPY ./internal/py .
# # Copy the requirements.txt file from the base of the project to the working
# # directory
# COPY requirements.txt .

# # Install Python
# RUN apk add --no-cache python3

# # We install the dependencies from the requirements.txt file at the base of the 
# # project
# RUN pip install -r requirements.txt

# # We set up the final container
# FROM alpine:3.14

# WORKDIR /app

# # We copy the binary from the builder container to the new container
# COPY --from=pybuilder /internal/py /app/internal/py
# COPY --from=builder /app/main /app/main

# RUN chmod +x /app/internal/py/main.py

EXPOSE 4400
CMD ["/app/main"]


