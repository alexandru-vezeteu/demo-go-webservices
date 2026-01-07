package handler

import (
	"eventManager/application/domain"
	"eventManager/application/repository"
	"eventManager/application/usecase"
	"eventManager/infrastructure/http/config"
	"eventManager/infrastructure/http/gin/middleware"
	"eventManager/infrastructure/http/httpdto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GinEventPacketHandler struct {
	usecase     usecase.EventPacketUseCase
	repo        repository.EventPacketRepository
	serviceURLs *config.ServiceURLs
}

func NewGinEventPacketHandler(usecase usecase.EventPacketUseCase, repo repository.EventPacketRepository, serviceURLs *config.ServiceURLs) *GinEventPacketHandler {
	return &GinEventPacketHandler{
		usecase:     usecase,
		repo:        repo,
		serviceURLs: serviceURLs,
	}
}

// CreateEventPacket godoc
// @Summary Create a new event packet
// @Description Create a new event packet (ticket package/bundle)
// @Tags event-packets
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param packet body httpdto.HttpCreateEventPacket true "Event packet details"
// @Success 201 {object} httpdto.HttpResponseEventPacket "Event packet created successfully"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 401 {object} map[string]string "Unauthorized - missing or invalid token"
// @Failure 409 {object} map[string]string "Event packet already exists"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /event-packets [post]
func (h *GinEventPacketHandler) CreateEventPacket(c *gin.Context) {
	token, ok := requireAuth(c)
	if !ok {
		return
	}

	var req httpdto.HttpCreateEventPacket
	if err := middleware.StrictBindJSON(c, &req); err != nil {
		handleError(c, err)
		return
	}

	eventPacket := req.ToEventPacket()

	ret, err := h.usecase.CreateEventPacket(c.Request.Context(), token, eventPacket)
	if handleError(c, err) {
		return
	}

	resp := httpdto.ToHttpResponseEventPacket(ret, h.serviceURLs)
	c.JSON(http.StatusCreated, resp)
}

// GetEventPacketByID godoc
// @Summary Get event packet by ID
// @Description Retrieve a specific event packet by its unique identifier
// @Tags event-packets
// @Accept json
// @Produce json
// @Param id path string true "Event Packet ID (UUID)"
// @Param Authorization header string false "Bearer token (optional)"
// @Success 200 {object} httpdto.HttpResponseEventPacket "Event packet details"
// @Failure 400 {object} map[string]string "Invalid event packet ID format"
// @Failure 404 {object} map[string]string "Event packet not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /event-packets/{id} [get]
func (h *GinEventPacketHandler) GetEventPacketByID(c *gin.Context) {
	id, err := middleware.ParseIDParam(c, "id")
	if err != nil {
		handleError(c, err)
		return
	}

	token := getTokenFromHeader(c)
	ret, err := h.usecase.GetEventPacketByID(c.Request.Context(), token, id)
	if handleError(c, err) {
		return
	}

	resp := httpdto.ToHttpResponseEventPacket(ret, h.serviceURLs)
	c.JSON(http.StatusOK, resp)
}

// UpdateEventPacket godoc
// @Summary Update an existing event packet
// @Description Partially update event packet details (PATCH - only provided fields are updated)
// @Tags event-packets
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Event Packet ID (UUID)"
// @Param packet body httpdto.HttpUpdateEventPacket true "Fields to update"
// @Success 200 {object} httpdto.HttpResponseEventPacket "Event packet updated successfully"
// @Failure 400 {object} map[string]string "Invalid request body or event packet ID"
// @Failure 401 {object} map[string]string "Unauthorized - missing or invalid token"
// @Failure 404 {object} map[string]string "Event packet not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /event-packets/{id} [patch]
func (h *GinEventPacketHandler) UpdateEventPacket(c *gin.Context) {
	token, ok := requireAuth(c)
	if !ok {
		return
	}

	id, err := middleware.ParseIDParam(c, "id")
	if err != nil {
		handleError(c, err)
		return
	}

	var req httpdto.HttpUpdateEventPacket
	if err := middleware.StrictBindJSON(c, &req); err != nil {
		handleError(c, err)
		return
	}

	updates := req.ToUpdateMap()

	event, err := h.usecase.UpdateEventPacket(c.Request.Context(), token, id, updates)
	if handleError(c, err) {
		return
	}

	resp := httpdto.ToHttpResponseEventPacket(event, h.serviceURLs)
	c.JSON(http.StatusOK, resp)
}

// DeleteEventPacket godoc
// @Summary Delete an event packet
// @Description Delete an event packet by its ID
// @Tags event-packets
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Event Packet ID (UUID)"
// @Success 200 {object} httpdto.HttpResponseEventPacket "Event packet deleted successfully"
// @Failure 400 {object} map[string]string "Invalid event packet ID format"
// @Failure 401 {object} map[string]string "Unauthorized - missing or invalid token"
// @Failure 404 {object} map[string]string "Event packet not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /event-packets/{id} [delete]
func (h *GinEventPacketHandler) DeleteEventPacket(c *gin.Context) {
	token, ok := requireAuth(c)
	if !ok {
		return
	}

	id, err := middleware.ParseIDParam(c, "id")
	if err != nil {
		handleError(c, err)
		return
	}

	ret, err := h.usecase.DeleteEventPacket(c.Request.Context(), token, id)
	if handleError(c, err) {
		return
	}

	resp := httpdto.ToHttpResponseEventPacket(ret, h.serviceURLs)
	c.JSON(http.StatusOK, resp)
}

// FilterEventPackets godoc
// @Summary List and filter event packets
// @Description Get a paginated list of event packets with optional filters. No authentication required for public access.
// @Tags event-packets
// @Accept json
// @Produce json
// @Param Authorization header string false "Bearer token (optional)"
// @Param name query string false "Filter by packet name (partial match)"
// @Param location query string false "Filter by location (partial match)"
// @Param description query string false "Filter by description (partial match)"
// @Param min_seats query int false "Minimum allocated seats"
// @Param max_seats query int false "Maximum allocated seats"
// @Param page query int false "Page number (default: 1)"
// @Param per_page query int false "Items per page (default: 10, max: 100)"
// @Param order_by query string false "Sort field (e.g., 'name_asc', 'seats_desc')"
// @Success 200 {object} httpdto.HttpResponseEventPacketList "Paginated list of event packets"
// @Failure 400 {object} map[string]string "Invalid query parameters"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /event-packets [get]
func (h *GinEventPacketHandler) FilterEventPackets(c *gin.Context) {
	var filter httpdto.HttpFilterEventPacket
	if err := c.ShouldBindQuery(&filter); err != nil {
		handleError(c, &domain.ValidationError{Reason: "invalid query parameters"})
		return
	}

	domainFilter := filter.ToEventPacketFilter()
	domainFilter.Default()
	if err := domainFilter.Validate(); err != nil {
		handleError(c, err)
		return
	}

	token := getTokenFromHeader(c)
	packets, err := h.usecase.FilterEventPackets(c.Request.Context(), token, domainFilter)
	if handleError(c, err) {
		return
	}

	totalCount, err := h.repo.CountEventPackets(c.Request.Context(), domainFilter)
	if handleError(c, err) {
		return
	}

	resp := httpdto.ToHttpResponseEventPacketListWithPagination(packets, domainFilter, totalCount, h.serviceURLs)
	c.JSON(http.StatusOK, resp)
}
