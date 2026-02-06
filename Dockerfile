# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o sidecar ./cmd/sidecar

# Runtime stage
FROM alpine:3.19

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/sidecar .

# Copy default config
COPY internal/config/config.yaml /app/config.yaml

EXPOSE 8080

CMD ["./sidecar"]
