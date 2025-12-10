package usecase

import (
	"eventManager/application/domain"
	"eventManager/application/repository"
	"eventManager/application/service"
)

type TicketUseCase interface {
	CreateTicket(ticket *domain.Ticket) (*domain.Ticket, error)
	GetTicketByCode(code string) (*domain.Ticket, error)
	UpdateTicket(code string, updates map[string]interface{}) (*domain.Ticket, error)
	DeleteTicket(code string) (*domain.Ticket, error)
}

type ticketUseCase struct {
	repo          repository.TicketRepository
	ticketService service.TicketService // For complex business logic (UUID generation, seat availability, constraint validation)
}

func NewTicketUseCase(repo repository.TicketRepository, ticketService service.TicketService) *ticketUseCase {
	return &ticketUseCase{
		repo:          repo,
		ticketService: ticketService,
	}
}

func (uc *ticketUseCase) CreateTicket(ticket *domain.Ticket) (*domain.Ticket, error) {
	// CreateTicket has complex business logic:
	// - UUID generation
	// - Complex seat availability validation (cross-entity calculations)
	// Delegate to service
	return uc.ticketService.CreateTicket(ticket)
}

func (uc *ticketUseCase) GetTicketByCode(code string) (*domain.Ticket, error) {
	if code == "" {
		return nil, &domain.ValidationError{Reason: "ticket code is required"}
	}
	return uc.repo.GetTicketByCode(code)
}

func (uc *ticketUseCase) UpdateTicket(code string, updates map[string]interface{}) (*domain.Ticket, error) {
	// UpdateTicket has complex constraint validation when moving tickets between events/packets
	// Delegate to service
	return uc.ticketService.UpdateTicket(code, updates)
}

func (uc *ticketUseCase) DeleteTicket(code string) (*domain.Ticket, error) {
	if code == "" {
		return nil, &domain.ValidationError{Reason: "ticket code is required"}
	}
	return uc.repo.DeleteEvent(code)
}
