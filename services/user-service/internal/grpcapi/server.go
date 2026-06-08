package grpcapi

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "user_service/internal/gen/user/v1"
	"user_service/internal/user"
)

type Server struct {
	pb.UnimplementedUserServiceServer
	repo user.Repository
}

func NewServer(repo user.Repository) *Server {
	return &Server{repo: repo}
}

func (s *Server) CreateUser(_ context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	identity, err := s.repo.Create(req.Email)
	if err != nil {
		if errors.Is(err, user.ErrConflict) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Errorf(codes.Internal, "create user: %v", err)
	}
	return &pb.CreateUserResponse{Id: int32(identity.ID), Email: identity.Email}, nil
}

func (s *Server) GetProfile(_ context.Context, req *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {
	profile, err := s.repo.GetProfile(int(req.Id))
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "get profile: %v", err)
	}
	return profileToProto(profile), nil
}

func (s *Server) UpdateProfile(_ context.Context, req *pb.UpdateProfileRequest) (*pb.UpdateProfileResponse, error) {
	updated, err := s.repo.UpdateProfile(int(req.Id), user.UpdateProfileRequest{
		Username: req.Username,
		Phone:    req.Phone,
		FullName: req.Fullname,
		Avatar:   req.Avatar,
		Gender:   int(req.Gender),
	})
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "update profile: %v", err)
	}
	return &pb.UpdateProfileResponse{Profile: profileToProto(updated)}, nil
}

func (s *Server) DeactivateUser(_ context.Context, req *pb.DeactivateUserRequest) (*pb.DeactivateUserResponse, error) {
	res, err := s.repo.Deactivate(int(req.Id))
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "deactivate: %v", err)
	}
	return &pb.DeactivateUserResponse{Id: int32(res.ID)}, nil
}

func profileToProto(p user.Profile) *pb.GetProfileResponse {
	return &pb.GetProfileResponse{
		Id:               int32(p.ID),
		Username:         p.Username,
		Email:            p.Email,
		Phone:            p.Phone,
		HiddenPhoneNumber: p.HiddenPhoneNumber,
		Fullname:         p.FullName,
		HiddenEmail:      p.HiddenEmail,
		Avatar:           p.Avatar,
		Gender:           int32(p.Gender),
		TwoFactorEnabled: p.TwoFactorEnabled,
		IsActive:         p.IsActive,
		CreatedAt:        timestamppb.New(p.CreatedAt),
	}
}
