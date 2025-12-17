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
	userInfo, err := uc.authNService.WhoIsUser(ctx, token)
	if err != nil {
		return nil, &domain.ValidationError{Reason: "authentication failed"}
	}

	canCreate, err := uc.authZService.CanUserCreateEvent(ctx, userInfo.UserID)
	if err != nil || !canCreate {
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

	userInfo, err := uc.authNService.WhoIsUser(ctx, token)
	if err != nil {
		return nil, &domain.ValidationError{Reason: "authentication failed"}
	}

	canSee, err := uc.authZService.CanUserSeeEvent(ctx, userInfo.UserID, fmt.Sprintf("%d", id))
	if err != nil || !canSee {
		return nil, &domain.ValidationError{Reason: "user not authorized to view this event"}
	}

	return uc.repo.GetByID(ctx, id)
}

func (uc *eventUseCase) UpdateEvent(ctx context.Context, token string, id int, updates map[string]interface{}) (*domain.Event, error) {
	userInfo, err := uc.authNService.WhoIsUser(ctx, token)
	if err != nil {
		return nil, &domain.ValidationError{Reason: "authentication failed"}
	}

	canUpdate, err := uc.authZService.CanUserUpdateEvent(ctx, userInfo.UserID, fmt.Sprintf("%d", id))
	if err != nil || !canUpdate {
		return nil, &domain.ValidationError{Reason: "user not authorized to update this event"}
	}

	return uc.eventService.UpdateEvent(ctx, id, updates)
}

func (uc *eventUseCase) DeleteEvent(ctx context.Context, token string, id int) (*domain.Event, error) {
	if id < 1 {
		return nil, &domain.ValidationError{Reason: fmt.Sprintf("id:%d must be positive", id)}
	}

	userInfo, err := uc.authNService.WhoIsUser(ctx, token)
	if err != nil {
		return nil, &domain.ValidationError{Reason: "authentication failed"}
	}

	canDelete, err := uc.authZService.CanUserDeleteEvent(ctx, userInfo.UserID, fmt.Sprintf("%d", id))
	if err != nil || !canDelete {
		return nil, &domain.ValidationError{Reason: "user not authorized to delete this event"}
	}

	return uc.repo.Delete(ctx, id)
}

func (uc *eventUseCase) FilterEvents(ctx context.Context, token string, filter *domain.EventFilter) ([]*domain.Event, error) {
	_, err := uc.authNService.WhoIsUser(ctx, token)
	if err != nil {
		return nil, &domain.ValidationError{Reason: "authentication failed"}
	}

	return uc.repo.FilterEvents(ctx, filter)
}
