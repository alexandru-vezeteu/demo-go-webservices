package controller

import (
	"eventManager/application/domain"
	"eventManager/application/service"
)

type IEventPacketController interface {
	CreateEventPacket(event *domain.EventPacket) (*domain.EventPacket, error)
	GetEventPacketByID(id int) (*domain.EventPacket, error)
	UpdateEventPacket(id int, updates map[string]interface{}) (*domain.EventPacket, error)
	DeleteEventPacket(id int) (*domain.EventPacket, error)
}

type eventPacketController struct {
	service service.IEventPacketService
}

func NewEventPacketController(service service.IEventPacketService) *eventPacketController {
	return &eventPacketController{service: service}
}

func (c *eventPacketController) CreateEventPacket(event *domain.EventPacket) (*domain.EventPacket, error) {
	return c.service.CreateEventPacket(event)
}

func (c *eventPacketController) GetEventPacketByID(id int) (*domain.EventPacket, error) {
	return c.service.GetEventPacketByID(id)
}

func (c *eventPacketController) UpdateEventPacket(id int, updates map[string]interface{}) (*domain.EventPacket, error) {
	return c.service.UpdateEventPacket(id, updates)
}

func (c *eventPacketController) DeleteEventPacket(id int) (*domain.EventPacket, error) {
	return c.service.DeleteEventPacket(id)
}
