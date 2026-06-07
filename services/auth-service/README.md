# Auth Service

Extracted service for credentials, registration, login, refresh, verification, reset, and session/device ownership.

## Scope

- `POST /internal/auth/register`
- `POST /internal/auth/login`
- `POST /internal/auth/refresh`
- `POST /internal/auth/verify-account`
- `POST /internal/auth/password-reset`

Firebase and profile/email delivery behavior are intentionally excluded from this service boundary.

Only registration is enabled in this foundation slice. Login, refresh, account verification, and password reset fail closed until durable credential, token, and password-history validation are implemented.

## Run

```bash
go test ./...
SERVICE_TOKEN=service-token go run ./cmd/server
```

The service listens on `PORT` or `8000` by default. Set `SERVICE_TOKEN` for internal bearer-token authentication.
