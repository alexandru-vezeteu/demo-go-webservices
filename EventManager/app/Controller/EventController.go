package Controller

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type EventController struct{}

func (ec *EventController) GetEvent(c *gin.Context) {
	fmt.Println("OK")

}
