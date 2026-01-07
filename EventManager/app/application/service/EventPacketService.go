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
	eventRepo     repository.EventRepository
	inclusionRepo repository.EventPacketInclusionRepository
}

func NewEventPacketService(repo repository.EventPacketRepository, eventRepo repository.EventRepository, inclusionRepo repository.EventPacketInclusionRepository) EventPacketService {
	return &eventPacketService{
		repo:          repo,
		eventRepo:     eventRepo,
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

			soldTickets, err := service.repo.CountSoldTickets(ctx, id)
			if err != nil {
				return nil, &domain.InternalError{Msg: "failed to count sold tickets", Err: err}
			}

			if seatsPtr < soldTickets {
				return nil, &domain.ValidationError{
					Reason: fmt.Sprintf("cannot reduce allocated_seats to %d: %d tickets have already been sold for this packet", seatsPtr, soldTickets),
				}
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

	for _, event := range events {
		if event.Seats == nil {
			continue
		}

		directSold, err := service.eventRepo.CountSoldTickets(ctx, event.ID)
		if err != nil {
			return &domain.InternalError{Msg: "failed to count sold tickets", Err: err}
		}

		allPackets, err := service.inclusionRepo.GetEventPacketsByEventID(ctx, event.ID)
		if err != nil {
			return err
		}

		otherPacketsAllocated := 0
		for _, packet := range allPackets {
			if packet.ID == packetID {
				continue
			}
			if packet.AllocatedSeats != nil {
				otherPacketsAllocated += *packet.AllocatedSeats
			}
		}

		availableSeats := *event.Seats - directSold - otherPacketsAllocated

		if requestedSeats > availableSeats {
			return &domain.ValidationError{
				Reason: fmt.Sprintf(
					"cannot allocate %d seats: event '%s' only has %d available seats (%d total - %d sold - %d in other packets)",
					requestedSeats, event.Name, availableSeats, *event.Seats, directSold, otherPacketsAllocated,
				),
			}
		}
	}

	return nil
}

func (service *eventPacketService) FilterEventPackets(ctx context.Context, filter *domain.EventPacketFilter) ([]*domain.EventPacket, error) {
	return service.repo.FilterEventPackets(ctx, filter)
}
