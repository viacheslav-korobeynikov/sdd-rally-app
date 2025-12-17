# Наблюдаемость системы

**Статус**: Черновик  
**Версия**: 0.1.0  
**Последнее обновление**: 2025-12-17

## Обзор

Этот документ определяет стандарты логирования, сбора метрик, мониторинга и трейсинга для SDD Rally App.

## Три столпа наблюдаемости

1. **Логи** - Что произошло?
2. **Метрики** - Как часто и насколько быстро?
3. **Трейсы** - Какой путь прошёл запрос?

## Логирование

### Библиотека

**Выбор**: `zerolog` или `slog` (стандартная библиотека Go 1.21+)

**Обоснование**: Структурированное логирование в JSON, высокая производительность, zero allocation.

### Уровни логирования

| Уровень | Когда использовать | Примеры |
|---------|-------------------|---------|
| **DEBUG** | Детальная отладочная информация (только в dev) | Значения переменных, промежуточные расчёты |
| **INFO** | Значимые бизнес-события | Создано соревнование, зарегистрирован экипаж, рассчитаны результаты |
| **WARN** | Восстанавливаемые ошибки, нештатные ситуации | Повторная попытка запроса, использование дефолтного значения |
| **ERROR** | Ошибки, требующие внимания | Не удалось сохранить в БД, внешний сервис недоступен |
| **FATAL** | Критические ошибки, после которых приложение не может продолжать работу | Невозможно подключиться к БД при старте |

### Структура лог-записи

**Обязательные поля**:
```json
{
  "timestamp": "2025-12-17T10:30:00Z",
  "level": "info",
  "message": "Competition created",
  "request_id": "abc-123-def-456",
  "user_id": "user_789",
  "service": "sdd-rally-app",
  "version": "1.2.3"
}
```

**Контекстные поля** (специфичные для операции):
```json
{
  "competition_id": "comp_001",
  "crew_id": "crew_042",
  "stage_id": "stage_05",
  "duration_ms": 123
}
```

### Примеры логирования

#### Успешная операция (INFO)
```go
log.Info().
    Str("request_id", requestID).
    Str("competition_id", comp.ID).
    Str("name", comp.Name).
    Msg("Competition created successfully")
```

#### Ошибка с контекстом (ERROR)
```go
log.Error().
    Err(err).
    Str("request_id", requestID).
    Str("crew_id", crewID).
    Str("stage_id", stageID).
    Msg("Failed to record finish time")
```

#### Измерение производительности (DEBUG/INFO)
```go
start := time.Now()
// ... выполнение операции
duration := time.Since(start)

log.Info().
    Str("operation", "calculate_results").
    Str("competition_id", compID).
    Dur("duration", duration).
    Int("crews_count", len(crews)).
    Msg("Results calculated")
```

### Что НЕ логировать

❌ **Запрещено**:
- Пароли, токены, API ключи
- Полные номера банковских карт
- Персональные данные (ФИО, телефоны, email) - только ID
- Большие объёмы данных (массивы с сотнями элементов)

✅ **Разрешено**:
- ID пользователей/экипажей/соревнований
- Количественные метрики (сколько, как долго)
- Статусы операций (успех/ошибка)
- Технические детали ошибок (без чувствительных данных)

## Метрики

### Инструменты

**Выбор**: Prometheus + Grafana

**Библиотека**: `github.com/prometheus/client_golang`

### Обязательные метрики

#### 1. HTTP метрики

```go
// Количество запросов
http_requests_total{method="GET", path="/api/competitions", status="200"}

// Длительность запросов (histogram)
http_request_duration_seconds{method="POST", path="/api/crews"}

// Размер ответов
http_response_size_bytes{method="GET", path="/api/results"}
```

#### 2. Бизнес-метрики

```go
// Количество соревнований по статусам
competitions_total{status="planned|active|completed"}

// Количество экипажей
crews_total{competition_id="comp_001"}

// Количество заявок за период
crew_registrations_total

// Время генерации протокола
protocol_generation_duration_seconds

// Количество рассчитанных результатов
results_calculated_total
```

#### 3. Метрики базы данных

```go
// Длительность запросов к БД
db_query_duration_seconds{operation="select|insert|update|delete"}

// Количество активных соединений
db_connections_active

// Количество ошибок БД
db_errors_total{type="connection|query|transaction"}
```

#### 4. Системные метрики

```go
// Использование памяти
go_memstats_alloc_bytes

// Количество горутин
go_goroutines

// Время работы сервиса
process_uptime_seconds
```

### Типы метрик

| Тип | Когда использовать | Пример |
|-----|-------------------|--------|
| **Counter** | Только увеличивается | Количество запросов, ошибок |
| **Gauge** | Может увеличиваться и уменьшаться | Количество активных соединений, память |
| **Histogram** | Распределение значений | Длительность запросов, размеры ответов |
| **Summary** | Процентили (p50, p95, p99) | Время ответа API |

### Пример реализации

```go
var (
    httpRequestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "path", "status"},
    )
    
    httpRequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "http_request_duration_seconds",
            Help: "HTTP request duration in seconds",
            Buckets: []float64{0.01, 0.05, 0.1, 0.5, 1, 5},
        },
        []string{"method", "path"},
    )
)

// В middleware
func MetricsMiddleware(c *fiber.Ctx) error {
    start := time.Now()
    
    err := c.Next()
    
    duration := time.Since(start).Seconds()
    status := c.Response().StatusCode()
    
    httpRequestsTotal.WithLabelValues(
        c.Method(), 
        c.Path(), 
        fmt.Sprintf("%d", status),
    ).Inc()
    
    httpRequestDuration.WithLabelValues(
        c.Method(), 
        c.Path(),
    ).Observe(duration)
    
    return err
}
```

## Трейсинг (Опционально на старте)

### Инструменты

**Выбор**: OpenTelemetry + Jaeger

### Что трейсить

- HTTP запросы (входящие)
- Запросы к БД
- Межсервисные вызовы (в будущем)
- Критичные операции (расчёт результатов)

### Пример span

```go
ctx, span := tracer.Start(ctx, "CalculateResults")
defer span.End()

span.SetAttributes(
    attribute.String("competition.id", compID),
    attribute.Int("crews.count", len(crews)),
)

// ... бизнес-логика

if err != nil {
    span.RecordError(err)
    span.SetStatus(codes.Error, err.Error())
}
```

## Дашборды Grafana

### Dashboard 1: Обзор системы

**Панели**:
- Requests per second (RPS)
- Error rate (%)
- P95/P99 latency
- Active connections
- Memory usage
- Goroutines count

### Dashboard 2: Бизнес-метрики

**Панели**:
- Активные соревнования
- Зарегистрированные экипажи (по дням)
- Среднее время генерации протоколов
- Количество ошибок при расчёте результатов

### Dashboard 3: База данных

**Панели**:
- Query duration (по типам операций)
- Connection pool usage
- Slow queries (>1s)
- Database errors

## Алертинг

### Критичные алерты (PagerDuty/Email)

| Алерт | Условие | Действие |
|-------|---------|----------|
| Service Down | HTTP health check fails 3 times | Немедленное реагирование |
| Database Unavailable | db_errors > 10/min | Проверить подключение к БД |
| High Error Rate | error_rate > 5% | Проверить логи последних ошибок |
| Memory Leak | memory growth >50MB/hour | Проверить goroutine leaks |

### Предупреждения (Slack/Email)

| Алерт | Условие | Действие |
|-------|---------|----------|
| Slow Responses | P95 latency > 2s | Оптимизировать запросы |
| High Load | RPS > 1000 | Подготовиться к масштабированию |
| Disk Space Low | disk usage > 80% | Очистить старые логи/бэкапы |

## Health Check Endpoints

### `/health` - Liveness probe
Проверяет, что приложение запущено и может принимать запросы.

```json
{
  "status": "ok",
  "timestamp": "2025-12-17T10:30:00Z"
}
```

### `/ready` - Readiness probe
Проверяет, что все зависимости доступны.

```json
{
  "status": "ready",
  "checks": {
    "database": "ok",
    "migrations": "ok"
  },
  "timestamp": "2025-12-17T10:30:00Z"
}
```

## Хранение и ротация логов

### Development
- Вывод в stdout (для docker-compose logs)
- Уровень: DEBUG
- Формат: JSON

### Staging/Production
- Централизованное хранение (ELK/Loki)
- Уровень: INFO
- Ротация: 7 дней для INFO, 30 дней для ERROR
- Сжатие: gzip

## Checklist для нового feature

При добавлении нового функционала убедитесь:

- [ ] Добавлено INFO логирование для значимых событий
- [ ] Добавлено ERROR логирование для ошибок
- [ ] Добавлены метрики для критичных операций
- [ ] Обновлены дашборды Grafana (если нужно)
- [ ] Настроены алерты (если критичный функционал)
- [ ] Request ID пробрасывается через весь стек

## TODO

- [ ] Настроить централизованное хранилище логов
- [ ] Создать базовые Grafana dashboards
- [ ] Настроить алертинг в Prometheus
- [ ] Интегрировать OpenTelemetry для трейсинга
- [ ] Создать runbook для типичных инцидентов

