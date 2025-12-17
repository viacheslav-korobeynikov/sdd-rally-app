# Руководство по тестированию

**Статус**: Черновик  
**Версия**: 0.1.0  
**Последнее обновление**: 2025-12-17

## Обзор

Этот документ определяет требования к тестированию для SDD Rally App, включая обязательное покрытие критичных доменных правил.

## Пирамида тестирования

```
        /\
       /  \      E2E Tests (Небольшое количество)
      /    \     - Критичные пользовательские сценарии
     /------\    
    /        \   Integration Tests (Среднее количество)
   /          \  - Взаимодействие модулей
  /            \ - API endpoints
 /--------------\
/                \ Unit Tests (Большое количество)
\________________/ - Бизнес-логика, вычисления

```

## Обязательное покрытие Unit-тестами

### Критичные доменные правила (покрытие >80%)

#### 1. Модуль Results - Расчёт результатов

**Что тестировать**:
- Расчёт времени прохождения СУ (финиш - старт)
- Суммирование результатов по всем СУ
- Применение штрафов за опоздание на контроль времени
- Применение штрафов за опережение
- Применение технических штрафов
- Классификация по общему времени
- Классификация по классам автомобилей
- Классификация по группам

**Примеры тест-кейсов**:
```go
// Расчёт времени прохождения СУ
func TestCalculateStageTime_Success(t *testing.T) {
    start := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
    finish := time.Date(2025, 1, 1, 10, 15, 30, 0, time.UTC)
    
    result := CalculateStageTime(start, finish)
    
    assert.Equal(t, 930*time.Second, result) // 15:30
}

// Применение штрафа за опоздание
func TestApplyLatePenalty_10Seconds(t *testing.T) {
    baseTime := 900 * time.Second  // 15:00
    lateness := 10 * time.Second
    
    result := ApplyLatePenalty(baseTime, lateness)
    
    assert.Equal(t, 910*time.Second, result) // 15:10
}
```

#### 2. Модуль Timing - Валидация временных данных

**Что тестировать**:
- Финиш не может быть раньше старта
- Время должно быть в пределах временного окна соревнования
- Нельзя записать финиш без старта
- Нельзя перезаписать время без явного разрешения

#### 3. Модуль Crews - Валидация данных экипажа

**Что тестировать**:
- Уникальность стартовых номеров в рамках соревнования
- Корректность данных лицензий
- Соответствие класса автомобиля допустимым значениям
- Невозможность зарегистрировать один экипаж дважды

## Интеграционные тесты

### Критичные сценарии

#### 1. Полный цикл соревнования
```
Создание соревнования 
→ Добавление СУ 
→ Регистрация экипажей 
→ Фиксация времён 
→ Расчёт результатов 
→ Генерация протокола
```

#### 2. Сценарий с штрафами
```
Регистрация экипажа 
→ Прохождение СУ с опозданием 
→ Применение штрафа 
→ Проверка корректности итогового результата
```

#### 3. Сценарий апелляции
```
Фиксация некорректного времени 
→ Исправление времени 
→ Пересчёт результатов 
→ Проверка изменения позиции в классификации
```

### Структура интеграционных тестов

```go
func TestFullCompetitionWorkflow(t *testing.T) {
    // Setup: База данных, fixtures
    db := setupTestDB(t)
    defer cleanupDB(t, db)
    
    // 1. Создание соревнования
    comp := createTestCompetition(t, db)
    
    // 2. Добавление СУ
    stages := createTestStages(t, db, comp.ID)
    
    // 3. Регистрация экипажей
    crews := registerTestCrews(t, db, comp.ID)
    
    // 4. Фиксация времён
    recordTestTimes(t, db, crews, stages)
    
    // 5. Расчёт результатов
    results, err := CalculateResults(db, comp.ID)
    require.NoError(t, err)
    
    // 6. Проверка корректности
    assert.Len(t, results, len(crews))
    assert.True(t, isSortedByTime(results))
}
```

## Contract Tests (API)

Каждый HTTP endpoint должен иметь contract test, проверяющий:
- Корректный HTTP статус при успехе
- Формат ответа (JSON schema)
- Обработку ошибок (400, 401, 403, 404, 500)
- Валидацию входных данных

### Пример
```go
func TestCreateCompetition_ValidData(t *testing.T) {
    app := setupTestApp(t)
    
    req := httptest.NewRequest("POST", "/api/competitions", strings.NewReader(`{
        "name": "Тестовое ралли",
        "date": "2025-06-15",
        "location": "Москва"
    }`))
    req.Header.Set("Content-Type", "application/json")
    
    resp, _ := app.Test(req)
    
    assert.Equal(t, 201, resp.StatusCode)
    
    var result Competition
    json.NewDecoder(resp.Body).Decode(&result)
    assert.NotEmpty(t, result.ID)
    assert.Equal(t, "Тестовое ралли", result.Name)
}
```

## E2E Tests (Опционально)

E2E тесты через браузер для критичных UI-сценариев:
- Регистрация экипажа через веб-форму
- Ввод результатов хронометража
- Просмотр протокола результатов

## Организация тестов

```
tests/
├── unit/
│   ├── results/
│   │   ├── calculation_test.go
│   │   └── penalties_test.go
│   ├── timing/
│   │   └── validation_test.go
│   └── crews/
│       └── validation_test.go
├── integration/
│   ├── competition_workflow_test.go
│   ├── results_calculation_test.go
│   └── appeals_workflow_test.go
└── contract/
    ├── competitions_api_test.go
    ├── crews_api_test.go
    └── results_api_test.go
```

## Test Fixtures и Helper-функции

Создать переиспользуемые helper-функции для:
- Настройки тестовой базы данных
- Создания тестовых данных (competitions, crews, stages)
- Очистки данных после тестов
- Моки внешних зависимостей

## CI/CD Integration

### Требования для прохождения CI

1. **Unit tests**: Покрытие >80% для модулей Results, Timing, Crews
2. **Integration tests**: Все тесты должны проходить
3. **Contract tests**: Все API endpoints должны быть покрыты
4. **Linting**: golangci-lint без ошибок

### Запуск тестов локально

```bash
# Все тесты
make test

# Только unit
make test-unit

# С покрытием
make test-coverage

# Интеграционные (требуется БД)
make test-integration
```

## Критерии приёмки для Pull Request

- [ ] Новый код покрыт unit-тестами (>80% для критичных модулей)
- [ ] Добавлены интеграционные тесты для новых workflow
- [ ] Все существующие тесты проходят
- [ ] Contract tests обновлены для изменённых API
- [ ] Тесты документированы (что проверяется и почему)

## TODO

- [ ] Настроить автоматические отчёты о покрытии
- [ ] Добавить smoke tests для продакшена
- [ ] Настроить мутационное тестирование для критичных модулей
- [ ] Создать библиотеку общих test fixtures

