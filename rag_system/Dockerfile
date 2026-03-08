# -----------------------------
# Stage 1: Build the Go binary
# -----------------------------
FROM golang:1.25-alpine AS builder

# Install git (needed for Go modules)
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum first (leverage Docker cache)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of your source code
COPY . .

# Build a statically linked binary for Linux
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o rag_system main.go

# -----------------------------
# Stage 2: Minimal runtime image
# -----------------------------
FROM alpine:latest

# Add CA certificates (needed if your Go app makes HTTPS requests)
RUN apk add --no-cache ca-certificates

# Set working directory
WORKDIR /app

# Copy the built binary from the builder
COPY --from=builder /app/rag_system .

# Copy the .env file for environment variables
COPY --from=builder /app/.env .env

# Expose the port used by your RAG server (default 8080)
EXPOSE 8080

# Start the server
CMD ["./rag_system"]