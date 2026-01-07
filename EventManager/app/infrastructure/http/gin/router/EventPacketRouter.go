package router

import (
	"eventManager/infrastructure/http/gin/handler"

	"github.com/gin-gonic/gin"
)

func RegisterEventPacketRoutes(router *gin.RouterGroup, handler *handler.GinEventPacketHandler) {
	router.POST("/event-packets", handler.CreateEventPacket)
	router.POST("/event-packets/", handler.CreateEventPacket)

	router.GET("/event-packets", handler.FilterEventPackets)
	router.GET("/event-packets/", handler.FilterEventPackets)

	router.GET("/event-packets/:id", handler.GetEventPacketByID)

	router.PATCH("/event-packets/:id", handler.UpdateEventPacket)

	router.DELETE("/event-packets/:id", handler.DeleteEventPacket)
}
