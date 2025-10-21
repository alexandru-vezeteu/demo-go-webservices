package service

import (
	"eventManager/domain"
)

type IEventService interface {
	CreateEvent(event *domain.Event) (*domain.Event, error)
	GetEventByID(id int) (*domain.Event, error)
	UpdateEvent(event *domain.Event) (*domain.Event, error)
	DeleteEvent(id int) (*domain.Event, error)
}
