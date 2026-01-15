package usecase

import (
	"context"
	"fmt"

	"idmService/application/domain"
	"idmService/application/repository"
	"idmService/application/service"
)

type RegisterResult struct {
	UserID string
}

type RegisterUseCase interface {
	Execute(ctx context.Context, email, password, role string) (*RegisterResult, error)
}

type registerUseCase struct {
	userRepo          repository.UserRepository
	userServiceClient service.UserServiceClient
	passwordHasher    service.PasswordHasher
}

func NewRegisterUseCase(
	userRepo repository.UserRepository,
	userServiceClient service.UserServiceClient,
	passwordHasher service.PasswordHasher,
) RegisterUseCase {
	return &registerUseCase{
		userRepo:          userRepo,
		userServiceClient: userServiceClient,
		passwordHasher:    passwordHasher,
	}
}

func (uc *registerUseCase) Execute(ctx context.Context, email, password, role string) (*RegisterResult, error) {
	existing, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		return nil, &domain.ValidationError{
			Field:  "email",
			Reason: "User with this email already exists",
		}
	}

	// Hash password before storing
	hashedPassword, err := uc.passwordHasher.HashPassword(password)
	if err != nil {
		return nil, &domain.InternalError{
			Operation: "hash password",
			Err:       fmt.Errorf("failed to hash password: %w", err),
		}
	}

	userRole := domain.RoleOwnerEvent
	if role == "client" {
		userRole = domain.RoleClient
	} else if role == "owner" {
		userRole = domain.RoleOwnerEvent
	}

	user := &domain.User{
		Email:  email,
		Parola: hashedPassword,
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
		UserID: fmt.Sprintf("%d", user.ID),
	}, nil
}
