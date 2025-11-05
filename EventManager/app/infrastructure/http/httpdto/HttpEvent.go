package httpdto

import (
	"eventManager/application/domain"
	"eventManager/infrastructure/http"
	"fmt"
)

type httpResponseEvent struct {
	ID          int                  `json:"id"`
	OwnerID     int                  `json:"id_owner"`
	Name        string               `json:"name"`
	Location    *string              `json:"location,omitempty"`
	Description *string              `json:"description,omitempty"`
	Seats       *int                 `json:"seats,omitempty"`
	Links       map[string]http.Link `json:"_links"`
}
type HttpResponseEvent struct {
	Event *httpResponseEvent `json:"event"`
}

func ToHttpResponseEvent(event *domain.Event) *HttpResponseEvent {
	if event == nil {
		return &HttpResponseEvent{}
	}
	dto := &httpResponseEvent{
		ID:          event.ID,
		OwnerID:     event.OwnerID,
		Name:        event.Name,
		Location:    event.Location,
		Description: event.Description,
		Seats:       event.Seats,
	}
	prefix := "/api/event-manager"
	dto.Links = map[string]http.Link{
		"self": {
			Href: fmt.Sprintf("%s/events/%d", prefix, event.ID),
			Type: "GET",
		},
		"update": {
			Href: fmt.Sprintf("%s/events/%d", prefix, event.ID),
			Type: "PATCH",
		},
		"delete": {
			Href: fmt.Sprintf("%s/events/%d", prefix, event.ID),
			Type: "DELETE",
		},
	}
	return &HttpResponseEvent{
		Event: dto,
	}
}

type HttpResponseEventList struct {
	Events []*httpResponseEvent `json:"events"`
}

func ToHttpResponseEventList(events []*domain.Event) *HttpResponseEventList {
	if events == nil {
		return &HttpResponseEventList{Events: []*httpResponseEvent{}}
	}
	httpEvents := make([]*httpResponseEvent, 0, len(events))

	for _, event := range events {
		httpEvents = append(httpEvents, &httpResponseEvent{
			ID:          event.ID,
			OwnerID:     event.OwnerID,
			Name:        event.Name,
			Location:    event.Location,
			Description: event.Description,
			Seats:       event.Seats,
		})
	}

	return &HttpResponseEventList{Events: httpEvents}
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

type HttpFilterEvent struct {
	Name        *string `json:"name,omitempty"        form:"name"`
	Location    *string `json:"location,omitempty"    form:"location"`
	Description *string `json:"description,omitempty" form:"description"`
	MinSeats    *int    `json:"min_seats,omitempty"   form:"min_seats"`
	MaxSeats    *int    `json:"max_seats,omitempty"   form:"max_seats"`

	Page    *int `json:"page,omitempty"        form:"page"`
	PerPage *int `json:"per_page,omitempty"    form:"per_page"`

	OrderBy *string `json:"order_by,omitempty"    form:"order_by"`
}

func (filter *HttpFilterEvent) ToEventFilter() *domain.EventFilter {
	return &domain.EventFilter{
		Name:        filter.Name,
		Location:    filter.Location,
		Description: filter.Description,
		Page:        filter.Page,
		PerPage:     filter.PerPage,
		MinSeats:    filter.MinSeats,
		MaxSeats:    filter.MaxSeats,
	}
}
