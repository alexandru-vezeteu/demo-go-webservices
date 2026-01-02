package model

import (
	"userService/application/domain"
)

type MongoTicket struct {
	Code string `bson:"code"`
}

type MongoUser struct {
	ID               int           `bson:"id"`
	Email            string        `bson:"email"`
	FirstName        string        `bson:"first_name"`
	LastName         string        `bson:"last_name"`
	SocialMediaLinks *string       `bson:"social_media_links,omitempty"`
	TicketList       []MongoTicket `bson:"ticket_list,omitempty"`

	FirstNamePrivate bool `bson:"first_name_private"`
	LastNamePrivate  bool `bson:"last_name_private"`
}

func (mu *MongoUser) ToDomain() *domain.User {
	domainTickets := make([]domain.Ticket, len(mu.TicketList))
	for i, mt := range mu.TicketList {
		domainTickets[i] = domain.Ticket{
			Code: mt.Code,
		}
	}

	return &domain.User{
		ID:               mu.ID,
		Email:            mu.Email,
		FirstName:        mu.FirstName,
		LastName:         mu.LastName,
		SocialMediaLinks: mu.SocialMediaLinks,
		TicketList:       domainTickets,
		FirstNamePrivate: mu.FirstNamePrivate,
		LastNamePrivate:  mu.LastNamePrivate,
	}
}

func FromUser(u *domain.User) *MongoUser {
	mongoTickets := make([]MongoTicket, len(u.TicketList))
	for i, dt := range u.TicketList {
		mongoTickets[i] = MongoTicket{
			Code: dt.Code,
		}
	}

	return &MongoUser{
		ID:               u.ID,
		Email:            u.Email,
		FirstName:        u.FirstName,
		LastName:         u.LastName,
		SocialMediaLinks: u.SocialMediaLinks,
		TicketList:       mongoTickets,
		FirstNamePrivate: u.FirstNamePrivate,
		LastNamePrivate:  u.LastNamePrivate,
	}
}
