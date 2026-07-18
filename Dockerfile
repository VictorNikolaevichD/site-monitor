# Этап 1: сборка приложения
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /build/site-monitor \
    ./cmd/monitor

# Этап 2: запуск приложения
FROM alpine:3.22 AS runtime

LABEL org.opencontainers.image.title="site-monitor" \
      org.opencontainers.image.description="HTTP service for site availability monitoring" \
      org.opencontainers.image.source="https://gitlab.com/Dokuchaevvn/site-monitor"

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /build/site-monitor ./site-monitor

EXPOSE 8080

ENTRYPOINT ["./site-monitor"]

CMD ["-config", "/configs/sites.yaml"]
