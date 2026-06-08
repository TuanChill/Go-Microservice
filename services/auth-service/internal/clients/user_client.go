package clients

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "auth_service/internal/gen/user/v1"
)

type UserClient struct {
	c pb.UserServiceClient
}

func NewUserClient(addr string) (*UserClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("dial user-service %s: %w", addr, err)
	}
	return &UserClient{c: pb.NewUserServiceClient(conn)}, nil
}

func (u *UserClient) Create(ctx context.Context, email string) (int32, error) {
	resp, err := u.c.CreateUser(ctx, &pb.CreateUserRequest{Email: email})
	if err != nil {
		return 0, err
	}
	return resp.Id, nil
}
