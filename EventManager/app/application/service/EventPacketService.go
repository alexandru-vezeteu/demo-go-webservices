package service

import (
	"context"
	"eventManager/application/domain"
	"eventManager/application/repository"
	"fmt"
)

type EventPacketService interface {
	CreateEventPacket(ctx context.Context, event *domain.EventPacket) (*domain.EventPacket, error)
	GetEventPacketByID(ctx context.Context, id int) (*domain.EventPacket, error)
	UpdateEventPacket(ctx context.Context, id int, updates map[string]interface{}) (*domain.EventPacket, error)
	DeleteEventPacket(ctx context.Context, id int) (*domain.EventPacket, error)
	FilterEventPackets(ctx context.Context, filter *domain.EventPacketFilter) ([]*domain.EventPacket, error)
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

func (service *eventPacketService) CreateEventPacket(ctx context.Context, event *domain.EventPacket) (*domain.EventPacket, error) {
	if err := service.validateEventPacket(event); err != nil {
		return nil, err
	}

	return service.repo.Create(ctx, event)
}

func (service *eventPacketService) GetEventPacketByID(ctx context.Context, id int) (*domain.EventPacket, error) {
	if id < 1 {
		return nil, &domain.ValidationError{Reason: fmt.Sprintf("id:%d must be positive", id)}
	}
	return service.repo.GetByID(ctx, id)
}

func (service *eventPacketService) UpdateEventPacket(ctx context.Context, id int, updates map[string]interface{}) (*domain.EventPacket, error) {

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

	if allocatedSeats, ok := updates["allocated_seats"]; ok {
		if seatsPtr, ok := allocatedSeats.(int); ok {
			if seatsPtr < 0 {
				return nil, &domain.ValidationError{Reason: "allocated_seats must be non-negative"}
			}

			if err := service.validateAllocatedSeatsConstraint(ctx, id, seatsPtr); err != nil {
				return nil, err
			}
		}
	}

	return service.repo.Update(ctx, id, updates)
}
func (service *eventPacketService) DeleteEventPacket(ctx context.Context, id int) (*domain.EventPacket, error) {
	if id < 1 {
		return nil, &domain.ValidationError{Reason: fmt.Sprintf("id:%d must be positive", id)}
	}
	return service.repo.Delete(ctx, id)
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

func (service *eventPacketService) validateAllocatedSeatsConstraint(ctx context.Context, packetID int, requestedSeats int) error {

	events, err := service.inclusionRepo.GetEventsByPacketID(ctx, packetID)
	if err != nil {
		return err
	}

	if len(events) == 0 {
		return nil
	}

	minSeats := findMinSeats(events)

	if minSeats == nil {
		return &domain.ValidationError{
			Reason: "cannot set allocated_seats: some included events don't have seats defined",
		}
	}

	if requestedSeats > *minSeats {
		return &domain.ValidationError{
			Reason: fmt.Sprintf("allocated_seats (%d) cannot exceed minimum seats of included events (%d)", requestedSeats, *minSeats),
		}
	}

	return nil
}

func (service *eventPacketService) FilterEventPackets(ctx context.Context, filter *domain.EventPacketFilter) ([]*domain.EventPacket, error) {
	return service.repo.FilterEventPackets(ctx, filter)
}
