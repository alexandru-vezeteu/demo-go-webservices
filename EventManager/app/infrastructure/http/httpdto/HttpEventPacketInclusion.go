package httpdto

import "eventManager/application/domain"

type HttpResponseEventPacketInclusion struct {
	PacketID int `json:"id_packet"`
	EventID  int `json:"id_event"`
}

func ToHttpResponseEventPacketInclusion(inclusion *domain.EventPacketInclusion) *HttpResponseEventPacketInclusion {
	return &HttpResponseEventPacketInclusion{
		PacketID: inclusion.PacketID,
		EventID:  inclusion.EventID,
	}
}

type HttpCreateEventPacketInclusion struct {
}

func (dto *HttpCreateEventPacketInclusion) ToEventPacketInclusion() *domain.EventPacketInclusion {
	return &domain.EventPacketInclusion{}
}

type HttpUpdateEventPacketInclusion struct {
}

func (dto *HttpUpdateEventPacketInclusion) ToUpdateMap() map[string]interface{} {
	updates := make(map[string]interface{})
	return updates
}
