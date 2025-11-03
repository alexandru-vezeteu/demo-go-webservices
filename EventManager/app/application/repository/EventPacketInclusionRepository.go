package repository

import "eventManager/domain"

type EventPacketInclusionRepository interface {
	Create(event *domain.EventPacketInclusion) (*domain.EventPacketInclusion, error)
	GetEventsByPacketID(packetID int) ([]*domain.Event, error)
	GetEventPacketsByEventID(eventID int) ([]*domain.EventPacket, error)
	Update(eventID, packetID int, updates map[string]interface{}) (*domain.EventPacketInclusion, error)
	Delete(eventID, packetID int) (*domain.EventPacketInclusion, error)
}
