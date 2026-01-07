package router

import (
	"eventManager/infrastructure/http/gin/handler"

	"github.com/gin-gonic/gin"
)

func RegisterTicketRoutes(router *gin.RouterGroup, handler *handler.GinTicketHandler) {
	router.POST("/tickets", handler.CreateTicket)
	router.POST("/tickets/", handler.CreateTicket)

	router.GET("/tickets/:code", handler.GetTicketByCode)
	router.PATCH("/tickets/:code", handler.UpdateTicket)
	router.PUT("/tickets/:code", handler.PutTicket)
	router.DELETE("/tickets/:code", handler.DeleteTicket)
}
