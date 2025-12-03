package usecase

import (
	"userService/application/domain"
	"userService/application/service"
)

type userUsecase struct {
	service service.IUserService
}

func NewUserUsecase(service service.IUserService) *userUsecase {
	return &userUsecase{service: service}
}

func (u *userUsecase) CreateUser(user *domain.User) (*domain.User, error) {
	return u.service.CreateUser(user)
}

func (u *userUsecase) GetUserByID(id int) (*domain.User, error) {
	return u.service.GetUserByID(id)
}

func (u *userUsecase) UpdateUser(id int, updates map[string]interface{}) (*domain.User, error) {
	return u.service.UpdateUser(id, updates)
}

func (u *userUsecase) DeleteUser(id int) (*domain.User, error) {
	return u.service.DeleteUser(id)
}
