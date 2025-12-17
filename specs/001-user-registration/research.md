# Исследование: Регистрация и аутентификация пользователей

**Feature**: 001-user-registration  
**Дата**: 2025-12-17  
**Статус**: Завершено

## Цель исследования

Определить оптимальные технологии, библиотеки и паттерны для реализации системы регистрации и аутентификации, соответствующей требованиям конституции проекта.

## Технологические решения

### 1. HTTP фреймворк

**Решение**: Fiber v2

**Обоснование**:
- Высокая производительность (основан на fasthttp, один из самых быстрых Go фреймворков)
- Express-подобный API (знакомый многим разработчикам)
- Встроенная поддержка middleware
- Отличная документация на русском языке
- Активное сообщество и регулярные обновления
- Легковесный (минимальные накладные расходы)

**Альтернативы рассмотрены**:
- **Gin**: Хорошая производительность, но Fiber быстрее по бенчмаркам
- **Echo**: Похожий функционал, но меньшее сообщество
- **Chi**: Более минималистичный, потребует больше кода для middleware
- **net/http (стандартная библиотека)**: Надёжно, но требует больше boilerplate кода

**Ссылки**:
- https://gofiber.io/
- https://github.com/gofiber/fiber

---

### 2. PostgreSQL драйвер

**Решение**: pgx v5

**Обоснование**:
- Нативный PostgreSQL драйвер (не использует database/sql)
- Лучшая производительность среди Go драйверов
- Поддержка prepared statements и batch операций
- Отличная поддержка PostgreSQL-специфичных типов
- Встроенная поддержка connection pooling
- Активная разработка и поддержка

**Альтернативы рассмотрены**:
- **pq**: Старый стандарт, но pgx производительнее
- **GORM**: ORM добавляет абстракцию, конституция требует прямого контроля над SQL
- **sqlx**: Хороший выбор, но pgx даёт больше контроля

**Ссылки**:
- https://github.com/jackc/pgx

---

### 3. Шаблонизатор

**Решение**: Templ

**Обоснование**:
- Типобезопасные шаблоны (компилируются в Go код)
- Автоматическое экранирование HTML (защита от XSS)
- Отличная интеграция с Go (автодополнение в IDE)
- Хорошая производительность (компиляция в Go код)
- Простой синтаксис
- Растущая популярность в Go экосистеме

**Альтернативы рассмотрены**:
- **html/template (стандартная библиотека)**: Медленнее, нет типобезопасности
- **Jet**: Хорошая производительность, но меньше community support
- **Pongo2**: Django-подобный синтаксис, но не типобезопасный

**Ссылки**:
- https://templ.guide/
- https://github.com/a-h/templ

---

### 4. Миграции базы данных

**Решение**: golang-migrate

**Обоснование**:
- Стандарт де-факто для миграций в Go
- Поддержка множества БД (PostgreSQL, MySQL и др.)
- CLI инструмент для управления миграциями
- Программный API для встраивания в приложение
- Поддержка up/down миграций
- Версионирование миграций

**Альтернативы рассмотрены**:
- **goose**: Хорошая альтернатива, но golang-migrate популярнее
- **sql-migrate**: Меньше функций
- **Atlas**: Более современный, но слишком сложный для наших нужд

**Ссылки**:
- https://github.com/golang-migrate/migrate

---

### 5. Логирование

**Решение**: zerolog

**Обоснование**:
- Высокая производительность (zero allocation)
- Структурированное логирование (JSON)
- Простой и понятный API
- Поддержка уровней логирования (debug, info, warn, error)
- Контекстное логирование
- Минимальные накладные расходы

**Альтернативы рассмотрены**:
- **zap (Uber)**: Очень быстрый, но более сложный API
- **logrus**: Популярный, но медленнее zerolog
- **slog (Go 1.21+)**: Часть стандартной библиотеки, но zerolog производительнее

**Ссылки**:
- https://github.com/rs/zerolog

---

### 6. Хеширование паролей

**Решение**: golang.org/x/crypto/bcrypt

**Обоснование**:
- Стандартная библиотека Go (часть golang.org/x)
- Надёжный, проверенный временем алгоритм
- Настраиваемый cost factor (будем использовать 12)
- Защита от rainbow table атак
- Автоматическое добавление salt

**Альтернативы рассмотрены**:
- **argon2**: Современнее, устойчивее к GPU атакам, но bcrypt проще и достаточен для наших нужд
- **scrypt**: Хорошо, но bcrypt стандарт для веб-приложений

**Конфигурация**: Cost factor = 12 (баланс между безопасностью и производительностью)

**Ссылки**:
- https://pkg.go.dev/golang.org/x/crypto/bcrypt

---

### 7. Валидация данных

**Решение**: validator v10 (go-playground/validator)

**Обоснование**:
- Декларативная валидация через struct tags
- Богатый набор встроенных валидаторов
- Поддержка кастомных валидаторов
- Хорошая производительность
- Интеграция с Fiber через middleware

**Альтернативы рассмотрены**:
- **ozzo-validation**: Программная валидация, но tags удобнее
- **govalidator**: Устарел, меньше функций

**Ссылки**:
- https://github.com/go-playground/validator

---

### 8. CSRF защита

**Решение**: fiber/middleware/csrf

**Обоснование**:
- Встроенный middleware Fiber
- Double-submit cookie pattern
- Поддержка кастомизации (cookie name, expiration)
- Простая интеграция

**Конфигурация**:
- Cookie name: `csrf_`
- Expiration: 1 hour
- SameSite: Lax

**Ссылки**:
- https://docs.gofiber.io/api/middleware/csrf

---

### 9. Rate Limiting

**Решение**: Кастомный middleware + in-memory store (sync.Map)

**Обоснование**:
- Простая реализация для MVP
- Достаточно для начального масштаба (< 1000 пользователей)
- Можно мигрировать на Redis позже
- Полный контроль над логикой

**Альтернативы рассмотрены**:
- **fiber/middleware/limiter**: Слишком простой, нужна более сложная логика (разные лимиты для разных эндпоинтов)
- **Redis-based**: Overkill для MVP, добавит зависимость

**Алгоритм**: Token bucket с разными бакетами для:
- Регистрация: 3 попытки / 15 минут на IP
- Вход: 5 попыток / 15 минут на IP
- После 10 попыток: блокировка на 30 минут

---

### 10. HTMX для интерактивности

**Решение**: HTMX 1.9+

**Обоснование**:
- Легковесная библиотека (14KB gzipped)
- Декларативный подход (HTML атрибуты)
- Прогрессивное улучшение (формы работают без JS)
- Поддержка partial updates
- Хорошая документация

**Паттерны использования**:
- `hx-post` для отправки форм
- `hx-target` для обновления конкретных элементов (ошибки валидации)
- `hx-swap` для контроля замены контента
- `hx-indicator` для loading состояния

**Ссылки**:
- https://htmx.org/

---

### 11. Session storage

**Решение**: Database (PostgreSQL) для MVP

**Обоснование**:
- Простота реализации (одна БД для всего)
- ACID гарантии для сессий
- Возможность сложных запросов (активные сессии пользователя)
- Нет дополнительных зависимостей
- Производительность достаточна для < 1000 пользователей

**План миграции**: При росте нагрузки (> 5000 пользователей) мигрировать на Redis:
- Redis для активных сессий (TTL)
- PostgreSQL для истории (audit trail)

**Альтернативы рассмотрены**:
- **Redis**: Быстрее, но добавляет зависимость и сложность на MVP этапе
- **In-memory**: Не переживёт рестарт приложения

**Структура session**:
```
sessions таблица:
- id (UUID, PK)
- user_id (UUID, FK to users)
- token (хеш сессии для cookie)
- created_at
- last_activity_at
- ip_address
- user_agent
- expires_at (индекс для автоочистки)
```

---

### 12. Captcha для rate limiting

**Решение**: Простая математическая captcha (для MVP)

**Обоснование**:
- Нет внешних зависимостей
- Достаточно для защиты от простых ботов
- Быстрая реализация
- Работает без JavaScript

**План улучшения**: В будущем добавить hCaptcha или Cloudflare Turnstile (privacy-friendly)

**Реализация**:
- Генерация простого примера (5 + 3 = ?)
- Сохранение ответа в сессии
- Проверка при сабмите формы

**Альтернативы рассмотрены**:
- **reCAPTCHA**: Privacy concerns, внешняя зависимость
- **hCaptcha**: Хорошо, но overkill для MVP
- **Image captcha**: Сложнее реализовать, хуже UX

---

## Паттерны и best practices

### Структура handler'ов

```go
// Паттерн: тонкие handlers, толстые services
func (h *AuthHandler) Register(c *fiber.Ctx) error {
    // 1. Парсинг и валидация входных данных
    var req RegisterRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(ErrorResponse{...})
    }
    
    if err := validate.Struct(req); err != nil {
        return c.Status(400).JSON(ValidationErrors(err))
    }
    
    // 2. Вызов service (вся бизнес-логика там)
    user, err := h.authService.Register(c.Context(), req)
    if err != nil {
        return handleServiceError(c, err)
    }
    
    // 3. Возврат ответа
    return c.Status(201).JSON(user)
}
```

### Структура service'ов

```go
// Паттерн: использование context, возврат доменных ошибок
func (s *AuthService) Register(ctx context.Context, req RegisterRequest) (*User, error) {
    // 1. Проверка бизнес-правил
    exists, err := s.userRepo.ExistsByUsername(ctx, req.Username)
    if err != nil {
        return nil, err
    }
    if exists {
        return nil, ErrUsernameTaken
    }
    
    // 2. Хеширование пароля
    hashedPassword, err := s.passwordService.Hash(req.Password)
    if err != nil {
        return nil, err
    }
    
    // 3. Создание пользователя
    user := &User{
        ID:       uuid.New(),
        Username: strings.ToLower(strings.TrimSpace(req.Username)),
        Password: hashedPassword,
        Role:     s.determineInitialRole(), // Первый = chief_organizer
        Status:   StatusActive,
    }
    
    // 4. Сохранение в БД (с transaction если нужно)
    if err := s.userRepo.Create(ctx, user); err != nil {
        return nil, err
    }
    
    // 5. Логирование
    s.logger.Info().
        Str("user_id", user.ID.String()).
        Str("username", user.Username).
        Msg("User registered successfully")
    
    return user, nil
}
```

### Обработка ошибок

```go
// Паттерн: доменные ошибки + HTTP mapping
var (
    ErrUsernameTaken = errors.New("username already taken")
    ErrInvalidCredentials = errors.New("invalid credentials")
    ErrAccountLocked = errors.New("account locked")
)

func handleServiceError(c *fiber.Ctx, err error) error {
    switch {
    case errors.Is(err, ErrUsernameTaken):
        return c.Status(409).JSON(ErrorResponse{
            Code: "USERNAME_TAKEN",
            Message: "Пользователь с таким логином уже существует",
        })
    case errors.Is(err, ErrInvalidCredentials):
        return c.Status(401).JSON(ErrorResponse{
            Code: "INVALID_CREDENTIALS",
            Message: "Неверный логин или пароль",
        })
    case errors.Is(err, ErrAccountLocked):
        return c.Status(423).JSON(ErrorResponse{
            Code: "ACCOUNT_LOCKED",
            Message: "Аккаунт временно заблокирован. Попробуйте через 30 минут",
        })
    default:
        return c.Status(500).JSON(ErrorResponse{
            Code: "INTERNAL_ERROR",
            Message: "Внутренняя ошибка сервера",
        })
    }
}
```

### Middleware chain

```go
// Порядок middleware (важно!)
app := fiber.New()

// 1. Логирование (первым, чтобы залогировать все запросы)
app.Use(middleware.Logging())

// 2. Security headers
app.Use(middleware.SecurityHeaders())

// 3. CORS
app.Use(middleware.CORS())

// 4. Rate limiting (до auth, чтобы блокировать брутфорс)
app.Use(middleware.RateLimit())

// 5. CSRF (только для форм, после auth)
app.Use(middleware.CSRF())

// 6. Auth (проверка сессии)
// Применяется селективно к защищённым роутам

// Пример защищённого роута
app.Get("/dashboard", middleware.RequireAuth(), dashboardHandler)
app.Get("/admin", middleware.RequireRole("chief_organizer"), adminHandler)
```

---

## Метрики производительности

### Целевые показатели

| Операция | Целевое время | Обоснование |
|----------|--------------|-------------|
| Hash пароля (bcrypt cost 12) | ~250ms | Приемлемо для регистрации/входа |
| Запрос к БД (SELECT user) | <10ms | SSD диск, индексы |
| Создание сессии | <20ms | INSERT + cookie |
| Полный flow регистрации | <500ms | Hash + DB insert + session |
| Полный flow входа | <300ms | DB select + bcrypt compare + session |

### Оптимизации

1. **Индексы БД**:
   - `users.username` (UNIQUE, lower()) - быстрый поиск при входе
   - `sessions.token` (UNIQUE) - быстрая проверка сессии
   - `sessions.user_id` - быстрый поиск всех сессий пользователя
   - `login_attempts.ip_address, created_at` - rate limiting queries

2. **Connection pooling**:
   - Max connections: 25
   - Max idle connections: 5
   - Connection lifetime: 5 minutes

3. **Кеширование** (для будущего):
   - User + Role в Redis на 5 минут
   - Permissions в memory cache

---

## Риски и митигация

| Риск | Вероятность | Влияние | Митигация |
|------|-------------|---------|-----------|
| Брутфорс атака на пароли | Высокая | Критичное | Rate limiting + captcha + account locking |
| Session fixation | Средняя | Критичное | Regenerate session ID после входа |
| Session hijacking | Средняя | Критичное | HTTP-only, Secure cookies + IP validation |
| DoS через создание сессий | Средняя | Высокое | Rate limiting на регистрацию + cleanup старых сессий |
| SQL injection | Низкая | Критичное | Параметризованные запросы (pgx) |
| XSS | Низкая | Высокое | Templ автоматически экранирует + CSP headers |
| CSRF | Средняя | Высокое | CSRF middleware с токенами |
| Утечка паролей в логах | Низкая | Критичное | Никогда не логировать пароли, только ID пользователя |

---

## Зависимости Go modules

```
require (
    github.com/gofiber/fiber/v2 v2.51.0
    github.com/jackc/pgx/v5 v5.5.0
    github.com/a-h/templ v0.2.543
    github.com/golang-migrate/migrate/v4 v4.17.0
    github.com/rs/zerolog v1.31.0
    github.com/go-playground/validator/v10 v10.16.0
    github.com/google/uuid v1.5.0
    golang.org/x/crypto v0.17.0
)
```

---

## Следующие шаги

1. ✅ Исследование завершено
2. → Создать data-model.md (детальная схема БД)
3. → Создать contracts/ (OpenAPI спецификация)
4. → Создать quickstart.md (инструкции для разработчиков)
5. → Перейти к /speckit.tasks (разбивка на задачи)

**Статус**: Все технические решения утверждены и готовы к реализации.

