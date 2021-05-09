package consumer

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/muktiarafi/ticketing-common/types"
	"github.com/muktiarafi/ticketing-tickets/internal/events/producer"
	"github.com/muktiarafi/ticketing-tickets/internal/repository"
)

type TicketConsumer struct {
	producer.TicketProducer
	repository.TicketRepostiory
}

func NewTicketConsumer(producer producer.TicketProducer, ticketRepo repository.TicketRepostiory) *TicketConsumer {
	return &TicketConsumer{
		TicketProducer:   producer,
		TicketRepostiory: ticketRepo,
	}
}

func (c *TicketConsumer) OrderCreated(msg *message.Message) error {
	orderCreatedData := new(types.OrderCreatedEvent)
	if err := orderCreatedData.Unmarshal(msg.Payload); err != nil {
		return err
	}

	ticket, err := c.FindOne(orderCreatedData.TicketID)
	if err != nil {
		return err
	}

	ticket.OrderID = orderCreatedData.ID
	updatedTicket, err := c.Update(ticket)
	if err != nil {
		return err
	}

	if err := c.Created(updatedTicket); err != nil {
		return err
	}

	msg.Ack()

	return nil
}

func (c *TicketConsumer) OrderCancelled(msg *message.Message) error {
	orderCreatedData := new(types.OrderCreatedEvent)
	if err := orderCreatedData.Unmarshal(msg.Payload); err != nil {
		return err
	}

	ticket, err := c.FindOne(orderCreatedData.TicketID)
	if err != nil {
		return err
	}

	ticket.OrderID = 0
	updatedTicket, err := c.Update(ticket)
	if err != nil {
		return err
	}

	if err := c.Created(updatedTicket); err != nil {
		return err
	}

	msg.Ack()

	return nil
}
