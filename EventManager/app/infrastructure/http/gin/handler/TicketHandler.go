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

type GinTicketHandler struct {
	usecase usecase.TicketUseCase
}

func NewGinTicketHandler(usecase usecase.TicketUseCase) *GinTicketHandler {
	return &GinTicketHandler{usecase: usecase}
}












func (h *GinTicketHandler) CreateTicket(c *gin.Context) {
	var req httpdto.HttpCreateTicket
	if err := middleware.StrictBindJSON(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ticket := req.ToTicket()

	ret, err := h.usecase.CreateTicket(c.Request.Context(), ticket)

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












func (h *GinTicketHandler) GetTicketByCode(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ticket code is required"})
		return
	}

	ret, err := h.usecase.GetTicketByCode(c.Request.Context(), code)

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

	ticket, err := h.usecase.UpdateTicket(c.Request.Context(), code, updates)

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













func (h *GinTicketHandler) ReplaceTicket(c *gin.Context) {
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

	ticket, err := h.usecase.UpdateTicket(c.Request.Context(), code, updates)

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












func (h *GinTicketHandler) DeleteTicket(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ticket code is required"})
		return
	}

	ret, err := h.usecase.DeleteTicket(c.Request.Context(), code)

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
