package repository

import (
	"context"
	"eventManager/application/domain"
)

type TicketRepository interface {
	CreateTicket(ctx context.Context, event *domain.Ticket) (*domain.Ticket, error)
	GetTicketByCode(ctx context.Context, code string) (*domain.Ticket, error)
	UpdateTicket(ctx context.Context, code string, updates map[string]interface{}) (*domain.Ticket, error)
	ReplaceTicket(ctx context.Context, ticket *domain.Ticket) (*domain.Ticket, error)
	DeleteEvent(ctx context.Context, code string) (*domain.Ticket, error)
}
