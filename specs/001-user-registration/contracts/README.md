# API Контракты: Регистрация и аутентификация

**Feature**: 001-user-registration  
**Дата**: 2025-12-17

## Обзор

Этот каталог содержит OpenAPI спецификацию для auth endpoints.

## Файлы

- `auth-api.yaml` - OpenAPI 3.0 спецификация для всех auth endpoints

## Endpoints

### Публичные (без авторизации)

- `POST /api/auth/register` - Регистрация нового пользователя
- `POST /api/auth/login` - Вход в систему  
- `GET /api/auth/me` - Получить текущего пользователя (требует auth)

### Защищённые (требуют авторизации)

- `POST /api/auth/logout` - Выход из системы
- `GET /api/users` - Список пользователей (только chief_organizer)
- `PATCH /api/users/:id/role` - Изменить роль пользователя (только chief_organizer)

## Форматы

- **Запросы**: JSON (Content-Type: application/json)
- **Ответы**: JSON
- **Даты**: ISO 8601 (RFC3339)
- **Ошибки**: Стандартный формат ErrorResponse

## Тестирование

Спецификацию можно импортировать в:
- Postman
- Insomnia
- Swagger UI
- OpenAPI Generator (для генерации клиентов)

## Просмотр спецификации

```bash
# Установить swagger-ui (если нет)
npm install -g swagger-ui-watcher

# Запустить
swagger-ui-watcher auth-api.yaml
```

Или использовать онлайн редактор: https://editor.swagger.io/

