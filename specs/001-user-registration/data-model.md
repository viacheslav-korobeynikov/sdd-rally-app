# Модель данных: Регистрация и аутентификация пользователей

**Feature**: 001-user-registration  
**Дата**: 2025-12-17  
**Статус**: Утверждено

## Обзор

Модель данных для системы аутентификации и авторизации включает 4 основные сущности:
1. **Users** - Пользователи системы
2. **Sessions** - Активные сессии пользователей
3. **Login Attempts** - История попыток входа (для безопасности и аналитики)
4. **Role Changes** - Audit trail изменений ролей

## Диаграмма ER

```
┌─────────────────────┐
│       Users         │
├─────────────────────┤
│ id (UUID, PK)       │
│ username (UNIQUE)   │◄─────┐
│ password_hash       │      │
│ role (ENUM)         │      │
│ status (ENUM)       │      │
│ created_at          │      │
│ updated_at          │      │
│ last_login_at       │      │
└─────────────────────┘      │
         △                   │
         │ user_id           │
         │                   │
┌────────┴────────────┐      │
│     Sessions        │      │
├─────────────────────┤      │
│ id (UUID, PK)       │      │
│ user_id (FK)        │──────┘
│ token_hash          │
│ created_at          │
│ last_activity_at    │
│ expires_at          │
│ ip_address          │
│ user_agent          │
└─────────────────────┘

┌─────────────────────┐
│   Login Attempts    │
├─────────────────────┤
│ id (BIGSERIAL, PK)  │
│ username_attempt    │
│ ip_address          │
│ user_agent          │
│ success (BOOLEAN)   │
│ failure_reason      │
│ created_at          │
└─────────────────────┘

┌─────────────────────┐      ┌──────────────────┐
│   Role Changes      │      │    Enum Types    │
├─────────────────────┤      ├──────────────────┤
│ id (BIGSERIAL, PK)  │      │ user_role:       │
│ user_id (FK)        │      │ - chief_organizer│
│ changed_by_id (FK)  │      │ - secretary      │
│ old_role            │      │ - timing         │
│ new_role            │      │ - observer       │
│ reason (TEXT)       │      │                  │
│ created_at          │      │ user_status:     │
└─────────────────────┘      │ - active         │
                              │ - locked         │
                              └──────────────────┘
```

## Детальное описание таблиц

### 1. users

**Назначение**: Хранение данных пользователей системы

**Структура**:

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role user_role NOT NULL DEFAULT 'observer',
    status user_status NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_login_at TIMESTAMPTZ,
    
    -- Constraints
    CONSTRAINT username_min_length CHECK (LENGTH(username) >= 3),
    CONSTRAINT username_format CHECK (username ~ '^[a-z0-9_-]+$')
);

-- Indexes
CREATE UNIQUE INDEX idx_users_username ON users (LOWER(username));
CREATE INDEX idx_users_status ON users (status) WHERE status = 'active';
CREATE INDEX idx_users_role ON users (role);
CREATE INDEX idx_users_created_at ON users (created_at DESC);
```

**Поля**:

| Поле | Тип | Обязательное | Описание |
|------|-----|--------------|----------|
| `id` | UUID | Да | Уникальный идентификатор пользователя |
| `username` | VARCHAR(50) | Да | Логин пользователя (нечувствительный к регистру) |
| `password_hash` | VARCHAR(255) | Да | Bcrypt хеш пароля (60 символов, но оставляем запас) |
| `role` | user_role | Да | Роль пользователя (ENUM) |
| `status` | user_status | Да | Статус аккаунта (ENUM) |
| `created_at` | TIMESTAMPTZ | Да | Дата создания аккаунта |
| `updated_at` | TIMESTAMPTZ | Да | Дата последнего обновления |
| `last_login_at` | TIMESTAMPTZ | Нет | Дата последнего успешного входа |

**Бизнес-правила**:
- Логин: 3-50 символов, только латинские буквы, цифры, дефис, подчеркивание
- Логин нечувствителен к регистру (хранится в нижнем регистре)
- Первый зарегистрированный пользователь автоматически получает роль `chief_organizer`
- Все последующие получают роль `observer` до явного изменения

**Enum типы**:

```sql
CREATE TYPE user_role AS ENUM (
    'chief_organizer',  -- Главный организатор (полный доступ)
    'secretary',        -- Секретарь (управление экипажами)
    'timing',           -- Хронометраж (ввод результатов)
    'observer'          -- Наблюдатель (только чтение)
);

CREATE TYPE user_status AS ENUM (
    'active',   -- Активный аккаунт
    'locked'    -- Заблокированный (после брутфорса или вручную)
);
```

---

### 2. sessions

**Назначение**: Хранение активных сессий пользователей

**Структура**:

```sql
CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_activity_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL,
    ip_address INET NOT NULL,
    user_agent TEXT NOT NULL
);

-- Indexes
CREATE UNIQUE INDEX idx_sessions_token_hash ON sessions (token_hash);
CREATE INDEX idx_sessions_user_id ON sessions (user_id);
CREATE INDEX idx_sessions_expires_at ON sessions (expires_at) WHERE expires_at > NOW();
CREATE INDEX idx_sessions_last_activity ON sessions (last_activity_at DESC);
```

**Поля**:

| Поле | Тип | Обязательное | Описание |
|------|-----|--------------|----------|
| `id` | UUID | Да | Уникальный идентификатор сессии |
| `user_id` | UUID | Да | ID пользователя (FK to users) |
| `token_hash` | VARCHAR(255) | Да | Хеш токена сессии (для cookie) |
| `created_at` | TIMESTAMPTZ | Да | Дата создания сессии |
| `last_activity_at` | TIMESTAMPTZ | Да | Дата последней активности |
| `expires_at` | TIMESTAMPTZ | Да | Дата истечения сессии |
| `ip_address` | INET | Да | IP-адрес клиента |
| `user_agent` | TEXT | Да | User-Agent браузера |

**Бизнес-правила**:
- Время жизни сессии: 1 час с момента последней активности
- `last_activity_at` и `expires_at` обновляются при каждом запросе
- Токен сессии генерируется случайно (32 байта), хешируется SHA-256
- В cookie хранится оригинальный токен, в БД - хеш
- При выходе сессия удаляется из БД
- Expired сессии автоматически очищаются фоновой задачей
- Сессии с `last_activity_at` старше 1 часа считаются неактивными и завершаются

**Автоочистка**:

```sql
-- Cronjob или фоновая задача раз в час
DELETE FROM sessions 
WHERE expires_at < NOW() - INTERVAL '1 day';
```

---

### 3. login_attempts

**Назначение**: Логирование всех попыток входа для безопасности и аналитики

**Структура**:

```sql
CREATE TABLE login_attempts (
    id BIGSERIAL PRIMARY KEY,
    username_attempt VARCHAR(255) NOT NULL,
    ip_address INET NOT NULL,
    user_agent TEXT NOT NULL,
    success BOOLEAN NOT NULL,
    failure_reason VARCHAR(100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_login_attempts_ip_created ON login_attempts (ip_address, created_at DESC);
CREATE INDEX idx_login_attempts_username_created ON login_attempts (username_attempt, created_at DESC);
CREATE INDEX idx_login_attempts_success ON login_attempts (success);
CREATE INDEX idx_login_attempts_created_at ON login_attempts (created_at DESC);
```

**Поля**:

| Поле | Тип | Обязательное | Описание |
|------|-----|--------------|----------|
| `id` | BIGSERIAL | Да | Автоинкрементный ID |
| `username_attempt` | VARCHAR(255) | Да | Логин, с которым пытались войти |
| `ip_address` | INET | Да | IP-адрес клиента |
| `user_agent` | TEXT | Да | User-Agent браузера |
| `success` | BOOLEAN | Да | Успешна ли попытка |
| `failure_reason` | VARCHAR(100) | Нет | Причина неудачи (если success=false) |
| `created_at` | TIMESTAMPTZ | Да | Дата попытки |

**Failure reasons**:
- `user_not_found` - Пользователь не найден
- `invalid_password` - Неверный пароль
- `account_locked` - Аккаунт заблокирован
- `rate_limited` - Превышен лимит попыток

**Бизнес-правила**:
- Логируются ВСЕ попытки входа (успешные и неуспешные)
- Используется для rate limiting (подсчёт попыток за 15 минут)
- Данные хранятся 30 дней, затем архивируются/удаляются

**Автоочистка**:

```sql
-- Cronjob раз в сутки
DELETE FROM login_attempts 
WHERE created_at < NOW() - INTERVAL '30 days';
```

**Запросы для rate limiting**:

```sql
-- Подсчёт неудачных попыток с IP за 15 минут
SELECT COUNT(*) 
FROM login_attempts 
WHERE ip_address = $1 
  AND success = false 
  AND created_at > NOW() - INTERVAL '15 minutes';

-- Подсчёт неудачных попыток для логина за 15 минут
SELECT COUNT(*) 
FROM login_attempts 
WHERE username_attempt = LOWER($1)
  AND success = false 
  AND created_at > NOW() - INTERVAL '15 minutes';
```

---

### 4. role_changes

**Назначение**: Audit trail изменений ролей пользователей

**Структура**:

```sql
CREATE TABLE role_changes (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    changed_by_id UUID NOT NULL REFERENCES users(id) ON DELETE SET NULL,
    old_role user_role NOT NULL,
    new_role user_role NOT NULL,
    reason TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_role_changes_user_id ON role_changes (user_id, created_at DESC);
CREATE INDEX idx_role_changes_changed_by ON role_changes (changed_by_id, created_at DESC);
CREATE INDEX idx_role_changes_created_at ON role_changes (created_at DESC);
```

**Поля**:

| Поле | Тип | Обязательное | Описание |
|------|-----|--------------|----------|
| `id` | BIGSERIAL | Да | Автоинкрементный ID |
| `user_id` | UUID | Да | ID пользователя, чью роль изменили |
| `changed_by_id` | UUID | Да | ID пользователя, который изменил роль |
| `old_role` | user_role | Да | Старая роль |
| `new_role` | user_role | Да | Новая роль |
| `reason` | TEXT | Нет | Причина изменения (опционально) |
| `created_at` | TIMESTAMPTZ | Да | Дата изменения |

**Бизнес-правила**:
- Запись создаётся КАЖДЫЙ раз при изменении роли
- Используется для аудита (кто, когда, почему изменил роль)
- Данные не удаляются (вечное хранение для комплаенса)
- При удалении пользователя `changed_by_id` становится NULL

---

## Миграции

### Порядок применения

1. **001_create_enum_types.up.sql** - Создание enum типов
2. **002_create_users_table.up.sql** - Создание таблицы users
3. **003_create_sessions_table.up.sql** - Создание таблицы sessions
4. **004_create_login_attempts_table.up.sql** - Создание таблицы login_attempts
5. **005_create_role_changes_table.up.sql** - Создание таблицы role_changes

### Примеры миграций

**001_create_enum_types.up.sql**:
```sql
CREATE TYPE user_role AS ENUM (
    'chief_organizer',
    'secretary',
    'timing',
    'observer'
);

CREATE TYPE user_status AS ENUM (
    'active',
    'locked'
);
```

**001_create_enum_types.down.sql**:
```sql
DROP TYPE IF EXISTS user_status;
DROP TYPE IF EXISTS user_role;
```

---

## Go модели

### User

```go
package models

import (
    "time"
    "github.com/google/uuid"
)

type UserRole string

const (
    RoleChiefOrganizer UserRole = "chief_organizer"
    RoleSecretary      UserRole = "secretary"
    RoleTiming         UserRole = "timing"
    RoleObserver       UserRole = "observer"
)

type UserStatus string

const (
    StatusActive UserStatus = "active"
    StatusLocked UserStatus = "locked"
)

type User struct {
    ID           uuid.UUID  `json:"id"`
    Username     string     `json:"username"`
    PasswordHash string     `json:"-"` // Никогда не сериализуем
    Role         UserRole   `json:"role"`
    Status       UserStatus `json:"status"`
    CreatedAt    time.Time  `json:"created_at"`
    UpdatedAt    time.Time  `json:"updated_at"`
    LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
}

// HasPermission проверяет, есть ли у пользователя право
func (u *User) HasPermission(permission string) bool {
    return rolePermissions[u.Role][permission]
}
```

### Session

```go
type Session struct {
    ID             uuid.UUID `json:"id"`
    UserID         uuid.UUID `json:"user_id"`
    TokenHash      string    `json:"-"` // Не сериализуем
    CreatedAt      time.Time `json:"created_at"`
    LastActivityAt time.Time `json:"last_activity_at"`
    ExpiresAt      time.Time `json:"expires_at"`
    IPAddress      string    `json:"ip_address"`
    UserAgent      string    `json:"user_agent"`
}

// IsExpired проверяет, истекла ли сессия
func (s *Session) IsExpired() bool {
    return time.Now().After(s.ExpiresAt)
}

// UpdateActivity обновляет время последней активности
func (s *Session) UpdateActivity() {
    s.LastActivityAt = time.Now()
    s.ExpiresAt = time.Now().Add(1 * time.Hour)
}
```

### LoginAttempt

```go
type LoginAttempt struct {
    ID              int64      `json:"id"`
    UsernameAttempt string     `json:"username_attempt"`
    IPAddress       string     `json:"ip_address"`
    UserAgent       string     `json:"user_agent"`
    Success         bool       `json:"success"`
    FailureReason   *string    `json:"failure_reason,omitempty"`
    CreatedAt       time.Time  `json:"created_at"`
}
```

### RoleChange

```go
type RoleChange struct {
    ID          int64     `json:"id"`
    UserID      uuid.UUID `json:"user_id"`
    ChangedByID uuid.UUID `json:"changed_by_id"`
    OldRole     UserRole  `json:"old_role"`
    NewRole     UserRole  `json:"new_role"`
    Reason      *string   `json:"reason,omitempty"`
    CreatedAt   time.Time `json:"created_at"`
}
```

---

## Права доступа (Permissions Matrix)

| Действие | chief_organizer | secretary | timing | observer |
|----------|-----------------|-----------|--------|----------|
| Регистрация | ✅ | ✅ | ✅ | ✅ |
| Вход | ✅ | ✅ | ✅ | ✅ |
| Выход | ✅ | ✅ | ✅ | ✅ |
| Просмотр списка пользователей | ✅ | ❌ | ❌ | ❌ |
| Изменение роли пользователя | ✅ | ❌ | ❌ | ❌ |
| Блокировка пользователя | ✅ | ❌ | ❌ | ❌ |
| Просмотр audit логов | ✅ | ❌ | ❌ | ❌ |

---

## Оптимизации

### Connection Pooling

```go
// config/database.go
config := &pgxpool.Config{
    MaxConns:          25,
    MinConns:          5,
    MaxConnLifetime:   5 * time.Minute,
    MaxConnIdleTime:   1 * time.Minute,
    HealthCheckPeriod: 30 * time.Second,
}
```

### Индексы для производительности

Все критичные запросы покрыты индексами:
- ✅ Поиск пользователя по username (UNIQUE INDEX)
- ✅ Проверка сессии по token (UNIQUE INDEX)
- ✅ Подсчёт попыток входа (COMPOSITE INDEX на ip_address + created_at)
- ✅ История активности (INDEX на created_at DESC для всех таблиц)

### Автоматическая очистка

```go
// internal/tasks/cleanup.go
func (t *CleanupTask) Run() {
    // Каждый час удаляем expired сессии
    t.db.Exec(ctx, "DELETE FROM sessions WHERE expires_at < NOW() - INTERVAL '1 day'")
    
    // Каждый день удаляем старые login_attempts
    t.db.Exec(ctx, "DELETE FROM login_attempts WHERE created_at < NOW() - INTERVAL '30 days'")
}
```

---

## Безопасность данных

### Хранение паролей

- ❌ **Никогда**: Plain text, MD5, SHA-1
- ✅ **Всегда**: bcrypt с cost 12+ или argon2id

### Хранение сессий

- ✅ Token хешируется SHA-256 перед сохранением в БД
- ✅ Cookie HTTP-only, Secure (HTTPS), SameSite=Lax
- ✅ IP validation (опционально, можно отключить для мобильных сетей)

### Audit Trail

- ✅ Все изменения ролей логируются в `role_changes`
- ✅ Все попытки входа логируются в `login_attempts`
- ✅ Никакие данные аудита не удаляются (comплаенс)

---

## Следующие шаги

1. ✅ Модель данных утверждена
2. → Создать SQL миграции
3. → Создать Go модели и репозитории
4. → Написать тесты для репозиториев
5. → Перейти к API контрактам (contracts/)

**Статус**: Модель данных готова к реализации.

