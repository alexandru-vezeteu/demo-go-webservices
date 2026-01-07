package service

import "context"

type CreateUserRequest struct {
	ID        int
	Email     string
	FirstName string
	LastName  string
}

type CreateUserResponse struct {
	ID        int
	Email     string
	FirstName string
	LastName  string
}

type UserServiceClient interface {
	CreateUser(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error)
}
