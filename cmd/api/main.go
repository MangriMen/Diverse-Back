package main

import (
	"github.com/MangriMen/Diverse-Back/api/server"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	server.InitApi()
}
