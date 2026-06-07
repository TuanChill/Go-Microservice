# Notification OTP Service

Extracted service for OTP workflow mechanics and notification scheduling.

## Scope

- `POST /internal/otp/request`
- `POST /internal/otp/verify`
- `POST /internal/notifications/email-verification`
- `POST /internal/notifications/password-reset`
- Idempotent event handling scaffold for RabbitMQ consumers

Current idempotency storage is in-memory for the first contract-testable slice; production cutover needs durable storage in `notification.email_logs` or an equivalent table.

Auth remains owner of credential, verification eligibility, reset authorization, lockout, and token/session decisions.

## Run

```bash
go test ./...
go run ./cmd/server
```

The service listens on `PORT` or `8000` by default. Set `SERVICE_TOKEN` for internal bearer-token authentication.
