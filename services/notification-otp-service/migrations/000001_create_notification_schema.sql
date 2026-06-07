CREATE SCHEMA IF NOT EXISTS notification;

CREATE TABLE IF NOT EXISTS notification.otps (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    email TEXT NOT NULL,
    otp_code TEXT NOT NULL UNIQUE,
    purpose TEXT NOT NULL CHECK (purpose IN ('login', 'email_update', 'password_reset')),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS notification.email_logs (
    id BIGSERIAL PRIMARY KEY,
    event_id TEXT NOT NULL UNIQUE,
    idempotency_key TEXT NOT NULL,
    correlation_id TEXT NOT NULL,
    recipient_email TEXT NOT NULL,
    notification_type TEXT NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
