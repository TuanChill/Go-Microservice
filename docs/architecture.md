# System Architecture

## Overview

This repository is transitioning into a lean Go microservice template. The target is a reusable single-service foundation with explicit runtime bootstrap, optional infrastructure adapters, transitional entrypoints, and starter deployment assets.

## Current State

The old auth application code still exists during transition. New template-facing code is additive and lives in:

```text
cmd/api
cmd/worker
cmd/migrate
internal/app
internal/config
internal/platform
internal/health
```

Legacy paths remain temporarily:

```text
cmd/server
cmd/queue
cmd/cronjob
global
internal/controllers
internal/service
internal/repo
```

## Target Runtime Flow

```text
cmd/<runtime>/main.go
  → internal/config.Load
  → internal/app.Bootstrap
  → internal/platform adapters
  → runtime-specific server/worker/migrate behavior
  → internal/app.Shutdown
```

## Runtime Entrypoints

| Entrypoint | Status | Purpose |
|---|---|---|
| `cmd/api` | transitional | Starts API with new bootstrap plus existing router. |
| `cmd/worker` | transitional | Boots worker dependencies; real queue consumer remains `cmd/queue`. |
| `cmd/migrate` | transitional | Boots migration dependencies; schema execution remains existing migration workflow. |
| `cmd/server` | legacy | Existing API server. |
| `cmd/queue` | legacy | Existing RabbitMQ consumer. |
| `cmd/cronjob` | legacy | Existing scheduled job process. |

## Infrastructure Adapters

```text
internal/platform/logger
internal/platform/postgres
internal/platform/redis
internal/platform/rabbitmq
```

Adapters expose explicit constructors and return errors to callers. They do not open network connections from package `init()`.

## Dependency Ownership

`internal/app.App` owns resources created through the new bootstrap path:

- config,
- logger,
- optional PostgreSQL connection,
- optional Redis client,
- optional RabbitMQ connection.

`App.Shutdown` closes owned resources and aggregates close errors.

## Health Endpoints

`internal/health` exposes:

- `/health/live`
- `/health/ready`

Current readiness is transitional. It supports dependency-aware callbacks, but the new API command currently registers a nil readiness callback while legacy router dependencies are still global.

## Deployment Assets

| Path | Purpose |
|---|---|
| `Dockerfile` | Multi-stage build for selected command via `APP_CMD`. |
| `docker-compose.yml` | Local stack: PostgreSQL, Redis, RabbitMQ, API. |
| `.env.example` | Safe placeholder local config. |
| `deployments/k8s/` | Starter Kubernetes manifests. |

Kubernetes files are starter manifests only. Production users must customize image registry, secrets backend, resource sizing, ingress/TLS, autoscaling, network policy, and worker probes.

## Lean Template Boundary

Firebase/auth are not default template core. They should be removed or isolated as optional material in later cleanup.
