FROM golang:1.24-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /notification-otp-service ./cmd/server

FROM alpine:3.20
WORKDIR /app
COPY --from=build /notification-otp-service /notification-otp-service
EXPOSE 8000
CMD ["/notification-otp-service"]
