package repository

import "eventManager/application/domain"

type EventPacketRepository interface {
	Create(event *domain.EventPacket) (*domain.EventPacket, error)
	GetByID(id int) (*domain.EventPacket, error)
	Update(id int, updates map[string]interface{}) (*domain.EventPacket, error)
	Delete(id int) (*domain.EventPacket, error)
	CountSoldTickets(id int) (int, error)
}
