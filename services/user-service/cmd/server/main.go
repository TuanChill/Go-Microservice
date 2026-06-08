package main

import (
	"log"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"

	pb "user_service/internal/gen/user/v1"
	"user_service/internal/grpcapi"
	"user_service/internal/user"
)

func main() {
	serviceToken := mustEnv("SERVICE_TOKEN")

	repo := user.NewMemoryRepository(time.Now)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcapi.RecoveryInterceptor,
			grpcapi.LoggingInterceptor,
			grpcapi.AuthInterceptor(serviceToken),
			grpcapi.CorrelationInterceptor,
		),
	)
	pb.RegisterUserServiceServer(grpcServer, grpcapi.NewServer(repo))

	addr := ":" + env("GRPC_PORT", "9002")
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("listen %s: %v", addr, err)
	}
	log.Printf("user-service gRPC listening on %s", addr)
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
