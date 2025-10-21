package httpdto

import (
	"errors"
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

func (event *HttpEvent) ToEvent() (*domain.Event, error) {
	if event == nil {
		return nil, errors.New("httpevent is nil")
	}
	return &domain.Event{
		ID:          event.ID,
		OwnerID:     event.OwnerID,
		Name:        event.Name,
		Location:    event.Location,
		Description: event.Description,
		Seats:       event.Seats,
	}, nil
}

func ToHttpEvent(event *domain.Event) (*HttpEvent, error) {
	if event == nil {
		return nil, errors.New("domain.event is nil")
	}
	return &HttpEvent{
		ID:          event.ID,
		OwnerID:     event.OwnerID,
		Name:        event.Name,
		Location:    event.Location,
		Description: event.Description,
		Seats:       event.Seats,
	}, nil
}
