-- +goose Up
-- +goose StatementBegin
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
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS role_changes;
-- +goose StatementEnd
