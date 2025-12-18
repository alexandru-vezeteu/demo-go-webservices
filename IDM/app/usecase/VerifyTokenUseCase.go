package usecase

import (
	"context"
	"fmt"
	"strconv"

	"idmService/domain"
	"idmService/infrastructure/blacklist"
	"idmService/service"
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
	userRepo     domain.UserRepository
	tokenService service.TokenService
	blacklist    *blacklist.InMemoryBlacklist
}

func NewVerifyTokenUseCase(userRepo domain.UserRepository, tokenService service.TokenService, blacklist *blacklist.InMemoryBlacklist) VerifyTokenUseCase {
	return &verifyTokenUseCase{
		userRepo:     userRepo,
		tokenService: tokenService,
		blacklist:    blacklist,
	}
}

func (uc *verifyTokenUseCase) Execute(ctx context.Context, token string) (*VerifyTokenResult, error) {
	
	if isBlacklisted, reason := uc.blacklist.IsBlacklisted(token); isBlacklisted {
		return &VerifyTokenResult{
			Valid:       false,
			Message:     fmt.Sprintf("Token is blacklisted: %s", reason),
			Blacklisted: true,
			Expired:     false,
		}, nil
	}

	
	claims, err := uc.tokenService.ParseToken(token)
	if err != nil {
		
		uc.blacklist.Add(token, fmt.Sprintf("corrupted: %v", err))
		return &VerifyTokenResult{
			Valid:       false,
			Message:     fmt.Sprintf("Token is corrupted and has been blacklisted: %v", err),
			Blacklisted: true,
			Expired:     false,
		}, nil
	}

	
	email := ""
	if claims.UserID != "" {
		userID, err := strconv.ParseUint(claims.UserID, 10, 32)
		if err == nil {
			user, err := uc.userRepo.FindByID(ctx, uint(userID))
			if err == nil && user != nil {
				email = user.Email
			}
		}
	}

	
	if claims.IsExpired {
		uc.blacklist.Add(token, "expired")
		return &VerifyTokenResult{
			Valid:       false,
			Email:       email,
			Message:     "Token has expired and has been blacklisted",
			UserID:      claims.UserID,
			Role:        claims.Role,
			Issuer:      claims.Issuer,
			ExpiresAt:   claims.ExpiresAt,
			Expired:     true,
			Blacklisted: true,
		}, nil
	}

	
	return &VerifyTokenResult{
		Valid:       true,
		Email:       email,
		Message:     "Token is valid",
		UserID:      claims.UserID,
		Role:        claims.Role,
		Issuer:      claims.Issuer,
		ExpiresAt:   claims.ExpiresAt,
		Expired:     false,
		Blacklisted: false,
	}, nil
}
