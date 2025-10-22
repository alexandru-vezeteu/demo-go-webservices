package controller

import (
	"errors"
	"eventManager/application/service"
	"eventManager/domain"
)

// aici nu prea are sens dar l am lasat asa (are sens in momentul in care am mai multe "servicii" si le agreg? - flowuri mai complexe)

type IEventController interface {
	CreateEvent(event *domain.Event) (*domain.Event, error)
	GetEventByID(id int) (*domain.Event, error)
	UpdateEvent(event *domain.Event) (*domain.Event, error)
	DeleteEvent(id int) (*domain.Event, error)
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

func (c *eventController) UpdateEvent(event *domain.Event) (*domain.Event, error) {
	return c.service.UpdateEvent(event)
}

func (c *eventController) DeleteEvent(id int) (*domain.Event, error) {
	return nil, errors.New("TODO")
}
