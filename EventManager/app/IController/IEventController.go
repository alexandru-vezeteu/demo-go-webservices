package IController

import "github.com/gin-gonic/gin"

type IEventController interface {
	GetEvent(c *gin.Context)
	PostEvent(c *gin.Context)
	//and other stuff
}
