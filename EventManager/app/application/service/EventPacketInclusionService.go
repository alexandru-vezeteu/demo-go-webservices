package service

import (
	"eventManager/application/domain"
	"eventManager/application/repository"
)

type eventPacketInclusionService struct {
	repo repository.EventPacketInclusionRepository
}

func NewEventPacketInclusionService(repo repository.EventPacketInclusionRepository) *eventPacketInclusionService {
	return &eventPacketInclusionService{repo: repo}
}

func (service *eventPacketInclusionService) CreateEventPacketInclusion(event *domain.EventPacketInclusion) (*domain.EventPacketInclusion, error) {

	if event == nil {
		return nil, &domain.ValidationError{Msg: "invalid object"}
	}
	if event.AllocatedSeats != nil && *event.AllocatedSeats < 0 {
		return nil, &domain.ValidationError{Msg: "allocated seats must be >= 0"}
	}
	if event.EventID < 0 {
		return nil, &domain.ValidationError{Msg: "event id must be >= 0"}
	}
	if event.PacketID < 0 {
		return nil, &domain.ValidationError{Msg: "packet id must be >= 0"}
	}
	return service.repo.Create(event)

}
func (service *eventPacketInclusionService) GetEventsByPacketID(packetID int) ([]*domain.Event, error) {
	if packetID < 0 {
		return nil, &domain.ValidationError{Msg: "packet id must be >= 0"}
	}
	return service.repo.GetEventsByPacketID(packetID)
}
func (service *eventPacketInclusionService) GetEventPacketsByEventID(eventID int) ([]*domain.EventPacket, error) {
	if eventID < 0 {
		return nil, &domain.ValidationError{Msg: "event id must be >= 0"}
	}
	return service.repo.GetEventPacketsByEventID(eventID)
}
func (service *eventPacketInclusionService) Update(eventID, packetID int, updates map[string]interface{}) (*domain.EventPacketInclusion, error) {
	if eventID < 0 {
		return nil, &domain.ValidationError{Msg: "event id must be >= 0"}
	}
	if packetID < 0 {
		return nil, &domain.ValidationError{Msg: "packet id must be >= 0"}
	}
	if len(updates) == 0 {
		return nil, &domain.ValidationError{Msg: "update must contain at least one field"}
	}
	return service.repo.Update(eventID, packetID, updates)
}
func (service *eventPacketInclusionService) DeleteEventPacketInclusion(eventID, packetID int) (*domain.EventPacketInclusion, error) {
	if eventID < 0 {
		return nil, &domain.ValidationError{Msg: "event id must be >= 0"}
	}
	if packetID < 0 {
		return nil, &domain.ValidationError{Msg: "packet id must be >= 0"}
	}

	return service.repo.Delete(eventID, packetID)
}
