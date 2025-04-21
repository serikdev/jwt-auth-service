-- +goose Up

CREATE TABLE IF NOT EXISTS sessions (
    id UUID PRIMARY KEY,
    user_guid UUID NOT NULL,
    refresh_hash TEXT NOT NULL,
    jti TEXT NOT NULL UNIQUE,
    ip_address TEXT NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    used BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS sessions;


