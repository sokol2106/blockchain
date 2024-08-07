package app

import (
	"github.com/ivan/blockchain/api-server/internal/handlers"
	"github.com/ivan/blockchain/api-server/internal/server"
	"github.com/ivan/blockchain/api-server/internal/service"
	"log"
)

func Run(addServer string, addDB string) {
	srvBlockchain := service.NewBlockchain()
	srvVerify := service.NewVerification()

	ser := server.NewServer(handlers.Router(handlers.NewHandlers(srvBlockchain, srvVerify)), addServer)
	err := ser.Start()
	if err != nil {
		log.Printf("Starting server error: %s", err)
	}
}
