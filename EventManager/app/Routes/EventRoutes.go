package Routes

import (
	"eventManager/IController"

	"github.com/gin-gonic/gin"
)

func AddEventRoutes(router *gin.RouterGroup, ec IController.IEventController) {
	var events = router.Group("events")
	events.GET("/:id", ec.GetEvent)
}
