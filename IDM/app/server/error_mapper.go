package server

import (
	"errors"
	"idmService/application/domain"
	pb "idmService/proto"
)

// mapRegisterError maps domain errors to RegisterResponse
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

// mapLoginError maps domain errors to LoginResponse
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

// mapVerifyTokenError maps domain errors to VerifyTokenResponse
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

// mapRevokeTokenError maps domain errors to RevokeTokenResponse
func mapRevokeTokenError(err error) *pb.RevokeTokenResponse {
	return &pb.RevokeTokenResponse{
		Success: false,
		Message: "Internal error",
	}
}
