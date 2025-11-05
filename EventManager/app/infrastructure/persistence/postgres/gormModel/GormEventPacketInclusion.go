package gormmodel

import (
	"eventManager/application/domain"
)

type GormEventPacketInclusion struct {
	PacketID       int             `gorm:"primaryKey;column:packet_id"`
	EventID        int             `gorm:"primaryKey;column:event_id"`
	AllocatedSeats *int            `gorm:"column:allocated_seats"`
	Packet         GormEventPacket `gorm:"foreignKey:PacketID;references:ID"`
	Event          GormEvent       `gorm:"foreignKey:EventID;references:ID"`
}

func (GormEventPacketInclusion) TableName() string {
	return "events_packet_inclusion"
}

func (ge *GormEventPacketInclusion) ToDomain() *domain.EventPacketInclusion {
	return &domain.EventPacketInclusion{
		PacketID:       ge.PacketID,
		EventID:        ge.EventID,
		AllocatedSeats: ge.AllocatedSeats,
	}
}

func FromEventPacketInclusion(e *domain.EventPacketInclusion) *GormEventPacketInclusion {

	return &GormEventPacketInclusion{
		PacketID:       e.PacketID,
		EventID:        e.EventID,
		AllocatedSeats: e.AllocatedSeats,
	}
}
