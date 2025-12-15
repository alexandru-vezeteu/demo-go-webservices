package router

import (
	"userService/infrastructure/http/gin/handler"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(router *gin.RouterGroup, handler *handler.GinUserHandler) {
	usersGroup := router.Group("/users")
	{
		usersGroup.POST("/", handler.CreateUser)

		usersGroup.GET("/:id", handler.GetUserByID)

		usersGroup.PATCH("/:id", handler.UpdateUser)

		usersGroup.DELETE("/:id", handler.DeleteUser)
	}

	clientsGroup := router.Group("/clients")
	{
		clientsGroup.POST("/:user_id/tickets", handler.CreateTicketForUser)
	}
}
