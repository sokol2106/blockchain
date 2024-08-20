package main

import (
	"github.com/ivan/blockchain/block-client/internal/app"
	"os"
)

const CServerURL = "http://localhost:8080"
const CNoncePattern = "0000"

type params struct {
	ServerURL    string
	NoncePattern string
}

func main() {
	p := params{
		ServerURL:    os.Getenv("SERVER_URL"),
		NoncePattern: os.Getenv("NONCE_PATTERN"),
	}

	if p.ServerURL == "" {
		p.ServerURL = CServerURL
	}

	if p.NoncePattern == "" {
		p.NoncePattern = CNoncePattern
	}

	ParseFlags(&p)
	app.Run(p.ServerURL, p.NoncePattern)
}
