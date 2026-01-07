package handler

import (
	"eventManager/application/usecase"
	"eventManager/infrastructure/http/config"
	"eventManager/infrastructure/http/gin/middleware"
	"eventManager/infrastructure/http/httpdto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GinTicketHandler struct {
	usecase     usecase.TicketUseCase
	serviceURLs *config.ServiceURLs
}

func NewGinTicketHandler(usecase usecase.TicketUseCase, serviceURLs *config.ServiceURLs) *GinTicketHandler {
	return &GinTicketHandler{
		usecase:     usecase,
		serviceURLs: serviceURLs,
	}
}

// CreateTicket godoc
// @Summary Create a new ticket
// @Description Create a new ticket for an event and packet
// @Tags tickets
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param ticket body httpdto.HttpCreateTicket true "Ticket details"
// @Success 201 {object} httpdto.HttpResponseTicket "Ticket created successfully"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 401 {object} map[string]string "Unauthorized - missing or invalid token"
// @Failure 404 {object} map[string]string "Event or packet not found"
// @Failure 409 {object} map[string]string "Ticket already exists"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /tickets [post]
func (h *GinTicketHandler) CreateTicket(c *gin.Context) {
	token, ok := requireAuth(c)
	if !ok {
		return
	}

	var req httpdto.HttpCreateTicket
	if err := middleware.StrictBindJSON(c, &req); err != nil {
		handleError(c, err)
		return
	}

	ticket := req.ToTicket()

	ret, err := h.usecase.CreateTicket(c.Request.Context(), token, ticket)
	if handleError(c, err) {
		return
	}

	resp := httpdto.ToHttpResponseTicket(ret, h.serviceURLs)
	c.JSON(http.StatusCreated, resp)
}

// PutTicket godoc
// @Summary Create or replace a ticket with specific code
// @Description Create a new ticket with a specified code or replace an existing one (idempotent)
// @Tags tickets
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param code path string true "Ticket code (UUID)"
// @Param ticket body httpdto.HttpCreateTicket true "Ticket details"
// @Success 200 {object} httpdto.HttpResponseTicket "Ticket created or updated successfully"
// @Failure 400 {object} map[string]string "Invalid request body or ticket code"
// @Failure 401 {object} map[string]string "Unauthorized - missing or invalid token"
// @Failure 404 {object} map[string]string "Event or packet not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /tickets/{code} [put]
func (h *GinTicketHandler) PutTicket(c *gin.Context) {
	token, ok := requireAuth(c)
	if !ok {
		return
	}

	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ticket code is required"})
		return
	}

	var req httpdto.HttpCreateTicket
	if err := middleware.StrictBindJSON(c, &req); err != nil {
		handleError(c, err)
		return
	}

	ticket := req.ToTicket()

	ret, err := h.usecase.PutTicket(c.Request.Context(), token, code, ticket)
	if handleError(c, err) {
		return
	}

	resp := httpdto.ToHttpResponseTicket(ret, h.serviceURLs)
	c.JSON(http.StatusOK, resp)
}

// GetTicketByCode godoc
// @Summary Get ticket by code
// @Description Retrieve a specific ticket by its unique code
// @Tags tickets
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param code path string true "Ticket code (UUID)"
// @Success 200 {object} httpdto.HttpResponseTicket "Ticket details"
// @Failure 400 {object} map[string]string "Invalid ticket code"
// @Failure 401 {object} map[string]string "Unauthorized - missing or invalid token"
// @Failure 404 {object} map[string]string "Ticket not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /tickets/{code} [get]
func (h *GinTicketHandler) GetTicketByCode(c *gin.Context) {
	token, ok := requireAuth(c)
	if !ok {
		return
	}

	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ticket code is required"})
		return
	}

	ret, err := h.usecase.GetTicketByCode(c.Request.Context(), token, code)
	if handleError(c, err) {
		return
	}

	resp := httpdto.ToHttpResponseTicket(ret, h.serviceURLs)
	c.JSON(http.StatusOK, resp)
}

// UpdateTicket godoc
// @Summary Update a ticket
// @Description Partially update ticket details (PATCH - only provided fields are updated)
// @Tags tickets
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param code path string true "Ticket code (UUID)"
// @Param ticket body httpdto.HttpUpdateTicket true "Fields to update"
// @Success 200 {object} httpdto.HttpResponseTicket "Ticket updated successfully"
// @Failure 400 {object} map[string]string "Invalid request body or ticket code"
// @Failure 401 {object} map[string]string "Unauthorized - missing or invalid token"
// @Failure 404 {object} map[string]string "Ticket not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /tickets/{code} [patch]
func (h *GinTicketHandler) UpdateTicket(c *gin.Context) {
	token, ok := requireAuth(c)
	if !ok {
		return
	}

	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ticket code is required"})
		return
	}

	var req httpdto.HttpUpdateTicket
	if err := middleware.StrictBindJSON(c, &req); err != nil {
		handleError(c, err)
		return
	}

	updates := req.ToUpdateMap()

	ticket, err := h.usecase.UpdateTicket(c.Request.Context(), token, code, updates)
	if handleError(c, err) {
		return
	}

	resp := httpdto.ToHttpResponseTicket(ticket, h.serviceURLs)
	c.JSON(http.StatusOK, resp)
}

// DeleteTicket godoc
// @Summary Delete a ticket
// @Description Delete a ticket by its code
// @Tags tickets
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param code path string true "Ticket code (UUID)"
// @Success 200 {object} httpdto.HttpResponseTicket "Ticket deleted successfully"
// @Failure 400 {object} map[string]string "Invalid ticket code"
// @Failure 401 {object} map[string]string "Unauthorized - missing or invalid token"
// @Failure 404 {object} map[string]string "Ticket not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /tickets/{code} [delete]
func (h *GinTicketHandler) DeleteTicket(c *gin.Context) {
	token, ok := requireAuth(c)
	if !ok {
		return
	}

	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ticket code is required"})
		return
	}

	ret, err := h.usecase.DeleteTicket(c.Request.Context(), token, code)
	if handleError(c, err) {
		return
	}

	resp := httpdto.ToHttpResponseTicket(ret, h.serviceURLs)
	c.JSON(http.StatusOK, resp)
}
