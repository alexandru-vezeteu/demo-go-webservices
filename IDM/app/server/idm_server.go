package server

import (
	"context"
	"errors"

	"idmService/application/domain"
	"idmService/application/usecase"
	pb "idmService/proto"
)

type IdentityServer struct {
	pb.UnimplementedIdentityServiceServer
	registerUseCase    usecase.RegisterUseCase
	loginUseCase       usecase.LoginUseCase
	verifyTokenUseCase usecase.VerifyTokenUseCase
	revokeTokenUseCase usecase.RevokeTokenUseCase
}

func NewIdentityServer(
	registerUseCase usecase.RegisterUseCase,
	loginUseCase usecase.LoginUseCase,
	verifyTokenUseCase usecase.VerifyTokenUseCase,
	revokeTokenUseCase usecase.RevokeTokenUseCase,
) *IdentityServer {
	return &IdentityServer{
		registerUseCase:    registerUseCase,
		loginUseCase:       loginUseCase,
		verifyTokenUseCase: verifyTokenUseCase,
		revokeTokenUseCase: revokeTokenUseCase,
	}
}

func (s *IdentityServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	result, err := s.registerUseCase.Execute(ctx, req.Email, req.Password, req.Role)
	if err != nil {
		var internalErr *domain.InternalError
		if errors.As(err, &internalErr) {
			return &pb.RegisterResponse{
				Success: false,
				Message: "Internal error during registration",
			}, nil
		}

		return &pb.RegisterResponse{
			Success: false,
			Message: "Registration failed",
		}, nil
	}

	return &pb.RegisterResponse{
		Success: result.Success,
		Message: result.Message,
		UserId:  result.UserID,
	}, nil
}

func (s *IdentityServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	result, err := s.loginUseCase.Execute(ctx, req.Email, req.Password)
	if err != nil {
		var authErr *domain.AuthenticationError
		if errors.As(err, &authErr) {
			return &pb.LoginResponse{
				Success: false,
				Token:   "",
				Message: authErr.Error(),
			}, nil
		}

		var configErr *domain.ConfigurationError
		if errors.As(err, &configErr) {
			return &pb.LoginResponse{
				Success: false,
				Token:   "",
				Message: "Service configuration error",
			}, nil
		}

		return &pb.LoginResponse{
			Success: false,
			Token:   "",
			Message: "Internal error",
		}, nil
	}

	return &pb.LoginResponse{
		Success: result.Success,
		Token:   result.Token,
		Message: result.Message,
		UserId:  result.UserID,
		Role:    result.Role,
		Email:   result.Email,
	}, nil
}

func (s *IdentityServer) VerifyToken(ctx context.Context, req *pb.VerifyTokenRequest) (*pb.VerifyTokenResponse, error) {
	result, err := s.verifyTokenUseCase.Execute(ctx, req.Token)
	if err != nil {
		var tokenErr *domain.TokenError
		if errors.As(err, &tokenErr) {
			return &pb.VerifyTokenResponse{
				Valid:       false,
				Message:     tokenErr.Error(),
				Expired:     tokenErr.Expired,
				Blacklisted: tokenErr.Blacklisted,
			}, nil
		}

		return &pb.VerifyTokenResponse{
			Valid:   false,
			Message: "Internal error",
		}, nil
	}

	return &pb.VerifyTokenResponse{
		Valid:     result.Valid,
		Email:     result.Email,
		Message:   result.Message,
		UserId:    result.UserID,
		Role:      result.Role,
		Issuer:    result.Issuer,
		ExpiresAt: result.ExpiresAt,
	}, nil
}

func (s *IdentityServer) RevokeToken(ctx context.Context, req *pb.RevokeTokenRequest) (*pb.RevokeTokenResponse, error) {
	result, err := s.revokeTokenUseCase.Execute(ctx, req.Token)
	if err != nil {
		return &pb.RevokeTokenResponse{
			Success: false,
			Message: "Internal error",
		}, nil
	}

	return &pb.RevokeTokenResponse{
		Success: result.Success,
		Message: result.Message,
	}, nil
}
