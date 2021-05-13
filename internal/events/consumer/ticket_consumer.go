package consumer

import (
	"log"

	"github.com/ThreeDotsLabs/watermill/message"
	common "github.com/muktiarafi/ticketing-common"
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
	log.Println("received event from topic:", common.OrderCreated)
	orderCreatedData := new(types.OrderCreatedEvent)
	if err := orderCreatedData.Unmarshal(msg.Payload); err != nil {
		msg.Nack()
		return err
	}

	ticket, err := c.FindOne(orderCreatedData.TicketID)
	if err != nil {
		msg.Ack()
		return err
	}

	ticket.OrderID = orderCreatedData.ID
	ticket.Version++
	updatedTicket, err := c.Update(ticket)
	if err != nil {
		er, _ := err.(*common.Error)
		if er.Code == common.ENOTFOUND {
			msg.Ack()
		} else {
			msg.Nack()
		}
		return err
	}

	if err := c.Updated(updatedTicket); err != nil {
		msg.Nack()
		return err
	}

	msg.Ack()

	return nil
}

func (c *TicketConsumer) OrderCancelled(msg *message.Message) error {
	log.Println("received event from topic:", common.OrderCancelled)
	orderCreatedData := new(types.OrderCreatedEvent)
	if err := orderCreatedData.Unmarshal(msg.Payload); err != nil {
		msg.Nack()
		return err
	}

	ticket, err := c.FindOne(orderCreatedData.TicketID)
	if err != nil {
		er, _ := err.(*common.Error)
		if er.Code == common.ENOTFOUND {
			msg.Ack()
		} else {
			msg.Nack()
		}
		return err
	}

	ticket.OrderID = 0
	updatedTicket, err := c.Update(ticket)
	if err != nil {
		msg.Nack()
		return err
	}

	if err := c.Created(updatedTicket); err != nil {
		msg.Nack()
		return err
	}

	msg.Ack()

	return nil
}
