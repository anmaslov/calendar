#!/bin/bash
set -e

# Calendar CLI - единый скрипт управления проектом
# Использование: ./calendar.sh <команда> [опции]

BINARY="calendar"
BUILD_DIR="bin"
CONFIG_PATH="${CONFIG_PATH:-configs/config.yaml}"
COMPOSE_FILE="docker-compose.yml"
COMPOSE_APP_FILE="docker-compose.app.yml"

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_header() {
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${BLUE}  Calendar - $1${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

print_info() {
    echo -e "${YELLOW}→ $1${NC}"
}

# Команда: build
cmd_build() {
    print_header "Сборка приложения"
    
    mkdir -p "$BUILD_DIR"
    print_info "Компиляция..."
    go build -o "$BUILD_DIR/$BINARY" ./cmd/calendar
    
    print_success "Сборка завершена: $BUILD_DIR/$BINARY"
}

# Команда: run
cmd_run() {
    print_header "Запуск приложения"
    
    local config="${1:-$CONFIG_PATH}"
    print_info "Конфигурация: $config"
    
    go run ./cmd/calendar --config="$config"
}

# Команда: test
cmd_test() {
    print_header "Запуск тестов"
    
    go test -v -race ./...
    
    print_success "Тесты завершены"
}

# Команда: test-coverage
cmd_test_coverage() {
    print_header "Тесты с покрытием"
    
    go test -v -race -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html
    
    print_success "Отчёт о покрытии: coverage.html"
}

# Команда: lint
cmd_lint() {
    print_header "Линтер"
    
    if ! command -v golangci-lint &> /dev/null; then
        print_error "golangci-lint не установлен"
        print_info "Установите: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
        exit 1
    fi
    
    golangci-lint run
    
    print_success "Линтер завершён"
}

# Команда: deps
cmd_deps() {
    print_header "Обновление зависимостей"
    
    print_info "Загрузка зависимостей..."
    go mod download
    
    print_info "Очистка go.mod..."
    go mod tidy
    
    print_success "Зависимости обновлены"
}

# Команда: db-up (только PostgreSQL)
cmd_db_up() {
    print_header "Запуск PostgreSQL"
    
    docker-compose -f "$COMPOSE_FILE" up -d
    
    print_success "PostgreSQL запущен"
    docker-compose -f "$COMPOSE_FILE" ps
}

# Команда: db-down
cmd_db_down() {
    print_header "Остановка PostgreSQL"
    
    docker-compose -f "$COMPOSE_FILE" down
    
    print_success "PostgreSQL остановлен"
}

# Команда: db-logs
cmd_db_logs() {
    print_header "Логи PostgreSQL"
    
    docker-compose -f "$COMPOSE_FILE" logs -f postgres
}

# Команда: app-up (полный стек)
cmd_app_up() {
    print_header "Запуск полного стека (PostgreSQL + App)"
    
    docker-compose -f "$COMPOSE_APP_FILE" up -d
    
    print_success "Сервисы запущены"
    docker-compose -f "$COMPOSE_APP_FILE" ps
}

# Команда: app-down
cmd_app_down() {
    print_header "Остановка полного стека"
    
    docker-compose -f "$COMPOSE_APP_FILE" down
    
    print_success "Сервисы остановлены"
}

# Команда: app-build
cmd_app_build() {
    print_header "Сборка Docker образа"
    
    docker-compose -f "$COMPOSE_APP_FILE" build
    
    print_success "Docker образ собран"
}

# Команда: app-logs
cmd_app_logs() {
    print_header "Логи сервисов"
    
    local service="${1:-}"
    if [ -n "$service" ]; then
        docker-compose -f "$COMPOSE_APP_FILE" logs -f "$service"
    else
        docker-compose -f "$COMPOSE_APP_FILE" logs -f
    fi
}

# Команда: app-restart
cmd_app_restart() {
    print_header "Перезапуск сервисов"
    
    docker-compose -f "$COMPOSE_APP_FILE" restart
    
    print_success "Сервисы перезапущены"
}

# Команда: clean
cmd_clean() {
    print_header "Очистка"
    
    print_info "Удаление бинарников..."
    rm -rf "$BUILD_DIR"
    
    print_info "Удаление файлов покрытия..."
    rm -f coverage.out coverage.html
    
    print_success "Очистка завершена"
}

# Команда: help
cmd_help() {
    echo -e "${BLUE}Calendar CLI${NC} - управление проектом"
    echo ""
    echo -e "${YELLOW}Использование:${NC}"
    echo "  ./calendar.sh <команда> [опции]"
    echo ""
    echo -e "${YELLOW}Команды разработки:${NC}"
    echo "  build              Сборка бинарника"
    echo "  run [config]       Запуск приложения (опционально: путь к конфигу)"
    echo "  test               Запуск тестов"
    echo "  test-coverage      Тесты с отчётом о покрытии"
    echo "  lint               Запуск линтера"
    echo "  deps               Обновление зависимостей"
    echo "  clean              Очистка артефактов сборки"
    echo ""
    echo -e "${YELLOW}База данных (только PostgreSQL для разработки):${NC}"
    echo "  db-up              Запуск PostgreSQL"
    echo "  db-down            Остановка PostgreSQL"
    echo "  db-logs            Логи PostgreSQL"
    echo ""
    echo -e "${YELLOW}Полный стек (PostgreSQL + приложение):${NC}"
    echo "  app-up             Запуск полного стека"
    echo "  app-down           Остановка полного стека"
    echo "  app-build          Сборка Docker образа"
    echo "  app-restart        Перезапуск сервисов"
    echo "  app-logs [svc]     Логи (опционально: app или postgres)"
    echo ""
    echo -e "${YELLOW}Переменные окружения:${NC}"
    echo "  CONFIG_PATH        Путь к конфигурации (по умолчанию: configs/config.yaml)"
    echo "  DB_PASSWORD        Пароль базы данных"
    echo ""
    echo -e "${YELLOW}Примеры:${NC}"
    echo "  ./calendar.sh build"
    echo "  ./calendar.sh db-up && ./calendar.sh run"
    echo "  ./calendar.sh run configs/config.local.yaml"
    echo "  DB_PASSWORD=secret ./calendar.sh app-up"
    echo "  ./calendar.sh app-logs app"
}

# Основная логика
main() {
    local command="${1:-help}"
    shift || true
    
    case "$command" in
        build)
            cmd_build "$@"
            ;;
        run)
            cmd_run "$@"
            ;;
        test)
            cmd_test "$@"
            ;;
        test-coverage|coverage)
            cmd_test_coverage "$@"
            ;;
        lint)
            cmd_lint "$@"
            ;;
        deps)
            cmd_deps "$@"
            ;;
        # База данных
        db-up|db)
            cmd_db_up "$@"
            ;;
        db-down)
            cmd_db_down "$@"
            ;;
        db-logs)
            cmd_db_logs "$@"
            ;;
        # Полный стек
        app-up|up)
            cmd_app_up "$@"
            ;;
        app-down|down)
            cmd_app_down "$@"
            ;;
        app-build|docker-build)
            cmd_app_build "$@"
            ;;
        app-restart|restart)
            cmd_app_restart "$@"
            ;;
        app-logs|logs)
            cmd_app_logs "$@"
            ;;
        clean)
            cmd_clean "$@"
            ;;
        help|--help|-h)
            cmd_help
            ;;
        *)
            print_error "Неизвестная команда: $command"
            echo ""
            cmd_help
            exit 1
            ;;
    esac
}

main "$@"
