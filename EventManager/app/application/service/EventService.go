package service

import (
	"eventManager/application/domain"
	"eventManager/application/repository"
	"fmt"
)

// nu poate fi utilizat direct trb initializat cu NewEventService:p
type eventService struct {
	repo          repository.EventRepository
	inclusionRepo repository.EventPacketInclusionRepository
}

func NewEventService(repo repository.EventRepository, inclusionRepo repository.EventPacketInclusionRepository) *eventService {
	return &eventService{
		repo:          repo,
		inclusionRepo: inclusionRepo,
	}
}

func (service *eventService) validateEvent(event *domain.Event) error {
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

func (service *eventService) CreateEvent(event *domain.Event) (*domain.Event, error) {
	if err := service.validateEvent(event); err != nil {
		return nil, err
	}
	return service.repo.Create(event)
}

func (service *eventService) GetEventByID(id int) (*domain.Event, error) {
	if id < 1 {
		return nil, &domain.ValidationError{Msg: fmt.Sprintf("id:%d must be positive", id)}
	}
	return service.repo.GetByID(id)
}

func (service *eventService) UpdateEvent(id int, updates map[string]interface{}) (*domain.Event, error) {

	if len(updates) == 0 {
		return nil, &domain.ValidationError{Msg: "no fields to update"}
	}

	if seats, ok := updates["seats"]; ok {
		if seatsPtr, ok := seats.(int); ok {
			if seatsPtr < 0 {
				return nil, &domain.ValidationError{Msg: "seats cannot be negative"}
			}
			// Validate that reducing seats doesn't break packet constraints
			if err := service.validateSeatsAgainstPackets(id, seatsPtr); err != nil {
				return nil, err
			}
		}
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

// validateSeatsAgainstPackets validates that reducing event seats doesn't violate packet constraints
func (service *eventService) validateSeatsAgainstPackets(eventID int, newSeats int) error {
	// Get all packets this event is included in
	packets, err := service.inclusionRepo.GetEventPacketsByEventID(eventID)
	if err != nil {
		return err
	}

	// Check each packet's allocated seats constraint
	for _, packet := range packets {
		if packet.AllocatedSeats != nil && newSeats < *packet.AllocatedSeats {
			return &domain.ValidationError{
				Msg: fmt.Sprintf("cannot reduce seats to %d: event is in packet '%s' which requires %d allocated seats", newSeats, packet.Name, *packet.AllocatedSeats),
			}
		}
	}

	return nil
}

func (service *eventService) DeleteEvent(id int) (*domain.Event, error) {
	if id < 1 {
		return nil, &domain.ValidationError{Msg: fmt.Sprintf("id:%d must be positive", id)}
	}
	return service.repo.Delete(id)
}

func (service *eventService) FilterEvents(filter *domain.EventFilter) ([]*domain.Event, error) {
	return service.repo.FilterEvents(filter)
}
