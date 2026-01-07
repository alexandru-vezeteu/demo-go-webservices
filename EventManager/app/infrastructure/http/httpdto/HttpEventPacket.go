package httpdto

import (
	"eventManager/application/domain"
	"eventManager/infrastructure/http/config"
	"eventManager/infrastructure/http/hateoas"
	"fmt"
	"net/url"
	"strconv"
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

type HttpResponseEventPacketList struct {
	EventPackets []*httpResponseEventPacket `json:"event_packets"`
	Links        map[string]hateoas.Link    `json:"_links"`
	Metadata     *PaginationMetadata        `json:"_metadata,omitempty"`
}

func ToHttpResponseEventPacketList(packets []*domain.EventPacket, selfPath string, serviceURLs *config.ServiceURLs) *HttpResponseEventPacketList {
	if packets == nil {
		packets = []*domain.EventPacket{}
	}

	httpPackets := make([]*httpResponseEventPacket, 0, len(packets))

	for _, packet := range packets {
		resourcePath := fmt.Sprintf("/packets/%d", packet.ID)
		httpPackets = append(httpPackets, &httpResponseEventPacket{
			ID:             packet.ID,
			OwnerID:        packet.OwnerID,
			Name:           packet.Name,
			Location:       packet.Location,
			Description:    packet.Description,
			AllocatedSeats: packet.AllocatedSeats,
			Links: map[string]hateoas.Link{
				"self":   hateoas.BuildSelfLink(serviceURLs.EventManager, resourcePath),
				"update": hateoas.BuildUpdateLink(serviceURLs.EventManager, resourcePath),
				"delete": hateoas.BuildDeleteLink(serviceURLs.EventManager, resourcePath),
				"owner": hateoas.BuildRelatedLink(
					fmt.Sprintf("%s/users/%d", serviceURLs.UserManager, packet.OwnerID),
					"owner",
					"GET",
					"Get packet owner",
				),
			},
		})
	}

	return &HttpResponseEventPacketList{
		EventPackets: httpPackets,
		Links: map[string]hateoas.Link{
			"self": hateoas.BuildSelfLink(serviceURLs.EventManager, selfPath),
		},
	}
}

func buildEventPacketFilterQuery(filter *domain.EventPacketFilter, page int) string {
	if filter == nil {
		return ""
	}

	params := url.Values{}

	if filter.Name != nil {
		params.Add("name", *filter.Name)
	}
	if filter.Location != nil {
		params.Add("location", *filter.Location)
	}
	if filter.Description != nil {
		params.Add("description", *filter.Description)
	}
	if filter.MinSeats != nil {
		params.Add("min_seats", strconv.Itoa(*filter.MinSeats))
	}
	if filter.MaxSeats != nil {
		params.Add("max_seats", strconv.Itoa(*filter.MaxSeats))
	}
	if filter.OrderBy != nil {
		params.Add("order_by", *filter.OrderBy)
	}

	if filter.PerPage != nil {
		params.Add("per_page", strconv.Itoa(*filter.PerPage))
	}

	params.Add("page", strconv.Itoa(page))

	return params.Encode()
}

func ToHttpResponseEventPacketListWithPagination(packets []*domain.EventPacket, filter *domain.EventPacketFilter, serviceURLs *config.ServiceURLs) *HttpResponseEventPacketList {
	if packets == nil {
		packets = []*domain.EventPacket{}
	}

	httpPackets := make([]*httpResponseEventPacket, 0, len(packets))

	for _, packet := range packets {
		resourcePath := fmt.Sprintf("/packets/%d", packet.ID)
		httpPackets = append(httpPackets, &httpResponseEventPacket{
			ID:             packet.ID,
			OwnerID:        packet.OwnerID,
			Name:           packet.Name,
			Location:       packet.Location,
			Description:    packet.Description,
			AllocatedSeats: packet.AllocatedSeats,
			Links: map[string]hateoas.Link{
				"self":   hateoas.BuildSelfLink(serviceURLs.EventManager, resourcePath),
				"parent": hateoas.BuildParentLink(serviceURLs.EventManager, "/event-packets"),
				"update": hateoas.BuildUpdateLink(serviceURLs.EventManager, resourcePath),
				"delete": hateoas.BuildDeleteLink(serviceURLs.EventManager, resourcePath),
				"owner": hateoas.BuildRelatedLink(
					fmt.Sprintf("%s/users/%d", serviceURLs.UserManager, packet.OwnerID),
					"owner",
					"GET",
					"Get packet owner",
				),
			},
		})
	}

	links := map[string]hateoas.Link{
		"create": hateoas.BuildCreateLink(serviceURLs.EventManager, "/event-packets"),
	}

	currentPage := 1
	perPage := 10
	if filter != nil && filter.Page != nil {
		currentPage = *filter.Page
	}
	if filter != nil && filter.PerPage != nil {
		perPage = *filter.PerPage
	}

	selfQuery := buildEventPacketFilterQuery(filter, currentPage)
	links["self"] = hateoas.BuildPaginationLink(serviceURLs.EventManager, "/event-packets", selfQuery, "self", "Current page")

	firstQuery := buildEventPacketFilterQuery(filter, 1)
	links["first"] = hateoas.BuildPaginationLink(serviceURLs.EventManager, "/event-packets", firstQuery, "first", "First page")

	if currentPage > 1 {
		prevQuery := buildEventPacketFilterQuery(filter, currentPage-1)
		links["prev"] = hateoas.BuildPaginationLink(serviceURLs.EventManager, "/event-packets", prevQuery, "prev", "Previous page")
	}

	if len(packets) == perPage {
		nextQuery := buildEventPacketFilterQuery(filter, currentPage+1)
		links["next"] = hateoas.BuildPaginationLink(serviceURLs.EventManager, "/event-packets", nextQuery, "next", "Next page")
	}

	metadata := &PaginationMetadata{
		Page:    currentPage,
		PerPage: perPage,
	}

	return &HttpResponseEventPacketList{
		EventPackets: httpPackets,
		Links:        links,
		Metadata:     metadata,
	}
}

type HttpCreateEventPacket struct {
	OwnerID        int     `json:"id_owner" binding:"required,min=1"`
	Name           string  `json:"name" binding:"required,min=1,max=255"`
	Location       *string `json:"location" binding:"omitempty,max=500"`
	Description    *string `json:"description" binding:"omitempty,max=1000"`
	AllocatedSeats *int    `json:"allocated_seats" binding:"omitempty,min=1"`
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
	OwnerID        *int    `json:"id_owner" binding:"omitempty,min=1"`
	Name           *string `json:"name" binding:"omitempty,min=1,max=255"`
	Location       *string `json:"location" binding:"omitempty,max=500"`
	Description    *string `json:"description" binding:"omitempty,max=1000"`
	AllocatedSeats *int    `json:"allocated_seats" binding:"omitempty,min=1"`
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

type HttpFilterEventPacket struct {
	Name        *string `json:"name,omitempty"        form:"name"`
	Location    *string `json:"location,omitempty"    form:"location"`
	Description *string `json:"description,omitempty" form:"description"`
	MinSeats    *int    `json:"min_seats,omitempty"   form:"min_seats"`
	MaxSeats    *int    `json:"max_seats,omitempty"   form:"max_seats"`

	Page    *int `json:"page,omitempty"        form:"page"`
	PerPage *int `json:"per_page,omitempty"    form:"per_page"`

	OrderBy *string `json:"order_by,omitempty"    form:"order_by"`
}

func (filter *HttpFilterEventPacket) ToEventPacketFilter() *domain.EventPacketFilter {
	return &domain.EventPacketFilter{
		Name:        filter.Name,
		Location:    filter.Location,
		Description: filter.Description,
		Page:        filter.Page,
		PerPage:     filter.PerPage,
		MinSeats:    filter.MinSeats,
		MaxSeats:    filter.MaxSeats,
		OrderBy:     filter.OrderBy,
	}
}
