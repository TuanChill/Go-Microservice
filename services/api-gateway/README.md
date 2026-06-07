# API Gateway

Lightweight strangler gateway for routing public traffic to legacy or extracted services.

## Routes

- `/v1/auth/*` → `AUTH_SERVICE_URL`
- `/v1/user/*` → `USER_SERVICE_URL`
- `/v1/otp/*` → `NOTIFICATION_OTP_SERVICE_URL`
- unmatched paths → `LEGACY_URL`

The gateway routes only exact prefix bundles or their path children, strips client-supplied internal identity headers, and preserves or creates `X-Correlation-ID` / `X-Request-ID`.

## Run

```bash
go test ./...
go run ./cmd/server
```

Defaults: gateway on `:8080`, legacy `http://localhost:8000`, auth `:8001`, user `:8002`, notification OTP `:8003`.
