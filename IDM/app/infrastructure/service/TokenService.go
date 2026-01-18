package service

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"idmService/application/domain"
	appservice "idmService/application/service"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type tokenService struct{}

func getJWTSecret() ([]byte, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is not set")
	}
	return []byte(secret), nil
}

func NewTokenService() appservice.TokenService {
	return &tokenService{}
}

func resolveIDMServiceURL() (string, error) {
	existing := os.Getenv("IDM_SERVICE_URL")
	if existing != "" {
		return existing, nil
	}

	host := os.Getenv("IDM_HOST")
	port := os.Getenv("IDM_PORT")
	if host == "" || port == "" {
		return "", fmt.Errorf("missing IDM_SERVICE_URL or IDM_HOST/IDM_PORT configuration")
	}

	return fmt.Sprintf("http://%s:%s", host, port), nil
}

func (s *tokenService) GenerateJWT(user *domain.User) (string, error) {
	issuer, err := resolveIDMServiceURL()
	if err != nil {
		return "", &domain.ConfigurationError{Key: "IDM_SERVICE_URL"}
	}

	jwtSecret, err := getJWTSecret()
	if err != nil {
		return "", &domain.ConfigurationError{Key: "JWT_SECRET"}
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
	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", &domain.InternalError{Operation: "token signing", Err: err}
	}
	return signedToken, nil
}

func (s *tokenService) ParseToken(tokenString string) (*appservice.TokenClaims, error) {
	jwtSecret, err := getJWTSecret()
	if err != nil {
		return nil, &domain.ConfigurationError{Key: "JWT_SECRET"}
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil && err != jwt.ErrTokenExpired {
		return nil, &domain.TokenError{Corrupted: true, Reason: err.Error()}
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, &domain.TokenError{Corrupted: true, Reason: "invalid claims format"}
	}

	userIDStr, _ := claims["sub"].(string)
	role, _ := claims["role"].(string)
	issuer, _ := claims["iss"].(string)
	exp, _ := claims["exp"].(float64)

	expirationTime := time.Unix(int64(exp), 0)
	isExpired := time.Now().After(expirationTime)

	if userIDStr == "" || role == "" || issuer == "" {
		return nil, &domain.TokenError{Corrupted: true, Reason: "missing required claims"}
	}

	return &appservice.TokenClaims{
		UserID:    userIDStr,
		Role:      role,
		Issuer:    issuer,
		ExpiresAt: int64(exp),
		IsValid:   !isExpired,
		IsExpired: isExpired,
	}, nil
}
