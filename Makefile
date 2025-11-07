.PHONY: help dev dev-auth dev-user dev-product dev-project dev-contact dev-portfolio stop clean

help: ## Показать справку
	@echo "Доступные команды:"
	@echo "  make dev          - Запустить все микросервисы в фоне"
	@echo "  make dev-auth     - Запустить только auth-service"
	@echo "  make dev-user     - Запустить только user-service"
	@echo "  make dev-product  - Запустить только product-service"
	@echo "  make dev-project  - Запустить только project-service"
	@echo "  make dev-contact  - Запустить только contact-service"
	@echo "  make dev-portfolio - Запустить только portfolio-service"
	@echo "  make stop         - Остановить все запущенные сервисы"
	@echo "  make clean        - Очистить логи и временные файлы"

dev: ## Запустить все микросервисы
	@echo "🚀 Запуск всех микросервисов..."
	@./scripts/dev.sh

dev-auth: ## Запустить auth-service
	@echo "🚀 Запуск auth-service..."
	@cd auth-service && go run main.go

dev-user: ## Запустить user-service
	@echo "🚀 Запуск user-service..."
	@cd user-service && go run main.go

dev-product: ## Запустить product-service
	@echo "🚀 Запуск product-service..."
	@cd product-service && go run main.go

dev-project: ## Запустить project-service
	@echo "🚀 Запуск project-service..."
	@cd project-service && go run main.go

dev-contact: ## Запустить contact-service
	@echo "🚀 Запуск contact-service..."
	@cd contact-service && go run cmd/main.go

dev-portfolio: ## Запустить portfolio-service
	@echo "🚀 Запуск portfolio-service..."
	@cd portfolio-service && go run main.go

stop: ## Остановить все сервисы
	@echo "🛑 Остановка всех сервисов..."
	@pkill -f "go run.*main.go" || true
	@pkill -f "go run.*cmd/main.go" || true
	@echo "✅ Все сервисы остановлены"

clean: ## Очистить логи и временные файлы
	@echo "🧹 Очистка..."
	@rm -f logs/*.log
	@rm -f *.pid
	@echo "✅ Очистка завершена"

