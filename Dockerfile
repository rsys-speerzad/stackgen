# Start from the official Golang image for building
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app
RUN go build -o server ./cmd/server/main.go

# Use a minimal image for running
FROM alpine:latest

WORKDIR /app

# Copy the built binary from builder
COPY --from=builder /app/server .

# Expose port (change if your app uses a different port)
EXPOSE 8080

# Run the binary
CMD ["./server"]