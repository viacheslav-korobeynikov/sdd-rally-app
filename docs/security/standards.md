# Стандарты безопасности

**Статус**: Черновик  
**Версия**: 0.1.0  
**Последнее обновление**: 2025-12-17

## Обзор

Этот документ определяет требования безопасности для SDD Rally App, включая управление паролями, ролями, секретами и персональными данными участников.

## Аутентификация

### Хранение паролей

**Алгоритм хеширования**: bcrypt (cost 12) или argon2id

**Обоснование**: 
- bcrypt - проверенный временем, медленный (защита от брутфорса)
- argon2id - современный, устойчив к GPU/ASIC атакам

**Реализация (bcrypt)**:
```go
import "golang.org/x/crypto/bcrypt"

// Хеширование при регистрации
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
    return string(bytes), err
}

// Проверка при входе
func CheckPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
```

### Политика паролей

**Требования**:
- Минимальная длина: 12 символов
- Обязательно: минимум 1 заглавная буква, 1 строчная, 1 цифра
- Рекомендовано: использование специальных символов
- Запрещено: пароли из списка распространённых (top 10000)

**Проверка надёжности**:
```go
func ValidatePasswordStrength(password string) error {
    if len(password) < 12 {
        return errors.New("пароль должен содержать минимум 12 символов")
    }
    
    var hasUpper, hasLower, hasDigit bool
    for _, char := range password {
        switch {
        case unicode.IsUpper(char):
            hasUpper = true
        case unicode.IsLower(char):
            hasLower = true
        case unicode.IsDigit(char):
            hasDigit = true
        }
    }
    
    if !hasUpper || !hasLower || !hasDigit {
        return errors.New("пароль должен содержать заглавные и строчные буквы, а также цифры")
    }
    
    return nil
}
```

### Управление сессиями

**Механизм**: HTTP-only, Secure cookies + server-side session storage

**Параметры cookie**:
```go
cookie := &fiber.Cookie{
    Name:     "session_id",
    Value:    sessionID,
    HTTPOnly: true,  // Защита от XSS
    Secure:   true,  // Только HTTPS (в продакшене)
    SameSite: "Lax", // Защита от CSRF
    MaxAge:   3600,  // 1 час
}
```

**Время жизни сессии**:
- Активная сессия: 1 час
- "Запомнить меня": 30 дней (с пометкой для критичных операций)
- Автоматический logout при неактивности: 1 час

**Хранение сессий**: PostgreSQL (таблица sessions) или Redis

### Защита от брутфорса

**Rate limiting на login**:
- Максимум 5 попыток в 15 минут с одного IP
- После 5 неудачных попыток - captcha
- После 10 неудачных попыток - блокировка на 30 минут

**Логирование**:
```go
log.Warn().
    Str("username", username).
    Str("ip", clientIP).
    Int("attempt", attemptCount).
    Msg("Failed login attempt")
```

## Авторизация (RBAC)

### Роли системы

| Роль | Описание | Права |
|------|----------|-------|
| **Главный организатор** | Полный доступ к соревнованию | Все операции |
| **Секретарь** | Управление экипажами и заявками | Создание/редактирование экипажей, просмотр результатов |
| **Хронометраж** | Ввод результатов | Ввод времени старта/финиша, просмотр результатов |
| **Наблюдатель** | Только чтение | Просмотр всех данных без изменений |

### Детальные права (permissions)

```go
type Permission string

const (
    // Competitions
    CompetitionCreate Permission = "competition:create"
    CompetitionRead   Permission = "competition:read"
    CompetitionUpdate Permission = "competition:update"
    CompetitionDelete Permission = "competition:delete"
    
    // Crews
    CrewCreate Permission = "crew:create"
    CrewRead   Permission = "crew:read"
    CrewUpdate Permission = "crew:update"
    CrewDelete Permission = "crew:delete"
    
    // Timing
    TimingRecord Permission = "timing:record"
    TimingRead   Permission = "timing:read"
    
    // Results
    ResultsRead      Permission = "results:read"
    ResultsCalculate Permission = "results:calculate"
    PenaltyApply     Permission = "penalty:apply"
)

// Права по ролям
var RolePermissions = map[Role][]Permission{
    RoleChiefOrganizer: {
        CompetitionCreate, CompetitionRead, CompetitionUpdate, CompetitionDelete,
        CrewCreate, CrewRead, CrewUpdate, CrewDelete,
        TimingRecord, TimingRead,
        ResultsRead, ResultsCalculate, PenaltyApply,
    },
    RoleSecretary: {
        CompetitionRead,
        CrewCreate, CrewRead, CrewUpdate,
        ResultsRead,
    },
    RoleTiming: {
        CompetitionRead,
        CrewRead,
        TimingRecord, TimingRead,
        ResultsRead,
    },
    RoleObserver: {
        CompetitionRead,
        CrewRead,
        TimingRead,
        ResultsRead,
    },
}
```

### Проверка прав (Middleware)

```go
func RequirePermission(perm Permission) fiber.Handler {
    return func(c *fiber.Ctx) error {
        userID := c.Locals("user_id").(string)
        
        user, err := getUserByID(c.Context(), userID)
        if err != nil {
            return fiber.ErrUnauthorized
        }
        
        if !user.HasPermission(perm) {
            log.Warn().
                Str("user_id", userID).
                Str("permission", string(perm)).
                Str("path", c.Path()).
                Msg("Permission denied")
            return fiber.ErrForbidden
        }
        
        return c.Next()
    }
}

// Использование
app.Post("/api/competitions", 
    RequireAuth(),
    RequirePermission(CompetitionCreate),
    handlers.CreateCompetition,
)
```

### Аудит изменения ролей

**Таблица role_assignments_audit**:
```sql
CREATE TABLE role_assignments_audit (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    competition_id UUID,  -- NULL для глобальных ролей
    old_role VARCHAR(50),
    new_role VARCHAR(50) NOT NULL,
    assigned_by UUID NOT NULL,  -- Кто назначил
    assigned_at TIMESTAMP NOT NULL DEFAULT NOW(),
    reason TEXT
);
```

## Управление секретами

### Запрещено хранить в репозитории

❌ **Никогда не коммитить**:
- Пароли баз данных
- Ключи шифрования
- API ключи внешних сервисов
- Session secrets
- JWT secrets
- Приватные ключи (SSH, TLS)

### Использование переменных окружения

**Development** (`.env.local` в .gitignore):
```env
DB_HOST=localhost
DB_PORT=5432
DB_NAME=rally_dev
DB_USER=rally_user
DB_PASSWORD=dev_password_here

SESSION_SECRET=dev_secret_key_32_chars_min

# API keys
EXTERNAL_API_KEY=test_key
```

**Production** (Kubernetes Secrets, AWS Secrets Manager, HashiCorp Vault):
```yaml
# k8s-secrets.yaml
apiVersion: v1
kind: Secret
metadata:
  name: rally-app-secrets
type: Opaque
stringData:
  db-password: <зашифрованный пароль>
  session-secret: <зашифрованный ключ>
```

### Ротация секретов

**Периодичность**:
- Database passwords: каждые 90 дней
- Session secrets: каждые 30 дней (с graceful transition)
- API keys: по требованию или при компрометации

## Защита от веб-уязвимостей

### 1. SQL Injection

✅ **Правильно** (параметризованные запросы):
```go
// С использованием pgx
row := db.QueryRow(ctx, 
    "SELECT * FROM crews WHERE id = $1", 
    crewID,
)

// С query builder (squirrel)
query, args, _ := sq.Select("*").
    From("crews").
    Where(sq.Eq{"id": crewID}).
    PlaceholderFormat(sq.Dollar).
    ToSql()
row := db.QueryRow(ctx, query, args...)
```

❌ **Неправильно** (конкатенация строк):
```go
// НИКОГДА ТАК НЕ ДЕЛАТЬ!
query := fmt.Sprintf("SELECT * FROM crews WHERE id = '%s'", crewID)
row := db.QueryRow(ctx, query)
```

### 2. XSS (Cross-Site Scripting)

✅ **Templ автоматически экранирует HTML**:
```templ
templ CrewCard(crew Crew) {
    <div class="crew-card">
        <h3>{ crew.Name }</h3>  // Автоматически безопасно
    </div>
}
```

❌ **Опасно** (если используется raw HTML):
```go
// Избегать! Только для доверенного контента
templ.Raw(userInput)  // НЕ ИСПОЛЬЗОВАТЬ с пользовательским вводом
```

### 3. CSRF (Cross-Site Request Forgery)

**Middleware для CSRF токенов**:
```go
import "github.com/gofiber/fiber/v2/middleware/csrf"

app.Use(csrf.New(csrf.Config{
    KeyLookup:      "form:csrf_token",
    CookieName:     "csrf_",
    CookieSameSite: "Lax",
    Expiration:     1 * time.Hour,
    KeyGenerator:   utils.UUIDs,
}))
```

**В формах**:
```templ
templ CrewRegistrationForm(csrfToken string) {
    <form method="POST" action="/api/crews">
        <input type="hidden" name="csrf_token" value={ csrfToken }/>
        // ... остальные поля
    </form>
}
```

### 4. CORS

**Конфигурация для разных окружений**:

```go
import "github.com/gofiber/fiber/v2/middleware/cors"

func ConfigureCORS(env string) fiber.Handler {
    config := cors.Config{
        AllowCredentials: true,
        AllowMethods: "GET,POST,PUT,DELETE",
    }
    
    switch env {
    case "development":
        config.AllowOrigins = "http://localhost:3000"
    case "staging":
        config.AllowOrigins = "https://staging.rally-app.example.com"
    case "production":
        config.AllowOrigins = "https://rally-app.example.com"
    }
    
    return cors.New(config)
}
```

### 5. Security Headers

```go
func SecurityHeaders() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Защита от clickjacking
        c.Set("X-Frame-Options", "DENY")
        
        // Защита от MIME sniffing
        c.Set("X-Content-Type-Options", "nosniff")
        
        // XSS protection (legacy)
        c.Set("X-XSS-Protection", "1; mode=block")
        
        // Content Security Policy
        c.Set("Content-Security-Policy", 
            "default-src 'self'; "+
            "script-src 'self' 'unsafe-inline'; "+
            "style-src 'self' 'unsafe-inline'; "+
            "img-src 'self' data:;")
        
        // HTTPS enforcement
        c.Set("Strict-Transport-Security", 
            "max-age=31536000; includeSubDomains")
        
        return c.Next()
    }
}
```

## Персональные данные (GDPR/152-ФЗ)

### Типы персональных данных в системе

| Данные | Категория | Требования |
|--------|-----------|------------|
| ФИО участника | Открытые | Могут публиковаться в протоколах |
| Дата рождения | Специальные | Только для проверки возраста, не публиковать |
| Номер лицензии | Открытые | Публикуются в официальных документах |
| Телефон, Email | Контактные | Только для организаторов, не публиковать |
| Адрес | Специальные | Только если требуется регламентом |

### Согласие на обработку

**При регистрации экипажа**:
- [ ] Чекбокс "Согласие на обработку персональных данных"
- [ ] Чекбокс "Согласие на публикацию ФИО и результатов в протоколах"
- [ ] Хранить timestamp согласия

### Права субъектов данных

Система должна поддерживать:
1. **Право на доступ** - пользователь может запросить свои данные
2. **Право на исправление** - пользователь может исправить неточности
3. **Право на удаление** - "право на забвение"
4. **Право на ограничение** - временная блокировка обработки

### Хранение и удаление

**Сроки хранения**:
- Данные активных соревнований: без ограничений
- Данные завершённых соревнований: 3 года (для апелляций)
- После 3 лет: псевдонимизация (заменить ФИО на "Участник №XXX")

**Процедура удаления**:
```sql
-- Псевдонимизация через 3 года после соревнования
UPDATE participants 
SET 
    first_name = 'Участник',
    last_name = CONCAT('№', id),
    email = NULL,
    phone = NULL,
    address = NULL,
    anonymized_at = NOW()
WHERE competition_id IN (
    SELECT id FROM competitions 
    WHERE end_date < NOW() - INTERVAL '3 years'
)
AND anonymized_at IS NULL;
```

## Логирование безопасности

### События для обязательного логирования

```go
// Успешный вход
log.Info().
    Str("user_id", user.ID).
    Str("ip", clientIP).
    Str("user_agent", userAgent).
    Msg("User logged in")

// Неудачная попытка входа
log.Warn().
    Str("username", username).
    Str("ip", clientIP).
    Msg("Failed login attempt")

// Изменение прав
log.Info().
    Str("admin_id", adminID).
    Str("target_user_id", userID).
    Str("old_role", oldRole).
    Str("new_role", newRole).
    Msg("User role changed")

// Доступ к чувствительным данным
log.Info().
    Str("user_id", userID).
    Str("action", "export_personal_data").
    Str("target_user_id", targetUserID).
    Msg("Personal data accessed")
```

## Checklist безопасности для Pull Request

При добавлении нового кода проверьте:

- [ ] Пароли хешируются перед сохранением
- [ ] Используются параметризованные SQL запросы
- [ ] Проверяется авторизация пользователя (middleware)
- [ ] CSRF токены для POST/PUT/DELETE операций
- [ ] Валидируется пользовательский ввод
- [ ] Логируются события безопасности
- [ ] Нет секретов в коде (используются env vars)
- [ ] Персональные данные обрабатываются согласно политике
- [ ] Security headers установлены
- [ ] Error messages не раскрывают внутреннюю структуру

## Процедура реагирования на инциденты

### При обнаружении уязвимости

1. **Немедленно**: Сообщить ответственному за безопасность
2. **В течение 1 часа**: Оценить масштаб проблемы
3. **В течение 4 часов**: Подготовить hotfix
4. **В течение 24 часов**: Развернуть исправление
5. **В течение 72 часов**: Уведомить пользователей (если данные скомпрометированы)

### При утечке данных

1. Немедленная блокировка доступа
2. Оценка объёма утечки
3. Уведомление регулятора (Роскомнадзор) в течение 24 часов
4. Уведомление затронутых пользователей
5. Публичное заявление (при необходимости)
6. Post-mortem и меры предотвращения

## TODO

- [ ] Настроить автоматическое сканирование зависимостей (Dependabot)
- [ ] Внедрить статический анализ безопасности (gosec)
- [ ] Настроить регулярный аудит безопасности
- [ ] Создать процедуру disaster recovery
- [ ] Подготовить шаблоны уведомлений при инцидентах

