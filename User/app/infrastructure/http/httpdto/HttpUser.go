package httpdto

import (
	"fmt"
	"userService/application/domain"
	"userService/infrastructure/http"
	"userService/infrastructure/http/config"
	"userService/infrastructure/http/hateoas"
)

type HttpTicket struct {
	Code string `json:"code"`
}

type httpResponseUser struct {
	ID               int                  `json:"id"`
	Email            string               `json:"email"`
	FirstName        string               `json:"first_name"`
	LastName         string               `json:"last_name"`
	SocialMediaLinks *string              `json:"social_media_links,omitempty"`
	TicketList       []HttpTicket         `json:"ticket_list,omitempty"`
	FirstNamePrivate bool                 `json:"first_name_private"`
	LastNamePrivate  bool                 `json:"last_name_private"`
	Links            map[string]http.Link `json:"_links"`
}

type HttpResponseUser struct {
	User *httpResponseUser `json:"user"`
}

func ToHttpResponseUser(user *domain.User, serviceURLs *config.ServiceURLs) *HttpResponseUser {
	if user == nil {
		return &HttpResponseUser{}
	}

	resourcePath := fmt.Sprintf("/users/%d", user.ID)

	httpTickets := make([]HttpTicket, len(user.TicketList))
	for i, ticket := range user.TicketList {
		httpTickets[i] = HttpTicket{
			Code: ticket.Code,
		}
	}

	dto := &httpResponseUser{
		ID:               user.ID,
		Email:            user.Email,
		FirstName:        user.FirstName,
		LastName:         user.LastName,
		SocialMediaLinks: user.SocialMediaLinks,
		TicketList:       httpTickets,
		FirstNamePrivate: user.FirstNamePrivate,
		LastNamePrivate:  user.LastNamePrivate,
		Links: map[string]http.Link{
			"self":   hateoas.BuildSelfLink(serviceURLs.UserManager, resourcePath),
			"update": hateoas.BuildUpdateLink(serviceURLs.UserManager, resourcePath),
			"delete": hateoas.BuildDeleteLink(serviceURLs.UserManager, resourcePath),
			"events": hateoas.BuildRelatedLink(
				fmt.Sprintf("%s/events?owner_id=%d", serviceURLs.EventManager, user.ID),
				"events",
				"GET",
				"Get events owned by this user",
			),
			"packets": hateoas.BuildRelatedLink(
				fmt.Sprintf("%s/packets?owner_id=%d", serviceURLs.EventManager, user.ID),
				"packets",
				"GET",
				"Get packets owned by this user",
			),
			"create-ticket": hateoas.BuildRelatedLink(
				fmt.Sprintf("%s/clients/%d/tickets", serviceURLs.UserManager, user.ID),
				"create-ticket",
				"POST",
				"Create a ticket for this user",
			),
			"tickets": hateoas.BuildRelatedLink(
				fmt.Sprintf("%s/clients/%d/tickets", serviceURLs.UserManager, user.ID),
				"tickets",
				"GET",
				"View purchased tickets",
			),
		},
	}

	return &HttpResponseUser{
		User: dto,
	}
}

type HttpCreateUser struct {
	Email            string  `json:"email" binding:"required"`
	FirstName        string  `json:"first_name" binding:"required"`
	LastName         string  `json:"last_name" binding:"required"`
	SocialMediaLinks *string `json:"social_media_links"`
}

func (user *HttpCreateUser) ToUser() *domain.User {
	return &domain.User{
		Email:            user.Email,
		FirstName:        user.FirstName,
		LastName:         user.LastName,
		SocialMediaLinks: user.SocialMediaLinks,
		TicketList:       []domain.Ticket{},
	}
}

type HttpUpdateUser struct {
	Email            *string `json:"email"`
	FirstName        *string `json:"first_name"`
	LastName         *string `json:"last_name"`
	SocialMediaLinks *string `json:"social_media_links"`
}

func (user *HttpUpdateUser) ToUpdateMap() map[string]interface{} {
	updates := make(map[string]interface{})

	if user.Email != nil {
		updates["email"] = *user.Email
	}
	if user.FirstName != nil {
		updates["first_name"] = *user.FirstName
	}
	if user.LastName != nil {
		updates["last_name"] = *user.LastName
	}
	if user.SocialMediaLinks != nil {
		updates["social_media_links"] = *user.SocialMediaLinks
	}

	return updates
}

type HttpCreateTicketForUser struct {
	PacketID *int `json:"packet_id"`
	EventID  *int `json:"event_id"`
}

type HttpCreateTicketResponse struct {
	TicketCode string `json:"ticket_code"`
}
