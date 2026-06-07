CREATE SCHEMA IF NOT EXISTS auth;

CREATE TABLE IF NOT EXISTS auth.credentials (
    user_id BIGINT PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT,
    is_active BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS auth.verifications (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    verified_token TEXT NOT NULL UNIQUE,
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    verified_at TIMESTAMPTZ,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS auth.password_histories (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    old_password_hash TEXT NOT NULL,
    reason_status TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS auth.sessions (
    device_id TEXT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    public_key TEXT,
    logged_out_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
