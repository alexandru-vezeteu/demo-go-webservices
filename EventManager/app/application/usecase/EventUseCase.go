package usecase

import (
	"context"
	"eventManager/application/domain"
	"eventManager/application/repository"
	"eventManager/application/service"
	"fmt"
)

type EventUseCase interface {
	CreateEvent(ctx context.Context, event *domain.Event) (*domain.Event, error)
	GetEventByID(ctx context.Context, id int) (*domain.Event, error)
	UpdateEvent(ctx context.Context, id int, updates map[string]interface{}) (*domain.Event, error)
	DeleteEvent(ctx context.Context, id int) (*domain.Event, error)
	FilterEvents(ctx context.Context, filter *domain.EventFilter) ([]*domain.Event, error)
}

type eventUseCase struct {
	repo          repository.EventRepository
	eventService  service.EventService 
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

func (uc *eventUseCase) CreateEvent(ctx context.Context, event *domain.Event) (*domain.Event, error) {
	if err := uc.validateEvent(event); err != nil {
		return nil, err
	}
	return uc.repo.Create(ctx, event)
}

func (uc *eventUseCase) GetEventByID(ctx context.Context, id int) (*domain.Event, error) {
	if id < 1 {
		return nil, &domain.ValidationError{Reason: fmt.Sprintf("id:%d must be positive", id)}
	}
	return uc.repo.GetByID(ctx, id)
}

func (uc *eventUseCase) UpdateEvent(ctx context.Context, id int, updates map[string]interface{}) (*domain.Event, error) {
	
	
	return uc.eventService.UpdateEvent(ctx, id, updates)
}

func (uc *eventUseCase) DeleteEvent(ctx context.Context, id int) (*domain.Event, error) {
	if id < 1 {
		return nil, &domain.ValidationError{Reason: fmt.Sprintf("id:%d must be positive", id)}
	}
	return uc.repo.Delete(ctx, id)
}

func (uc *eventUseCase) FilterEvents(ctx context.Context, filter *domain.EventFilter) ([]*domain.Event, error) {
	return uc.repo.FilterEvents(ctx, filter)
}
