package app

import (
	"github.com/ivan/blockchain/api-server/internal/handlers"
	"github.com/ivan/blockchain/api-server/internal/server"
	"github.com/ivan/blockchain/api-server/internal/service"
	"github.com/ivan/blockchain/api-server/internal/storage"
	"log"
)

func Run(addServer string, cnfDataBase string) {
	stor := storage.NewPostgresql(cnfDataBase)
	err := stor.Connect()
	if err != nil {
		log.Printf("Connect DB: %s", err)
	} else {
		err = stor.Migrations("file://migrations/postgresql")
		if err != nil {
			log.Printf("Could not run migrations: %s", err)
		}
	}
	srvBlockchain := service.NewBlockchain(stor)
	srvBlockchain.RunProcessBlockChain()
	srvBlockchain.RunBlockchainDBLoad()
	srvVerify := service.NewVerification()

	ser := server.NewServer(handlers.Router(handlers.NewHandlers(srvBlockchain, srvVerify)), addServer)
	err = ser.Start()
	if err != nil {
		log.Printf("Starting server error: %s", err)
	}
}
