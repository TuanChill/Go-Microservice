package main

import (
	"log"
	"log/slog"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"

	pb "auth_service/internal/gen/auth/v1"
	"auth_service/internal/auth"
	"auth_service/internal/clients"
	"auth_service/internal/grpcapi"
)

func main() {
	serviceToken := mustEnv("SERVICE_TOKEN")

	userClient, err := clients.NewUserClient(env("USER_SERVICE_GRPC_ADDR", "localhost:9002"))
	if err != nil {
		log.Fatalf("user client: %v", err)
	}
	otpClient, err := clients.NewOtpClient(env("OTP_SERVICE_GRPC_ADDR", "localhost:9003"))
	if err != nil {
		log.Fatalf("otp client: %v", err)
	}

	svc := auth.NewService(time.Now)
	srv := grpcapi.NewServer(svc, userClient, otpClient, slog.Default())

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcapi.RecoveryInterceptor,
			grpcapi.LoggingInterceptor,
			grpcapi.AuthInterceptor(serviceToken),
			grpcapi.CorrelationInterceptor,
		),
	)
	pb.RegisterAuthServiceServer(grpcServer, srv)

	addr := ":" + env("GRPC_PORT", "9001")
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("listen %s: %v", addr, err)
	}
	log.Printf("auth-service gRPC listening on %s", addr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("serve: %v", err)
	}
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("%s is required", key)
	}
	return v
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
