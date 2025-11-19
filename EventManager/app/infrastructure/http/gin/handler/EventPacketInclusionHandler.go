package handler

import (
	"errors"
	"eventManager/application/controller"
	"eventManager/application/domain"
	"eventManager/infrastructure/http/gin/middleware"
	"eventManager/infrastructure/http/httpdto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GinEventPacketInclusionHandler struct {
	controller controller.IEventPacketInclusionController
}

func NewGinEventPacketInclusionHandler(controller controller.IEventPacketInclusionController) *GinEventPacketInclusionHandler {
	return &GinEventPacketInclusionHandler{controller: controller}
}

// @Summary      Create an event packet inclusion
// @Description  Adds an event to an event packet. Validates that the event has enough seats to meet the packet's allocated_seats requirement.
// @Tags         event-packet-inclusions
// @Accept       json
// @Produce      json
// @Param        event_id   path      int  true  "Event ID"
// @Param        packet_id  path      int  true  "Event Packet ID"
// @Param        inclusion body httpdto.HttpCreateEventPacketInclusion true "Inclusion data"
// @Success      201  {object}  httpdto.HttpResponseEventPacketInclusion  "Inclusion created successfully"
// @Failure      400  {object}  map[string]interface{} "Invalid request format or validation error"
// @Failure      404  {object}  map[string]interface{} "Event or packet not found (foreign key error)"
// @Failure      409  {object}  map[string]interface{} "Inclusion already exists"
// @Failure      500  {object}  map[string]interface{} "Internal error"
// @Router       /event-packet-inclusions/event/{event_id}/packet/{packet_id} [post]
func (h *GinEventPacketInclusionHandler) CreateEventPacketInclusion(c *gin.Context) {
	eventID, err := middleware.ParseIDParam(c, "event_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	packetID, err := middleware.ParseIDParam(c, "packet_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var dto httpdto.HttpCreateEventPacketInclusion
	if err := middleware.StrictBindJSON(c, &dto); err != nil {
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, httpdto.ToHttpResponseEventPacketInclusion(created))
}

// @Summary      Get event packets by event ID
// @Description  Retrieves all event packets that include a specific event
// @Tags         event-packet-inclusions
// @Accept       json
// @Produce      json
// @Param        event_id   path      int  true  "Event ID"
// @Success      200  {array}   httpdto.HttpResponseEventPacket  "List of event packets"
// @Failure      400  {object}  map[string]interface{} "Invalid event ID format"
// @Failure      404  {object}  map[string]interface{} "Not found"
// @Failure      500  {object}  map[string]interface{} "Internal error"
// @Router       /event-packet-inclusions/event/{event_id} [get]
func (h *GinEventPacketInclusionHandler) GetEventPacketsByEventID(c *gin.Context) {
	eventID, err := middleware.ParseIDParam(c, "event_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	response := make([]*httpdto.HttpResponseEventPacket, len(packets))
	for i, packet := range packets {
		response[i] = httpdto.ToHttpResponseEventPacket(packet)
	}

	c.JSON(http.StatusOK, response)
}

// @Summary      Get events by packet ID
// @Description  Retrieves all events that are included in a specific event packet
// @Tags         event-packet-inclusions
// @Accept       json
// @Produce      json
// @Param        packet_id   path      int  true  "Event Packet ID"
// @Success      200  {array}   httpdto.HttpResponseEvent  "List of events"
// @Failure      400  {object}  map[string]interface{} "Invalid packet ID format"
// @Failure      404  {object}  map[string]interface{} "Not found"
// @Failure      500  {object}  map[string]interface{} "Internal error"
// @Router       /event-packet-inclusions/packet/{packet_id} [get]
func (h *GinEventPacketInclusionHandler) GetEventsByPacketID(c *gin.Context) {
	packetID, err := middleware.ParseIDParam(c, "packet_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	response := make([]*httpdto.HttpResponseEvent, len(events))
	for i, event := range events {
		response[i] = httpdto.ToHttpResponseEvent(event)
	}

	c.JSON(http.StatusOK, response)
}

// @Summary      Update an event packet inclusion
// @Description  Updates an existing inclusion relationship between an event and a packet
// @Tags         event-packet-inclusions
// @Accept       json
// @Produce      json
// @Param        event_id   path      int  true  "Event ID"
// @Param        packet_id  path      int  true  "Event Packet ID"
// @Param        updates body httpdto.HttpUpdateEventPacketInclusion true "Fields to update"
// @Success      200  {object}  httpdto.HttpResponseEventPacketInclusion  "Inclusion updated successfully"
// @Failure      400  {object}  map[string]interface{} "Invalid IDs, request body, or validation error"
// @Failure      404  {object}  map[string]interface{} "Inclusion not found or foreign key error"
// @Failure      500  {object}  map[string]interface{} "Internal error"
// @Router       /event-packet-inclusions/event/{event_id}/packet/{packet_id} [patch]
func (h *GinEventPacketInclusionHandler) UpdateEventPacketInclusion(c *gin.Context) {
	eventID, err := middleware.ParseIDParam(c, "event_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	packetID, err := middleware.ParseIDParam(c, "packet_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var dto httpdto.HttpUpdateEventPacketInclusion
	if err := middleware.StrictBindJSON(c, &dto); err != nil {
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, httpdto.ToHttpResponseEventPacketInclusion(updated))
}

// @Summary      Delete an event packet inclusion
// @Description  Removes an event from an event packet by deleting the inclusion relationship
// @Tags         event-packet-inclusions
// @Accept       json
// @Produce      json
// @Param        event_id   path      int  true  "Event ID"
// @Param        packet_id  path      int  true  "Event Packet ID"
// @Success      200  {object}  httpdto.HttpResponseEventPacketInclusion  "Inclusion deleted successfully"
// @Failure      400  {object}  map[string]interface{} "Invalid event or packet ID format"
// @Failure      404  {object}  map[string]interface{} "Inclusion not found"
// @Failure      500  {object}  map[string]interface{} "Internal error"
// @Router       /event-packet-inclusions/event/{event_id}/packet/{packet_id} [delete]
func (h *GinEventPacketInclusionHandler) DeleteEventPacketInclusion(c *gin.Context) {
	eventID, err := middleware.ParseIDParam(c, "event_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	packetID, err := middleware.ParseIDParam(c, "packet_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, httpdto.ToHttpResponseEventPacketInclusion(deleted))
}
