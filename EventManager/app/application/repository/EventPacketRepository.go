package repository

import (
	"context"
	"eventManager/application/domain"
)

type EventPacketRepository interface {
	Create(ctx context.Context, event *domain.EventPacket) (*domain.EventPacket, error)
	GetByID(ctx context.Context, id int) (*domain.EventPacket, error)
	Update(ctx context.Context, id int, updates map[string]interface{}) (*domain.EventPacket, error)
	Delete(ctx context.Context, id int) (*domain.EventPacket, error)
	CountSoldTickets(ctx context.Context, id int) (int, error)
}
