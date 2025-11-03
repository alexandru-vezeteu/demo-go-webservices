package service

import "eventManager/domain"

type IEventPacketInclusionService interface {
	CreateEventPacketInclusion(event *domain.EventPacketInclusion) (*domain.EventPacketInclusion, error)
	GetEventsByPacketID(packetID int) ([]*domain.Event, error)
	GetEventPacketsByEventID(eventID int) ([]*domain.EventPacket, error)
	Update(eventID, packetID int, updates map[string]interface{}) (*domain.EventPacketInclusion, error)
	DeleteEventPacketInclusion(eventID, packetID int) (*domain.EventPacketInclusion, error)
}
