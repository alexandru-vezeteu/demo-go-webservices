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

func resolveIDMServiceURL() (string, error) {
	host := os.Getenv("IDM_HOST")
	port := os.Getenv("IDM_PORT")
	if host == "" || port == "" {
		return "", fmt.Errorf("missing IDM_HOST/IDM_PORT configuration")
	}

	return fmt.Sprintf("%s:%s", host, port), nil
}

func NewIDMClient() (*IDMClient, error) {
	idmURL, err := resolveIDMServiceURL()
	if err != nil {
		return nil, err
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
