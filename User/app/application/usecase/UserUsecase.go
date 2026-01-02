package usecase

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"userService/application/domain"
	"userService/application/repository"
	"userService/application/service"

	"github.com/google/uuid"
)

type UserUsecase interface {
	CreateUser(ctx context.Context, token string, user *domain.User) (*domain.User, error)
	GetUserByID(ctx context.Context, token string, id int) (*domain.User, error)
	UpdateUser(ctx context.Context, token string, id int, updates map[string]interface{}) (*domain.User, error)
	DeleteUser(ctx context.Context, token string, id int) (*domain.User, error)
	CreateTicketForUser(ctx context.Context, userID int, token string, packetID *int, eventID *int) (string, error)
}

type userUsecase struct {
	repo                repository.UserRepository
	eventManagerService service.EventManagerService
	authNService        service.AuthenticationService
	authZService        service.AuthorizationService
}

func NewUserUsecase(
	repo repository.UserRepository,
	eventManagerService service.EventManagerService,
	authNService service.AuthenticationService,
	authZService service.AuthorizationService,
) UserUsecase {
	return &userUsecase{
		repo:                repo,
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

func (uc *userUsecase) validateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func (uc *userUsecase) validateUser(user *domain.User) error {
	if strings.TrimSpace(user.Email) == "" {
		return &domain.ValidationError{Field: "email", Reason: "Email is required"}
	}

	if !uc.validateEmail(user.Email) {
		return &domain.ValidationError{Field: "email", Reason: "Invalid email format"}
	}

	if strings.TrimSpace(user.FirstName) == "" {
		return &domain.ValidationError{Field: "first_name", Reason: "First name is required"}
	}

	if strings.TrimSpace(user.LastName) == "" {
		return &domain.ValidationError{Field: "last_name", Reason: "Last name is required"}
	}

	if user.SocialMediaLinks != nil {
		trimmed := strings.TrimSpace(*user.SocialMediaLinks)
		if trimmed == "" {
			return &domain.ValidationError{Field: "social_media_links", Reason: "Social media links cannot be empty string"}
		}
	}

	return nil
}

func (uc *userUsecase) CreateUser(ctx context.Context, token string, user *domain.User) (*domain.User, error) {
	_, err := uc.authenticate(ctx, token)
	if err != nil {
		return nil, err
	}

	if err := uc.validateUser(user); err != nil {
		return nil, err
	}
	return uc.repo.Create(ctx, user)
}

func (uc *userUsecase) GetUserByID(ctx context.Context, token string, id int) (*domain.User, error) {
	identity, err := uc.authenticate(ctx, token)
	if err != nil {
		return nil, err
	}

	user, err := uc.repo.GetByID(ctx, id)
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

	user, err := uc.repo.GetByID(ctx, id)
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

	if email, ok := updates["email"]; ok {
		emailStr, isString := email.(string)
		if !isString || strings.TrimSpace(emailStr) == "" {
			return nil, &domain.ValidationError{Field: "email", Reason: "Email cannot be empty"}
		}
		if !uc.validateEmail(emailStr) {
			return nil, &domain.ValidationError{Field: "email", Reason: "Invalid email format"}
		}
	}

	if firstName, ok := updates["first_name"]; ok {
		firstNameStr, isString := firstName.(string)
		if !isString || strings.TrimSpace(firstNameStr) == "" {
			return nil, &domain.ValidationError{Field: "first_name", Reason: "First name cannot be empty"}
		}
	}

	if lastName, ok := updates["last_name"]; ok {
		lastNameStr, isString := lastName.(string)
		if !isString || strings.TrimSpace(lastNameStr) == "" {
			return nil, &domain.ValidationError{Field: "last_name", Reason: "Last name cannot be empty"}
		}
	}

	if socialMediaLinks, ok := updates["social_media_links"]; ok {
		if socialMediaLinks != nil {
			linksStr, isString := socialMediaLinks.(string)
			if isString && strings.TrimSpace(linksStr) == "" {
				return nil, &domain.ValidationError{Field: "social_media_links", Reason: "Social media links cannot be empty string"}
			}
		}
	}

	delete(updates, "ticket_list")

	return uc.repo.Update(ctx, id, updates)
}

func (uc *userUsecase) DeleteUser(ctx context.Context, token string, id int) (*domain.User, error) {
	identity, err := uc.authenticate(ctx, token)
	if err != nil {
		return nil, err
	}

	user, err := uc.repo.GetByID(ctx, id)
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

	return uc.repo.Delete(ctx, id)
}

func (uc *userUsecase) CreateTicketForUser(ctx context.Context, userID int, token string, packetID *int, eventID *int) (string, error) {
	identity, err := uc.authNService.WhoIsUser(ctx, token)
	if err != nil {
		return "", &domain.ValidationError{Field: "token", Reason: "invalid or expired token"}
	}

	user, err := uc.repo.GetByID(ctx, userID)
	if err != nil {
		return "", err
	}

	if user.Email != identity.Email {
		return "", &domain.ForbiddenError{Reason: "token email does not match user email"}
	}

	ticketCode := uuid.New().String()

	ticketResp, err := uc.eventManagerService.CreateTicket(ctx, ticketCode, packetID, eventID)
	if err != nil {
		return "", err
	}

	newTicket := domain.Ticket{
		Code: ticketCode,
	}
	_ = ticketResp

	updatedTicketList := append(user.TicketList, newTicket)

	updates := map[string]interface{}{
		"ticket_list": updatedTicketList,
	}

	_, err = uc.repo.Update(ctx, userID, updates)
	if err != nil {
		return "", &domain.InternalError{Msg: "failed to update user ticket list", Err: err}
	}

	return ticketCode, nil
}
