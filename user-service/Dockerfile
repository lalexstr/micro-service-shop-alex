# syntax=docker/dockerfile:1.4
# Build stage
FROM golang:1.24-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Настраиваем быстрые зеркала Alpine перед установкой пакетов
RUN echo "http://mirrors.edge.kernel.org/alpine/v3.22/main" > /etc/apk/repositories && \
    echo "http://mirrors.edge.kernel.org/alpine/v3.22/community" >> /etc/apk/repositories && \
    echo "http://dl-cdn.alpinelinux.org/alpine/v3.22/main" >> /etc/apk/repositories && \
    echo "http://dl-cdn.alpinelinux.org/alpine/v3.22/community" >> /etc/apk/repositories

# Устанавливаем C компилятор для CGO (нужен для SQLite)
RUN apk update --no-cache && \
    apk add --no-cache gcc musl-dev

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
    CGO_ENABLED=1 GOOS=linux go build -installsuffix cgo -o user-service ./main.go

# Final stage
FROM alpine:latest

# Настраиваем быстрые зеркала Alpine
RUN echo "http://mirrors.edge.kernel.org/alpine/v3.22/main" > /etc/apk/repositories && \
    echo "http://mirrors.edge.kernel.org/alpine/v3.22/community" >> /etc/apk/repositories && \
    echo "http://dl-cdn.alpinelinux.org/alpine/v3.22/main" >> /etc/apk/repositories && \
    echo "http://dl-cdn.alpinelinux.org/alpine/v3.22/community" >> /etc/apk/repositories

# Устанавливаем необходимые пакеты для SQLite и wget для healthcheck
RUN apk update --no-cache && \
    apk add --no-cache ca-certificates sqlite wget

WORKDIR /app

# Копируем бинарный файл из builder stage
COPY --from=builder /app/user-service .

# Создаем директорию для базы данных
RUN mkdir -p /app/data

# Указываем порт
EXPOSE 8085

# Запускаем приложение
CMD ["./user-service"]
