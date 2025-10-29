package controller

import (
	"eventManager/application/service"
	"eventManager/domain"
)

type IEventPacketInclusionController interface {
	CreateEventPacketInclusion(event *domain.EventPacketInclusion) (*domain.EventPacketInclusion, error)
	GetEventsInPacketbyID(id int) (*domain.EventPacketInclusion, error)
	GetEventPacketsByEventID(id int) (*domain.EventPacketInclusion, error)
	DeleteEventPacketInclusion(event *domain.EventPacket) (*domain.EventPacket, error)
}

type eventPacketInclusionController struct {
	service service.IEventPacketInclusionService
}

func NewEventPacketInclusionController(service service.IEventPacketInclusionService) *eventPacketInclusionController {
	return &eventPacketInclusionController{service: service}
}

func (controller *eventPacketInclusionController) CreateEventPacketInclusion(event *domain.EventPacketInclusion) (*domain.EventPacketInclusion, error) {
	return nil, nil

}

func (controller *eventPacketInclusionController) GetEventsInPacketbyID(id int) (*domain.EventPacketInclusion, error) {
	return nil, nil
}
func (controller *eventPacketInclusionController) GetEventPacketsByEventID(id int) (*domain.EventPacketInclusion, error) {
	return nil, nil
}

func (controller *eventPacketInclusionController) DeleteEventPacketInclusion(event *domain.EventPacketInclusion) (*domain.EventPacketInclusion, error) {
	return nil, nil
}
