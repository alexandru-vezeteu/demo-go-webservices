package service

import (
	"context"
)

type UserRole string

type UserIdentity struct {
	UserID    uint
	Email     string
	Role      UserRole
	ExpiresAt int64
}

const (
	RoleAdmin         UserRole = "admin"
	RoleOwnerEvent    UserRole = "owner-event"
	RoleClient        UserRole = "client"
	RoleServiceClient UserRole = "serviciu_clienti"
)

type AuthenticationService interface {
	WhoIsUser(ctx context.Context, token string) (*UserIdentity, error)
}
