package usecase

import (
	"context"
	"eventManager/application/domain"
	"eventManager/application/repository"
	"eventManager/application/service"
	"fmt"
)

type EventUseCase interface {
	CreateEvent(ctx context.Context, token string, event *domain.Event) (*domain.Event, error)
	GetEventByID(ctx context.Context, token string, id int) (*domain.Event, error)
	UpdateEvent(ctx context.Context, token string, id int, updates map[string]interface{}) (*domain.Event, error)
	DeleteEvent(ctx context.Context, token string, id int) (*domain.Event, error)
	FilterEvents(ctx context.Context, token string, filter *domain.EventFilter) ([]*domain.Event, error)
}

type eventUseCase struct {
	repo         repository.EventRepository
	eventService service.EventService
	authNService service.AuthenticationService
	authZService service.AuthorizationService
}

func NewEventUseCase(
	repo repository.EventRepository,
	eventService service.EventService,
	authNService service.AuthenticationService,
	authZService service.AuthorizationService,
) *eventUseCase {
	return &eventUseCase{
		repo:         repo,
		eventService: eventService,
		authNService: authNService,
		authZService: authZService,
	}
}

func (uc *eventUseCase) authenticate(ctx context.Context, token string) (*service.UserIdentity, error) {
	identity, err := uc.authNService.WhoIsUser(ctx, token)
	if err != nil {
		return nil, &domain.ValidationError{Reason: "invalid or expired token"}
	}
	return identity, nil
}

func (uc *eventUseCase) validateEvent(event *domain.Event) error {
	if event == nil {
		return &domain.ValidationError{Reason: "invalid object received"}
	}

	if event.OwnerID < 1 {
		return &domain.ValidationError{Reason: "owner_id must be positive"}
	}

	if event.Name == "" {
		return &domain.ValidationError{Reason: "name must be set"}
	}
	return nil
}

func (uc *eventUseCase) CreateEvent(ctx context.Context, token string, event *domain.Event) (*domain.Event, error) {
	identity, err := uc.authenticate(ctx, token)
	if err != nil {
		return nil, err
	}

	allowed, err := uc.authZService.CanUserCreateEvent(ctx, identity.UserID)
	if err != nil {
		return nil, &domain.ValidationError{Reason: fmt.Sprintf("authorization check failed: %v", err)}
	}
	if !allowed {
		return nil, &domain.ValidationError{Reason: "user not authorized to create events"}
	}

	if err := uc.validateEvent(event); err != nil {
		return nil, err
	}

	return uc.repo.Create(ctx, event)
}

func (uc *eventUseCase) GetEventByID(ctx context.Context, token string, id int) (*domain.Event, error) {
	if id < 1 {
		return nil, &domain.ValidationError{Reason: fmt.Sprintf("id:%d must be positive", id)}
	}

	identity, err := uc.authenticate(ctx, token)
	if err != nil {
		return nil, err
	}

	event, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	allowed, err := uc.authZService.CanUserViewEvent(ctx, identity.UserID, event)
	if err != nil {
		return nil, &domain.ValidationError{Reason: fmt.Sprintf("authorization check failed: %v", err)}
	}
	if !allowed {
		return nil, &domain.ValidationError{Reason: "user not authorized to view event"}
	}

	return event, nil
}

func (uc *eventUseCase) UpdateEvent(ctx context.Context, token string, id int, updates map[string]interface{}) (*domain.Event, error) {
	identity, err := uc.authenticate(ctx, token)
	if err != nil {
		return nil, err
	}

	event, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	allowed, err := uc.authZService.CanUserEditEvent(ctx, identity.UserID, event)
	if err != nil {
		return nil, &domain.ValidationError{Reason: fmt.Sprintf("authorization check failed: %v", err)}
	}
	if !allowed {
		return nil, &domain.ValidationError{Reason: "user not authorized to edit event"}
	}

	return uc.eventService.UpdateEvent(ctx, id, updates)
}

func (uc *eventUseCase) DeleteEvent(ctx context.Context, token string, id int) (*domain.Event, error) {
	if id < 1 {
		return nil, &domain.ValidationError{Reason: fmt.Sprintf("id:%d must be positive", id)}
	}

	identity, err := uc.authenticate(ctx, token)
	if err != nil {
		return nil, err
	}

	event, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	allowed, err := uc.authZService.CanUserDeleteEvent(ctx, identity.UserID, event)
	if err != nil {
		return nil, &domain.ValidationError{Reason: fmt.Sprintf("authorization check failed: %v", err)}
	}
	if !allowed {
		return nil, &domain.ValidationError{Reason: "user not authorized to delete event"}
	}

	return uc.repo.Delete(ctx, id)
}

func (uc *eventUseCase) FilterEvents(ctx context.Context, token string, filter *domain.EventFilter) ([]*domain.Event, error) {
	identity, err := uc.authenticate(ctx, token)
	if err != nil {
		return nil, err
	}

	events, err := uc.repo.FilterEvents(ctx, filter)
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
