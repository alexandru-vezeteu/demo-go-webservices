package controller

import (
	"eventManager/application/domain"
	"eventManager/application/service"
)

// aici nu prea are sens dar l am lasat asa (are sens in momentul in care am mai multe "servicii" si le agreg? - flowuri mai complexe)

type IEventController interface {
	CreateEvent(event *domain.Event) (*domain.Event, error)
	GetEventByID(id int) (*domain.Event, error)
	UpdateEvent(id int, updates map[string]interface{}) (*domain.Event, error)
	DeleteEvent(id int) (*domain.Event, error)
	FilterEvents(filter *domain.EventFilter) ([]*domain.Event, error)
}

type eventController struct {
	service service.IEventService
}

func NewEventController(service service.IEventService) *eventController {
	return &eventController{service: service}
}

func (c *eventController) CreateEvent(event *domain.Event) (*domain.Event, error) {
	return c.service.CreateEvent(event)
}

func (c *eventController) GetEventByID(id int) (*domain.Event, error) {
	return c.service.GetEventByID(id)
}

func (c *eventController) UpdateEvent(id int, updates map[string]interface{}) (*domain.Event, error) {
	return c.service.UpdateEvent(id, updates)
}

func (c *eventController) DeleteEvent(id int) (*domain.Event, error) {
	return c.service.DeleteEvent(id)
}

func (c *eventController) FilterEvents(filter *domain.EventFilter) ([]*domain.Event, error) {
	return c.service.FilterEvents(filter)
}
