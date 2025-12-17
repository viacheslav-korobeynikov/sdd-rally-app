.PHONY: help dev build run test test-unit test-integration test-coverage lint fmt migrate-up migrate-down migrate-create seed docker-up docker-down clean

# === Переменные ===
APP_NAME=sdd-rally-app
VERSION?=0.1.0
BUILD_DIR=./bin
MAIN_PATH=./cmd/server

# Go параметры
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt

# Database
DB_URL?=postgres://rally_user:changeme_dev_password@localhost:5432/rally_dev?sslmode=disable
MIGRATIONS_DIR=./migrations

# === Основные команды ===

help: ## Показать справку
	@echo "SDD Rally App - Доступные команды:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

dev: ## Запустить приложение с hot reload (требует air)
	@if ! command -v air > /dev/null; then \
		echo "Air не установлен. Установите: go install github.com/cosmtrek/air@latest"; \
		exit 1; \
	fi
	air

run: ## Запустить приложение
	$(GOCMD) run $(MAIN_PATH)/main.go

build: ## Собрать бинарный файл
	@echo "Сборка $(APP_NAME) v$(VERSION)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -ldflags="-X 'main.Version=$(VERSION)'" -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)
	@echo "Готово: $(BUILD_DIR)/$(APP_NAME)"

install-deps: ## Установить зависимости
	$(GOMOD) download
	$(GOMOD) verify

# === Тестирование ===

test: ## Запустить все тесты
	$(GOTEST) -v -race -timeout 30s ./...

test-unit: ## Запустить unit тесты
	$(GOTEST) -v -race -timeout 30s -tags=unit ./...

test-integration: ## Запустить интеграционные тесты
	$(GOTEST) -v -race -timeout 60s -tags=integration ./...

test-coverage: ## Запустить тесты с coverage
	$(GOTEST) -v -race -coverprofile=coverage.out -covermode=atomic ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage отчёт: coverage.html"

# === Линтеры и форматирование ===

lint: ## Запустить golangci-lint
	@if ! command -v golangci-lint > /dev/null; then \
		echo "golangci-lint не установлен. Установите: https://golangci-lint.run/usage/install/"; \
		exit 1; \
	fi
	golangci-lint run ./...

fmt: ## Форматировать код
	$(GOFMT) ./...
	@if command -v goimports > /dev/null; then \
		goimports -w .; \
	fi

# === База данных ===

migrate-up: ## Применить миграции
	@if ! command -v migrate > /dev/null; then \
		echo "migrate не установлен. Установите: https://github.com/golang-migrate/migrate"; \
		exit 1; \
	fi
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up

migrate-down: ## Откатить последнюю миграцию
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down 1

migrate-create: ## Создать новую миграцию (использование: make migrate-create name=create_users_table)
	@if [ -z "$(name)" ]; then \
		echo "Ошибка: укажите имя миграции. Пример: make migrate-create name=create_users_table"; \
		exit 1; \
	fi
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)

seed: ## Загрузить тестовые данные
	$(GOCMD) run ./scripts/seed/main.go

# === Docker ===

docker-up: ## Запустить Docker Compose
	docker-compose up -d

docker-down: ## Остановить Docker Compose
	docker-compose down

docker-logs: ## Показать логи Docker
	docker-compose logs -f

docker-rebuild: ## Пересобрать Docker образы
	docker-compose up -d --build

# === Утилиты ===

clean: ## Очистить артефакты сборки
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
	$(GOCMD) clean

clean-all: clean ## Очистить всё, включая зависимости
	rm -rf vendor/
	$(GOCMD) clean -modcache

generate: ## Запустить go generate
	$(GOCMD) generate ./...

mod-tidy: ## Привести go.mod в порядок
	$(GOMOD) tidy

security-check: ## Проверить уязвимости (требует govulncheck)
	@if ! command -v govulncheck > /dev/null; then \
		echo "govulncheck не установлен. Установите: go install golang.org/x/vuln/cmd/govulncheck@latest"; \
		exit 1; \
	fi
	govulncheck ./...

# === Разработка ===

setup: install-deps migrate-up seed ## Полная настройка для разработки
	@echo "Окружение готово! Запустите: make dev"

install-tools: ## Установить инструменты разработки
	go install github.com/cosmtrek/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install golang.org/x/tools/cmd/goimports@latest
	@echo "Инструменты установлены"

# === Production ===

build-prod: ## Собрать для продакшена (оптимизированный)
	@echo "Сборка для продакшена..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) \
		-ldflags="-w -s -X 'main.Version=$(VERSION)'" \
		-o $(BUILD_DIR)/$(APP_NAME)-linux-amd64 \
		$(MAIN_PATH)
	@echo "Готово: $(BUILD_DIR)/$(APP_NAME)-linux-amd64"

docker-build-prod: ## Собрать production Docker образ
	docker build -t $(APP_NAME):$(VERSION) -t $(APP_NAME):latest .

