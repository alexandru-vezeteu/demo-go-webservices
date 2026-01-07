package httpdto

import (
	"eventManager/application/domain"
	"eventManager/infrastructure/http/config"
	"eventManager/infrastructure/http/hateoas"
	"fmt"
	"net/url"
	"strconv"
)

type httpResponseEvent struct {
	ID          int                     `json:"id"`
	OwnerID     int                     `json:"id_owner"`
	Name        string                  `json:"name"`
	Location    *string                 `json:"location,omitempty"`
	Description *string                 `json:"description,omitempty"`
	Seats       *int                    `json:"seats,omitempty"`
	Links       map[string]hateoas.Link `json:"_links"`
}

type HttpResponseEvent struct {
	Event *httpResponseEvent `json:"event"`
}

func ToHttpResponseEvent(event *domain.Event, serviceURLs *config.ServiceURLs) *HttpResponseEvent {
	if event == nil {
		return &HttpResponseEvent{}
	}

	resourcePath := fmt.Sprintf("/events/%d", event.ID)

	dto := &httpResponseEvent{
		ID:          event.ID,
		OwnerID:     event.OwnerID,
		Name:        event.Name,
		Location:    event.Location,
		Description: event.Description,
		Seats:       event.Seats,
		Links: map[string]hateoas.Link{
			"self":   hateoas.BuildSelfLink(serviceURLs.EventManager, resourcePath),
			"parent": hateoas.BuildParentLink(serviceURLs.EventManager, "/events"),
			"update": hateoas.BuildUpdateLink(serviceURLs.EventManager, resourcePath),
			"delete": hateoas.BuildDeleteLink(serviceURLs.EventManager, resourcePath),
			"owner": hateoas.BuildRelatedLink(
				fmt.Sprintf("%s/users/%d", serviceURLs.UserManager, event.OwnerID),
				"owner",
				"GET",
				"Get event owner",
			),
			"packets": hateoas.BuildRelatedLink(
				fmt.Sprintf("%s/packets?event_id=%d", serviceURLs.EventManager, event.ID),
				"packets",
				"GET",
				"Get packets for this event",
			),
			"tickets": hateoas.BuildRelatedLink(
				fmt.Sprintf("%s/tickets?event_id=%d", serviceURLs.EventManager, event.ID),
				"tickets",
				"GET",
				"Get tickets for this event",
			),
		},
	}

	return &HttpResponseEvent{
		Event: dto,
	}
}

type HttpResponseEventList struct {
	Events   []*httpResponseEvent    `json:"events"`
	Links    map[string]hateoas.Link `json:"_links"`
	Metadata *PaginationMetadata     `json:"_metadata,omitempty"`
}

type PaginationMetadata struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	TotalItems int `json:"total_items,omitempty"`
	TotalPages int `json:"total_pages,omitempty"`
}

func buildEventFilterQuery(filter *domain.EventFilter, page int) string {
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

func ToHttpResponseEventList(events []*domain.Event, serviceURLs *config.ServiceURLs) *HttpResponseEventList {
	if events == nil {
		return &HttpResponseEventList{
			Events: []*httpResponseEvent{},
			Links: map[string]hateoas.Link{
				"self":   hateoas.BuildSelfLink(serviceURLs.EventManager, "/events"),
				"create": hateoas.BuildCreateLink(serviceURLs.EventManager, "/events"),
			},
		}
	}
	httpEvents := make([]*httpResponseEvent, 0, len(events))

	for _, event := range events {
		resourcePath := fmt.Sprintf("/events/%d", event.ID)
		httpEvents = append(httpEvents, &httpResponseEvent{
			ID:          event.ID,
			OwnerID:     event.OwnerID,
			Name:        event.Name,
			Location:    event.Location,
			Description: event.Description,
			Seats:       event.Seats,
			Links: map[string]hateoas.Link{
				"self":   hateoas.BuildSelfLink(serviceURLs.EventManager, resourcePath),
				"parent": hateoas.BuildParentLink(serviceURLs.EventManager, "/events"),
				"update": hateoas.BuildUpdateLink(serviceURLs.EventManager, resourcePath),
				"delete": hateoas.BuildDeleteLink(serviceURLs.EventManager, resourcePath),
				"owner": hateoas.BuildRelatedLink(
					fmt.Sprintf("%s/users/%d", serviceURLs.UserManager, event.OwnerID),
					"owner",
					"GET",
					"Get event owner",
				),
			},
		})
	}

	return &HttpResponseEventList{
		Events: httpEvents,
		Links: map[string]hateoas.Link{
			"self":   hateoas.BuildSelfLink(serviceURLs.EventManager, "/events"),
			"create": hateoas.BuildCreateLink(serviceURLs.EventManager, "/events"),
		},
	}
}

func ToHttpResponseEventListCustom(events []*domain.Event, selfPath string, serviceURLs *config.ServiceURLs) *HttpResponseEventList {
	if events == nil {
		events = []*domain.Event{}
	}

	httpEvents := make([]*httpResponseEvent, 0, len(events))

	for _, event := range events {
		resourcePath := fmt.Sprintf("/events/%d", event.ID)
		httpEvents = append(httpEvents, &httpResponseEvent{
			ID:          event.ID,
			OwnerID:     event.OwnerID,
			Name:        event.Name,
			Location:    event.Location,
			Description: event.Description,
			Seats:       event.Seats,
			Links: map[string]hateoas.Link{
				"self":   hateoas.BuildSelfLink(serviceURLs.EventManager, resourcePath),
				"parent": hateoas.BuildParentLink(serviceURLs.EventManager, "/events"),
				"update": hateoas.BuildUpdateLink(serviceURLs.EventManager, resourcePath),
				"delete": hateoas.BuildDeleteLink(serviceURLs.EventManager, resourcePath),
				"owner": hateoas.BuildRelatedLink(
					fmt.Sprintf("%s/users/%d", serviceURLs.UserManager, event.OwnerID),
					"owner",
					"GET",
					"Get event owner",
				),
			},
		})
	}

	return &HttpResponseEventList{
		Events: httpEvents,
		Links: map[string]hateoas.Link{
			"self": hateoas.BuildSelfLink(serviceURLs.EventManager, selfPath),
		},
	}
}

func ToHttpResponseEventListWithPagination(events []*domain.Event, filter *domain.EventFilter, totalCount int, serviceURLs *config.ServiceURLs) *HttpResponseEventList {
	if events == nil {
		events = []*domain.Event{}
	}

	httpEvents := make([]*httpResponseEvent, 0, len(events))

	for _, event := range events {
		resourcePath := fmt.Sprintf("/events/%d", event.ID)
		httpEvents = append(httpEvents, &httpResponseEvent{
			ID:          event.ID,
			OwnerID:     event.OwnerID,
			Name:        event.Name,
			Location:    event.Location,
			Description: event.Description,
			Seats:       event.Seats,
			Links: map[string]hateoas.Link{
				"self":   hateoas.BuildSelfLink(serviceURLs.EventManager, resourcePath),
				"parent": hateoas.BuildParentLink(serviceURLs.EventManager, "/events"),
				"update": hateoas.BuildUpdateLink(serviceURLs.EventManager, resourcePath),
				"delete": hateoas.BuildDeleteLink(serviceURLs.EventManager, resourcePath),
				"owner": hateoas.BuildRelatedLink(
					fmt.Sprintf("%s/users/%d", serviceURLs.UserManager, event.OwnerID),
					"owner",
					"GET",
					"Get event owner",
				),
			},
		})
	}

	links := map[string]hateoas.Link{
		"create": hateoas.BuildCreateLink(serviceURLs.EventManager, "/events"),
	}

	currentPage := 1
	perPage := 10
	if filter != nil && filter.Page != nil {
		currentPage = *filter.Page
	}
	if filter != nil && filter.PerPage != nil {
		perPage = *filter.PerPage
	}

	totalPages := (totalCount + perPage - 1) / perPage
	if totalPages < 1 {
		totalPages = 1
	}

	selfQuery := buildEventFilterQuery(filter, currentPage)
	links["self"] = hateoas.BuildPaginationLink(serviceURLs.EventManager, "/events", selfQuery, "self", "Current page")

	firstQuery := buildEventFilterQuery(filter, 1)
	links["first"] = hateoas.BuildPaginationLink(serviceURLs.EventManager, "/events", firstQuery, "first", "First page")

	if currentPage > 1 {
		prevQuery := buildEventFilterQuery(filter, currentPage-1)
		links["prev"] = hateoas.BuildPaginationLink(serviceURLs.EventManager, "/events", prevQuery, "prev", "Previous page")
	}

	if currentPage < totalPages {
		nextQuery := buildEventFilterQuery(filter, currentPage+1)
		links["next"] = hateoas.BuildPaginationLink(serviceURLs.EventManager, "/events", nextQuery, "next", "Next page")
	}

	if totalPages > 1 {
		lastQuery := buildEventFilterQuery(filter, totalPages)
		links["last"] = hateoas.BuildPaginationLink(serviceURLs.EventManager, "/events", lastQuery, "last", "Last page")
	}

	metadata := &PaginationMetadata{
		Page:       currentPage,
		PerPage:    perPage,
		TotalItems: totalCount,
		TotalPages: totalPages,
	}

	return &HttpResponseEventList{
		Events:   httpEvents,
		Links:    links,
		Metadata: metadata,
	}
}

type HttpCreateEvent struct {
	OwnerID     int     `json:"id_owner" binding:"required,min=1"`
	Name        string  `json:"name" binding:"required,min=1,max=255"`
	Location    *string `json:"location" binding:"omitempty,max=500"`
	Description *string `json:"description" binding:"omitempty,max=1000"`
	Seats       *int    `json:"seats" binding:"omitempty,min=1"`
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
	OwnerID     *int    `json:"id_owner" binding:"omitempty,min=1"`
	Name        *string `json:"name" binding:"omitempty,min=1,max=255"`
	Location    *string `json:"location" binding:"omitempty,max=500"`
	Description *string `json:"description" binding:"omitempty,max=1000"`
	Seats       *int    `json:"seats" binding:"omitempty,min=1"`
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
