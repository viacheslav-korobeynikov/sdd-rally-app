-- +goose Up
-- +goose StatementBegin
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
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS login_attempts;
-- +goose StatementEnd
