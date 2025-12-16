package usecase

import (
	"context"
	"eventManager/application/domain"
	"eventManager/application/repository"
	"eventManager/application/service"
)

type TicketUseCase interface {
	CreateTicket(ctx context.Context, ticket *domain.Ticket) (*domain.Ticket, error)
	GetTicketByCode(ctx context.Context, code string) (*domain.Ticket, error)
	UpdateTicket(ctx context.Context, code string, updates map[string]interface{}) (*domain.Ticket, error)
	DeleteTicket(ctx context.Context, code string) (*domain.Ticket, error)
}

type ticketUseCase struct {
	repo          repository.TicketRepository
	ticketService service.TicketService 
}

func NewTicketUseCase(repo repository.TicketRepository, ticketService service.TicketService) *ticketUseCase {
	return &ticketUseCase{
		repo:          repo,
		ticketService: ticketService,
	}
}

func (uc *ticketUseCase) CreateTicket(ctx context.Context, ticket *domain.Ticket) (*domain.Ticket, error) {
	
	
	
	
	return uc.ticketService.CreateTicket(ctx, ticket)
}

func (uc *ticketUseCase) GetTicketByCode(ctx context.Context, code string) (*domain.Ticket, error) {
	if code == "" {
		return nil, &domain.ValidationError{Reason: "ticket code is required"}
	}
	return uc.repo.GetTicketByCode(ctx, code)
}

func (uc *ticketUseCase) UpdateTicket(ctx context.Context, code string, updates map[string]interface{}) (*domain.Ticket, error) {
	
	
	return uc.ticketService.UpdateTicket(ctx, code, updates)
}

func (uc *ticketUseCase) DeleteTicket(ctx context.Context, code string) (*domain.Ticket, error) {
	if code == "" {
		return nil, &domain.ValidationError{Reason: "ticket code is required"}
	}
	return uc.repo.DeleteEvent(ctx, code)
}
