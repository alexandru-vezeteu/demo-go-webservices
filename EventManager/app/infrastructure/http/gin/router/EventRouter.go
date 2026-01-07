package router

import (
	"eventManager/infrastructure/http/gin/handler"

	"github.com/gin-gonic/gin"
)

func RegisterEventRoutes(router *gin.RouterGroup, handler *handler.GinEventHandler) {
	router.POST("/events", handler.CreateEvent)
	router.POST("/events/", handler.CreateEvent)

	router.GET("/events/:id", handler.GetEventByID)
	router.GET("/events", handler.FilterEvents)
	router.GET("/events/", handler.FilterEvents)

	router.PATCH("/events/:id", handler.UpdateEvent)

	router.DELETE("/events/:id", handler.DeleteEvent)
}
