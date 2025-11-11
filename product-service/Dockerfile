# syntax=docker/dockerfile:1.4
# Build stage
FROM golang:1.24-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Устанавливаем C компилятор для CGO (нужен для SQLite) с retry логикой
RUN for i in 1 2 3; do \
        apk update --no-cache && \
        apk add --no-cache gcc musl-dev && break || \
        (echo "Попытка $i не удалась, повтор через 5 секунд..." && sleep 5); \
    done

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
    CGO_ENABLED=1 GOOS=linux go build -installsuffix cgo -o product-service ./main.go

# Final stage
FROM alpine:latest

# Устанавливаем необходимые пакеты для SQLite и wget для healthcheck с retry логикой
RUN for i in 1 2 3; do \
        apk update --no-cache && \
        apk add --no-cache ca-certificates sqlite wget && break || \
        (echo "Попытка $i не удалась, повтор через 5 секунд..." && sleep 5); \
    done

WORKDIR /app

# Копируем бинарный файл из builder stage
COPY --from=builder /app/product-service .

# Копируем .env файл (если нужен)
# COPY .env .env

# Создаем директорию для базы данных
RUN mkdir -p /app/data

# Указываем порт
EXPOSE 8081

# Запускаем приложение
CMD ["./product-service"]
