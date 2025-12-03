package server

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"idmService/domain"
	"idmService/infrastructure/blacklist"
	pb "idmService/proto"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var jwtSecret = []byte("your-secret-key-change-this-in-production")

type IdentityServer struct {
	pb.UnimplementedIdentityServiceServer
	userRepo  domain.IUserRepository
	blacklist *blacklist.InMemoryBlacklist
}

func NewIdentityServer(userRepo domain.IUserRepository) *IdentityServer {
	return &IdentityServer{
		userRepo:  userRepo,
		blacklist: blacklist.NewInMemoryBlacklist(),
	}
}

func (s *IdentityServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return &pb.LoginResponse{
			Success: false,
			Token:   "",
			Message: "Database error",
		}, nil
	}

	if user == nil || user.Parola != req.Password {
		return &pb.LoginResponse{
			Success: false,
			Token:   "",
			Message: "Invalid email or password",
		}, nil
	}

	token, err := generateJWT(user)
	if err != nil {
		return &pb.LoginResponse{
			Success: false,
			Token:   "",
			Message: fmt.Sprintf("Failed to generate token: %v", err),
		}, nil
	}

	return &pb.LoginResponse{
		Success: true,
		Token:   token,
		Message: "Login successful",
		UserId:  strconv.FormatUint(uint64(user.ID), 10),
		Role:    string(user.Rol),
		Email:   user.Email,
	}, nil
}

// VerifyToken validates a JWT token's signature and validity period
// IMPORTANT: Both expired and corrupted tokens are added to blacklist
// Returns success (validation processing completed) with indication if token is invalid
func (s *IdentityServer) VerifyToken(ctx context.Context, req *pb.VerifyTokenRequest) (*pb.VerifyTokenResponse, error) {
	if isBlacklisted, reason := s.blacklist.IsBlacklisted(req.Token); isBlacklisted {
		return &pb.VerifyTokenResponse{
			Valid:       false,
			Message:     fmt.Sprintf("Token is blacklisted: %s", reason),
			Blacklisted: true,
			Expired:     false,
		}, nil
	}

	token, err := jwt.Parse(req.Token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil && err != jwt.ErrTokenExpired {
		s.blacklist.Add(req.Token, fmt.Sprintf("corrupted: %v", err))

		return &pb.VerifyTokenResponse{
			Valid:       false,
			Message:     fmt.Sprintf("Token is corrupted and has been blacklisted: %v", err),
			Blacklisted: true,
			Expired:     false,
		}, nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		s.blacklist.Add(req.Token, "invalid claims format")
		return &pb.VerifyTokenResponse{
			Valid:       false,
			Message:     "Token has invalid claims format and has been blacklisted",
			Blacklisted: true,
			Expired:     false,
		}, nil
	}

	userIDStr, _ := claims["sub"].(string)
	role, _ := claims["role"].(string)
	issuer, _ := claims["iss"].(string)
	exp, _ := claims["exp"].(float64)

	email := ""
	if userIDStr != "" {
		userID, err := strconv.ParseUint(userIDStr, 10, 32)
		if err == nil {
			user, err := s.userRepo.FindByID(uint(userID))
			if err == nil && user != nil {
				email = user.Email
			}
		}
	}

	expirationTime := time.Unix(int64(exp), 0)
	isExpired := time.Now().After(expirationTime)

	if isExpired {
		s.blacklist.Add(req.Token, "expired")

		return &pb.VerifyTokenResponse{
			Valid:       false,
			Email:       email,
			Message:     "Token has expired and has been blacklisted",
			UserId:      userIDStr,
			Role:        role,
			Issuer:      issuer,
			ExpiresAt:   int64(exp),
			Expired:     true,
			Blacklisted: true,
		}, nil
	}

	if userIDStr == "" || role == "" || issuer == "" {
		s.blacklist.Add(req.Token, "missing required claims")
		return &pb.VerifyTokenResponse{
			Valid:       false,
			Email:       email,
			Message:     "Token is missing required claims and has been blacklisted",
			UserId:      userIDStr,
			Role:        role,
			Issuer:      issuer,
			ExpiresAt:   int64(exp),
			Blacklisted: true,
			Expired:     false,
		}, nil
	}

	return &pb.VerifyTokenResponse{
		Valid:       true,
		Email:       email,
		Message:     "Token is valid",
		UserId:      userIDStr,
		Role:        role,
		Issuer:      issuer,
		ExpiresAt:   int64(exp),
		Expired:     false,
		Blacklisted: false,
	}, nil
}

func (s *IdentityServer) RevokeToken(ctx context.Context, req *pb.RevokeTokenRequest) (*pb.RevokeTokenResponse, error) {
	if isBlacklisted, _ := s.blacklist.IsBlacklisted(req.Token); isBlacklisted {
		return &pb.RevokeTokenResponse{
			Success: true,
			Message: "Token was already revoked",
		}, nil
	}

	err := s.blacklist.Add(req.Token, "manually revoked")
	if err != nil {
		return &pb.RevokeTokenResponse{
			Success: true,
			Message: fmt.Sprintf("Token revocation completed with warning: %v", err),
		}, nil
	}

	return &pb.RevokeTokenResponse{
		Success: true,
		Message: "Token has been successfully revoked and added to blacklist",
	}, nil
}

func generateJWT(user *domain.User) (string, error) {
	issuer := os.Getenv("IDM_SERVICE_URL")
	if issuer == "" {
		println("WTFF IDM_SERVICE_URL NOT SET")
		os.Exit(-1)
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
