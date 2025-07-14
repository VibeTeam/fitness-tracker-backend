package handler

import (
	"context"
	"errors"

	userv1 "github.com/VibeTeam/fitness-tracker-backend/proto/gen/go/user/v1"
	"github.com/VibeTeam/fitness-tracker-backend/user/repository"
	usecase "github.com/VibeTeam/fitness-tracker-backend/user/use_case"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GRPCServer implements user.v1.UserService and delegates work to AuthService.
type GRPCServer struct {
	userv1.UnimplementedUserServiceServer
	authSvc *usecase.AuthService
	repo    repository.UserRepository
}

// NewGRPCServer builds a GRPCServer instance.
func NewGRPCServer(authSvc *usecase.AuthService, repo repository.UserRepository) *GRPCServer {
	return &GRPCServer{authSvc: authSvc, repo: repo}
}

func (s *GRPCServer) Register(ctx context.Context, req *userv1.RegisterRequest) (*userv1.RegisterResponse, error) {
	if req.GetEmail() == "" || req.GetPasswordHash() == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password are required")
	}

	if _, _, err := s.authSvc.Register(ctx, req.GetEmail(), req.GetPasswordHash()); err != nil {
		if errors.Is(err, usecase.ErrEmailAlreadyUsed) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	user, err := s.repo.GetByEmail(ctx, req.GetEmail())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &userv1.RegisterResponse{UserId: int32(user.ID)}, nil
}

func (s *GRPCServer) Login(ctx context.Context, req *userv1.LoginRequest) (*userv1.LoginResponse, error) {
	if req.GetEmail() == "" || req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password are required")
	}

	access, refresh, err := s.authSvc.Login(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidCredentials) {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &userv1.LoginResponse{AccessToken: access, RefreshToken: refresh}, nil
}

func (s *GRPCServer) RefreshToken(ctx context.Context, req *userv1.RefreshTokenRequest) (*userv1.RefreshTokenResponse, error) {
	if req.GetRefreshToken() == "" {
		return nil, status.Error(codes.InvalidArgument, "refresh token is required")
	}

	access, refresh, err := s.authSvc.Refresh(ctx, req.GetRefreshToken())
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	return &userv1.RefreshTokenResponse{AccessToken: access, RefreshToken: refresh}, nil
}

func (s *GRPCServer) ValidateToken(ctx context.Context, req *userv1.ValidateTokenRequest) (*userv1.ValidateTokenResponse, error) {
	if req.GetAccessToken() == "" {
		return nil, status.Error(codes.InvalidArgument, "access token is required")
	}

	uid, err := s.authSvc.Validate(ctx, req.GetAccessToken())
	if err != nil {
		return &userv1.ValidateTokenResponse{Valid: false}, nil
	}

	return &userv1.ValidateTokenResponse{UserId: int32(uid), Valid: true}, nil
}
