package repository

import "eventManager/domain"

type EventPacketRepository interface {
	Create(event *domain.EventPacket) (*domain.EventPacket, error)
	GetByID(id int) (*domain.EventPacket, error)
	Update(event *domain.EventPacket) (*domain.EventPacket, error)
	Delete(event *domain.EventPacket) (*domain.EventPacket, error)
}
