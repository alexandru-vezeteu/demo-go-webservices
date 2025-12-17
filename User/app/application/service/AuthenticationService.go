package service

import "context"

// UserInfo represents the authenticated user information
type UserInfo struct {
	UserID string
	Email  string
	Role   string
}

// AuthenticationService defines methods for user authentication
type AuthenticationService interface {
	// WhoIsUser extracts and validates the token, returning user information
	WhoIsUser(ctx context.Context, token string) (*UserInfo, error)
}
