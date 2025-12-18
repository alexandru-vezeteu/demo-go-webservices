package handler

import (
	"errors"
	"net/http"
	"userService/application/usecase"
	"userService/application/domain"
	"userService/infrastructure/http/gin/middleware"
	"userService/infrastructure/http/httpdto"

	"github.com/gin-gonic/gin"
)

type GinUserHandler struct {
	usecase usecase.UserUsecase
}

func NewGinUserHandler(usecase usecase.UserUsecase) *GinUserHandler {
	return &GinUserHandler{usecase: usecase}
}

func getTokenFromHeader(c *gin.Context) string {
	token := c.GetHeader("Authorization")
	if token == "" {
		token = "dummy-token"
	}
	return token
}

func (h *GinUserHandler) CreateUser(c *gin.Context) {
	var req httpdto.HttpCreateUser
	if err := middleware.StrictBindJSON(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := req.ToUser()
	token := getTokenFromHeader(c)

	ret, err := h.usecase.CreateUser(c.Request.Context(), token, user)
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

	resp := httpdto.ToHttpResponseUser(ret)
	c.JSON(http.StatusCreated, resp)
}

func (h *GinUserHandler) GetUserByID(c *gin.Context) {
	id, err := middleware.ParseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token := getTokenFromHeader(c)
	ret, err := h.usecase.GetUserByID(c.Request.Context(), token, id)

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

	resp := httpdto.ToHttpResponseUser(ret)

	c.JSON(http.StatusOK, resp)
}

func (h *GinUserHandler) UpdateUser(c *gin.Context) {
	id, err := middleware.ParseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req httpdto.HttpUpdateUser
	if err := middleware.StrictBindJSON(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := req.ToUpdateMap()
	token := getTokenFromHeader(c)

	user, err := h.usecase.UpdateUser(c.Request.Context(), token, id, updates)

	var validationErr *domain.ValidationError
	var notFoundErr *domain.NotFoundError
	var existsErr *domain.AlreadyExistsError
	if errors.As(err, &validationErr) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if errors.As(err, &notFoundErr) {
		c.JSON(http.StatusNotFound, gin.H{"error": notFoundErr.Error()})
		return
	}
	if errors.As(err, &existsErr) {
		c.JSON(http.StatusConflict, gin.H{"error": existsErr.Error() + " is already taken"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An unexpected error occurred"})
		return
	}

	resp := httpdto.ToHttpResponseUser(user)
	c.JSON(http.StatusOK, resp)
}

func (h *GinUserHandler) DeleteUser(c *gin.Context) {
	id, err := middleware.ParseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token := getTokenFromHeader(c)
	ret, err := h.usecase.DeleteUser(c.Request.Context(), token, id)

	var notFoundErr *domain.NotFoundError
	if errors.As(err, &notFoundErr) {
		c.JSON(http.StatusNotFound, gin.H{"error": notFoundErr.Error()})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An unexpected error occurred"})
		return
	}

	resp := httpdto.ToHttpResponseUser(ret)

	c.JSON(http.StatusOK, resp)
}

func (h *GinUserHandler) CreateTicketForUser(c *gin.Context) {
	userID, err := middleware.ParseIDParam(c, "user_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization header is required"})
		return
	}

	token := authHeader
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
	}

	var req httpdto.HttpCreateTicketForUser
	if err := middleware.StrictBindJSON(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ticketCode, err := h.usecase.CreateTicketForUser(c.Request.Context(), userID, token, req.PacketID, req.EventID)

	var validationErr *domain.ValidationError
	if errors.As(err, &validationErr) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var forbiddenErr *domain.ForbiddenError
	if errors.As(err, &forbiddenErr) {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	var notFoundErr2 *domain.NotFoundError
	if errors.As(err, &notFoundErr2) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
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

	resp := &httpdto.HttpCreateTicketResponse{
		TicketCode: ticketCode,
	}
	c.JSON(http.StatusCreated, resp)
}
