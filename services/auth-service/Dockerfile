FROM golang:1.24-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /auth-service ./cmd/server

FROM alpine:3.20
WORKDIR /app
COPY --from=build /auth-service /auth-service
EXPOSE 8000
CMD ["/auth-service"]
