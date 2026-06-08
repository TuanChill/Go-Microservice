package grpcapi

import (
	"context"
	"log/slog"
	"os"
	"runtime/debug"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

type correlationKey struct{}

func RecoveryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	defer func() {
		if r := recover(); r != nil {
			logger.Error("panic in gRPC handler", "method", info.FullMethod, "panic", r, "stack", string(debug.Stack()))
			err = status.Error(codes.Internal, "internal server error")
		}
	}()
	return handler(ctx, req)
}

func LoggingInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	logger.Info("rpc",
		"method", info.FullMethod,
		"duration_ms", time.Since(start).Milliseconds(),
		"code", status.Code(err).String(),
		"correlation_id", correlationFromCtx(ctx),
	)
	return resp, err
}

func AuthInterceptor(serviceToken string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}
		tokens := md.Get("x-service-token")
		if len(tokens) == 0 || tokens[0] != serviceToken {
			return nil, status.Error(codes.Unauthenticated, "invalid service token")
		}
		return handler(ctx, req)
	}
}

func CorrelationInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if ids := md.Get("x-correlation-id"); len(ids) > 0 {
			ctx = context.WithValue(ctx, correlationKey{}, ids[0])
		}
	}
	return handler(ctx, req)
}

func correlationFromCtx(ctx context.Context) string {
	if v, ok := ctx.Value(correlationKey{}).(string); ok {
		return v
	}
	return ""
}
