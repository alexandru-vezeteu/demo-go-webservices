package usecase

import (
	"context"
	"fmt"

	"idmService/application/domain"
	"idmService/application/repository"
	"idmService/application/service"
)

type LoginResult struct {
	Token  string
	UserID string
	Role   string
	Email  string
}

type LoginUseCase interface {
	Execute(ctx context.Context, email, password string) (*LoginResult, error)
}

type loginUseCase struct {
	userRepo       repository.UserRepository
	tokenService   service.TokenService
	passwordHasher service.PasswordHasher
}

func NewLoginUseCase(
	userRepo repository.UserRepository,
	tokenService service.TokenService,
	passwordHasher service.PasswordHasher,
) LoginUseCase {
	return &loginUseCase{
		userRepo:       userRepo,
		tokenService:   tokenService,
		passwordHasher: passwordHasher,
	}
}

func (uc *loginUseCase) Execute(ctx context.Context, email, password string) (*LoginResult, error) {

	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, &domain.AuthenticationError{Reason: "invalid email or password"}
	}

	if err := uc.passwordHasher.CheckPassword(user.Parola, password); err != nil {
		return nil, &domain.AuthenticationError{Reason: "invalid email or password"}
	}

	token, err := uc.tokenService.GenerateJWT(user)
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		Token:  token,
		UserID: fmt.Sprintf("%d", user.ID),
		Role:   string(user.Rol),
		Email:  user.Email,
	}, nil
}
