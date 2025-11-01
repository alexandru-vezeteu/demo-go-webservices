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

type GinEventHandler struct {
	controller controller.IEventController
}

func NewGinEventHandler(controller controller.IEventController) *GinEventHandler {
	return &GinEventHandler{controller: controller}
}

func (h *GinEventHandler) CreateEvent(c *gin.Context) {
	var req httpdto.HttpCreateEvent
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	event := req.ToEvent()

	ret, err := h.controller.CreateEvent(event)
	var validationErr *domain.EventValidationError
	if errors.As(err, &validationErr) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existsErr *domain.EventAlreadyExistsError
	if errors.As(err, &existsErr) {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	var internalErr *domain.InternalError
	if errors.As(err, &internalErr) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return
	}

	resp := httpdto.ToHttpResponseEvent(ret)
	c.JSON(http.StatusCreated, resp)
}

func (h *GinEventHandler) GetEventByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID format"})
		return
	}

	ret, err := h.controller.GetEventByID(int(id))

	var notFoundErr *domain.EventNotFoundError
	if errors.As(err, &notFoundErr) {
		c.JSON(http.StatusNotFound, gin.H{"error": notFoundErr.Error()})
		return
	}

	var validationErr *domain.EventValidationError
	if errors.As(err, &validationErr) {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	resp := httpdto.ToHttpResponseEvent(ret)

	c.JSON(http.StatusOK, resp)

}

func (h *GinEventHandler) UpdateEvent(c *gin.Context) {

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID format"})
		return
	}
	var req httpdto.HttpUpdateEvent
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	updates := req.ToUpdateMap()

	event, err := h.controller.UpdateEvent(int(id), updates)

	var validationErr *domain.EventValidationError
	var notFoundErr *domain.EventNotFoundError
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

	resp := httpdto.ToHttpResponseEvent(event)
	c.JSON(http.StatusOK, resp)
}

func (h *GinEventHandler) DeleteEvent(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID format"})
		return
	}

	ret, err := h.controller.DeleteEvent(int(id))

	var notFoundErr *domain.EventNotFoundError
	if errors.As(err, &notFoundErr) {
		c.JSON(http.StatusNotFound, gin.H{"error": notFoundErr.Error()})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An unexpected error occurred"})
		return
	}

	resp := httpdto.ToHttpResponseEvent(ret)

	c.JSON(http.StatusOK, resp)
}
