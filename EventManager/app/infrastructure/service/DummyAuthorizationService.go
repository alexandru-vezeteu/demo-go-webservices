package service

import (
	"context"
	"eventManager/application/domain"
	"eventManager/application/service"
)

type DummyAuthorizationService struct{}

func NewDummyAuthorizationService() service.AuthorizationService {
	return &DummyAuthorizationService{}
}

func (s *DummyAuthorizationService) CanUserCreateEvent(ctx context.Context, userID uint) (bool, error) {
	return true, nil
}

func (s *DummyAuthorizationService) CanUserViewEvent(ctx context.Context, userID uint, event *domain.Event) (bool, error) {
	return true, nil
}

func (s *DummyAuthorizationService) CanUserViewEvents(ctx context.Context, userID uint, events []*domain.Event) ([]bool, error) {
	results := make([]bool, len(events))
	for i := range events {
		results[i] = true
	}
	return results, nil
}

func (s *DummyAuthorizationService) CanUserEditEvent(ctx context.Context, userID uint, event *domain.Event) (bool, error) {
	return true, nil
}

func (s *DummyAuthorizationService) CanUserDeleteEvent(ctx context.Context, userID uint, event *domain.Event) (bool, error) {
	return true, nil
}

func (s *DummyAuthorizationService) CanUserViewTicket(ctx context.Context, userID uint, ticket *domain.Ticket) (bool, error) {
	return true, nil
}

func (s *DummyAuthorizationService) CanUserViewTickets(ctx context.Context, userID uint, tickets []*domain.Ticket) ([]bool, error) {
	results := make([]bool, len(tickets))
	for i := range tickets {
		results[i] = true
	}
	return results, nil
}

func (s *DummyAuthorizationService) CanUserEditTicket(ctx context.Context, userID uint, ticket *domain.Ticket) (bool, error) {
	return true, nil
}

func (s *DummyAuthorizationService) CanUserDeleteTicket(ctx context.Context, userID uint, ticket *domain.Ticket) (bool, error) {
	return true, nil
}

func (s *DummyAuthorizationService) CanUserBuyTicket(ctx context.Context, userID uint, event *domain.Event) (bool, error) {
	return true, nil
}

func (s *DummyAuthorizationService) CanUserCreateEventPacket(ctx context.Context, userID uint) (bool, error) {
	return true, nil
}

func (s *DummyAuthorizationService) CanUserViewEventPacket(ctx context.Context, userID uint, packet *domain.EventPacket) (bool, error) {
	return true, nil
}

func (s *DummyAuthorizationService) CanUserViewEventPackets(ctx context.Context, userID uint, packets []*domain.EventPacket) ([]bool, error) {
	results := make([]bool, len(packets))
	for i := range packets {
		results[i] = true
	}
	return results, nil
}

func (s *DummyAuthorizationService) CanUserEditEventPacket(ctx context.Context, userID uint, packet *domain.EventPacket) (bool, error) {
	return true, nil
}

func (s *DummyAuthorizationService) CanUserDeleteEventPacket(ctx context.Context, userID uint, packet *domain.EventPacket) (bool, error) {
	return true, nil
}

func (s *DummyAuthorizationService) CanUserCreateEventPacketInclusion(ctx context.Context, userID uint, eventID int, packetID int) (bool, error) {
	return true, nil
}

func (s *DummyAuthorizationService) CanUserViewEventPacketInclusion(ctx context.Context, userID uint, eventID int, packetID int) (bool, error) {
	return true, nil
}

func (s *DummyAuthorizationService) CanUserUpdateEventPacketInclusion(ctx context.Context, userID uint, eventID int, packetID int) (bool, error) {
	return true, nil
}

func (s *DummyAuthorizationService) CanUserDeleteEventPacketInclusion(ctx context.Context, userID uint, eventID int, packetID int) (bool, error) {
	return true, nil
}
