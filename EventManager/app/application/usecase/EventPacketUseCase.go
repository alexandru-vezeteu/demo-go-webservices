package usecase

import (
	"context"
	"eventManager/application/domain"
	"eventManager/application/repository"
	"eventManager/application/service"
	"fmt"
)

type EventPacketUseCase interface {
	CreateEventPacket(ctx context.Context, event *domain.EventPacket) (*domain.EventPacket, error)
	GetEventPacketByID(ctx context.Context, id int) (*domain.EventPacket, error)
	UpdateEventPacket(ctx context.Context, id int, updates map[string]interface{}) (*domain.EventPacket, error)
	DeleteEventPacket(ctx context.Context, id int) (*domain.EventPacket, error)
}

type eventPacketUseCase struct {
	repo                 repository.EventPacketRepository
	eventPacketService   service.EventPacketService 
}

func NewEventPacketUseCase(repo repository.EventPacketRepository, eventPacketService service.EventPacketService) *eventPacketUseCase {
	return &eventPacketUseCase{
		repo:               repo,
		eventPacketService: eventPacketService,
	}
}

func (uc *eventPacketUseCase) validateEventPacket(event *domain.EventPacket) error {
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

func (uc *eventPacketUseCase) CreateEventPacket(ctx context.Context, event *domain.EventPacket) (*domain.EventPacket, error) {
	if err := uc.validateEventPacket(event); err != nil {
		return nil, err
	}
	return uc.repo.Create(ctx, event)
}

func (uc *eventPacketUseCase) GetEventPacketByID(ctx context.Context, id int) (*domain.EventPacket, error) {
	if id < 1 {
		return nil, &domain.ValidationError{Reason: fmt.Sprintf("id:%d must be positive", id)}
	}
	return uc.repo.GetByID(ctx, id)
}

func (uc *eventPacketUseCase) UpdateEventPacket(ctx context.Context, id int, updates map[string]interface{}) (*domain.EventPacket, error) {
	
	
	return uc.eventPacketService.UpdateEventPacket(ctx, id, updates)
}

func (uc *eventPacketUseCase) DeleteEventPacket(ctx context.Context, id int) (*domain.EventPacket, error) {
	if id < 1 {
		return nil, &domain.ValidationError{Reason: fmt.Sprintf("id:%d must be positive", id)}
	}
	return uc.repo.Delete(ctx, id)
}
