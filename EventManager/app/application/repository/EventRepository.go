package repository

import "eventManager/domain"

type EventRepository interface {
	Create(event *domain.Event) (*domain.Event, error)
	GetByID(id int) (*domain.Event, error)
	Update(event *domain.Event) (*domain.Event, error)
	Delete(event *domain.Event) (*domain.Event, error)
}
