package producer

import "github.com/muktiarafi/ticketing-tickets/internal/entity"

type TicketProducer interface {
	Created(ticket *entity.Ticket) error
	Updated(ticket *entity.Ticket) error
}
