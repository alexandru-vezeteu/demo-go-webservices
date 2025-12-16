package repository

import (
	"context"
	"eventManager/application/domain"
)

type EventPacketInclusionRepository interface {
	Create(ctx context.Context, event *domain.EventPacketInclusion) (*domain.EventPacketInclusion, error)
	GetEventsByPacketID(ctx context.Context, packetID int) ([]*domain.Event, error)
	GetEventPacketsByEventID(ctx context.Context, eventID int) ([]*domain.EventPacket, error)
	Update(ctx context.Context, eventID, packetID int, updates map[string]interface{}) (*domain.EventPacketInclusion, error)
	Delete(ctx context.Context, eventID, packetID int) (*domain.EventPacketInclusion, error)
}
