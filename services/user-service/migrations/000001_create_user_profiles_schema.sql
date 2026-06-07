CREATE SCHEMA IF NOT EXISTS user_profiles;

CREATE TABLE IF NOT EXISTS user_profiles.users (
    id BIGSERIAL PRIMARY KEY,
    username TEXT,
    email TEXT NOT NULL UNIQUE,
    phone TEXT,
    hidden_phone_number TEXT,
    fullname TEXT,
    hidden_email TEXT,
    avatar TEXT,
    gender INTEGER NOT NULL DEFAULT 0,
    two_factor_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
