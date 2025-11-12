package httpdto

import "eventManager/application/domain"

// HttpCreateTicket is the DTO for creating a ticket
type HttpCreateTicket struct {
	PacketID *int `json:"packet_id"`
	EventID  *int `json:"event_id"`
}

// HttpResponseTicket is the DTO for ticket responses
type HttpResponseTicket struct {
	Code     string `json:"code"`
	PacketID *int   `json:"packet_id"`
	EventID  *int   `json:"event_id"`
}

// HttpUpdateTicket is the DTO for updating a ticket
type HttpUpdateTicket struct {
	PacketID *int `json:"packet_id"`
	EventID  *int `json:"event_id"`
}

// ToTicket converts HttpCreateTicket to domain.Ticket
func (dto *HttpCreateTicket) ToTicket() *domain.Ticket {
	return &domain.Ticket{
		PacketID: dto.PacketID,
		EventID:  dto.EventID,
	}
}

// ToHttpResponseTicket converts domain.Ticket to HttpResponseTicket
func ToHttpResponseTicket(ticket *domain.Ticket) *HttpResponseTicket {
	return &HttpResponseTicket{
		Code:     ticket.Code,
		PacketID: ticket.PacketID,
		EventID:  ticket.EventID,
	}
}

// ToUpdateMap converts HttpUpdateTicket to a map for partial updates
func (dto *HttpUpdateTicket) ToUpdateMap() map[string]interface{} {
	updates := make(map[string]interface{})

	if dto.PacketID != nil {
		updates["packet_id"] = *dto.PacketID
	}

	if dto.EventID != nil {
		updates["event_id"] = *dto.EventID
	}

	return updates
}
