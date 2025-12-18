package service

import (
	"context"
	"userService/application/domain"
)

type AuthorizationService interface {
	CanUserViewUser(ctx context.Context, actorID uint, targetUser *domain.User) (bool, error)
	CanUserEditUser(ctx context.Context, actorID uint, targetUser *domain.User) (bool, error)
	CanUserDeleteUser(ctx context.Context, actorID uint, targetUser *domain.User) (bool, error)

	CanUserViewTicket(ctx context.Context, userID uint, ticketCode string) (bool, error)
	CanUserEditTicket(ctx context.Context, userID uint, ticketCode string) (bool, error)
	CanUserBuyTicket(ctx context.Context, userID uint, eventID int) (bool, error)
}
