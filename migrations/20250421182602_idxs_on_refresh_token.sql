-- +goose Up
CREATE INDEX IF NOT EXISTS refresh_tokens_user_idx ON refresh_tokens(user_id);

CREATE INDEX IF NOT EXISTS refresh_tokens_token_hash_idx ON refresh_tokens(token_hash);

-- +goose Down
DROP INDEX IF EXISTS refresh_tokens_user_idx;
DROP INDEX IF EXISTS refresh_tokens_token_hash_idx;
