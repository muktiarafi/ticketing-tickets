package repository

import (
	"github.com/muktiarafi/ticketing-tickets/internal/entity"
	"github.com/muktiarafi/ticketing-tickets/internal/model"
)

type TicketRepostiory interface {
	Insert(userID int64, ticketDTO *model.TicketDTO) (*entity.Ticket, error)
	Find() ([]*entity.Ticket, error)
	FindOne(ticketId int64) (*entity.Ticket, error)
	Update(ticket *entity.Ticket) (*entity.Ticket, error)
}
