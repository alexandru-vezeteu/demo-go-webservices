package service

import (
	"eventManager/application/domain"
	"eventManager/application/repository"
)

type ticketService struct {
	repo repository.TicketRepository
}

func NewTicketService(repo repository.TicketRepository) *ticketService {
	return &ticketService{repo: repo}
}

func (service *ticketService) CreateTicket(ticket *domain.Ticket) (*domain.Ticket, error) {
	return nil, nil
}

func (service *ticketService) GetTicketByCode(code string) (*domain.Ticket, error) {
	return nil, nil
}
func (service *ticketService) UpdateTicket(code string, updates map[string]interface{}) (*domain.Ticket, error) {
	return nil, nil
}
func (service *ticketService) DeleteEvent(code string) (*domain.Ticket, error) {
	return nil, nil
}
