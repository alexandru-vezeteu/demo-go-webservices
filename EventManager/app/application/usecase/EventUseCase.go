package usecase

import (
	"eventManager/application/domain"
	"eventManager/application/repository"
	"eventManager/application/service"
	"fmt"
)

type EventUseCase interface {
	CreateEvent(event *domain.Event) (*domain.Event, error)
	GetEventByID(id int) (*domain.Event, error)
	UpdateEvent(id int, updates map[string]interface{}) (*domain.Event, error)
	DeleteEvent(id int) (*domain.Event, error)
	FilterEvents(filter *domain.EventFilter) ([]*domain.Event, error)
}

type eventUseCase struct {
	repo          repository.EventRepository
	eventService  service.EventService // For complex business logic like seat validation
}

func NewEventUseCase(repo repository.EventRepository, eventService service.EventService) *eventUseCase {
	return &eventUseCase{
		repo:         repo,
		eventService: eventService,
	}
}

func (uc *eventUseCase) validateEvent(event *domain.Event) error {
	if event == nil {
		return &domain.ValidationError{Reason: "invalid object received"}
	}

	if event.OwnerID < 1 {
		return &domain.ValidationError{Reason: "owner_id must be positive"}
	}

	if event.Name == "" {
		return &domain.ValidationError{Reason: "name must be set"}
	}
	return nil
}

func (uc *eventUseCase) CreateEvent(event *domain.Event) (*domain.Event, error) {
	if err := uc.validateEvent(event); err != nil {
		return nil, err
	}
	return uc.repo.Create(event)
}

func (uc *eventUseCase) GetEventByID(id int) (*domain.Event, error) {
	if id < 1 {
		return nil, &domain.ValidationError{Reason: fmt.Sprintf("id:%d must be positive", id)}
	}
	return uc.repo.GetByID(id)
}

func (uc *eventUseCase) UpdateEvent(id int, updates map[string]interface{}) (*domain.Event, error) {
	// UpdateEvent has complex business logic (validateSeatsAgainstPackets)
	// so we delegate to the service
	return uc.eventService.UpdateEvent(id, updates)
}

func (uc *eventUseCase) DeleteEvent(id int) (*domain.Event, error) {
	if id < 1 {
		return nil, &domain.ValidationError{Reason: fmt.Sprintf("id:%d must be positive", id)}
	}
	return uc.repo.Delete(id)
}

func (uc *eventUseCase) FilterEvents(filter *domain.EventFilter) ([]*domain.Event, error) {
	return uc.repo.FilterEvents(filter)
}
