FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main ./cmd/api

FROM alpine:latest

WORKDIR /app

RUN mkdir -p frontend/public/images && chmod 755 frontend/public/images

COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/config.yaml .

EXPOSE 8000

CMD ["./main"] 