package server

import (
	"context"

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
		return mapRegisterError(err), nil
	}

	return &pb.RegisterResponse{
		Success: true,
		Message: "Registration successful",
		UserId:  result.UserID,
	}, nil
}

func (s *IdentityServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	result, err := s.loginUseCase.Execute(ctx, req.Email, req.Password)
	if err != nil {
		return mapLoginError(err), nil
	}

	return &pb.LoginResponse{
		Success: true,
		Token:   result.Token,
		Message: "Login successful",
		UserId:  result.UserID,
		Role:    result.Role,
		Email:   result.Email,
	}, nil
}

func (s *IdentityServer) VerifyToken(ctx context.Context, req *pb.VerifyTokenRequest) (*pb.VerifyTokenResponse, error) {
	result, err := s.verifyTokenUseCase.Execute(ctx, req.Token)
	if err != nil {
		return mapVerifyTokenError(err), nil
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
		return mapRevokeTokenError(err), nil
	}

	return &pb.RevokeTokenResponse{
		Success: result.Success,
		Message: result.Message,
	}, nil
}
