package router

import (
	"eventManager/infrastructure/http/gin/handler"

	"github.com/gin-gonic/gin"
)

func RegisterEventPacketInclusionRoutes(router *gin.RouterGroup, handler *handler.GinEventPacketInclusionHandler) {

	router.POST("/events/:event_id/event-packets/:packet_id", handler.CreateEventPacketInclusion)
	router.POST("/event-packets/:packet_id/events/:event_id", handler.CreateEventPacketInclusion)

	router.GET("/events/:event_id/event-packets", handler.GetEventPacketsByEventID)
	router.GET("/event-packets/:packet_id/events", handler.GetEventsByPacketID)

	router.DELETE("/events/:event_id/event-packets/:packet_id", handler.DeleteEventPacketInclusion)
	router.DELETE("/event-packets/:packet_id/events/:event_id", handler.DeleteEventPacketInclusion)

	router.PATCH("/events/:event_id/event-packets/:packet_id", handler.UpdateEventPacketInclusion)
	router.PATCH("/event-packets/:packet_id/events/:event_id", handler.UpdateEventPacketInclusion)

}
