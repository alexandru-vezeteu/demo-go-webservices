package httpdto

import (
	"eventManager/application/domain"
	"eventManager/infrastructure/http/config"
	"eventManager/infrastructure/http/hateoas"
	"fmt"
)

type HttpCreateTicket struct {
	PacketID *int `json:"packet_id" binding:"omitempty,min=1"`
	EventID  *int `json:"event_id" binding:"omitempty,min=1"`
}

type HttpResponseTicket struct {
	Code     string                  `json:"code"`
	PacketID *int                    `json:"packet_id"`
	EventID  *int                    `json:"event_id"`
	Links    map[string]hateoas.Link `json:"_links"`
}

type HttpUpdateTicket struct {
	PacketID *int `json:"packet_id" binding:"omitempty,min=1"`
	EventID  *int `json:"event_id" binding:"omitempty,min=1"`
}

func (dto *HttpCreateTicket) ToTicket() *domain.Ticket {
	return &domain.Ticket{
		PacketID: dto.PacketID,
		EventID:  dto.EventID,
	}
}

func ToHttpResponseTicket(ticket *domain.Ticket, serviceURLs *config.ServiceURLs) *HttpResponseTicket {
	resourcePath := fmt.Sprintf("/tickets/%s", ticket.Code)

	links := map[string]hateoas.Link{
		"self":   hateoas.BuildSelfLink(serviceURLs.EventManager, resourcePath),
		"update": hateoas.BuildUpdateLink(serviceURLs.EventManager, resourcePath),
		"delete": hateoas.BuildDeleteLink(serviceURLs.EventManager, resourcePath),
	}

	if ticket.EventID != nil {
		links["event"] = hateoas.BuildRelatedLink(
			fmt.Sprintf("%s/events/%d", serviceURLs.EventManager, *ticket.EventID),
			"event",
			"GET",
			"Get event for this ticket",
		)
	}

	if ticket.PacketID != nil {
		links["packet"] = hateoas.BuildRelatedLink(
			fmt.Sprintf("%s/packets/%d", serviceURLs.EventManager, *ticket.PacketID),
			"packet",
			"GET",
			"Get packet for this ticket",
		)
	}

	return &HttpResponseTicket{
		Code:     ticket.Code,
		PacketID: ticket.PacketID,
		EventID:  ticket.EventID,
		Links:    links,
	}
}

func (dto *HttpUpdateTicket) ToUpdateMap() map[string]interface{} {
	updates := make(map[string]interface{})

	if dto.PacketID != nil {
		updates["packet_id"] = *dto.PacketID
	}

	if dto.EventID != nil {
		updates["event_id"] = *dto.EventID
	}

	return updates
}
