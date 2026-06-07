package main

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"api_gateway/internal/gateway"
)

func main() {
	legacyURL := mustURL(env("LEGACY_URL", "http://localhost:8000"))
	routes := []gateway.Route{
		{Prefix: "/v1/auth", Target: mustURL(env("AUTH_SERVICE_URL", "http://localhost:8001"))},
		{Prefix: "/v1/user", Target: mustURL(env("USER_SERVICE_URL", "http://localhost:8002"))},
		{Prefix: "/v1/otp", Target: mustURL(env("NOTIFICATION_OTP_SERVICE_URL", "http://localhost:8003"))},
	}

	port := env("PORT", "8080")
	server := &http.Server{Addr: ":" + port, Handler: gateway.New(routes, legacyURL), ReadHeaderTimeout: 5 * time.Second}
	log.Printf("api-gateway listening on :%s", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func env(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func mustURL(raw string) *url.URL {
	parsed, err := gateway.ParseURL(raw)
	if err != nil {
		log.Fatal(err)
	}
	return parsed
}
