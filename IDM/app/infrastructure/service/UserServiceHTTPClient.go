package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"idmService/application/domain"
	"idmService/application/service"
	"net/http"
	"os"
	"time"
)

type HttpCreateUserRequest struct {
	ID               int     `json:"id"`
	Email            string  `json:"email"`
	FirstName        string  `json:"first_name"`
	LastName         string  `json:"last_name"`
	SocialMediaLinks *string `json:"social_media_links,omitempty"`
}

type HttpResponseUser struct {
	User struct {
		ID        int    `json:"id"`
		Email     string `json:"email"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	} `json:"user"`
}

type UserServiceHTTPClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewUserServiceHTTPClient() *UserServiceHTTPClient {
	userServiceHost := os.Getenv("USER_SERVICE_HOST")
	if userServiceHost == "" {
		userServiceHost = "localhost"
	}
	userServicePort := os.Getenv("USER_SERVICE_PORT")
	if userServicePort == "" {
		userServicePort = "12346"
	}

	baseURL := fmt.Sprintf("http://%s:%s/api/user-manager", userServiceHost, userServicePort)

	return &UserServiceHTTPClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *UserServiceHTTPClient) CreateUser(ctx context.Context, req *service.CreateUserRequest) (*service.CreateUserResponse, error) {
	httpReq := HttpCreateUserRequest{
		ID:        req.ID,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	jsonData, err := json.Marshal(httpReq)
	if err != nil {
		return nil, &domain.InternalError{Operation: "create user profile", Err: fmt.Errorf("failed to marshal request: %w", err)}
	}

	url := fmt.Sprintf("%s/users", c.baseURL)
	httpRequest, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, &domain.InternalError{Operation: "create user profile", Err: fmt.Errorf("failed to create HTTP request: %w", err)}
	}

	httpRequest.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpRequest)
	if err != nil {
		return nil, &domain.InternalError{Operation: "create user profile", Err: fmt.Errorf("failed to call User service: %w", err)}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		var errorResp map[string]string
		json.NewDecoder(resp.Body).Decode(&errorResp)
		return nil, &domain.InternalError{
			Operation: "create user profile",
			Err:       fmt.Errorf("User service returned status %d: %v", resp.StatusCode, errorResp),
		}
	}

	var userResp HttpResponseUser
	if err := json.NewDecoder(resp.Body).Decode(&userResp); err != nil {
		return nil, &domain.InternalError{Operation: "create user profile", Err: fmt.Errorf("failed to decode response: %w", err)}
	}

	return &service.CreateUserResponse{
		ID:        userResp.User.ID,
		Email:     userResp.User.Email,
		FirstName: userResp.User.FirstName,
		LastName:  userResp.User.LastName,
	}, nil
}
