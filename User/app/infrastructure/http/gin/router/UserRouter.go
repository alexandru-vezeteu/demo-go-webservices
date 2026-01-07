package router

import (
	"userService/infrastructure/http/gin/handler"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(router *gin.RouterGroup, handler *handler.GinUserHandler) {
	router.POST("/users", handler.CreateUser)
	router.POST("/users/", handler.CreateUser)

	router.GET("/users/:id", handler.GetUserByID)

	router.PATCH("/users/:id", handler.UpdateUser)

	router.DELETE("/users/:id", handler.DeleteUser)

	router.GET("/events/:event_id/customers", handler.GetCustomersByEventID)
	router.GET("/packets/:packet_id/customers", handler.GetCustomersByPacketID)

	router.POST("/clients/:user_id/tickets", handler.CreateTicketForUser)
}
