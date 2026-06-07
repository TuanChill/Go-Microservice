# User Service

Extracted service for profile and user-account metadata.

## Scope

- `POST /internal/users`
- `GET /internal/users/{id}`
- `PATCH /internal/users/{id}`
- `DELETE /internal/users/{id}`

This service does not own password verification, token issuing, refresh sessions, or login behavior.

## Run

```bash
go test ./...
SERVICE_TOKEN=service-token go run ./cmd/server
```

The service listens on `PORT` or `8000` by default. Set `SERVICE_TOKEN` for internal bearer-token authentication.
