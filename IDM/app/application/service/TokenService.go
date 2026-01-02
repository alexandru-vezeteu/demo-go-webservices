package service

import "idmService/application/domain"

type TokenClaims struct {
	UserID    string
	Role      string
	Issuer    string
	ExpiresAt int64
	IsValid   bool
	IsExpired bool
}

type TokenService interface {
	GenerateJWT(user *domain.User) (string, error)
	ParseToken(tokenString string) (*TokenClaims, error)
}
