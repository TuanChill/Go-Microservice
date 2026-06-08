package main

import (
	"context"
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"

	pb "notification_otp_service/internal/gen/otp/v1"
	"notification_otp_service/internal/events"
	"notification_otp_service/internal/grpcapi"
	"notification_otp_service/internal/messaging"
	"notification_otp_service/internal/notification"
	"notification_otp_service/internal/otp"
)

func main() {
	serviceToken := mustEnv("SERVICE_TOKEN")
	logger := slog.Default()

	otpRepo := otp.NewMemoryRepository()
	otpSvc := otp.NewService(otpRepo, otp.GenerateSixDigitCode, time.Now, 10*time.Minute)
	notifSvc := notification.NewService()

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcapi.RecoveryInterceptor,
			grpcapi.LoggingInterceptor,
			grpcapi.AuthInterceptor(serviceToken),
			grpcapi.CorrelationInterceptor,
		),
	)
	pb.RegisterOtpServiceServer(grpcServer, grpcapi.NewOtpServer(otpSvc))
	pb.RegisterNotificationServiceServer(grpcServer, grpcapi.NewNotificationServer(notifSvc))

	addr := ":" + env("GRPC_PORT", "9003")
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("listen %s: %v", addr, err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if mqURL := os.Getenv("RABBITMQ_URL"); mqURL != "" {
		conn, err := amqp.Dial(mqURL)
		if err != nil {
			log.Fatalf("rabbitmq dial: %v", err)
		}
		handler := events.NewHandler(events.NewMemoryStore())
		consumer := messaging.NewConsumer(conn, env("OTP_QUEUE", "otp.requests"), handler, logger)
		go func() {
			if err := consumer.Run(ctx); err != nil && ctx.Err() == nil {
				logger.Error("consumer exited", "err", err)
			}
		}()
	}

	go func() {
		<-ctx.Done()
		grpcServer.GracefulStop()
	}()

	log.Printf("notification-otp-service gRPC listening on %s", addr)
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
