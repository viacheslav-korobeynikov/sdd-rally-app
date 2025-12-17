# Быстрый старт: Регистрация и аутентификация

**Feature**: 001-user-registration  
**Дата**: 2025-12-17

## Цель

Эта инструкция поможет быстро настроить локальное окружение и начать разработку функции регистрации и аутентификации.

## Предварительные требования

### Обязательно

- Go 1.21+ ([установка](https://go.dev/doc/install))
- PostgreSQL 15+ ([установка](https://www.postgresql.org/download/))
- Docker и Docker Compose ([установка](https://docs.docker.com/get-docker/))
- Git
- Make

### Рекомендуется

- Go IDE (VSCode с Go extension, GoLand)
- Postman или Insomnia (для тестирования API)
- psql (PostgreSQL CLI client)

## Шаг 1: Клонирование и настройка

```bash
# Клонировать репозиторий (если ещё не сделано)
git clone https://github.com/your-org/sdd-rally-app.git
cd sdd-rally-app

# Переключиться на feature ветку
git checkout 1-user-registration

# Создать локальную конфигурацию
cp .env.example .env.local

# Отредактировать .env.local
# (установить пароли, порты и т.д.)
```

### Пример .env.local

```env
# Окружение
ENV=development

# Сервер
PORT=3000
HOST=localhost

# База данных
DB_HOST=localhost
DB_PORT=5432
DB_NAME=rally_dev
DB_USER=rally_user
DB_PASSWORD=your_secure_password_here
DB_SSL_MODE=disable

# Безопасность
SESSION_SECRET=your_random_32_char_secret_here
CSRF_KEY=your_random_32_char_csrf_key_here
BCRYPT_COST=12

# CORS
CORS_ALLOWED_ORIGINS=http://localhost:3000

# Логирование
LOG_LEVEL=debug
LOG_FORMAT=console

# Rate Limiting
RATE_LIMIT_ENABLED=true
RATE_LIMIT_REQUESTS_PER_MINUTE=100
RATE_LIMIT_LOGIN_ATTEMPTS=5
```

## Шаг 2: Установка зависимостей

### Через Docker Compose (рекомендуется)

```bash
# Запустить все сервисы
make docker-up

# Это запустит:
# - PostgreSQL 15 (порт 5432)
# - Backend приложение (порт 3000)
```

### Локально (без Docker)

```bash
# Установить Go зависимости
make install-deps

# Установить инструменты разработки
make install-tools

# Запустить PostgreSQL (если не используете Docker)
# На Linux:
sudo systemctl start postgresql

# На macOS:
brew services start postgresql

# На Windows:
# Запустить PostgreSQL через pgAdmin или services
```

## Шаг 3: Инициализация базы данных

```bash
# Создать базу данных
createdb rally_dev

# Или через psql:
psql -U postgres -c "CREATE DATABASE rally_dev;"
psql -U postgres -c "CREATE USER rally_user WITH PASSWORD 'your_password';"
psql -U postgres -c "GRANT ALL PRIVILEGES ON DATABASE rally_dev TO rally_user;"

# Применить миграции
make migrate-up

# Проверить статус миграций
migrate -path internal/database/migrations -database "postgresql://rally_user:your_password@localhost:5432/rally_dev?sslmode=disable" version
```

## Шаг 4: Загрузка тестовых данных (опционально)

```bash
# Загрузить seed данные
make seed

# Это создаст:
# - Тестового пользователя: admin / Admin123456
# - Несколько примеров пользователей с разными ролями
```

## Шаг 5: Запуск приложения

### Вариант 1: С hot reload (разработка)

```bash
# Запустить с автоперезагрузкой при изменении кода
make dev

# Приложение запустится на http://localhost:3000
# При изменении .go файлов будет автоматически перезапускаться
```

### Вариант 2: Обычный запуск

```bash
# Запустить без hot reload
make run
```

### Вариант 3: Через Docker Compose

```bash
# Всё уже запущено из Шага 2
# Логи можно посмотреть через:
make docker-logs
```

## Шаг 6: Проверка работоспособности

### Через curl

```bash
# Health check
curl http://localhost:3000/health

# Ожидаемый ответ:
# {"status":"ok","timestamp":"2025-12-17T10:30:00Z"}

# Регистрация пользователя
curl -X POST http://localhost:3000/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "TestPassword123",
    "password_confirm": "TestPassword123"
  }'

# Вход в систему
curl -X POST http://localhost:3000/api/auth/login \
  -H "Content-Type: application/json" \
  -c cookies.txt \
  -d '{
    "username": "testuser",
    "password": "TestPassword123"
  }'

# Получить текущего пользователя
curl http://localhost:3000/api/auth/me \
  -b cookies.txt
```

### Через браузер

1. Откройте http://localhost:3000
2. Перейдите на страницу регистрации: http://localhost:3000/register
3. Заполните форму:
   - Логин: `testuser`
   - Пароль: `TestPassword123`
   - Подтверждение пароля: `TestPassword123`
4. Нажмите "Зарегистрироваться"
5. Вы должны быть автоматически авторизованы и перенаправлены на главную страницу

### Через Postman

1. Импортируйте коллекцию: `specs/001-user-registration/contracts/auth-api.yaml`
2. Выберите environment "Local"
3. Запустите request "Register"
4. Проверьте, что cookie `session_id` установлен
5. Запустите request "Get Current User"

## Шаг 7: Запуск тестов

```bash
# Все тесты
make test

# Только unit тесты
make test-unit

# Только integration тесты
make test-integration

# С coverage
make test-coverage

# Открыть coverage report в браузере
open coverage.html
```

## Структура проекта

```
sdd-rally-app/
├── cmd/
│   └── server/
│       └── main.go              # Точка входа
├── internal/
│   ├── auth/                    # Модуль аутентификации
│   │   ├── handlers/            # HTTP handlers
│   │   ├── services/            # Бизнес-логика
│   │   ├── repositories/        # Работа с БД
│   │   └── models/              # Доменные модели
│   ├── middleware/              # HTTP middleware
│   ├── database/                # Подключение к БД
│   │   └── migrations/          # SQL миграции
│   ├── shared/                  # Общие утилиты
│   └── config/                  # Конфигурация
├── web/
│   ├── templates/               # Templ шаблоны
│   └── static/                  # CSS, JS, изображения
├── tests/
│   ├── unit/                    # Unit тесты
│   ├── integration/             # Integration тесты
│   └── fixtures/                # Тестовые данные
├── specs/
│   └── 001-user-registration/   # Документация feature
├── docs/                        # Общая документация
├── .env.example                 # Пример конфигурации
├── docker-compose.yml           # Docker окружение
├── Makefile                     # Команды для разработки
├── go.mod                       # Go модуль
└── README.md                    # Основной README
```

## Полезные команды

### Разработка

```bash
# Запустить с hot reload
make dev

# Собрать бинарник
make build

# Запустить линтеры
make lint

# Форматировать код
make fmt

# Обновить зависимости
make mod-tidy
```

### База данных

```bash
# Применить миграции
make migrate-up

# Откатить последнюю миграцию
make migrate-down

# Создать новую миграцию
make migrate-create name=add_something

# Загрузить seed данные
make seed

# Подключиться к БД через psql
psql -U rally_user -d rally_dev
```

### Docker

```bash
# Запустить все сервисы
make docker-up

# Остановить все сервисы
make docker-down

# Показать логи
make docker-logs

# Пересобрать образы
make docker-rebuild
```

### Тестирование

```bash
# Все тесты
make test

# Unit тесты
make test-unit

# Integration тесты
make test-integration

# Coverage
make test-coverage

# Security check
make security-check
```

### Очистка

```bash
# Очистить build артефакты
make clean

# Очистить всё (включая зависимости)
make clean-all
```

## Troubleshooting

### Проблема: Не удаётся подключиться к БД

```bash
# Проверить, что PostgreSQL запущен
pg_isready -h localhost -p 5432

# Проверить права доступа
psql -U rally_user -d rally_dev -c "SELECT 1;"

# Проверить настройки в .env.local
cat .env.local | grep DB_
```

### Проблема: Миграции не применяются

```bash
# Проверить текущую версию БД
migrate -path internal/database/migrations -database "your_db_url" version

# Форсировать версию (ОСТОРОЖНО!)
migrate -path internal/database/migrations -database "your_db_url" force <version>

# Откатить всё и применить заново
make migrate-down
make migrate-up
```

### Проблема: Port 3000 уже занят

```bash
# Найти процесс, занимающий порт
lsof -i :3000  # Linux/macOS
netstat -ano | findstr :3000  # Windows

# Убить процесс или изменить PORT в .env.local
```

### Проблема: Hot reload не работает

```bash
# Проверить, что air установлен
which air

# Переустановить air
go install github.com/cosmtrek/air@latest

# Проверить .air.toml конфигурацию
cat .air.toml
```

### Проблема: Тесты падают

```bash
# Запустить тесты с verbose
go test -v ./...

# Проверить тестовую БД
make test-db-setup

# Очистить тестовый кеш
go clean -testcache
```

## Следующие шаги

После успешного запуска локального окружения:

1. 📖 Изучите [spec.md](./spec.md) - требования feature
2. 📊 Изучите [data-model.md](./data-model.md) - схема БД
3. 🔌 Изучите [contracts/](./contracts/) - API спецификация
4. ✅ Перейдите к [tasks.md](./tasks.md) - список задач для реализации (создаётся командой `/speckit.tasks`)

## Полезные ресурсы

### Документация

- [Fiber документация](https://docs.gofiber.io/)
- [pgx документация](https://pkg.go.dev/github.com/jackc/pgx/v5)
- [Templ документация](https://templ.guide/)
- [HTMX документация](https://htmx.org/docs/)
- [golang-migrate](https://github.com/golang-migrate/migrate)

### Внутренние ресурсы

- [Конституция проекта](../../.specify/memory/constitution.md)
- [Каталог модулей](../../docs/architecture/module-catalog.md)
- [Руководство по тестированию](../../docs/testing/guidelines.md)
- [Стандарты безопасности](../../docs/security/standards.md)
- [Наблюдаемость](../../docs/operations/observability.md)

## Помощь

Если возникли проблемы:

1. Проверьте [Troubleshooting](#troubleshooting) выше
2. Проверьте логи: `make docker-logs` или `tail -f logs/app.log`
3. Создайте issue в GitHub с описанием проблемы
4. Обратитесь в команду разработки

---

**Статус**: Готово к разработке  
**Следующий шаг**: Запустите `/speckit.tasks` для создания детального списка задач

