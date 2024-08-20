package app

import (
	"fmt"
	"github.com/ivan/blockchain/block-client/internal/service"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func Run(serverURL, noncePattern string) {

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	srvMiner := service.NewBlockMiner(serverURL, noncePattern)
	srvVerify := service.NewBlockVerification(serverURL, noncePattern)

	for {
		select {
		case sig := <-signalChan:
			fmt.Printf("Received signal: %s. Exiting...\n", sig)
			return
		default:
			err := srvMiner.RequestMiningData()
			if err == nil {
				srvMiner.MineData()
				err = srvMiner.SendMiningBlock()
				if err != nil {
					log.Printf("send block: %s", err)

				}
			}

			err = srvVerify.RequestVerificationData()
			if err == nil {
				srvVerify.VerifyData()
				err = srvVerify.UpdateStatus()
				if err != nil {
					log.Printf("verify data: %s", err)
				}
			}
		}
	}

}
