package domain

import (
	"regexp"
	"strings"
)

type Ticket struct {
	PacketID *int
	EventID  *int
	Code     string
}

type User struct {
	ID               int
	Email            string
	FirstName        string
	LastName         string
	SocialMediaLinks *string
	TicketList       []Ticket

	FirstNamePrivate bool
	LastNamePrivate  bool
}

func validateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func validateSocialMediaLinks(links *string) bool {
	if links == nil {
		return true
	}

	trimmed := strings.TrimSpace(*links)
	return trimmed != ""
}

func validateTicketList(tickets []Ticket) bool {
	return true
}
