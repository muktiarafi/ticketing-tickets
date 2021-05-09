package service

import (
	"github.com/muktiarafi/ticketing-tickets/internal/entity"
	"github.com/muktiarafi/ticketing-tickets/internal/model"
)

type TicketService interface {
	Create(int, *model.TicketDTO) (*entity.Ticket, error)
	Find() ([]*entity.Ticket, error)
	FindOne(ticketID int) (*entity.Ticket, error)
	Update(userID, ticketID int, ticketDTO *model.TicketDTO) (*entity.Ticket, error)
}
