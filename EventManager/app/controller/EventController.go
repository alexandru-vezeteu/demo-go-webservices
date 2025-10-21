package controller

import (
	"errors"
	"eventManager/application/service"
	"eventManager/domain"
)

// aici nu prea are sens dar l am lasat asa (are sens in momentul in care am mai multe "servicii" si le agreg)
type EventController struct {
	service service.IEventService
}

func NewEventController(service service.IEventService) *EventController {
	return &EventController{service: service}
}

func (c *EventController) CreateEvent(event *domain.Event) (*domain.Event, error) {
	return c.service.CreateEvent(event)
}

func (c *EventController) GetEventByID(id int) (*domain.Event, error) {
	return nil, errors.New("TODO")
}

func (c *EventController) UpdateEvent(id int, event *domain.Event) (*domain.Event, error) {
	return nil, errors.New("TODO")
}

func (c *EventController) DeleteEvent(id int) (*domain.Event, error) {
	return nil, errors.New("TODO")
}
