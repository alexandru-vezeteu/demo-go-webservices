package httpdto

import (
	"eventManager/domain"
)

type HttpResponseEvent struct {
	ID          int     `json:"id"`
	OwnerID     int     `json:"id_owner"`
	Name        string  `json:"name"`
	Location    *string `json:"location,omitempty"`
	Description *string `json:"description,omitempty"`
	Seats       *int    `json:"seats,omitempty"`
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
	if event == nil {
		return nil
	}
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
	OwnerID     int     `json:"id_owner" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Location    *string `json:"location"`
	Description *string `json:"description"`
	Seats       *int    `json:"seats"`
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
	OwnerID     *int    `json:"id_owner"`
	Name        *string `json:"name"`
	Location    *string `json:"location"`
	Description *string `json:"description"`
	Seats       *int    `json:"seats"`
}

func (event *HttpUpdateEvent) ToUpdateMap() map[string]interface{} {
	updates := make(map[string]interface{})

	if event.OwnerID != nil {
		updates["id_owner"] = *event.OwnerID
	}
	if event.Name != nil {
		updates["name"] = *event.Name
	}
	if event.Location != nil {
		updates["location"] = *event.Location
	}
	if event.Description != nil {
		updates["description"] = *event.Description
	}
	if event.Seats != nil {
		updates["seats"] = *event.Seats
	}

	return updates
}
