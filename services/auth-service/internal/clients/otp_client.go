package clients

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "auth_service/internal/gen/otp/v1"
)

type OtpClient struct {
	c pb.OtpServiceClient
}

func NewOtpClient(addr string) (*OtpClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("dial otp-service %s: %w", addr, err)
	}
	return &OtpClient{c: pb.NewOtpServiceClient(conn)}, nil
}

func (o *OtpClient) RequestOtp(ctx context.Context, userID int32, email, purpose string) error {
	_, err := o.c.RequestOtp(ctx, &pb.RequestOtpRequest{
		UserId:  userID,
		Email:   email,
		Purpose: purpose,
	})
	return err
}
