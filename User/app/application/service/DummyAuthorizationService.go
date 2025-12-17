package service

import "context"

// DummyAuthorizationService is a dummy implementation that always returns true (allowed)
type DummyAuthorizationService struct{}

// NewDummyAuthorizationService creates a new dummy authorization service
func NewDummyAuthorizationService() *DummyAuthorizationService {
	return &DummyAuthorizationService{}
}

// CanUserSeeUser always returns true
func (s *DummyAuthorizationService) CanUserSeeUser(ctx context.Context, requestingUserID string, targetUserID string) (bool, error) {
	return true, nil
}

// CanUserUpdateUser always returns true
func (s *DummyAuthorizationService) CanUserUpdateUser(ctx context.Context, requestingUserID string, targetUserID string) (bool, error) {
	return true, nil
}

// CanUserDeleteUser always returns true
func (s *DummyAuthorizationService) CanUserDeleteUser(ctx context.Context, requestingUserID string, targetUserID string) (bool, error) {
	return true, nil
}

// CanUserCreateTicket always returns true
func (s *DummyAuthorizationService) CanUserCreateTicket(ctx context.Context, requestingUserID string, targetUserID string) (bool, error) {
	return true, nil
}

// CanUserSeeTicket always returns true
func (s *DummyAuthorizationService) CanUserSeeTicket(ctx context.Context, userID string, ticketID string) (bool, error) {
	return true, nil
}
