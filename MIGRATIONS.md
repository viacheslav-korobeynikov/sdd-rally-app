# Работа с миграциями БД

## Быстрая справка

### Применить все миграции
```powershell
.\migrate.cmd up
```

### Проверить статус миграций
```powershell
.\migrate.cmd status
```

### Откатить последнюю миграцию
```powershell
.\migrate.cmd down
```

### Создать новую миграцию
```powershell
goose -dir internal/database/migrations create migration_name sql
```

## Подключение к БД

### Через Docker
```powershell
docker exec -it sdd-rally-postgres psql -U rally_user -d rally_dev
```

### Через psql (если установлен локально)
```powershell
psql -h 127.0.0.1 -p 5433 -U rally_user -d rally_dev
```
Пароль: `rally_dev_password`

## Полезные SQL команды в psql

- `\dt` - показать все таблицы
- `\d table_name` - описание таблицы
- `\dT` - показать все типы (enums)
- `\du` - список пользователей
- `\q` - выйти

## Настройки подключения

- **Host:** 127.0.0.1
- **Port:** 5433 (не 5432!)
- **Database:** rally_dev
- **User:** rally_user
- **Password:** rally_dev_password

**Примечание:** Порт 5433 используется, чтобы избежать конфликта с локальным PostgreSQL.

## Структура миграций

Миграции находятся в `internal/database/migrations/` и используют формат goose:

```sql
-- +goose Up
-- +goose StatementBegin
-- Ваш SQL код для применения миграции
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Ваш SQL код для отката миграции
-- +goose StatementEnd
```

## Troubleshooting

### Ошибка подключения
1. Проверьте, что контейнер запущен: `docker ps`
2. Проверьте порт в `.env` файле (должен быть 5433)
3. Убедитесь, что не подключаетесь к локальному PostgreSQL

### Контейнер не запускается
```powershell
docker-compose down -v
docker-compose up -d postgres
```

Подробнее см. в [SETUP.md](SETUP.md)

