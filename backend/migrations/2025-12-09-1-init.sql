-- Создаем enum для status
CREATE TYPE users_status AS ENUM ('pending_confirmation', 'active', 'blocked');

-- Таблица администраторов
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    status users_status NOT NULL DEFAULT 'pending_confirmation',
    terms_accepted_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    activated_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_users_email ON users (LOWER(email));
CREATE INDEX idx_users_status ON users (status);

-- Trigger для updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();


-- Таблица confirmation_tokens 
CREATE TABLE confirmation_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    used_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT check_token_active CHECK (used_at IS NULL OR expires_at > NOW())
);

CREATE INDEX idx_confirmation_tokens_admin_id ON confirmation_tokens (user_id);
CREATE INDEX idx_confirmation_tokens_token ON confirmation_tokens (LOWER(token));
CREATE INDEX idx_confirmation_tokens_expires_at ON confirmation_tokens (expires_at) WHERE used_at IS NULL;
