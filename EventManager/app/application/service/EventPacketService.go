package service

import (
	"eventManager/application/domain"
	"eventManager/application/repository"
	"fmt"
)

type EventPacketService interface {
	CreateEventPacket(event *domain.EventPacket) (*domain.EventPacket, error)
	GetEventPacketByID(id int) (*domain.EventPacket, error)
	UpdateEventPacket(id int, updates map[string]interface{}) (*domain.EventPacket, error)
	DeleteEventPacket(id int) (*domain.EventPacket, error)
}

type eventPacketService struct {
	repo          repository.EventPacketRepository
	inclusionRepo repository.EventPacketInclusionRepository
}

func NewEventPacketService(repo repository.EventPacketRepository, inclusionRepo repository.EventPacketInclusionRepository) EventPacketService {
	return &eventPacketService{
		repo:          repo,
		inclusionRepo: inclusionRepo,
	}
}

func (service *eventPacketService) CreateEventPacket(event *domain.EventPacket) (*domain.EventPacket, error) {
	if err := service.validateEventPacket(event); err != nil {
		return nil, err
	}
	// Note: At creation time, no events are included yet, so no constraint validation needed
	return service.repo.Create(event)
}

func (service *eventPacketService) GetEventPacketByID(id int) (*domain.EventPacket, error) {
	if id < 1 {
		return nil, &domain.ValidationError{Reason: fmt.Sprintf("id:%d must be positive", id)}
	}
	return service.repo.GetByID(id)
}

func (service *eventPacketService) UpdateEventPacket(id int, updates map[string]interface{}) (*domain.EventPacket, error) {

	if len(updates) == 0 {
		return nil, &domain.ValidationError{Reason: "no fields to update"}
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

	// Validate allocated_seats constraint
	if allocatedSeats, ok := updates["allocated_seats"]; ok {
		if seatsPtr, ok := allocatedSeats.(int); ok {
			if seatsPtr < 0 {
				return nil, &domain.ValidationError{Reason: "allocated_seats must be non-negative"}
			}
			// Validate against included events' seats
			if err := service.validateAllocatedSeatsConstraint(id, seatsPtr); err != nil {
				return nil, err
			}
		}
	}

	return service.repo.Update(id, updates)
}
func (service *eventPacketService) DeleteEventPacket(id int) (*domain.EventPacket, error) {
	if id < 1 {
		return nil, &domain.ValidationError{Reason: fmt.Sprintf("id:%d must be positive", id)}
	}
	return service.repo.Delete(id)
}

func (service *eventPacketService) validateEventPacket(event *domain.EventPacket) error {
	if event == nil {
		return &domain.ValidationError{Reason: "invalid object received"}
	}

	if event.OwnerID < 1 {
		return &domain.ValidationError{Reason: "owner_id must be positive"}
	}

	if event.Name == "" {
		return &domain.ValidationError{Reason: "name must be set"}
	}

	if event.AllocatedSeats != nil && *event.AllocatedSeats < 0 {
		return &domain.ValidationError{Reason: "allocated_seats must be non-negative"}
	}

	return nil
}

// findMinSeats finds the minimum seats among a list of events
// Returns nil if any event has nil seats
func findMinSeats(events []*domain.Event) *int {
	var min *int
	for _, event := range events {
		if event.Seats == nil {
			return nil
		}
		if min == nil || *event.Seats < *min {
			min = event.Seats
		}
	}
	return min
}

// validateAllocatedSeatsConstraint validates that allocated seats doesn't exceed minimum event seats
func (service *eventPacketService) validateAllocatedSeatsConstraint(packetID int, requestedSeats int) error {
	// Get all events in this packet
	events, err := service.inclusionRepo.GetEventsByPacketID(packetID)
	if err != nil {
		return err
	}

	// If no events included, any value is OK
	if len(events) == 0 {
		return nil
	}

	// Find minimum seats among included events
	minSeats := findMinSeats(events)

	// If any event has nil seats, we can't validate
	if minSeats == nil {
		return &domain.ValidationError{
			Reason: "cannot set allocated_seats: some included events don't have seats defined",
		}
	}

	// Validate constraint
	if requestedSeats > *minSeats {
		return &domain.ValidationError{
			Reason: fmt.Sprintf("allocated_seats (%d) cannot exceed minimum seats of included events (%d)", requestedSeats, *minSeats),
		}
	}

	return nil
}
