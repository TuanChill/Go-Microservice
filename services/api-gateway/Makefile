test:
	go test ./...

race:
	go test -race ./...

build:
	go build -o /tmp/api-gateway ./cmd/server

docker-build:
	docker build -t api-gateway:local .
