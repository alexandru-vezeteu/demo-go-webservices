package gormmodel

import (
	"eventManager/application/domain"
)

type GormEventPacket struct {
	ID             int     `gorm:"primaryKey;autoIncrement"`
	OwnerID        int     `gorm:"column:id_owner;not null"`
	Name           string  `gorm:"column:name;unique;not null"`
	Location       *string `gorm:"column:location"`
	Description    *string `gorm:"column:description"`
	AllocatedSeats *int    `gorm:"column:allocated_seats"`
}

func (GormEventPacket) TableName() string {
	return "eventsPacket"
}

func (ge *GormEventPacket) ToDomain() *domain.EventPacket {
	return &domain.EventPacket{
		ID:             ge.ID,
		OwnerID:        ge.OwnerID,
		Name:           ge.Name,
		Location:       ge.Location,
		Description:    ge.Description,
		AllocatedSeats: ge.AllocatedSeats,
	}
}

func FromEventPacket(e *domain.EventPacket) *GormEventPacket {

	return &GormEventPacket{
		ID:             e.ID,
		OwnerID:        e.OwnerID,
		Name:           e.Name,
		Location:       e.Location,
		Description:    e.Description,
		AllocatedSeats: e.AllocatedSeats,
	}
}
