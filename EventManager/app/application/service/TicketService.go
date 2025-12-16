package service

import (
	"context"
	"eventManager/application/domain"
	"eventManager/application/repository"
	"fmt"

	"github.com/google/uuid"
)

type TicketService interface {
	CreateTicket(ctx context.Context, ticket *domain.Ticket) (*domain.Ticket, error)
	GetTicketByCode(ctx context.Context, code string) (*domain.Ticket, error)
	UpdateTicket(ctx context.Context, code string, updates map[string]interface{}) (*domain.Ticket, error)
	DeleteTicket(ctx context.Context, code string) (*domain.Ticket, error)
}

type ticketService struct {
	repo          repository.TicketRepository
	eventRepo     repository.EventRepository
	packetRepo    repository.EventPacketRepository
	inclusionRepo repository.EventPacketInclusionRepository
}

func NewTicketService(
	repo repository.TicketRepository,
	eventRepo repository.EventRepository,
	packetRepo repository.EventPacketRepository,
	inclusionRepo repository.EventPacketInclusionRepository,
) TicketService {
	return &ticketService{
		repo:          repo,
		eventRepo:     eventRepo,
		packetRepo:    packetRepo,
		inclusionRepo: inclusionRepo,
	}
}


func (service *ticketService) validateTicket(ticket *domain.Ticket) error {
	if ticket == nil {
		return &domain.ValidationError{Reason: "invalid ticket object"}
	}

	
	if ticket.PacketID == nil && ticket.EventID == nil {
		return &domain.ValidationError{Reason: "ticket must be associated with either a packet or an event"}
	}

	if ticket.PacketID != nil && ticket.EventID != nil {
		return &domain.ValidationError{Reason: "ticket cannot be associated with both a packet and an event"}
	}

	return nil
}

func (service *ticketService) CreateTicket(ctx context.Context, ticket *domain.Ticket) (*domain.Ticket, error) {
	if err := service.validateTicket(ticket); err != nil {
		return nil, err
	}

	
	ticket.Code = uuid.New().String()

	
	if err := service.validateSeatAvailability(ctx, ticket); err != nil {
		return nil, err
	}

	return service.repo.CreateTicket(ctx, ticket)
}





func (service *ticketService) validateSeatAvailability(ctx context.Context, ticket *domain.Ticket) error {
	if ticket.EventID != nil {
		
		event, err := service.eventRepo.GetByID(ctx, *ticket.EventID)
		if err != nil {
			return err
		}

		
		if event.Seats == nil {
			return &domain.ValidationError{
				Reason: fmt.Sprintf("event %d does not have seats defined", *ticket.EventID),
			}
		}

		
		

		
		directTicketsSold, err := service.eventRepo.CountSoldTickets(ctx, *ticket.EventID)
		if err != nil {
			return err
		}

		
		packets, err := service.inclusionRepo.GetEventPacketsByEventID(ctx, *ticket.EventID)
		if err != nil {
			return err
		}

		
		totalPacketSeats := 0
		for _, packet := range packets {
			if packet.AllocatedSeats != nil {
				totalPacketSeats += *packet.AllocatedSeats
			}
		}

		
		totalSeats := *event.Seats
		reservedSeats := directTicketsSold + totalPacketSeats
		availableSeats := totalSeats - reservedSeats

		
		if availableSeats <= 0 {
			return &domain.ValidationError{
				Reason: fmt.Sprintf("event '%s' has no available seats (total: %d, direct tickets: %d, packet allocations: %d)",
					event.Name, totalSeats, directTicketsSold, totalPacketSeats),
			}
		}
	}

	if ticket.PacketID != nil {
		
		packet, err := service.packetRepo.GetByID(ctx, *ticket.PacketID)
		if err != nil {
			return err
		}

		
		if packet.AllocatedSeats == nil {
			return &domain.ValidationError{
				Reason: fmt.Sprintf("packet %d does not have allocated seats defined", *ticket.PacketID),
			}
		}

		
		

		
		soldTickets, err := service.packetRepo.CountSoldTickets(ctx, *ticket.PacketID)
		if err != nil {
			return err
		}

		
		availableSeats := *packet.AllocatedSeats - soldTickets

		
		if availableSeats <= 0 {
			return &domain.ValidationError{
				Reason: fmt.Sprintf("packet '%s' is sold out (%d/%d tickets sold)",
					packet.Name, soldTickets, *packet.AllocatedSeats),
			}
		}
	}

	return nil
}

func (service *ticketService) GetTicketByCode(ctx context.Context, code string) (*domain.Ticket, error) {
	if code == "" {
		return nil, &domain.ValidationError{Reason: "ticket code is required"}
	}
	return service.repo.GetTicketByCode(ctx, code)
}

func (service *ticketService) UpdateTicket(ctx context.Context, code string, updates map[string]interface{}) (*domain.Ticket, error) {
	if code == "" {
		return nil, &domain.ValidationError{Reason: "ticket code is required"}
	}

	if len(updates) == 0 {
		return nil, &domain.ValidationError{Reason: "no fields to update"}
	}

	
	packetID, hasPacketID := updates["packet_id"]
	eventID, hasEventID := updates["event_id"]

	if hasPacketID || hasEventID {
		
		current, err := service.repo.GetTicketByCode(ctx, code)
		if err != nil {
			return nil, err
		}

		
		newPacketID := current.PacketID
		newEventID := current.EventID

		if hasPacketID {
			if packetID == nil {
				newPacketID = nil
			} else {
				val := packetID.(int)
				newPacketID = &val
			}
		}

		if hasEventID {
			if eventID == nil {
				newEventID = nil
			} else {
				val := eventID.(int)
				newEventID = &val
			}
		}

		
		if newPacketID == nil && newEventID == nil {
			return nil, &domain.ValidationError{Reason: "ticket must be associated with either a packet or an event"}
		}

		if newPacketID != nil && newEventID != nil {
			return nil, &domain.ValidationError{Reason: "ticket cannot be associated with both a packet and an event"}
		}

		
		isMovingToNewEvent := hasEventID && (current.EventID == nil || (newEventID != nil && *newEventID != *current.EventID))
		isMovingToNewPacket := hasPacketID && (current.PacketID == nil || (newPacketID != nil && *newPacketID != *current.PacketID))

		if isMovingToNewEvent || isMovingToNewPacket {
			
			tempTicket := &domain.Ticket{
				Code:     code,
				PacketID: newPacketID,
				EventID:  newEventID,
			}

			if err := service.validateSeatAvailability(ctx, tempTicket); err != nil {
				return nil, err
			}
		}
	}

	return service.repo.UpdateTicket(ctx, code, updates)
}

func (service *ticketService) DeleteTicket(ctx context.Context, code string) (*domain.Ticket, error) {
	if code == "" {
		return nil, &domain.ValidationError{Reason: "ticket code is required"}
	}
	return service.repo.DeleteEvent(ctx, code)
}
