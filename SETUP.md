# SDD Rally App - Setup Guide

## Быстрый старт (Quick Start)

```powershell
# 1. Установите goose
go install github.com/pressly/goose/v3/cmd/goose@latest

# 2. Запустите PostgreSQL в Docker
docker-compose up -d postgres

# 3. Примените миграции
.\migrate.cmd up

# 4. Проверьте статус
.\migrate.cmd status

# 5. Запустите приложение
go run cmd/server/main.go
```

**Важно:** PostgreSQL работает на порту **5433** (не 5432), чтобы избежать конфликта с локальным PostgreSQL.

---

## Prerequisites

- **Go 1.21+** installed
- **Docker Desktop** installed and running
- **Git** for version control

## Initial Setup

### 1. Install Goose Migration Tool

Goose is used to manage database migrations. Install it globally:

```powershell
go install github.com/pressly/goose/v3/cmd/goose@latest
```

Verify installation:

```powershell
goose --version
```

### 2. Start PostgreSQL Database

The database runs in Docker. Start it with:

```powershell
docker-compose up -d postgres
```

Verify PostgreSQL is running:

```powershell
docker ps
```

You should see `sdd-rally-postgres` container running.

Check database health:

```powershell
docker logs sdd-rally-postgres
```

## Database Migrations

### Connection String

The default database connection string is:

```
postgresql://rally_user:rally_dev_password@127.0.0.1:5433/rally_dev?sslmode=disable
```

**Примечание:** Используется порт **5433**, чтобы избежать конфликта с локальным PostgreSQL на порту 5432.

### Важно: Настройка переменных окружения

Чтобы не вводить каждый раз длинную строку подключения, настройте переменные окружения в PowerShell:

```powershell
$env:GOOSE_DRIVER="postgres"
$env:GOOSE_DBSTRING="postgresql://rally_user:rally_dev_password@127.0.0.1:5433/rally_dev?sslmode=disable"
$env:GOOSE_MIGRATION_DIR="internal/database/migrations"
```

После этого команды станут проще:

```powershell
goose up
goose status
goose down
```

**Для постоянного использования** добавьте эти переменные в ваш PowerShell профиль:

```powershell
# Откройте профиль
notepad $PROFILE

# Добавьте эти строки:
$env:GOOSE_DRIVER="postgres"
$env:GOOSE_DBSTRING="postgresql://rally_user:rally_dev_password@127.0.0.1:5433/rally_dev?sslmode=disable"
$env:GOOSE_MIGRATION_DIR="internal/database/migrations"
```

**ВАЖНО:** Для этого проекта используйте скрипт `migrate.cmd` для упрощения работы с миграциями.

### Common Migration Commands

#### Apply All Migrations (Up)

**Рекомендуется:** Используйте скрипт `migrate.cmd`:

```powershell
.\migrate.cmd up
```

Или полная команда:

```powershell
goose -dir internal/database/migrations postgres "postgresql://rally_user:rally_dev_password@127.0.0.1:5433/rally_dev?sslmode=disable" up
```

#### Check Migration Status

```powershell
.\migrate.cmd status
```

#### Rollback Last Migration (Down)

```powershell
.\migrate.cmd down
```

#### Rollback to Specific Version

```powershell
goose -dir internal/database/migrations postgres "postgresql://rally_user:rally_dev_password@localhost:5432/rally_dev?sslmode=disable" down-to VERSION
```

Replace `VERSION` with the target migration version (e.g., `1`).

#### Reset Database (Down to Zero, then Up)

```powershell
goose -dir internal/database/migrations postgres "postgresql://rally_user:rally_dev_password@localhost:5432/rally_dev?sslmode=disable" reset
```

#### Create New Migration

```powershell
goose -dir internal/database/migrations create migration_name sql
```

Replace `migration_name` with a descriptive name (e.g., `create_events_table`).

This will create two files:
- `NNNNNN_migration_name.sql` (where NNNNNN is a timestamp)

Edit the file to add your SQL with goose annotations:

```sql
-- +goose Up
-- +goose StatementBegin
-- Your UP migration SQL here
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Your DOWN migration SQL here
-- +goose StatementEnd
```

## Development Workflow

### First Time Setup

1. **Start Docker services:**

   ```powershell
   docker-compose up -d postgres
   ```

2. **Install dependencies:**

   ```powershell
   go mod download
   ```

3. **Run migrations:**

   ```powershell
   .\migrate.cmd up
   ```

4. **Run the application:**

   ```powershell
   go run cmd/server/main.go
   ```

### Daily Development

1. **Start services (if not running):**

   ```powershell
   docker-compose up -d postgres
   ```

2. **Run application:**

   ```powershell
   go run cmd/server/main.go
   ```

### Running Tests

```powershell
go test -v -race ./...
```

For unit tests only:

```powershell
go test -v -race -tags=unit ./...
```

For integration tests:

```powershell
go test -v -race -tags=integration ./...
```

### Building the Application

```powershell
go build -o bin/sdd-rally-app.exe cmd/server/main.go
```

Run the built binary:

```powershell
.\bin\sdd-rally-app.exe
```

## Docker Commands

### Start All Services

```powershell
docker-compose up -d
```

### Start Only Database

```powershell
docker-compose up -d postgres
```

### View Logs

```powershell
docker-compose logs -f
```

Or for specific service:

```powershell
docker-compose logs -f postgres
```

### Stop Services

```powershell
docker-compose down
```

### Stop and Remove Volumes (CAUTION: Deletes Data)

```powershell
docker-compose down -v
```

### Rebuild Services

```powershell
docker-compose up -d --build
```

## Database Access

### Connect to PostgreSQL via Docker

```powershell
docker exec -it sdd-rally-postgres psql -U rally_user -d rally_dev
```

Common PostgreSQL commands inside psql:

- `\dt` - List tables
- `\d table_name` - Describe table structure
- `\du` - List users
- `\l` - List databases
- `\q` - Quit

### Direct Connection (with psql installed locally)

```powershell
psql "postgresql://rally_user:rally_dev_password@127.0.0.1:5433/rally_dev?sslmode=disable"
```

## Troubleshooting

### Port 5432 Already in Use

If you have PostgreSQL running locally, stop it or change the port in `docker-compose.yml`:

```yaml
ports:
  - "5433:5432"  # Use port 5433 on host
```

Then update connection strings to use `localhost:5433`.

### Goose Command Not Found

Make sure `%GOPATH%\bin` is in your PATH. Check with:

```powershell
$env:Path
```

Add to PATH if needed:

```powershell
$env:Path += ";$env:USERPROFILE\go\bin"
```

### Docker Container Won't Start

Check logs:

```powershell
docker logs sdd-rally-postgres
```

Remove container and try again:

```powershell
docker-compose down -v
docker-compose up -d postgres
```

### Migration Fails

Check migration status:

```powershell
goose -dir internal/database/migrations postgres $env:GOOSE_DBSTRING status
```

Roll back if needed:

```powershell
goose -dir internal/database/migrations postgres $env:GOOSE_DBSTRING down
```

## Project Structure

```
sdd-rally-app/
├── cmd/
│   └── server/          # Application entry point
├── internal/
│   ├── config/          # Configuration management
│   ├── database/        # Database connection
│   │   └── migrations/  # SQL migration files
│   └── ...
├── web/
│   ├── templates/       # HTML templates
│   └── static/          # CSS, JS, images
├── docker-compose.yml   # Docker services definition
├── go.mod              # Go dependencies
└── SETUP.md            # This file
```

## Additional Tools (Optional)

### Air (Hot Reload)

For automatic reloading during development:

```powershell
go install github.com/cosmtrek/air@latest
air
```

### golangci-lint (Linting)

```powershell
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
golangci-lint run ./...
```

### templ (Template Engine)

If using templ for Go templates:

```powershell
go install github.com/a-h/templ/cmd/templ@latest
```

## Support

For issues or questions, check:
- Project documentation in `specs/` directory
- Go documentation: `go doc <package>`
- Goose documentation: https://github.com/pressly/goose

