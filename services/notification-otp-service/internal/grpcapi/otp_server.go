package grpcapi

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "notification_otp_service/internal/gen/otp/v1"
	"notification_otp_service/internal/otp"
)

type OtpServer struct {
	pb.UnimplementedOtpServiceServer
	svc *otp.Service
}

func NewOtpServer(svc *otp.Service) *OtpServer {
	return &OtpServer{svc: svc}
}

func (s *OtpServer) RequestOtp(_ context.Context, req *pb.RequestOtpRequest) (*pb.RequestOtpResponse, error) {
	resp, err := s.svc.Request(otp.Request{
		UserID:  int(req.UserId),
		Email:   req.Email,
		Purpose: otp.Purpose(req.Purpose),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "request otp: %v", err)
	}
	return &pb.RequestOtpResponse{
		Id:        int32(resp.ID),
		UserId:    int32(resp.UserID),
		ExpiresAt: timestamppb.New(resp.ExpiresAt),
	}, nil
}

func (s *OtpServer) VerifyOtp(_ context.Context, req *pb.VerifyOtpRequest) (*pb.VerifyOtpResponse, error) {
	resp, err := s.svc.Verify(otp.VerifyRequest{Otp: req.Otp})
	if err != nil {
		if errors.Is(err, otp.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "otp not found or expired")
		}
		return nil, status.Errorf(codes.Internal, "verify otp: %v", err)
	}
	return &pb.VerifyOtpResponse{
		UserId: int32(resp.UserID),
		Email:  resp.Email,
	}, nil
}
