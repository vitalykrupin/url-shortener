# Makefile для URL Shortener

# Переменные
APP_NAME = shortener
BINARY_NAME = shortener
MAIN_FILE = cmd/shortener/main.go
TEST_FLAGS = -v

# Цвета для вывода
RED = \033[0;31m
GREEN = \033[0;32m
YELLOW = \033[1;33m
NC = \033[0m # No Color

# Форматирование вывода
INFO = @echo "$(GREEN)[INFO]$(NC)"
WARN = @echo "$(YELLOW)[WARN]$(NC)"
ERR = @echo "$(RED)[ERROR]$(NC)"

# Цели по умолчанию
.PHONY: all
all: build

# Сборка приложения
.PHONY: build
build:
	$(INFO) Building $(APP_NAME)...
	go build -o $(BINARY_NAME) $(MAIN_FILE)
	$(INFO) Build completed successfully.

# Запуск приложения
.PHONY: run
run:
	$(INFO) Running $(APP_NAME)...
	go run $(MAIN_FILE)

# Установка зависимостей
.PHONY: deps
deps:
	$(INFO) Installing dependencies...
	go mod tidy
	$(INFO) Dependencies installed.

# Тестирование
.PHONY: test
test:
	$(INFO) Running tests...
	go test $(TEST_FLAGS) ./...

# Тестирование с покрытием
.PHONY: test-cover
test-cover:
	$(INFO) Running tests with coverage...
	go test $(TEST_FLAGS) -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	$(INFO) Coverage report generated: coverage.html

# Очистка
.PHONY: clean
clean:
	$(INFO) Cleaning up...
	rm -f $(BINARY_NAME)
	rm -f coverage.out
	rm -f coverage.html
	$(INFO) Cleanup completed.

# Линтинг
.PHONY: lint
lint:
	$(INFO) Running linter...
	golangci-lint run
	$(INFO) Linting completed.

# Форматирование кода
.PHONY: fmt
fmt:
	$(INFO) Formatting code...
	go fmt ./...
	$(INFO) Code formatting completed.

# Проверка кода
.PHONY: vet
vet:
	$(INFO) Vetting code...
	go vet ./...
	$(INFO) Code vetting completed.

# Запуск всех проверок
.PHONY: check
check: fmt vet lint test

# Установка инструментов разработки
.PHONY: install-tools
install-tools:
	$(INFO) Installing development tools...
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(INFO) Development tools installed.

# Помощь
.PHONY: help
help:
	@echo "Доступные команды:"
	@echo "  all          - Сборка приложения (по умолчанию)"
	@echo "  build        - Сборка приложения"
	@echo "  run          - Запуск приложения"
	@echo "  deps         - Установка зависимостей"
	@echo "  test         - Запуск тестов"
	@echo "  test-cover   - Запуск тестов с покрытием"
	@echo "  clean        - Очистка"
	@echo "  lint         - Линтинг кода"
	@echo "  fmt          - Форматирование кода"
	@echo "  vet          - Проверка кода"
	@echo "  check        - Запуск всех проверок"
	@echo "  install-tools - Установка инструментов разработки"
	@echo "  help         - Показать эту справку"