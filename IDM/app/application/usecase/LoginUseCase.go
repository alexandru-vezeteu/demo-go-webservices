package usecase

import (
	"context"
	"fmt"

	"idmService/application/domain"
	"idmService/application/repository"
	"idmService/application/service"
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
	Execute(ctx context.Context, email, password string) (*LoginResult, error)
}

type loginUseCase struct {
	userRepo     repository.UserRepository
	tokenService service.TokenService
}

func NewLoginUseCase(
	userRepo repository.UserRepository,
	tokenService service.TokenService,
) LoginUseCase {
	return &loginUseCase{
		userRepo:     userRepo,
		tokenService: tokenService,
	}
}

func (uc *loginUseCase) Execute(ctx context.Context, email, password string) (*LoginResult, error) {

	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if user == nil || user.Parola != password {
		return nil, &domain.AuthenticationError{Reason: "invalid email or password"}
	}

	token, err := uc.tokenService.GenerateJWT(user)
	if err != nil {
		return nil, err
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
