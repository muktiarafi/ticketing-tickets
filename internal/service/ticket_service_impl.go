package service

import (
	"errors"

	common "github.com/muktiarafi/ticketing-common"
	"github.com/muktiarafi/ticketing-tickets/internal/entity"
	"github.com/muktiarafi/ticketing-tickets/internal/model"
	"github.com/muktiarafi/ticketing-tickets/internal/repository"
)

type TicketServiceImpl struct {
	repository.TicketRepostiory
}

func NewTicketService(ticketRepo repository.TicketRepostiory) TicketService {
	return &TicketServiceImpl{
		TicketRepostiory: ticketRepo,
	}
}

func (s *TicketServiceImpl) Create(userID int64, ticketDTO *model.TicketDTO) (*entity.Ticket, error) {
	return s.Insert(userID, ticketDTO)
}

func (s *TicketServiceImpl) Find() ([]*entity.Ticket, error) {
	return s.TicketRepostiory.Find()
}

func (s *TicketServiceImpl) FindOne(ticketID int64) (*entity.Ticket, error) {
	return s.TicketRepostiory.FindOne(ticketID)
}

func (s *TicketServiceImpl) Update(userID int64, ticketID int64, ticketDTO *model.TicketDTO) (*entity.Ticket, error) {
	ticket, err := s.TicketRepostiory.FindOne(ticketID)
	if err != nil {
		return nil, err
	}

	if userID != ticket.UserID {
		return nil, &common.Error{
			Op:      "TicketServiceImpl.Update",
			Message: "Not Authorized",
			Err:     errors.New("Trying to update not owned ticket"),
		}
	}

	ticket.Title = ticketDTO.Title
	ticket.Price = ticketDTO.Price
	ticket.Version++

	updatedTicket, err := s.TicketRepostiory.Update(ticket)
	if err != nil {
		return nil, err
	}

	return updatedTicket, nil
}
