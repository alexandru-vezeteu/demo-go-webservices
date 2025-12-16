package httpdto

import "eventManager/application/domain"


type HttpCreateTicket struct {
	PacketID *int `json:"packet_id"`
	EventID  *int `json:"event_id"`
}


type HttpResponseTicket struct {
	Code     string `json:"code"`
	PacketID *int   `json:"packet_id"`
	EventID  *int   `json:"event_id"`
}


type HttpUpdateTicket struct {
	PacketID *int `json:"packet_id"`
	EventID  *int `json:"event_id"`
}


func (dto *HttpCreateTicket) ToTicket() *domain.Ticket {
	return &domain.Ticket{
		PacketID: dto.PacketID,
		EventID:  dto.EventID,
	}
}


func ToHttpResponseTicket(ticket *domain.Ticket) *HttpResponseTicket {
	return &HttpResponseTicket{
		Code:     ticket.Code,
		PacketID: ticket.PacketID,
		EventID:  ticket.EventID,
	}
}


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
