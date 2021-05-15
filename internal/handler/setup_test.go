package handler

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-playground/validator/v10"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	common "github.com/muktiarafi/ticketing-common"
	"github.com/muktiarafi/ticketing-tickets/internal/driver"
	"github.com/muktiarafi/ticketing-tickets/internal/entity"
	"github.com/muktiarafi/ticketing-tickets/internal/repository"
	"github.com/muktiarafi/ticketing-tickets/internal/service"
	"github.com/ory/dockertest/v3"
)

var (
	pool     *dockertest.Pool
	resource *dockertest.Resource
)

var router *echo.Echo

func TestMain(m *testing.M) {
	db := &driver.DB{
		SQL: newTestDatabase(),
	}

	router = echo.New()
	router.Use(middleware.Logger())

	val := validator.New()
	trans := common.NewDefaultTranslator(val)
	customValidator := &common.CustomValidator{val, trans}
	router.Validator = customValidator
	router.HTTPErrorHandler = common.CustomErrorHandler

	ticketRepository := repository.NewTicketRepository(db)

	ticketPublisher := &ticketPublisherStub{}
	ticketService := service.NewTicketService(ticketPublisher, ticketRepository)
	ticketHandler := NewTicketHandler(ticketService)
	ticketHandler.Route(router)

	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func newTestDatabase() *sql.DB {
	var err error
	pool, err = dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err = pool.Run("postgres", "alpine", []string{"POSTGRES_PASSWORD=secret", "POSTGRES_DB=postgres"})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	var db *sql.DB
	if err = pool.Retry(func() error {
		db, err = sql.Open(
			"pgx",
			fmt.Sprintf("host=localhost port=%s dbname=postgres user=postgres password=secret", resource.GetPort("5432/tcp")))
		if err != nil {
			return err
		}

		migrationFilePath := filepath.Join("..", "..", "db", "migrations")
		return driver.Migration(migrationFilePath, db)
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	return db
}

type ticketPublisherStub struct{}

func (p *ticketPublisherStub) Created(ticket *entity.Ticket) error {
	fmt.Println("Ticket publisher publish ticket created event")

	return nil
}

func (p *ticketPublisherStub) Updated(ticket *entity.Ticket) error {
	fmt.Println("Ticket publisher publish ticket created event")

	return nil
}
