package handler

import (
	"errors"
	"eventManager/application/domain"
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

func (h *GinTicketHandler) CreateTicket(c *gin.Context) {
	var req httpdto.HttpCreateTicket
	if err := middleware.StrictBindJSON(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ticket := req.ToTicket()
	token := getTokenFromHeader(c)

	ret, err := h.usecase.CreateTicket(c.Request.Context(), token, ticket)

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

	resp := httpdto.ToHttpResponseTicket(ret, h.serviceURLs)
	c.JSON(http.StatusCreated, resp)
}

func (h *GinTicketHandler) PutTicket(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ticket code is required"})
		return
	}

	var req httpdto.HttpCreateTicket
	if err := middleware.StrictBindJSON(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ticket := req.ToTicket()
	token := getTokenFromHeader(c)

	ret, err := h.usecase.PutTicket(c.Request.Context(), token, code, ticket)

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

	resp := httpdto.ToHttpResponseTicket(ret, h.serviceURLs)
	c.JSON(http.StatusOK, resp)
}

func (h *GinTicketHandler) GetTicketByCode(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ticket code is required"})
		return
	}

	token := getTokenFromHeader(c)
	ret, err := h.usecase.GetTicketByCode(c.Request.Context(), token, code)

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

	resp := httpdto.ToHttpResponseTicket(ret, h.serviceURLs)
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
	token := getTokenFromHeader(c)

	ticket, err := h.usecase.UpdateTicket(c.Request.Context(), token, code, updates)

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

	resp := httpdto.ToHttpResponseTicket(ticket, h.serviceURLs)
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
	token := getTokenFromHeader(c)

	ticket, err := h.usecase.UpdateTicket(c.Request.Context(), token, code, updates)

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

	resp := httpdto.ToHttpResponseTicket(ticket, h.serviceURLs)
	c.JSON(http.StatusOK, resp)
}

func (h *GinTicketHandler) DeleteTicket(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ticket code is required"})
		return
	}

	token := getTokenFromHeader(c)
	ret, err := h.usecase.DeleteTicket(c.Request.Context(), token, code)

	var notFoundErr *domain.NotFoundError
	if errors.As(err, &notFoundErr) {
		c.JSON(http.StatusNotFound, gin.H{"error": notFoundErr.Error()})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An unexpected error occurred"})
		return
	}

	resp := httpdto.ToHttpResponseTicket(ret, h.serviceURLs)
	c.JSON(http.StatusOK, resp)
}
