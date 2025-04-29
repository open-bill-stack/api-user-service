BEGIN;

CREATE EXTENSION IF NOT EXISTS "pgcrypto"; -- для gen_random_uuid()

CREATE TABLE users
(
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email             VARCHAR(255) UNIQUE NOT NULL,
    phone             VARCHAR(20)         NULL,
    is_email_verified BOOLEAN          DEFAULT FALSE,
    is_phone_verified BOOLEAN          DEFAULT FALSE,
    created_at        TIMESTAMPTZ      DEFAULT now(),
    updated_at        TIMESTAMPTZ      DEFAULT now()
);

CREATE INDEX idx_users_phone ON users (phone);

CREATE TABLE user_metadata
(
    user_id       UUID PRIMARY KEY REFERENCES users (id) ON DELETE CASCADE,
    last_login_at TIMESTAMPTZ,
    last_login_ip INET,
    timezone      VARCHAR(50),
    locale        VARCHAR(10)
);

CREATE TABLE credentials
(
    user_id       UUID PRIMARY KEY REFERENCES users (id) ON DELETE CASCADE,
    password_hash TEXT NOT NULL,
    created_at    TIMESTAMPTZ DEFAULT now()
);

COMMIT;