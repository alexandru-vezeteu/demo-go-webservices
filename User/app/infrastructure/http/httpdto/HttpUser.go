package httpdto

import (
	"fmt"
	"userService/application/domain"
	"userService/infrastructure/http"
)

type httpResponseUser struct {
	ID               int                  `json:"id"`
	Email            string               `json:"email"`
	FirstName        string               `json:"first_name"`
	LastName         string               `json:"last_name"`
	SocialMediaLinks *string              `json:"social_media_links,omitempty"`
	TicketList       *string              `json:"ticket_list,omitempty"`
	Links            map[string]http.Link `json:"_links"`
}

type HttpResponseUser struct {
	User *httpResponseUser `json:"user"`
}

func ToHttpResponseUser(user *domain.User) *HttpResponseUser {
	if user == nil {
		return &HttpResponseUser{}
	}
	dto := &httpResponseUser{
		ID:               user.ID,
		Email:            user.Email,
		FirstName:        user.FirstName,
		LastName:         user.LastName,
		SocialMediaLinks: user.SocialMediaLinks,
		TicketList:       user.TicketList,
	}
	prefix := "/api/user-manager"
	dto.Links = map[string]http.Link{
		"self": {
			Href: fmt.Sprintf("%s/users/%d", prefix, user.ID),
			Type: "GET",
		},
		"update": {
			Href: fmt.Sprintf("%s/users/%d", prefix, user.ID),
			Type: "PATCH",
		},
		"delete": {
			Href: fmt.Sprintf("%s/users/%d", prefix, user.ID),
			Type: "DELETE",
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
	TicketList       *string `json:"ticket_list"`
}

func (user *HttpCreateUser) ToUser() *domain.User {
	return &domain.User{
		Email:            user.Email,
		FirstName:        user.FirstName,
		LastName:         user.LastName,
		SocialMediaLinks: user.SocialMediaLinks,
		TicketList:       user.TicketList,
	}
}

type HttpUpdateUser struct {
	Email            *string `json:"email"`
	FirstName        *string `json:"first_name"`
	LastName         *string `json:"last_name"`
	SocialMediaLinks *string `json:"social_media_links"`
	TicketList       *string `json:"ticket_list"`
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
	if user.TicketList != nil {
		updates["ticket_list"] = *user.TicketList
	}

	return updates
}
