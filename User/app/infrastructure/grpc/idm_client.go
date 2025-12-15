package grpc

import (
	"context"
	"fmt"
	"os"
	"time"

	pb "idmService/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type IDMClient struct {
	client pb.IdentityServiceClient
	conn   *grpc.ClientConn
}

func NewIDMClient() (*IDMClient, error) {
	idmURL := os.Getenv("IDM_SERVICE_URL")
	if idmURL == "" {
		idmURL = "localhost:50051"
	}

	conn, err := grpc.NewClient(idmURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to IDM service: %w", err)
	}

	client := pb.NewIdentityServiceClient(conn)

	return &IDMClient{
		client: client,
		conn:   conn,
	}, nil
}

func (c *IDMClient) VerifyToken(token string) (*pb.VerifyTokenResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.VerifyTokenRequest{Token: token}
	resp, err := c.client.VerifyToken(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to verify token: %w", err)
	}

	return resp, nil
}

func (c *IDMClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
