package service

import (
	"context"
	"eventManager/application/domain"
)

type AuthorizationService interface {
	CanUserCreateEvent(ctx context.Context, user UserIdentity) (bool, error)
	CanUserViewEvent(ctx context.Context, user UserIdentity, event *domain.Event) (bool, error)
	CanUserViewEvents(ctx context.Context, user UserIdentity, events []*domain.Event) ([]bool, error)
	CanUserEditEvent(ctx context.Context, user UserIdentity, event *domain.Event) (bool, error)
	CanUserDeleteEvent(ctx context.Context, user UserIdentity, event *domain.Event) (bool, error)

	CanUserViewTicket(ctx context.Context, user UserIdentity, ticket *domain.Ticket) (bool, error)
	CanUserViewTickets(ctx context.Context, user UserIdentity, tickets []*domain.Ticket) ([]bool, error)
	CanUserEditTicket(ctx context.Context, user UserIdentity, ticket *domain.Ticket) (bool, error)
	CanUserDeleteTicket(ctx context.Context, user UserIdentity, ticket *domain.Ticket) (bool, error)
	CanUserBuyTicket(ctx context.Context, user UserIdentity, event *domain.Event) (bool, error)
	CanUserCreateTicket(ctx context.Context, user UserIdentity) (bool, error)

	CanUserCreateEventPacket(ctx context.Context, user UserIdentity) (bool, error)
	CanUserViewEventPacket(ctx context.Context, user UserIdentity, packet *domain.EventPacket) (bool, error)
	CanUserViewEventPackets(ctx context.Context, user UserIdentity, packets []*domain.EventPacket) ([]bool, error)
	CanUserEditEventPacket(ctx context.Context, user UserIdentity, packet *domain.EventPacket) (bool, error)
	CanUserDeleteEventPacket(ctx context.Context, user UserIdentity, packet *domain.EventPacket) (bool, error)

	CanUserCreateEventPacketInclusion(ctx context.Context, user UserIdentity, eventID int, packetID int) (bool, error)
	CanUserViewEventPacketInclusion(ctx context.Context, user UserIdentity, eventID int, packetID int) (bool, error)
	CanUserUpdateEventPacketInclusion(ctx context.Context, user UserIdentity, eventID int, packetID int) (bool, error)
	CanUserDeleteEventPacketInclusion(ctx context.Context, user UserIdentity, eventID int, packetID int) (bool, error)
}
