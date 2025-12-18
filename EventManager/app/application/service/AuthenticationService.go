package service

import (
	"context"
)

type UserIdentity struct {
	UserID    uint
	Email     string
	Role      string
	ExpiresAt int64
}

type AuthenticationService interface {
	WhoIsUser(ctx context.Context, token string) (*UserIdentity, error)
}
