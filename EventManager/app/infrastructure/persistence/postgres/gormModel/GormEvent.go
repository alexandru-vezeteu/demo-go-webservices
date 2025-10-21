package gormmodel

import (
	"errors"
	"eventManager/domain"
)

type GormEvent struct {
	ID          int    `gorm:"primaryKey;autoIncrement"`
	OwnerID     int    `gorm:"column:id_owner;not null"`
	Name        string `gorm:"column:name;unique;not null"`
	Location    string `gorm:"column:location"`
	Description string `gorm:"column:description"`
	Seats       int    `gorm:"column:seats"`
}

func (GormEvent) TableName() string {
	return "events"
}

func (ge *GormEvent) ToDomain() (*domain.Event, error) {
	if ge == nil {
		return nil, errors.New("gormevent is nil")
	}
	return &domain.Event{
		ID:          ge.ID,
		OwnerID:     ge.OwnerID,
		Name:        ge.Name,
		Location:    ge.Location,
		Description: ge.Description,
		Seats:       ge.Seats,
	}, nil
}

func FromDomain(e *domain.Event) (*GormEvent, error) {
	if e == nil {
		return nil, errors.New("domain event is nil")
	}
	return &GormEvent{
		ID:          e.ID,
		OwnerID:     e.OwnerID,
		Name:        e.Name,
		Location:    e.Location,
		Description: e.Description,
		Seats:       e.Seats,
	}, nil
}
