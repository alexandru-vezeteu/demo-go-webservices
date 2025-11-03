package httpdto

import "eventManager/domain"

type HttpResponseEventPacketInclusion struct {
	PacketID int  `json:"id_packet"`
	EventID  int  `json:"id_event"`
	Seats    *int `json:"seats"`
}

func ToHttpResponseEventPacketInclusion(inclusion *domain.EventPacketInclusion) *HttpResponseEventPacketInclusion {
	return &HttpResponseEventPacketInclusion{
		PacketID: inclusion.PacketID,
		EventID:  inclusion.EventID,
		Seats:    inclusion.AllocatedSeats,
	}
}

type HttpCreateEventPacketInclusion struct {
	Seats *int `json:"seats"`
}

func (dto *HttpCreateEventPacketInclusion) ToEventPacketInclusion() *domain.EventPacketInclusion {
	return &domain.EventPacketInclusion{
		AllocatedSeats: dto.Seats,
	}
}

type HttpUpdateEventPacketInclusion struct {
	Seats *int `json:"seats"`
}

func (dto *HttpUpdateEventPacketInclusion) ToUpdateMap() map[string]interface{} {
	updates := make(map[string]interface{})
	if dto.Seats != nil {
		updates["allocated_seats"] = *dto.Seats
	}
	return updates
}
