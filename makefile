################# TODO: CONSTANTS #################

#* GET FILE ENV
include .env
export $(shell sed 's/=.*//' .env)

# * FILE RUN GO
GO_API := ./cmd/api
GO_WORKER := ./cmd/worker
GO_MIGRATE := ./cmd/migrate
GO_SERVER_PRO := ./cmd/server/main.go
GO_SERVER_DEV:= ./fsnotify.go
GO_SERVER_CRON := ./cmd/cronjob/main.go
GO_CONSUMER := ./cmd/queue/main.go

# * DOCKER COMPOSE
DOCKER_COMPOSE_DEV := docker-compose.dev.yml
DOCKER_COMPOSE_PRO := docker-compose.pro.yml

# * DOCKER IMAGES
SERVER_IMAGE_NAME := go-service-api:local
CRON_IMAGE_NAME := go-service-cron:local
QUEUE_IMAGE_NAME := go-service-worker:local


# * DOCKER FILE
DOCKER_FILE_PATH := ./third_party/docker/go/Dockerfile
DOCKER_FILE_CRON_PATH := ./third_party/docker/go/Dockerfile-cron
DOCKER_FILE_QUEUE_PATH := ./third_party/docker/go/Dockerfile-queue

#* DOCKER CONTAINER
CONTAINER_SERVICE_AUTH := service_auth
CONTAINER_SERVICE_CRON := service_cron
CONTAINER_SERVICE_QUEUE := service_queue


# * FOLDER
SWAGGER_DIR=./docs/swagger

################# TODO: GO #################
start:
	go run $(GO_SERVER_PRO)

api:
	go run $(GO_API)

worker:
	go run $(GO_WORKER)

migrate:
	go run $(GO_MIGRATE)

dev:
	go run $(GO_SERVER_DEV)

air:
	air

cron:
	go run $(GO_SERVER_CRON)

consumer:
	go run $(GO_CONSUMER)
    
################# TODO: DOCKER #################
docker-build-api:
	docker build --build-arg APP_CMD=api -t go-service-api:local .

docker-build-worker:
	docker build --build-arg APP_CMD=worker -t go-service-worker:local .

compose-up:
	docker compose up -d --build

compose-down:
	docker compose down

compose-config:
	docker compose config

build-pro:
	docker-compose -f $(DOCKER_COMPOSE_PRO) up -d --build

down-pro:
	docker-compose -f $(DOCKER_COMPOSE_PRO) down

build-dev:
	docker-compose -f $(DOCKER_COMPOSE_DEV) up -d --build

down-dev:
	docker-compose -f $(DOCKER_COMPOSE_DEV) down

update-server:
	docker-compose -f $(DOCKER_COMPOSE_PRO) pull $(CONTAINER_SERVICE_AUTH)
	docker-compose -f $(DOCKER_COMPOSE_PRO) up -d --no-deps $(CONTAINER_SERVICE_AUTH)

update-cron:
	docker-compose -f $(DOCKER_COMPOSE_PRO) pull $(CONTAINER_SERVICE_QUEUE)
	docker-compose -f $(DOCKER_COMPOSE_PRO) up -d --no-deps $(CONTAINER_SERVICE_QUEUE)

update-image: update-server update-cron
	@echo "Both server and cron images updated successfully."

################# TODO: SQLC #################
# Generate SQLC code
sqlc:
	sqlc generate

# Generate Swagger documentation
swagger:
	@echo "Generating Swagger documentation..."
	swag init --parseDependency -g $(GO_SERVER_PRO) -o $(SWAGGER_DIR)




