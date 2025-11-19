package handler

import (
	"errors"
	"net/http"
	"userService/application/controller"
	"userService/application/domain"
	"userService/infrastructure/http/gin/middleware"
	"userService/infrastructure/http/httpdto"

	"github.com/gin-gonic/gin"
)

type GinUserHandler struct {
	controller controller.IUserController
}

func NewGinUserHandler(controller controller.IUserController) *GinUserHandler {
	return &GinUserHandler{controller: controller}
}

// @Summary      Create a new user
// @Description  Adds a new user to the system.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user body httpdto.HttpCreateUser true "User to create"
// @Success      201  {object}  httpdto.HttpResponseUser  "User created successfully"
// @Failure      400  {object}  map[string]interface{} "Invalid request format or validation error"
// @Failure      409  {object}  map[string]interface{} "User already exists"
// @Failure      500  {object}  map[string]interface{} "Internal error"
// @Router       /users [post]
func (h *GinUserHandler) CreateUser(c *gin.Context) {
	var req httpdto.HttpCreateUser
	if err := middleware.StrictBindJSON(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := req.ToUser()

	ret, err := h.controller.CreateUser(user)
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

// @Summary      Get a user by ID
// @Description  Retrieves a single user using its unique integer ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  httpdto.HttpResponseUser  "The requested user"
// @Failure      400  {object}  map[string]interface{} "Invalid user ID format or validation error"
// @Failure      404  {object}  map[string]interface{} "User not found"
// @Failure      500  {object}  map[string]interface{} "Internal error"
// @Router       /users/{id} [get]
func (h *GinUserHandler) GetUserByID(c *gin.Context) {
	id, err := middleware.ParseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ret, err := h.controller.GetUserByID(id)

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

// @Summary      Update a user
// @Description  Partially updates an existing user by its ID.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Param        updates body httpdto.HttpUpdateUser true "Fields to update"
// @Success      200  {object}  httpdto.HttpResponseUser "User updated successfully"
// @Failure      400  {object}  map[string]interface{} "Invalid user ID format or request body"
// @Failure      404  {object}  map[string]interface{} "User not found"
// @Failure      409  {object}  map[string]interface{} "Resource already exists (e.g., email already taken)"
// @Failure      500  {object}  map[string]interface{} "An unexpected error occurred"
// @Router       /users/{id} [patch]
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

	user, err := h.controller.UpdateUser(id, updates)

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

// @Summary      Delete a user
// @Description  Deletes a user by its ID and returns the deleted user.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  httpdto.HttpResponseUser "User deleted successfully"
// @Failure      400  {object}  map[string]interface{} "Invalid user ID format"
// @Failure      404  {object}  map[string]interface{} "User not found"
// @Failure      500  {object}  map[string]interface{} "An unexpected error occurred"
// @Router       /users/{id} [delete]
func (h *GinUserHandler) DeleteUser(c *gin.Context) {
	id, err := middleware.ParseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ret, err := h.controller.DeleteUser(id)

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
