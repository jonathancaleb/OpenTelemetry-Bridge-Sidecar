# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY main.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o sample-app .

# Runtime stage
FROM alpine:3.19

WORKDIR /app
COPY --from=builder /app/sample-app .

EXPOSE 3000

CMD ["./sample-app"]
