package handler

import (
	"errors"
	"eventManager/application/domain"
	"eventManager/application/usecase"
	"eventManager/infrastructure/http/gin/middleware"
	"eventManager/infrastructure/http/httpdto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GinEventPacketInclusionHandler struct {
	usecase usecase.EventPacketInclusionUseCase
}

func NewGinEventPacketInclusionHandler(usecase usecase.EventPacketInclusionUseCase) *GinEventPacketInclusionHandler {
	return &GinEventPacketInclusionHandler{usecase: usecase}
}

func (h *GinEventPacketInclusionHandler) CreateEventPacketInclusion(c *gin.Context) {
	eventID, err := middleware.ParseIDParam(c, "event_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	packetID, err := middleware.ParseIDParam(c, "packet_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var dto httpdto.HttpCreateEventPacketInclusion
	if err := middleware.StrictBindJSON(c, &dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	inclusion := dto.ToEventPacketInclusion()
	inclusion.EventID = eventID
	inclusion.PacketID = packetID
	token := getTokenFromHeader(c)

	created, err := h.usecase.CreateEventPacketInclusion(c.Request.Context(), token, inclusion)
	if err != nil {
		var validationErr *domain.ValidationError
		var alreadyExistsErr *domain.AlreadyExistsError
		var foreignKeyErr *domain.ForeignKeyError
		var internalErr *domain.InternalError

		if errors.As(err, &validationErr) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else if errors.As(err, &alreadyExistsErr) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else if errors.As(err, &foreignKeyErr) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if errors.As(err, &internalErr) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, httpdto.ToHttpResponseEventPacketInclusion(created))
}

func (h *GinEventPacketInclusionHandler) GetEventPacketsByEventID(c *gin.Context) {
	eventID, err := middleware.ParseIDParam(c, "event_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token := getTokenFromHeader(c)
	packets, err := h.usecase.GetEventPacketsByEventID(c.Request.Context(), token, eventID)
	if err != nil {
		var notFoundErr *domain.NotFoundError
		var internalErr *domain.InternalError

		if errors.As(err, &notFoundErr) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if errors.As(err, &internalErr) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	response := make([]*httpdto.HttpResponseEventPacket, len(packets))
	for i, packet := range packets {
		response[i] = httpdto.ToHttpResponseEventPacket(packet)
	}

	c.JSON(http.StatusOK, response)
}

func (h *GinEventPacketInclusionHandler) GetEventsByPacketID(c *gin.Context) {
	packetID, err := middleware.ParseIDParam(c, "packet_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token := getTokenFromHeader(c)
	events, err := h.usecase.GetEventsByPacketID(c.Request.Context(), token, packetID)
	if err != nil {
		var notFoundErr *domain.NotFoundError
		var internalErr *domain.InternalError

		if errors.As(err, &notFoundErr) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if errors.As(err, &internalErr) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	response := make([]*httpdto.HttpResponseEvent, len(events))
	for i, event := range events {
		response[i] = httpdto.ToHttpResponseEvent(event)
	}

	c.JSON(http.StatusOK, response)
}

func (h *GinEventPacketInclusionHandler) UpdateEventPacketInclusion(c *gin.Context) {
	eventID, err := middleware.ParseIDParam(c, "event_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	packetID, err := middleware.ParseIDParam(c, "packet_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var dto httpdto.HttpUpdateEventPacketInclusion
	if err := middleware.StrictBindJSON(c, &dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := dto.ToUpdateMap()
	token := getTokenFromHeader(c)
	updated, err := h.usecase.Update(c.Request.Context(), token, eventID, packetID, updates)
	if err != nil {
		var notFoundErr *domain.NotFoundError
		var validationErr *domain.ValidationError
		var foreignKeyErr *domain.ForeignKeyError
		var internalErr *domain.InternalError

		if errors.As(err, &notFoundErr) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if errors.As(err, &validationErr) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else if errors.As(err, &foreignKeyErr) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if errors.As(err, &internalErr) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, httpdto.ToHttpResponseEventPacketInclusion(updated))
}

func (h *GinEventPacketInclusionHandler) DeleteEventPacketInclusion(c *gin.Context) {
	eventID, err := middleware.ParseIDParam(c, "event_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	packetID, err := middleware.ParseIDParam(c, "packet_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token := getTokenFromHeader(c)
	deleted, err := h.usecase.DeleteEventPacketInclusion(c.Request.Context(), token, eventID, packetID)
	if err != nil {
		var notFoundErr *domain.NotFoundError
		var internalErr *domain.InternalError

		if errors.As(err, &notFoundErr) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if errors.As(err, &internalErr) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, httpdto.ToHttpResponseEventPacketInclusion(deleted))
}
