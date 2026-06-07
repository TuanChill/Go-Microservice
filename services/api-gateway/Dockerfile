FROM golang:1.24-alpine AS build
WORKDIR /app
COPY go.mod ./
COPY . .
RUN go build -o /api-gateway ./cmd/server

FROM alpine:3.20
WORKDIR /app
COPY --from=build /api-gateway /api-gateway
EXPOSE 8080
CMD ["/api-gateway"]
