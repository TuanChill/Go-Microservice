test:
	go test ./...

race:
	go test -race ./...

build:
	go build -o /tmp/auth-service ./cmd/server

docker-build:
	docker build -t auth-service:local .
