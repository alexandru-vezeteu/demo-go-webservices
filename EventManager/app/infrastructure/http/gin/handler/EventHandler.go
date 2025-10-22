package handler

import (
	"eventManager/controller"
	"eventManager/infrastructure/http/httpdto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GinEventHandler struct {
	controller controller.IEventController
}

func NewGinEventHandler(controller controller.IEventController) *GinEventHandler {
	return &GinEventHandler{controller: controller}
}

func (h *GinEventHandler) CreateEvent(c *gin.Context) {

	var req httpdto.HttpEvent
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event := req.ToEvent()
	ret, err := h.controller.CreateEvent(event)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp := httpdto.ToHttpEvent(ret)

	c.JSON(http.StatusCreated, resp)
}

func (h *GinEventHandler) GetEventByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID format"})
		return
	}

	ret, err := h.controller.GetEventByID(int(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else if ret == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No event found!"})
	}
	resp := httpdto.ToHttpEvent(ret)

	c.JSON(http.StatusOK, resp)

}

func (h *GinEventHandler) UpdateEvent(c *gin.Context) {
	var req httpdto.HttpEvent
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event := req.ToEvent()
	event, err := h.controller.UpdateEvent(event)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "TODO": "TODO"})
		return
	}

	resp := httpdto.ToHttpEvent(event)
	c.JSON(http.StatusNoContent, resp)
}

func (h *GinEventHandler) DeleteEvent(c *gin.Context) {

}
