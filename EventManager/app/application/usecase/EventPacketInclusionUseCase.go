package usecase

import (
	"context"
	"eventManager/application/domain"
	"eventManager/application/repository"
	"fmt"
	"strings"
)

type EventPacketInclusionUseCase interface {
	CreateEventPacketInclusion(ctx context.Context, inclusion *domain.EventPacketInclusion) (*domain.EventPacketInclusion, error)
	GetEventsByPacketID(ctx context.Context, packetID int) ([]*domain.Event, error)
	GetEventPacketsByEventID(ctx context.Context, eventID int) ([]*domain.EventPacket, error)
	Update(ctx context.Context, eventID, packetID int, updates map[string]interface{}) (*domain.EventPacketInclusion, error)
	DeleteEventPacketInclusion(ctx context.Context, eventID, packetID int) (*domain.EventPacketInclusion, error)
}

type eventPacketInclusionUseCase struct {
	repo       repository.EventPacketInclusionRepository
	eventRepo  repository.EventRepository
	packetRepo repository.EventPacketRepository
}

func NewEventPacketInclusionUseCase(
	repo repository.EventPacketInclusionRepository,
	eventRepo repository.EventRepository,
	packetRepo repository.EventPacketRepository,
) *eventPacketInclusionUseCase {
	return &eventPacketInclusionUseCase{
		repo:       repo,
		eventRepo:  eventRepo,
		packetRepo: packetRepo,
	}
}

func (uc *eventPacketInclusionUseCase) CreateEventPacketInclusion(ctx context.Context, inclusion *domain.EventPacketInclusion) (*domain.EventPacketInclusion, error) {
	if inclusion == nil {
		return nil, &domain.ValidationError{Reason: "invalid object"}
	}
	if inclusion.EventID < 0 {
		return nil, &domain.ValidationError{Reason: "event id must be >= 0"}
	}
	if inclusion.PacketID < 0 {
		return nil, &domain.ValidationError{Reason: "packet id must be >= 0"}
	}

	
	if err := uc.validateInclusionConstraints(ctx, inclusion.EventID, inclusion.PacketID); err != nil {
		return nil, err
	}

	return uc.repo.Create(ctx, inclusion)
}


func (uc *eventPacketInclusionUseCase) validateInclusionConstraints(ctx context.Context, eventID int, packetID int) error {
	
	event, err := uc.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return &domain.NotFoundError{ID: eventID}
		}
		return err
	}

	
	packet, err := uc.packetRepo.GetByID(ctx, packetID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return &domain.NotFoundError{ID: packetID}
		}
		return err
	}

	
	if packet.AllocatedSeats == nil {
		return nil
	}

	
	if event.Seats == nil {
		return &domain.ValidationError{
			Reason: fmt.Sprintf("event %d doesn't have seats defined, cannot be added to packet requiring %d seats", event.ID, *packet.AllocatedSeats),
		}
	}

	
	if *event.Seats < *packet.AllocatedSeats {
		return &domain.ValidationError{
			Reason: fmt.Sprintf("event has %d seats but packet requires %d allocated seats", *event.Seats, *packet.AllocatedSeats),
		}
	}

	return nil
}

func (uc *eventPacketInclusionUseCase) GetEventsByPacketID(ctx context.Context, packetID int) ([]*domain.Event, error) {
	if packetID < 0 {
		return nil, &domain.ValidationError{Reason: "packet id must be >= 0"}
	}
	return uc.repo.GetEventsByPacketID(ctx, packetID)
}

func (uc *eventPacketInclusionUseCase) GetEventPacketsByEventID(ctx context.Context, eventID int) ([]*domain.EventPacket, error) {
	if eventID < 0 {
		return nil, &domain.ValidationError{Reason: "event id must be >= 0"}
	}
	return uc.repo.GetEventPacketsByEventID(ctx, eventID)
}

func (uc *eventPacketInclusionUseCase) Update(ctx context.Context, eventID, packetID int, updates map[string]interface{}) (*domain.EventPacketInclusion, error) {
	if eventID < 0 {
		return nil, &domain.ValidationError{Reason: "event id must be >= 0"}
	}
	if packetID < 0 {
		return nil, &domain.ValidationError{Reason: "packet id must be >= 0"}
	}
	if len(updates) == 0 {
		return nil, &domain.ValidationError{Reason: "update must contain at least one field"}
	}
	return uc.repo.Update(ctx, eventID, packetID, updates)
}

func (uc *eventPacketInclusionUseCase) DeleteEventPacketInclusion(ctx context.Context, eventID, packetID int) (*domain.EventPacketInclusion, error) {
	if eventID < 0 {
		return nil, &domain.ValidationError{Reason: "event id must be >= 0"}
	}
	if packetID < 0 {
		return nil, &domain.ValidationError{Reason: "packet id must be >= 0"}
	}
	return uc.repo.Delete(ctx, eventID, packetID)
}
