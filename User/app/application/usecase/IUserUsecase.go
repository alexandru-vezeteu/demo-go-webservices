package usecase

import "userService/application/domain"

type IUserUsecase interface {
	CreateUser(user *domain.User) (*domain.User, error)
	GetUserByID(id int) (*domain.User, error)
	UpdateUser(id int, updates map[string]interface{}) (*domain.User, error)
	DeleteUser(id int) (*domain.User, error)
}
