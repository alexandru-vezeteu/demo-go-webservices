package usecase

import (
	"context"
	"fmt"
	"userService/application/domain"
	"userService/application/service"
)

type UserUsecase interface {
	CreateUser(ctx context.Context, token string, user *domain.User) (*domain.User, error)
	GetUserByID(ctx context.Context, token string, id int) (*domain.User, error)
	UpdateUser(ctx context.Context, token string, id int, updates map[string]interface{}) (*domain.User, error)
	DeleteUser(ctx context.Context, token string, id int) (*domain.User, error)
	CreateTicketForUser(ctx context.Context, userID int, token string, packetID *int, eventID *int) (string, error)

	GetCustomersByEventID(ctx context.Context, token string, eventID int) ([]*domain.User, error)
	GetCustomersByPacketID(ctx context.Context, token string, packetID int) ([]*domain.User, error)
}

type userUsecase struct {
	userService         service.UserService
	eventManagerService service.EventManagerService
	authNService        service.AuthenticationService
	authZService        service.AuthorizationService
}

func NewUserUsecase(
	userService service.UserService,
	eventManagerService service.EventManagerService,
	authNService service.AuthenticationService,
	authZService service.AuthorizationService,
) UserUsecase {
	return &userUsecase{
		userService:         userService,
		eventManagerService: eventManagerService,
		authNService:        authNService,
		authZService:        authZService,
	}
}

func (uc *userUsecase) authenticate(ctx context.Context, token string) (*service.UserIdentity, error) {
	identity, err := uc.authNService.WhoIsUser(ctx, token)
	if err != nil {
		return nil, &domain.ValidationError{Field: "token", Reason: "invalid or expired token"}
	}
	return identity, nil
}

func (uc *userUsecase) CreateUser(ctx context.Context, token string, user *domain.User) (*domain.User, error) {
	return uc.userService.CreateUser(ctx, user)
}

func (uc *userUsecase) GetUserByID(ctx context.Context, token string, id int) (*domain.User, error) {
	identity, err := uc.authenticate(ctx, token)
	if err != nil {
		return nil, err
	}

	user, err := uc.userService.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	allowed, err := uc.authZService.CanUserViewUser(ctx, identity.UserID, user)
	if err != nil {
		return nil, &domain.ForbiddenError{Reason: fmt.Sprintf("authorization check failed: %v", err)}
	}
	if !allowed {
		return nil, &domain.ForbiddenError{Reason: "user not authorized to view user"}
	}

	return user, nil
}

func (uc *userUsecase) UpdateUser(ctx context.Context, token string, id int, updates map[string]interface{}) (*domain.User, error) {
	identity, err := uc.authenticate(ctx, token)
	if err != nil {
		return nil, err
	}

	user, err := uc.userService.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	allowed, err := uc.authZService.CanUserEditUser(ctx, identity.UserID, user)
	if err != nil {
		return nil, &domain.ForbiddenError{Reason: fmt.Sprintf("authorization check failed: %v", err)}
	}
	if !allowed {
		return nil, &domain.ForbiddenError{Reason: "user not authorized to edit user"}
	}

	return uc.userService.UpdateUser(ctx, id, updates)
}

func (uc *userUsecase) DeleteUser(ctx context.Context, token string, id int) (*domain.User, error) {
	identity, err := uc.authenticate(ctx, token)
	if err != nil {
		return nil, err
	}

	user, err := uc.userService.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	allowed, err := uc.authZService.CanUserDeleteUser(ctx, identity.UserID, user)
	if err != nil {
		return nil, &domain.ForbiddenError{Reason: fmt.Sprintf("authorization check failed: %v", err)}
	}
	if !allowed {
		return nil, &domain.ForbiddenError{Reason: "user not authorized to delete user"}
	}

	return uc.userService.DeleteUser(ctx, id)
}

func (uc *userUsecase) CreateTicketForUser(ctx context.Context, userID int, token string, packetID *int, eventID *int) (string, error) {
	identity, err := uc.authNService.WhoIsUser(ctx, token)
	if err != nil {
		return "", &domain.ValidationError{Field: "token", Reason: "invalid or expired token"}
	}

	user, err := uc.userService.GetUserByID(ctx, userID)
	if err != nil {
		return "", err
	}

	if user.Email != identity.Email {
		return "", &domain.ForbiddenError{Reason: "token email does not match user email"}
	}

	return uc.userService.CreateTicketForUser(ctx, userID, packetID, eventID, uc.eventManagerService)
}

func (uc *userUsecase) GetCustomersByEventID(ctx context.Context, token string, eventID int) ([]*domain.User, error) {
	identity, err := uc.authenticate(ctx, token)
	if err != nil {
		return nil, err
	}

	allowed, err := uc.authZService.CanUserViewEventCustomers(ctx, identity, eventID)
	if err != nil {
		return nil, &domain.ForbiddenError{Reason: fmt.Sprintf("authorization check failed: %v", err)}
	}
	if !allowed {
		return nil, &domain.ForbiddenError{Reason: "only event owners can view customers"}
	}

	return uc.userService.GetCustomersByEventID(ctx, eventID)
}

func (uc *userUsecase) GetCustomersByPacketID(ctx context.Context, token string, packetID int) ([]*domain.User, error) {
	identity, err := uc.authenticate(ctx, token)
	if err != nil {
		return nil, err
	}

	allowed, err := uc.authZService.CanUserViewPacketCustomers(ctx, identity, packetID)
	if err != nil {
		return nil, &domain.ForbiddenError{Reason: fmt.Sprintf("authorization check failed: %v", err)}
	}
	if !allowed {
		return nil, &domain.ForbiddenError{Reason: "only packet owners can view customers"}
	}

	return uc.userService.GetCustomersByPacketID(ctx, packetID)
}
