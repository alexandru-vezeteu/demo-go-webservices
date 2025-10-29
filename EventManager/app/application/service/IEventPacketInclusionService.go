package service

import "eventManager/domain"

type IEventPacketInclusionService interface {
	CreateEventPacketInclusion(event *domain.EventPacketInclusion) (*domain.EventPacketInclusion, error)
	GetEventsInPacketbyID(id int) (*domain.EventPacketInclusion, error)
	GetEventPacketsByEventID(id int) (*domain.EventPacketInclusion, error)
	//Update(event *domain.EventPacketInclusion) (*domain.EventPacketInclusion, error)
	DeleteEventPacketInclusion(event *domain.EventPacket) (*domain.EventPacket, error)
}
