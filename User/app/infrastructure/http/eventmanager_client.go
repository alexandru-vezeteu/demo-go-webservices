package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
	"userService/application/domain"
)

type ServiceTokenProvider interface {
	GetServiceToken(ctx context.Context) (string, error)
	IsConfigured() bool
}

type EventManagerClient struct {
	baseURL       string
	httpClient    *http.Client
	tokenProvider ServiceTokenProvider
}

type CreateTicketRequest struct {
	PacketID *int `json:"packet_id"`
	EventID  *int `json:"event_id"`
}

type TicketResponse struct {
	Code     string `json:"code"`
	PacketID *int   `json:"packet_id"`
	EventID  *int   `json:"event_id"`
}

func resolveEventManagerURL() (string, error) {
	if url := os.Getenv("EVENT_MANAGER_URL"); url != "" {
		return url, nil
	}

	host := os.Getenv("EVENT_MANAGER_HOST")
	port := os.Getenv("EVENT_MANAGER_PORT")
	if host == "" || port == "" {
		return "", fmt.Errorf("missing EVENT_MANAGER host/port")
	}

	return fmt.Sprintf("http://%s:%s", host, port), nil
}

func NewEventManagerClient(tokenProvider ServiceTokenProvider) *EventManagerClient {
	baseURL, err := resolveEventManagerURL()
	if err != nil {
		panic(err)
	}

	return &EventManagerClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		tokenProvider: tokenProvider,
	}
}

func (c *EventManagerClient) CreateTicket(ctx context.Context, code string, packetID *int, eventID *int) (*TicketResponse, error) {
	url := fmt.Sprintf("%s/api/event-manager/tickets/%s", c.baseURL, code)

	reqBody := CreateTicketRequest{
		PacketID: packetID,
		EventID:  eventID,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, &domain.InternalError{Msg: "failed to marshal request", Err: err}
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, &domain.InternalError{Msg: "failed to create request", Err: err}
	}

	req.Header.Set("Content-Type", "application/json")

	if c.tokenProvider != nil && c.tokenProvider.IsConfigured() {
		serviceToken, err := c.tokenProvider.GetServiceToken(ctx)
		if err != nil {
			return nil, &domain.InternalError{Msg: "failed to get service token", Err: err}
		}
		req.Header.Set("Authorization", "Bearer "+serviceToken)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, &domain.InternalError{Msg: "event manager service unavailable", Err: err}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &domain.InternalError{Msg: "failed to read response", Err: err}
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, &domain.UnauthorizedError{Reason: fmt.Sprintf("service authentication failed: %s", string(body))}
	}

	if resp.StatusCode == http.StatusForbidden {
		return nil, &domain.ForbiddenError{Reason: fmt.Sprintf("service not authorized: %s", string(body))}
	}

	if resp.StatusCode >= 500 {
		return nil, &domain.InternalError{Msg: "event manager service error", Err: fmt.Errorf("status %d: %s", resp.StatusCode, string(body))}
	}

	if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		return nil, &domain.ValidationError{Field: "ticket", Reason: fmt.Sprintf("failed to create ticket: %s", string(body))}
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		return nil, &domain.InternalError{Msg: "unexpected response from event manager", Err: fmt.Errorf("status %d", resp.StatusCode)}
	}

	if resp.StatusCode == http.StatusNoContent {
		return &TicketResponse{Code: code}, nil
	}

	var ticketResp TicketResponse
	if err := json.Unmarshal(body, &ticketResp); err != nil {
		return nil, &domain.InternalError{Msg: "failed to parse response", Err: err}
	}

	return &ticketResp, nil
}
