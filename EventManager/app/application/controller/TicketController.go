package controller

import (
	"eventManager/application/domain"
	"eventManager/application/service"
)

type ITicketController interface {
	CreateTicket(ticket *domain.Ticket) (*domain.Ticket, error)
	GetTicketByCode(code string) (*domain.Ticket, error)
	UpdateTicket(code string, updates map[string]interface{}) (*domain.Ticket, error)
	DeleteTicket(code string) (*domain.Ticket, error)
}

type ticketController struct {
	service service.ITicketService
}

func NewTicketController(service service.ITicketService) *ticketController {
	return &ticketController{service: service}
}

func (c *ticketController) CreateTicket(ticket *domain.Ticket) (*domain.Ticket, error) {
	return c.service.CreateTicket(ticket)
}

func (c *ticketController) GetTicketByCode(code string) (*domain.Ticket, error) {
	return c.service.GetTicketByCode(code)
}

func (c *ticketController) UpdateTicket(code string, updates map[string]interface{}) (*domain.Ticket, error) {
	return c.service.UpdateTicket(code, updates)
}

func (c *ticketController) DeleteTicket(code string) (*domain.Ticket, error) {
	return c.service.DeleteTicket(code)
}
