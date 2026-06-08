# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Go 1.24 microservice template in transition from an auth-heavy backend to a lean service foundation. The current default direction is to keep infrastructure/runtime scaffolding, make auth/Firebase optional or isolated over time, and prefer explicit dependency wiring over package-level globals.

Stack: Gin, PostgreSQL, Redis, RabbitMQ, Viper/godotenv, Swagger, Docker Compose, and starter Kubernetes manifests.

## Common Commands

```bash
# Setup and local services
cp .env.example .env
make compose-up                 # start PostgreSQL, Redis, RabbitMQ, and API from docker-compose.yml
make compose-down               # stop local compose stack
make compose-config             # validate docker-compose.yml

# Current transitional commands
make api                        # run cmd/api with explicit bootstrap and health routes
make worker                     # run cmd/worker bootstrap placeholder
make migrate                    # run cmd/migrate bootstrap placeholder

# Legacy commands still present during transition
make start                      # run legacy cmd/server
make dev                        # run legacy fsnotify hot reload wrapper
make air                        # run Air hot reload if installed
make cron                       # run legacy cmd/cronjob
make consumer                   # run legacy cmd/queue RabbitMQ consumer

# Code generation
make sqlc                       # generate sqlc code from migrations/query into migrations/repo
make swagger                    # generate Swagger docs into docs/swagger

# Validation
go test ./...                   # run all tests
go test -race ./...             # run all tests with race detector
go test -cover ./...            # run all tests with coverage
go test ./internal/health -run TestName  # run a single package/test
go vet ./...                    # static analysis
docker compose config             # validates using values from .env
docker build --build-arg APP_CMD=api -t go-service-api:local .
```

Microservice strangler stack:

```bash
SERVICE_TOKEN=local-dev-token docker compose -f docker-compose.microservices.yml config
SERVICE_TOKEN=local-dev-token docker compose -f docker-compose.microservices.yml up --build
```

## Architecture

### Transitional runtime paths

- `cmd/api` is the preferred new API entry point. It calls `internal/app.Bootstrap`, registers only health routes from `internal/health`, and owns explicit shutdown.
- `cmd/worker` and `cmd/migrate` also use `internal/app.Bootstrap`, but currently only verify bootstrap; real queue consumption and schema execution remain in legacy flows.
- `cmd/server`, `cmd/queue`, and `cmd/cronjob` are legacy entry points kept while the template transition stabilizes.

### New foundation packages

- `internal/app` owns runtime bootstrap and shutdown orchestration.
- `internal/config` loads and validates runtime config.
- `internal/platform/*` contains adapters for PostgreSQL, Redis, RabbitMQ, logging, notification, and user-related integrations.
- `internal/health` exposes `/health/live` and `/health/ready`.

### Legacy API path

The legacy API still uses:

```text
cmd/server/main.go → internal/routers → controllers → services/repositories → global connections
```

`global/global.go` initializes config, PostgreSQL, Redis, Firebase, RabbitMQ, AWS SES, and AWS S3 in package `init()`. Avoid expanding this pattern in new code; prefer explicit bootstrap and dependency wiring under `internal/app` and `internal/platform`.

Legacy routes are grouped under `/v1` in `internal/routers/router.go`; `/ping` and Swagger at `/docs/swagger/*any` are also registered there.

### Extracted service modules

Nested service modules live under:

```text
services/api-gateway
services/auth-service
services/user-service
services/notification-otp-service
contracts
```

`docker-compose.microservices.yml` runs the strangler stack. The gateway routes `/v1/auth/*`, `/v1/user/*`, and `/v1/otp/*` to extracted services and keeps the legacy app as fallback for non-migrated routes. Do not remove legacy auth/user/OTP paths until gateway metrics show no fallback usage for migrated flows during the approved deprecation window.

## Data and Generated Code

- SQL migrations and SQLC queries live under `migrations/`.
- `sqlc.yaml` generates package `migrations` into `migrations/repo` with interfaces and JSON tags.
- Swagger generation uses `cmd/server/main.go` as the current `swag init` entry point and writes to `docs/swagger`.

## Configuration

- Local env starts from `.env.example`; the makefile includes `.env` directly.
- Main compose stack reads `.env.example` for the API container but infrastructure services use environment values such as `POSTGRES_DB`, `POSTGRES_USER`, `POSTGRES_PASSWORD`, `RABBIT_USER`, and `RABBIT_PASSWORD`.
- Legacy config loads from `configs/yaml` through Viper and is exposed as `global.Cfg`.

## Working Guidelines

- Match the transitional direction in `README.md`: new runtime code should move toward explicit bootstrap and dependency wiring, not new package-level global initialization.
- Keep legacy compatibility unless a task explicitly removes it; old entry points and routes are still part of the transition plan.
- Use existing response helpers and error codes in `response/` when editing legacy controllers/services.
- For endpoint work, check whether the change belongs in the new extracted service stack, the legacy router, or both during the transition.
