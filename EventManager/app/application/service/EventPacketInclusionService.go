package service

import (
	"eventManager/application/domain"
	"eventManager/application/repository"
	"fmt"
	"strings"
)

type EventPacketInclusionService interface {
	CreateEventPacketInclusion(inclusion *domain.EventPacketInclusion) (*domain.EventPacketInclusion, error)
	GetEventsByPacketID(packetID int) ([]*domain.Event, error)
	GetEventPacketsByEventID(eventID int) ([]*domain.EventPacket, error)
	Update(eventID, packetID int, updates map[string]interface{}) (*domain.EventPacketInclusion, error)
	DeleteEventPacketInclusion(eventID, packetID int) (*domain.EventPacketInclusion, error)
}

type eventPacketInclusionService struct {
	repo       repository.EventPacketInclusionRepository
	eventRepo  repository.EventRepository
	packetRepo repository.EventPacketRepository
}

func NewEventPacketInclusionService(
	repo repository.EventPacketInclusionRepository,
	eventRepo repository.EventRepository,
	packetRepo repository.EventPacketRepository,
) EventPacketInclusionService {
	return &eventPacketInclusionService{
		repo:       repo,
		eventRepo:  eventRepo,
		packetRepo: packetRepo,
	}
}

func (service *eventPacketInclusionService) CreateEventPacketInclusion(inclusion *domain.EventPacketInclusion) (*domain.EventPacketInclusion, error) {

	if inclusion == nil {
		return nil, &domain.ValidationError{Reason: "invalid object"}
	}
	if inclusion.EventID < 0 {
		return nil, &domain.ValidationError{Reason: "event id must be >= 0"}
	}
	if inclusion.PacketID < 0 {
		return nil, &domain.ValidationError{Reason: "packet id must be >= 0"}
	}

	// Validate seat constraints
	if err := service.validateInclusionConstraints(inclusion.EventID, inclusion.PacketID); err != nil {
		return nil, err
	}

	return service.repo.Create(inclusion)

}

// validateInclusionConstraints validates that an event can be added to a packet
func (service *eventPacketInclusionService) validateInclusionConstraints(eventID int, packetID int) error {
	// Get the event
	event, err := service.eventRepo.GetByID(eventID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return &domain.NotFoundError{ID: eventID}
		}
		return err
	}

	// Get the packet
	packet, err := service.packetRepo.GetByID(packetID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return &domain.NotFoundError{ID: packetID}
		}
		return err
	}

	// If packet doesn't specify allocated seats, any event is OK
	if packet.AllocatedSeats == nil {
		return nil
	}

	// Event must have seats defined
	if event.Seats == nil {
		return &domain.ValidationError{
			Reason: fmt.Sprintf("event %d doesn't have seats defined, cannot be added to packet requiring %d seats", event.ID, *packet.AllocatedSeats),
		}
	}

	// Event must have enough seats
	if *event.Seats < *packet.AllocatedSeats {
		return &domain.ValidationError{
			Reason: fmt.Sprintf("event has %d seats but packet requires %d allocated seats", *event.Seats, *packet.AllocatedSeats),
		}
	}

	return nil
}
func (service *eventPacketInclusionService) GetEventsByPacketID(packetID int) ([]*domain.Event, error) {
	if packetID < 0 {
		return nil, &domain.ValidationError{Reason: "packet id must be >= 0"}
	}
	return service.repo.GetEventsByPacketID(packetID)
}
func (service *eventPacketInclusionService) GetEventPacketsByEventID(eventID int) ([]*domain.EventPacket, error) {
	if eventID < 0 {
		return nil, &domain.ValidationError{Reason: "event id must be >= 0"}
	}
	return service.repo.GetEventPacketsByEventID(eventID)
}
func (service *eventPacketInclusionService) Update(eventID, packetID int, updates map[string]interface{}) (*domain.EventPacketInclusion, error) {
	if eventID < 0 {
		return nil, &domain.ValidationError{Reason: "event id must be >= 0"}
	}
	if packetID < 0 {
		return nil, &domain.ValidationError{Reason: "packet id must be >= 0"}
	}
	if len(updates) == 0 {
		return nil, &domain.ValidationError{Reason: "update must contain at least one field"}
	}
	return service.repo.Update(eventID, packetID, updates)
}
func (service *eventPacketInclusionService) DeleteEventPacketInclusion(eventID, packetID int) (*domain.EventPacketInclusion, error) {
	if eventID < 0 {
		return nil, &domain.ValidationError{Reason: "event id must be >= 0"}
	}
	if packetID < 0 {
		return nil, &domain.ValidationError{Reason: "packet id must be >= 0"}
	}

	return service.repo.Delete(eventID, packetID)
}
