package repository

import (
	"context"
	"eventManager/application/domain"
)

type EventPacketRepository interface {
	Create(ctx context.Context, eventPacket *domain.EventPacket) (*domain.EventPacket, error)
	GetByID(ctx context.Context, id int) (*domain.EventPacket, error)
	Update(ctx context.Context, id int, updates map[string]interface{}) (*domain.EventPacket, error)
	Delete(ctx context.Context, id int) (*domain.EventPacket, error)
	CountSoldTickets(ctx context.Context, id int) (int, error)
	FilterEventPackets(ctx context.Context, filter *domain.EventPacketFilter) ([]*domain.EventPacket, error)
	CountEventPackets(ctx context.Context, filter *domain.EventPacketFilter) (int, error)
}
