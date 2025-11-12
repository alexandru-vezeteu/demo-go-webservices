package router

import (
	"eventManager/infrastructure/http/gin/handler"

	"github.com/gin-gonic/gin"
)

func RegisterEventPacketInclusionRoutes(router *gin.RouterGroup, handler *handler.GinEventPacketInclusionHandler) {

	router.POST("/event-packet-inclusions/event/:event_id/packet/:packet_id", handler.CreateEventPacketInclusion)

	router.GET("/event-packet-inclusions/event/:event_id", handler.GetEventPacketsByEventID)
	router.GET("/event-packet-inclusions/packet/:packet_id", handler.GetEventsByPacketID)

	router.DELETE("/event-packet-inclusions/event/:event_id/packet/:packet_id", handler.DeleteEventPacketInclusion)

	router.PATCH("/event-packet-inclusions/event/:event_id/packet/:packet_id", handler.UpdateEventPacketInclusion)

}
