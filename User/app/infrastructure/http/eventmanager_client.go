package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
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
	eventManagerURL := os.Getenv("EVENT_MANAGER_URL")
	if eventManagerURL == "" {
		eventManagerURL = "http://localhost:12345"
	}

	return &EventManagerClient{
		baseURL: eventManagerURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *EventManagerClient) CreateTicket(code string, packetID *int, eventID *int) (*TicketResponse, int, error) {
	url := fmt.Sprintf("%s/api/event-manager/tickets/%s", c.baseURL, code)

	reqBody := CreateTicketRequest{
		PacketID: packetID,
		EventID:  eventID,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, resp.StatusCode, fmt.Errorf("ticket creation failed with status %d: %s", resp.StatusCode, string(body))
	}

	var ticketResp TicketResponse
	if err := json.Unmarshal(body, &ticketResp); err != nil {
		return nil, resp.StatusCode, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &ticketResp, resp.StatusCode, nil
}
