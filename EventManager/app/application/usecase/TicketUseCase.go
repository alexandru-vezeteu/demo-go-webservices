package usecase

import (
	"context"
	"eventManager/application/domain"
	"eventManager/application/repository"
	"eventManager/application/service"
	"fmt"
)

type TicketUseCase interface {
	CreateTicket(ctx context.Context, token string, ticket *domain.Ticket) (*domain.Ticket, error)
	PutTicket(ctx context.Context, token string, code string, ticket *domain.Ticket) (*domain.Ticket, error)
	GetTicketByCode(ctx context.Context, token string, code string) (*domain.Ticket, error)
	UpdateTicket(ctx context.Context, token string, code string, updates map[string]interface{}) (*domain.Ticket, error)
	DeleteTicket(ctx context.Context, token string, code string) (*domain.Ticket, error)
}

type ticketUseCase struct {
	repo          repository.TicketRepository
	ticketService service.TicketService
	authNService  service.AuthenticationService
	authZService  service.AuthorizationService
}

func NewTicketUseCase(
	repo repository.TicketRepository,
	ticketService service.TicketService,
	authNService service.AuthenticationService,
	authZService service.AuthorizationService,
) *ticketUseCase {
	return &ticketUseCase{
		repo:          repo,
		ticketService: ticketService,
		authNService:  authNService,
		authZService:  authZService,
	}
}

func (uc *ticketUseCase) authenticate(ctx context.Context, token string) (*service.UserIdentity, error) {
	identity, err := uc.authNService.WhoIsUser(ctx, token)
	if err != nil {
		return nil, &domain.ValidationError{Reason: "invalid or expired token"}
	}
	return identity, nil
}

func (uc *ticketUseCase) CreateTicket(ctx context.Context, token string, ticket *domain.Ticket) (*domain.Ticket, error) {
	identity, err := uc.authenticate(ctx, token)
	if err != nil {
		return nil, err
	}

	allowed, err := uc.authZService.CanUserCreateTicket(ctx, *identity)
	if err != nil {
		return nil, &domain.InternalError{Msg: fmt.Sprintf("authorization check failed: %v", err)}
	}
	if !allowed {
		return nil, &domain.ForbiddenError{Reason: "you don't have permission to create tickets"}
	}

	return uc.ticketService.CreateTicket(ctx, ticket)
}

func (uc *ticketUseCase) PutTicket(ctx context.Context, token string, code string, ticket *domain.Ticket) (*domain.Ticket, error) {
	if code == "" {
		return nil, &domain.ValidationError{Reason: "ticket code is required"}
	}

	identity, err := uc.authenticate(ctx, token)
	if err != nil {
		return nil, err
	}

	ticket.Code = code

	allowed, err := uc.authZService.CanUserCreateTicket(ctx, *identity)
	if err != nil {
		return nil, &domain.InternalError{Msg: fmt.Sprintf("authorization check failed: %v", err)}
	}
	if !allowed {
		return nil, &domain.ForbiddenError{Reason: "you don't have permission to create this ticket"}
	}
	return uc.ticketService.ReplaceTicket(ctx, ticket)
}

func (uc *ticketUseCase) GetTicketByCode(ctx context.Context, token string, code string) (*domain.Ticket, error) {
	identity, err := uc.authenticate(ctx, token)
	if err != nil {
		return nil, err
	}

	ticket, err := uc.ticketService.GetTicketByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	allowed, err := uc.authZService.CanUserViewTicket(ctx, *identity, ticket)
	if err != nil {
		return nil, &domain.InternalError{Msg: fmt.Sprintf("authorization check failed: %v", err)}
	}
	if !allowed {
		return nil, &domain.ForbiddenError{Reason: "you don't have permission to view this ticket"}
	}

	return ticket, nil
}

func (uc *ticketUseCase) UpdateTicket(ctx context.Context, token string, code string, updates map[string]interface{}) (*domain.Ticket, error) {
	identity, err := uc.authenticate(ctx, token)
	if err != nil {
		return nil, err
	}

	ticket, err := uc.ticketService.GetTicketByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	allowed, err := uc.authZService.CanUserEditTicket(ctx, *identity, ticket)
	if err != nil {
		return nil, &domain.InternalError{Msg: fmt.Sprintf("authorization check failed: %v", err)}
	}
	if !allowed {
		return nil, &domain.ForbiddenError{Reason: "you don't have permission to edit this ticket"}
	}

	return uc.ticketService.UpdateTicket(ctx, code, updates)
}

func (uc *ticketUseCase) DeleteTicket(ctx context.Context, token string, code string) (*domain.Ticket, error) {
	identity, err := uc.authenticate(ctx, token)
	if err != nil {
		return nil, err
	}

	ticket, err := uc.ticketService.GetTicketByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	allowed, err := uc.authZService.CanUserDeleteTicket(ctx, *identity, ticket)
	if err != nil {
		return nil, &domain.InternalError{Msg: fmt.Sprintf("authorization check failed: %v", err)}
	}
	if !allowed {
		return nil, &domain.ForbiddenError{Reason: "you don't have permission to delete this ticket"}
	}

	return uc.ticketService.DeleteTicket(ctx, code)
}
