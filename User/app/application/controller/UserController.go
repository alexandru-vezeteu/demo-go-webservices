package controller

import (
	"userService/application/domain"
	"userService/application/service"
)

type userController struct {
	service service.IUserService
}

func NewUserController(service service.IUserService) *userController {
	return &userController{service: service}
}

func (c *userController) CreateUser(user *domain.User) (*domain.User, error) {
	return c.service.CreateUser(user)
}

func (c *userController) GetUserByID(id int) (*domain.User, error) {
	return c.service.GetUserByID(id)
}

func (c *userController) UpdateUser(id int, updates map[string]interface{}) (*domain.User, error) {
	return c.service.UpdateUser(id, updates)
}

func (c *userController) DeleteUser(id int) (*domain.User, error) {
	return c.service.DeleteUser(id)
}
