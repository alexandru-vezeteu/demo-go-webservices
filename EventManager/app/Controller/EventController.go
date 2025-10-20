package Controller

import (
	"eventManager/Config"
	"eventManager/Model"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type EventController struct{}

func (ec *EventController) GetEvent(c *gin.Context) {
	id := c.Param("id")
	var event Model.Event

	if err := Config.DB.First(&event, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	c.JSON(http.StatusOK, event)

}

func (ec *EventController) PostEvent(c *gin.Context) {
	var event Model.Event
	fmt.Println("LOL")
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if event.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Event name is required"})
		return
	}
	if event.OwnerID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Owner ID is required"})
		return
	}

	if err := Config.DB.Create(&event).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, event)
}
