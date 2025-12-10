package usecase

import (
	"fmt"
	"idmService/domain"
	"idmService/service"
)

type LoginResult struct {
	Success bool
	Token   string
	Message string
	UserID  string
	Role    string
	Email   string
}

type LoginUseCase interface {
	Execute(email, password string) (*LoginResult, error)
}

type loginUseCase struct {
	userRepo     domain.UserRepository
	tokenService service.TokenService
}

func NewLoginUseCase(userRepo domain.UserRepository, tokenService service.TokenService) LoginUseCase {
	return &loginUseCase{
		userRepo:     userRepo,
		tokenService: tokenService,
	}
}

func (uc *loginUseCase) Execute(email, password string) (*LoginResult, error) {
	// Find user by email
	user, err := uc.userRepo.FindByEmail(email)
	if err != nil {
		return &LoginResult{
			Success: false,
			Token:   "",
			Message: "Database error",
		}, nil
	}

	// Validate credentials
	if user == nil || user.Parola != password {
		return &LoginResult{
			Success: false,
			Token:   "",
			Message: "Invalid email or password",
		}, nil
	}

	// Generate JWT token
	token, err := uc.tokenService.GenerateJWT(user)
	if err != nil {
		return &LoginResult{
			Success: false,
			Token:   "",
			Message: "Failed to generate token",
		}, nil
	}

	return &LoginResult{
		Success: true,
		Token:   token,
		Message: "Login successful",
		UserID:  fmt.Sprintf("%d", user.ID),
		Role:    string(user.Rol),
		Email:   user.Email,
	}, nil
}
