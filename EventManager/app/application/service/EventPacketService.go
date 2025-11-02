package service

import (
	"eventManager/application/repository"
	"eventManager/domain"
	"fmt"
)

type eventPacketService struct {
	repo repository.EventPacketRepository
}

func NewEventPacketService(repo repository.EventPacketRepository) *eventPacketService {
	return &eventPacketService{repo: repo}
}

func (service *eventPacketService) CreateEventPacket(event *domain.EventPacket) (*domain.EventPacket, error) {
	if err := service.validateEventPacket(event); err != nil {
		return nil, err
	}
	return service.repo.Create(event)
}

func (service *eventPacketService) GetEventPacketByID(id int) (*domain.EventPacket, error) {
	if id < 1 {
		return nil, domain.NewEventPacketValidationError(fmt.Sprintf("id:%d must be positive", id))
	}
	return service.repo.GetByID(id)
}

func (service *eventPacketService) UpdateEventPacket(id int, updates map[string]interface{}) (*domain.EventPacket, error) {

	if len(updates) == 0 {
		return nil, domain.NewEventPacketValidationError("no fields to update")
	}

	if owner_id, ok := updates["id_owner"]; ok {
		if owner_idPtr, ok := owner_id.(int); ok && owner_idPtr < 1 {
			return nil, domain.NewEventPacketValidationError("owner_id must be positive")

		}
	}

	if name, ok := updates["name"]; ok {
		if namePtr, ok := name.(string); ok && namePtr == "" {
			return nil, domain.NewEventPacketValidationError("name must be set")
		}
	}

	return service.repo.Update(id, updates)
}
func (service *eventPacketService) DeleteEventPacket(id int) (*domain.EventPacket, error) {
	if id < 1 {
		return nil, domain.NewEventPacketValidationError(fmt.Sprintf("id:%d must be positive", id))
	}
	return service.repo.Delete(id)
}

func (service *eventPacketService) validateEventPacket(event *domain.EventPacket) error {
	if event == nil {
		return domain.NewEventValidationError("invalid object received")
	}

	if event.OwnerID < 1 {
		return domain.NewEventValidationError("owner_id must be positive")
	}

	if event.Name == "" {
		return domain.NewEventValidationError("name must be set")
	}
	return nil
}
