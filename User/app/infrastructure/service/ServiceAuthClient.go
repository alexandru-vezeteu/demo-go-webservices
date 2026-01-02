package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	pb "idmService/proto"
)

type ServiceAuthClient struct {
	idmClient     pb.IdentityServiceClient
	email         string
	password      string
	cachedToken   string
	expiresAt     time.Time
	mutex         sync.RWMutex
	refreshBefore time.Duration
}

func NewServiceAuthClient(idmClient pb.IdentityServiceClient, email, password string) *ServiceAuthClient {
	return &ServiceAuthClient{
		idmClient:     idmClient,
		email:         email,
		password:      password,
		refreshBefore: 5 * time.Minute,
	}
}

func (c *ServiceAuthClient) GetServiceToken(ctx context.Context) (string, error) {
	c.mutex.RLock()
	if c.cachedToken != "" && time.Now().Add(c.refreshBefore).Before(c.expiresAt) {
		token := c.cachedToken
		c.mutex.RUnlock()
		return token, nil
	}
	c.mutex.RUnlock()

	return c.refreshToken(ctx)
}

func (c *ServiceAuthClient) refreshToken(ctx context.Context) (string, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.cachedToken != "" && time.Now().Add(c.refreshBefore).Before(c.expiresAt) {
		return c.cachedToken, nil
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	req := &pb.LoginRequest{
		Email:    c.email,
		Password: c.password,
	}

	resp, err := c.idmClient.Login(ctxWithTimeout, req)
	if err != nil {
		return "", fmt.Errorf("failed to authenticate service account: %w", err)
	}

	if !resp.Success {
		return "", fmt.Errorf("service account authentication failed: %s", resp.Message)
	}

	if resp.Role != "serviciu_clienti" {
		return "", fmt.Errorf("service account has incorrect role: expected 'serviciu_clienti', got '%s'", resp.Role)
	}

	c.cachedToken = resp.Token
	c.expiresAt = time.Now().Add(1 * time.Hour)

	return c.cachedToken, nil
}

func (c *ServiceAuthClient) IsConfigured() bool {
	return c.email != "" && c.password != ""
}
