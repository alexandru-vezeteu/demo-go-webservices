package service

import "context"

// DummyAuthenticationService is a dummy implementation that always returns success
type DummyAuthenticationService struct{}

// NewDummyAuthenticationService creates a new dummy authentication service
func NewDummyAuthenticationService() *DummyAuthenticationService {
	return &DummyAuthenticationService{}
}

// WhoIsUser returns a dummy user with allowed access
func (s *DummyAuthenticationService) WhoIsUser(ctx context.Context, token string) (*UserInfo, error) {
	return &UserInfo{
		UserID: "dummy-user-id",
		Email:  "dummy@example.com",
		Role:   "admin",
	}, nil
}
