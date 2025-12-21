package usecase

import (
	"context"
	"eventManager/application/domain"
	"eventManager/application/repository"
	"eventManager/application/service"
	"fmt"
	"strings"
)

type EventPacketInclusionUseCase interface {
	CreateEventPacketInclusion(ctx context.Context, token string, inclusion *domain.EventPacketInclusion) (*domain.EventPacketInclusion, error)
	GetEventsByPacketID(ctx context.Context, token string, packetID int) ([]*domain.Event, error)
	GetEventPacketsByEventID(ctx context.Context, token string, eventID int) ([]*domain.EventPacket, error)
	Update(ctx context.Context, token string, eventID, packetID int, updates map[string]interface{}) (*domain.EventPacketInclusion, error)
	DeleteEventPacketInclusion(ctx context.Context, token string, eventID, packetID int) (*domain.EventPacketInclusion, error)
}

type eventPacketInclusionUseCase struct {
	repo         repository.EventPacketInclusionRepository
	eventRepo    repository.EventRepository
	packetRepo   repository.EventPacketRepository
	authNService service.AuthenticationService
	authZService service.AuthorizationService
}

func NewEventPacketInclusionUseCase(
	repo repository.EventPacketInclusionRepository,
	eventRepo repository.EventRepository,
	packetRepo repository.EventPacketRepository,
	authNService service.AuthenticationService,
	authZService service.AuthorizationService,
) *eventPacketInclusionUseCase {
	return &eventPacketInclusionUseCase{
		repo:         repo,
		eventRepo:    eventRepo,
		packetRepo:   packetRepo,
		authNService: authNService,
		authZService: authZService,
	}
}

func (uc *eventPacketInclusionUseCase) authenticate(ctx context.Context, token string) (*service.UserIdentity, error) {
	identity, err := uc.authNService.WhoIsUser(ctx, token)
	if err != nil {
		return nil, &domain.ValidationError{Reason: "invalid or expired token"}
	}
	return identity, nil
}

func (uc *eventPacketInclusionUseCase) CreateEventPacketInclusion(ctx context.Context, token string, inclusion *domain.EventPacketInclusion) (*domain.EventPacketInclusion, error) {
	identity, err := uc.authenticate(ctx, token)
	if err != nil {
		return nil, err
	}
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

	allowed, err := uc.authZService.CanUserCreateEventPacketInclusion(ctx, identity.UserID, inclusion.EventID, inclusion.PacketID)
	if err != nil {
		return nil, &domain.ValidationError{Reason: fmt.Sprintf("authorization check failed: %v", err)}
	}
	if !allowed {
		return nil, &domain.ValidationError{Reason: "user not authorized to add packets to this event"}
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

func (uc *eventPacketInclusionUseCase) GetEventsByPacketID(ctx context.Context, token string, packetID int) ([]*domain.Event, error) {
	if packetID < 0 {
		return nil, &domain.ValidationError{Reason: "packet id must be >= 0"}
	}

	identity, err := uc.authenticate(ctx, token)
	if err != nil {
		return nil, err
	}

	events, err := uc.repo.GetEventsByPacketID(ctx, packetID)
	if err != nil {
		return nil, err
	}

	permissions, err := uc.authZService.CanUserViewEvents(ctx, identity.UserID, events)
	if err != nil {
		return nil, &domain.InternalError{Msg: fmt.Sprintf("authorization check failed: %v", err)}
	}

	authorizedEvents := make([]*domain.Event, 0, len(events))
	for i, event := range events {
		if permissions[i] {
			authorizedEvents = append(authorizedEvents, event)
		}
	}

	return authorizedEvents, nil
}

func (uc *eventPacketInclusionUseCase) GetEventPacketsByEventID(ctx context.Context, token string, eventID int) ([]*domain.EventPacket, error) {
	if eventID < 0 {
		return nil, &domain.ValidationError{Reason: "event id must be >= 0"}
	}

	identity, err := uc.authenticate(ctx, token)
	if err != nil {
		return nil, err
	}

	packets, err := uc.repo.GetEventPacketsByEventID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	permissions, err := uc.authZService.CanUserViewEventPackets(ctx, identity.UserID, packets)
	if err != nil {
		return nil, &domain.InternalError{Msg: fmt.Sprintf("authorization check failed: %v", err)}
	}

	authorizedPackets := make([]*domain.EventPacket, 0, len(packets))
	for i, packet := range packets {
		if permissions[i] {
			authorizedPackets = append(authorizedPackets, packet)
		}
	}

	return authorizedPackets, nil
}

func (uc *eventPacketInclusionUseCase) Update(ctx context.Context, token string, eventID, packetID int, updates map[string]interface{}) (*domain.EventPacketInclusion, error) {
	if eventID < 0 {
		return nil, &domain.ValidationError{Reason: "event id must be >= 0"}
	}
	if packetID < 0 {
		return nil, &domain.ValidationError{Reason: "packet id must be >= 0"}
	}
	if len(updates) == 0 {
		return nil, &domain.ValidationError{Reason: "update must contain at least one field"}
	}

	identity, err := uc.authenticate(ctx, token)
	if err != nil {
		return nil, err
	}

	allowed, err := uc.authZService.CanUserUpdateEventPacketInclusion(ctx, identity.UserID, eventID, packetID)
	if err != nil {
		return nil, &domain.ValidationError{Reason: fmt.Sprintf("authorization check failed: %v", err)}
	}
	if !allowed {
		return nil, &domain.ValidationError{Reason: "user not authorized to update this inclusion"}
	}

	return uc.repo.Update(ctx, eventID, packetID, updates)
}

func (uc *eventPacketInclusionUseCase) DeleteEventPacketInclusion(ctx context.Context, token string, eventID, packetID int) (*domain.EventPacketInclusion, error) {
	if eventID < 0 {
		return nil, &domain.ValidationError{Reason: "event id must be >= 0"}
	}
	if packetID < 0 {
		return nil, &domain.ValidationError{Reason: "packet id must be >= 0"}
	}

	identity, err := uc.authenticate(ctx, token)
	if err != nil {
		return nil, err
	}

	allowed, err := uc.authZService.CanUserDeleteEventPacketInclusion(ctx, identity.UserID, eventID, packetID)
	if err != nil {
		return nil, &domain.ValidationError{Reason: fmt.Sprintf("authorization check failed: %v", err)}
	}
	if !allowed {
		return nil, &domain.ValidationError{Reason: "user not authorized to delete this inclusion"}
	}

	return uc.repo.Delete(ctx, eventID, packetID)
}
