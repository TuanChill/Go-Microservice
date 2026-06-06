FROM golang:1.24-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG APP_CMD=api
RUN go build -o /out/app ./cmd/${APP_CMD}

FROM alpine:3.20
RUN adduser -D -H appuser
WORKDIR /app
COPY --from=build /out/app ./app
USER appuser
EXPOSE 8000
CMD ["./app"]
