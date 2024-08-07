package main

import (
	"fmt"
	"github.com/ivan/blockchain/api-server/internal/app"
)

const RUN_ADDRESS = "localhost:8080"
const DATABASE_URI = "localhost:5432"

type params struct {
	ServerAddress   string
	BaseAddress     string
	FileStoragePath string
	DatabaseDSN     string
}

func main() {

	fmt.Println("Hello World")
	app.Run("localhost:8080", "localhost:5432")
}
