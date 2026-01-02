package usecase

import (
	"context"
	"fmt"
	"time"

	"idmService/application/domain"
	"idmService/application/service"
)

type RevokeTokenResult struct {
	Success bool
	Message string
}

type RevokeTokenUseCase interface {
	Execute(ctx context.Context, token string) (*RevokeTokenResult, error)
}

type revokeTokenUseCase struct {
	blacklist    domain.TokenBlacklist
	tokenService service.TokenService
}

func NewRevokeTokenUseCase(blacklist domain.TokenBlacklist, tokenService service.TokenService) RevokeTokenUseCase {
	return &revokeTokenUseCase{
		blacklist:    blacklist,
		tokenService: tokenService,
	}
}

func (uc *revokeTokenUseCase) Execute(ctx context.Context, token string) (*RevokeTokenResult, error) {

	if isBlacklisted, _ := uc.blacklist.IsBlacklisted(token); isBlacklisted {
		return &RevokeTokenResult{
			Success: true,
			Message: "Token was already revoked",
		}, nil
	}

	claims, parseErr := uc.tokenService.ParseToken(token)
	expiresAt := time.Now().Add(24 * time.Hour)
	if parseErr == nil && claims != nil {
		expiresAt = time.Unix(claims.ExpiresAt, 0)
	}

	err := uc.blacklist.Add(token, "manually revoked", expiresAt)
	if err != nil {
		return &RevokeTokenResult{
			Success: true,
			Message: fmt.Sprintf("Token revocation completed with warning: %v", err),
		}, nil
	}

	return &RevokeTokenResult{
		Success: true,
		Message: "Token has been successfully revoked and added to blacklist",
	}, nil
}
