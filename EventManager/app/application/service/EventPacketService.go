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
		return nil, &domain.ValidationError{Msg: fmt.Sprintf("id:%d must be positive", id)}
	}
	return service.repo.GetByID(id)
}

func (service *eventPacketService) UpdateEventPacket(id int, updates map[string]interface{}) (*domain.EventPacket, error) {

	if len(updates) == 0 {
		return nil, &domain.ValidationError{Msg: "no fields to update"}
	}

	if owner_id, ok := updates["id_owner"]; ok {
		if owner_idPtr, ok := owner_id.(int); ok && owner_idPtr < 1 {
			return nil, &domain.ValidationError{Msg: "owner_id must be positive"}

		}
	}

	if name, ok := updates["name"]; ok {
		if namePtr, ok := name.(string); ok && namePtr == "" {
			return nil, &domain.ValidationError{Msg: "name must be set"}
		}
	}

	return service.repo.Update(id, updates)
}
func (service *eventPacketService) DeleteEventPacket(id int) (*domain.EventPacket, error) {
	if id < 1 {
		return nil, &domain.ValidationError{Msg: fmt.Sprintf("id:%d must be positive", id)}
	}
	return service.repo.Delete(id)
}

func (service *eventPacketService) validateEventPacket(event *domain.EventPacket) error {
	if event == nil {
		return &domain.ValidationError{Msg: "invalid object received"}
	}

	if event.OwnerID < 1 {
		return &domain.ValidationError{Msg: "owner_id must be positive"}
	}

	if event.Name == "" {
		return &domain.ValidationError{Msg: "name must be set"}
	}
	return nil
}
