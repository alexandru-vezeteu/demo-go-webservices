package usecase

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"idmService/application/domain"
	"idmService/application/repository"
	"idmService/application/service"
)

type VerifyTokenResult struct {
	Valid       bool
	Email       string
	Message     string
	UserID      string
	Role        string
	Issuer      string
	ExpiresAt   int64
	Expired     bool
	Blacklisted bool
}

type VerifyTokenUseCase interface {
	Execute(ctx context.Context, token string) (*VerifyTokenResult, error)
}

type verifyTokenUseCase struct {
	userRepo     repository.UserRepository
	tokenService service.TokenService
	blacklist    domain.TokenBlacklist
}

func NewVerifyTokenUseCase(
	userRepo repository.UserRepository,
	tokenService service.TokenService,
	blacklist domain.TokenBlacklist,
) VerifyTokenUseCase {
	return &verifyTokenUseCase{
		userRepo:     userRepo,
		tokenService: tokenService,
		blacklist:    blacklist,
	}
}

func (uc *verifyTokenUseCase) Execute(ctx context.Context, token string) (*VerifyTokenResult, error) {

	if isBlacklisted, reason := uc.blacklist.IsBlacklisted(token); isBlacklisted {
		return nil, &domain.TokenError{Blacklisted: true, Reason: reason}
	}

	claims, err := uc.tokenService.ParseToken(token)
	if err != nil {
		uc.blacklist.Add(token, fmt.Sprintf("corrupted: %v", err), time.Now().Add(24*time.Hour))
		return nil, err
	}

	email := ""
	if claims.UserID != "" {
		userID, parseErr := strconv.ParseUint(claims.UserID, 10, 32)
		if parseErr == nil {
			user, repoErr := uc.userRepo.FindByID(ctx, uint(userID))
			if repoErr == nil && user != nil {
				email = user.Email
			}
		}
	}

	if claims.IsExpired {
		expiresAt := time.Unix(claims.ExpiresAt, 0)
		uc.blacklist.Add(token, "expired", expiresAt)
		return nil, &domain.TokenError{Expired: true, Reason: "token has expired"}
	}

	return &VerifyTokenResult{
		Valid:     true,
		Email:     email,
		Message:   "Token is valid",
		UserID:    claims.UserID,
		Role:      claims.Role,
		Issuer:    claims.Issuer,
		ExpiresAt: claims.ExpiresAt,
	}, nil
}
