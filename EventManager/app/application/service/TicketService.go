package service

import (
	"eventManager/application/domain"
	"eventManager/application/repository"
	"fmt"

	"github.com/google/uuid"
)

type TicketService interface {
	CreateTicket(ticket *domain.Ticket) (*domain.Ticket, error)
	GetTicketByCode(code string) (*domain.Ticket, error)
	UpdateTicket(code string, updates map[string]interface{}) (*domain.Ticket, error)
	DeleteTicket(code string) (*domain.Ticket, error)
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

// validateTicket enforces the business rule: ticket must be associated with either a packet OR an event (not both, not neither)
func (service *ticketService) validateTicket(ticket *domain.Ticket) error {
	if ticket == nil {
		return &domain.ValidationError{Reason: "invalid ticket object"}
	}

	// Business rule: must have exactly one of PacketID or EventID
	if ticket.PacketID == nil && ticket.EventID == nil {
		return &domain.ValidationError{Reason: "ticket must be associated with either a packet or an event"}
	}

	if ticket.PacketID != nil && ticket.EventID != nil {
		return &domain.ValidationError{Reason: "ticket cannot be associated with both a packet and an event"}
	}

	return nil
}

func (service *ticketService) CreateTicket(ticket *domain.Ticket) (*domain.Ticket, error) {
	if err := service.validateTicket(ticket); err != nil {
		return nil, err
	}

	// Generate UUID for ticket code
	ticket.Code = uuid.New().String()

	// Validate seat availability
	if err := service.validateSeatAvailability(ticket); err != nil {
		return nil, err
	}

	return service.repo.CreateTicket(ticket)
}

// validateSeatAvailability checks if there are enough seats available for the ticket
// Complex logic:
// - For EVENT tickets: available = event.Seats - (direct_tickets_sold + sum(packets_allocated_seats))
// - For PACKET tickets: available = packet.AllocatedSeats - packet_tickets_sold
func (service *ticketService) validateSeatAvailability(ticket *domain.Ticket) error {
	if ticket.EventID != nil {
		// Ticket is for a direct event (not through a packet)
		event, err := service.eventRepo.GetByID(*ticket.EventID)
		if err != nil {
			return err
		}

		// Check if event has seats defined
		if event.Seats == nil {
			return &domain.ValidationError{
				Reason: fmt.Sprintf("event %d does not have seats defined", *ticket.EventID),
			}
		}

		// Calculate available seats for direct event tickets
		// Formula: event.Seats - (direct_tickets_sold + sum(all_packets_allocated_seats_for_this_event))

		// 1. Count tickets already sold directly to this event
		directTicketsSold, err := service.eventRepo.CountSoldTickets(*ticket.EventID)
		if err != nil {
			return err
		}

		// 2. Get all packets that include this event
		packets, err := service.inclusionRepo.GetEventPacketsByEventID(*ticket.EventID)
		if err != nil {
			return err
		}

		// 3. Sum all allocated seats from packets that include this event
		totalPacketSeats := 0
		for _, packet := range packets {
			if packet.AllocatedSeats != nil {
				totalPacketSeats += *packet.AllocatedSeats
			}
		}

		// 4. Calculate available seats
		totalSeats := *event.Seats
		reservedSeats := directTicketsSold + totalPacketSeats
		availableSeats := totalSeats - reservedSeats

		// 5. Check if seats are available
		if availableSeats <= 0 {
			return &domain.ValidationError{
				Reason: fmt.Sprintf("event '%s' has no available seats (total: %d, direct tickets: %d, packet allocations: %d)",
					event.Name, totalSeats, directTicketsSold, totalPacketSeats),
			}
		}
	}

	if ticket.PacketID != nil {
		// Ticket is for a packet
		packet, err := service.packetRepo.GetByID(*ticket.PacketID)
		if err != nil {
			return err
		}

		// Check if packet has allocated seats defined
		if packet.AllocatedSeats == nil {
			return &domain.ValidationError{
				Reason: fmt.Sprintf("packet %d does not have allocated seats defined", *ticket.PacketID),
			}
		}

		// Calculate available seats for packet tickets
		// Formula: packet.AllocatedSeats - packet_tickets_sold

		// Count tickets already sold for this packet
		soldTickets, err := service.packetRepo.CountSoldTickets(*ticket.PacketID)
		if err != nil {
			return err
		}

		// Calculate available seats
		availableSeats := *packet.AllocatedSeats - soldTickets

		// Check if seats are available
		if availableSeats <= 0 {
			return &domain.ValidationError{
				Reason: fmt.Sprintf("packet '%s' is sold out (%d/%d tickets sold)",
					packet.Name, soldTickets, *packet.AllocatedSeats),
			}
		}
	}

	return nil
}

func (service *ticketService) GetTicketByCode(code string) (*domain.Ticket, error) {
	if code == "" {
		return nil, &domain.ValidationError{Reason: "ticket code is required"}
	}
	return service.repo.GetTicketByCode(code)
}

func (service *ticketService) UpdateTicket(code string, updates map[string]interface{}) (*domain.Ticket, error) {
	if code == "" {
		return nil, &domain.ValidationError{Reason: "ticket code is required"}
	}

	if len(updates) == 0 {
		return nil, &domain.ValidationError{Reason: "no fields to update"}
	}

	// Validate the constraint if updating packet_id or event_id
	packetID, hasPacketID := updates["packet_id"]
	eventID, hasEventID := updates["event_id"]

	if hasPacketID || hasEventID {
		// Get current ticket to check the constraint
		current, err := service.repo.GetTicketByCode(code)
		if err != nil {
			return nil, err
		}

		// Determine the new values
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

		// Validate the constraint
		if newPacketID == nil && newEventID == nil {
			return nil, &domain.ValidationError{Reason: "ticket must be associated with either a packet or an event"}
		}

		if newPacketID != nil && newEventID != nil {
			return nil, &domain.ValidationError{Reason: "ticket cannot be associated with both a packet and an event"}
		}

		// If the ticket is being moved to a different event or packet, validate seat availability
		isMovingToNewEvent := hasEventID && (current.EventID == nil || (newEventID != nil && *newEventID != *current.EventID))
		isMovingToNewPacket := hasPacketID && (current.PacketID == nil || (newPacketID != nil && *newPacketID != *current.PacketID))

		if isMovingToNewEvent || isMovingToNewPacket {
			// Create a temporary ticket object with new values to validate seat availability
			tempTicket := &domain.Ticket{
				Code:     code,
				PacketID: newPacketID,
				EventID:  newEventID,
			}

			if err := service.validateSeatAvailability(tempTicket); err != nil {
				return nil, err
			}
		}
	}

	return service.repo.UpdateTicket(code, updates)
}

func (service *ticketService) DeleteTicket(code string) (*domain.Ticket, error) {
	if code == "" {
		return nil, &domain.ValidationError{Reason: "ticket code is required"}
	}
	return service.repo.DeleteEvent(code)
}
