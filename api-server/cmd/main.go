package main

import (
	"github.com/ivan/blockchain/api-server/internal/app"
	"os"
)

const CRunAddress = "localhost:8080"
const CDataBaseDSN = "host=localhost port=5432 user=pia password=12345678 dbname=yandex sslmode=disable"

type params struct {
	RunAddress  string
	DatabaseDSN string
}

func main() {

	p := params{
		RunAddress:  os.Getenv("RUN_ADDRESS"),
		DatabaseDSN: os.Getenv("DATABASE_DSN"),
	}

	if p.RunAddress == "" {
		p.RunAddress = CRunAddress
	}

	if p.DatabaseDSN == "" {
		p.DatabaseDSN = CDataBaseDSN
	}

	ParseFlags(&p)

	app.Run(p.RunAddress, p.DatabaseDSN)
}
