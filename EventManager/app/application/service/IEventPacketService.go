package service

import "eventManager/domain"

type IEventPacketService interface {
	CreateEventPacket(event *domain.EventPacket) (*domain.EventPacket, error)
	GetEventPacketByID(id int) (*domain.EventPacket, error)
	UpdateEventPacket(id int, updates map[string]interface{}) (*domain.EventPacket, error)
	DeleteEventPacket(id int) (*domain.EventPacket, error)
}
