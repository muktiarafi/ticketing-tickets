package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	common "github.com/muktiarafi/ticketing-common"
	"github.com/muktiarafi/ticketing-tickets/internal/model"
	"github.com/muktiarafi/ticketing-tickets/internal/service"
)

type TicketHandler struct {
	service.TicketService
}

func NewTicketHandler(ticketService service.TicketService) *TicketHandler {
	return &TicketHandler{
		TicketService: ticketService,
	}
}

func (h *TicketHandler) Route(router *echo.Echo) {
	tickets := router.Group("/api/tickets")
	tickets.POST("", h.NewTicket, common.RequireAuth)
	tickets.GET("", h.GetAll)
	tickets.PUT("/:ticketID", h.Update, common.RequireAuth)
	tickets.GET("/:ticketID", h.Show)
}

func (h *TicketHandler) NewTicket(c echo.Context) error {
	userPayload, ok := c.Get("userPayload").(*common.UserPayload)
	const op = "TicketHandler.Update"
	if !ok {
		return &common.Error{
			Op:  op,
			Err: errors.New("missing payload in context"),
		}
	}

	ticketDTO := new(model.TicketDTO)
	if err := c.Bind(ticketDTO); err != nil {
		return &common.Error{Op: op, Err: err}
	}

	if err := c.Validate(ticketDTO); err != nil {
		return err
	}

	ticket, err := h.TicketService.Create(int64(userPayload.ID), ticketDTO)
	if err != nil {
		return err
	}

	return common.NewResponse(http.StatusCreated, "Created", ticket).SendJSON(c)
}

func (h *TicketHandler) GetAll(c echo.Context) error {
	tickets, err := h.Find()
	if err != nil {
		return err
	}

	return common.NewResponse(http.StatusOK, "OK", tickets).SendJSON(c)
}

func (h *TicketHandler) Show(c echo.Context) error {
	ticketIDStr := c.Param("ticketID")
	ticketID, err := strconv.ParseInt(ticketIDStr, 10, 64)
	if err != nil {
		return &common.Error{
			Code:    common.EINVALID,
			Op:      "TicketHandler.Show",
			Message: "Invalid Ticket Id",
			Err:     err,
		}
	}

	ticket, err := h.FindOne(ticketID)
	if err != nil {
		return err
	}

	return common.NewResponse(http.StatusOK, "OK", ticket).SendJSON(c)
}

func (h *TicketHandler) Update(c echo.Context) error {
	userPayload, ok := c.Get("userPayload").(*common.UserPayload)
	if !ok {
		return &common.Error{
			Op:  "TicketHandler.Update",
			Err: errors.New("missing payload in context"),
		}
	}

	ticketIDStr := c.Param("ticketID")
	ticketID, err := strconv.ParseInt(ticketIDStr, 10, 64)
	const op = "TicketHandler.Update"
	if err != nil {
		return &common.Error{
			Code:    common.EINVALID,
			Op:      op,
			Message: "Invalid Ticket Id",
			Err:     err,
		}
	}

	ticketDTO := new(model.TicketDTO)
	if err := c.Bind(ticketDTO); err != nil {
		return &common.Error{Op: op, Err: err}
	}

	if err := c.Validate(ticketDTO); err != nil {
		return err
	}

	updatedTicket, err := h.TicketService.Update(int64(userPayload.ID), ticketID, ticketDTO)
	if err != nil {
		return err
	}

	return common.NewResponse(http.StatusOK, "OK", updatedTicket).SendJSON(c)
}
