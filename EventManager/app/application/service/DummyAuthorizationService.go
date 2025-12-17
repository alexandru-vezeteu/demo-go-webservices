package service

import "context"

// DummyAuthorizationService is a dummy implementation that always returns true (allowed)
type DummyAuthorizationService struct{}

// NewDummyAuthorizationService creates a new dummy authorization service
func NewDummyAuthorizationService() *DummyAuthorizationService {
	return &DummyAuthorizationService{}
}

// CanUserSeeEvent always returns true
func (s *DummyAuthorizationService) CanUserSeeEvent(ctx context.Context, userID string, eventID string) (bool, error) {
	return true, nil
}

// CanUserSeeTicket always returns true
func (s *DummyAuthorizationService) CanUserSeeTicket(ctx context.Context, userID string, ticketID string) (bool, error) {
	return true, nil
}

// CanUserCreateEvent always returns true
func (s *DummyAuthorizationService) CanUserCreateEvent(ctx context.Context, userID string) (bool, error) {
	return true, nil
}

// CanUserUpdateEvent always returns true
func (s *DummyAuthorizationService) CanUserUpdateEvent(ctx context.Context, userID string, eventID string) (bool, error) {
	return true, nil
}

// CanUserDeleteEvent always returns true
func (s *DummyAuthorizationService) CanUserDeleteEvent(ctx context.Context, userID string, eventID string) (bool, error) {
	return true, nil
}

// CanUserSeeEventPacket always returns true
func (s *DummyAuthorizationService) CanUserSeeEventPacket(ctx context.Context, userID string, packetID string) (bool, error) {
	return true, nil
}

// CanUserCreateEventPacket always returns true
func (s *DummyAuthorizationService) CanUserCreateEventPacket(ctx context.Context, userID string) (bool, error) {
	return true, nil
}

// CanUserUpdateEventPacket always returns true
func (s *DummyAuthorizationService) CanUserUpdateEventPacket(ctx context.Context, userID string, packetID string) (bool, error) {
	return true, nil
}

// CanUserDeleteEventPacket always returns true
func (s *DummyAuthorizationService) CanUserDeleteEventPacket(ctx context.Context, userID string, packetID string) (bool, error) {
	return true, nil
}
