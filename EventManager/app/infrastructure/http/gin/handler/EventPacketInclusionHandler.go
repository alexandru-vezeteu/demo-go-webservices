package handler

import (
	"eventManager/controller"

	"github.com/gin-gonic/gin"
)

type GinEventPacketInclusionHandler struct {
	controller controller.IEventPacketInclusionController
}

func NewGinEventPacketInclusionHandler(controller controller.IEventPacketInclusionController) *GinEventPacketInclusionHandler {
	return &GinEventPacketInclusionHandler{controller: controller}
}

func (h *GinEventPacketInclusionHandler) CreateEventPacketInclusion(c *gin.Context) {
	// var req httpdto.HttpCreateEventPacket
	// if err := c.ShouldBindJSON(&req); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	// eventPacket := req.ToEventPacket()

	// ret, err := h.controller.CreateEventPacket(eventPacket)

	// if errors.Is(err, domain.NewEventValidationError("sda")) {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	// resp := httpdto.ToHttpResponseEventPacket(ret)
	// c.JSON(http.StatusCreated, resp)
}

func (h *GinEventPacketInclusionHandler) GetEventPacketsByEventID(c *gin.Context) {
	// idStr := c.Param("id")
	// id, err := strconv.ParseInt(idStr, 10, 64)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID format"})
	// 	return
	// }

	// ret, err := h.controller.GetEventPacketByID(int(id))
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// } else if ret == nil {
	// 	c.JSON(http.StatusNotFound, gin.H{"error": "No event found!"})
	// }
	// resp := httpdto.ToHttpResponseEventPacket(ret)

	// c.JSON(http.StatusOK, resp)

}

// func (h *GinEventPacketHandler) UpdateEventPacket(c *gin.Context) {
// 	var req httpdto.HttpUpdateEventPacket
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	eventPacket := req.ToEventPacket()
// 	eventPacket, err := h.controller.UpdateEventPacket(eventPacket)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "TODO": "TODO"})
// 		return
// 	}

// 	resp := httpdto.ToHttpResponseEventPacket(eventPacket)
// 	c.JSON(http.StatusNoContent, resp)
// }

func (h *GinEventPacketInclusionHandler) DeleteEventPacketInclusion(c *gin.Context) {
	// idStr := c.Param("id")
	// id, err := strconv.ParseInt(idStr, 10, 64)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID format"})
	// 	return
	// }

	// ret, err := h.controller.DeleteEventPacket(int(id))
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// } else if ret == nil {
	// 	c.JSON(http.StatusNotFound, gin.H{"error": "No event found!"})
	// }
	// resp := httpdto.ToHttpResponseEventPacket(ret)

	// c.JSON(http.StatusOK, resp)
}

func (h *GinEventPacketInclusionHandler) GetEventsByPacketID(c *gin.Context) {
	// idStr := c.Param("id")
	// id, err := strconv.ParseInt(idStr, 10, 64)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID format"})
	// 	return
	// }

	// ret, err := h.controller.DeleteEventPacket(int(id))
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// } else if ret == nil {
	// 	c.JSON(http.StatusNotFound, gin.H{"error": "No event found!"})
	// }
	// resp := httpdto.ToHttpResponseEventPacket(ret)

	// c.JSON(http.StatusOK, resp)
}
