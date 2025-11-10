# Build stage
FROM golang:1.24-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Устанавливаем C компилятор для CGO (нужен для SQLite)
RUN apk add --no-cache gcc musl-dev

# Копируем go mod файлы
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь код проекта
COPY . .

# Обновляем зависимости и собираем приложение (CGO нужен для SQLite)
RUN go mod tidy && CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o product-service ./main.go

# Final stage
FROM alpine:latest

# Устанавливаем необходимые пакеты для SQLite и wget для healthcheck
RUN apk --no-cache add ca-certificates sqlite wget

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

