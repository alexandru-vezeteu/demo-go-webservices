package handler

import (
	"errors"
	"net/http"
	"strings"
	"userService/application/domain"
	"userService/application/usecase"
	"userService/infrastructure/http/config"
	"userService/infrastructure/http/gin/middleware"
	"userService/infrastructure/http/httpdto"

	"github.com/gin-gonic/gin"
)

type GinUserHandler struct {
	usecase     usecase.UserUsecase
	serviceURLs *config.ServiceURLs
}

func NewGinUserHandler(usecase usecase.UserUsecase, serviceURLs *config.ServiceURLs) *GinUserHandler {
	return &GinUserHandler{
		usecase:     usecase,
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

	var internalErr *domain.InternalError
	if errors.As(err, &internalErr) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return true
	}

	c.JSON(http.StatusInternalServerError, gin.H{"error": "An unexpected error occurred"})
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

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user account with the provided details
// @Tags users
// @Accept json
// @Produce json
// @Param Authorization header string false "Bearer token (optional for user creation)"
// @Param user body httpdto.HttpCreateUser true "User details"
// @Success 201 {object} httpdto.HttpResponseUser "User created successfully"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 409 {object} map[string]string "User already exists"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users [post]
func (h *GinUserHandler) CreateUser(c *gin.Context) {
	var req httpdto.HttpCreateUser
	if err := middleware.StrictBindJSON(c, &req); err != nil {
		handleError(c, err)
		return
	}

	user := req.ToUser()
	token := getTokenFromHeader(c)

	ret, err := h.usecase.CreateUser(c.Request.Context(), token, user)
	if handleError(c, err) {
		return
	}

	resp := httpdto.ToHttpResponseUser(ret, h.serviceURLs)
	c.JSON(http.StatusCreated, resp)
}

// GetUserByID godoc
// @Summary Get user by ID
// @Description Retrieve a specific user by their unique identifier
// @Tags users
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "User ID"
// @Success 200 {object} httpdto.HttpResponseUser "User details"
// @Failure 400 {object} map[string]string "Invalid user ID"
// @Failure 401 {object} map[string]string "Unauthorized - missing or invalid token"
// @Failure 403 {object} map[string]string "Forbidden - insufficient permissions"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/{id} [get]
func (h *GinUserHandler) GetUserByID(c *gin.Context) {
	token, ok := requireAuth(c)
	if !ok {
		return
	}

	id, err := middleware.ParseIDParam(c, "id")
	if err != nil {
		handleError(c, err)
		return
	}

	ret, err := h.usecase.GetUserByID(c.Request.Context(), token, id)
	if handleError(c, err) {
		return
	}

	resp := httpdto.ToHttpResponseUser(ret, h.serviceURLs)
	c.JSON(http.StatusOK, resp)
}

// UpdateUser godoc
// @Summary Update an existing user
// @Description Partially update user details (PATCH - only provided fields are updated)
// @Tags users
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "User ID"
// @Param user body httpdto.HttpUpdateUser true "Fields to update"
// @Success 200 {object} httpdto.HttpResponseUser "User updated successfully"
// @Failure 400 {object} map[string]string "Invalid request body or user ID"
// @Failure 401 {object} map[string]string "Unauthorized - missing or invalid token"
// @Failure 403 {object} map[string]string "Forbidden - cannot update other users"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/{id} [patch]
func (h *GinUserHandler) UpdateUser(c *gin.Context) {
	token, ok := requireAuth(c)
	if !ok {
		return
	}

	id, err := middleware.ParseIDParam(c, "id")
	if err != nil {
		handleError(c, err)
		return
	}

	var req httpdto.HttpUpdateUser
	if err := middleware.StrictBindJSON(c, &req); err != nil {
		handleError(c, err)
		return
	}

	updates := req.ToUpdateMap()

	user, err := h.usecase.UpdateUser(c.Request.Context(), token, id, updates)
	if handleError(c, err) {
		return
	}

	resp := httpdto.ToHttpResponseUser(user, h.serviceURLs)
	c.JSON(http.StatusOK, resp)
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Delete a user account by its ID
// @Tags users
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "User ID"
// @Success 200 {object} httpdto.HttpResponseUser "User deleted successfully"
// @Failure 400 {object} map[string]string "Invalid user ID"
// @Failure 401 {object} map[string]string "Unauthorized - missing or invalid token"
// @Failure 403 {object} map[string]string "Forbidden - cannot delete other users"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/{id} [delete]
func (h *GinUserHandler) DeleteUser(c *gin.Context) {
	token, ok := requireAuth(c)
	if !ok {
		return
	}

	id, err := middleware.ParseIDParam(c, "id")
	if err != nil {
		handleError(c, err)
		return
	}

	ret, err := h.usecase.DeleteUser(c.Request.Context(), token, id)
	if handleError(c, err) {
		return
	}

	resp := httpdto.ToHttpResponseUser(ret, h.serviceURLs)
	c.JSON(http.StatusOK, resp)
}

// CreateTicketForUser godoc
// @Summary Create ticket for user
// @Description Purchase a ticket for a user by creating a ticket through the EventManager service
// @Tags clients
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param user_id path int true "User ID"
// @Param ticket body httpdto.HttpCreateTicketForUser true "Ticket purchase details (packet_id and event_id)"
// @Success 201 {object} httpdto.HttpCreateTicketResponse "Ticket created successfully with ticket code"
// @Failure 400 {object} map[string]string "Invalid request body or user ID"
// @Failure 401 {object} map[string]string "Unauthorized - missing or invalid token"
// @Failure 404 {object} map[string]string "User, event, or packet not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /clients/{user_id}/tickets [post]
func (h *GinUserHandler) CreateTicketForUser(c *gin.Context) {
	token, ok := requireAuth(c)
	if !ok {
		return
	}

	userID, err := middleware.ParseIDParam(c, "user_id")
	if err != nil {
		handleError(c, err)
		return
	}

	var req httpdto.HttpCreateTicketForUser
	if err := middleware.StrictBindJSON(c, &req); err != nil {
		handleError(c, err)
		return
	}

	ticketCode, err := h.usecase.CreateTicketForUser(c.Request.Context(), userID, token, req.PacketID, req.EventID)
	if handleError(c, err) {
		return
	}

	resp := &httpdto.HttpCreateTicketResponse{
		TicketCode: ticketCode,
	}
	c.JSON(http.StatusCreated, resp)
}

// GetCustomersByEventID godoc
// @Summary Get customers who purchased tickets for an event
// @Description Retrieve all customers who have purchased tickets for a specific event (owner only)
// @Tags customers
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param event_id path int true "Event ID"
// @Success 200 {object} httpdto.HttpResponseUserList "List of customers"
// @Failure 400 {object} map[string]string "Invalid event ID"
// @Failure 401 {object} map[string]string "Unauthorized - missing or invalid token"
// @Failure 403 {object} map[string]string "Forbidden - only event owners can view customers"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /events/{event_id}/customers [get]
func (h *GinUserHandler) GetCustomersByEventID(c *gin.Context) {
	token, ok := requireAuth(c)
	if !ok {
		return
	}

	eventID, err := middleware.ParseIDParam(c, "event_id")
	if err != nil {
		handleError(c, err)
		return
	}

	customers, err := h.usecase.GetCustomersByEventID(c.Request.Context(), token, eventID)
	if handleError(c, err) {
		return
	}

	resp := httpdto.ToHttpResponseUserList(customers, h.serviceURLs)
	c.JSON(http.StatusOK, resp)
}

// GetCustomersByPacketID godoc
// @Summary Get customers who purchased tickets for a packet
// @Description Retrieve all customers who have purchased tickets for a specific packet (owner only)
// @Tags customers
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param packet_id path int true "Packet ID"
// @Success 200 {object} httpdto.HttpResponseUserList "List of customers"
// @Failure 400 {object} map[string]string "Invalid packet ID"
// @Failure 401 {object} map[string]string "Unauthorized - missing or invalid token"
// @Failure 403 {object} map[string]string "Forbidden - only packet owners can view customers"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /packets/{packet_id}/customers [get]
func (h *GinUserHandler) GetCustomersByPacketID(c *gin.Context) {
	token, ok := requireAuth(c)
	if !ok {
		return
	}

	packetID, err := middleware.ParseIDParam(c, "packet_id")
	if err != nil {
		handleError(c, err)
		return
	}

	customers, err := h.usecase.GetCustomersByPacketID(c.Request.Context(), token, packetID)
	if handleError(c, err) {
		return
	}

	resp := httpdto.ToHttpResponseUserList(customers, h.serviceURLs)
	c.JSON(http.StatusOK, resp)
}
