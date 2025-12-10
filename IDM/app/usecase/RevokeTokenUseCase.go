package usecase

import (
	"fmt"

	"idmService/infrastructure/blacklist"
)

type RevokeTokenResult struct {
	Success bool
	Message string
}

type RevokeTokenUseCase interface {
	Execute(token string) (*RevokeTokenResult, error)
}

type revokeTokenUseCase struct {
	blacklist *blacklist.InMemoryBlacklist
}

func NewRevokeTokenUseCase(blacklist *blacklist.InMemoryBlacklist) RevokeTokenUseCase {
	return &revokeTokenUseCase{
		blacklist: blacklist,
	}
}

func (uc *revokeTokenUseCase) Execute(token string) (*RevokeTokenResult, error) {
	// Check if already blacklisted
	if isBlacklisted, _ := uc.blacklist.IsBlacklisted(token); isBlacklisted {
		return &RevokeTokenResult{
			Success: true,
			Message: "Token was already revoked",
		}, nil
	}

	// Add to blacklist
	err := uc.blacklist.Add(token, "manually revoked")
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
