package gormmodel

import (
	"eventManager/application/domain"
)

type GormTicket struct {
	Code     string           `gorm:"primaryKey;column:code"`
	PacketID *int             `gorm:"column:packet_id"`
	EventID  *int             `gorm:"column:event_id"`
	Packet   *GormEventPacket `gorm:"foreignKey:PacketID;references:ID"`
	Event    *GormEvent       `gorm:"foreignKey:EventID;references:ID"`
}

func (GormTicket) TableName() string {
	return "tickets"
}

func (gt *GormTicket) ToDomain() *domain.Ticket {
	return &domain.Ticket{
		Code:     gt.Code,
		PacketID: gt.PacketID,
		EventID:  gt.EventID,
	}
}

func FromTicket(t *domain.Ticket) *GormTicket {
	return &GormTicket{
		Code:     t.Code,
		PacketID: t.PacketID,
		EventID:  t.EventID,
	}
}
