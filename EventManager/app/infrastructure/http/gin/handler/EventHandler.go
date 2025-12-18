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

func getTokenFromHeader(c *gin.Context) string {
	token := c.GetHeader("Authorization")
	if token == "" {
		token = "dummy-token"
	}
	return token
}

func (h *GinEventHandler) CreateEvent(c *gin.Context) {
	var req httpdto.HttpCreateEvent
	if err := middleware.StrictBindJSON(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event := req.ToEvent()
	token := getTokenFromHeader(c)

	ret, err := h.usecase.CreateEvent(c.Request.Context(), token, event)
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

func (h *GinEventHandler) GetEventByID(c *gin.Context) {
	id, err := middleware.ParseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token := getTokenFromHeader(c)
	ret, err := h.usecase.GetEventByID(c.Request.Context(), token, id)

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
	token := getTokenFromHeader(c)

	event, err := h.usecase.UpdateEvent(c.Request.Context(), token, id, updates)

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

func (h *GinEventHandler) DeleteEvent(c *gin.Context) {
	id, err := middleware.ParseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token := getTokenFromHeader(c)
	ret, err := h.usecase.DeleteEvent(c.Request.Context(), token, id)

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

	token := getTokenFromHeader(c)
	events, err := h.usecase.FilterEvents(c.Request.Context(), token, domainFilter)

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
