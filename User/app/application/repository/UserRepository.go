package repository

import "userService/application/domain"

type UserRepository interface {
	Create(user *domain.User) (*domain.User, error)
	GetByID(id int) (*domain.User, error)
	Update(id int, updates map[string]interface{}) (*domain.User, error)
	Delete(id int) (*domain.User, error)
}
