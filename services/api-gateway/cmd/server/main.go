package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"


	"api_gateway/internal/clients"
	"api_gateway/internal/handler"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	authAddr := env("AUTH_SERVICE_GRPC_ADDR", "localhost:9001")
	userAddr := env("USER_SERVICE_GRPC_ADDR", "localhost:9002")
	otpAddr := env("OTP_SERVICE_GRPC_ADDR", "localhost:9003")

	c, err := clients.Dial(authAddr, userAddr, otpAddr)
	if err != nil {
		log.Fatalf("dial gRPC services: %v", err)
	}
	logger.Info("connected to gRPC backends",
		"auth", authAddr, "user", userAddr, "otp", otpAddr)

	port := env("PORT", "8080")
	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           handler.NewRouter(c.Auth, c.User, c.Otp),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		logger.Info("api-gateway listening", "port", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down")
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
