package service

import (
	"eventManager/application/repository"
	"eventManager/domain"
	"fmt"
)

// nu poate fi utilizat direct trb initializat cu NewEventService:p
type eventService struct {
	repo repository.EventRepository
}

func NewEventService(repo repository.EventRepository) *eventService {
	return &eventService{repo: repo}
}

func (service *eventService) CreateEvent(event *domain.Event) (*domain.Event, error) {
	if event == nil {
		return nil, domain.NewEventValidationError("invalid object received")
	}

	event.ID = 0
	//nu dau check la id ca e creat automat oricum
	if event.OwnerID < 1 {
		return nil, domain.NewEventValidationError("owner_id must be positive")
	}

	if event.Name == "" {
		return nil, domain.NewEventValidationError("name must be set")
	}

	return service.repo.Create(event)

}

func (service *eventService) GetEventByID(id int) (*domain.Event, error) {
	if id < 1 {
		return nil, domain.NewEventValidationError(fmt.Sprintf("id:%d must be positive", id))
	}
	return service.repo.GetByID(id)
}

func (service *eventService) UpdateEvent(event *domain.Event) (*domain.Event, error) {
	return nil, nil
}
func (service *eventService) DeleteEvent(id int) (*domain.Event, error) {
	return nil, nil
}
