# =========================
# 1. Build stage
# =========================
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Install CA certs (needed for HTTPS calls if any)
RUN apk add --no-cache ca-certificates

# Copy go mod files first (better caching)
COPY go.mod go.sum ./

RUN go mod download

# Copy source code
COPY . .

# Build binary (static, production-ready)
ARG VERSION=dev
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags "-X library-api/internal/version.Version=${VERSION}" \
    -o server ./cmd/api

# =========================
# 2. Runtime stage
# =========================
FROM alpine:3.20

WORKDIR /app

# CA certificates for HTTPS
RUN apk add --no-cache ca-certificates

# Run as a non-root user
RUN addgroup -S app && adduser -S -G app app

# Copy binary from builder
COPY --from=builder /app/server .

# Copy HTML templates
COPY --from=builder /app/internal/views ./internal/views

# Copy .env if you want local container config (optional)
# COPY .env .

RUN chown -R app:app /app
USER app

# Production defaults (override via env vars / Cloud Run / Secret Manager)
ENV GIN_MODE=release

# Expose port
EXPOSE 8080

# Run app
CMD ["./server"]
