package service

import (
	"eventManager/application/repository"
	"eventManager/domain"
)

type eventPacketService struct {
	repo repository.EventPacketRepository
}

func NewEventPacketService(repo repository.EventPacketRepository) *eventPacketService {
	return &eventPacketService{repo: repo}
}

func (service *eventPacketService) CreateEventPacket(event *domain.EventPacket) (*domain.EventPacket, error) {
	if event == nil {
		return nil, domain.NewEventValidationError("invalid object received")
	}

	if event.OwnerID < 1 {
		return nil, domain.NewEventValidationError("owner id must be positive")
	}

	if event.Name == "" {
		return nil, domain.NewEventValidationError("name must be set")
	}

	return service.repo.Create(event)

}

func (service *eventPacketService) GetEventPacketByID(id int) (*domain.EventPacket, error) {
	return nil, nil
}

func (service *eventPacketService) UpdateEventPacket(event *domain.EventPacket) (*domain.EventPacket, error) {
	return nil, nil
}
func (service *eventPacketService) DeleteEventPacket(id int) (*domain.EventPacket, error) {
	return nil, nil
}
