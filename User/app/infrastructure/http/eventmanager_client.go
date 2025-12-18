package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
	"userService/application/domain"
)

type EventManagerClient struct {
	baseURL    string
	httpClient *http.Client
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

func NewEventManagerClient() *EventManagerClient {
	host := os.Getenv("EVENT_MANAGER_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("EVENT_MANAGER_PORT")
	if port == "" {
		port = "8080"
	}

	baseURL := fmt.Sprintf("http:%s%s:%s", "/", "/", host, port)

	return &EventManagerClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// CreateTicket creates a ticket in the Event Manager service
// Returns domain errors - no HTTP status codes exposed to use case layer
func (c *EventManagerClient) CreateTicket(code string, packetID *int, eventID *int) (*TicketResponse, error) {
	url := fmt.Sprintf("%s/api/event-manager/tickets/%s", c.baseURL, code)

	reqBody := CreateTicketRequest{
		PacketID: packetID,
		EventID:  eventID,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, &domain.InternalError{Msg: "failed to marshal request", Err: err}
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, &domain.InternalError{Msg: "failed to create request", Err: err}
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, &domain.InternalError{Msg: "event manager service unavailable", Err: err}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &domain.InternalError{Msg: "failed to read response", Err: err}
	}

	// Translate HTTP status codes to domain errors (infrastructure layer responsibility)
	if resp.StatusCode >= 500 {
		return nil, &domain.InternalError{Msg: "event manager service error", Err: fmt.Errorf("status %d: %s", resp.StatusCode, string(body))}
	}

	if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		// Client errors - could be validation, not found, forbidden, etc.
		return nil, &domain.ValidationError{Field: "ticket", Reason: fmt.Sprintf("failed to create ticket: %s", string(body))}
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, &domain.InternalError{Msg: "unexpected response from event manager", Err: fmt.Errorf("status %d", resp.StatusCode)}
	}

	var ticketResp TicketResponse
	if err := json.Unmarshal(body, &ticketResp); err != nil {
		return nil, &domain.InternalError{Msg: "failed to parse response", Err: err}
	}

	return &ticketResp, nil
}
