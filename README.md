# Calendar

Calendar synchronization service that syncs events from corporate Exchange server, stores them in PostgreSQL, and provides a REST API.

## Features

- Sync calendar events from Microsoft Exchange server
- Store events in PostgreSQL database
- RESTful API for managing events
- Docker support for easy deployment

## Architecture

The project follows Clean Architecture principles:

```
├── cmd/calendar/          # Application entry point
├── internal/
│   ├── config/           # Configuration management
│   ├── domain/           # Domain models and errors
│   ├── handler/          # HTTP handlers (delivery layer)
│   ├── repository/       # Data access layer
│   │   └── postgres/     # PostgreSQL implementation
│   └── service/          # Business logic layer
├── migrations/           # Database migrations
└── .cursor/              # Cursor IDE configuration
```

## Prerequisites

- Go 1.23+
- PostgreSQL 16+
- Docker & Docker Compose (optional)

## Quick Start

### Using Docker Compose

1. Copy environment file and configure:
   ```bash
   cp .env.example .env
   # Edit .env with your settings
   ```

2. Start services:
   ```bash
   docker-compose up -d
   ```

3. The API will be available at `http://localhost:8080`

### Local Development

1. Copy environment file:
   ```bash
   cp .env.example .env
   ```

2. Start PostgreSQL:
   ```bash
   docker-compose up -d postgres
   ```

3. Run the application:
   ```bash
   go run ./cmd/calendar
   ```

## API Endpoints

### Health Check
```
GET /health
```

### Events

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/events` | List events |
| POST | `/api/v1/events` | Create event |
| GET | `/api/v1/events/{id}` | Get event by ID |
| PUT | `/api/v1/events/{id}` | Update event |
| DELETE | `/api/v1/events/{id}` | Delete event |

### Sync

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/sync` | Trigger Exchange sync |

### Query Parameters for List Events

- `limit` - Number of events to return (default: 20)
- `offset` - Offset for pagination
- `start_date` - Filter by start date (RFC3339 format)
- `end_date` - Filter by end date (RFC3339 format)
- `subject` - Filter by subject (partial match)
- `status` - Filter by status

### Example Requests

**Create Event:**
```bash
curl -X POST http://localhost:8080/api/v1/events \
  -H "Content-Type: application/json" \
  -d '{
    "subject": "Team Meeting",
    "start_time": "2024-01-15T10:00:00Z",
    "end_time": "2024-01-15T11:00:00Z",
    "location": "Conference Room A"
  }'
```

**List Events:**
```bash
curl "http://localhost:8080/api/v1/events?limit=10&start_date=2024-01-01T00:00:00Z"
```

## Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_PORT` | HTTP server port | 8080 |
| `DB_HOST` | PostgreSQL host | localhost |
| `DB_PORT` | PostgreSQL port | 5432 |
| `DB_USER` | Database user | calendar |
| `DB_PASSWORD` | Database password | (required) |
| `DB_NAME` | Database name | calendar |
| `EXCHANGE_URL` | Exchange EWS endpoint | - |
| `EXCHANGE_USERNAME` | Exchange username | - |
| `EXCHANGE_PASSWORD` | Exchange password | - |
| `EXCHANGE_DOMAIN` | Exchange domain | - |

## Development

### Running Tests
```bash
go test ./...
```

### Building
```bash
go build -o calendar ./cmd/calendar
```

## License

MIT

