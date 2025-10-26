package httpdto

import (
	"eventManager/domain"
)

type HttpResponseEvent struct {
	ID          int    `json:"id"`
	OwnerID     int    `json:"id_owner"`
	Name        string `json:"name"`
	Location    string `json:"location"`
	Description string `json:"description"`
	Seats       int    `json:"seats"`
}

func (event *HttpResponseEvent) ToEvent() *domain.Event {

	return &domain.Event{
		ID:          event.ID,
		OwnerID:     event.OwnerID,
		Name:        event.Name,
		Location:    event.Location,
		Description: event.Description,
		Seats:       event.Seats,
	}
}
func ToHttpResponseEvent(event *domain.Event) *HttpResponseEvent {
	return &HttpResponseEvent{
		ID:          event.ID,
		OwnerID:     event.OwnerID,
		Name:        event.Name,
		Location:    event.Location,
		Description: event.Description,
		Seats:       event.Seats,
	}
}

type HttpCreateEvent struct {
	OwnerID     int    `json:"id_owner"`
	Name        string `json:"name"`
	Location    string `json:"location"`
	Description string `json:"description"`
	Seats       int    `json:"seats"`
}

func (event *HttpCreateEvent) ToEvent() *domain.Event {

	return &domain.Event{
		OwnerID:     event.OwnerID,
		Name:        event.Name,
		Location:    event.Location,
		Description: event.Description,
		Seats:       event.Seats,
	}
}

type HttpUpdateEvent struct {
	ID          int    `json:"id"`
	OwnerID     int    `json:"id_owner"`
	Name        string `json:"name"`
	Location    string `json:"location"`
	Description string `json:"description"`
	Seats       int    `json:"seats"`
}

func (event *HttpUpdateEvent) ToEvent() *domain.Event {

	return &domain.Event{
		ID:          event.ID,
		OwnerID:     event.OwnerID,
		Name:        event.Name,
		Location:    event.Location,
		Description: event.Description,
		Seats:       event.Seats,
	}
}
