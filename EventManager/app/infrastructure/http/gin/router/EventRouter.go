package router

import (
	"eventManager/infrastructure/http/gin/handler"

	"github.com/gin-gonic/gin"
)

func RegisterEventRoutes(router *gin.RouterGroup, handler *handler.GinEventHandler) {
	eventsGroup := router.Group("/events")
	{
		eventsGroup.POST("/", handler.CreateEvent)

		eventsGroup.GET("/:id", handler.GetEventByID)

		eventsGroup.PATCH("/:id", handler.UpdateEvent)

		eventsGroup.DELETE("/:id", handler.DeleteEvent)
	}
}
