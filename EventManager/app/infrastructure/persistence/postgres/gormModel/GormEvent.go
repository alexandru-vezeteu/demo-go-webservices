package gormmodel

import (
	"eventManager/domain"
)

type GormEvent struct {
	ID          int     `gorm:"primaryKey;autoIncrement"`
	OwnerID     int     `gorm:"column:id_owner;not null"`
	Name        string  `gorm:"column:name;unique;not null"`
	Location    *string `gorm:"column:location"`
	Description *string `gorm:"column:description"`
	Seats       *int    `gorm:"column:seats"`
}

func (GormEvent) TableName() string {
	return "events"
}

func (ge *GormEvent) ToDomain() *domain.Event {
	return &domain.Event{
		ID:          ge.ID,
		OwnerID:     ge.OwnerID,
		Name:        ge.Name,
		Location:    ge.Location,
		Description: ge.Description,
		Seats:       ge.Seats,
	}
}

func FromEvent(e *domain.Event) *GormEvent {

	return &GormEvent{
		ID:          e.ID,
		OwnerID:     e.OwnerID,
		Name:        e.Name,
		Location:    e.Location,
		Description: e.Description,
		Seats:       e.Seats,
	}
}
