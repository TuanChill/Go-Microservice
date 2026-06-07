# Getting Started

## Prerequisites

| Tool | Purpose |
|---|---|
| Go 1.24+ | Build and test service commands |
| Docker + Compose | Local PostgreSQL, Redis, RabbitMQ, API image |
| make | Command shortcuts |

Optional:

```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install github.com/swaggo/swag/cmd/swag@latest
```

## Configure

```bash
cp .env.example .env
```

Edit `.env` for local values. Placeholder secrets in `.env.example` must be replaced before any shared environment.

## Start Local Infrastructure

```bash
make compose-up
```

Validate compose without starting containers:

```bash
make compose-config
```

## Run Commands

New transitional commands:

```bash
make api
make worker
make migrate
```

Legacy commands still available during transition:

```bash
make start
make consumer
make cron
```

## Health Checks

```bash
curl http://localhost:8000/health/live
curl http://localhost:8000/health/ready
```

The readiness endpoint supports dependency-aware callbacks, but current API wiring is transitional.

## Run Microservice Decomposition Stack

The local strangler stack uses nested service modules at `services/api-gateway`, `services/auth-service`, `services/user-service`, `services/notification-otp-service`, and `contracts`.

```bash
SERVICE_TOKEN=local-dev-token docker compose -f docker-compose.microservices.yml config
SERVICE_TOKEN=local-dev-token docker compose -f docker-compose.microservices.yml up --build
```

The gateway listens on port `8080`, routes migrated `/v1/auth/*`, `/v1/user/*`, and `/v1/otp/*` traffic to extracted services, and keeps the legacy app as fallback for routes not migrated yet.

## Legacy Cleanup Guard

Do not delete legacy auth/user/OTP controllers, services, repos, models, or global runtime dependencies until staging or production gateway metrics prove migrated flows no longer hit legacy fallback during the approved deprecation window.

## Build Container Image

```bash
make docker-build-api
make docker-build-worker
```

Or directly:

```bash
docker build --build-arg APP_CMD=api -t go-service-api:local .
```

## Kubernetes Starter Manifests

Files live in `deployments/k8s/`:

```text
api-deployment.yaml
api-service.yaml
worker-deployment.yaml
configmap.yaml
secret.example.yaml
```

Before production use, customize:

- image registry and immutable tags,
- secret manager integration,
- resource requests/limits,
- ingress and TLS,
- autoscaling,
- network policies,
- worker health/heartbeat probe.

## Validation

Run before committing:

```bash
go test ./...
go test -race ./...
go vet ./...
docker compose --env-file .env.example config
docker build --build-arg APP_CMD=api -t go-service-api:local .
```

## Adding a Domain

Preferred lean template shape:

```text
internal/<domain>/handler.go
internal/<domain>/service.go
internal/<domain>/repository.go
internal/<domain>/model.go
internal/<domain>/dto.go
```

Keep handlers at transport boundaries, services for business logic, and repositories for data access. Avoid adding new package-level globals.
