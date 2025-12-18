-- +goose Up
-- +goose StatementBegin
-- Создание необходимых расширений
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

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
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TYPE IF EXISTS user_status;
DROP TYPE IF EXISTS user_role;
-- +goose StatementEnd
