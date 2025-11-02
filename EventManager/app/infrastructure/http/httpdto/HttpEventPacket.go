package httpdto

import (
	"eventManager/domain"
)

type HttpResponseEventPacket struct {
	ID          int     `json:"id"`
	OwnerID     int     `json:"id_owner"`
	Name        string  `json:"name"`
	Location    *string `json:"location"`
	Description *string `json:"description"`
}

func (event *HttpResponseEventPacket) ToEventPacket() *domain.EventPacket {
	return &domain.EventPacket{
		ID:          event.ID,
		OwnerID:     event.OwnerID,
		Name:        event.Name,
		Location:    event.Location,
		Description: event.Description,
	}
}

func ToHttpResponseEventPacket(event *domain.EventPacket) *HttpResponseEventPacket {
	return &HttpResponseEventPacket{
		ID:          event.ID,
		OwnerID:     event.OwnerID,
		Name:        event.Name,
		Location:    event.Location,
		Description: event.Description,
	}
}

type HttpCreateEventPacket struct {
	OwnerID     int     `json:"id_owner" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Location    *string `json:"location"`
	Description *string `json:"description"`
}

func (event *HttpCreateEventPacket) ToEventPacket() *domain.EventPacket {

	return &domain.EventPacket{
		OwnerID:     event.OwnerID,
		Name:        event.Name,
		Location:    event.Location,
		Description: event.Description,
	}
}

type HttpUpdateEventPacket struct {
	OwnerID     *int    `json:"id_owner"`
	Name        *string `json:"name"`
	Location    *string `json:"location"`
	Description *string `json:"description"`
}

func (event *HttpUpdateEventPacket) ToUpdateMap() map[string]interface{} {
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

	return updates
}
