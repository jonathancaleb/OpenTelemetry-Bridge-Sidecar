# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod files first for layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary with version info
ARG VERSION=dev
ARG COMMIT=unknown
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w -X main.version=${VERSION} -X main.commit=${COMMIT}" \
    -o sidecar ./cmd/sidecar

# Runtime stage
FROM alpine:3.19

# OCI standard labels
LABEL org.opencontainers.image.title="OpenTelemetry Bridge Sidecar" \
      org.opencontainers.image.description="Transparent reverse-proxy that adds W3C tracing to any HTTP service" \
      org.opencontainers.image.source="https://github.com/jonathancaleb/OpenTelemetry-Bridge-Sidecar-" \
      org.opencontainers.image.vendor="jonathancaleb"

# Install ca-certificates for TLS and curl for healthcheck
RUN apk add --no-cache ca-certificates curl

# Create non-root user
RUN addgroup -S sidecar && adduser -S sidecar -G sidecar

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/sidecar .

# Copy default config
COPY internal/config/config.yaml /app/config.yaml

# Set ownership
RUN chown -R sidecar:sidecar /app

USER sidecar

EXPOSE 8080

HEALTHCHECK --interval=15s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/ || exit 1

CMD ["./sidecar"]
