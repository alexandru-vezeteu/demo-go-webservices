package service

import "context"

// AuthorizationService defines methods for authorization checks
type AuthorizationService interface {
	// CanUserSeeUser checks if a user can view another user's profile
	CanUserSeeUser(ctx context.Context, requestingUserID string, targetUserID string) (bool, error)

	// CanUserUpdateUser checks if a user can update another user's profile
	CanUserUpdateUser(ctx context.Context, requestingUserID string, targetUserID string) (bool, error)

	// CanUserDeleteUser checks if a user can delete another user
	CanUserDeleteUser(ctx context.Context, requestingUserID string, targetUserID string) (bool, error)

	// CanUserCreateTicket checks if a user can create a ticket for a specific user
	CanUserCreateTicket(ctx context.Context, requestingUserID string, targetUserID string) (bool, error)

	// CanUserSeeTicket checks if a user can view a specific ticket
	CanUserSeeTicket(ctx context.Context, userID string, ticketID string) (bool, error)
}
