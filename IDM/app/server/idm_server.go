package server

import (
	"context"

	pb "idmService/proto"
	"idmService/usecase"
)

type IdentityServer struct {
	pb.UnimplementedIdentityServiceServer
	loginUseCase       usecase.LoginUseCase
	verifyTokenUseCase usecase.VerifyTokenUseCase
	revokeTokenUseCase usecase.RevokeTokenUseCase
}

func NewIdentityServer(
	loginUseCase usecase.LoginUseCase,
	verifyTokenUseCase usecase.VerifyTokenUseCase,
	revokeTokenUseCase usecase.RevokeTokenUseCase,
) *IdentityServer {
	return &IdentityServer{
		loginUseCase:       loginUseCase,
		verifyTokenUseCase: verifyTokenUseCase,
		revokeTokenUseCase: revokeTokenUseCase,
	}
}

func (s *IdentityServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	result, err := s.loginUseCase.Execute(req.Email, req.Password)
	if err != nil {
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
	result, err := s.verifyTokenUseCase.Execute(req.Token)
	if err != nil {
		return &pb.VerifyTokenResponse{
			Valid:   false,
			Message: "Internal error",
		}, nil
	}

	return &pb.VerifyTokenResponse{
		Valid:       result.Valid,
		Email:       result.Email,
		Message:     result.Message,
		UserId:      result.UserID,
		Role:        result.Role,
		Issuer:      result.Issuer,
		ExpiresAt:   result.ExpiresAt,
		Expired:     result.Expired,
		Blacklisted: result.Blacklisted,
	}, nil
}

func (s *IdentityServer) RevokeToken(ctx context.Context, req *pb.RevokeTokenRequest) (*pb.RevokeTokenResponse, error) {
	result, err := s.revokeTokenUseCase.Execute(req.Token)
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
