package main

import (
	"log"

	"github.com/muktiarafi/ticketing-tickets/internal/server"
)

func main() {
	e := server.SetupServer()

	log.Fatal(e.Start(":8080"))
}
