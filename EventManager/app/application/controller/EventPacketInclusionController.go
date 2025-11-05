package controller

import (
	"eventManager/application/domain"
	"eventManager/application/service"
)

type IEventPacketInclusionController interface {
	CreateEventPacketInclusion(event *domain.EventPacketInclusion) (*domain.EventPacketInclusion, error)
	GetEventsByPacketID(packetID int) ([]*domain.Event, error)
	GetEventPacketsByEventID(eventID int) ([]*domain.EventPacket, error)
	DeleteEventPacketInclusion(eventID, packetID int) (*domain.EventPacketInclusion, error)
	Update(eventID, packetID int, updates map[string]interface{}) (*domain.EventPacketInclusion, error)
}

type eventPacketInclusionController struct {
	service service.IEventPacketInclusionService
}

func NewEventPacketInclusionController(service service.IEventPacketInclusionService) *eventPacketInclusionController {
	return &eventPacketInclusionController{service: service}
}

func (controller *eventPacketInclusionController) CreateEventPacketInclusion(event *domain.EventPacketInclusion) (*domain.EventPacketInclusion, error) {
	return controller.service.CreateEventPacketInclusion(event)
}
func (controller *eventPacketInclusionController) GetEventsByPacketID(packetID int) ([]*domain.Event, error) {
	return controller.service.GetEventsByPacketID(packetID)
}
func (controller *eventPacketInclusionController) GetEventPacketsByEventID(eventID int) ([]*domain.EventPacket, error) {
	return controller.service.GetEventPacketsByEventID(eventID)
}
func (controller *eventPacketInclusionController) DeleteEventPacketInclusion(eventID, packetID int) (*domain.EventPacketInclusion, error) {
	return controller.service.DeleteEventPacketInclusion(eventID, packetID)
}
func (controller *eventPacketInclusionController) Update(eventID, packetID int, updates map[string]interface{}) (*domain.EventPacketInclusion, error) {
	return controller.service.Update(eventID, packetID, updates)
}
