package service

import (
	"eventManager/application/domain"
)

type IEventService interface {
	CreateEvent(event *domain.Event) (*domain.Event, error)
	GetEventByID(id int) (*domain.Event, error)
	UpdateEvent(id int, updates map[string]interface{}) (*domain.Event, error)
	DeleteEvent(id int) (*domain.Event, error)
	FilterEvents(filter *domain.EventFilter) ([]*domain.Event, error)
}
