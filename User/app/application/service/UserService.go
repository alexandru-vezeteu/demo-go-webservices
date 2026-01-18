package service

import (
	"context"
	"regexp"
	"strings"
	"userService/application/domain"
	"userService/application/repository"

	"github.com/google/uuid"
)

type UserService interface {
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	GetUserByID(ctx context.Context, id int) (*domain.User, error)
	UpdateUser(ctx context.Context, id int, updates map[string]interface{}) (*domain.User, error)
	DeleteUser(ctx context.Context, id int) (*domain.User, error)
	CreateTicketForUser(ctx context.Context, userID int, packetID *int, eventID *int, ticketCreator TicketCreator) (string, error)
	GetCustomersByEventID(ctx context.Context, eventID int) ([]*domain.User, error)
	GetCustomersByPacketID(ctx context.Context, packetID int) ([]*domain.User, error)
}

type TicketCreator interface {
	CreateTicket(ctx context.Context, code string, packetID *int, eventID *int) (*TicketResponse, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) validateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func (s *userService) validateUser(user *domain.User) error {
	if strings.TrimSpace(user.Email) == "" {
		return &domain.ValidationError{Field: "email", Reason: "Email is required"}
	}

	if !s.validateEmail(user.Email) {
		return &domain.ValidationError{Field: "email", Reason: "Invalid email format"}
	}

	if user.FirstName != "" && len(strings.TrimSpace(user.FirstName)) > 100 {
		return &domain.ValidationError{Field: "first_name", Reason: "First name too long (max 100 characters)"}
	}

	if user.LastName != "" && len(strings.TrimSpace(user.LastName)) > 100 {
		return &domain.ValidationError{Field: "last_name", Reason: "Last name too long (max 100 characters)"}
	}

	if user.SocialMediaLinks != nil {
		trimmed := strings.TrimSpace(*user.SocialMediaLinks)
		if trimmed == "" {
			return &domain.ValidationError{Field: "social_media_links", Reason: "Social media links cannot be empty string"}
		}
	}

	return nil
}

func (s *userService) validateUpdates(updates map[string]interface{}) error {
	if email, ok := updates["email"]; ok {
		emailStr, isString := email.(string)
		if !isString || strings.TrimSpace(emailStr) == "" {
			return &domain.ValidationError{Field: "email", Reason: "Email cannot be empty"}
		}
		if !s.validateEmail(emailStr) {
			return &domain.ValidationError{Field: "email", Reason: "Invalid email format"}
		}
	}

	if firstName, ok := updates["first_name"]; ok {
		firstNameStr, isString := firstName.(string)
		if !isString || strings.TrimSpace(firstNameStr) == "" {
			return &domain.ValidationError{Field: "first_name", Reason: "First name cannot be empty"}
		}
	}

	if lastName, ok := updates["last_name"]; ok {
		lastNameStr, isString := lastName.(string)
		if !isString || strings.TrimSpace(lastNameStr) == "" {
			return &domain.ValidationError{Field: "last_name", Reason: "Last name cannot be empty"}
		}
	}

	if socialMediaLinks, ok := updates["social_media_links"]; ok {
		if socialMediaLinks != nil {
			linksStr, isString := socialMediaLinks.(string)
			if isString && strings.TrimSpace(linksStr) == "" {
				return &domain.ValidationError{Field: "social_media_links", Reason: "Social media links cannot be empty string"}
			}
		}
	}

	return nil
}

func (s *userService) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	if err := s.validateUser(user); err != nil {
		return nil, err
	}
	return s.repo.Create(ctx, user)
}

func (s *userService) GetUserByID(ctx context.Context, id int) (*domain.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *userService) UpdateUser(ctx context.Context, id int, updates map[string]interface{}) (*domain.User, error) {
	if err := s.validateUpdates(updates); err != nil {
		return nil, err
	}

	delete(updates, "ticket_list")

	return s.repo.Update(ctx, id, updates)
}

func (s *userService) DeleteUser(ctx context.Context, id int) (*domain.User, error) {
	return s.repo.Delete(ctx, id)
}

func (s *userService) CreateTicketForUser(ctx context.Context, userID int, packetID *int, eventID *int, ticketCreator TicketCreator) (string, error) {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return "", err
	}

	ticketCode := uuid.New().String()

	ticketResp, err := ticketCreator.CreateTicket(ctx, ticketCode, packetID, eventID)
	if err != nil {
		return "", err
	}

	newTicket := domain.Ticket{
		PacketID: packetID,
		EventID:  eventID,
		Code:     ticketCode,
	}
	_ = ticketResp

	updatedTicketList := append(user.TicketList, newTicket)

	updates := map[string]interface{}{
		"ticket_list": updatedTicketList,
	}

	_, err = s.repo.Update(ctx, userID, updates)
	if err != nil {
		return "", &domain.InternalError{Msg: "failed to update user ticket list", Err: err}
	}

	return ticketCode, nil
}

func (s *userService) filterPrivateFields(users []*domain.User) []*domain.User {
	filteredUsers := make([]*domain.User, 0, len(users))
	for _, u := range users {
		filtered := &domain.User{
			ID:               u.ID,
			Email:            u.Email,
			FirstName:        u.FirstName,
			LastName:         u.LastName,
			SocialMediaLinks: u.SocialMediaLinks,
			TicketList:       nil,
		}

		if u.FirstNamePrivate {
			filtered.FirstName = "[Private]"
		}
		if u.LastNamePrivate {
			filtered.LastName = "[Private]"
		}

		filteredUsers = append(filteredUsers, filtered)
	}
	return filteredUsers
}

func (s *userService) GetCustomersByEventID(ctx context.Context, eventID int) ([]*domain.User, error) {
	users, err := s.repo.GetUsersByEventID(ctx, eventID)
	if err != nil {
		return nil, err
	}
	return s.filterPrivateFields(users), nil
}

func (s *userService) GetCustomersByPacketID(ctx context.Context, packetID int) ([]*domain.User, error) {
	users, err := s.repo.GetUsersByPacketID(ctx, packetID)
	if err != nil {
		return nil, err
	}
	return s.filterPrivateFields(users), nil
}
