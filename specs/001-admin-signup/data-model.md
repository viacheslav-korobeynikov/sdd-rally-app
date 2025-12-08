# Data Model: Регистрация администратора ралли

**Branch**: `001-admin-signup`  
**Date**: 2025-12-08  
**Source**: `specs/001-admin-signup/spec.md`

## Entities

### Администратор
- id (uuid)
- full_name (string, required)
- email (string, unique, required, lowercase)
- password_hash (string, required)
- status (enum: pending_confirmation, active, blocked)
- terms_accepted_at (timestamp, required)
- created_at / updated_at (timestamp)
- activated_at (timestamp, nullable)

Constraints & Rules:
- email уникален.
- status transitions: pending_confirmation → active (после подтверждения); active ↔ blocked (админом системы).
- пароль хранится только как хэш (bcrypt/argon2).

### Токен подтверждения
- id (uuid)
- admin_id (fk → Администратор.id)
- token (string, unique)
- expires_at (timestamp)
- used_at (timestamp, nullable)
- created_at (timestamp)

Constraints & Rules:
- Один активный токен на администратора; при повторной отправке предыдущий помечается использованным/истекшим.
- Токен недействителен после expires_at или used_at.

## Relationships
- Администратор 1—N Токен подтверждения (активный максимум один).

## Validation & Derived Rules
- Пароль: минимальная длина и базовая сложность (см. FR-003).
- Лимиты попыток: 5 попыток за 15 минут на email+IP (FR-010).
- Вход возможен только при status=active (FR-009).

