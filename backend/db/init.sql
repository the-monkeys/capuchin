CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS todos (
    id UUID PRIMARY KEY,
    item TEXT NOT NULL,
    completed BOOLEAN DEFAULT FALSE,
    user_id UUID REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS blacklisted_tokens (
    token TEXT PRIMARY KEY,
    expired_at TIMESTAMP NOT NULL
);
