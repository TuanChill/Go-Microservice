# Go Microservice Template

Lean Go service template for production-oriented systems. It provides a reusable single-service foundation with Gin, PostgreSQL, Redis, RabbitMQ, Docker, Kubernetes starter manifests, explicit runtime bootstrap, and health endpoints.

## Current Template Status

This repository is being refactored from an auth-heavy backend into a lean microservice template.

Default template direction:

- Keep infrastructure and runtime scaffolding.
- Keep auth/Firebase only as optional/non-default material during transition.
- Prefer explicit dependency wiring over package-level globals.
- Keep old entrypoints temporarily for compatibility while new commands stabilize.

## Stack

| Layer | Technology |
|---|---|
| Language | Go 1.24 |
| HTTP | Gin |
| Database | PostgreSQL |
| Cache | Redis |
| Queue | RabbitMQ |
| Config | Viper + godotenv |
| API Docs | Swagger |
| Containers | Docker + Docker Compose |
| Orchestration starter | Kubernetes manifests |

## Key Paths

```text
cmd/
  api/       # transitional API command
  worker/    # transitional worker bootstrap
  migrate/   # transitional migration bootstrap
  server/    # legacy API command
  queue/     # legacy RabbitMQ consumer
  cronjob/   # legacy cron command
internal/
  app/       # runtime bootstrap and shutdown ownership
  config/    # runtime config loading and validation
  platform/  # postgres, redis, rabbitmq, logger adapters
  health/    # /health/live and /health/ready handlers
  routers/   # legacy router wiring during transition
deployments/k8s/ # starter Kubernetes manifests
```

## Quick Start

```bash
cp .env.example .env
make compose-up
make api
```

Health endpoints on the new API command:

```bash
curl http://localhost:8000/health/live
curl http://localhost:8000/health/ready
```

Legacy app route still exists while transition is in progress:

```bash
curl http://localhost:8000/ping
```

## Commands

```bash
make api              # run transitional cmd/api
make worker           # run transitional worker bootstrap
make migrate          # run transitional migrate bootstrap
make start            # run legacy cmd/server
make consumer         # run legacy cmd/queue consumer
make cron             # run legacy cmd/cronjob
make compose-up       # start local stack
make compose-down     # stop local stack
make compose-config   # validate compose config
make docker-build-api # build local API image
```

## Validation

```bash
go test ./...
go test -race ./...
go vet ./...
docker compose --env-file .env.example config
docker build --build-arg APP_CMD=api -t go-service-api:local .
```

## Deployment Scaffolding

- `Dockerfile` builds a selected command via `APP_CMD`.
- `docker-compose.yml` is for local development.
- `deployments/k8s/*.yaml` are starter manifests, not a complete production platform.

Before production use, customize:

- image registry and tags,
- external secret management,
- resource requests/limits,
- ingress/TLS,
- autoscaling,
- network policies,
- worker health/heartbeat probes.

## Transitional Notes

- `cmd/api` still uses existing router wiring, and that router still imports legacy globals.
- `cmd/worker` currently bootstraps worker dependencies but real queue consumption remains in `cmd/queue`.
- `cmd/migrate` currently bootstraps migration dependencies but schema execution remains in the existing migration workflow.
- Firebase/auth are not default Lean template core and should be removed or isolated in later cleanup.
