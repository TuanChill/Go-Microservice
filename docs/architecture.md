# System Architecture

## Overview

Go backend template providing secure authentication with JWT, OTP, two-factor auth, social login, and device management. Built on Gin, PostgreSQL, Redis, RabbitMQ, and Firebase.

---

## Stack

| Layer | Technology |
|---|---|
| HTTP Framework | Gin |
| Database | PostgreSQL (via `lib/pq`) |
| Cache | Redis (`go-redis/v9`) |
| Message Queue | RabbitMQ (`amqp091-go`) |
| Auth Provider | Firebase Admin SDK |
| Config | Viper + `.env` (godotenv) |
| Code Generation | SQLC |
| API Docs | Swagger (swaggo) |
| Hot Reload | Air / fsnotify |

---

## Directory Structure

```
.
├── cmd/
│   ├── server/main.go       # HTTP API server entry point
│   ├── cronjob/main.go      # Scheduled tasks entry point
│   ├── queue/main.go        # RabbitMQ consumer entry point
│   └── cli/                 # CLI tools
│
├── configs/
│   ├── config.go            # Viper config loader
│   ├── yaml/
│   │   ├── config.dev.yml
│   │   └── config.prod.yml
│   └── common/
│       ├── constants/       # App-wide constants, RabbitMQ keys
│       └── utils/           # Shared utilities (AsyncHandler, etc.)
│
├── global/
│   └── global.go            # Singleton connections: DB, Cache, AdminSdk, MessageQueue, Cfg
│
├── internal/
│   ├── controllers/         # HTTP handlers — validate input, call service, return response
│   │   └── initialization/  # DB/Redis/RabbitMQ connection factories
│   ├── middlewares/         # Gin middleware chain (auth, rate limit, CORS, etc.)
│   ├── models/              # Request/response structs, DB models
│   ├── repo/                # Data access layer — SQL queries, Redis ops
│   │   └── redis/           # Redis-specific repo operations
│   ├── routers/             # Route registration and middleware wiring
│   ├── service/             # Business logic layer
│   └── messaging/           # RabbitMQ producer and consumer
│
├── migrations/
│   ├── init/                # DDL migrations (run once)
│   └── query/               # SQLC query files
│
├── pkg/
│   ├── helpers/             # Shared helpers (JWT, hash, validate)
│   ├── mail/                # SMTP email sender
│   └── setting/             # Firebase SDK init
│
├── response/                # Unified HTTP response helpers and error codes
├── third_party/
│   └── telegram/            # Telegram bot integration
├── tests/                   # Integration/manual test scripts
└── docs/
    └── swagger/             # Auto-generated Swagger docs
```

---

## Request Lifecycle

```
Client Request
     │
     ▼
Middleware Chain (in order):
  1. IPBlackList          — block banned IPs
  2. CORSMiddleware        — cross-origin headers
  3. SecurityHeaders       — Helmet-style headers
  4. HeadersMiddlewares    — custom request headers
  5. RequestSizeLimiter    — max 1 MB body
  6. RateLimiter           — 5 req/s, burst 10
  7. RequestLogging        — structured access log
  8. PathTraversal         — block `../` attacks
  9. ContentTypeValidation — enforce JSON
 10. SanitizeParams        — strip dangerous chars
     │
     ▼
Router (routers/router.go)
     │
     ▼
Controller (internal/controllers/)
  • Bind & validate request
  • Call service function
  • Call response helper
     │
     ▼
Service (internal/service/)
  • Business logic
  • Cache-aside (Redis)
  • Call repo for DB ops
  • Publish MQ events
     │
     ├── Repository (internal/repo/)
     │     • Raw SQL via database/sql
     │     • Returns (result, error)
     │
     └── Cache (internal/repo/redis/)
           • HGetAll / HSet / Expire
```

---

## Global Connections (`global/global.go`)

Initialized once via `init()` at process start. Panic on failure — fail fast.

```go
global.DB           // *sql.DB       — PostgreSQL
global.Cache        // *redis.Client — Redis
global.AdminSdk     // *firebase.App — Firebase
global.MessageQueue // *amqp.Connection — RabbitMQ
global.Cfg          // models.Config — loaded from YAML + .env
```

---

## Layered Rules

| Layer | Allowed imports | Forbidden |
|---|---|---|
| Controller | service, response, gin | repo, global.DB |
| Service | repo, global, models, pkg | gin (HTTP types) |
| Repository | database/sql, global.DB, models | service, gin |
| Middleware | response, global, pkg | service, repo |

---

## Auth Flow

```
Register → email verification link → VerificationAccount
Login    → JWT access token (header) + refresh token (cookie) + device_id (header)
Refresh  → RefetchTokenMiddleware validates cookie → /renew-token
Logout   → revokes device session from DB + clears Redis cache
2FA      → EnableTwoFactor → SendOtpUpdateEmail → VerifyOtp
Social   → Firebase ID token → /login-social
```

**Token delivery:**
- `Authorization: Bearer <access_token>` — request header
- `user_login` cookie — refresh token (httpOnly)
- `X-Device-Id` header — required on protected routes

---

## Cache Strategy

Cache-aside pattern per user profile:

```
Read:  HGetAll(cacheKey) → hit? return cached : query DB → HSet cache → return
Write: update DB → Del(cacheKey)
```

Cache key format defined in `configs/common/constants/`.

---

## Message Queue

| Event | Exchange/Queue | Direction |
|---|---|---|
| Email send | `mail.*` | producer → consumer |
| Account events | `user.*` | producer → consumer |

Consumer runs as a separate process (`cmd/queue/`). Producer lives inside service layer.

---

## Database Schema (migrations/init/)

| Table | Purpose |
|---|---|
| `users` | Core user record |
| `password_history` | Prevent password reuse |
| `devices` | Per-device session tracking |
| `social_logins` | OAuth provider links |
| `otp` | One-time passwords |
| `verifications` | Email verification tokens |

---

## Error Code Ranges (`response/customErrorCode.go`)

| Range | Domain |
|---|---|
| 1000s | General |
| 2000s | Database |
| 3000s | Validation |
| 4000s | Auth |
| 5000s+ | Resource / external / user-specific |

Full table: `docs/CODETABLE.md`
