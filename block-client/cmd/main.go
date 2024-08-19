package main

import (
	"github.com/ivan/blockchain/block-client/internal/app"
	"os"
)

const CServerAddress = "http://localhost:8080"
const CNoncePattern = "0000"

type params struct {
	ServerAddress string
	NoncePattern  string
}

func main() {
	p := params{
		ServerAddress: os.Getenv("SERVER_ADDRESS"),
		NoncePattern:  os.Getenv("NONCE_PATTERN"),
	}

	if p.ServerAddress == "" {
		p.ServerAddress = CServerAddress
	}

	if p.NoncePattern == "" {
		p.NoncePattern = CNoncePattern
	}

	ParseFlags(&p)
	app.Run("http://localhost:8080", p.NoncePattern)
}
