-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id          BIGSERIAL   PRIMARY KEY,
    email       TEXT        UNIQUE NOT NULL,
    password_hash TEXT      NOT NULL
);

CREATE TABLE IF NOT EXISTS todos (
    id          BIGSERIAL   PRIMARY KEY,
    item        TEXT        NOT NULL,
    completed   BOOLEAN     NOT NULL DEFAULT FALSE,
    user_id     BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS blacklisted_tokens (
    token       TEXT        PRIMARY KEY,
    expired_at  TIMESTAMPTZ NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS blacklisted_tokens;
DROP TABLE IF EXISTS todos;
DROP TABLE IF EXISTS users;
