package httpdto

import (
	"eventManager/application/domain"
	"eventManager/infrastructure/http/config"
	"eventManager/infrastructure/http/hateoas"
	"fmt"
)

type HttpResponseEventPacket struct {
	EventPacket *httpResponseEventPacket `json:"event_packet"`
}

type httpResponseEventPacket struct {
	ID             int                     `json:"id"`
	OwnerID        int                     `json:"id_owner"`
	Name           string                  `json:"name"`
	Location       *string                 `json:"location"`
	Description    *string                 `json:"description"`
	AllocatedSeats *int                    `json:"allocated_seats"`
	Links          map[string]hateoas.Link `json:"_links"`
}

func ToHttpResponseEventPacket(event *domain.EventPacket, serviceURLs *config.ServiceURLs) *HttpResponseEventPacket {
	resourcePath := fmt.Sprintf("/packets/%d", event.ID)

	dto := &httpResponseEventPacket{
		ID:             event.ID,
		OwnerID:        event.OwnerID,
		Name:           event.Name,
		Location:       event.Location,
		Description:    event.Description,
		AllocatedSeats: event.AllocatedSeats,
		Links: map[string]hateoas.Link{
			"self":   hateoas.BuildSelfLink(serviceURLs.EventManager, resourcePath),
			"update": hateoas.BuildUpdateLink(serviceURLs.EventManager, resourcePath),
			"delete": hateoas.BuildDeleteLink(serviceURLs.EventManager, resourcePath),
			"owner": hateoas.BuildRelatedLink(
				fmt.Sprintf("%s/users/%d", serviceURLs.UserManager, event.OwnerID),
				"owner",
				"GET",
				"Get packet owner",
			),
			"events": hateoas.BuildRelatedLink(
				fmt.Sprintf("%s/events?packet_id=%d", serviceURLs.EventManager, event.ID),
				"events",
				"GET",
				"Get events in this packet",
			),
			"tickets": hateoas.BuildRelatedLink(
				fmt.Sprintf("%s/tickets?packet_id=%d", serviceURLs.EventManager, event.ID),
				"tickets",
				"GET",
				"Get tickets for this packet",
			),
		},
	}
	return &HttpResponseEventPacket{
		EventPacket: dto,
	}
}

type HttpCreateEventPacket struct {
	OwnerID        int     `json:"id_owner" binding:"required"`
	Name           string  `json:"name" binding:"required"`
	Location       *string `json:"location"`
	Description    *string `json:"description"`
	AllocatedSeats *int    `json:"allocated_seats"`
}

func (event *HttpCreateEventPacket) ToEventPacket() *domain.EventPacket {

	return &domain.EventPacket{
		OwnerID:        event.OwnerID,
		Name:           event.Name,
		Location:       event.Location,
		Description:    event.Description,
		AllocatedSeats: event.AllocatedSeats,
	}
}

type HttpUpdateEventPacket struct {
	OwnerID        *int    `json:"id_owner"`
	Name           *string `json:"name"`
	Location       *string `json:"location"`
	Description    *string `json:"description"`
	AllocatedSeats *int    `json:"allocated_seats"`
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
	if event.AllocatedSeats != nil {
		updates["allocated_seats"] = *event.AllocatedSeats
	}

	return updates
}
