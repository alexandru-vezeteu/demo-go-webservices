package service

import (
	"context"
	"userService/application/service"
)

type DummyAuthenticationService struct{}

func NewDummyAuthenticationService() service.AuthenticationService {
	return &DummyAuthenticationService{}
}

func (s *DummyAuthenticationService) WhoIsUser(ctx context.Context, token string) (*service.UserIdentity, error) {
	return &service.UserIdentity{
		UserID: 1,
	}, nil
}
