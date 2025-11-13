FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy shared module first
COPY shared/ ./shared/

# Copy worker module files
COPY workers/data/go.mod workers/data/go.sum ./
COPY workers/data/ ./

# Add replace directive for local shared module (overrides the one in go.mod for Docker build)
RUN go mod edit -replace github.com/RicardoCenci/iot-distributed-architecture/shared=./shared

RUN go mod download
RUN go mod verify

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates curl
WORKDIR /root/

COPY --from=builder /app/main .

COPY workers/data/config.json ./config.json

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
	CMD pgrep main || exit 1

CMD ["./main"]
