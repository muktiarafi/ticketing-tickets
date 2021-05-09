package producer

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	common "github.com/muktiarafi/ticketing-common"
	"github.com/muktiarafi/ticketing-common/types"
	"github.com/muktiarafi/ticketing-tickets/internal/entity"
)

type TicketPublisherImpl struct {
	message.Publisher
}

func NewTicketProducer(publisher message.Publisher) TicketProducer {
	return &TicketPublisherImpl{
		Publisher: publisher,
	}
}

func (p *TicketPublisherImpl) Created(ticket *entity.Ticket) error {
	ticketCreatedEvent := types.TicketCreatedEvent{
		ID:      ticket.ID,
		Version: ticket.Version,
		Title:   ticket.Title,
		Price:   ticket.Price,
		UserID:  ticket.UserID,
	}
	ticketBytes, err := ticketCreatedEvent.Marshal()
	if err != nil {
		return err
	}

	msg := message.NewMessage(watermill.NewUUID(), ticketBytes)
	return p.Publish(common.TicketCreated, msg)
}

func (p *TicketPublisherImpl) Updated(ticket *entity.Ticket) error {
	ticketUpdatedEvent := types.TicketUpdatedEvent{
		ID:      ticket.ID,
		Version: ticket.Version,
		Title:   ticket.Title,
		Price:   ticket.Price,
		UserID:  ticket.UserID,
		OrderID: ticket.OrderID,
	}
	ticketBytes, err := ticketUpdatedEvent.Marshal()
	if err != nil {
		return err
	}

	msg := message.NewMessage(watermill.NewUUID(), ticketBytes)
	return p.Publish(common.TIcketUpdated, msg)
}
