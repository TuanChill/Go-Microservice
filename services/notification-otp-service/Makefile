test:
	go test ./...

race:
	go test -race ./...

build:
	go build -o /tmp/notification-otp-service ./cmd/server

docker-build:
	docker build -t notification-otp-service:local .
