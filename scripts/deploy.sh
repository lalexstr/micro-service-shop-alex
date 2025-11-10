#!/bin/bash

# Скрипт для деплоя всех микросервисов на сервер
# Использование: ./scripts/deploy.sh [environment]
# environment: production (по умолчанию) или staging

set -e

# Цвета для вывода
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

ENVIRONMENT=${1:-production}

echo -e "${GREEN}🚀 Начало деплоя микросервисов (${ENVIRONMENT})...${NC}\n"

# Проверяем наличие docker-compose
if ! command -v docker-compose &> /dev/null; then
    echo -e "${RED}❌ docker-compose не установлен${NC}"
    exit 1
fi

# Проверяем наличие Docker
if ! command -v docker &> /dev/null; then
    echo -e "${RED}❌ Docker не установлен${NC}"
    exit 1
fi

# Переходим в корневую директорию проекта
cd "$(dirname "$0")/.."

echo -e "${YELLOW}📦 Остановка существующих контейнеров...${NC}"
docker-compose down || true

echo -e "${YELLOW}🔨 Сборка Docker образов...${NC}"
docker-compose build --no-cache

echo -e "${YELLOW}🚀 Запуск всех сервисов...${NC}"
docker-compose up -d

echo -e "${YELLOW}⏳ Ожидание запуска сервисов...${NC}"
sleep 10

echo -e "${YELLOW}📊 Статус контейнеров:${NC}"
docker-compose ps

echo -e "\n${GREEN}✅ Деплой завершен!${NC}\n"

echo -e "${YELLOW}Порты сервисов:${NC}"
echo "  - auth-service:     http://localhost:8080"
echo "  - product-service:  http://localhost:8081"
echo "  - project-service:  http://localhost:8082"
echo "  - portfolio-service: http://localhost:8083"
echo "  - contact-service:  http://localhost:8084"
echo "  - user-service:    http://localhost:8085"
echo ""

echo -e "${YELLOW}Проверка health checks:${NC}"
docker-compose ps --format "table {{.Name}}\t{{.Status}}\t{{.Ports}}"

echo -e "\n${GREEN}Для просмотра логов используйте:${NC}"
echo "  docker-compose logs -f [service-name]"
echo ""
echo -e "${GREEN}Для остановки всех сервисов:${NC}"
echo "  docker-compose down"

