# Getting Started

Step-by-step guide to set up, run, and extend this template.

---

## Prerequisites

| Tool | Version | Purpose |
|---|---|---|
| Go | 1.22.3+ | Runtime |
| Docker + Compose | any | PostgreSQL, Redis, RabbitMQ |
| `make` | any | Task runner |
| `golangci-lint` | latest | Linting |
| `gosec` | latest | Security scan |

Optional but recommended:

```bash
go install github.com/air-verse/air@latest         # hot reload
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest # code gen
go install github.com/swaggo/swag/cmd/swag@latest   # swagger gen
```

---

## 1. Clone & Configure

```bash
git clone <repo-url>
cd go-backend-template

cp .env.example .env
# Edit .env — fill in DB, Redis, RabbitMQ, SMTP credentials
```

Key `.env` values:

```
ENV=dev
PORT=8080
PORT_FRONTEND=http://localhost:3000

POSTGRES_DB=mydb
POSTGRES_USER=postgres
POSTGRES_PASSWORD=secret
POSTGRES_HOST=localhost
POSTGRES_PORT=5432

REDIS_HOST=localhost
REDIS_PORT=6379

RABBIT_URL=amqp://guest:guest@localhost:5672/
```

Config YAML lives in `configs/yaml/config.dev.yml` — edit server, JWT, and SMTP settings there.

---

## 2. Start Infrastructure

```bash
make build-dev    # starts PostgreSQL, Redis, RabbitMQ via Docker Compose
```

Verify containers are up:

```bash
docker ps
```

---

## 3. Run Migrations

```bash
psql -U postgres -d mydb -f migrations/init/1_create_table_user.sql
psql -U postgres -d mydb -f migrations/init/2_create_table_password_history.sql
psql -U postgres -d mydb -f migrations/init/3_create_table_devices.sql
psql -U postgres -d mydb -f migrations/init/4_create_table_social_logins.sql
psql -U postgres -d mydb -f migrations/init/5_create_table_otp.sql
psql -U postgres -d mydb -f migrations/init/6_create_table_verifications.sql
```

---

## 4. Start the Server

```bash
make dev      # hot reload via fsnotify (recommended for development)
make air      # hot reload via Air (alternative)
make start    # production build, no reload
```

Health check:

```bash
curl http://localhost:8080/ping
# → {"message":"pong"}
```

Swagger UI: `http://localhost:8080/docs/swagger/index.html`

---

## 5. Start Background Services (optional)

```bash
make cron      # scheduled jobs (cmd/cronjob/main.go)
make consumer  # RabbitMQ consumer (cmd/queue/main.go)
```

---

## Adding a New Feature

Follow the layered pattern: **model → repo → service → controller → router**.

### Step 1 — Define the model (`internal/models/`)

```go
// internal/models/postModel.go
type Post struct {
    ID        int       `json:"id"`
    Title     string    `json:"title"`
    UserID    int       `json:"user_id"`
    CreatedAt time.Time `json:"created_at"`
}

type CreatePostRequest struct {
    Title string `json:"title" binding:"required,min=1,max=255"`
}
```

### Step 2 — Write the repository (`internal/repo/`)

```go
// internal/repo/postRepo.go
func CreatePost(db *sql.DB, userID int, title string) (models.Post, error) {
    row := db.QueryRow(
        "INSERT INTO posts (user_id, title) VALUES ($1, $2) RETURNING id, title, user_id, created_at",
        userID, title,
    )
    var p models.Post
    err := row.Scan(&p.ID, &p.Title, &p.UserID, &p.CreatedAt)
    return p, err
}
```

### Step 3 — Write the service (`internal/service/`)

```go
// internal/service/postService.go
func CreatePost(c *gin.Context) *models.Post {
    var req models.CreatePostRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.BadRequestError(c, response.ErrCodeInvalidFormat)
        return nil
    }

    userID, _ := c.Get("user_id") // set by AuthorizationMiddleware

    post, err := repo.CreatePost(global.DB, userID.(int), req.Title)
    if err != nil {
        response.InternalServerError(c, response.ErrCodeDBQuery)
        return nil
    }
    return &post
}
```

### Step 4 — Write the controller (`internal/controllers/`)

```go
// internal/controllers/postController.go
func CreatePost(c *gin.Context) error {
    result := service.CreatePost(c)
    if result == nil {
        return nil
    }
    response.Created(c, "Create Post", result)
    return nil
}
```

### Step 5 — Register the route (`internal/routers/router.go`)

```go
post := v1.Group("/post")
{
    post.Use(middlewares.AuthorizationMiddleware())
    post.POST("/", utils.AsyncHandler(controller.CreatePost))
}
```

### Step 6 — Add error codes if needed (`response/customErrorCode.go`)

```go
ErrCodePostNotFound = 5100
ErrCodePostForbidden = 5101
```

---

## Adding a Migration

Create a numbered SQL file in `migrations/init/`:

```bash
touch migrations/init/7_create_table_posts.sql
```

```sql
CREATE TABLE IF NOT EXISTS posts (
    id         SERIAL PRIMARY KEY,
    user_id    INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title      VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_posts_user_id ON posts(user_id);
```

Apply:

```bash
psql -U postgres -d mydb -f migrations/init/7_create_table_posts.sql
```

---

## Adding a Middleware

```go
// internal/middlewares/myMiddleware.go
func MyMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // pre-processing
        c.Next()
        // post-processing
    }
}
```

Register in `routers/router.go`:

```go
r.Use(middlewares.MyMiddleware())   // global
// or
user.Use(middlewares.MyMiddleware()) // route group only
```

---

## Adding a RabbitMQ Event

**Publish** (from service):

```go
import "go_template/internal/messaging"

messaging.Publish(global.MessageQueue, "post.created", PostCreatedEvent{
    PostID: post.ID,
    UserID: post.UserID,
})
```

**Consume** (in `internal/messaging/consumer.go` or new file):

```go
func HandlePostCreated(event PostCreatedEvent) error {
    // send notification, update search index, etc.
    return nil
}
```

Register the handler in `cmd/queue/main.go`.

---

## Code Generation

```bash
make sqlc     # regenerate repo code from migrations/query/*.sql
make swagger  # regenerate docs/swagger/ from controller annotations
```

---

## Linting & Tests

```bash
golangci-lint run ./...   # lint
gosec ./...               # security scan
go test -race ./...       # tests with race detector
go test -cover ./...      # coverage report
```

Run before every commit.

---

## Production Build

```bash
make build-pro   # builds and starts all containers (API + workers)
make down-dev    # stop dev containers
```

---

## Environment Reference

| Variable | Description |
|---|---|
| `ENV` | `dev` or `prod` — selects config YAML |
| `PORT` | HTTP listen port |
| `PORT_FRONTEND` | Allowed CORS origin |
| `POSTGRES_*` | PostgreSQL connection |
| `REDIS_*` | Redis connection |
| `RABBIT_URL` | RabbitMQ AMQP URL |
| `SMTP_*` | Email sender config |
| `RANDOM_PASSWORD` | Salt for internal hashing |
| `CSRF_TOKEN` | CSRF secret (if enabled) |
