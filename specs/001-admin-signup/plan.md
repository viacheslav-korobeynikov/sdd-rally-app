# Implementation Plan: Регистрация администратора ралли

**Branch**: `001-admin-signup` | **Date**: 2025-12-08 | **Spec**: `specs/001-admin-signup/spec.md`
**Input**: Feature specification from `/specs/001-admin-signup/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

Админ должен сам зарегистрироваться, подтвердить email и войти, чтобы управлять мероприятиями/участниками. Технически: Go/Fiber backend с чистой архитектурой, SSR Templ + HTMX, PostgreSQL как единственный сторедж, email-подтверждение с токеном, ограничения попыток (5 за 15 минут на email+IP).

## Technical Context

**Language/Version**: Go 1.21  
**Primary Dependencies**: Fiber, Templ, HTMX, bcrypt/argon2 hashing, стандартный `net/smtp`  
**Storage**: PostgreSQL (единственный сторедж)  
**Testing**: go test + table-driven/unit; критичные доменные сервисы покрываются юнит-тестами  
**Target Platform**: Linux containers (docker-compose локально, готовность к k8s)  
**Project Type**: web (server-rendered, без SPA)  
**Performance Goals**: P95 ответа форм/подтверждения < 500 мс при нагрузке до ~50 rps; регистрация+письмо < 3 мин (SC-001/002)  
**Constraints**: Чистая архитектура handler→service→repository; только параметризованные запросы; без логики в БД; HTMX-only интерактивность; секреты только из env/secret storage (не в репо); tracing/context propagation через middleware  
**Scale/Scope**: Десятки-сотни администраторов; одновременные сессии низкие (до ~50); рост покрывается горизонтальным масштабированием веб-сервера и БД

## Проверка конституции

*GATE: обязательна до исследования (Phase 0) и повторно после дизайна (Phase 1).*

- Модульный монолит с четкими доменами (соревнования, СУ, экипажи, результаты, документы); без пересечений ответственности.
- Backend на Go (Fiber) с чистой архитектурой: handler → service → repository; хэндлеры не ходят в БД; все запросы с context; middleware для логов/трейсов/аутентификации/авторизации.
- UI: серверный рендеринг Templ + HTMX; без SPA; деградация без JS допустима.
- Данные: единственный сторедж PostgreSQL, строгая нормализация; логика в БД запрещена (кроме технических/audit); изменения только миграциями.
- Безопасность: login/password c bcrypt/argon2; минимальные роли (Главный организатор, Секретарь, Хронометраж, Наблюдатель); секреты не в репо; параметризованные запросы.
- Наблюдаемость: структурированные JSON‑логи, базовые метрики (ошибки API, задержка отчетов, количество заявок и т.п.).
- Качество: unit‑тесты для критичных доменных сервисов (расчеты результатов/штрафов); линтер перед мерджем.
- Версионирование и документация: semver + changelog; актуальные README, диаграммы доменов, API‑контракты, каталоги модулей, QA/observability/security‑гайды.

## Project Structure

### Documentation (this feature)

```text
specs/[###-feature]/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)
<!--
  ACTION REQUIRED: Replace the placeholder tree below with the concrete layout
  for this feature. Delete unused options and expand the chosen structure with
  real paths (e.g., apps/admin, packages/something). The delivered plan must
  not include Option labels.
-->

```text
backend/
├── cmd/server/                 # точка входа
├── internal/
│   ├── http/handlers/          # Fiber handlers (SSR/HTMX)
│   ├── services/               # бизнес-логика регистрации/подтверждения/входа
│   ├── repository/             # PostgreSQL доступ (без логики в БД)
│   ├── domain/                 # модели: Admin, ConfirmationToken
│   ├── templates/              # Templ шаблоны для SSR/HTMX
│   └── config/                 # конфигурация, env
├── migrations/                 # golang-migrate
└── pkg/                        # общие утилиты (логирование, почта)

tests/
├── integration/                # e2e потоки регистрации/подтверждения/входа
└── unit/                       # доменные сервисы, rate-limit, токены
```

**Structure Decision**: Web backend с серверным рендерингом; весь фронт в Templ/HTMX в `templates/`, без отдельного SPA.

## Complexity Tracking

Нет заявленных нарушений гейтов; таблица не требуется.
