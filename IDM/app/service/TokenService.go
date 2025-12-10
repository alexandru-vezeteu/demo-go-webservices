package service

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"idmService/domain"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var jwtSecret = []byte("your-secret-key-change-this-in-production")

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

type tokenService struct{}

func NewTokenService() TokenService {
	return &tokenService{}
}

func (s *tokenService) GenerateJWT(user *domain.User) (string, error) {
	issuer := os.Getenv("IDM_SERVICE_URL")
	if issuer == "" {
		return "", fmt.Errorf("IDM_SERVICE_URL environment variable not set")
	}

	tokenID := uuid.New().String()

	claims := jwt.MapClaims{
		"iss":  issuer,
		"sub":  strconv.FormatUint(uint64(user.ID), 10),
		"exp":  time.Now().Add(time.Hour * 1).Unix(),
		"jti":  tokenID,
		"role": string(user.Rol),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func (s *tokenService) ParseToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil && err != jwt.ErrTokenExpired {
		return nil, fmt.Errorf("token parsing failed: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims format")
	}

	userIDStr, _ := claims["sub"].(string)
	role, _ := claims["role"].(string)
	issuer, _ := claims["iss"].(string)
	exp, _ := claims["exp"].(float64)

	expirationTime := time.Unix(int64(exp), 0)
	isExpired := time.Now().After(expirationTime)

	// Check for missing required claims
	if userIDStr == "" || role == "" || issuer == "" {
		return nil, fmt.Errorf("missing required claims")
	}

	return &TokenClaims{
		UserID:    userIDStr,
		Role:      role,
		Issuer:    issuer,
		ExpiresAt: int64(exp),
		IsValid:   !isExpired,
		IsExpired: isExpired,
	}, nil
}
