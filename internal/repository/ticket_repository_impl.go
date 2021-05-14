package repository

import (
	"database/sql"

	common "github.com/muktiarafi/ticketing-common"
	"github.com/muktiarafi/ticketing-tickets/internal/driver"
	"github.com/muktiarafi/ticketing-tickets/internal/entity"
	"github.com/muktiarafi/ticketing-tickets/internal/model"
)

type TicketRepositoryImpl struct {
	*driver.DB
}

func NewTicketRepository(db *driver.DB) TicketRepostiory {
	return &TicketRepositoryImpl{
		DB: db,
	}
}

func (r *TicketRepositoryImpl) Insert(userID int64, ticketDTO *model.TicketDTO) (*entity.Ticket, error) {
	ctx, cancel := newDBContext()
	defer cancel()

	stmt := `INSERT INTO tickets (title, price, user_id)
	VALUES ($1, $2, $3)
	RETURNING *`

	var ticket entity.Ticket
	if err := r.SQL.QueryRowContext(ctx, stmt, ticketDTO.Title, ticketDTO.Price, userID).Scan(
		&ticket.ID,
		&ticket.Title,
		&ticket.Price,
		&ticket.UserID,
		&ticket.OrderID,
		&ticket.Version,
	); err != nil {
		return nil, &common.Error{Op: "TicketRepositoryImpl.Insert", Err: err}
	}

	return &ticket, nil
}

func (r *TicketRepositoryImpl) Find() ([]*entity.Ticket, error) {
	ctx, cancel := newDBContext()
	defer cancel()

	stmt := `SELECT * FROM tickets`

	tickets := make([]*entity.Ticket, 0)
	const op = "TicketRepository.Find"
	rows, err := r.SQL.QueryContext(ctx, stmt)
	if err != nil {
		return tickets, &common.Error{Op: op, Err: err}
	}
	defer rows.Close()

	for rows.Next() {
		ticket := new(entity.Ticket)
		if err := rows.Scan(
			&ticket.ID,
			&ticket.Title,
			&ticket.Price,
			&ticket.UserID,
			&ticket.OrderID,
			&ticket.Version,
		); err != nil {
			return tickets, &common.Error{Op: op, Err: err}
		}

		tickets = append(tickets, ticket)
	}

	return tickets, nil
}

func (r *TicketRepositoryImpl) FindOne(ticketID int64) (*entity.Ticket, error) {
	ctx, cancel := newDBContext()
	defer cancel()

	stmt := `SELECT * FROM tickets
	WHERE id = $1`

	ticket := new(entity.Ticket)
	if err := r.SQL.QueryRowContext(ctx, stmt, ticketID).Scan(
		&ticket.ID,
		&ticket.Title,
		&ticket.Price,
		&ticket.UserID,
		&ticket.OrderID,
		&ticket.Version,
	); err != nil {
		const op = "TicketRepositoryImpl.FindOne"
		if err == sql.ErrNoRows {
			return nil, &common.Error{
				Code:    common.ENOTFOUND,
				Op:      op,
				Message: "Ticket Not Found",
				Err:     err,
			}
		}
		return nil, &common.Error{Op: op, Err: err}
	}

	return ticket, nil
}

func (r *TicketRepositoryImpl) Update(ticket *entity.Ticket) (*entity.Ticket, error) {
	ctx, cancel := newDBContext()
	defer cancel()

	stmt := `UPDATE tickets 
	SET title=$1, price=$2, version=$3
	WHERE id = $4 AND version =$5
	RETURNING *`

	updatedTicket := new(entity.Ticket)
	if err := r.SQL.QueryRowContext(
		ctx,
		stmt,
		ticket.Title,
		ticket.Price,
		ticket.Version,
		ticket.ID,
		ticket.Version-1,
	).Scan(
		&updatedTicket.ID,
		&updatedTicket.Title,
		&updatedTicket.Price,
		&updatedTicket.UserID,
		&updatedTicket.OrderID,
		&updatedTicket.Version,
	); err != nil {
		const op = "TicketRepositoryImpl.Update"
		if err == sql.ErrNoRows {
			return nil, &common.Error{
				Code:    common.ENOTFOUND,
				Op:      op,
				Message: "Ticket Not Found",
				Err:     err,
			}
		}
		return nil, &common.Error{Op: op, Err: err}
	}

	return updatedTicket, nil
}
