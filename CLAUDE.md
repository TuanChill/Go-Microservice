# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Go backend template with authentication features. Built with Go 1.22.3, Gin web framework, PostgreSQL, Redis, and RabbitMQ.

## Common Commands

```bash
# Development
make dev          # Run server with fsnotify hot reload (monitors .go file changes)
make air          # Run with Air hot reload (alternative)
make start        # Run production server

# Background services
make cron         # Run cron job service
make consumer     # Run RabbitMQ message queue consumer

# Docker
make build-dev    # Start PostgreSQL, Redis, RabbitMQ containers
make down-dev     # Stop dev containers
make build-pro    # Build production containers

# Code generation
make sqlc         # Generate SQLC code from SQL queries
make swagger      # Generate Swagger documentation

# Linting
golangci-lint run ./...   # Run all linters (config: .golangci.yml)
go test -race ./...       # Run tests with race detector
go test -cover ./...      # Run tests with coverage
```

## Architecture

### Layered Structure
```
cmd/server/main.go → routers → controllers → services → repositories
                                           ↓
                                      global (DB, Cache, MQ)
```

### Entry Points (`cmd/`)
- `cmd/server/` - Main API server
- `cmd/cronjob/` - Scheduled tasks
- `cmd/queue/` - RabbitMQ consumer
- `cmd/cli/` - CLI tools

### Request Flow
1. Middleware chain (IP blacklist → CORS → Rate limiting → Logging → etc.)
2. Router groups route to controllers
3. Controllers call service layer
4. Services call repositories (DB) or cache (Redis)
5. Response helpers in `response/` package

### Global Connections (`global/global.go`)
Initialized in `init()`: PostgreSQL (`global.DB`), Redis (`global.Cache`), Firebase (`global.AdminSdk`), RabbitMQ (`global.MessageQueue`)

## Code Conventions

### Error Handling
Error codes are defined in `response/customErrorCode.go` with numeric prefixes:
- 1000s: General errors
- 2000s: Database errors
- 3000s: Validation errors
- 4000s: Auth errors
- 5000s+: Resource, external service, user-specific errors

Full table in `docs/CODETABLE.md`.

### Response Format
Use helpers from `response/` package:
```go
response.Ok(c, "Action Name", data)
response.Created(c, "Action Name", data)
response.BadRequestError(c, errorCode)
response.InternalServerError(c, errorCode)
```

### Service Pattern
Services receive `*gin.Context`, return `nil` on success (response sent internally), return `*Type` on error where caller handles response:
```go
func GetProfileUser(c *gin.Context) *models.ProfileResponseJSON {
    // ... logic
    if err != nil {
        response.InternalServerError(c, response.ErrCodeDBQuery)
        return nil  // error case
    }
    return &profileResponse  // success case - caller doesn't respond
}
```

### Repository Pattern
Repository functions take `*sql.DB` and params struct, return (result, error):
```go
func GetUserId(db *sql.DB, params GetUserIdParams) (*User, error)
```

### Cache Key Pattern
```go
fmt.Sprintf(constants.CacheProfileUser, strconv.Itoa(userId))
```

## Configuration

- Environment: `.env` (copy from `.env.example`)
- Config files: `configs/yaml/config.dev.yml` or `config.prod.yml`
- Config loader: `configs/config.go` using Viper
- Port from `global.Cfg.Server.Port`

## Dependencies

Module: `go_template`

Key dependencies:
- `github.com/gin-gonic/gin` - Web framework
- `github.com/lib/pq` - PostgreSQL driver
- `github.com/redis/go-redis/v9` - Redis client
- `github.com/rabbitmq/amqp091-go` - RabbitMQ
- `firebase.google.com/go` - Firebase Auth
- `github.com/swaggo/swag` - Swagger docs

## Go Coding Rules

### Formatting
- **gofmt** and **goimports** are mandatory — no style debates
- Run `go fmt ./...` before commits

### Design Principles
- Accept interfaces, return structs
- Keep interfaces small (1-3 methods)
- Define interfaces where they are used, not where they are implemented

### Error Handling
Always wrap errors with context:
```go
if err != nil {
    return fmt.Errorf("failed to create user: %w", err)
}
```

### Context & Timeouts
Always use `context.Context` for timeout control:
```go
ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
defer cancel()
```

### Security
- Use **gosec** for static security analysis: `gosec ./...`
- Never hardcode secrets — use environment variables

## Go Patterns

### Functional Options
Use for optional config on structs — avoids large constructor signatures:
```go
type Option func(*Server)

func WithTimeout(d time.Duration) Option {
    return func(s *Server) { s.timeout = d }
}

func NewServer(opts ...Option) *Server {
    s := &Server{timeout: 30 * time.Second}
    for _, opt := range opts {
        opt(s)
    }
    return s
}
```

### Dependency Injection via Constructors
```go
func NewUserService(repo UserRepository, cache CacheClient) *UserService {
    return &UserService{repo: repo, cache: cache}
}
```

### Table-Driven Tests
```go
func TestGetUser(t *testing.T) {
    tests := []struct {
        name    string
        id      int
        want    *User
        wantErr bool
    }{
        {"valid id", 1, &User{ID: 1}, false},
        {"not found", 999, nil, true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := GetUser(tt.id)
            if (err != nil) != tt.wantErr {
                t.Fatalf("err = %v, wantErr %v", err, tt.wantErr)
            }
            // assert got == tt.want
        })
    }
}
```

### Goroutine + Channel Pattern
Fan-out with bounded workers:
```go
func processItems(items []Item) []Result {
    ch := make(chan Result, len(items))
    sem := make(chan struct{}, 10) // max 10 concurrent

    var wg sync.WaitGroup
    for _, item := range items {
        wg.Add(1)
        go func(it Item) {
            defer wg.Done()
            sem <- struct{}{}
            defer func() { <-sem }()
            ch <- process(it)
        }(item)
    }

    go func() {
        wg.Wait()
        close(ch)
    }()

    var results []Result
    for r := range ch {
        results = append(results, r)
    }
    return results
}
```

### Graceful Shutdown
```go
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit

ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
if err := srv.Shutdown(ctx); err != nil {
    log.Fatal("forced shutdown:", err)
}
```

### Middleware Pattern (Gin)
```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            response.UnauthorizedError(c, response.ErrCodeUnauthorized)
            c.Abort()
            return
        }
        c.Next()
    }
}
```

### Sentinel Errors
```go
var (
    ErrNotFound   = errors.New("not found")
    ErrUnauthorized = errors.New("unauthorized")
)

// check with errors.Is
if errors.Is(err, ErrNotFound) { ... }
```

### Init Guards (fail fast at startup)
```go
func mustGetEnv(key string) string {
    v := os.Getenv(key)
    if v == "" {
        log.Fatalf("required env var %s not set", key)
    }
    return v
}
```

### Testing
```bash
go test -race ./...    # always run with race detector
go test -cover ./...   # verify coverage >= 80%
```

## Go Best Practices

### Package Design
- One purpose per package — name it by what it provides, not what it contains (`auth` not `authutils`)
- Avoid `util`, `common`, `helpers` packages — they become dumping grounds
- Internal packages (`internal/`) enforce import boundaries within the module
- `main` package only wires dependencies; no business logic

### Struct & Interface Design
- Embed interfaces for partial implementation (decorator pattern)
- Use struct embedding for code reuse, not inheritance
- Export only what external callers need — keep internals unexported
- Prefer value receivers unless mutation or large struct:
```go
// value receiver — immutable, safe to copy
func (u User) FullName() string { return u.First + " " + u.Last }

// pointer receiver — mutates or avoids copy of large struct
func (u *User) SetName(name string) { u.First = name }
```

### Concurrency Best Practices
- Never start goroutines without a way to stop them
- Pass context for cancellation — not a done channel
- Share memory by communicating; do not communicate by sharing memory
- Use `sync.Mutex` to protect shared state, `sync/atomic` for counters:
```go
type SafeCounter struct {
    mu sync.Mutex
    v  map[string]int
}

func (c *SafeCounter) Inc(key string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.v[key]++
}
```

### Resource Management
- Always `defer` cleanup immediately after acquiring a resource:
```go
f, err := os.Open(name)
if err != nil { return err }
defer f.Close()

rows, err := db.QueryContext(ctx, query)
if err != nil { return err }
defer rows.Close()
```
- Check `rows.Err()` after iterating SQL rows
- Close request bodies: `defer resp.Body.Close()`

### Logging
- Use structured logging (`log/slog` or `zap`) — never `fmt.Println` in production
- Include context fields: `user_id`, `request_id`, `trace_id`
- Log at the boundary (controller/handler), not deep in business logic

### Performance
- Pre-allocate slices when length is known: `make([]T, 0, n)`
- Use `strings.Builder` for string concatenation in loops
- Avoid interface boxing in hot paths
- Profile before optimizing: `go tool pprof`

## Design Patterns (Project-Specific)

### Clean Layered Architecture
```
Handler (controller) → validates input, calls service, returns HTTP response
Service              → business logic, orchestrates repos/cache
Repository           → data access only, no business logic
```
Rules:
- Services NEVER import `gin` — they are HTTP-agnostic
- Repositories NEVER call other repositories
- Controllers NEVER contain business logic

### Builder Pattern (for complex query params)
```go
type UserQuery struct {
    filters []string
    args    []any
}

func (q *UserQuery) WithEmail(email string) *UserQuery {
    q.filters = append(q.filters, "email = $"+strconv.Itoa(len(q.args)+1))
    q.args = append(q.args, email)
    return q
}
```

### Strategy Pattern (swappable algorithms)
```go
type Hasher interface {
    Hash(password string) (string, error)
    Verify(password, hash string) bool
}

// bcrypt, argon2 implement Hasher — swap without changing callers
type AuthService struct{ hasher Hasher }
```

### Observer / Event Pattern (RabbitMQ)
Producers publish events; consumers handle them independently:
```go
// producer — just emit, does not care about consumers
mq.Publish("user.created", UserCreatedEvent{ID: user.ID})

// consumer — registered handler, single responsibility
func HandleUserCreated(event UserCreatedEvent) error { ... }
```

### Retry with Backoff
```go
func withRetry(attempts int, delay time.Duration, fn func() error) error {
    for i := 0; i < attempts; i++ {
        if err := fn(); err == nil {
            return nil
        }
        time.Sleep(delay * time.Duration(i+1))
    }
    return fmt.Errorf("all %d attempts failed", attempts)
}
```

### Circuit Breaker (external service calls)
Wrap external HTTP/gRPC calls with a circuit breaker to prevent cascade failures. Open after N failures, half-open after timeout, close on success.

### Anti-Patterns to Avoid
| Anti-Pattern | Problem | Fix |
|---|---|---|
| `init()` side effects | Hidden order dependency | Explicit `Initialize()` calls in `main` |
| Returning `interface{}` | Loses type safety | Return concrete types |
| Naked `return` in long func | Hard to read | Always name return values explicitly |
| `panic` in library code | Crashes callers | Return `error` instead |
| Global mutable state | Race conditions | Inject as dependencies |
| Ignoring errors with `_` | Silent failures | Handle or explicitly document why |
