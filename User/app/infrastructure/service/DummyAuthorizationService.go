package service

import (
	"context"
	"userService/application/domain"
	"userService/application/repository"
	"userService/application/service"
)

type DummyAuthorizationService struct {
	userRepo repository.UserRepository
}

func NewDummyAuthorizationService(userRepo repository.UserRepository) service.AuthorizationService {
	return &DummyAuthorizationService{userRepo: userRepo}
}

func (s *DummyAuthorizationService) CanUserViewUser(ctx context.Context, actorID uint, targetUser *domain.User) (bool, error) {
	return actorID == uint(targetUser.ID), nil
}

func (s *DummyAuthorizationService) CanUserEditUser(ctx context.Context, actorID uint, targetUser *domain.User) (bool, error) {
	return actorID == uint(targetUser.ID), nil
}

func (s *DummyAuthorizationService) CanUserDeleteUser(ctx context.Context, actorID uint, targetUser *domain.User) (bool, error) {
	return actorID == uint(targetUser.ID), nil
}

func (s *DummyAuthorizationService) CanUserViewTicket(ctx context.Context, userID uint, ticketCode string) (bool, error) {
	return s.userOwnsTicket(ctx, userID, ticketCode)
}

func (s *DummyAuthorizationService) CanUserEditTicket(ctx context.Context, userID uint, ticketCode string) (bool, error) {
	return s.userOwnsTicket(ctx, userID, ticketCode)
}

func (s *DummyAuthorizationService) CanUserBuyTicket(ctx context.Context, userID uint, eventID int) (bool, error) {
	return true, nil
}

func (s *DummyAuthorizationService) CanUserViewEventCustomers(ctx context.Context, identity *service.UserIdentity, eventID int) (bool, error) {
	return identity.Role == "owner-event", nil
}

func (s *DummyAuthorizationService) CanUserViewPacketCustomers(ctx context.Context, identity *service.UserIdentity, packetID int) (bool, error) {
	return identity.Role == "owner-event", nil
}

func (s *DummyAuthorizationService) userOwnsTicket(ctx context.Context, userID uint, ticketCode string) (bool, error) {
	user, err := s.userRepo.GetByID(ctx, int(userID))
	if err != nil {
		return false, nil
	}

	for _, ticket := range user.TicketList {
		if ticket.Code == ticketCode {
			return true, nil
		}
	}
	return false, nil
}
