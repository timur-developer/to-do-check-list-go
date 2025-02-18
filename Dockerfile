# Используем образ с поддержкой CGO
FROM golang:1.23 AS builder

WORKDIR /app

# Устанавливаем зависимости
RUN apt-get update && apt-get install -y git gcc g++ librdkafka-dev

# Копируем модули и загружаем зависимости
COPY go.mod go.sum ./
RUN go mod tidy

# Копируем код
COPY . .

# Компилируем с CGO
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o main ./cmd

# Финальный контейнер на Debian (не Alpine!)
FROM debian:bookworm-slim

WORKDIR /app

# Устанавливаем librdkafka
RUN apt-get update && apt-get install -y librdkafka1

COPY --from=builder /app/main /app/main
COPY --from=builder /app/.env .env

# Делаем исполняемым
RUN chmod +x /app/main

CMD ["/app/main"]
