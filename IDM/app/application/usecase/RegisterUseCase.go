package usecase

import (
	"context"
	"fmt"

	"idmService/application/domain"
	"idmService/application/repository"
	"idmService/application/service"
)

type RegisterResult struct {
	Success bool
	Message string
	UserID  string
}

type RegisterUseCase interface {
	Execute(ctx context.Context, email, password, role string) (*RegisterResult, error)
}

type registerUseCase struct {
	userRepo          repository.UserRepository
	userServiceClient service.UserServiceClient
}

func NewRegisterUseCase(
	userRepo repository.UserRepository,
	userServiceClient service.UserServiceClient,
) RegisterUseCase {
	return &registerUseCase{
		userRepo:          userRepo,
		userServiceClient: userServiceClient,
	}
}

func (uc *registerUseCase) Execute(ctx context.Context, email, password, role string) (*RegisterResult, error) {
	existing, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		return &RegisterResult{
			Success: false,
			Message: "User with this email already exists",
		}, nil
	}

	userRole := domain.RoleOwnerEvent
	if role == "client" {
		userRole = domain.RoleClient
	} else if role == "owner" {
		userRole = domain.RoleOwnerEvent
	}

	user := &domain.User{
		Email:  email,
		Parola: password,
		Rol:    userRole,
	}

	err = uc.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	userProfileReq := &service.CreateUserRequest{
		ID:        int(user.ID),
		Email:     email,
		FirstName: "",
		LastName:  "",
	}

	_, err = uc.userServiceClient.CreateUser(ctx, userProfileReq)
	if err != nil {
		return nil, &domain.InternalError{
			Operation: "register user",
			Err:       fmt.Errorf("auth account created but user profile creation failed: %w", err),
		}
	}

	return &RegisterResult{
		Success: true,
		Message: "Registration successful",
		UserID:  fmt.Sprintf("%d", user.ID),
	}, nil
}
