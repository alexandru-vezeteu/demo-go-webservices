package httpdto

import (
	"eventManager/domain"
)

// poate ar trb sa l sparg in req response per fiecare caz dar cred ca s ar aduna prea mule?
type HttpEvent struct {
	ID          int    `json:"id"`
	OwnerID     int    `json:"id_owner"`
	Name        string `json:"name"`
	Location    string `json:"location"`
	Description string `json:"description"`
	Seats       int    `json:"seats"`
}

func (event *HttpEvent) ToEvent() *domain.Event {

	return &domain.Event{
		ID:          event.ID,
		OwnerID:     event.OwnerID,
		Name:        event.Name,
		Location:    event.Location,
		Description: event.Description,
		Seats:       event.Seats,
	}
}

func ToHttpEvent(event *domain.Event) *HttpEvent {
	return &HttpEvent{
		ID:          event.ID,
		OwnerID:     event.OwnerID,
		Name:        event.Name,
		Location:    event.Location,
		Description: event.Description,
		Seats:       event.Seats,
	}
}
