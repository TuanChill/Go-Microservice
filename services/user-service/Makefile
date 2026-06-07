test:
	go test ./...

race:
	go test -race ./...

build:
	go build -o /tmp/user-service ./cmd/server

docker-build:
	docker build -t user-service:local .
