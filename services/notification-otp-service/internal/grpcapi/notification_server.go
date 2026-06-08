package grpcapi

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "notification_otp_service/internal/gen/otp/v1"
	"notification_otp_service/internal/notification"
)

type NotificationServer struct {
	pb.UnimplementedNotificationServiceServer
	svc *notification.Service
}

func NewNotificationServer(svc *notification.Service) *NotificationServer {
	return &NotificationServer{svc: svc}
}

func (s *NotificationServer) SendEmailVerification(_ context.Context, req *pb.SendEmailVerificationRequest) (*pb.SendEmailVerificationResponse, error) {
	result := s.svc.AcceptEmailVerification(req.IdempotencyKey, notification.EmailVerificationRequest{
		UserID:         int(req.UserId),
		Email:          req.Email,
		Token:          req.Token,
		ExpiresAtToken: req.ExpiresAtToken,
	})
	if result.IdempotencyKey == "" {
		return nil, status.Error(codes.Internal, "send email verification failed")
	}
	return &pb.SendEmailVerificationResponse{IdempotencyKey: result.IdempotencyKey}, nil
}

func (s *NotificationServer) SendPasswordReset(_ context.Context, req *pb.SendPasswordResetRequest) (*pb.SendPasswordResetResponse, error) {
	result := s.svc.AcceptPasswordReset(req.IdempotencyKey, notification.PasswordResetRequest{
		UserID: int(req.UserId),
		Email:  req.Email,
		Token:  req.Token,
	})
	if result.IdempotencyKey == "" {
		return nil, status.Error(codes.Internal, "send password reset failed")
	}
	return &pb.SendPasswordResetResponse{IdempotencyKey: result.IdempotencyKey}, nil
}
