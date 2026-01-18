package server

import (
	"errors"
	"idmService/application/domain"
	pb "idmService/proto"
)

func mapRegisterError(err error) *pb.RegisterResponse {
	var validationErr *domain.ValidationError
	if errors.As(err, &validationErr) {
		return &pb.RegisterResponse{
			Success: false,
			Message: validationErr.Error(),
		}
	}

	var internalErr *domain.InternalError
	if errors.As(err, &internalErr) {
		return &pb.RegisterResponse{
			Success: false,
			Message: "Internal error during registration",
		}
	}

	return &pb.RegisterResponse{
		Success: false,
		Message: "Registration failed",
	}
}

func mapLoginError(err error) *pb.LoginResponse {
	var authErr *domain.AuthenticationError
	if errors.As(err, &authErr) {
		return &pb.LoginResponse{
			Success: false,
			Token:   "",
			Message: authErr.Error(),
		}
	}

	var configErr *domain.ConfigurationError
	if errors.As(err, &configErr) {
		return &pb.LoginResponse{
			Success: false,
			Token:   "",
			Message: "Service configuration error",
		}
	}

	return &pb.LoginResponse{
		Success: false,
		Token:   "",
		Message: "Internal error",
	}
}

func mapVerifyTokenError(err error) *pb.VerifyTokenResponse {
	var tokenErr *domain.TokenError
	if errors.As(err, &tokenErr) {
		return &pb.VerifyTokenResponse{
			Valid:       false,
			Message:     tokenErr.Error(),
			Expired:     tokenErr.Expired,
			Blacklisted: tokenErr.Blacklisted,
		}
	}

	return &pb.VerifyTokenResponse{
		Valid:   false,
		Message: "Internal error",
	}
}

func mapRevokeTokenError(err error) *pb.RevokeTokenResponse {
	return &pb.RevokeTokenResponse{
		Success: false,
		Message: "Internal error",
	}
}
