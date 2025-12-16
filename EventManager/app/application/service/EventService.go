package service

import (
	"context"
	"eventManager/application/domain"
	"eventManager/application/repository"
	"fmt"
)

type EventService interface {
	CreateEvent(ctx context.Context, event *domain.Event) (*domain.Event, error)
	GetEventByID(ctx context.Context, id int) (*domain.Event, error)
	UpdateEvent(ctx context.Context, id int, updates map[string]interface{}) (*domain.Event, error)
	DeleteEvent(ctx context.Context, id int) (*domain.Event, error)
	FilterEvents(ctx context.Context, filter *domain.EventFilter) ([]*domain.Event, error)
}


type eventService struct {
	repo          repository.EventRepository
	inclusionRepo repository.EventPacketInclusionRepository
}

func NewEventService(repo repository.EventRepository, inclusionRepo repository.EventPacketInclusionRepository) EventService {
	return &eventService{
		repo:          repo,
		inclusionRepo: inclusionRepo,
	}
}

func (service *eventService) validateEvent(event *domain.Event) error {
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

func (service *eventService) CreateEvent(ctx context.Context, event *domain.Event) (*domain.Event, error) {
	if err := service.validateEvent(event); err != nil {
		return nil, err
	}
	return service.repo.Create(ctx, event)
}

func (service *eventService) GetEventByID(ctx context.Context, id int) (*domain.Event, error) {
	if id < 1 {
		return nil, &domain.ValidationError{Reason: fmt.Sprintf("id:%d must be positive", id)}
	}
	return service.repo.GetByID(ctx, id)
}

func (service *eventService) UpdateEvent(ctx context.Context, id int, updates map[string]interface{}) (*domain.Event, error) {

	if len(updates) == 0 {
		return nil, &domain.ValidationError{Reason: "no fields to update"}
	}

	if seats, ok := updates["seats"]; ok {
		if seatsPtr, ok := seats.(int); ok {
			if seatsPtr < 0 {
				return nil, &domain.ValidationError{Reason: "seats cannot be negative"}
			}
			
			if err := service.validateSeatsAgainstPackets(ctx, id, seatsPtr); err != nil {
				return nil, err
			}
		}
	}

	if owner_id, ok := updates["id_owner"]; ok {
		if owner_idPtr, ok := owner_id.(int); ok && owner_idPtr < 1 {
			return nil, &domain.ValidationError{Reason: "owner_id must be positive"}

		}
	}

	if name, ok := updates["name"]; ok {
		if namePtr, ok := name.(string); ok && namePtr == "" {
			return nil, &domain.ValidationError{Reason: "name must be set"}
		}
	}

	return service.repo.Update(ctx, id, updates)
}


func (service *eventService) validateSeatsAgainstPackets(ctx context.Context, eventID int, newSeats int) error {
	
	packets, err := service.inclusionRepo.GetEventPacketsByEventID(ctx, eventID)
	if err != nil {
		return err
	}

	
	for _, packet := range packets {
		if packet.AllocatedSeats != nil && newSeats < *packet.AllocatedSeats {
			return &domain.ValidationError{
				Reason: fmt.Sprintf("cannot reduce seats to %d: event is in packet '%s' which requires %d allocated seats", newSeats, packet.Name, *packet.AllocatedSeats),
			}
		}
	}

	return nil
}

func (service *eventService) DeleteEvent(ctx context.Context, id int) (*domain.Event, error) {
	if id < 1 {
		return nil, &domain.ValidationError{Reason: fmt.Sprintf("id:%d must be positive", id)}
	}
	return service.repo.Delete(ctx, id)
}

func (service *eventService) FilterEvents(ctx context.Context, filter *domain.EventFilter) ([]*domain.Event, error) {
	return service.repo.FilterEvents(ctx, filter)
}
