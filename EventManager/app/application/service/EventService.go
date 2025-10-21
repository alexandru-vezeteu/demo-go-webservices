package service

import (
	"errors"
	"eventManager/application/repository"
	"eventManager/domain"
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
		return nil, errors.New("event is nil")
	}

	//nu dau check la id ca e creat automat oricum
	if event.OwnerID < 1 {
		return nil, errors.New("ownerid must be positive")
	}

	if event.Name == "" {
		return nil, errors.New("name must be set and unique")
	}

	//location, description, seats can be null
	ret, err := service.repo.Create(event)
	if err != nil {
		return nil, err
	}
	return ret, nil
}
func (service *eventService) GetEventByID(id int) (*domain.Event, error) {
	return nil, nil
}
func (service *eventService) UpdateEvent(event *domain.Event) (*domain.Event, error) {
	return nil, nil
}
func (service *eventService) DeleteEvent(id int) (*domain.Event, error) {
	return nil, nil
}
