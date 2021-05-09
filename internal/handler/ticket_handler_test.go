package handler

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	common "github.com/muktiarafi/ticketing-common"
	"github.com/muktiarafi/ticketing-tickets/internal/entity"
	"github.com/muktiarafi/ticketing-tickets/internal/model"
)

func TestTicketHandlerNew(t *testing.T) {
	userPayload := &common.UserPayload{1, "bambank@gmail.com"}
	cookie := signIn(userPayload)

	t.Run("create ticket with passing cookie", func(t *testing.T) {
		ticketDTO := &model.TicketDTO{"konser", 14000}
		ticketJSON, _ := json.Marshal(ticketDTO)

		request := httptest.NewRequest(http.MethodPost, "/api/tickets", bytes.NewBuffer(ticketJSON))
		request.Header.Set("Content-Type", "application/json")
		request.AddCookie(cookie)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)
		assertResponseCode(t, http.StatusCreated, response.Code)

		responseBody, _ := ioutil.ReadAll(response.Body)
		apiResponse := struct {
			Data *entity.Ticket `json:"data"`
		}{}
		json.Unmarshal(responseBody, &apiResponse)

		got := apiResponse.Data.Title
		want := ticketDTO.Title
		if got != want {
			t.Errorf("Expecting title to be %q, but got %q instead", got, want)
		}

		if apiResponse.Data.Version != 1 {
			t.Error("Expecting version to be 1")
		}
	})

	t.Run("create ticket without passing cookie", func(t *testing.T) {
		ticketDTO := &model.TicketDTO{"konser", 14000}
		ticketJSON, _ := json.Marshal(ticketDTO)

		request := httptest.NewRequest(http.MethodPost, "/api/tickets", bytes.NewBuffer(ticketJSON))
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)
		assertResponseCode(t, http.StatusBadRequest, response.Code)
	})

	t.Run("create ticket with invalid title", func(t *testing.T) {
		ticketDTO := &model.TicketDTO{"1", 14000}
		ticketJSON, _ := json.Marshal(ticketDTO)

		request := httptest.NewRequest(http.MethodPost, "/api/tickets", bytes.NewBuffer(ticketJSON))
		request.Header.Set("Content-Type", "application/json")
		request.AddCookie(cookie)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)
		assertResponseCode(t, http.StatusBadRequest, response.Code)
	})

	t.Run("create ticket with minus price", func(t *testing.T) {
		ticketDTO := &model.TicketDTO{"konser", -1}
		ticketJSON, _ := json.Marshal(ticketDTO)

		request := httptest.NewRequest(http.MethodPost, "/api/tickets", bytes.NewBuffer(ticketJSON))
		request.Header.Set("Content-Type", "application/json")
		request.AddCookie(cookie)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)
		assertResponseCode(t, http.StatusBadRequest, response.Code)
	})

	t.Run("create ticket with invalid data", func(t *testing.T) {
		ticketDTO := &model.TicketDTO{"1", -1}
		ticketJSON, _ := json.Marshal(ticketDTO)

		request := httptest.NewRequest(http.MethodPost, "/api/tickets", bytes.NewBuffer(ticketJSON))
		request.Header.Set("Content-Type", "application/json")
		request.AddCookie(cookie)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)
		assertResponseCode(t, http.StatusBadRequest, response.Code)
	})
}

func TestTicketHandlerGetAll(t *testing.T) {
	userPayload := &common.UserPayload{2, "paijo@gmail.com"}
	cookie := signIn(userPayload)

	ticketDTO := &model.TicketDTO{"konser", 14000}
	ticketJSON, _ := json.Marshal(ticketDTO)

	request := httptest.NewRequest(http.MethodPost, "/api/tickets", bytes.NewBuffer(ticketJSON))
	request.Header.Set("Content-Type", "application/json")
	request.AddCookie(cookie)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assertResponseCode(t, http.StatusCreated, response.Code)

	request = httptest.NewRequest(http.MethodGet, "/api/tickets", nil)
	response = httptest.NewRecorder()
	router.ServeHTTP(response, request)
	assertResponseCode(t, http.StatusOK, response.Code)

	responseBody, _ := ioutil.ReadAll(response.Body)
	apiResponse := struct {
		Data []*entity.Ticket `json:"data"`
	}{}
	json.Unmarshal(responseBody, &apiResponse)

	if len(apiResponse.Data) == 0 {
		t.Error("Expecting to get list of tickets")
	}
}

func TestTicketHandlerShow(t *testing.T) {
	userPayload := &common.UserPayload{3, "yoloo@gmail.com"}
	cookie := signIn(userPayload)

	t.Run("show created ticket", func(t *testing.T) {
		ticketDTO := &model.TicketDTO{"konser", 14000}
		ticketJSON, _ := json.Marshal(ticketDTO)

		request := httptest.NewRequest(http.MethodPost, "/api/tickets", bytes.NewBuffer(ticketJSON))
		request.Header.Set("Content-Type", "application/json")
		request.AddCookie(cookie)
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)

		assertResponseCode(t, http.StatusCreated, response.Code)

		responseBody, _ := ioutil.ReadAll(response.Body)
		apiResponse := struct {
			Data *entity.Ticket `json:"data"`
		}{}
		json.Unmarshal(responseBody, &apiResponse)

		ticket := apiResponse.Data
		if ticket.Title != ticketDTO.Title {
			t.Errorf("Expecting title to be %q, but got %q instead", ticketDTO.Title, ticket.Title)
		}

		ticketID := strconv.FormatInt(apiResponse.Data.ID, 10)
		request = httptest.NewRequest(http.MethodGet, "/api/tickets/"+ticketID, nil)
		response = httptest.NewRecorder()
		router.ServeHTTP(response, request)
		assertResponseCode(t, http.StatusOK, response.Code)

		responseBody, _ = ioutil.ReadAll(response.Body)
		apiResponse = struct {
			Data *entity.Ticket `json:"data"`
		}{}
		json.Unmarshal(responseBody, &apiResponse)

		if apiResponse.Data.ID != ticket.ID {
			t.Errorf("Expecting id to be %d, but got %d instead", ticket.ID, apiResponse.Data.ID)
		}
	})

	t.Run("requesting nonexistent ticket", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/api/tickets/9999", nil)
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)

		assertResponseCode(t, http.StatusNotFound, response.Code)
	})

	t.Run("invalid path param", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/api/tickets/werwer", nil)
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)

		assertResponseCode(t, http.StatusBadRequest, response.Code)
	})
}

func TestTicketHandlerUpdate(t *testing.T) {
	userPayload := &common.UserPayload{4, "werwer@gmail.com"}
	cookie := signIn(userPayload)

	t.Run("update already created ticket", func(t *testing.T) {
		ticketDTO := &model.TicketDTO{"konser", 14000}
		ticketJSON, _ := json.Marshal(ticketDTO)

		request := httptest.NewRequest(http.MethodPost, "/api/tickets", bytes.NewBuffer(ticketJSON))
		request.Header.Set("Content-Type", "application/json")
		request.AddCookie(cookie)
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)

		assertResponseCode(t, http.StatusCreated, response.Code)
		responseBody, _ := ioutil.ReadAll(response.Body)
		apiResponse := struct {
			Data *entity.Ticket `json:"data"`
		}{}
		json.Unmarshal(responseBody, &apiResponse)
		ticket := apiResponse.Data

		ticketDTO.Price = 18000
		ticketJSON, _ = json.Marshal(ticketDTO)

		ticketID := strconv.FormatInt(ticket.ID, 10)
		request = httptest.NewRequest(http.MethodPut, "/api/tickets/"+ticketID, bytes.NewBuffer(ticketJSON))
		request.Header.Set("Content-Type", "application/json")
		request.AddCookie(cookie)
		response = httptest.NewRecorder()
		router.ServeHTTP(response, request)

		assertResponseCode(t, http.StatusOK, response.Code)

		responseBody, _ = ioutil.ReadAll(response.Body)
		apiResponse = struct {
			Data *entity.Ticket `json:"data"`
		}{}
		json.Unmarshal(responseBody, &apiResponse)

		if apiResponse.Data.ID != ticket.ID {
			t.Errorf("Expecting id to be %d, but got %d instead", ticket.ID, apiResponse.Data.ID)
		}

		if apiResponse.Data.Version != ticket.Version+1 {
			t.Errorf("Expecting version to be %d, but got %d instead", ticket.Version+1, apiResponse.Data.Version)
		}
	})

	t.Run("update nonexistent ticket", func(t *testing.T) {
		ticketDTO := &model.TicketDTO{"konser", 14000}
		ticketJSON, _ := json.Marshal(ticketDTO)

		request := httptest.NewRequest(http.MethodPut, "/api/tickets/9999", bytes.NewBuffer(ticketJSON))
		request.Header.Set("Content-Type", "application/json")
		request.AddCookie(cookie)
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)

		assertResponseCode(t, http.StatusNotFound, response.Code)
	})

	// t.Run("update with not sync version", func(t *testing.T) {
	// 	ticketDTO := &model.TicketDTO{"konser", 14000}
	// 	ticketJSON, _ := json.Marshal(ticketDTO)

	// 	request := httptest.NewRequest(http.MethodPost, "/api/tickets", bytes.NewBuffer(ticketJSON))
	// 	request.Header.Set("Content-Type", "application/json")
	// 	request.AddCookie(cookie)
	// 	response := httptest.NewRecorder()
	// 	router.ServeHTTP(response, request)

	// 	assertResponseCode(t, http.StatusCreated, response.Code)
	// 	responseBody, _ := ioutil.ReadAll(response.Body)
	// 	apiResponse := struct {
	// 		Data *entity.Ticket `json:"data"`
	// 	}{}
	// 	json.Unmarshal(responseBody, &apiResponse)
	// 	ticket := apiResponse.Data

	// 	ticketDTO.Title = "2"
	// 	ticketDTO.Price = -1
	// 	ticketJSON, _ = json.Marshal(ticketDTO)

	// 	ticketID := strconv.Itoa(ticket.ID)
	// 	request = httptest.NewRequest(http.MethodPut, "/api/tickets/"+ticketID, bytes.NewBuffer(ticketJSON))
	// 	request.Header.Set("Content-Type", "application/json")
	// 	request.AddCookie(cookie)
	// 	response = httptest.NewRecorder()
	// 	router.ServeHTTP(response, request)

	// 	assertResponseCode(t, http.StatusBadRequest, response.Code)
	// })
}

func assertResponseCode(t testing.TB, want, got int) {
	t.Helper()

	if got != want {
		t.Errorf("Expected status code %d, but got %d instead", want, got)
	}
}

func signIn(userPayload *common.UserPayload) *http.Cookie {

	token, _ := common.CreateToken(userPayload)

	cookie := http.Cookie{
		Name:    "session",
		Value:   token,
		Expires: time.Now().Add(10 * time.Minute),
		Path:    "/auth",
	}

	return &cookie
}
