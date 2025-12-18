package service

import (
	"context"
	"eventManager/application/domain"
)

type AuthorizationService interface {
	CanUserCreateEvent(ctx context.Context, userID uint) (bool, error)
	CanUserViewEvent(ctx context.Context, userID uint, event *domain.Event) (bool, error)
	CanUserViewEvents(ctx context.Context, userID uint, events []*domain.Event) ([]bool, error)
	CanUserEditEvent(ctx context.Context, userID uint, event *domain.Event) (bool, error)
	CanUserDeleteEvent(ctx context.Context, userID uint, event *domain.Event) (bool, error)

	CanUserViewTicket(ctx context.Context, userID uint, ticket *domain.Ticket) (bool, error)
	CanUserViewTickets(ctx context.Context, userID uint, tickets []*domain.Ticket) ([]bool, error)
	CanUserEditTicket(ctx context.Context, userID uint, ticket *domain.Ticket) (bool, error)
	CanUserDeleteTicket(ctx context.Context, userID uint, ticket *domain.Ticket) (bool, error)
	CanUserBuyTicket(ctx context.Context, userID uint, event *domain.Event) (bool, error)

	CanUserCreateEventPacket(ctx context.Context, userID uint) (bool, error)
	CanUserViewEventPacket(ctx context.Context, userID uint, packet *domain.EventPacket) (bool, error)
	CanUserViewEventPackets(ctx context.Context, userID uint, packets []*domain.EventPacket) ([]bool, error)
	CanUserEditEventPacket(ctx context.Context, userID uint, packet *domain.EventPacket) (bool, error)
	CanUserDeleteEventPacket(ctx context.Context, userID uint, packet *domain.EventPacket) (bool, error)

	CanUserCreateEventPacketInclusion(ctx context.Context, userID uint, eventID int, packetID int) (bool, error)
	CanUserViewEventPacketInclusion(ctx context.Context, userID uint, eventID int, packetID int) (bool, error)
	CanUserUpdateEventPacketInclusion(ctx context.Context, userID uint, eventID int, packetID int) (bool, error)
	CanUserDeleteEventPacketInclusion(ctx context.Context, userID uint, eventID int, packetID int) (bool, error)
}
