package handler

import (
	"errors"
	"eventManager/application/usecase"
	"eventManager/application/domain"
	"eventManager/infrastructure/http/gin/middleware"
	"eventManager/infrastructure/http/httpdto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GinEventHandler struct {
	usecase usecase.EventUseCase
}

func NewGinEventHandler(usecase usecase.EventUseCase) *GinEventHandler {
	return &GinEventHandler{usecase: usecase}
}

// @Summary      Create a new event
// @Description  Adds a new event to the system.
// @Tags         events
// @Accept       json
// @Produce      json
// @Param        event body httpdto.HttpCreateEvent true "Event to create"
// @Success      201  {object}  httpdto.HttpResponseEvent  "Event created successfully"
// @Failure      400  {object}  map[string]interface{} "Invalid request format or validation error"
// @Failure      409  {object}  map[string]interface{} "Event already exists"
// @Failure      500  {object}  map[string]interface{} "Internal server error"
// @Router       /events [post]
func (h *GinEventHandler) CreateEvent(c *gin.Context) {
	var req httpdto.HttpCreateEvent
	if err := middleware.StrictBindJSON(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event := req.ToEvent()

	ret, err := h.usecase.CreateEvent(event)
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	resp := httpdto.ToHttpResponseEvent(ret)
	c.JSON(http.StatusCreated, resp)
}

// @Summary      Get an event by ID
// @Description  Retrieves a single event using its unique integer ID
// @Tags         events
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Event ID"
// @Success      200  {object}  httpdto.HttpResponseEvent  "The requested event"
// @Failure      400  {object}  map[string]interface{} "Invalid event ID format or validation error"
// @Failure      404  {object}  map[string]interface{} "Event not found"
// @Failure      500  {object}  map[string]interface{} "Internal server error"
// @Router       /events/{id} [get]
func (h *GinEventHandler) GetEventByID(c *gin.Context) {
	id, err := middleware.ParseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ret, err := h.usecase.GetEventByID(id)

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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	resp := httpdto.ToHttpResponseEvent(ret)

	c.JSON(http.StatusOK, resp)

}

// @Summary      Update an event
// @Description  Partially updates an existing event by its ID.
// @Tags         events
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Event ID"
// @Param        updates body httpdto.HttpUpdateEvent true "Fields to update"
// @Success      200  {object}  httpdto.HttpResponseEvent "Event updated successfully"
// @Failure      400  {object}  map[string]interface{} "Invalid event ID format or request body"
// @Failure      404  {object}  map[string]interface{} "Event not found"
// @Failure      409  {object}  map[string]interface{} "Name already taken"
// @Failure      500  {object}  map[string]interface{} "Internal server error"
// @Router       /events/{id} [patch]
func (h *GinEventHandler) UpdateEvent(c *gin.Context) {
	id, err := middleware.ParseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req httpdto.HttpUpdateEvent
	if err := middleware.StrictBindJSON(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := req.ToUpdateMap()

	event, err := h.usecase.UpdateEvent(id, updates)

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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	resp := httpdto.ToHttpResponseEvent(event)
	c.JSON(http.StatusOK, resp)
}

// @Summary      Delete an event
// @Description  Deletes an event by its ID and returns the deleted event.
// @Tags         events
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Event ID"
// @Success      200  {object}  httpdto.HttpResponseEvent "Event deleted successfully"
// @Failure      400  {object}  map[string]interface{} "Invalid event ID format"
// @Failure      404  {object}  map[string]interface{} "Event not found"
// @Failure      500  {object}  map[string]interface{} "Internal server error"
// @Router       /events/{id} [delete]
func (h *GinEventHandler) DeleteEvent(c *gin.Context) {
	id, err := middleware.ParseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ret, err := h.usecase.DeleteEvent(id)

	var notFoundErr *domain.NotFoundError
	if errors.As(err, &notFoundErr) {
		c.JSON(http.StatusNotFound, gin.H{"error": notFoundErr.Error()})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	resp := httpdto.ToHttpResponseEvent(ret)

	c.JSON(http.StatusOK, resp)
}

// @Summary      Filter and list events
// @Description  Retrieves a list of events with optional filters, sorting, and pagination.
// @Tags         events
// @Accept       json
// @Produce      json
// @Param        name query string false "Filter by name (partial match)"
// @Param        location query string false "Filter by location (partial match)"
// @Param        description query string false "Filter by description (partial match)"
// @Param        min_seats query int false "Minimum seats"
// @Param        max_seats query int false "Maximum seats"
// @Param        page query int false "Page number" default(1)
// @Param        per_page query int false "Items per page" default(10)
// @Param        order_by query string false "Sort order: name_asc, name_desc, seats_asc, seats_desc"
// @Success      200  {object}  httpdto.HttpResponseEventList "A list of events"
// @Failure      400  {object}  map[string]interface{} "Invalid query parameters or validation error"
// @Failure      500  {object}  map[string]interface{} "Internal server error"
// @Router       /events [get]
func (h *GinEventHandler) FilterEvents(c *gin.Context) {
	var filter httpdto.HttpFilterEvent
	allowedParams := []string{"name", "location", "description", "min_seats", "max_seats", "page", "per_page", "order_by"}
	if err := middleware.StrictBindQuery(c, &filter, allowedParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	domainFilter := filter.ToEventFilter()
	domainFilter.Default()
	if err := domainFilter.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	events, err := h.usecase.FilterEvents(domainFilter)

	var internalErr *domain.InternalError
	if errors.As(err, &internalErr) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	resp := httpdto.ToHttpResponseEventList(events)
	c.JSON(http.StatusOK, resp)

}
