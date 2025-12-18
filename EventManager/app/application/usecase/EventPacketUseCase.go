package usecase

import (
	"context"
	"eventManager/application/domain"
	"eventManager/application/repository"
	"eventManager/application/service"
	"fmt"
)

type EventPacketUseCase interface {
	CreateEventPacket(ctx context.Context, token string, event *domain.EventPacket) (*domain.EventPacket, error)
	GetEventPacketByID(ctx context.Context, token string, id int) (*domain.EventPacket, error)
	UpdateEventPacket(ctx context.Context, token string, id int, updates map[string]interface{}) (*domain.EventPacket, error)
	DeleteEventPacket(ctx context.Context, token string, id int) (*domain.EventPacket, error)
}

type eventPacketUseCase struct {
	repo               repository.EventPacketRepository
	eventPacketService service.EventPacketService
	authNService       service.AuthenticationService
	authZService       service.AuthorizationService
}

func NewEventPacketUseCase(
	repo repository.EventPacketRepository,
	eventPacketService service.EventPacketService,
	authNService service.AuthenticationService,
	authZService service.AuthorizationService,
) *eventPacketUseCase {
	return &eventPacketUseCase{
		repo:               repo,
		eventPacketService: eventPacketService,
		authNService:       authNService,
		authZService:       authZService,
	}
}

func (uc *eventPacketUseCase) authenticate(ctx context.Context, token string) (*service.UserIdentity, error) {
	identity, err := uc.authNService.WhoIsUser(ctx, token)
	if err != nil {
		return nil, &domain.ValidationError{Reason: "invalid or expired token"}
	}
	return identity, nil
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

func (uc *eventPacketUseCase) CreateEventPacket(ctx context.Context, token string, event *domain.EventPacket) (*domain.EventPacket, error) {
	identity, err := uc.authenticate(ctx, token)
	if err != nil {
		return nil, err
	}

	allowed, err := uc.authZService.CanUserCreateEventPacket(ctx, identity.UserID)
	if err != nil {
		return nil, &domain.ValidationError{Reason: fmt.Sprintf("authorization check failed: %v", err)}
	}
	if !allowed {
		return nil, &domain.ValidationError{Reason: "user not authorized to create event packets"}
	}

	if err := uc.validateEventPacket(event); err != nil {
		return nil, err
	}

	return uc.repo.Create(ctx, event)
}

func (uc *eventPacketUseCase) GetEventPacketByID(ctx context.Context, token string, id int) (*domain.EventPacket, error) {
	if id < 1 {
		return nil, &domain.ValidationError{Reason: fmt.Sprintf("id:%d must be positive", id)}
	}

	identity, err := uc.authenticate(ctx, token)
	if err != nil {
		return nil, err
	}

	packet, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	allowed, err := uc.authZService.CanUserViewEventPacket(ctx, identity.UserID, packet)
	if err != nil {
		return nil, &domain.ValidationError{Reason: fmt.Sprintf("authorization check failed: %v", err)}
	}
	if !allowed {
		return nil, &domain.ValidationError{Reason: "user not authorized to view event packet"}
	}

	return packet, nil
}

func (uc *eventPacketUseCase) UpdateEventPacket(ctx context.Context, token string, id int, updates map[string]interface{}) (*domain.EventPacket, error) {
	if id < 1 {
		return nil, &domain.ValidationError{Reason: fmt.Sprintf("id:%d must be positive", id)}
	}

	identity, err := uc.authenticate(ctx, token)
	if err != nil {
		return nil, err
	}

	packet, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	allowed, err := uc.authZService.CanUserEditEventPacket(ctx, identity.UserID, packet)
	if err != nil {
		return nil, &domain.ValidationError{Reason: fmt.Sprintf("authorization check failed: %v", err)}
	}
	if !allowed {
		return nil, &domain.ValidationError{Reason: "user not authorized to edit event packet"}
	}

	return uc.eventPacketService.UpdateEventPacket(ctx, id, updates)
}

func (uc *eventPacketUseCase) DeleteEventPacket(ctx context.Context, token string, id int) (*domain.EventPacket, error) {
	if id < 1 {
		return nil, &domain.ValidationError{Reason: fmt.Sprintf("id:%d must be positive", id)}
	}

	identity, err := uc.authenticate(ctx, token)
	if err != nil {
		return nil, err
	}

	packet, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	allowed, err := uc.authZService.CanUserDeleteEventPacket(ctx, identity.UserID, packet)
	if err != nil {
		return nil, &domain.ValidationError{Reason: fmt.Sprintf("authorization check failed: %v", err)}
	}
	if !allowed {
		return nil, &domain.ValidationError{Reason: "user not authorized to delete event packet"}
	}

	return uc.repo.Delete(ctx, id)
}
