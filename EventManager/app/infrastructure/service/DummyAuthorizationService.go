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

// oricine poate vedea/crea un event
func (s *DummyAuthorizationService) CanUserCreateEvent(ctx context.Context, user service.UserIdentity) (bool, error) {
	return user.Role == service.RoleOwnerEvent, nil
}

func (s *DummyAuthorizationService) CanUserViewEvent(ctx context.Context, user service.UserIdentity, event *domain.Event) (bool, error) {
	return true, nil
}

func (s *DummyAuthorizationService) CanUserViewEvents(ctx context.Context, user service.UserIdentity, events []*domain.Event) ([]bool, error) {
	results := make([]bool, len(events))
	for i := range events {
		results[i] = true
	}
	return results, nil
}

// doar ownerul poate sterge/edita
func (s *DummyAuthorizationService) CanUserEditEvent(ctx context.Context, user service.UserIdentity, event *domain.Event) (bool, error) {
	return user.UserID == uint(event.OwnerID), nil
}

func (s *DummyAuthorizationService) CanUserDeleteEvent(ctx context.Context, user service.UserIdentity, event *domain.Event) (bool, error) {
	return user.UserID == uint(event.OwnerID), nil
}

// NIMENI NU POATE VEDEA NICIUN BILET. BILETE SE CONSULTA DIN SERVICIUL CLIENTI
func (s *DummyAuthorizationService) CanUserViewTicket(ctx context.Context, user service.UserIdentity, ticket *domain.Ticket) (bool, error) {
	return false, nil
}

func (s *DummyAuthorizationService) CanUserViewTickets(ctx context.Context, user service.UserIdentity, tickets []*domain.Ticket) ([]bool, error) {
	results := make([]bool, len(tickets))
	for i := range tickets {
		results[i] = false
	}
	return results, nil
}

// ce sens are sa editez un bilet??
func (s *DummyAuthorizationService) CanUserEditTicket(ctx context.Context, user service.UserIdentity, ticket *domain.Ticket) (bool, error) {
	return false, nil
}

// nu se poate sterge
func (s *DummyAuthorizationService) CanUserDeleteTicket(ctx context.Context, user service.UserIdentity, ticket *domain.Ticket) (bool, error) {
	return false, nil
}

func (s *DummyAuthorizationService) CanUserBuyTicket(ctx context.Context, user service.UserIdentity, event *domain.Event) (bool, error) {
	return true, nil
}

func (s *DummyAuthorizationService) CanUserCreateTicket(ctx context.Context, user service.UserIdentity) (bool, error) {
	return user.Role == service.RoleServiceClient, nil
}

func (s *DummyAuthorizationService) CanUserCreateEventPacket(ctx context.Context, user service.UserIdentity) (bool, error) {
	return user.Role == service.RoleOwnerEvent, nil
}

func (s *DummyAuthorizationService) CanUserViewEventPacket(ctx context.Context, user service.UserIdentity, packet *domain.EventPacket) (bool, error) {
	return true, nil
}

func (s *DummyAuthorizationService) CanUserViewEventPackets(ctx context.Context, user service.UserIdentity, packets []*domain.EventPacket) ([]bool, error) {
	results := make([]bool, len(packets))
	for i := range packets {
		results[i] = true
	}
	return results, nil
}

func (s *DummyAuthorizationService) CanUserEditEventPacket(ctx context.Context, user service.UserIdentity, packet *domain.EventPacket) (bool, error) {
	return user.UserID == uint(packet.OwnerID), nil
}

func (s *DummyAuthorizationService) CanUserDeleteEventPacket(ctx context.Context, user service.UserIdentity, packet *domain.EventPacket) (bool, error) {
	return user.UserID == uint(packet.OwnerID), nil
}

func (s *DummyAuthorizationService) CanUserCreateEventPacketInclusion(ctx context.Context, user service.UserIdentity, eventID int, packetID int) (bool, error) {
	return user.Role == service.RoleOwnerEvent, nil
}

func (s *DummyAuthorizationService) CanUserViewEventPacketInclusion(ctx context.Context, user service.UserIdentity, eventID int, packetID int) (bool, error) {
	return true, nil
}

func (s *DummyAuthorizationService) CanUserUpdateEventPacketInclusion(ctx context.Context, user service.UserIdentity, eventID int, packetID int) (bool, error) {
	return user.Role == service.RoleOwnerEvent, nil
}

func (s *DummyAuthorizationService) CanUserDeleteEventPacketInclusion(ctx context.Context, user service.UserIdentity, eventID int, packetID int) (bool, error) {
	return user.Role == service.RoleOwnerEvent, nil
}
