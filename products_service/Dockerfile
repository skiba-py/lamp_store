FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o products_service ./cmd/api

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/products_service .
COPY --from=builder /app/config.yaml .

RUN ls -la /app/config.yaml

RUN adduser -D -g '' appuser
USER appuser

CMD ["./products_service"] 