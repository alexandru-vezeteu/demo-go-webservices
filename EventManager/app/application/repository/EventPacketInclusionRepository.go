package repository

import "eventManager/domain"

type EventPacketInclusionRepository interface {
	Create(event *domain.EventPacketInclusion) (*domain.EventPacketInclusion, error)
	GetEventsInPacketbyID(id int) (*domain.EventPacketInclusion, error)
	GetEventPacketsByEventID(id int) (*domain.EventPacketInclusion, error)
	//Update(event *domain.EventPacketInclusion) (*domain.EventPacketInclusion, error)
	Delete(event *domain.EventPacketInclusion) (*domain.EventPacketInclusion, error)
}
