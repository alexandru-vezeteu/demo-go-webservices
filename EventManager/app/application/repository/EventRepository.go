package repository

import (
	"context"
	"eventManager/application/domain"
)

type EventRepository interface {
	Create(ctx context.Context, event *domain.Event) (*domain.Event, error)
	GetByID(ctx context.Context, id int) (*domain.Event, error)
	Update(ctx context.Context, id int, updates map[string]interface{}) (*domain.Event, error)
	Delete(ctx context.Context, id int) (*domain.Event, error)
	FilterEvents(ctx context.Context, filter *domain.EventFilter) ([]*domain.Event, error)
	CountEvents(ctx context.Context, filter *domain.EventFilter) (int, error)
	CountSoldTickets(ctx context.Context, id int) (int, error)
}
