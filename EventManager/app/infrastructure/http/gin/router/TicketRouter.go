package router

import (
	"eventManager/infrastructure/http/gin/handler"

	"github.com/gin-gonic/gin"
)

func RegisterTicketRoutes(router *gin.RouterGroup, handler *handler.GinTicketHandler) {
	ticketsGroup := router.Group("/tickets")
	{
		ticketsGroup.POST("/", handler.CreateTicket)
		ticketsGroup.GET("/:code", handler.GetTicketByCode)
		ticketsGroup.PATCH("/:code", handler.UpdateTicket)
		ticketsGroup.PUT("/:code", handler.ReplaceTicket)
		ticketsGroup.DELETE("/:code", handler.DeleteTicket)
	}
}
