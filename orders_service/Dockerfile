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

CMD ["./orders_service"] 