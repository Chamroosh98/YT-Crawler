
# ===== Build stage =====
FROM golang:1.22-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build Telegram bot
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /out/bot ./cmd/bot

# ===== Runtime stage =====
FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata \
    && adduser -D -H -s /sbin/nologin appuser

COPY --from=builder /out/bot /app/bot
COPY config /app/config

USER appuser

ENV TZ=UTC

CMD ["/app/bot"]