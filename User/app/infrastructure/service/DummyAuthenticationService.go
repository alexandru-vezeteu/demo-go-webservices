package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"userService/application/service"
)

type DummyAuthenticationService struct{}

func NewDummyAuthenticationService() service.AuthenticationService {
	return &DummyAuthenticationService{}
}

func (s *DummyAuthenticationService) WhoIsUser(ctx context.Context, token string) (*service.UserIdentity, error) {
	token = strings.TrimPrefix(token, "Bearer ")
	token = strings.TrimSpace(token)

	if strings.HasPrefix(token, "user-") {
		userIDStr := strings.TrimPrefix(token, "user-")
		userID, err := strconv.ParseUint(userIDStr, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid token format: expected 'user-{id}'")
		}
		return &service.UserIdentity{
			UserID: uint(userID),
			Role:   "owner-event",
		}, nil
	}

	return &service.UserIdentity{
		UserID: 1,
		Role:   "owner-event",
	}, nil
}
