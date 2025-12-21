package handler

import (
	"errors"
	"eventManager/application/usecase"
	"eventManager/application/domain"
	"eventManager/infrastructure/http/config"
	"eventManager/infrastructure/http/gin/middleware"
	"eventManager/infrastructure/http/httpdto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GinEventPacketHandler struct {
	usecase     usecase.EventPacketUseCase
	serviceURLs *config.ServiceURLs
}

func NewGinEventPacketHandler(usecase usecase.EventPacketUseCase, serviceURLs *config.ServiceURLs) *GinEventPacketHandler {
	return &GinEventPacketHandler{
		usecase:     usecase,
		serviceURLs: serviceURLs,
	}
}












func (h *GinEventPacketHandler) CreateEventPacket(c *gin.Context) {
	var req httpdto.HttpCreateEventPacket
	if err := middleware.StrictBindJSON(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	eventPacket := req.ToEventPacket()
	token := getTokenFromHeader(c)

	ret, err := h.usecase.CreateEventPacket(c.Request.Context(), token, eventPacket)

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

	resp := httpdto.ToHttpResponseEventPacket(ret, h.serviceURLs)
	c.JSON(http.StatusCreated, resp)
}












func (h *GinEventPacketHandler) GetEventPacketByID(c *gin.Context) {
	id, err := middleware.ParseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token := getTokenFromHeader(c)
	ret, err := h.usecase.GetEventPacketByID(c.Request.Context(), token, id)

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

	resp := httpdto.ToHttpResponseEventPacket(ret, h.serviceURLs)

	c.JSON(http.StatusOK, resp)

}














func (h *GinEventPacketHandler) UpdateEventPacket(c *gin.Context) {
	id, err := middleware.ParseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req httpdto.HttpUpdateEventPacket
	if err := middleware.StrictBindJSON(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := req.ToUpdateMap()
	token := getTokenFromHeader(c)

	event, err := h.usecase.UpdateEventPacket(c.Request.Context(), token, id, updates)

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

	resp := httpdto.ToHttpResponseEventPacket(event, h.serviceURLs)
	c.JSON(http.StatusOK, resp)
}












func (h *GinEventPacketHandler) DeleteEventPacket(c *gin.Context) {
	id, err := middleware.ParseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token := getTokenFromHeader(c)
	ret, err := h.usecase.DeleteEventPacket(c.Request.Context(), token, id)

	var notFoundErr *domain.NotFoundError
	if errors.As(err, &notFoundErr) {
		c.JSON(http.StatusNotFound, gin.H{"error": notFoundErr.Error()})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	resp := httpdto.ToHttpResponseEventPacket(ret, h.serviceURLs)

	c.JSON(http.StatusOK, resp)
}
