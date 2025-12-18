package service

// EventManagerService provides event management operations
// This interface abstracts the Event Manager client implementation
type EventManagerService interface {
	// CreateTicket creates a ticket with the specified code
	// Returns the created ticket or a domain error
	CreateTicket(code string, packetID *int, eventID *int) (*TicketResponse, error)
}

// TicketResponse represents a ticket returned from the Event Manager
type TicketResponse struct {
	Code     string
	PacketID *int
	EventID  *int
}
