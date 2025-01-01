# Use Golang base image
FROM golang:1.20 AS builder

# Set the working directory
WORKDIR /app

# Copy the Go modules and other necessary files first
COPY go.mod go.sum ./

# Install dependencies
RUN go mod download

# Copy the Go source code
COPY . .

# Build the Go application
RUN --mount=type=cache,id=go-build-cache,target=/root/.cache/go-build go build -o out ./cmd/api

# Final image setup (optional)
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/out /root/
