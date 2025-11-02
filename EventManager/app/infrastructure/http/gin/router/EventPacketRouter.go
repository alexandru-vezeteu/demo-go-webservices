package router

import (
	"eventManager/infrastructure/http/gin/handler"

	"github.com/gin-gonic/gin"
)

func RegisterEventPacketRoutes(router *gin.RouterGroup, handler *handler.GinEventPacketHandler) {
	eventsGroup := router.Group("/event-packets")
	{
		eventsGroup.POST("/", handler.CreateEventPacket)

		eventsGroup.GET("/:id", handler.GetEventPacketByID)

		eventsGroup.PATCH("/:id", handler.UpdateEventPacket)

		eventsGroup.DELETE("/:id", handler.DeleteEventPacket)
	}
}
