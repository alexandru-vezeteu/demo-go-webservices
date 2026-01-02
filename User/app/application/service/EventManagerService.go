package service

import "context"

type EventManagerService interface {
	CreateTicket(ctx context.Context, code string, packetID *int, eventID *int) (*TicketResponse, error)
}

type TicketResponse struct {
	Code     string
	PacketID *int
	EventID  *int
}
