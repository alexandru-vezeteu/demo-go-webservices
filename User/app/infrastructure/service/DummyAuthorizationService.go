package service

import (
	"context"
	"userService/application/domain"
	"userService/application/service"
)

type DummyAuthorizationService struct{}

func NewDummyAuthorizationService() service.AuthorizationService {
	return &DummyAuthorizationService{}
}

func (s *DummyAuthorizationService) CanUserViewUser(ctx context.Context, actorID uint, targetUser *domain.User) (bool, error) {
	return true, nil
}

func (s *DummyAuthorizationService) CanUserEditUser(ctx context.Context, actorID uint, targetUser *domain.User) (bool, error) {
	return true, nil
}

func (s *DummyAuthorizationService) CanUserDeleteUser(ctx context.Context, actorID uint, targetUser *domain.User) (bool, error) {
	return true, nil
}

func (s *DummyAuthorizationService) CanUserViewTicket(ctx context.Context, userID uint, ticketCode string) (bool, error) {
	return true, nil
}

func (s *DummyAuthorizationService) CanUserEditTicket(ctx context.Context, userID uint, ticketCode string) (bool, error) {
	return true, nil
}

func (s *DummyAuthorizationService) CanUserBuyTicket(ctx context.Context, userID uint, eventID int) (bool, error) {
	return true, nil
}
