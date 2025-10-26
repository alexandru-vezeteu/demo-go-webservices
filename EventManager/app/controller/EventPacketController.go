package controller

import (
	"errors"
	"eventManager/application/service"
	"eventManager/domain"
)

type IEventPacketController interface {
	CreateEventPacket(event *domain.EventPacket) (*domain.EventPacket, error)
	GetEventPacketByID(id int) (*domain.EventPacket, error)
	UpdateEventPacket(event *domain.EventPacket) (*domain.EventPacket, error)
	DeleteEventPacket(id int) (*domain.EventPacket, error)
}

type eventPacketController struct {
	service service.IEventPacketService
}

func NewEventPacketController(service service.IEventPacketService) *eventPacketController {
	return &eventPacketController{service: service}
}

func (c *eventPacketController) CreateEventPacket(event *domain.EventPacket) (*domain.EventPacket, error) {
	return nil, errors.New("TODO")
}

func (c *eventPacketController) GetEventPacketByID(id int) (*domain.EventPacket, error) {
	return nil, errors.New("TODO")
}

func (c *eventPacketController) UpdateEventPacket(event *domain.EventPacket) (*domain.EventPacket, error) {
	return nil, errors.New("TODO")
}

func (c *eventPacketController) DeleteEventPacket(id int) (*domain.EventPacket, error) {
	return nil, errors.New("TODO")
}
