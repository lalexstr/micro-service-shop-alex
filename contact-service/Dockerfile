# Build stage
FROM golang:1.23-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем весь код проекта
COPY . .

# Инициализируем go.mod если его нет и устанавливаем зависимости
RUN if [ ! -f go.mod ]; then \
      go mod init ooolalex/contact-service; \
    fi && \
    go mod tidy && \
    go mod download

# Собираем приложение (CGO нужен для SQLite)
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o contact-service ./cmd/main.go

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

