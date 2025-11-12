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

type GinEventPacketHandler struct {
	controller controller.IEventPacketController
}

func NewGinEventPacketHandler(controller controller.IEventPacketController) *GinEventPacketHandler {
	return &GinEventPacketHandler{controller: controller}
}

// @Summary      Create a new event packet
// @Description  Creates a new event packet. Event packets can include multiple events and have allocated seats that must not exceed the minimum seats of included events.
// @Tags         event-packets
// @Accept       json
// @Produce      json
// @Param        event_packet body httpdto.HttpCreateEventPacket true "Event packet to create"
// @Success      201  {object}  httpdto.HttpResponseEventPacket  "Event packet created successfully"
// @Failure      400  {object}  map[string]interface{} "Invalid request format or validation error"
// @Failure      409  {object}  map[string]interface{} "Event packet already exists"
// @Failure      500  {object}  map[string]interface{} "Internal error"
// @Router       /event-packets [post]
func (h *GinEventPacketHandler) CreateEventPacket(c *gin.Context) {
	var req httpdto.HttpCreateEventPacket
	if err := middleware.StrictBindJSON(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	eventPacket := req.ToEventPacket()

	ret, err := h.controller.CreateEventPacket(eventPacket)

	var validationErr *domain.ValidationError
	if errors.As(err, &validationErr) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existsErr *domain.AlreadyExistsError
	if errors.As(err, &existsErr) {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	var internalErr *domain.InternalError
	if errors.As(err, &internalErr) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return
	}

	resp := httpdto.ToHttpResponseEventPacket(ret)
	c.JSON(http.StatusCreated, resp)
}

// @Summary      Get an event packet by ID
// @Description  Retrieves a single event packet using its unique integer ID
// @Tags         event-packets
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Event Packet ID"
// @Success      200  {object}  httpdto.HttpResponseEventPacket  "The requested event packet"
// @Failure      400  {object}  map[string]interface{} "Invalid event packet ID format or validation error"
// @Failure      404  {object}  map[string]interface{} "Event packet not found"
// @Failure      500  {object}  map[string]interface{} "Internal error"
// @Router       /event-packets/{id} [get]
func (h *GinEventPacketHandler) GetEventPacketByID(c *gin.Context) {
	id, err := middleware.ParseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ret, err := h.controller.GetEventPacketByID(id)

	var notFoundErr *domain.NotFoundError
	if errors.As(err, &notFoundErr) {
		c.JSON(http.StatusNotFound, gin.H{"error": notFoundErr.Error()})
		return
	}

	var validationErr *domain.ValidationError
	if errors.As(err, &validationErr) {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	resp := httpdto.ToHttpResponseEventPacket(ret)

	c.JSON(http.StatusOK, resp)

}

// @Summary      Update an event packet
// @Description  Partially updates an existing event packet by its ID. Validates that allocated_seats does not exceed minimum seats of included events.
// @Tags         event-packets
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Event Packet ID"
// @Param        updates body httpdto.HttpUpdateEventPacket true "Fields to update"
// @Success      200  {object}  httpdto.HttpResponseEventPacket "Event packet updated successfully"
// @Failure      400  {object}  map[string]interface{} "Invalid event packet ID format, request body, or validation error"
// @Failure      404  {object}  map[string]interface{} "Event packet not found"
// @Failure      409  {object}  map[string]interface{} "Name already taken"
// @Failure      500  {object}  map[string]interface{} "An unexpected error occurred"
// @Router       /event-packets/{id} [patch]
func (h *GinEventPacketHandler) UpdateEventPacket(c *gin.Context) {
	id, err := middleware.ParseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req httpdto.HttpUpdateEventPacket
	if err := middleware.StrictBindJSON(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := req.ToUpdateMap()

	event, err := h.controller.UpdateEventPacket(id, updates)

	var validationErr *domain.ValidationError
	var notFoundErr *domain.NotFoundError
	var uniqueName *domain.UniqueNameError
	if errors.As(err, &validationErr) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if errors.As(err, &notFoundErr) {
		c.JSON(http.StatusNotFound, gin.H{"error": notFoundErr.Error()})
		return
	}
	if errors.As(err, &uniqueName) {
		c.JSON(http.StatusConflict, gin.H{"error": uniqueName.Error() + " is already taken"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An unexpected error occurred"})
		return
	}

	resp := httpdto.ToHttpResponseEventPacket(event)
	c.JSON(http.StatusOK, resp)
}

// @Summary      Delete an event packet
// @Description  Deletes an event packet by its ID and returns the deleted event packet.
// @Tags         event-packets
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Event Packet ID"
// @Success      200  {object}  httpdto.HttpResponseEventPacket "Event packet deleted successfully"
// @Failure      400  {object}  map[string]interface{} "Invalid event packet ID format"
// @Failure      404  {object}  map[string]interface{} "Event packet not found"
// @Failure      500  {object}  map[string]interface{} "An unexpected error occurred"
// @Router       /event-packets/{id} [delete]
func (h *GinEventPacketHandler) DeleteEventPacket(c *gin.Context) {
	id, err := middleware.ParseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ret, err := h.controller.DeleteEventPacket(id)

	var notFoundErr *domain.NotFoundError
	if errors.As(err, &notFoundErr) {
		c.JSON(http.StatusNotFound, gin.H{"error": notFoundErr.Error()})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An unexpected error occurred"})
		return
	}

	resp := httpdto.ToHttpResponseEventPacket(ret)

	c.JSON(http.StatusOK, resp)
}
