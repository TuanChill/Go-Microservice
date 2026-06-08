package clients

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	authpb "api_gateway/internal/gen/auth/v1"
	otppb "api_gateway/internal/gen/otp/v1"
	userpb "api_gateway/internal/gen/user/v1"
)

type Clients struct {
	Auth authpb.AuthServiceClient
	User userpb.UserServiceClient
	Otp  otppb.OtpServiceClient
}

func Dial(authAddr, userAddr, otpAddr string) (*Clients, error) {
	authConn, err := dial(authAddr)
	if err != nil {
		return nil, fmt.Errorf("auth-service: %w", err)
	}
	userConn, err := dial(userAddr)
	if err != nil {
		return nil, fmt.Errorf("user-service: %w", err)
	}
	otpConn, err := dial(otpAddr)
	if err != nil {
		return nil, fmt.Errorf("otp-service: %w", err)
	}
	return &Clients{
		Auth: authpb.NewAuthServiceClient(authConn),
		User: userpb.NewUserServiceClient(userConn),
		Otp:  otppb.NewOtpServiceClient(otpConn),
	}, nil
}

func dial(addr string) (*grpc.ClientConn, error) {
	return grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
}
