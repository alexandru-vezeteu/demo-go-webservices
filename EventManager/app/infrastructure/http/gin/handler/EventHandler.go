package handler

import (
	"eventManager/controller"
	"eventManager/infrastructure/http/httpdto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GinEventHandler struct {
	controller *controller.EventController
}

func NewGinEventHandler(controller *controller.EventController) *GinEventHandler {
	return &GinEventHandler{controller: controller}
}

func (h *GinEventHandler) CreateEvent(c *gin.Context) {
	var req httpdto.HttpEvent

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event, err := req.ToEvent()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	ret, err := h.controller.CreateEvent(event)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := httpdto.ToHttpEvent(ret)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *GinEventHandler) GetEventByID(c *gin.Context) {

}

func (h *GinEventHandler) UpdateEvent(c *gin.Context) {

}

func (h *GinEventHandler) DeleteEvent(c *gin.Context) {

}
