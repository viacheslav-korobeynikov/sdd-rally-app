# API Contracts: Регистрация администратора ралли

**Branch**: `001-admin-signup`  
**Date**: 2025-12-08

## POST /admin/signup
- Purpose: Саморегистрация администратора.
- Request (application/json):
  - full_name (string, required)
  - email (string, required, email)
  - password (string, required, policy per FR-003)
  - accept_terms (boolean, required)
- Responses:
  - 201 Created: `{ "status": "pending_confirmation" }`
  - 400 ValidationError: список полей с ошибками
  - 409 Conflict: email уже существует
  - 429 Too Many Requests: превышен лимит (5 попыток/15 мин на email+IP)

## POST /admin/signup/resend-confirmation
- Purpose: Повторная отправка письма подтверждения.
- Request: `{ "email": "user@example.com" }`
- Responses:
  - 200 OK: письмо поставлено в очередь
  - 400 ValidationError
  - 404 Not Found: учетной записи нет (опционально возвращать 200 для неразглашения)
  - 429 Too Many Requests

## GET /admin/confirm
- Purpose: Подтверждение email по токену.
- Query: `token` (string, required)
- Responses:
  - 302/200 Success: email подтвержден, статус active
  - 400/410 InvalidOrExpired: токен истек или уже использован

## POST /admin/login
- Purpose: Первый вход администратора.
- Request: `{ "email": "user@example.com", "password": "secret" }`
- Responses:
  - 200 OK: сессия установлена/токен выдан; тело без лишних деталей
  - 401 Unauthorized: неверные креденшелы или неактивный статус
  - 429 Too Many Requests: превышен лимит попыток

## Common considerations
- Все запросы используют HTTPS, параметризованные обращения к БД.
- Логи: структурированные JSON, ошибки/метрики (успех/ошибка, задержка, попытки).
- Безопасность: хэш паролей (bcrypt/argon2), не раскрывать, существует ли email при ошибке входа.

