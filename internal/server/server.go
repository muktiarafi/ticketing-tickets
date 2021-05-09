package server

import (
	"log"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	common "github.com/muktiarafi/ticketing-common"
	"github.com/muktiarafi/ticketing-tickets/internal/config"
	"github.com/muktiarafi/ticketing-tickets/internal/driver"
	"github.com/muktiarafi/ticketing-tickets/internal/events/consumer"
	"github.com/muktiarafi/ticketing-tickets/internal/events/producer"
	"github.com/muktiarafi/ticketing-tickets/internal/handler"
	custommiddleware "github.com/muktiarafi/ticketing-tickets/internal/middleware"
	"github.com/muktiarafi/ticketing-tickets/internal/repository"
	"github.com/muktiarafi/ticketing-tickets/internal/service"
)

func SetupServer() *echo.Echo {
	e := echo.New()
	p := custommiddleware.NewPrometheus("echo", nil)
	p.Use(e)

	val := validator.New()
	trans := common.NewDefaultTranslator(val)
	customValidator := &common.CustomValidator{val, trans}
	e.Validator = customValidator
	e.HTTPErrorHandler = common.CustomErrorHandler
	e.Use(middleware.Logger())

	db, err := driver.ConnectSQL(config.PostgresDSN())
	if err != nil {
		log.Fatal(err)
	}
	ticketRepository := repository.NewTicketRepository(db)
	ticketService := service.NewTicketService(ticketRepository)

	producerBrokers := []string{config.NewProducerBroker()}
	commonPublisher, err := common.CreatePublisher(producerBrokers, watermill.NewStdLogger(false, false))
	if err != nil {
		log.Fatal(err)
	}
	ticketPublisher := producer.NewTicketProducer(commonPublisher)
	ticketHandler := handler.NewTicketHandler(ticketPublisher, ticketService)
	ticketHandler.Route(e)

	consumerBrokers := []string{config.NewConsumerBroker()}
	subscriber, err := common.CreateSubscriber(consumerBrokers, "tickets-service", watermill.NewStdLogger(false, false))
	if err != nil {
		log.Fatal(err)
	}

	ticketConsumer := consumer.NewTicketConsumer(ticketPublisher, ticketRepository)
	consumer := common.NewConsumer(subscriber)
	consumer.On(common.OrderCreated, ticketConsumer.OrderCreated)
	consumer.On(common.OrderCancelled, ticketConsumer.OrderCancelled)

	return e
}
