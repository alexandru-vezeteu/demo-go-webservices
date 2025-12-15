package usecase

import (
	"fmt"
	"regexp"
	"strings"
	"userService/application/domain"
	"userService/application/repository"
	"userService/infrastructure/grpc"
	"userService/infrastructure/http"

	"github.com/google/uuid"
)

type UserUsecase interface {
	CreateUser(user *domain.User) (*domain.User, error)
	GetUserByID(id int) (*domain.User, error)
	UpdateUser(id int, updates map[string]interface{}) (*domain.User, error)
	DeleteUser(id int) (*domain.User, error)
	CreateTicketForUser(userID int, token string, packetID *int, eventID *int) (string, error)
}

type userUsecase struct {
	repo               repository.UserRepository
	idmClient          *grpc.IDMClient
	eventManagerClient *http.EventManagerClient
}

func NewUserUsecase(repo repository.UserRepository, idmClient *grpc.IDMClient, eventManagerClient *http.EventManagerClient) UserUsecase {
	return &userUsecase{
		repo:               repo,
		idmClient:          idmClient,
		eventManagerClient: eventManagerClient,
	}
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

	if user.TicketList != nil {
		trimmed := strings.TrimSpace(*user.TicketList)
		if trimmed == "" {
			return &domain.ValidationError{Field: "ticket_list", Reason: "Ticket list cannot be empty string"}
		}
	}

	return nil
}

func (uc *userUsecase) CreateUser(user *domain.User) (*domain.User, error) {
	if err := uc.validateUser(user); err != nil {
		return nil, err
	}
	return uc.repo.Create(user)
}

func (uc *userUsecase) GetUserByID(id int) (*domain.User, error) {
	return uc.repo.GetByID(id)
}

func (uc *userUsecase) UpdateUser(id int, updates map[string]interface{}) (*domain.User, error) {
	// Validate partial updates
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

	if ticketList, ok := updates["ticket_list"]; ok {
		if ticketList != nil {
			ticketListStr, isString := ticketList.(string)
			if isString && strings.TrimSpace(ticketListStr) == "" {
				return nil, &domain.ValidationError{Field: "ticket_list", Reason: "Ticket list cannot be empty string"}
			}
		}
	}

	return uc.repo.Update(id, updates)
}

func (uc *userUsecase) DeleteUser(id int) (*domain.User, error) {
	return uc.repo.Delete(id)
}

func (uc *userUsecase) CreateTicketForUser(userID int, token string, packetID *int, eventID *int) (string, error) {
	// a) Validate token
	tokenResp, err := uc.idmClient.VerifyToken(token)
	if err != nil {
		return "", &domain.InternalError{Msg: "failed to verify token", Err: err}
	}

	// Check if token is valid
	if !tokenResp.Valid {
		return "", &domain.ValidationError{Field: "token", Reason: fmt.Sprintf("invalid token: %s", tokenResp.Message)}
	}

	// b) EMAIL FROM TOKEN CLAIMS
	tokenEmail := tokenResp.Email

	// c) EMAIL FROM DB BASED ON user_id
	user, err := uc.repo.GetByID(userID)
	if err != nil {
		return "", err
	}

	// d) EMAIL TOKEN == EMAIL DB => IF NOT 403
	if user.Email != tokenEmail {
		return "", &domain.ForbiddenError{Reason: "token email does not match user email"}
	}

	// e) GENERATE UUID FOR TICKET CODE
	ticketCode := uuid.New().String()

	// f) PUT /TICKETS/{CODE} -> IF 201 -> APPEND TO TICKETS LIST ELSE => FAIL
	_, statusCode, err := uc.eventManagerClient.CreateTicket(ticketCode, packetID, eventID)
	if err != nil || statusCode != 201 {
		if statusCode >= 500 {
			return "", &domain.InternalError{Msg: "event manager service unavailable", Err: err}
		}
		return "", &domain.InternalError{Msg: "failed to create ticket", Err: err}
	}

	// Append ticket code to user's ticket_list
	var newTicketList string
	if user.TicketList == nil || strings.TrimSpace(*user.TicketList) == "" {
		newTicketList = ticketCode
	} else {
		newTicketList = *user.TicketList + "," + ticketCode
	}

	updates := map[string]interface{}{
		"ticket_list": newTicketList,
	}

	_, err = uc.repo.Update(userID, updates)
	if err != nil {
		return "", &domain.InternalError{Msg: "failed to update user ticket list", Err: err}
	}

	return ticketCode, nil
}
