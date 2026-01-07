package handler

import (
	"errors"
	"eventManager/application/domain"
	"eventManager/application/repository"
	"eventManager/application/usecase"
	"eventManager/infrastructure/http/config"
	"eventManager/infrastructure/http/gin/middleware"
	"eventManager/infrastructure/http/httpdto"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type GinEventHandler struct {
	usecase     usecase.EventUseCase
	repo        repository.EventRepository
	serviceURLs *config.ServiceURLs
}

func NewGinEventHandler(usecase usecase.EventUseCase, repo repository.EventRepository, serviceURLs *config.ServiceURLs) *GinEventHandler {
	return &GinEventHandler{
		usecase:     usecase,
		repo:        repo,
		serviceURLs: serviceURLs,
	}
}

func handleError(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}

	var validationErr *domain.ValidationError
	if errors.As(err, &validationErr) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return true
	}

	var invalidReqErr *domain.InvalidRequestError
	if errors.As(err, &invalidReqErr) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return true
	}

	var unauthorizedErr *domain.UnauthorizedError
	if errors.As(err, &unauthorizedErr) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": unauthorizedErr.Error()})
		return true
	}

	var forbiddenErr *domain.ForbiddenError
	if errors.As(err, &forbiddenErr) {
		c.JSON(http.StatusForbidden, gin.H{"error": forbiddenErr.Error()})
		return true
	}

	var notFoundErr *domain.NotFoundError
	if errors.As(err, &notFoundErr) {
		c.JSON(http.StatusNotFound, gin.H{"error": notFoundErr.Error()})
		return true
	}

	var existsErr *domain.AlreadyExistsError
	if errors.As(err, &existsErr) {
		c.JSON(http.StatusConflict, gin.H{"error": existsErr.Error()})
		return true
	}

	var uniqueNameErr *domain.UniqueNameError
	if errors.As(err, &uniqueNameErr) {
		c.JSON(http.StatusConflict, gin.H{"error": uniqueNameErr.Error()})
		return true
	}

	var foreignKeyErr *domain.ForeignKeyError
	if errors.As(err, &foreignKeyErr) {
		c.JSON(http.StatusNotFound, gin.H{"error": foreignKeyErr.Error()})
		return true
	}

	var internalErr *domain.InternalError
	if errors.As(err, &internalErr) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return true
	}

	c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
	return true
}

func getTokenFromHeader(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}

	token := authHeader
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
	}

	return strings.TrimSpace(token)
}

func requireAuth(c *gin.Context) (string, bool) {
	token := getTokenFromHeader(c)
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		return "", false
	}
	return token, true
}

// CreateEvent godoc
// @Summary Create a new event
// @Description Create a new event with the provided details
// @Tags events
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param event body httpdto.HttpCreateEvent true "Event details"
// @Success 201 {object} httpdto.HttpResponseEvent "Event created successfully"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 401 {object} map[string]string "Unauthorized - missing or invalid token"
// @Failure 403 {object} map[string]string "Forbidden - insufficient permissions"
// @Failure 409 {object} map[string]string "Conflict - event already exists"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /events [post]
func (h *GinEventHandler) CreateEvent(c *gin.Context) {
	token, ok := requireAuth(c)
	if !ok {
		return
	}

	var req httpdto.HttpCreateEvent
	if err := middleware.StrictBindJSON(c, &req); err != nil {
		handleError(c, err)
		return
	}

	event := req.ToEvent()

	ret, err := h.usecase.CreateEvent(c.Request.Context(), token, event)
	if handleError(c, err) {
		return
	}

	resp := httpdto.ToHttpResponseEvent(ret, h.serviceURLs)
	c.JSON(http.StatusCreated, resp)
}

// GetEventByID godoc
// @Summary Get event by ID
// @Description Retrieve a specific event by its unique identifier
// @Tags events
// @Accept json
// @Produce json
// @Param id path string true "Event ID (UUID)"
// @Param Authorization header string false "Bearer token (optional)"
// @Success 200 {object} httpdto.HttpResponseEvent "Event details"
// @Failure 400 {object} map[string]string "Invalid event ID format"
// @Failure 404 {object} map[string]string "Event not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /events/{id} [get]
func (h *GinEventHandler) GetEventByID(c *gin.Context) {
	id, err := middleware.ParseIDParam(c, "id")
	if err != nil {
		handleError(c, err)
		return
	}

	token := getTokenFromHeader(c)
	ret, err := h.usecase.GetEventByID(c.Request.Context(), token, id)
	if handleError(c, err) {
		return
	}

	resp := httpdto.ToHttpResponseEvent(ret, h.serviceURLs)
	c.JSON(http.StatusOK, resp)
}

// UpdateEvent godoc
// @Summary Update an existing event
// @Description Partially update event details (PATCH - only provided fields are updated)
// @Tags events
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Event ID (UUID)"
// @Param event body httpdto.HttpUpdateEvent true "Fields to update"
// @Success 200 {object} httpdto.HttpResponseEvent "Event updated successfully"
// @Failure 400 {object} map[string]string "Invalid request body or event ID"
// @Failure 401 {object} map[string]string "Unauthorized - missing or invalid token"
// @Failure 403 {object} map[string]string "Forbidden - not the event owner"
// @Failure 404 {object} map[string]string "Event not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /events/{id} [patch]
func (h *GinEventHandler) UpdateEvent(c *gin.Context) {
	token, ok := requireAuth(c)
	if !ok {
		return
	}

	id, err := middleware.ParseIDParam(c, "id")
	if err != nil {
		handleError(c, err)
		return
	}

	var req httpdto.HttpUpdateEvent
	if err := middleware.StrictBindJSON(c, &req); err != nil {
		handleError(c, err)
		return
	}

	updates := req.ToUpdateMap()

	event, err := h.usecase.UpdateEvent(c.Request.Context(), token, id, updates)
	if handleError(c, err) {
		return
	}

	resp := httpdto.ToHttpResponseEvent(event, h.serviceURLs)
	c.JSON(http.StatusOK, resp)
}

// DeleteEvent godoc
// @Summary Delete an event
// @Description Delete an event by its ID (only the owner can delete)
// @Tags events
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Event ID (UUID)"
// @Success 200 {object} httpdto.HttpResponseEvent "Event deleted successfully"
// @Failure 400 {object} map[string]string "Invalid event ID format"
// @Failure 401 {object} map[string]string "Unauthorized - missing or invalid token"
// @Failure 403 {object} map[string]string "Forbidden - not the event owner"
// @Failure 404 {object} map[string]string "Event not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /events/{id} [delete]
func (h *GinEventHandler) DeleteEvent(c *gin.Context) {
	token, ok := requireAuth(c)
	if !ok {
		return
	}

	id, err := middleware.ParseIDParam(c, "id")
	if err != nil {
		handleError(c, err)
		return
	}

	ret, err := h.usecase.DeleteEvent(c.Request.Context(), token, id)
	if handleError(c, err) {
		return
	}

	resp := httpdto.ToHttpResponseEvent(ret, h.serviceURLs)
	c.JSON(http.StatusOK, resp)
}

// FilterEvents godoc
// @Summary List and filter events
// @Description Get a paginated list of events with optional filters
// @Tags events
// @Accept json
// @Produce json
// @Param Authorization header string false "Bearer token (optional)"
// @Param name query string false "Filter by event name (partial match)"
// @Param location query string false "Filter by location (partial match)"
// @Param description query string false "Filter by description (partial match)"
// @Param min_seats query int false "Minimum number of seats"
// @Param max_seats query int false "Maximum number of seats"
// @Param page query int false "Page number (default: 1)"
// @Param per_page query int false "Items per page (default: 10, max: 100)"
// @Param order_by query string false "Sort field (e.g., 'name', 'created_at')"
// @Success 200 {object} map[string]interface{} "Paginated list of events"
// @Failure 400 {object} map[string]string "Invalid query parameters"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /events [get]
func (h *GinEventHandler) FilterEvents(c *gin.Context) {
	var filter httpdto.HttpFilterEvent
	allowedParams := []string{"name", "location", "description", "min_seats", "max_seats", "page", "per_page", "order_by"}
	if err := middleware.StrictBindQuery(c, &filter, allowedParams); err != nil {
		handleError(c, err)
		return
	}

	domainFilter := filter.ToEventFilter()
	domainFilter.Default()
	if err := domainFilter.Validate(); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	token := getTokenFromHeader(c)
	events, err := h.usecase.FilterEvents(c.Request.Context(), token, domainFilter)
	if handleError(c, err) {
		return
	}

	// Get total count for pagination
	totalCount, err := h.repo.CountEvents(c.Request.Context(), domainFilter)
	if handleError(c, err) {
		return
	}

	resp := httpdto.ToHttpResponseEventListWithPagination(events, domainFilter, totalCount, h.serviceURLs)
	c.JSON(http.StatusOK, resp)
}
