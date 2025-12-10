package usecase

import (
	"eventManager/application/domain"
	"eventManager/application/repository"
	"eventManager/application/service"
)

type EventPacketInclusionUseCase interface {
	CreateEventPacketInclusion(inclusion *domain.EventPacketInclusion) (*domain.EventPacketInclusion, error)
	GetEventsByPacketID(packetID int) ([]*domain.Event, error)
	GetEventPacketsByEventID(eventID int) ([]*domain.EventPacket, error)
	Update(eventID, packetID int, updates map[string]interface{}) (*domain.EventPacketInclusion, error)
	DeleteEventPacketInclusion(eventID, packetID int) (*domain.EventPacketInclusion, error)
}

type eventPacketInclusionUseCase struct {
	repo                         repository.EventPacketInclusionRepository
	eventPacketInclusionService  service.EventPacketInclusionService // For complex constraint validation
}

func NewEventPacketInclusionUseCase(repo repository.EventPacketInclusionRepository, eventPacketInclusionService service.EventPacketInclusionService) *eventPacketInclusionUseCase {
	return &eventPacketInclusionUseCase{
		repo:                        repo,
		eventPacketInclusionService: eventPacketInclusionService,
	}
}

func (uc *eventPacketInclusionUseCase) CreateEventPacketInclusion(inclusion *domain.EventPacketInclusion) (*domain.EventPacketInclusion, error) {
	// CreateEventPacketInclusion has complex constraint validation
	// (validates event and packet exist, checks seat constraints)
	// Delegate to service
	return uc.eventPacketInclusionService.CreateEventPacketInclusion(inclusion)
}

func (uc *eventPacketInclusionUseCase) GetEventsByPacketID(packetID int) ([]*domain.Event, error) {
	if packetID < 0 {
		return nil, &domain.ValidationError{Reason: "packet id must be >= 0"}
	}
	return uc.repo.GetEventsByPacketID(packetID)
}

func (uc *eventPacketInclusionUseCase) GetEventPacketsByEventID(eventID int) ([]*domain.EventPacket, error) {
	if eventID < 0 {
		return nil, &domain.ValidationError{Reason: "event id must be >= 0"}
	}
	return uc.repo.GetEventPacketsByEventID(eventID)
}

func (uc *eventPacketInclusionUseCase) Update(eventID, packetID int, updates map[string]interface{}) (*domain.EventPacketInclusion, error) {
	if eventID < 0 {
		return nil, &domain.ValidationError{Reason: "event id must be >= 0"}
	}
	if packetID < 0 {
		return nil, &domain.ValidationError{Reason: "packet id must be >= 0"}
	}
	if len(updates) == 0 {
		return nil, &domain.ValidationError{Reason: "update must contain at least one field"}
	}
	return uc.repo.Update(eventID, packetID, updates)
}

func (uc *eventPacketInclusionUseCase) DeleteEventPacketInclusion(eventID, packetID int) (*domain.EventPacketInclusion, error) {
	if eventID < 0 {
		return nil, &domain.ValidationError{Reason: "event id must be >= 0"}
	}
	if packetID < 0 {
		return nil, &domain.ValidationError{Reason: "packet id must be >= 0"}
	}
	return uc.repo.Delete(eventID, packetID)
}
