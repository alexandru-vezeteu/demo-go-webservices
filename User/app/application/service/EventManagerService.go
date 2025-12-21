package service

type EventManagerService interface {
	CreateTicket(code string, packetID *int, eventID *int) (*TicketResponse, error)
}

type TicketResponse struct {
	Code     string
	PacketID *int
	EventID  *int
}
