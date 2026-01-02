package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
	"userService/application/domain"
	"userService/application/service"

	pb "idmService/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RealAuthenticationService struct {
	idmHost string
	idmPort string
	client  pb.IdentityServiceClient
	conn    *grpc.ClientConn
}

func NewRealAuthenticationService(idmHost, idmPort string) (service.AuthenticationService, error) {
	address := fmt.Sprintf("%s:%s", idmHost, idmPort)

	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to IDM service at %s: %w", address, err)
	}

	client := pb.NewIdentityServiceClient(conn)

	return &RealAuthenticationService{
		idmHost: idmHost,
		idmPort: idmPort,
		client:  client,
		conn:    conn,
	}, nil
}

func (s *RealAuthenticationService) WhoIsUser(ctx context.Context, token string) (*service.UserIdentity, error) {
	token = strings.TrimPrefix(token, "Bearer ")
	token = strings.TrimSpace(token)

	if token == "" {
		return nil, &domain.UnauthorizedError{
			Reason: "missing authentication token",
		}
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req := &pb.VerifyTokenRequest{Token: token}
	resp, err := s.client.VerifyToken(ctxWithTimeout, req)
	if err != nil {
		return nil, &domain.UnauthorizedError{
			Reason: "failed to verify token with IDM service",
		}
	}

	if !resp.Valid {
		if resp.Expired {
			return nil, &domain.UnauthorizedError{
				Reason: "token has expired",
			}
		}
		if resp.Blacklisted {
			return nil, &domain.UnauthorizedError{
				Reason: "token has been revoked",
			}
		}
		return nil, &domain.UnauthorizedError{
			Reason: fmt.Sprintf("token is invalid: %s", resp.Message),
		}
	}

	userID, err := strconv.ParseUint(resp.UserId, 10, 32)
	if err != nil {
		return nil, &domain.UnauthorizedError{
			Reason: "invalid user ID in token",
		}
	}

	return &service.UserIdentity{
		UserID:    uint(userID),
		Email:     resp.Email,
		Role:      resp.Role,
		ExpiresAt: resp.ExpiresAt,
	}, nil
}

func (s *RealAuthenticationService) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}
