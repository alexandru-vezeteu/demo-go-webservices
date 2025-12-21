package httpdto

import (
	"eventManager/application/domain"
	"eventManager/infrastructure/http/config"
	"eventManager/infrastructure/http/hateoas"
	"fmt"
)

type HttpResponseEventPacketInclusion struct {
	PacketID int                     `json:"id_packet"`
	EventID  int                     `json:"id_event"`
	Links    map[string]hateoas.Link `json:"_links"`
}

func ToHttpResponseEventPacketInclusion(inclusion *domain.EventPacketInclusion, serviceURLs *config.ServiceURLs) *HttpResponseEventPacketInclusion {
	resourcePath := fmt.Sprintf("/packets/%d/events/%d", inclusion.PacketID, inclusion.EventID)

	return &HttpResponseEventPacketInclusion{
		PacketID: inclusion.PacketID,
		EventID:  inclusion.EventID,
		Links: map[string]hateoas.Link{
			"self": hateoas.BuildSelfLink(serviceURLs.EventManager, resourcePath),
			"delete": hateoas.BuildDeleteLink(serviceURLs.EventManager, resourcePath),
			"event": hateoas.BuildRelatedLink(
				fmt.Sprintf("%s/events/%d", serviceURLs.EventManager, inclusion.EventID),
				"event",
				"GET",
				"Get the event in this inclusion",
			),
			"packet": hateoas.BuildRelatedLink(
				fmt.Sprintf("%s/packets/%d", serviceURLs.EventManager, inclusion.PacketID),
				"packet",
				"GET",
				"Get the packet in this inclusion",
			),
		},
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
