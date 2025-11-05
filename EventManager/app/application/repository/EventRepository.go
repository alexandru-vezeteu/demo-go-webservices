package repository

import "eventManager/application/domain"

type EventRepository interface {
	Create(event *domain.Event) (*domain.Event, error)
	GetByID(id int) (*domain.Event, error)
	Update(id int, updates map[string]interface{}) (*domain.Event, error)
	Delete(id int) (*domain.Event, error)
	FilterEvents(filter *domain.EventFilter) ([]*domain.Event, error)
}
