package service

import (
	"context"
	"userService/application/service"
	"userService/infrastructure/http"
)

type EventManagerHTTPAdapter struct {
	client *http.EventManagerClient
}

func NewEventManagerHTTPAdapter(client *http.EventManagerClient) service.EventManagerService {
	return &EventManagerHTTPAdapter{
		client: client,
	}
}

func (a *EventManagerHTTPAdapter) CreateTicket(ctx context.Context, code string, packetID *int, eventID *int) (*service.TicketResponse, error) {
	resp, err := a.client.CreateTicket(ctx, code, packetID, eventID)
	if err != nil {
		return nil, err
	}

	return &service.TicketResponse{
		Code:     resp.Code,
		PacketID: resp.PacketID,
		EventID:  resp.EventID,
	}, nil
}
