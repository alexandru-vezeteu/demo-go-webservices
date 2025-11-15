package model

import (
	"userService/application/domain"
)

type MongoUser struct {
	ID               int     `bson:"id"`
	Email            string  `bson:"email"`
	FirstName        string  `bson:"first_name"`
	LastName         string  `bson:"last_name"`
	SocialMediaLinks *string `bson:"social_media_links,omitempty"`
	TicketList       *string `bson:"ticket_list,omitempty"`
}

func (mu *MongoUser) ToDomain() *domain.User {
	return &domain.User{
		ID:               mu.ID,
		Email:            mu.Email,
		FirstName:        mu.FirstName,
		LastName:         mu.LastName,
		SocialMediaLinks: mu.SocialMediaLinks,
		TicketList:       mu.TicketList,
	}
}

func FromUser(u *domain.User) *MongoUser {
	return &MongoUser{
		ID:               u.ID,
		Email:            u.Email,
		FirstName:        u.FirstName,
		LastName:         u.LastName,
		SocialMediaLinks: u.SocialMediaLinks,
		TicketList:       u.TicketList,
	}
}
