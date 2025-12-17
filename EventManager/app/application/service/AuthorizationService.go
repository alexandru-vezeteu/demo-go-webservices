package service

import "context"

// AuthorizationService defines methods for authorization checks
type AuthorizationService interface {
	// CanUserSeeEvent checks if a user can view a specific event
	CanUserSeeEvent(ctx context.Context, userID string, eventID string) (bool, error)

	// CanUserSeeTicket checks if a user can view a specific ticket
	CanUserSeeTicket(ctx context.Context, userID string, ticketID string) (bool, error)

	// CanUserCreateEvent checks if a user can create events
	CanUserCreateEvent(ctx context.Context, userID string) (bool, error)

	// CanUserUpdateEvent checks if a user can update a specific event
	CanUserUpdateEvent(ctx context.Context, userID string, eventID string) (bool, error)

	// CanUserDeleteEvent checks if a user can delete a specific event
	CanUserDeleteEvent(ctx context.Context, userID string, eventID string) (bool, error)

	// CanUserSeeEventPacket checks if a user can view a specific event packet
	CanUserSeeEventPacket(ctx context.Context, userID string, packetID string) (bool, error)

	// CanUserCreateEventPacket checks if a user can create event packets
	CanUserCreateEventPacket(ctx context.Context, userID string) (bool, error)

	// CanUserUpdateEventPacket checks if a user can update a specific event packet
	CanUserUpdateEventPacket(ctx context.Context, userID string, packetID string) (bool, error)

	// CanUserDeleteEventPacket checks if a user can delete a specific event packet
	CanUserDeleteEventPacket(ctx context.Context, userID string, packetID string) (bool, error)
}
