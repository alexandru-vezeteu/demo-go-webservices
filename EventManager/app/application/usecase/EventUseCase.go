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
	FilterEvents(ctx context.Context, token string, filter *domain.EventFilter) ([]*domain.Event, int, error)
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

func (uc *eventUseCase) CreateEvent(ctx context.Context, token string, event *domain.Event) (*domain.Event, error) {
	identity, err := uc.authenticate(ctx, token)
	if err != nil {
		return nil, err
	}

	allowed, err := uc.authZService.CanUserCreateEvent(ctx, *identity)
	if err != nil {
		return nil, &domain.InternalError{Msg: fmt.Sprintf("authorization check failed: %v", err)}
	}
	if !allowed {
		return nil, &domain.ForbiddenError{Reason: "you don't have permission to create events"}
	}

	return uc.eventService.CreateEvent(ctx, event)
}

func (uc *eventUseCase) GetEventByID(ctx context.Context, token string, id int) (*domain.Event, error) {
	return uc.eventService.GetEventByID(ctx, id)
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

	allowed, err := uc.authZService.CanUserEditEvent(ctx, *identity, event)
	if err != nil {
		return nil, &domain.InternalError{Msg: fmt.Sprintf("authorization check failed: %v", err)}
	}
	if !allowed {
		return nil, &domain.ForbiddenError{Reason: "you don't have permission to edit this event"}
	}

	return uc.eventService.UpdateEvent(ctx, id, updates)
}

func (uc *eventUseCase) DeleteEvent(ctx context.Context, token string, id int) (*domain.Event, error) {
	identity, err := uc.authenticate(ctx, token)
	if err != nil {
		return nil, err
	}

	event, err := uc.eventService.GetEventByID(ctx, id)
	if err != nil {
		return nil, err
	}

	allowed, err := uc.authZService.CanUserDeleteEvent(ctx, *identity, event)
	if err != nil {
		return nil, &domain.InternalError{Msg: fmt.Sprintf("authorization check failed: %v", err)}
	}
	if !allowed {
		return nil, &domain.ForbiddenError{Reason: "you don't have permission to delete this event"}
	}

	return uc.eventService.DeleteEvent(ctx, id)
}

func (uc *eventUseCase) FilterEvents(ctx context.Context, token string, filter *domain.EventFilter) ([]*domain.Event, int, error) {
	events, err := uc.eventService.FilterEvents(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	totalCount, err := uc.repo.CountEvents(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return events, totalCount, nil
}
