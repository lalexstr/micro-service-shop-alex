#!/bin/bash

# Скрипт для запуска всех микросервисов в режиме разработки

set -e

# Цвета для вывода
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Функция для остановки всех процессов при выходе
cleanup() {
    echo -e "\n${YELLOW}Остановка всех сервисов...${NC}"
    pkill -f "go run.*main.go" || true
    pkill -f "go run.*cmd/main.go" || true
    exit 0
}

# Перехватываем сигналы для корректной остановки
trap cleanup SIGINT SIGTERM

echo -e "${GREEN}🚀 Запуск всех микросервисов в режиме разработки...${NC}\n"

# Проверяем, что все сервисы имеют необходимые зависимости
echo -e "${YELLOW}Проверка зависимостей...${NC}"
for service in auth-service user-service product-service project-service contact-service portfolio-service; do
    if [ -d "$service" ]; then
        echo "  ✓ $service"
    else
        echo -e "  ${RED}✗ $service не найден${NC}"
    fi
done

echo ""

# Запускаем каждый сервис в отдельном терминале/процессе
echo -e "${GREEN}Запуск сервисов...${NC}\n"

# Auth Service (порт 8080) - должен запуститься первым
echo -e "${YELLOW}→ auth-service (порт 8080)${NC}"
cd auth-service && go run main.go > ../logs/auth-service.log 2>&1 &
AUTH_PID=$!
cd ..
sleep 2

# User Service (порт 8085)
echo -e "${YELLOW}→ user-service (порт 8085)${NC}"
cd user-service && go run main.go > ../logs/user-service.log 2>&1 &
USER_PID=$!
cd ..
sleep 1

# Product Service (порт 8081)
echo -e "${YELLOW}→ product-service (порт 8081)${NC}"
cd product-service && go run main.go > ../logs/product-service.log 2>&1 &
PRODUCT_PID=$!
cd ..
sleep 1

# Project Service (порт 8082)
echo -e "${YELLOW}→ project-service (порт 8082)${NC}"
cd project-service && go run main.go > ../logs/project-service.log 2>&1 &
PROJECT_PID=$!
cd ..
sleep 1

# Contact Service (порт 8084)
echo -e "${YELLOW}→ contact-service (порт 8084)${NC}"
cd contact-service && go run cmd/main.go > ../logs/contact-service.log 2>&1 &
CONTACT_PID=$!
cd ..
sleep 1

# Portfolio Service (порт 8083)
echo -e "${YELLOW}→ portfolio-service (порт 8083)${NC}"
cd portfolio-service && go run main.go > ../logs/portfolio-service.log 2>&1 &
PORTFOLIO_PID=$!
cd ..
sleep 1

echo ""
echo -e "${GREEN}✅ Все сервисы запущены!${NC}"
echo ""
echo "Порты сервисов:"
echo "  - auth-service:     http://localhost:8080"
echo "  - product-service:  http://localhost:8081"
echo "  - project-service:  http://localhost:8082"
echo "  - portfolio-service: http://localhost:8083"
echo "  - contact-service:  http://localhost:8084"
echo "  - user-service:    http://localhost:8085"
echo ""
echo -e "${YELLOW}Логи находятся в директории logs/${NC}"
echo -e "${YELLOW}Нажмите Ctrl+C для остановки всех сервисов${NC}"
echo ""

# Ждем завершения всех процессов
wait





