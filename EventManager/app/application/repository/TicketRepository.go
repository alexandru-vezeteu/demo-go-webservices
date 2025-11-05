package repository

import "eventManager/application/domain"

type TicketRepository interface {
	CreateTicket(event *domain.Ticket) (*domain.Ticket, error)
	GetTicketByCode(code string) (*domain.Ticket, error)
	UpdateTicket(code string, updates map[string]interface{}) (*domain.Ticket, error)
	DeleteEvent(code string) (*domain.Ticket, error)
}
