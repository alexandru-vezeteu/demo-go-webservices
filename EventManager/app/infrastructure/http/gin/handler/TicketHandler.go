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

type GinTicketHandler struct {
	controller controller.ITicketController
}

func NewGinTicketHandler(controller controller.ITicketController) *GinTicketHandler {
	return &GinTicketHandler{controller: controller}
}

// @Summary      Create a new ticket
// @Description  Creates a new ticket for either an event or an event packet. The ticket code is automatically generated as a UUID. Validates seat availability before creation.
// @Tags         tickets
// @Accept       json
// @Produce      json
// @Param        ticket body httpdto.HttpCreateTicket true "Ticket to create (must specify either event_id or packet_id, not both)"
// @Success      201  {object}  httpdto.HttpResponseTicket  "Ticket created successfully with auto-generated code"
// @Failure      400  {object}  map[string]interface{} "Invalid request format, validation error, or no seats available"
// @Failure      409  {object}  map[string]interface{} "Ticket already exists"
// @Failure      500  {object}  map[string]interface{} "Internal error or unexpected error occurred"
// @Router       /tickets [post]
func (h *GinTicketHandler) CreateTicket(c *gin.Context) {
	var req httpdto.HttpCreateTicket
	if err := middleware.StrictBindJSON(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ticket := req.ToTicket()

	ret, err := h.controller.CreateTicket(ticket)

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

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An unexpected error occurred"})
		return
	}

	resp := httpdto.ToHttpResponseTicket(ret)
	c.JSON(http.StatusCreated, resp)
}

// @Summary      Get a ticket by code
// @Description  Retrieves a single ticket using its unique UUID code
// @Tags         tickets
// @Accept       json
// @Produce      json
// @Param        code   path      string  true  "Ticket Code (UUID)"
// @Success      200  {object}  httpdto.HttpResponseTicket  "The requested ticket"
// @Failure      400  {object}  map[string]interface{} "Invalid ticket code format or validation error"
// @Failure      404  {object}  map[string]interface{} "Ticket not found"
// @Failure      500  {object}  map[string]interface{} "Internal error"
// @Router       /tickets/{code} [get]
func (h *GinTicketHandler) GetTicketByCode(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ticket code is required"})
		return
	}

	ret, err := h.controller.GetTicketByCode(code)

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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return
	}

	resp := httpdto.ToHttpResponseTicket(ret)
	c.JSON(http.StatusOK, resp)
}

// @Summary      Update a ticket
// @Description  Partially updates an existing ticket by its code. Validates seat availability when reassigning to a different event or packet.
// @Tags         tickets
// @Accept       json
// @Produce      json
// @Param        code   path      string  true  "Ticket Code (UUID)"
// @Param        updates body httpdto.HttpUpdateTicket true "Fields to update (event_id or packet_id)"
// @Success      200  {object}  httpdto.HttpResponseTicket "Ticket updated successfully"
// @Failure      400  {object}  map[string]interface{} "Invalid ticket code format, request body, or validation error (e.g., no seats available)"
// @Failure      404  {object}  map[string]interface{} "Ticket not found"
// @Failure      500  {object}  map[string]interface{} "An unexpected error occurred"
// @Router       /tickets/{code} [patch]
func (h *GinTicketHandler) UpdateTicket(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ticket code is required"})
		return
	}

	var req httpdto.HttpUpdateTicket
	if err := middleware.StrictBindJSON(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := req.ToUpdateMap()

	ticket, err := h.controller.UpdateTicket(code, updates)

	var validationErr *domain.ValidationError
	var notFoundErr *domain.NotFoundError
	if errors.As(err, &validationErr) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if errors.As(err, &notFoundErr) {
		c.JSON(http.StatusNotFound, gin.H{"error": notFoundErr.Error()})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An unexpected error occurred"})
		return
	}

	resp := httpdto.ToHttpResponseTicket(ticket)
	c.JSON(http.StatusOK, resp)
}

// @Summary      Delete a ticket
// @Description  Deletes a ticket by its code and returns the deleted ticket.
// @Tags         tickets
// @Accept       json
// @Produce      json
// @Param        code   path      string  true  "Ticket Code (UUID)"
// @Success      200  {object}  httpdto.HttpResponseTicket "Ticket deleted successfully"
// @Failure      400  {object}  map[string]interface{} "Invalid ticket code format"
// @Failure      404  {object}  map[string]interface{} "Ticket not found"
// @Failure      500  {object}  map[string]interface{} "An unexpected error occurred"
// @Router       /tickets/{code} [delete]
func (h *GinTicketHandler) DeleteTicket(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ticket code is required"})
		return
	}

	ret, err := h.controller.DeleteTicket(code)

	var notFoundErr *domain.NotFoundError
	if errors.As(err, &notFoundErr) {
		c.JSON(http.StatusNotFound, gin.H{"error": notFoundErr.Error()})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An unexpected error occurred"})
		return
	}

	resp := httpdto.ToHttpResponseTicket(ret)
	c.JSON(http.StatusOK, resp)
}
