# Tasks: Регистрация администратора ралли

**Input**: plan.md, spec.md, data-model.md, contracts/http.md, research.md, quickstart.md  
**Prerequisites**: plan.md (заполнен), spec.md, research.md, data-model.md, contracts/, quickstart.md  

## Format: `[ID] [P?] [Story] Description`

## Phase 1: Setup (Shared Infrastructure)

- [ ] T001 Инициализировать Go модуль (`backend/go.mod`, `backend/go.sum`)
- [ ] T002 Создать `.env.example` с POSTGRES_DSN, SMTP_HOST/PORT/USER/PASSWORD, APP_BASE_URL (`.env.example`)
- [ ] T003 Обновить `docker-compose.yml` для mailhog (smtp:1025, ui:8025) и postgres (`docker-compose.yml`)
- [ ] T004 [P] Добавить зависимости: Fiber, Templ, HTMX helper, bcrypt/argon2, pg driver, migrate, rate-limit, `net/smtp` (`backend/go.mod`)

## Phase 2: Foundational (Blocking Prerequisites)

- [ ] T005 Добавить миграции для таблиц `admins` и `confirmation_tokens` (`migrations/*.sql`)
- [ ] T006 Настроить конфиги: Postgres DSN, SMTP (`net/smtp`), base URL, rate-limit (5/15m email+IP) (`backend/internal/config/config.go`)
- [ ] T007 [P] Настроить логирование JSON и метрики (счетчики попыток, ошибки) (`backend/pkg/logger/*.go`, `backend/pkg/metrics/*.go`)
- [ ] T008 [P] Добавить tracing/context middleware (включая request-id) (`backend/internal/http/middleware/tracing.go`)
- [ ] T009 [P] Реализовать почтовый сервис через стандартный `net/smtp` с интерфейсом и адаптером для mailhog/prod (`backend/pkg/mail/service.go`, `backend/pkg/mail/smtp.go`)
- [ ] T010 [P] Добавить middleware rate-limit по email+IP (5 попыток/15 минут) для signup/login (`backend/internal/http/middleware/ratelimit.go`)
- [ ] T011 Подключить Templ/HTMX базовый layout и статику (`backend/internal/templates/layouts/base.templ`)
- [ ] T012 Обеспечить хранение секретов только через env/secret storage, исключить секреты из репо (`backend/internal/config/config.go`, `.gitignore`)

## Phase 3: User Story 1 - Самостоятельная регистрация администратора (P1) 🎯 MVP

**Goal**: Регистрация администратора с формой и валидацией, статус pending_confirmation.  
**Independent Test**: Новый email проходит регистрацию, получает письмо, учетная запись pending_confirmation.

### Implementation
- [ ] T013 [US1] Добавить доменную модель Admin/Token и репозитории (create/find-by-email) (`backend/internal/domain/admin.go`, `backend/internal/repository/admin_repo.go`)
- [ ] T014 [P] [US1] Сервис регистрации: валидация, хэш пароля, создание admin pending, создание токена, вызов mail сервиса (`backend/internal/services/signup_service.go`)
- [ ] T015 [P] [US1] Handler + HTMX/SSR формы регистрации, ошибки валидации, success экран (`backend/internal/http/handlers/signup_handler.go`, `backend/internal/templates/signup.templ`)
- [ ] T016 [US1] Handler/endpoint для повторной отправки подтверждения (resend), без создания дублей (`backend/internal/http/handlers/resend_handler.go`)
- [ ] T017 [US1] Интеграционный тест потока регистрации + resend + письмо (mailhog) (`tests/integration/signup_flow_test.go`)

## Phase 4: User Story 2 - Подтверждение email и активация (P1)

**Goal**: Подтверждение email по токену, статус active.  
**Independent Test**: Переход по валидной ссылке активирует учетку; истекшая/повторная ссылка отклоняется.

### Implementation
- [ ] T018 [US2] Репозиторий токенов: выборка/пометка использован/истек (`backend/internal/repository/token_repo.go`)
- [ ] T019 [P] [US2] Сервис подтверждения: проверка срока/повторного использования, активация admin, idempotency (`backend/internal/services/confirm_service.go`)
- [ ] T020 [P] [US2] Handler GET /admin/confirm с экранами успех/ошибка (`backend/internal/http/handlers/confirm_handler.go`, `backend/internal/templates/confirm.templ`)
- [ ] T021 [US2] Интеграционный тест подтверждения (валидный, истекший, повторно использованный) (`tests/integration/confirm_flow_test.go`)

## Phase 5: User Story 3 - Первый вход и доступ к управлению (P2)

**Goal**: Вход активированного администратора, доступ к разделам управления.  
**Independent Test**: Active аккаунт входит; неактивный/неверный пароль — отказ; лимит попыток действует.

### Implementation
- [ ] T022 [US3] Сервис аутентификации: проверка статуса active, сравнение хэша, учет лимита попыток (`backend/internal/services/auth_service.go`)
- [ ] T023 [P] [US3] Handler POST /admin/login с ответами без утечки, установка сессии/токена (`backend/internal/http/handlers/login_handler.go`, `backend/internal/templates/login.templ`)
- [ ] T024 [US3] Интеграционный тест входа: активный/неактивный/неверный пароль/лимит (`tests/integration/login_flow_test.go`)

## Phase N: Polish & Cross-Cutting Concerns

- [ ] T025 [P] Обновить quickstart с примерами curl/HTMX и mailhog ссылкой (`specs/001-admin-signup/quickstart.md`)
- [ ] T026 [P] Документация API: актуализировать `specs/001-admin-signup/contracts/http.md` и корневые `docs/api-contracts/` по фактическим handler’ам (`specs/001-admin-signup/contracts/http.md`, `docs/api-contracts/`)
- [ ] T027 [P] Добавить unit-тесты доменных сервисов (signup/confirm/auth) как обязательные (table-driven) (`tests/unit/services/*.go`)
- [ ] T028 [P] Обновить README бэкенда и диаграмму доменных сущностей (сущности Admin/Token) (`README.md`, `docs/domain-diagram.png`)
- [ ] T029 [P] Привести ответы/ошибки к единому формату и человеко-понятным сообщениям (FR-011) (`backend/internal/http/handlers/*`)

## Dependencies & Execution Order

- Setup (Phase 1) → Foundational (Phase 2) → User stories (US1 P1 → US2 P1 → US3 P2) → Polish.
- US2 зависит от токенов, созданных в US1; US3 зависит от активированных аккаунтов из US2.

## Parallel Opportunities

- В Phase 2 параллельно: T006 (логи/метрики), T007 (mail), T008 (rate-limit), T009 (templ) после конфигов.
- В US1: T011 (сервис) и T012 (handler) параллельно после T010; тест T013 после реализации.
- В US2: T015 и T016 параллельно после T014; тест T017 после реализации.
- В US3: T018 и T019 параллельно; тест T020 после реализации.

## Implementation Strategy

- MVP: завершить US1 (регистрация + письмо).
- Затем US2 (активация), далее US3 (вход).
- После — Polish и документирование.

