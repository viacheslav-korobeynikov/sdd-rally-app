CREATE TABLE confirmation_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    admin_id UUID NOT NULL REFERENCES admins(id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    used_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT check_token_active CHECK (used_at IS NULL OR expires_at > NOW())
);

CREATE INDEX idx_confirmation_tokens_admin_id ON confirmation_tokens (admin_id);
CREATE INDEX idx_confirmation_tokens_token ON confirmation_tokens (LOWER(token));
CREATE INDEX idx_confirmation_tokens_expires_at ON confirmation_tokens (expires_at) WHERE used_at IS NULL;
