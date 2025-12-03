package domain

import (
	"context"
	pb "idmService/proto"
)

type IIdentityService interface {
	Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error)

	VerifyToken(ctx context.Context, req *pb.VerifyTokenRequest) (*pb.VerifyTokenResponse, error)

	RevokeToken(ctx context.Context, req *pb.RevokeTokenRequest) (*pb.RevokeTokenResponse, error)
}

type TokenBlacklist interface {
	Add(token string, reason string) error

	IsBlacklisted(token string) (bool, string)

	Remove(token string) error
}
