# Calendar

Сервис синхронизации календаря с корпоративным Exchange сервером. Сохраняет события в PostgreSQL и предоставляет REST API для чтения данных.

## Возможности

- Фоновая синхронизация событий с Microsoft Exchange сервером
- Хранение событий в PostgreSQL
- REST API для получения событий (только чтение)
- Kubernetes-ready (liveness/readiness probes)
- Graceful shutdown
- Docker поддержка

## Архитектура

Проект построен по принципам Clean Architecture:

```
├── cmd/calendar/          # Точка входа приложения
├── configs/               # Конфигурационные файлы
├── internal/
│   ├── config/           # Загрузка конфигурации
│   ├── domain/           # Доменные модели и ошибки
│   ├── handler/          # HTTP handlers (delivery layer)
│   ├── repository/       # Слой доступа к данным
│   │   └── postgres/     # PostgreSQL реализация
│   ├── service/          # Бизнес-логика
│   └── sync/             # Фоновая синхронизация с Exchange
├── migrations/           # SQL миграции
├── calendar.sh           # CLI скрипт управления проектом
└── .cursor/              # Конфигурация Cursor IDE
```

## Требования

- Go 1.23+
- PostgreSQL 16+
- Docker & Docker Compose (опционально)

## Быстрый старт

### Через Docker Compose

1. Создайте файл `.env` с паролем БД:
   ```bash
   echo "DB_PASSWORD=your_secure_password" > .env
   ```

2. Запустите полный стек:
   ```bash
   ./calendar.sh app-up
   ```

3. API доступен по адресу `http://localhost:8080`

### Локальная разработка

1. Скопируйте конфигурацию:
   ```bash
   cp configs/config.local.yaml.example configs/config.local.yaml
   # Отредактируйте файл, укажите пароль БД
   ```

2. Запустите PostgreSQL:
   ```bash
   ./calendar.sh db-up
   ```

3. Запустите приложение:
   ```bash
   ./calendar.sh run
   # или с указанием конфига:
   ./calendar.sh run configs/config.local.yaml
   ```

## Конфигурация

Конфигурация загружается из YAML файла. Путь по умолчанию: `configs/config.yaml`

### Структура конфигурации

```yaml
server:
  port: 8080
  read_timeout: 15s
  write_timeout: 15s
  idle_timeout: 60s
  shutdown_timeout: 30s

database:
  host: localhost
  port: 5432
  user: calendar
  password: ""  # Или через переменную DB_PASSWORD
  name: calendar
  ssl_mode: disable
  max_open_conns: 25
  max_idle_conns: 5

exchange:
  url: "https://your-exchange-server.com/ews/exchange.asmx"
  username: ""
  password: ""  # Или через переменную EXCHANGE_PASSWORD
  domain: ""

sync:
  enabled: false   # Включить фоновую синхронизацию
  interval: 5m     # Интервал синхронизации
  sync_days: 30    # На сколько дней вперёд синхронизировать

logging:
  level: info      # debug, info, warn, error
  format: json     # json, console
```

### Переменные окружения

Чувствительные данные можно передать через переменные окружения:
- `DB_PASSWORD` — пароль базы данных
- `EXCHANGE_PASSWORD` — пароль Exchange сервера

## CLI Скрипт

Все операции выполняются через единый скрипт `./calendar.sh`:

```bash
./calendar.sh <команда> [опции]
```

### Команды разработки

| Команда | Описание |
|---------|----------|
| `build` | Сборка бинарника в `bin/calendar` |
| `run [config]` | Запуск приложения (опционально: путь к конфигу) |
| `test` | Запуск тестов |
| `test-coverage` | Тесты с отчётом о покрытии |
| `lint` | Запуск линтера golangci-lint |
| `deps` | Обновление зависимостей |
| `clean` | Очистка артефактов сборки |

### База данных (только PostgreSQL)

| Команда | Описание |
|---------|----------|
| `db-up` / `db` | Запуск PostgreSQL |
| `db-down` | Остановка PostgreSQL |
| `db-logs` | Логи PostgreSQL |

### Полный стек (PostgreSQL + приложение)

| Команда | Описание |
|---------|----------|
| `app-up` / `up` | Запуск полного стека |
| `app-down` / `down` | Остановка полного стека |
| `app-build` | Сборка Docker образа |
| `app-restart` / `restart` | Перезапуск сервисов |
| `app-logs [svc]` / `logs` | Просмотр логов |

### Примеры использования

```bash
# Сборка
./calendar.sh build

# Локальная разработка (PostgreSQL в Docker, приложение локально)
./calendar.sh db-up
./calendar.sh run configs/config.local.yaml

# Полный стек в Docker
DB_PASSWORD=secret ./calendar.sh app-up
./calendar.sh app-logs app
./calendar.sh app-down

# Тесты
./calendar.sh test
./calendar.sh test-coverage

# Через переменную окружения
CONFIG_PATH=configs/config.local.yaml ./calendar.sh run
```

## API Endpoints

### Health Checks

| Endpoint | Описание |
|----------|----------|
| `GET /health` | Базовая проверка здоровья |
| `GET /healthz` | Kubernetes liveness probe |
| `GET /readyz` | Kubernetes readiness probe |

### События (только чтение)

| Метод | Endpoint | Описание |
|-------|----------|----------|
| GET | `/api/v1/events` | Список событий |
| GET | `/api/v1/events/{id}` | Получение события по ID |

### Параметры запроса для списка событий

| Параметр | Описание | По умолчанию |
|----------|----------|--------------|
| `limit` | Количество событий | 20 |
| `offset` | Смещение для пагинации | 0 |
| `start_date` | Фильтр по дате начала (RFC3339) | — |
| `end_date` | Фильтр по дате окончания (RFC3339) | — |
| `subject` | Поиск по теме (частичное совпадение) | — |
| `status` | Фильтр по статусу | — |

### Примеры запросов

**Получить список событий:**
```bash
curl "http://localhost:8080/api/v1/events?limit=10"
```

**Получить события за период:**
```bash
curl "http://localhost:8080/api/v1/events?start_date=2024-01-01T00:00:00Z&end_date=2024-01-31T23:59:59Z"
```

**Получить событие по ID:**
```bash
curl "http://localhost:8080/api/v1/events/550e8400-e29b-41d4-a716-446655440000"
```

### Формат ответа

**Список событий:**
```json
{
  "events": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "exchange_id": "AAMkAGI2...",
      "subject": "Совещание команды",
      "body": "Обсуждение планов на квартал",
      "location": "Переговорная А",
      "start_time": "2024-01-15T10:00:00Z",
      "end_time": "2024-01-15T11:00:00Z",
      "is_all_day": false,
      "organizer": "user@company.com",
      "status": "confirmed",
      "created_at": "2024-01-10T12:00:00Z",
      "updated_at": "2024-01-10T12:00:00Z",
      "synced_at": "2024-01-10T12:00:00Z"
    }
  ],
  "total": 1,
  "limit": 20,
  "offset": 0
}
```

**Ошибка:**
```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "Event not found"
  }
}
```

## Синхронизация с Exchange

События синхронизируются автоматически в фоновом режиме, если включена опция `sync.enabled`.

### Как работает синхронизация:

1. Воркер запускается при старте приложения
2. Каждые N минут (настраивается через `sync.interval`) запрашивает события из Exchange
3. События за период от текущей даты + `sync_days` дней вперёд
4. Новые события добавляются, существующие обновляются (по `exchange_id`)
5. События, удалённые из Exchange, удаляются из локальной БД

### Включение синхронизации:

```yaml
sync:
  enabled: true
  interval: 5m
  sync_days: 30

exchange:
  url: https://mail.company.com/ews/exchange.asmx
  username: calendar_service
  password: ""  # Через EXCHANGE_PASSWORD
  domain: COMPANY
```

> **Примечание:** Текущая реализация Exchange клиента — заглушка. Требуется реализовать интерфейс `ExchangeClient` с использованием EWS API.

## Kubernetes

### Пример манифеста Deployment:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: calendar
spec:
  replicas: 1
  template:
    spec:
      containers:
        - name: calendar
          image: calendar:latest
          ports:
            - containerPort: 8080
          env:
            - name: DB_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: calendar-secrets
                  key: db-password
          livenessProbe:
            httpGet:
              path: /healthz
              port: 8080
            initialDelaySeconds: 5
            periodSeconds: 10
          readinessProbe:
            httpGet:
              path: /readyz
              port: 8080
            initialDelaySeconds: 5
            periodSeconds: 5
          volumeMounts:
            - name: config
              mountPath: /app/configs
      volumes:
        - name: config
          configMap:
            name: calendar-config
```

## Лицензия

MIT
