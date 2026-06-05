# Go Template Project

A secure authentication system built with Go, providing comprehensive authentication features with advanced security measures.

## 🚀 Key Features

- **User Authentication**: Registration, login with JWT tokens
- **Two-Factor Authentication (2FA)**: OTP via email
- **Session Management**: Session management with Redis
- **Security**: Rate limiting, CSRF protection, IP blacklist
- **Social Login**: Firebase Authentication integration
- **Message Queue**: RabbitMQ for asynchronous processing
- **Cron Jobs**: Automated scheduled tasks
- **API Documentation**: Swagger/OpenAPI
- **Database**: PostgreSQL with migrations
- **Containerization**: Docker support

## 📁 Directory Structure

```
Go_Secure_Auth_Pro/
├── cmd/                          # Entry points
│   ├── server/                   # Main API server
│   ├── cronjob/                  # Cron job service
│   ├── queue/                    # Message queue consumer
│   └── cli/                      # Command line tools
├── internal/                     # Private application code
│   ├── controllers/              # HTTP handlers
│   ├── middlewares/              # HTTP middlewares
│   ├── models/                   # Data models
│   ├── repo/                     # Repository layer
│   ├── service/                  # Business logic
│   ├── routers/                  # Route definitions
│   └── messaging/                # Message queue handlers
├── pkg/                          # Public libraries
│   ├── helpers/                  # Utility functions
│   ├── mail/                     # Email services
│   └── setting/                  # Configuration helpers
├── configs/                      # Configuration files
├── migrations/                   # Database migrations
├── docs/                         # Documentation & Swagger
├── response/                     # Response handlers
├── third_party/                  # External integrations
│   ├── docker/                   # Docker configurations
│   ├── firebase/                 # Firebase config
│   └── telegram/                 # Telegram integration
├── templates/                    # HTML templates
└── tests/                        # Test files
```

## 🛠️ Technologies Used

- **Backend**: Go 1.22.3
- **Web Framework**: Gin
- **Database**: PostgreSQL
- **Cache**: Redis
- **Message Queue**: RabbitMQ
- **Authentication**: JWT, Firebase Auth
- **Documentation**: Swagger
- **Containerization**: Docker & Docker Compose
- **Process Management**: Air (hot reload)

## 📋 System Requirements

- Go 1.22.3 or higher
- Docker & Docker Compose
- PostgreSQL
- Redis
- RabbitMQ

## ⚡ Installation and Setup

### 1. Clone the Project

```bash
git clone https://github.com/TuanChill/Go-Backend-Template-.git
cd Go-Backend-Template-
```

### 2. Environment Configuration

```bash
# Copy the example configuration file
cp .env.example .env

# Edit the environment variables in the .env file
# Required configurations:
# - Database credentials
# - Redis connection
# - RabbitMQ settings
# - SMTP settings for email
# - Firebase credentials
```

### 3. Development Mode

#### Using Docker (Recommended)

```bash
# Start services (PostgreSQL, Redis, RabbitMQ)
make build-dev

# Run application with hot reload
make dev
# or
make air
```

#### Run Directly

```bash
# Install dependencies
go mod download

# Run server
make start

# Run cron job (in another terminal)
make cron

# Run message queue consumer (in another terminal)
make consumer
```

### 4. Production Mode

```bash
# Build and run all services with Docker
make build-pro
```

## 🔧 Useful Makefile Commands

### Development

```bash
make start          # Run production server
make dev            # Run development server with file watcher
make air            # Run with Air hot reload
make cron           # Run cron job service
make consumer       # Run message queue consumer
```

### Docker

```bash
make build-dev      # Build development environment
make build-pro      # Build production environment
make down-dev       # Stop development containers
make down-pro       # Stop production containers
```

### Database & Documentation

```bash
make sqlc           # Generate SQLC code
make swagger        # Generate Swagger documentation
```

## 📚 API Documentation

After running the server, access Swagger UI at:

- Development: `http://localhost:PORT/swagger/index.html`
- Production: `http://your-domain/swagger/index.html`

## 🔐 Security

The project integrates multiple security layers:

- **Rate Limiting**: Request rate limiting
- **CSRF Protection**: Cross-Site Request Forgery protection
- **Input Sanitization**: Input data sanitization
- **IP Blacklist**: Malicious IP blocking
- **JWT Security**: Token-based authentication
- **Password Hashing**: Secure password encryption
- **Request Size Limiting**: Request size limitation

## 🗄️ Database

### Migrations

```bash
# Migrations are automatically run when starting PostgreSQL container
# Migration files are located in the migrations/ directory
```

### Main Schema

- `users`: User information
- `password_history`: Password history
- `devices`: Login devices
- `otp`: OTP codes
- `verifications`: Email/phone verification
- `social_logins`: Social login

## 🚀 Deployment

### Docker Hub Images

```bash
# Build and push images
make build-and-push-all

# Update production images
make update-image
```

## 🤝 Contributing

1. Fork the project
2. Create a feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Create a Pull Request

## 🙏 Acknowledgments

- [Gin Web Framework](https://gin-gonic.com/)
- [Firebase](https://firebase.google.com/)
- [PostgreSQL](https://www.postgresql.org/)
- [Redis](https://redis.io/)
- [RabbitMQ](https://www.rabbitmq.com/)

## Project Notes

This project focuses on authentication and security features for Go backend services. It is intended as a reusable backend template for production-oriented systems.

<!-- go run gif -->
<div align="center">
  <a href="https://go.dev/"><img src="https://raw.githubusercontent.com/TuanChill/TuanChill/main/assets/go_run.gif"></a>
</div>

---

<!-- go run gif -->

## 🗂 **Folder Structure**

```plaintext
.
├── .dockerignore
├── .env
├── .env.example
├── .gitignore
├── .vscode/
│   └── settings.json
├── cmd/
│   ├── cli/
│   ├── cronjob/
│   ├── queue/
│   └── server/
├── configs/
│   ├── common/
│   ├── config.go
│   └── yaml/
├── docker-compose.dev.yml
├── docker-compose.pro.yml
├── docs/
│   ├── assets/
│   ├── CODE.md
│   ├── CODETABLE.md
│   ├── GO.md
│   ├── postman/
│   └── swagger/
├── fsnotify.go
├── global/
├── go.mod
├── go.sum
├── internal/
│   ├── controllers/
│   ├── messaging/
│   ├── middlewares/
│   ├── models/
│   └── repo/
├── makefile
├── migrations/
├── pkg/
├── response/
├── scripts/
├── sqlc.yaml
├── templates/
├── tests/
├── third_party/
└── tmp/
```




- `.dockerignore`: Contains a list of files and directories that Docker should ignore when building an image.
- `.env`: Contains environment variables for the project.
- `.env.example`: An example `.env` file containing necessary environment variables, meant to guide setup.
- `.github/`: Contains configuration files for GitHub, like `FUNDING.yml` for sponsorship settings.
- `.gitignore`: Contains a list of files and directories that git should ignore.
- `.vscode/`: Contains configurations for Visual Studio Code, such as `settings.json`.
- `cmd/`: Contains the application's entry points like CLI, cronjob, queue, and server.
- `configs/`: Contains configuration files for the application, including common configurations and configurations in YAML format.
- `docker-compose.dev.yml` and `docker-compose.pro.yml`: Contain Docker Compose configurations for development and production environments.
- `docs/`: Contains project documentation, including coding standards, code tables, Go guidelines, Postman collections, and Swagger files.
- `fsnotify.go`: This file may contain code to monitor file system changes.
- `global/`: Contains global variables for the application.
- `go.mod` and `go.sum`: Manage the project's Go dependencies.
- `GUILD.md`: May contain guidelines or information on how to join and contribute to the project.
- `internal/`: Contains the application's internal source code, not intended for external reuse.
- `makefile`: Contains automation commands for building and managing the project.
- `migrations/`: Contains database migration files.
- `pkg/`: Contains libraries and packages that can be reused outside the project.
- `README.md`: This file contains an overview and instructions for the project.
- `response/`: May contain code for handling and returning HTTP responses.
- `scripts/`: Contains support scripts for development and deployment.
- `sqlc.yaml`: Configuration for sqlc, a tool for generating code from SQL.
- `templates/`: Contains templates for user interfaces or other files.
- `tests/`: Contains automated tests for the project.
- `third_party/`: Contains code from third-party projects.
- `tmp/`: A temporary directory for files created during development.
