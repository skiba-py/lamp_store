FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o orders_service ./cmd/api

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/orders_service .
COPY --from=builder /app/config.yaml .
COPY --from=builder /app/migrations ./migrations
COPY ./migrations /app/migrations

RUN adduser -D -g '' appuser && \
    mkdir -p /app/migrations && \
    chown -R appuser:appuser /app

USER appuser

CMD ["./orders_service"] 