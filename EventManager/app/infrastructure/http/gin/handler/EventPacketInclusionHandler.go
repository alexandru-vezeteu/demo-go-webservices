package handler

import (
	"errors"
	"eventManager/controller"
	"eventManager/domain"
	"eventManager/infrastructure/http/httpdto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GinEventPacketInclusionHandler struct {
	controller controller.IEventPacketInclusionController
}

func NewGinEventPacketInclusionHandler(controller controller.IEventPacketInclusionController) *GinEventPacketInclusionHandler {
	return &GinEventPacketInclusionHandler{controller: controller}
}

func (h *GinEventPacketInclusionHandler) CreateEventPacketInclusion(c *gin.Context) {
	eventID, err := strconv.Atoi(c.Param("event_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event_id"})
		return
	}

	packetID, err := strconv.Atoi(c.Param("packet_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid packet_id"})
		return
	}

	var dto httpdto.HttpCreateEventPacketInclusion
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	inclusion := dto.ToEventPacketInclusion()
	inclusion.EventID = eventID
	inclusion.PacketID = packetID

	created, err := h.controller.CreateEventPacketInclusion(inclusion)
	if err != nil {
		var validationErr *domain.ValidationError
		var alreadyExistsErr *domain.AlreadyExistsError
		var foreignKeyErr *domain.ForeignKeyError
		var internalErr *domain.InternalError

		if errors.As(err, &validationErr) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else if errors.As(err, &alreadyExistsErr) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else if errors.As(err, &foreignKeyErr) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if errors.As(err, &internalErr) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, httpdto.ToHttpResponseEventPacketInclusion(created))
}

func (h *GinEventPacketInclusionHandler) GetEventPacketsByEventID(c *gin.Context) {
	eventID, err := strconv.Atoi(c.Param("event_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event_id"})
		return
	}

	packets, err := h.controller.GetEventPacketsByEventID(eventID)
	if err != nil {
		var notFoundErr *domain.NotFoundError
		var internalErr *domain.InternalError

		if errors.As(err, &notFoundErr) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if errors.As(err, &internalErr) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	response := make([]*httpdto.HttpResponseEventPacket, len(packets))
	for i, packet := range packets {
		response[i] = httpdto.ToHttpResponseEventPacket(packet)
	}

	c.JSON(http.StatusOK, response)
}
func (h *GinEventPacketInclusionHandler) GetEventsByPacketID(c *gin.Context) {
	packetID, err := strconv.Atoi(c.Param("packet_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid packet_id"})
		return
	}

	events, err := h.controller.GetEventsByPacketID(packetID)
	if err != nil {
		var notFoundErr *domain.NotFoundError
		var internalErr *domain.InternalError

		if errors.As(err, &notFoundErr) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if errors.As(err, &internalErr) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	response := make([]*httpdto.HttpResponseEvent, len(events))
	for i, event := range events {
		response[i] = httpdto.ToHttpResponseEvent(event)
	}

	c.JSON(http.StatusOK, response)
}
func (h *GinEventPacketInclusionHandler) UpdateEventPacketInclusion(c *gin.Context) {
	eventID, err := strconv.Atoi(c.Param("event_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event_id"})
		return
	}

	packetID, err := strconv.Atoi(c.Param("packet_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid packet_id"})
		return
	}

	var dto httpdto.HttpUpdateEventPacketInclusion
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := dto.ToUpdateMap()
	updated, err := h.controller.Update(eventID, packetID, updates)
	if err != nil {
		var notFoundErr *domain.NotFoundError
		var validationErr *domain.ValidationError
		var foreignKeyErr *domain.ForeignKeyError
		var internalErr *domain.InternalError

		if errors.As(err, &notFoundErr) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if errors.As(err, &validationErr) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else if errors.As(err, &foreignKeyErr) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if errors.As(err, &internalErr) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, httpdto.ToHttpResponseEventPacketInclusion(updated))
}

func (h *GinEventPacketInclusionHandler) DeleteEventPacketInclusion(c *gin.Context) {
	eventID, err := strconv.Atoi(c.Param("event_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event_id"})
		return
	}

	packetID, err := strconv.Atoi(c.Param("packet_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid packet_id"})
		return
	}

	deleted, err := h.controller.DeleteEventPacketInclusion(eventID, packetID)
	if err != nil {
		var notFoundErr *domain.NotFoundError
		var internalErr *domain.InternalError

		if errors.As(err, &notFoundErr) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if errors.As(err, &internalErr) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, httpdto.ToHttpResponseEventPacketInclusion(deleted))
}
