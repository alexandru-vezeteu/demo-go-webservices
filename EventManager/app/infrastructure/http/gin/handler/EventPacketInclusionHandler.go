package handler

import (
	"eventManager/application/usecase"
	"eventManager/infrastructure/http/config"
	"eventManager/infrastructure/http/gin/middleware"
	"eventManager/infrastructure/http/httpdto"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GinEventPacketInclusionHandler struct {
	usecase     usecase.EventPacketInclusionUseCase
	serviceURLs *config.ServiceURLs
}

func NewGinEventPacketInclusionHandler(usecase usecase.EventPacketInclusionUseCase, serviceURLs *config.ServiceURLs) *GinEventPacketInclusionHandler {
	return &GinEventPacketInclusionHandler{
		usecase:     usecase,
		serviceURLs: serviceURLs,
	}
}

// CreateEventPacketInclusion godoc
// @Summary Create event packet inclusion
// @Description Link an event packet to an event (many-to-many relationship)
// @Tags event-packet-inclusions
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param event_id path int true "Event ID"
// @Param packet_id path int true "Packet ID"
// @Param inclusion body httpdto.HttpCreateEventPacketInclusion true "Inclusion details"
// @Success 201 {object} httpdto.HttpResponseEventPacketInclusion "Inclusion created successfully"
// @Failure 400 {object} map[string]string "Invalid parameters or request body"
// @Failure 401 {object} map[string]string "Unauthorized - missing or invalid token"
// @Failure 404 {object} map[string]string "Event or packet not found"
// @Failure 409 {object} map[string]string "Inclusion already exists"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /event-packet-inclusions/event/{event_id}/packet/{packet_id} [post]
func (h *GinEventPacketInclusionHandler) CreateEventPacketInclusion(c *gin.Context) {
	token, ok := requireAuth(c)
	if !ok {
		return
	}

	eventID, err := middleware.ParseIDParam(c, "event_id")
	if err != nil {
		handleError(c, err)
		return
	}

	packetID, err := middleware.ParseIDParam(c, "packet_id")
	if err != nil {
		handleError(c, err)
		return
	}

	var dto httpdto.HttpCreateEventPacketInclusion
	if err := middleware.StrictBindJSON(c, &dto); err != nil {
		handleError(c, err)
		return
	}

	inclusion := dto.ToEventPacketInclusion()
	inclusion.EventID = eventID
	inclusion.PacketID = packetID

	created, err := h.usecase.CreateEventPacketInclusion(c.Request.Context(), token, inclusion)
	if handleError(c, err) {
		return
	}

	c.JSON(http.StatusCreated, httpdto.ToHttpResponseEventPacketInclusion(created, h.serviceURLs))
}

// GetEventPacketsByEventID godoc
// @Summary Get all packet inclusions for an event
// @Description Retrieve all event packets linked to a specific event
// @Tags event-packet-inclusions
// @Accept json
// @Produce json
// @Param Authorization header string false "Bearer token (optional)"
// @Param event_id path int true "Event ID"
// @Success 200 {object} map[string]interface{} "List of event packets"
// @Failure 400 {object} map[string]string "Invalid event ID"
// @Failure 404 {object} map[string]string "Event not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /event-packet-inclusions/event/{event_id} [get]
func (h *GinEventPacketInclusionHandler) GetEventPacketsByEventID(c *gin.Context) {
	eventID, err := middleware.ParseIDParam(c, "event_id")
	if err != nil {
		handleError(c, err)
		return
	}

	token := getTokenFromHeader(c)
	packets, err := h.usecase.GetEventPacketsByEventID(c.Request.Context(), token, eventID)
	if handleError(c, err) {
		return
	}

	selfPath := fmt.Sprintf("/events/%d/packets", eventID)
	response := httpdto.ToHttpResponseEventPacketList(packets, selfPath, h.serviceURLs)

	c.JSON(http.StatusOK, response)
}

// GetEventsByPacketID godoc
// @Summary Get all event inclusions for a packet
// @Description Retrieve all events linked to a specific event packet
// @Tags event-packet-inclusions
// @Accept json
// @Produce json
// @Param Authorization header string false "Bearer token (optional)"
// @Param packet_id path int true "Packet ID"
// @Success 200 {object} map[string]interface{} "List of events"
// @Failure 400 {object} map[string]string "Invalid packet ID"
// @Failure 404 {object} map[string]string "Packet not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /event-packet-inclusions/packet/{packet_id} [get]
func (h *GinEventPacketInclusionHandler) GetEventsByPacketID(c *gin.Context) {
	packetID, err := middleware.ParseIDParam(c, "packet_id")
	if err != nil {
		handleError(c, err)
		return
	}

	token := getTokenFromHeader(c)
	events, err := h.usecase.GetEventsByPacketID(c.Request.Context(), token, packetID)
	if handleError(c, err) {
		return
	}

	selfPath := fmt.Sprintf("/packets/%d/events", packetID)
	response := httpdto.ToHttpResponseEventListCustom(events, selfPath, h.serviceURLs)

	c.JSON(http.StatusOK, response)
}

// UpdateEventPacketInclusion godoc
// @Summary Update event packet inclusion
// @Description Update the details of an event-packet inclusion relationship
// @Tags event-packet-inclusions
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param event_id path int true "Event ID"
// @Param packet_id path int true "Packet ID"
// @Param inclusion body httpdto.HttpUpdateEventPacketInclusion true "Fields to update"
// @Success 200 {object} httpdto.HttpResponseEventPacketInclusion "Inclusion updated successfully"
// @Failure 400 {object} map[string]string "Invalid parameters or request body"
// @Failure 401 {object} map[string]string "Unauthorized - missing or invalid token"
// @Failure 404 {object} map[string]string "Inclusion not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /event-packet-inclusions/event/{event_id}/packet/{packet_id} [patch]
func (h *GinEventPacketInclusionHandler) UpdateEventPacketInclusion(c *gin.Context) {
	token, ok := requireAuth(c)
	if !ok {
		return
	}

	eventID, err := middleware.ParseIDParam(c, "event_id")
	if err != nil {
		handleError(c, err)
		return
	}

	packetID, err := middleware.ParseIDParam(c, "packet_id")
	if err != nil {
		handleError(c, err)
		return
	}

	var dto httpdto.HttpUpdateEventPacketInclusion
	if err := middleware.StrictBindJSON(c, &dto); err != nil {
		handleError(c, err)
		return
	}

	updates := dto.ToUpdateMap()
	updated, err := h.usecase.Update(c.Request.Context(), token, eventID, packetID, updates)
	if handleError(c, err) {
		return
	}

	c.JSON(http.StatusOK, httpdto.ToHttpResponseEventPacketInclusion(updated, h.serviceURLs))
}

// DeleteEventPacketInclusion godoc
// @Summary Delete event packet inclusion
// @Description Remove the link between an event and an event packet
// @Tags event-packet-inclusions
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param event_id path int true "Event ID"
// @Param packet_id path int true "Packet ID"
// @Success 200 {object} httpdto.HttpResponseEventPacketInclusion "Inclusion deleted successfully"
// @Failure 400 {object} map[string]string "Invalid parameters"
// @Failure 401 {object} map[string]string "Unauthorized - missing or invalid token"
// @Failure 404 {object} map[string]string "Inclusion not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /event-packet-inclusions/event/{event_id}/packet/{packet_id} [delete]
func (h *GinEventPacketInclusionHandler) DeleteEventPacketInclusion(c *gin.Context) {
	token, ok := requireAuth(c)
	if !ok {
		return
	}

	eventID, err := middleware.ParseIDParam(c, "event_id")
	if err != nil {
		handleError(c, err)
		return
	}

	packetID, err := middleware.ParseIDParam(c, "packet_id")
	if err != nil {
		handleError(c, err)
		return
	}

	deleted, err := h.usecase.DeleteEventPacketInclusion(c.Request.Context(), token, eventID, packetID)
	if handleError(c, err) {
		return
	}

	c.JSON(http.StatusOK, httpdto.ToHttpResponseEventPacketInclusion(deleted, h.serviceURLs))
}
