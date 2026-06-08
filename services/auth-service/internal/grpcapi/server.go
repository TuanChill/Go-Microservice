package grpcapi

import (
	"context"
	"errors"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "auth_service/internal/gen/auth/v1"
	"auth_service/internal/auth"
	"auth_service/internal/clients"
)

type Server struct {
	pb.UnimplementedAuthServiceServer
	svc        *auth.Service
	userClient *clients.UserClient
	otpClient  *clients.OtpClient
	logger     *slog.Logger
}

func NewServer(svc *auth.Service, userClient *clients.UserClient, otpClient *clients.OtpClient, logger *slog.Logger) *Server {
	return &Server{svc: svc, userClient: userClient, otpClient: otpClient, logger: logger}
}

func (s *Server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	res, err := s.svc.Register(auth.RegisterRequest{Email: req.Email})
	if errors.Is(err, auth.ErrConflict) {
		return nil, status.Error(codes.AlreadyExists, "credential already exists")
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "register: %v", err)
	}

	userID, err := s.userClient.Create(ctx, req.Email)
	if err != nil {
		if status.Code(err) != codes.AlreadyExists {
			return nil, status.Errorf(codes.Internal, "create user profile: %v", err)
		}
	}

	if otpErr := s.otpClient.RequestOtp(ctx, userID, req.Email, "email_verification"); otpErr != nil {
		s.logger.Warn("otp dispatch failed after register", "email", req.Email, "err", otpErr)
	}

	return &pb.RegisterResponse{
		Id:             int32(res.ID),
		Email:          res.Email,
		Token:          res.Token,
		ExpiresAtToken: timestamppb.New(res.ExpiresAtToken),
	}, nil
}

func (s *Server) Login(_ context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	res, err := s.svc.Login(auth.LoginRequest{
		Identifier: req.Identifier,
		Password:   req.Password,
		DeviceID:   req.DeviceId,
	})
	if errors.Is(err, auth.ErrInvalidCredentials) {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "login: %v", err)
	}
	return &pb.LoginResponse{
		Id:          int32(res.ID),
		DeviceId:    res.DeviceID,
		Email:       res.Email,
		AccessToken: res.AccessToken,
	}, nil
}

func (s *Server) RefreshToken(_ context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	res, err := s.svc.Refresh(auth.RefreshRequest{
		RefreshToken: req.RefreshToken,
		DeviceID:     req.DeviceId,
	})
	if errors.Is(err, auth.ErrInvalidToken) {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "refresh: %v", err)
	}
	return &pb.RefreshTokenResponse{AccessToken: res.AccessToken}, nil
}

func (s *Server) VerifyAccount(_ context.Context, req *pb.VerifyAccountRequest) (*pb.VerifyAccountResponse, error) {
	_, err := s.svc.VerifyAccount(auth.VerifyAccountRequest{
		UserID:        int(req.UserId),
		VerifiedToken: req.VerifiedToken,
		Email:         req.Email,
	})
	if errors.Is(err, auth.ErrInvalidToken) {
		return nil, status.Error(codes.Unauthenticated, "invalid or expired token")
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "verify account: %v", err)
	}
	return &pb.VerifyAccountResponse{Id: req.UserId}, nil
}

func (s *Server) ResetPassword(_ context.Context, req *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error) {
	res, err := s.svc.ResetPassword(auth.ResetPasswordRequest{
		Token:    req.Token,
		UserID:   int(req.UserId),
		Password: req.Password,
	})
	if errors.Is(err, auth.ErrInvalidToken) {
		return nil, status.Error(codes.Unauthenticated, "invalid or expired token")
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "reset password: %v", err)
	}
	return &pb.ResetPasswordResponse{Id: int32(res.ID)}, nil
}
