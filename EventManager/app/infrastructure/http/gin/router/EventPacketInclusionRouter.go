package router

import (
	"eventManager/infrastructure/http/gin/handler"

	"github.com/gin-gonic/gin"
)

func RegisterEventPacketInclusionRoutes(router *gin.RouterGroup, handler *handler.GinEventPacketInclusionHandler) {

	router.POST("/events/event-packets", handler.CreateEventPacketInclusion)
	router.POST("/event-packets/events", handler.CreateEventPacketInclusion)

	router.GET("/events/:id/events-packets", handler.GetEventPacketsByEventID)
	router.GET("/events-packets/:id/events", handler.GetEventsByPacketID)

	//router.DELETE("/:id", handler.DeleteEventPacket)

}
