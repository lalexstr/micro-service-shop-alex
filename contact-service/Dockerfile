# syntax=docker/dockerfile:1.4
# Build stage
FROM golang:1.24-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Устанавливаем C компилятор для CGO (нужен для SQLite)
RUN apk add --no-cache gcc musl-dev

# Копируем go mod файлы
COPY go.mod go.sum ./

# Используем BuildKit cache mount для ускорения загрузки зависимостей
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Копируем весь код проекта
COPY . .

# Собираем приложение с использованием кеша для модулей
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=1 GOOS=linux go build -installsuffix cgo -o contact-service ./cmd/main.go

# Final stage
FROM alpine:latest

# Устанавливаем необходимые пакеты для SQLite и wget для healthcheck
RUN apk --no-cache add ca-certificates sqlite wget

WORKDIR /app

# Копируем бинарный файл из builder stage
COPY --from=builder /app/contact-service .

# Создаем директорию для базы данных
RUN mkdir -p /app/data

# Указываем порт
EXPOSE 8084

# Запускаем приложение
CMD ["./contact-service"]

