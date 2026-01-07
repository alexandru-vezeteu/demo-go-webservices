package httpdto

import (
	"fmt"
	"userService/application/domain"
	"userService/infrastructure/http"
	"userService/infrastructure/http/config"
	"userService/infrastructure/http/hateoas"
)

type HttpTicket struct {
	PacketID *int   `json:"packet_id,omitempty"`
	EventID  *int   `json:"event_id,omitempty"`
	Code     string `json:"code"`
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
			PacketID: ticket.PacketID,
			EventID:  ticket.EventID,
			Code:     ticket.Code,
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
	ID               int     `json:"id" binding:"required,min=1"`
	Email            string  `json:"email" binding:"required,email,min=3,max=255"`
	FirstName        string  `json:"first_name" binding:"omitempty,max=100"`
	LastName         string  `json:"last_name" binding:"omitempty,max=100"`
	SocialMediaLinks *string `json:"social_media_links" binding:"omitempty,max=500"`
}

func (user *HttpCreateUser) ToUser() *domain.User {
	return &domain.User{
		ID:               user.ID,
		Email:            user.Email,
		FirstName:        user.FirstName,
		LastName:         user.LastName,
		SocialMediaLinks: user.SocialMediaLinks,
		TicketList:       []domain.Ticket{},
	}
}

type HttpUpdateUser struct {
	Email            *string `json:"email" binding:"omitempty,email,min=3,max=255"`
	FirstName        *string `json:"first_name" binding:"omitempty,min=1,max=100"`
	LastName         *string `json:"last_name" binding:"omitempty,min=1,max=100"`
	SocialMediaLinks *string `json:"social_media_links" binding:"omitempty,max=500"`
	FirstNamePrivate *bool   `json:"first_name_private" binding:"omitempty"`
	LastNamePrivate  *bool   `json:"last_name_private" binding:"omitempty"`
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
	if user.FirstNamePrivate != nil {
		updates["first_name_private"] = *user.FirstNamePrivate
	}
	if user.LastNamePrivate != nil {
		updates["last_name_private"] = *user.LastNamePrivate
	}

	return updates
}

type HttpCreateTicketForUser struct {
	PacketID *int `json:"packet_id" binding:"omitempty,min=1"`
	EventID  *int `json:"event_id" binding:"omitempty,min=1"`
}

type HttpCreateTicketResponse struct {
	TicketCode string `json:"ticket_code"`
}

type HttpResponseUserList struct {
	Users []*httpResponseUser `json:"users"`
}

func ToHttpResponseUserList(users []*domain.User, serviceURLs *config.ServiceURLs) *HttpResponseUserList {
	if users == nil {
		users = []*domain.User{}
	}

	httpUsers := make([]*httpResponseUser, 0, len(users))

	for _, user := range users {
		resourcePath := fmt.Sprintf("/users/%d", user.ID)

		var httpTickets []HttpTicket
		if user.TicketList != nil {
			httpTickets = make([]HttpTicket, len(user.TicketList))
			for i, ticket := range user.TicketList {
				httpTickets[i] = HttpTicket{
					PacketID: ticket.PacketID,
					EventID:  ticket.EventID,
					Code:     ticket.Code,
				}
			}
		}

		httpUsers = append(httpUsers, &httpResponseUser{
			ID:               user.ID,
			Email:            user.Email,
			FirstName:        user.FirstName,
			LastName:         user.LastName,
			SocialMediaLinks: user.SocialMediaLinks,
			TicketList:       httpTickets,
			FirstNamePrivate: user.FirstNamePrivate,
			LastNamePrivate:  user.LastNamePrivate,
			Links: map[string]http.Link{
				"self": hateoas.BuildSelfLink(serviceURLs.UserManager, resourcePath),
			},
		})
	}

	return &HttpResponseUserList{
		Users: httpUsers,
	}
}
