# Multi-stage build for production
FROM golang:1.23-alpine AS builder

# Install git for go modules
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod ./
COPY go.sum* ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o discord-bot ./cmd/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create a non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Set working directory
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/discord-bot .

# Change ownership
RUN chown appuser:appgroup discord-bot

# Switch to non-root user
USER appuser

# Install wget for health checks
USER root
RUN apk add --no-cache wget
USER appuser

# Health check
HEALTHCHECK --interval=300s --timeout=10s --start-period=120s --retries=3 \
  CMD wget --quiet --tries=1 --spider http://localhost:${BOT_PORT:-8082}/health || exit 1

# Run the application
CMD ["./discord-bot"]